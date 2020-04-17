package gauge

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
		gaugeOptions []GaugeOptions
	}{
		{
			golden: "empty",
			gaugeOptions: []GaugeOptions{{
				Size:             500,
				Padding:          30,
				GapRadians:       1,
				BackgroundColour: "white",
				Colour:           "green",
				LineWidth:        30,
				FillProportion:   0,
				Label:            "empty",
				LabelFont:        "monospace",
				LabelColour:      "white",
				LabelSize:        50,
			}},
		}, {
			golden: "some",
			gaugeOptions: []GaugeOptions{{
				Size:             500,
				Padding:          30,
				GapRadians:       1,
				BackgroundColour: "white",
				Colour:           "green",
				LineWidth:        30,
				FillProportion:   0.1,
				Label:            "some",
				LabelFont:        "monospace",
				LabelColour:      "white",
				LabelSize:        50,
			}},
		}, {
			golden: "half",
			gaugeOptions: []GaugeOptions{{
				Size:             500,
				Padding:          30,
				GapRadians:       1,
				BackgroundColour: "white",
				Colour:           "green",
				LineWidth:        30,
				FillProportion:   0.5,
				Label:            "half",
				LabelFont:        "monospace",
				LabelColour:      "white",
				LabelSize:        50,
			}},
		}, {
			golden: "most",
			gaugeOptions: []GaugeOptions{{
				Size:             500,
				Padding:          30,
				GapRadians:       1,
				BackgroundColour: "white",
				Colour:           "green",
				LineWidth:        30,
				FillProportion:   0.9,
				Label:            "most",
				LabelFont:        "monospace",
				LabelColour:      "white",
				LabelSize:        50,
			}},
		}, {
			golden: "full",
			gaugeOptions: []GaugeOptions{{
				Size:             500,
				Padding:          30,
				GapRadians:       1,
				BackgroundColour: "white",
				Colour:           "green",
				LineWidth:        30,
				FillProportion:   1,
				Label:            "full",
				LabelFont:        "monospace",
				LabelColour:      "white",
				LabelSize:        50,
			}},
		}, {
			golden: "joined",
			gaugeOptions: []GaugeOptions{{
				Size:             200,
				Padding:          10,
				GapRadians:       1,
				BackgroundColour: "white",
				Colour:           "green",
				LineWidth:        10,
				FillProportion:   0.1,
				Label:            "some",
				LabelFont:        "monospace",
				LabelColour:      "white",
				LabelSize:        20,
			}, {
				Size:             200,
				Padding:          10,
				GapRadians:       1,
				BackgroundColour: "white",
				Colour:           "green",
				LineWidth:        10,
				FillProportion:   0.5,
				Label:            "half",
				LabelFont:        "monospace",
				LabelColour:      "white",
				LabelSize:        20,
			}, {
				Size:             200,
				Padding:          10,
				GapRadians:       1,
				BackgroundColour: "white",
				Colour:           "green",
				LineWidth:        10,
				FillProportion:   0.9,
				Label:            "most",
				LabelFont:        "monospace",
				LabelColour:      "white",
				LabelSize:        20,
			}},
		},
	} {
		t.Run(testcase.golden, func(t *testing.T) {
			builder := &strings.Builder{}
			Gauges(builder, testcase.gaugeOptions)
			got := builder.String()
			want := visualtest.GoldenValue(t, testcase.golden, got, *update)
			if got != want {
				t.Errorf("mismatched output:\n%s", diff.Diff(want, got))
			}
		})
	}
}
