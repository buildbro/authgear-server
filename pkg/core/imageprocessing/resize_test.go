package imageprocessing

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestResize(t *testing.T) {
	Convey("Resize", t, func() {
		Convey("ResolveTargetDimension", func() {
			// 400:300 => 4:3
			imageW := 400
			imageH := 300
			cases := []struct {
				Resize    Resize
				ExpectedW int
				ExpectedH int
			}{
				// Aspect-ratio-respecting scaling mode: lfit, mfit and fixed.

				// No dimensions are specified; Default to image dimensions.
				{
					Resize{
						ScalingMode: ResizeScalingModeLfit,
					},
					400, 300,
				},
				// Only w is specified; h is derived from w and aspect ratio.
				{
					Resize{
						ScalingMode: ResizeScalingModeLfit,
						Width:       200,
					},
					200, 150,
				},
				// Only h is specified; w is derived from h and aspect ratio.
				{
					Resize{
						ScalingMode: ResizeScalingModeLfit,
						Height:      150,
					},
					200, 150,
				},
				// l is w.
				{
					Resize{
						ScalingMode: ResizeScalingModeLfit,
						LongerSide:  200,
					},
					200, 150,
				},
				// s is h.
				{
					Resize{
						ScalingMode: ResizeScalingModeLfit,
						ShorterSide: 150,
					},
					200, 150,
				},
				// w and h are specified; Use them directly.
				{
					Resize{
						ScalingMode: ResizeScalingModeLfit,
						Width:       120,
						Height:      130,
					},
					120, 130,
				},
				// w and h have precedence over l and s.
				{
					Resize{
						ScalingMode: ResizeScalingModeLfit,
						Width:       1,
						Height:      2,
						LongerSide:  3,
						ShorterSide: 4,
					},
					1, 2,
				},

				// Non-aspect-ratio-respecting scaling mode: pad and fill.

				// No dimensions are specified; Default to image dimensions.
				{
					Resize{
						ScalingMode: ResizeScalingModePad,
					},
					400, 300,
				},
				// Only w is specified; h is equal to w.
				{
					Resize{
						ScalingMode: ResizeScalingModePad,
						Width:       200,
					},
					200, 200,
				},
				// Only h is specified; w is equal to h.
				{
					Resize{
						ScalingMode: ResizeScalingModePad,
						Height:      200,
					},
					200, 200,
				},
				// l and s are w and h respectively.
				{
					Resize{
						ScalingMode: ResizeScalingModePad,
						LongerSide:  100,
						ShorterSide: 200,
					},
					100, 200,
				},
				// w and h are specified; Use them directly.
				{
					Resize{
						ScalingMode: ResizeScalingModePad,
						Width:       100,
						Height:      200,
					},
					100, 200,
				},
				// w and h have precedence over l and s.
				{
					Resize{
						ScalingMode: ResizeScalingModePad,
						Width:       1,
						Height:      2,
						LongerSide:  3,
						ShorterSide: 4,
					},
					1, 2,
				},
			}
			for _, c := range cases {
				actualW, actualH := c.Resize.ResolveTargetDimension(
					imageW,
					imageH,
				)
				So(actualW, ShouldEqual, c.ExpectedW)
				So(actualH, ShouldEqual, c.ExpectedH)
			}
		})
	})
}
