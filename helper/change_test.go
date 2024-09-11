// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper_test

import (
	"testing"

	"github.com/miromax42/indicator/v2/helper"
)

func TestChange(t *testing.T) {
	input := helper.SliceToChan([]int{1, 2, 5, 5, 8, 2, 1, 1, 3, 4})
	expected := helper.SliceToChan([]int{4, 3, 3, -3, -7, -1, 2, 3})

	actual := helper.Change(input, 2)

	err := helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}
