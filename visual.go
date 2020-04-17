package visualisations

// this file should only be used for maths-y/svg-y utilities
// that all visualisations use. every other visualisation should
// be in its own file in this package

// TODO: maybe each vis should be in it's own sub package

import (
	"fmt"
	"math"
	"strings"
)

// PointOnCircum calculates the coords (float64, float64) of a point
// on a circle with origin (`cx`, `cy`), radius `r`, at angle `a`
func PointOnCircum(cx, cy, r, a float64) (float64, float64) {
	rad := a * math.Pi / 180 // deg to rad
	return cx + r*math.Cos(rad), cy + r*math.Sin(rad)
}

// ScaleRange scales the value `n` which originally was part of the
// range `nMin` `nMax`, to a new range `rMin` `rMax`.
// eg. (n: 5, rMin: 0, rMax: 10, tMin: 0, tMax: 100) -> 50
func ScaleRange(n, rMin, rMax, tMin, tMax float64) float64 {
	return ((n-rMin)/(rMax-rMin))*(tMax-tMin) + tMin
}

func ParseFill(colour string) string {
	return fmt.Sprintf("fill:%s", colour)
}

func ParseFillOpacity(opacity float64) string {
	return fmt.Sprintf("fill-opacity:%f", opacity)
}

func ParseStroke(stroke string) string {
	return fmt.Sprintf("stroke:%s", stroke)
}

func ParseStrokeWidth(width float64) string {
	return fmt.Sprintf("stroke-width:%.1f", width)
}

func ParseStrokeOpacity(opacity float64) string {
	return fmt.Sprintf("stroke-opacity:%f", opacity)
}

type CapStyle string

const (
	CapStyleButt   = CapStyle("butt")
	CapStyleRound  = CapStyle("round")
	CapStyleSquare = CapStyle("square")
)

func ParseStrokeLineCap(lineCap CapStyle) string {
	return fmt.Sprintf("stroke-linecap:%s", string(lineCap))
}

func ParseTextAnchor(anchor string) string {
	return fmt.Sprintf("text-anchor:%s", anchor)
}

func ParseDominantBaseline(anchor string) string {
	return fmt.Sprintf("dominant-baseline:%s", anchor)
}

func ParseFontFamily(font string) string {
	return fmt.Sprintf("font-family:%s", font)
}

func ParseFontSize(fontSize int) string {
	return fmt.Sprintf("font-size:%vpx", fontSize)
}

func ParseStyles(styles ...string) string {
	return strings.Join(styles, ";")
}
