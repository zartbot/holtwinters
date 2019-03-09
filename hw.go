//reference by https://github.com/DmitrySerg/DataStart
//rewrite in Golang
package holtwinters

import "math"

type HoltWintersT struct {
	Series             []float64
	SeasonLenth        int
	Npreds             int
	Alpha              float64
	Beta               float64
	Gamma              float64
	ScalingFactor      float64
	Result             []float64
	Smooth             []float64
	Season             []float64
	Trend              []float64
	PredictedDeviation []float64
	UpperBond          []float64
	LowerBond          []float64
}

func New(data []float64, slen int, npreds int, alpha float64, beta float64, gamma float64, scalingfactor float64) *HoltWintersT {
	return &HoltWintersT{
		Series:        data,
		SeasonLenth:   slen,
		Npreds:        npreds,
		Alpha:         alpha,
		Beta:          beta,
		Gamma:         gamma,
		ScalingFactor: scalingfactor,
	}
}

func (h *HoltWintersT) initTrend() float64 {
	sum := float64(0.0)
	for idx := 0; idx < h.SeasonLenth; idx++ {
		sum = sum + float64(h.Series[idx+h.SeasonLenth]-h.Series[idx])/float64(h.SeasonLenth)
	}
	return sum / float64(h.SeasonLenth)
}

func (h *HoltWintersT) seriesSum(start int, end int) float64 {
	sum := float64(0.0)
	for idx := start; idx < end; idx++ {
		sum = sum + h.Series[idx]
	}
	return sum
}

func (h *HoltWintersT) initSeasonalComponets() map[int]float64 {
	seasonals := make(map[int]float64)
	seasonAverages := make([]float64, 0, 1)
	Nseasons := len(h.Series) / h.SeasonLenth
	SumOfValsOverAvg := float64(0.0)
	for j := 0; j < Nseasons; j++ {
		seasonAverages = append(seasonAverages, h.seriesSum(h.SeasonLenth*j, h.SeasonLenth*j+h.SeasonLenth)/float64(h.SeasonLenth))
	}

	for i := 0; i < h.SeasonLenth; i++ {
		SumOfValsOverAvg = float64(0.0)
		for k := 0; k < Nseasons; k++ {
			SumOfValsOverAvg = SumOfValsOverAvg + h.Series[h.SeasonLenth*k+i] - seasonAverages[k]
		}
		seasonals[i] = SumOfValsOverAvg / float64(Nseasons)
	}
	return seasonals
}

func (h *HoltWintersT) TripleExponentialSmoothing() {
	seasonals := h.initSeasonalComponets()

	SeriesLenth := len(h.Series)

	totalLength := SeriesLenth + h.Npreds
	var smooth, lastsmooth, trend float64
	for i := 0; i < totalLength; i++ {
		if i == 0 {
			smooth = h.Series[0]
			trend = h.initTrend()
			h.Result = append(h.Result, h.Series[0])
			h.Smooth = append(h.Smooth, smooth)
			h.Trend = append(h.Trend, trend)
			h.Season = append(h.Season, seasonals[i%h.SeasonLenth])
			h.PredictedDeviation = append(h.PredictedDeviation, 0)
			h.UpperBond = append(h.UpperBond, h.Result[0]+h.ScalingFactor*h.PredictedDeviation[0])
			h.LowerBond = append(h.LowerBond, h.Result[0]-h.ScalingFactor*h.PredictedDeviation[0])
			continue
		}
		if i >= SeriesLenth { //Prediction
			m := float64(i - SeriesLenth + 1)
			h.Result = append(h.Result, (smooth+m*trend)+seasonals[i%h.SeasonLenth])
			h.PredictedDeviation = append(h.PredictedDeviation, h.PredictedDeviation[i-1]*1.01)
		} else {
			val := h.Series[i]
			lastsmooth, smooth = smooth, h.Alpha*(val-seasonals[i%h.SeasonLenth])+(1-h.Alpha)*(smooth+trend)
			trend = h.Beta*(smooth-lastsmooth) + (1-h.Beta)*trend
			seasonals[i%h.SeasonLenth] = h.Gamma*(val-smooth) + (1-h.Gamma)*seasonals[i%h.SeasonLenth]
			h.Result = append(h.Result, smooth+trend+seasonals[i%h.SeasonLenth])

			//BrutLag
			h.PredictedDeviation = append(h.PredictedDeviation, h.Gamma*math.Abs(h.Series[i]-h.Result[i])+(1-h.Gamma)*h.PredictedDeviation[i-1])
		}
		h.UpperBond = append(h.UpperBond, h.Result[i]+h.ScalingFactor*h.PredictedDeviation[0])
		h.LowerBond = append(h.UpperBond, h.Result[i]+h.ScalingFactor*h.PredictedDeviation[0])
		h.Smooth = append(h.Smooth, smooth)
		h.Trend = append(h.Trend, trend)
		h.Season = append(h.Season, seasonals[i%h.SeasonLenth])
	}
}
