package holtwinters

import (
	"errors"
	"math"

	"gonum.org/v1/gonum/optimize"
)

func meanSquareError(a []float64, b []float64) (float64, error) {
	var mse float64
	if len(a) != len(b) {
		return 0, errors.New("Length Mismatch")
	}

	if len(a) == 0 {
		return 0, nil
	}

	for idx, v := range a {
		delta := v - b[idx]
		mse += (delta * delta)
	}
	return mse / float64(len(a)), nil
}

type tsSplitT struct {
	Start int
	Mid   int
	End   int
}

func tsSplit(data []float64, kfold int) ([]*tsSplitT, error) {
	result := make([]*tsSplitT, 0, 1)

	numSample := len(data)
	testSize := numSample / (kfold + 1)
	if testSize < 1 {
		return nil, errors.New("fold number too large cause testsize invalid")
	}

	for idx := 1; idx < kfold; idx++ {
		item := &tsSplitT{
			Start: 0,
			Mid:   idx * testSize,
			End:   (idx + 1) * testSize,
		}
		result = append(result, item)
	}
	return result, nil
}

func TimeSplitMSE(data []float64, alpha float64, beta float64, gamma float64, seasonLength int, kfold int) (float64, error) {
	var result float64
	numSample := len(data)
	testSize := numSample / (kfold + 1)

	tsSymbol, err := tsSplit(data, kfold)
	if err != nil {
		return 0, err
	}

	for _, tsItem := range tsSymbol {
		train := data[tsItem.Start:tsItem.Mid]
		test := data[tsItem.Mid:tsItem.End]
		model := New(train, seasonLength, testSize, alpha, beta, gamma, 3)
		model.TripleExponentialSmoothing()
		predict := model.Result[len(train):]
		mse, err := meanSquareError(predict, test)
		if err != nil {
			return 0, err
		}
		result = result + mse
	}
	return result / float64(kfold), nil
}

type HoltWinterParameter struct {
	Alpha float64
	Beta  float64
	Gamma float64
}

func Fit(data []float64, seasonLength int, kfold int) (*HoltWinterParameter, error) {
	parameter := &HoltWinterParameter{
		Alpha: 0,
		Beta:  0,
		Gamma: 0,
	}
	problem := optimize.Problem{
		Func: func(x []float64) float64 {
			tsMSE := 0.0
			alpha := x[0]
			beta := x[1]
			gamma := x[2]
			tsMSE, err := TimeSplitMSE(data, alpha, beta, gamma, seasonLength, kfold)
			if err != nil {
				return math.MaxFloat64
			}
			return tsMSE
		},
	}
	result, err := optimize.Minimize(problem, []float64{1, 1, 1}, nil, &optimize.NelderMead{})
	if err != nil {
		return parameter, err
	}
	parameter.Alpha = result.X[0]
	parameter.Beta = result.X[1]
	parameter.Gamma = result.X[2]
	return parameter, nil
}
