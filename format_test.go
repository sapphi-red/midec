package midec_test

import (
	"bufio"
	"os"
	"testing"

	"github.com/sapphi-red/midec"
	_ "github.com/sapphi-red/midec/gif"
	_ "github.com/sapphi-red/midec/png"
	_ "github.com/sapphi-red/midec/webp"
	_ "github.com/sapphi-red/midec/isobmff"
)

const testdataFolder = "testdata/"

func Test_IsAnimated(t *testing.T) {
	t.Parallel()

	runIsAnimated := func(filename string) (bool, error) {
		fp, err := os.Open(testdataFolder + filename)
		if err != nil {
			panic(err)
		}
		return midec.IsAnimated(fp)
	}

	testcases := []struct {
		filename           string
		expectedIsAnimated bool
		expectedHasError   bool
	}{
		{"gif/animated.gif", true, false},
		{"png/animated.png", true, false},
		{"webp/animated.webp", true, false},
		{"isobmff/animated.avif", true, false},
		{"invalid.txt", false, true},
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

func Test_IsAnimated_WithBuffer(t *testing.T) {
	fp, err := os.Open(testdataFolder + "gif/animated.gif")
	if err != nil {
		panic(err)
	}

	bfp := bufio.NewReader(fp)

	actualIsAnimated, actualErr := midec.IsAnimated(bfp)
	if !actualIsAnimated {
		t.Errorf("IsAnimated = %t; want true", actualIsAnimated)
	}
	if actualErr != nil {
		t.Errorf("Error = %v; want HasError = false", actualErr)
	}
}
