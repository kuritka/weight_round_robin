package wrr

import "math/rand"

// I want to avoid floating point computing, using percentage as integer instead of float numbers
func prepareCDF(pdf []int) (cdf [][]int) {
	var probabilities [][]int
	//validation
	// ....

	// init
	for i := 0; i < len(pdf); i++ {
		cdf = append(cdf, []int{})
		probabilities = append(probabilities, []int{})
		for p := 0; p < len(pdf)-i; p++ {
			cdf[i] = append(cdf[i], 0)
			probabilities[i] = append(probabilities[i], pdf[p+i])
		}
	}

	// normalize. Each probability pdf has to have  100% in total
	for i := 0; i < len(probabilities); i++ {
		sum := 0
		for p := 0; p < len(probabilities[i]); p++ {
			sum += probabilities[i][p]
		}
		x := (100 - sum) / len(probabilities[i])
		for p := 0; p < len(probabilities[i]); p++ {
			probabilities[i][p] += x
		}
		if probabilities[i][len(probabilities[i])-1]%2 == 1 {
			probabilities[i][len(probabilities[i])-1]++
		}
	}

	// count cdf
	for i := 0; i < len(pdf); i++ {
		cdf[i][0] = probabilities[i][0]
		for p := 1; p < len(probabilities[i]); p++ {
			cdf[i][p] = cdf[i][p-1] + probabilities[i][p]
		}
	}
	return cdf
}

func pickIndexes(cdf [][]int) (p []int) {
	bucket := func(cdf []int) int {
		r := rand.Intn(100)
		bucket := 0
		for r > cdf[bucket] {
			bucket++
		}
		return bucket
	}
	for _, v := range cdf {
		//p = append(p, bucket(v) + i)
		p = append(p, bucket(v))
	}
	return p
}
