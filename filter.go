package matchqueue

import "math"

const (
	defaultScoreInitial float64 = 25.0
	defaultScoreMid     float64 = defaultScoreInitial
	defaultScoreRange   float64 = defaultScoreInitial
	scoreBoundMin       float64 = 0.0
	scoreBoundMax       float64 = 100.0
	scoreBoundMid       float64 = (scoreBoundMax + scoreBoundMin) / 2.0
	scoreBoundRange     float64 = (scoreBoundMax - scoreBoundMin) / 2.0
	m2DivPi             float64 = 2.0 / math.Pi
	mPiDiv2             float64 = math.Pi / 2.0
	slope               float64 = 3.0

	partyScoreMod float64 = 1.5
)

// implementation of matchFilter
type matchFilter struct {
	scoreBoundFilter     func(float64) float64
	scoreModFilter       []func(float64) float64
	matchingWindowFilter func(float64, float64) float64
}

func newMatchFilter(scoreBoundFilter, matchingWindowFilter string, scoreModRatio []float64) *matchFilter {
	f := &matchFilter{}

	switch scoreBoundFilter {
	case "curve":
		f.scoreBoundFilter = scoreBoundFilterCurve
	case "simple":
		fallthrough
	default:
		f.scoreBoundFilter = scoreBoundFilterSimple
	}

	switch matchingWindowFilter {
	case "calculated":
		f.matchingWindowFilter = matchingWindowSizeCalculated
	case "simple":
		fallthrough
	default:
		f.matchingWindowFilter = matchingWindowSizeSimple
	}

	for _, ratio := range scoreModRatio {
		f.scoreModFilter = append(f.scoreModFilter, func(score float64) float64 {
			return scoreBoundMid + (score-scoreBoundMid)*ratio
		})
	}

	return f
}

func (f *matchFilter) AdjustPartyScore(cnt int, score float64) float64 {
	if cnt > 1 {
		return partyFilter(score)
	}
	return score
}

func (f *matchFilter) BoundScore(score float64) float64 {
	return f.scoreBoundFilter(score)
}

func (f *matchFilter) ModifyScore(retry int, score float64) float64 {
	return f.scoreModFilter[clamp(retry, 0, len(f.scoreModFilter)-1)](score)
}

func (f *matchFilter) AdjustWindow(score, window float64) float64 {
	return f.matchingWindowFilter(score, window)
}

// filters
func scoreBoundFilterSimple(score float64) float64 {
	return clamp(2.0*score, scoreBoundMin, scoreBoundMax)
}

func scoreBoundFilterCurve(score float64) float64 {
	return scoreBoundMid + scoreBoundRange*m2DivPi*math.Atan((score-defaultScoreMid)*(slope/defaultScoreRange))
}

func inverseBound(t float64) float64 {
	return defaultScoreMid + (defaultScoreRange/slope)*math.Tan((t-scoreBoundMid)*mPiDiv2/scoreBoundRange)
}

func derivativeBound(score float64) float64 {
	return (m2DivPi * slope * scoreBoundRange / defaultScoreRange) * (1.0 / (1.0 + math.Pow((slope*(score-defaultScoreMid)/defaultScoreRange), 2)))
}

// window size for bounding domain
func matchingWindowSizeSimple(score, window float64) float64 {
	return window
}

func matchingWindowSizeCalculated(score, window float64) float64 {
	// 3.82 is calculated considering base filter function with slope 3.0
	return window * derivativeBound(inverseBound(score)) / 3.82
}

// party filter
func partyFilter(score float64) float64 {
	return partyScoreMod * score
}
