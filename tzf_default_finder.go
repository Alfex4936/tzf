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

	// Setup for concurrent offset checks
	var wg sync.WaitGroup
	results := make(chan string, 1) // Use a buffered channel

	checkAndSend := func(lng, lat float64) {
		defer wg.Done()
		if name := f.fuzzyFinder.GetTimezoneName(lng, lat); name != "" {
			select {
			case results <- name:
			default:
			}
		} else if name := f.finder.GetTimezoneName(lng, lat); name != "" {
			select {
			case results <- name:
			default:
			}
		}
	}

	// Launch goroutines for offset checks
	offsets := []float64{-0.02, 0.02} // Only need to check non-zero offsets
	for _, dx := range offsets {
		for _, dy := range offsets {
			wg.Add(1)
			go checkAndSend(lng+dx, lat+dy)
		}
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	// Return the first result received
	if result, ok := <-results; ok {
		return result
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
