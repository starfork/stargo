package number

import (
	"strconv"

	"golang.org/x/exp/constraints"
)

func String[T constraints.Signed](i T) string {
	return strconv.FormatInt(int64(i), 10)
}
