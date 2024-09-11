// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper_test

import (
	"testing"

	"github.com/miromax42/indicator/v2/helper"
)

func TestBuffered(_ *testing.T) {
	c := make(chan int, 1)
	b := helper.Buffered(c, 4)

	c <- 1
	c <- 2
	c <- 3
	c <- 4

	close(c)

	helper.Drain(b)
}
