// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper_test

import (
	"testing"

	"github.com/miromax42/indicator/v2/helper"
)

func TestDecrementBy(t *testing.T) {
	input := []int{2, 3, 4, 5}
	expected := helper.SliceToChan([]int{1, 2, 3, 4})

	actual := helper.DecrementBy(helper.SliceToChan(input), 1)

	err := helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}
