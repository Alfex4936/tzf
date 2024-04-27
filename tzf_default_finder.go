package tzf

import (
	"fmt"
	"runtime"
	"sync"

	"github.com/Alfex4936/tzf/pb"
	tzfrellite "github.com/ringsaturn/tzf-rel-lite"
	"google.golang.org/protobuf/proto"
)

// DefaultFinder is a finder impl combine both [FuzzyFinder] and [Finder].
//
// It's designed for performance first and allow some not so correct return at some area.
type DefaultFinder struct {
	fuzzyFinder F
	finder      F
}

func NewDefaultFinder() (F, error) {
	var (
		fuzzyFinder F
		finder      F
		err         error
	)

	if err != nil {
		return nil, err
	}

	// Initiate both finder setups concurrently to improve setup time
	ch := make(chan error, 2)

	go func() {
		fuzzyFinder, err = initFuzzyFinder(tzfrellite.PreindexData)
		ch <- err
	}()

	go func() {
		finder, err = initCompressedFinder(tzfrellite.LiteCompressData)
		ch <- err
	}()

	// Wait for both goroutines to complete
	err1, err2 := <-ch, <-ch
	if err1 != nil || err2 != nil {
		return nil, fmt.Errorf("failed to initialize finders: %v, %v", err1, err2)
	}

	if finder.DataVersion() != fuzzyFinder.DataVersion() {
		return nil, fmt.Errorf(
			"tzf: DefaultFinder only support same data version for Finder(version=%v) and FuzzyFinder(version=%v)",
			finder.DataVersion(),
			fuzzyFinder.DataVersion(),
		)
	}

	f := &DefaultFinder{
		fuzzyFinder: fuzzyFinder,
		finder:      finder,
	}

	// Force free mem by probuf, about 50~80MB
	runtime.GC()

	return f, nil
}

// initFuzzyFinder initializes a fuzzy finder from preindex data
func initFuzzyFinder(data []byte) (F, error) {
	input := &pb.PreindexTimezones{}
	if err := proto.Unmarshal(data, input); err != nil {
		return nil, err
	}
	return NewFuzzyFinderFromPB(input)
}

// initCompressedFinder initializes a compressed finder from compressed data
func initCompressedFinder(data []byte) (F, error) {
	input := &pb.CompressedTimezones{}
	if err := proto.Unmarshal(data, input); err != nil {
		return nil, err
	}
	return NewFinderFromCompressed(input, SetDropPBTZ)
}

// Optimized GetTimezoneName considers initial direct lookups before falling back to concurrent offset checks
func (f *DefaultFinder) GetTimezoneName(lng float64, lat float64) string {
	// Immediate checks without offsets
	if name := f.fuzzyFinder.GetTimezoneName(lng, lat); name != "" {
		return name
	}
	if name := f.finder.GetTimezoneName(lng, lat); name != "" {
		return name
	}

	// Use concurrency for offset checks only if immediate lookups fail
	offsets := []float64{-0.02, 0.02}
	var result string
	var once sync.Once

	var wg sync.WaitGroup
	for _, dx := range offsets {
		for _, dy := range offsets {
			wg.Add(1)
			go func(dx, dy float64) {
				defer wg.Done()
				// Offsets applied
				offsetLng, offsetLat := lng+dx, lat+dy
				if name := f.fuzzyFinder.GetTimezoneName(offsetLng, offsetLat); name != "" {
					once.Do(func() { result = name })
				} else if name := f.finder.GetTimezoneName(offsetLng, offsetLat); name != "" {
					once.Do(func() { result = name })
				}
			}(dx, dy)
		}
	}
	wg.Wait()

	return result
}

func (f *DefaultFinder) GetTimezoneNames(lng float64, lat float64) ([]string, error) {
	return f.finder.GetTimezoneNames(lng, lat)
}

func (f *DefaultFinder) TimezoneNames() []string {
	return f.finder.TimezoneNames()
}

func (f *DefaultFinder) DataVersion() string {
	return f.finder.DataVersion()
}
