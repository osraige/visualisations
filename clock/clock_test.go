package clock

import (
	"flag"
	"os"
	"strings"
	"testing"

	"github.com/kylelemons/godebug/diff"
	"github.com/osraige/visualisations/visualtest"
)

var (
	update = flag.Bool("update", false, "update the golden files of this test")
)

func TestMain(m *testing.M) {
	flag.Parse()
	os.Exit(m.Run())
}

func TestGenerate(t *testing.T) {
	for _, testcase := range []struct {
		golden       string
		clockOptions ClockOptions
	}{
		{
			golden: "empty",
			clockOptions: ClockOptions{
				Size:               500,
				CenterRadius:       100,
				HandGap:            3,
				Segments:           24,
				Colour:             "#33065d",
				ColourAccent:       "#ad9bbe",
				ColourAverage:      "orange",
				AverageStrokeWidth: 3,
				AveragePointRadius: 5.5,
				DataHands:          []int{},
				DataAverage:        []int{},
			},
		}, {
			golden: "complex",
			clockOptions: ClockOptions{
				Size:               500,
				CenterRadius:       100,
				HandGap:            3,
				Segments:           24,
				Colour:             "#33065d",
				ColourAccent:       "#ad9bbe",
				ColourAverage:      "orange",
				AverageStrokeWidth: 3,
				AveragePointRadius: 5.5,
				DataHands: []int{
					1, 2, 3, 4, 5, 6,
					7, 8, 9, 10, 11, 12,
					1, 2, 3, 4, 5, 6,
					7, 8, 9, 10, 11, 12,
				},
				DataAverage: []int{
					12, 11, 10, 9, 8, 7,
					6, 5, 4, 3, 2, 1,
					12, 11, 10, 9, 8, 7,
					6, 5, 4, 3, 2, 1,
				},
			},
		},
	} {
		t.Run(testcase.golden, func(t *testing.T) {
			// t.Parallel()
			builder := &strings.Builder{}
			Clock(builder, testcase.clockOptions)
			got := builder.String()
			want := visualtest.GoldenValue(t, testcase.golden, got, *update)
			if got != want {
				t.Errorf("mismatched output:\n%s", diff.Diff(want, got))
			}
		})
	}
}
