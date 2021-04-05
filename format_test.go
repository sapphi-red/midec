package midec_test

import (
	"os"
	"testing"

	"github.com/sapphi-red/midec"
	_ "github.com/sapphi-red/midec/gif"
	_ "github.com/sapphi-red/midec/png"
)

const fixtureFolder = "fixtures/"

func Test_IsAnimated(t *testing.T) {
	t.Parallel()

	runIsAnimated := func(filename string) (bool, error) {
		fp, err := os.Open(fixtureFolder + filename)
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
		{"gif/1.gif", true, false},
		{"png/1.png", true, false},
		{"invalid.txt", false, true},
	}

	for _, tc := range testcases {
		t.Run(tc.filename, func(t *testing.T) {
			actualIsAnimated, actualErr := runIsAnimated(tc.filename)
			if tc.expectedIsAnimated != actualIsAnimated {
				t.Errorf("IsAnimated = %t; want %t", actualIsAnimated, tc.expectedIsAnimated)
			}
			if tc.expectedHasError != (actualErr != nil) {
				t.Errorf("HasError = %v; want %t", actualErr, tc.expectedHasError)
			}
		})
	}
}
