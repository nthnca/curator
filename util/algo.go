package util

import (
	"github.com/nthnca/curator/mediainfo/message"
)

type Game struct {
	Opponent string
	Result   int
}

type Data struct {
	Key                                 string
	Score, Views, Count, Next, Min, Max int
	Skip                                bool
	Games                               []Game
}

func maxint(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func minint(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func CalculateRankings(comparisons []*message.ComparisonEntry) map[string]Data {
	rv := make(map[string]Data)

	add := func(photo, opponent string, result int) {
		d, ok := rv[photo]
		if !ok {
			d = Data{Key: photo, Score: 5000}
		}
		d.Views++
		d.Games = append(d.Games, Game{
			Opponent: opponent,
			Result:   result})
		rv[photo] = d

	}
	for _, e := range comparisons {
		add(e.GetPhoto1(), e.GetPhoto2(), int(e.GetScore()))
		add(e.GetPhoto2(), e.GetPhoto1(), int(-e.GetScore()))
	}

	for k, v := range rv {
		i := 1
		for ; i < len(v.Games); i++ {
			if (v.Games[i-1].Result > 0) != (v.Games[i].Result > 0) {
				break
			}
		}
		if i != len(v.Games) {
			continue
		}

		v.Skip = true
		if v.Games[0].Result > 0 {
			v.Score = 10000
		} else {
			v.Score = 0
		}
		rv[k] = v
	}

	for i := 0; i < 20; i++ {
		for k, v := range rv {
			if v.Skip {
				continue
			}

			best := 0
			worst := 10000
			for _, e := range v.Games {
				if e.Result > 0 {
					best = maxint(best, rv[e.Opponent].Score)
				} else {
					worst = minint(worst, rv[e.Opponent].Score)
				}

			}
			v.Next = (best + worst) / 2
			rv[k] = v
		}
		for k, v := range rv {
			if v.Skip {
				continue
			}

			v.Score, v.Next = v.Next, v.Score
			rv[k] = v
		}
	}

	return rv
}
