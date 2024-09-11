// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper_test

import (
	"testing"

	"github.com/miromax42/indicator/v2/helper"
)

func TestSqrt(t *testing.T) {
	input := helper.SliceToChan([]int{9, 81, 16, 100})
	expected := helper.SliceToChan([]int{3, 9, 4, 10})

	actual := helper.Sqrt(input)

	err := helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}
