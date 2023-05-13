package slice

import (
	"fmt"
	"testing"
)

func TestRangeFunc(t *testing.T) {
	ids := []uint64{123, 4343, 343, 343, 555}
	index := "abcd"
	testJob := func(ids []uint64) {
		fmt.Println(index, ids)
	}

	RangeJob(ids, 2, testJob)
}
