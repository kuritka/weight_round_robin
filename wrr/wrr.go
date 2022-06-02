package wrr

import (
	"fmt"
	"math/rand"
	"time"
)

// WRR Weight Round Robin Alghoritm
type WRR struct {
	pdf []int
}

// NewWRR instantiate weight round robin
func NewWRR(pdf []int) (wrr *WRR, err error) {
	r := 0
	for _, v := range pdf {
		r += v
		if v < 0 || v > 100 {
			return wrr, fmt.Errorf("value %v out of range [0;100]", v)
		}
	}
	if r != 100 {
		return wrr, fmt.Errorf("sum of pdf elements must be equal to 100 perent")
	}
	rand.Seed(time.Now().UnixNano())
	wrr = new(WRR)
	wrr.pdf = pdf
	return wrr, nil
}

// Pick returns slice shuffled by pdf distribution.
// The item with the highest probability will occur more often
// at the position that has the highest probability in the PDF
func (w *WRR) Pick() (indexes []int) {
	pdf := make([]int, len(w.pdf))
	copy(pdf, w.pdf)
	balance := 100
	for i := 0; i < len(pdf); i++ {
		cdf := w.getCDF(pdf)
		index := w.pick(cdf, balance)
		indexes = append(indexes, index)

		balance -= pdf[index]
		pdf[index] = 0

		// Summary of new pdf must be 100%. Need to add missing percentage
		for q, v := range pdf {
			if v != 0 {
				pdf[q] = v
			}
		}
	}
	return indexes
}

// pick one index
func (w *WRR) pick(cdf []int, n int) int {
	r := rand.Intn(n)
	index := 0
	for r >= cdf[index] {
		index++
	}
	return index
}

func (w *WRR) getCDF(pdf []int) (cdf []int) {
	// prepare cdf
	for i := 0; i < len(pdf); i++ {
		cdf = append(cdf, 0)
	}
	cdf[0] = pdf[0]
	for i := 1; i < 4; i++ {
		cdf[i] = cdf[i-1] + pdf[i]
	}
	return cdf
}
