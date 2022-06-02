package wrr

import (
	"fmt"
	"testing"
)

func TestRRWeight(t *testing.T) {
	pdf := []int{30, 40, 20, 10}
	wrr, err := NewWRR(pdf)
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}
	result := map[int][]int{}
	for i := 0; i < len(pdf); i++ {
		result[i] = []int{0, 0, 0, 0}
	}

	for i := 0; i < 1000; i++ {
		indexes := wrr.Pick()
		fmt.Println(indexes)
		result[0][indexes[0]]++
		result[1][indexes[1]]++
		result[2][indexes[2]]++
		result[3][indexes[3]]++
	}

	fmt.Printf("    [10.0.0.1],[10.1.0.1],[10.2.0.1],[10.3.0.1]\n")
	fmt.Printf("    %v\n", pdf)
	fmt.Printf("    -----------------\n")
	for i := 0; i < len(pdf); i++ {
		fmt.Printf(" %v. %v \n", i, result[i])
	}
}
