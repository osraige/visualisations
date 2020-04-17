package timeline

import (
	"fmt"
	"hash/fnv"
	"io"
	"sort"

	svg "github.com/ajstarks/svgo/float"
	visual "github.com/osraige/visualisations"
)

type TimelineOptions struct {
	canvas  *svg.SVG
	width   float64
	height  float64
	rows    float64
	columns float64
	// Length of horizontal flat sections
	SegmentLength float64
	// Width of the lines
	LineWidth float64
	// Opacity of the lines when dropping out off the timeline
	DropoutOpacity float64
	LineCap        visual.CapStyle
	baseLineStyle  string
	// Radius of the dots at the begining of the segments
	DotRadius float64
	// Vertical distance between the lines
	GapHeight float64
	// Distance between the flat horizontal segments
	GapWidth float64
	// Ratio of the gap between segments the the bezier handles are offset
	HandleGapRatio float64
	handleOffset   float64
	// Outer horizontal paddingx
	PaddingX float64
	// Outer vertical padding
	PaddingY float64
	// Some function that when given the name of an entry will return a colour
	GetColour func(string) string
	// Some function that when given the name of an entry will return a colour
	GetLabelColour func(string) string
	// Some function that when given the name of an entry will return the label text
	GetEntryLabel func(string) string
	// Timeline entries
	Entries [][]string
	entries []entry
	// Labels for the columns
	ColumnLabels []string
	// Colour for the column labels
	ColumnLabelColour string
	// Font size for the column labels
	ColumnLabelFontSize int
	// Font for the labels
	LabelFont string
	// Font size for the entry labels
	LabelFontSize int
	// Whether or not to centre the text over the start of the entries
	CentreText bool
	// Vertical distance from the entry start to the entry label
	EntryLabelGap float64
	// Text to display when no entries are provided
	NoEntryText   string
	baseTextStyle string
}

type entry struct {
	name       string
	occurences []occurence
}

type entrySlice []entry

func (e entrySlice) Len() int {
	return len(e)
}

func (e entrySlice) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

func (e entrySlice) Less(i, j int) bool {
	return e[i].name < e[j].name
}

type occurence struct {
	row    float64
	column float64
}

func (t *TimelineOptions) drawEntries() {
	for _, e := range t.entries {
		t.canvas.Group()
		t.canvas.Title(e.name)
		t.drawEntry(e)
		t.canvas.Gend()
	}
}

func (t *TimelineOptions) drawEntry(e entry) {
	bottomY := t.rows*(t.GapHeight) + t.PaddingY
	colWidth := t.SegmentLength + t.GapWidth
	entryColour := t.GetColour(e.name)
	var prevLineX, prevLineY float64
	var prevColumn float64
	style := visual.ParseStyles(t.baseLineStyle,
		visual.ParseStroke(entryColour))
	fadedStyle := visual.ParseStyles(style,
		visual.ParseStrokeOpacity(t.DropoutOpacity))
	dotStyle := visual.ParseStyles(visual.ParseStroke(entryColour),
		visual.ParseFill(entryColour))
	textStyle := visual.ParseStyles(t.baseTextStyle,
		visual.ParseFill(t.GetLabelColour(e.name)))
	for i, o := range e.occurences {
		// Draw flat segment
		startX := o.column*colWidth + t.PaddingX
		segmentY := o.row*(t.GapHeight) + t.PaddingY
		endX := startX + t.SegmentLength
		t.canvas.Line(startX, segmentY, endX, segmentY, style)
		if i == 0 {
			// Draw start dot
			t.canvas.Circle(startX, segmentY, t.DotRadius, dotStyle)
			// Draw entry label
			t.canvas.Text(startX, segmentY-t.EntryLabelGap,
				t.GetEntryLabel(e.name), textStyle)
			prevLineX = endX
			prevLineY = segmentY
			prevColumn = o.column
			continue
		}
		// Draw drop curve and flat segment until next occurence
		if prevColumn+1 != o.column {
			x := (prevColumn+1)*colWidth + t.PaddingX
			t.connectorHelper(prevLineX, prevLineY, x,
				bottomY, fadedStyle)
			bottomEndX := x + (o.column-prevColumn-1)*
				colWidth - t.GapWidth
			t.canvas.Line(x, bottomY, bottomEndX,
				bottomY, fadedStyle)
			prevLineX = bottomEndX
			prevLineY = bottomY
		}
		curveStyle := style
		if prevLineY == bottomY {
			curveStyle = fadedStyle
		}
		// Connect previous occurence
		t.connectorHelper(prevLineX, prevLineY,
			startX, segmentY, curveStyle)
		prevLineX = endX
		prevLineY = segmentY
		prevColumn = o.column
	}
}

func (t *TimelineOptions) connectorHelper(startX, startY,
	endX, endY float64, style string) {
	t.canvas.Bezier(startX, startY,
		startX+t.handleOffset, startY,
		endX-t.handleOffset, endY,
		endX, endY, style,
	)
}

func (t *TimelineOptions) drawColumnLabels() {
	y := (t.rows+1)*(t.GapHeight) + t.PaddingY
	colWidth := t.SegmentLength + t.GapWidth
	for i := 0.0; i < t.columns; i++ {
		x := t.PaddingX + i*colWidth
		t.canvas.Text(x, y, t.ColumnLabels[int(i)],
			visual.ParseStyles(
				visual.ParseFontFamily(t.LabelFont),
				visual.ParseFontSize(t.LabelFontSize),
				visual.ParseFill(t.ColumnLabelColour),
			),
		)
	}
}

func (t *TimelineOptions) drawNoEntryText() {
	t.canvas.Text(t.width/2, t.height/2, t.NoEntryText, visual.ParseStyles(
		visual.ParseFillOpacity(0.5),
		visual.ParseFontFamily(t.LabelFont),
		visual.ParseFontSize(t.LabelFontSize),
		visual.ParseDominantBaseline("central"),
		visual.ParseTextAnchor("middle"),
	))
}

func flattenEntries(entries [][]string) []entry {
	ret := []entry{}
	entryMap := map[string][]occurence{}
	for columnIndex, column := range entries {
		for rowIndex, row := range column {
			if _, ok := entryMap[row]; !ok {
				entryMap[row] = []occurence{}
			}
			entryMap[row] = append(entryMap[row],
				occurence{
					row:    float64(rowIndex),
					column: float64(columnIndex),
				})
		}
	}
	for name, occurences := range entryMap {
		ret = append(ret, entry{
			name:       name,
			occurences: occurences,
		})
	}
	sortable := entrySlice(ret)
	sort.Sort(sortable)
	return ret
}

func defaultGetColour(name string) string {
	h := fnv.New32a()
	_, _ = h.Write([]byte(name))
	return fmt.Sprintf("#%06x", h.Sum32()%0xFFFFFF)
}

func defaultGetEntryLabel(name string) string {
	return name
}

func GetTruncatedEntryLabel(n int) func(string) string {
	return func(name string) string {
		i := n
		ret := name
		if len(name) > i {
			if i > 3 {
				i -= 3
			}
			ret = name[0:i] + "..."
		}
		return ret
	}
}

func getMaxEntryColumnLength(entries [][]string) int {
	if len(entries) == 0 {
		return 0
	}
	max := 0
	for _, column := range entries {
		if len(column) > max {
			max = len(column)
		}
	}
	return max
}

// Timeline gnerates a timeline from the given options
func Timeline(out io.Writer, opts TimelineOptions) {
	canvas := svg.New(out)
	opts.columns = float64(len(opts.Entries))
	opts.rows = float64(getMaxEntryColumnLength(opts.Entries))
	opts.width = opts.columns*
		(opts.SegmentLength+opts.GapWidth) -
		opts.GapWidth + opts.PaddingX*2
	opts.height = (opts.rows+1)*opts.GapHeight +
		opts.PaddingY*2
	opts.handleOffset = opts.GapWidth * opts.HandleGapRatio
	if opts.GetColour == nil {
		opts.GetColour = defaultGetColour
	}
	if opts.GetEntryLabel == nil {
		opts.GetEntryLabel = defaultGetEntryLabel
	}
	if opts.GetLabelColour == nil {
		opts.GetLabelColour = defaultGetColour
	}
	opts.baseTextStyle = visual.ParseStyles(
		visual.ParseFontFamily(opts.LabelFont),
		visual.ParseFontSize(opts.LabelFontSize),
		visual.ParseDominantBaseline("central"),
	)
	if opts.CentreText {
		opts.baseTextStyle = visual.ParseStyles(
			opts.baseTextStyle,
			visual.ParseTextAnchor("middle"),
		)
	}
	opts.baseLineStyle = visual.ParseStyles(
		visual.ParseStrokeWidth(opts.LineWidth),
		visual.ParseStrokeLineCap(opts.LineCap),
		visual.ParseFill("none"),
	)
	opts.entries = flattenEntries(opts.Entries)
	if len(opts.Entries) == 0 {
		opts.height = 100
		opts.width = 200
		canvas.Start(opts.width, opts.height)
		defer canvas.End()
		opts.canvas = canvas
		opts.canvas.Gid("root")
		opts.drawNoEntryText()
		opts.canvas.Gend()
		return
	}
	canvas.Start(opts.width, opts.height)
	defer canvas.End()
	opts.canvas = canvas
	opts.canvas.Gid("root")
	opts.drawEntries()
	opts.drawColumnLabels()
	opts.canvas.Gend()
}
