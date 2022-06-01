package wrr

import (
	"fmt"
	"testing"
)

func TestRRWeight2(t *testing.T) {
	prob := []int{10, 10, 60, 20}
	cdf := prepareCDF(prob)
	fmt.Println(prob)
	for i := 0; i < 30; i++ {
		fmt.Println(pickIndexes(cdf))
	}
}
