// Package tzf is a package convert (lng,lat) to timezone.
//
// Inspired by timezonefinder https://github.com/jannikmi/timezonefinder,
// fast python package for finding the timezone of any point on earth (coordinates) offline.
package tzf

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sync"

	"github.com/Alfex4936/tzf/convert"
	"github.com/Alfex4936/tzf/pb"
	"github.com/Alfex4936/tzf/reduce"
	"github.com/tidwall/geojson/geometry"
	"github.com/tidwall/rtree"
	"golang.org/x/exp/slices"
)

var ErrNoTimezoneFound = errors.New("tzf: no timezone found")

type Option struct {
	DropPBTZ bool
}

type OptionFunc = func(opt *Option)

// SetDropPBTZ will make Finder not save [github.com/Alfex4936/tzf/pb.Timezone] in memory
func SetDropPBTZ(opt *Option) {
	opt.DropPBTZ = true
}

type tzitem struct {
	pbtz  *pb.Timezone
	name  string
	polys []*geometry.Poly
	min   [2]float64
	max   [2]float64
}

func newNotFoundErr(lng float64, lat float64) error {
	return fmt.Errorf("tzf: not found for %v,%v", lng, lat)
}

func (i *tzitem) ContainsPoint(p geometry.Point) bool {
	// Quick rejection if the point is outside the tzitem's bounding box
	if p.X < i.min[0] || p.X > i.max[0] || p.Y < i.min[1] || p.Y > i.max[1] {
		return false
	}

	if len(i.polys) > 10 {
		return i.containsPointConcurrent(p)
	}

	// Sequential check for small number of polygons
	for _, poly := range i.polys {
		if poly.ContainsPoint(p) {
			return true
		}
	}
	return false
}

func (i *tzitem) containsPointConcurrent(p geometry.Point) bool {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	resultChan := make(chan bool, 1)

	for _, poly := range i.polys {
		wg.Add(1)
		go func(poly *geometry.Poly) {
			defer wg.Done()
			select {
			case <-ctx.Done():
				// Another goroutine found the point, no need to continue
				return
			default:
				if poly.ContainsPoint(p) {
					select {
					case resultChan <- true:
						cancel() // Signal to other goroutines to stop processing
					default:
					}
				}
			}
		}(poly)
	}

	go func() {
		wg.Wait()
		close(resultChan) // Close the channel once all checks are done
	}()

	select {
	case result := <-resultChan:
		return result
	default:
		return false
	}
}

func (i *tzitem) getMinMax() ([2]float64, [2]float64) {
	retmin := [2]float64{
		i.polys[0].Rect().Min.X,
		i.polys[0].Rect().Min.Y,
	}
	retmax := [2]float64{
		i.polys[0].Rect().Max.X,
		i.polys[0].Rect().Max.Y,
	}

	for _, poly := range i.polys {
		minx := poly.Rect().Min.X
		miny := poly.Rect().Min.Y
		if minx < retmin[0] {
			retmin[0] = minx
		}
		if miny < retmin[1] {
			retmin[1] = miny
		}

		maxx := poly.Rect().Max.X
		maxy := poly.Rect().Max.Y
		if maxx > retmax[0] {
			retmax[0] = maxx

		}
		if maxy > retmax[1] {
			retmax[1] = maxy
		}
	}
	return retmin, retmax
}

// Finder is based on point-in-polygon search algo.
//
// Memory will use about 100MB if lite data and 1G if full data.
// Performance is very stable and very accuate.
type Finder struct {
	items   []*tzitem
	names   []string
	reduced bool
	tr      *rtree.RTreeG[*tzitem]
	opt     *Option
	version string
}

func NewFinderFromRawJSON(input *convert.BoundaryFile, opts ...OptionFunc) (F, error) {
	timezones, err := convert.Do(input)
	if err != nil {
		return nil, err
	}
	return NewFinderFromPB(timezones, opts...)
}

func NewFinderFromPB(input *pb.Timezones, opts ...OptionFunc) (F, error) {
	items := make([]*tzitem, 0)
	names := make([]string, 0)

	opt := &Option{}
	for _, optFunc := range opts {
		optFunc(opt)
	}

	tr := &rtree.RTreeG[*tzitem]{}
	for _, timezone := range input.Timezones {
		names = append(names, timezone.Name)

		newItem := &tzitem{
			name: timezone.Name,
		}
		if !opt.DropPBTZ {
			newItem.pbtz = timezone
		}
		for _, polygon := range timezone.Polygons {

			newPoints := make([]geometry.Point, 0)
			for _, point := range polygon.Points {
				newPoints = append(newPoints, geometry.Point{
					X: float64(point.Lng),
					Y: float64(point.Lat),
				})
			}

			holes := [][]geometry.Point{}
			for _, holePoly := range polygon.Holes {
				newHolePoints := make([]geometry.Point, 0)
				for _, point := range holePoly.Points {
					newHolePoints = append(newHolePoints, geometry.Point{
						X: float64(point.Lng),
						Y: float64(point.Lat),
					})
				}
				holes = append(holes, newHolePoints)
			}

			newPoly := geometry.NewPoly(newPoints, holes, &geometry.IndexOptions{Kind: geometry.RTree, MinPoints: 64})
			newItem.polys = append(newItem.polys, newPoly)
		}
		minp, maxp := newItem.getMinMax()

		newItem.min = minp
		newItem.max = maxp

		items = append(items, newItem)
		tr.Insert(minp, maxp, newItem)
	}
	finder := &Finder{}
	finder.items = items
	finder.names = names
	finder.reduced = input.Reduced
	finder.tr = tr
	finder.opt = opt
	finder.version = input.Version
	return finder, nil
}

func NewFinderFromCompressed(input *pb.CompressedTimezones, opts ...OptionFunc) (F, error) {
	tzs, err := reduce.Decompress(input)
	if err != nil {
		return nil, err
	}
	return NewFinderFromPB(tzs, opts...)
}

func getRTreeRangeShifted(lng float64, lat float64) float64 {
	if 73 < lng && lng < 140 && 8 < lat && lat < 54 {
		return 70.0
	}
	return 30.0
}

func (f *Finder) getItemInRanges(lng float64, lat float64) []*tzitem {
	candidates := []*tzitem{}

	// Adjust the range based on more nuanced logic or additional conditions
	shifted := getRTreeRangeShifted(lng, lat)

	// Calculate bounds and ensure they're within valid geographic ranges
	minLng := math.Max(lng-shifted, -180)
	maxLng := math.Min(lng+shifted, 180)
	minLat := math.Max(lat-shifted, -90)
	maxLat := math.Min(lat+shifted, 90)

	f.tr.Search([2]float64{minLng, minLat}, [2]float64{maxLng, maxLat}, func(min, max [2]float64, data *tzitem) bool {
		candidates = append(candidates, data)
		return true // Continue searching
	})

	// Fallback to all items if no candidates found within the adjusted range
	if len(candidates) == 0 {
		candidates = f.items
	}

	return candidates
}

func (f *Finder) getItem(lng float64, lat float64) ([]*tzitem, error) {
	p := geometry.Point{
		X: float64(lng),
		Y: float64(lat),
	}
	ret := []*tzitem{}
	candidates := f.getItemInRanges(lng, lat)
	if len(candidates) == 0 {
		return nil, ErrNoTimezoneFound
	}
	for _, item := range candidates {
		if item.ContainsPoint(p) {
			ret = append(ret, item)
		}
	}
	if len(ret) == 0 {
		return nil, newNotFoundErr(lng, lat)
	}
	return ret, nil
}

// GetTimezoneName will use alphabet order and return first matched result.
func (f *Finder) GetTimezoneName(lng float64, lat float64) string {
	p := geometry.Point{X: lng, Y: lat}
	for _, item := range f.items {
		if item.ContainsPoint(p) {
			return item.name
		}
	}
	return ""
}

func (f *Finder) GetTimezoneName2(lng, lat float64) string {
	point := geometry.Point{X: lng, Y: lat}
	candidates := f.getCandidatesByBoundingBox(lng, lat)
	for _, item := range candidates {
		for _, poly := range item.polys {
			if poly.ContainsPoint(point) {
				return item.name
			}
		}
	}
	return ""
}

func (f *Finder) GetTimezoneNames(lng float64, lat float64) ([]string, error) {
	item, err := f.getItem(lng, lat)
	if err != nil {
		return nil, err
	}
	ret := []string{}
	for i := 0; i < len(item); i++ {
		ret = append(ret, item[i].name)
	}
	slices.Sort(ret)
	return ret, nil
}

func (f *Finder) TimezoneNames() []string {
	return f.names
}

func (f *Finder) DataVersion() string {
	return f.version
}

func (i *tzitem) getBoundingBox() (min, max [2]float64) {
	return i.min, i.max
}

// Fast initial candidate selection based on bounding box overlap.
func (f *Finder) getCandidatesByBoundingBox(lng, lat float64) []*tzitem {
	candidates := make([]*tzitem, 0)
	point := [2]float64{lng, lat}
	for _, item := range f.items {
		min, max := item.getBoundingBox()
		if point[0] >= min[0] && point[0] <= max[0] && point[1] >= min[1] && point[1] <= max[1] {
			candidates = append(candidates, item)
		}
	}
	return candidates
}
