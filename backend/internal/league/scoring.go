package league

import (
	"ecpc-league/internal/messages"
	"fmt"
	"sort"
)

func checkScoringIsSane(scoring []messages.LeagueScoring) error {

	sort.Slice(scoring, func(i, j int) bool {
		return scoring[i].Rank < scoring[j].Rank
	})

	for i, v := range scoring {
		if i+1 != v.Rank {
			return fmt.Errorf("scoring is missing a rank")
		}
	}

	for i := 1; i < len(scoring); i++ {
		if scoring[i].Score > scoring[i-1].Score {
			return fmt.Errorf("scoring isn't sorted correctly")
		}
	}

	return nil
}
