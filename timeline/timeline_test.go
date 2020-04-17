package timeline

import (
	"flag"
	"os"
	"strings"
	"testing"

	"github.com/kylelemons/godebug/diff"
	visual "github.com/osraige/visualisations"
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
		golden          string
		timelineOptions TimelineOptions
	}{
		{
			golden: "empty",
			timelineOptions: TimelineOptions{
				NoEntryText: "No tags found in the last 7 days :(",
			},
		}, {
			golden: "complex",
			timelineOptions: TimelineOptions{
				SegmentLength:       40,
				LineWidth:           3,
				DropoutOpacity:      0.25,
				DotRadius:           3,
				GapHeight:           30,
				GapWidth:            60,
				HandleGapRatio:      0.3,
				PaddingX:            40,
				PaddingY:            20,
				LabelFontSize:       7,
				LabelFont:           "monospace",
				LineCap:             visual.CapStyleRound,
				CentreText:          true,
				EntryLabelGap:       15,
				GetLabelColour:      func(_ string) string { return "black" },
				GetEntryLabel:       GetTruncatedEntryLabel(14),
				ColumnLabelColour:   "grey",
				ColumnLabelFontSize: 9,
				Entries: [][]string{
					{"1", "2", "3"},
					{"2", "1", "4"},
					{"4", "1", "3"},
				},
				ColumnLabels: []string{"a", "b", "c"},
				NoEntryText:  "No tags found in the last 7 days :(",
			},
		}, {
			golden: "flat-lines",
			timelineOptions: TimelineOptions{
				SegmentLength:       40,
				LineWidth:           3,
				DropoutOpacity:      0.25,
				DotRadius:           3,
				GapHeight:           30,
				GapWidth:            60,
				HandleGapRatio:      0.3,
				PaddingX:            40,
				PaddingY:            20,
				LabelFontSize:       7,
				LabelFont:           "monospace",
				LineCap:             visual.CapStyleRound,
				CentreText:          true,
				EntryLabelGap:       15,
				GetLabelColour:      func(_ string) string { return "black" },
				GetEntryLabel:       GetTruncatedEntryLabel(14),
				ColumnLabelColour:   "grey",
				ColumnLabelFontSize: 9,
				Entries: [][]string{
					{"1", "2", "3"},
					{"1", "2", "3"},
					{"1", "2", "3"},
					{"1", "2", "3"},
					{"1", "2", "3"},
					{"1", "2", "3"},
					{"1", "2", "3"},
				},
				ColumnLabels: []string{"a", "b", "c", "d", "e", "f", "g"},
				NoEntryText:  "No tags found in the last 7 days :(",
			},
		},
	} {
		t.Run(testcase.golden, func(t *testing.T) {
			builder := &strings.Builder{}
			Timeline(builder, testcase.timelineOptions)
			got := builder.String()
			want := visualtest.GoldenValue(t, testcase.golden, got, *update)
			if got != want {
				t.Errorf("mismatched output:\n%s", diff.Diff(want, got))
			}
		})
	}
}
