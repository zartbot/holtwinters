//Ref:https://mp.weixin.qq.com/s/-Mx3j_pgW1pVdQtEjVj1VQ
package main

import (
	"log"
	"math"

	"gonum.org/v1/gonum/optimize"
)

func main() {
	points := [][2]float64{
		{1, 9.801428},
		{2, 17.762811},
		{3, 20.222147},
		{4, 18.435252},
		{5, 12.570380},
		{6, 20.979064},
		{7, 24.313054},
		{8, 21.307317},
		{9, 26.555673},
		{10, 27.772882},
		{11, 41.202046},
		{12, 44.854088},
		{13, 40.916411},
		{14, 49.013679},
		{15, 37.969996},
		{16, 49.735623},
		{17, 48.259766},
		{18, 50.009173},
		{19, 61.297761},
		{20, 58.333159},
	}
	problem := optimize.Problem{
		Func: func(x []float64) float64 {
			sumOfResiduals := 0.0
			m := x[0] //slope
			b := x[1] //intercept
			for _, p := range points {
				actualY := p[1]
				testY := m*p[0] + b
				sumOfResiduals += math.Abs(testY - actualY)
			}
			return sumOfResiduals
		},
	}
	result, err := optimize.Minimize(problem, []float64{1, 1}, nil, &optimize.NelderMead{})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("result:", result.X)
}
