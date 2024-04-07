package tzf_test

import (
	"bytes"
	"fmt"
	"math/rand"
	"runtime"
	"testing"

	"github.com/Alfex4936/tzf"
	"github.com/Alfex4936/tzf/pb"
	"github.com/loov/hrtime/hrtesting"
	gocitiesjson "github.com/ringsaturn/go-cities.json"
	tzfrellite "github.com/ringsaturn/tzf-rel-lite"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/lotsa"
	"google.golang.org/protobuf/proto"
)

var finder tzf.F = func() tzf.F {
	input := &pb.Timezones{}
	if err := proto.Unmarshal(tzfrellite.LiteData, input); err != nil {
		panic(err)
	}
	finder, err := tzf.NewFinderFromPB(input)
	if err != nil {
		panic(err)
	}
	return finder
}()

func BenchmarkGetTimezoneName(b *testing.B) {
	bench := hrtesting.NewBenchmark(b)
	defer bench.Report()
	for bench.Next() {
		_ = finder.GetTimezoneName(116.6386, 40.0786)
	}
}

func BenchmarkGetTimezoneNameAtEdge(b *testing.B) {
	bench := hrtesting.NewBenchmark(b)
	defer bench.Report()
	for bench.Next() {
		_ = finder.GetTimezoneName(110.8571, 43.1483)
	}
}

func BenchmarkGetTimezoneName_Random_WorldCities(b *testing.B) {
	bench := hrtesting.NewBenchmark(b)
	defer bench.Report()
	for bench.Next() {
		p := gocitiesjson.Cities[rand.Intn(len(gocitiesjson.Cities))]
		_ = finder.GetTimezoneName(p.Lng, p.Lat)
	}
}

func ExampleFinder_GetTimezoneName() {
	input := &pb.Timezones{}
	if err := proto.Unmarshal(tzfrellite.LiteData, input); err != nil {
		panic(err)
	}
	finder, _ := tzf.NewFinderFromPB(input)
	fmt.Println(finder.GetTimezoneName(116.6386, 40.0786))
	// Output: Asia/Shanghai
}

func TestFinder_GetTimezoneName(t *testing.T) {
	input := &pb.Timezones{}
	err := proto.Unmarshal(tzfrellite.LiteData, input)
	assert.NoError(t, err) // Use assert to check that err is nil

	finder, err := tzf.NewFinderFromPB(input)
	assert.NoError(t, err) // Again, check that err is nil with assert

	// Use assert to check the function's output
	assert.Equal(t, "Asia/Shanghai", finder.GetTimezoneName(116.6386, 40.0786))
}

func ExampleFinder_GetTimezoneNames() {
	input := &pb.Timezones{}
	if err := proto.Unmarshal(tzfrellite.LiteData, input); err != nil {
		panic(err)
	}
	finder, _ := tzf.NewFinderFromPB(input)
	fmt.Println(finder.GetTimezoneNames(87.6168, 43.8254))
	// Output: [Asia/Shanghai Asia/Urumqi] <nil>
}

func TestFinder_GetTimezoneNames(t *testing.T) {
	input := &pb.Timezones{}
	err := proto.Unmarshal(tzfrellite.LiteData, input)
	assert.NoError(t, err)

	finder, err := tzf.NewFinderFromPB(input)
	assert.NoError(t, err)

	names, err := finder.GetTimezoneNames(87.6168, 43.8254)
	// First, assert that there was no error
	assert.NoError(t, err)
	// Then, assert that the output is as expected
	expectedNames := []string{"Asia/Shanghai", "Asia/Urumqi"}
	assert.Equal(t, expectedNames, names)
}

func Test_Finder_GetTimezoneName_Random_WorldCities_Alll(t *testing.T) {
	wri := bytes.NewBufferString("")
	lotsa.Output = wri
	lotsa.Ops(len(gocitiesjson.Cities), runtime.NumCPU(), func(i, _ int) {
		city := gocitiesjson.Cities[i]
		_ = finder.GetTimezoneName(city.Lng, city.Lat)
	})
	testing.Verbose()
	t.Log(wri.String())
}

// korea
func BenchmarkGeoGetTimezoneName(b *testing.B) {
	// Pre-initialize the timezone finder to exclude initialization from the benchmark
	finder, err := tzf.NewDefaultFinder()
	if err != nil {
		b.Fatalf("Failed to initialize timezone finder: %v", err)
	}

	// Define a slice of coordinates to benchmark with
	coordinates := []struct {
		name             string  // Name of the test case
		lat, lng         float64 // Latitude and Longitude
		expectedTimezone string  // Expected result
	}{
		{"서울", 37.5665, 126.9780, "Asia/Seoul"},
		{"제주 카카오오름", 33.45049302403202, 126.57055468146439, "Asia/Seoul"},
		{"해운대", 35.1581232984585, 129.1598440928477, "Asia/Seoul"},  // 해운대 해수욕장
		{"포항", 36.08502506194445, 129.55140108962055, "Asia/Seoul"}, // 포항
		{"세종", 36.481550006080006, 127.28920084353089, "Asia/Seoul"},
		{"제주도", 33.4890, 126.4983, "Asia/Seoul"},
		{"우도", 33.51412972779723, 126.97244569597137, "Asia/Seoul"},
		{"마라도", 33.11294701534852, 126.2662987980748, "Asia/Seoul"},
		{"독도", 37.2426, 131.8597, "Asia/Seoul"},
		{"울릉도", 37.4845, 130.9057, "Asia/Seoul"},
		{"차귀도", 33.311273820042125, 126.14345298508049, "Asia/Seoul"}, // 차귀도- 제주특별자치도 제주시 한경면 고산리
		{"대강리", 38.61453830741445, 128.35799152766955, "Asia/Seoul"},  // northernmost point, 강원특별자치도 고성군 현내면
		{"백령도", 37.96909906079667, 124.609983839757, "Asia/Seoul"},    // westernmost point, 코끼리바위, 인천 옹진군 백령면 연화리 1026-29
		{"백령도2", 37.98488937628463, 124.68608584402796, "Asia/Seoul"}, // 코끼리바위, 인천 옹진군 백령면 연화리 1026-29
		{"철원", 38.31374456713513, 127.13423745903036, "Asia/Seoul"},   // 강원특별자치도 철원군 철원읍 가단리 52
		{"거제도", 34.54419719852532, 128.43864110479205, "Asia/Seoul"},  // 거제도
		{"광도", 34.269977354595504, 127.53055654653483, "Asia/Seoul"},
		{"가거도", 34.077014440034155, 125.11863713970902, "Asia/Seoul"},
		// "Etc/GMT-8"
		{"이어도", 32.124463344828854, 125.18301360832207, "Etc/GMT-8"}, // southernmost point, 이어도. cannot build stuff.
		{"Los Angeles", 34.0522, -118.2437, "America/Los_Angeles"},
		{"Tokyo", 35.6895, 139.6917, "Asia/Tokyo"},
		{"Beijing", 39.9042, 116.4074, "Asia/Shanghai"},
		{"Uni Island", 34.707351308730146, 129.43478825264333, "Asia/Tokyo"},
		{"Uni Island2", 34.43217756058352, 129.33997781093186, "Asia/Tokyo"},
		{"Uni Island3", 34.636217082470296, 129.4828167691493, "Asia/Tokyo"},
		{"Uni Island4", 34.29666974505072, 129.3871993238883, "Asia/Tokyo"},
		{"Uni Island5", 34.0854739629158, 129.2154168085643, "Asia/Tokyo"},
		{"Fukuoka 1 (Japan)", 33.784029222960406, 130.53443527389945, "Asia/Tokyo"},
		{"Fukuoka 2 (Japan)", 34.296085822281455, 130.93051474444093, "Asia/Tokyo"},
		{"Fukuoka 3 (Japan)", 32.69461329871054, 128.79495039442563, "Asia/Tokyo"},
		{"Fukuoka 4 (Japan)", 32.95445481630956, 129.09330313600782, "Asia/Tokyo"},
		{"Fukuoka 5 (Japan)", 33.53700218298737, 130.3983824405139, "Asia/Tokyo"},
		{"Shimanae (Japan)", 35.03719336610837, 132.4915325911786, "Asia/Tokyo"},
		{"Okinoshimajo (Japan)", 36.27042331297408, 133.24889805463428, "Asia/Tokyo"},
		{"Shimayama Island (Japan)", 32.683327616680096, 128.64905526405005, "Asia/Tokyo"},
		{"Kyoto (Japan)", 35.277030942449066, 135.4727941919809, "Asia/Tokyo"},
		{"Yantai (China)", 37.45460313491269, 122.43159543394779, "Asia/Shanghai"},
		{"평양 (N.Korea)", 39.040122308158885, 125.75997459218848, "Asia/Pyongyang"},
	}

	// Run the benchmark
	for _, coord := range coordinates {
		b.Run("GetTimezoneName", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = finder.GetTimezoneName(coord.lng, coord.lat)
			}
		})
	}
}

func TestGetTimezoneName(t *testing.T) {
	// Define test cases
	var tests = []struct {
		name             string  // Name of the test case
		lat, lng         float64 // Latitude and Longitude
		expectedTimezone string  // Expected result
	}{
		{"서울", 37.5665, 126.9780, "Asia/Seoul"},
		{"제주 카카오오름", 33.45049302403202, 126.57055468146439, "Asia/Seoul"},
		{"해운대", 35.1581232984585, 129.1598440928477, "Asia/Seoul"},  // 해운대 해수욕장
		{"포항", 36.08502506194445, 129.55140108962055, "Asia/Seoul"}, // 포항
		{"세종", 36.481550006080006, 127.28920084353089, "Asia/Seoul"},
		{"제주도", 33.4890, 126.4983, "Asia/Seoul"},
		{"우도", 33.51412972779723, 126.97244569597137, "Asia/Seoul"},
		{"마라도", 33.11294701534852, 126.2662987980748, "Asia/Seoul"},
		{"독도", 37.2426, 131.8597, "Asia/Seoul"},
		{"울릉도", 37.4845, 130.9057, "Asia/Seoul"},
		{"차귀도", 33.311273820042125, 126.14345298508049, "Asia/Seoul"}, // 차귀도- 제주특별자치도 제주시 한경면 고산리
		{"대강리", 38.61453830741445, 128.35799152766955, "Asia/Seoul"},  // northernmost point, 강원특별자치도 고성군 현내면
		{"백령도", 37.96909906079667, 124.609983839757, "Asia/Seoul"},    // westernmost point, 코끼리바위, 인천 옹진군 백령면 연화리 1026-29
		{"백령도2", 37.98488937628463, 124.68608584402796, "Asia/Seoul"}, // 코끼리바위, 인천 옹진군 백령면 연화리 1026-29
		{"철원", 38.31374456713513, 127.13423745903036, "Asia/Seoul"},   // 강원특별자치도 철원군 철원읍 가단리 52
		{"거제도", 34.54419719852532, 128.43864110479205, "Asia/Seoul"},  // 거제도
		{"광도", 34.269977354595504, 127.53055654653483, "Asia/Seoul"},
		{"가거도", 34.077014440034155, 125.11863713970902, "Asia/Seoul"},
		// "Etc/GMT-8"
		{"이어도", 32.124463344828854, 125.18301360832207, "Etc/GMT-8"}, // southernmost point, 이어도. cannot build stuff.
		{"Los Angeles", 34.0522, -118.2437, "America/Los_Angeles"},
		{"Tokyo", 35.6895, 139.6917, "Asia/Tokyo"},
		{"Beijing", 39.9042, 116.4074, "Asia/Shanghai"},
		{"Uni Island", 34.707351308730146, 129.43478825264333, "Asia/Tokyo"},
		{"Uni Island2", 34.43217756058352, 129.33997781093186, "Asia/Tokyo"},
		{"Uni Island3", 34.636217082470296, 129.4828167691493, "Asia/Tokyo"},
		{"Uni Island4", 34.29666974505072, 129.3871993238883, "Asia/Tokyo"},
		{"Uni Island5", 34.0854739629158, 129.2154168085643, "Asia/Tokyo"},
		{"Fukuoka 1 (Japan)", 33.784029222960406, 130.53443527389945, "Asia/Tokyo"},
		{"Fukuoka 2 (Japan)", 34.296085822281455, 130.93051474444093, "Asia/Tokyo"},
		{"Fukuoka 3 (Japan)", 32.69461329871054, 128.79495039442563, "Asia/Tokyo"},
		{"Fukuoka 4 (Japan)", 32.95445481630956, 129.09330313600782, "Asia/Tokyo"},
		{"Fukuoka 5 (Japan)", 33.53700218298737, 130.3983824405139, "Asia/Tokyo"},
		{"Shimanae (Japan)", 35.03719336610837, 132.4915325911786, "Asia/Tokyo"},
		{"Okinoshimajo (Japan)", 36.27042331297408, 133.24889805463428, "Asia/Tokyo"},
		{"Shimayama Island (Japan)", 32.683327616680096, 128.64905526405005, "Asia/Tokyo"},
		{"Kyoto (Japan)", 35.277030942449066, 135.4727941919809, "Asia/Tokyo"},
		{"Yantai (China)", 37.45460313491269, 122.43159543394779, "Asia/Shanghai"},
		{"평양 (N.Korea)", 39.040122308158885, 125.75997459218848, "Asia/Pyongyang"},
	}

	var finder tzf.F
	finder, err := tzf.NewDefaultFinder()
	if err != nil {
		t.Fatalf("Failed to initialize timezone finder: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Get timezone name for the coordinates
			timezone := finder.GetTimezoneName(tt.lng, tt.lat)

			// Assert the timezone is as expected
			assert.Equal(t, tt.expectedTimezone, timezone, "The timezone for %v should be %v", tt.name, tt.expectedTimezone)
		})
	}
}
