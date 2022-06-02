package main

import (
	"fmt"

	"github.com/kuritka/weight_round_robin/wrr"
)

func main() {
	pdf := []int{30, 40, 20, 10}
	c, err := wrr.NewWRR(pdf)
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}
	result := []int{0, 0, 0, 0}
	for i := 0; i < 1000; i++ {
		indexes := c.Pick()
		i0 := indexes[0]
		result[i0]++
	}

	fmt.Printf("%v => %v \n", pdf, result)
}
