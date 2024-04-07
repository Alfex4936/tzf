package tzf

import (
	"fmt"
	"runtime"

	tzfrellite "github.com/ringsaturn/tzf-rel-lite"
	"github.com/ringsaturn/tzf/pb"
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
	fuzzyFinder, err := func() (F, error) {
		input := &pb.PreindexTimezones{}
		if err := proto.Unmarshal(tzfrellite.PreindexData, input); err != nil {
			panic(err)
		}
		return NewFuzzyFinderFromPB(input)
	}()
	if err != nil {
		return nil, err
	}

	finder, err := func() (F, error) {
		input := &pb.CompressedTimezones{}
		if err := proto.Unmarshal(tzfrellite.LiteCompressData, input); err != nil {
			panic(err)
		}
		return NewFinderFromCompressed(input, SetDropPBTZ)
	}()
	if err != nil {
		return nil, err
	}

	if finder.DataVersion() != fuzzyFinder.DataVersion() {
		return nil, fmt.Errorf(
			"tzf: DefaultFinder only support same data version for Finder(version=%v) and FuzzyFinder(version=%v)",
			finder.DataVersion(),
			fuzzyFinder.DataVersion(),
		)
	}

	f := &DefaultFinder{}
	f.fuzzyFinder = fuzzyFinder
	f.finder = finder

	// Force free mem by probuf, about 80MB
	runtime.GC()

	return f, nil
}

func (f *DefaultFinder) GetTimezoneName(lng float64, lat float64) string {
	// Immediate checks without offsets
	if name := f.fuzzyFinder.GetTimezoneName(lng, lat); name != "" {
		return name
	}
	if name := f.finder.GetTimezoneName(lng, lat); name != "" {
		return name
	}

	// Preparing for parallel checks
	type result struct {
		name string
		err  error
	}
	results := make(chan result, 9) // Buffer for potential 9 offset checks

	offsets := []float64{-0.02, 0, 0.02}
	for _, dx := range offsets {
		for _, dy := range offsets {
			// Avoid the central point (0, 0) offset, already checked
			if dx == 0 && dy == 0 {
				continue
			}
			dlng, dlat := dx+lng, dy+lat
			go func(lng, lat float64) {
				var res result
				if name := f.fuzzyFinder.GetTimezoneName(lng, lat); name != "" {
					res.name = name
				} else if name := f.finder.GetTimezoneName(lng, lat); name != "" {
					res.name = name
				}
				results <- res
			}(dlng, dlat)
		}
	}

	// Collecting results
	for i := 0; i < cap(results); i++ {
		res := <-results
		if res.name != "" {
			return res.name
		}
	}
	return ""
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
