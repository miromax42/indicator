// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper_test

import (
	"sync"
	"testing"

	"github.com/miromax42/indicator/v2/helper"
)

func TestWaitable(_ *testing.T) {
	wg := &sync.WaitGroup{}
	c := make(chan int)

	helper.Waitable[int](wg, c)
	close(c)

	wg.Wait()
}
