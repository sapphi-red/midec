package webp

import (
	"os"
	"testing"
)

const testdataFolder = "../testdata/webp/"

func Test_isAnimated(t *testing.T) {
	t.Parallel()

	runIsAnimated := func(filename string) (bool, error) {
		fp, err := os.Open(testdataFolder + filename)
		if err != nil {
			panic(err)
		}
		return isAnimated(fp)
	}

	testcases := []struct {
		filename           string
		expectedIsAnimated bool
		expectedHasError   bool
	}{
		{"animated.webp", true, false},
		{"static-vp8.webp", false, false},
		{"static-vp8x.webp", false, false},
		{"static-vp8x-1frame.webp", false, false},
		{"invalid-firstchunk-header1.webp", false, true},
		{"invalid-firstchunk-header2.webp", false, true},
		{"invalid-vp8x-chunk1.webp", false, true},
		{"invalid-vp8x-chunk2.webp", false, true},
		{"invalid-chunk-header1.webp", false, true},
		{"invalid-chunk-header2.webp", false, true},
		{"invalid-chunk-header3.webp", false, true},
		{"invalid-unknown-chunk.webp", false, true},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.filename, func(t *testing.T) {
			t.Parallel()

			actualIsAnimated, actualErr := runIsAnimated(tc.filename)
			if tc.expectedIsAnimated != actualIsAnimated {
				t.Errorf("IsAnimated = %t; want %t", actualIsAnimated, tc.expectedIsAnimated)
			}
			if tc.expectedHasError != (actualErr != nil) {
				t.Errorf("Error = %v; want HasError = %t", actualErr, tc.expectedHasError)
			}
		})
	}
}
