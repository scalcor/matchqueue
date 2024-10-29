package matchqueue

import (
	"math"
	"sort"
)

// ProcMatching does a matching process.
func (q *queue) ProcMatching() ([]*Group, error) {
	q.state.Round++
	q.state.PlayerQueued = q.playerCnt
	q.state.MatchWindow = q.matchWindow

	var created []*Group

	if q.playerCnt > 0 {
		oldCnt := len(q.parties)

		// check prerequisites
		procCreate := q.state.Round >= q.roundGroupCreated+uint64(q.config.NumRoundToCreateGroup) ||
			q.playerCnt >= q.config.NumPlayerToCreateGroup

		if procCreate {
			var err error
			if created, err = q.ProcCreate(); err != nil {
				return nil, err
			} else if len(created) > 0 {
				q.roundGroupCreated = q.state.Round
			}
		}

		// TODO implement other types of matching

		processed := procCreate

		if processed {
			// handle remaining players
			for _, p := range q.parties {
				p.waitCnt++

				// adjust the party's matching factors due to change of wait count
				p.AdjustMatchingFactor(q.matchWindow)
			}
		}

		// adjust queue's match window
		q.adjustMatchWindow(oldCnt)
	}

	return created, nil
}

// ProcCreate commits process to create groups.
func (q *queue) ProcCreate() ([]*Group, error) {
	if q.playerCnt < q.config.MinNumToCreateGroup {
		return nil, ErrNotEnoughPlayer
	}

	results := q.procCreateImpl(0, len(q.parties))

	groups := []*Group{}
	for _, candidates := range results {
		if g := q.newGroup(candidates); g != nil {
			// a group created
			groups = append(groups, g)

			// remove matched parties
			for _, cand := range candidates {
				q.RemovePlayer(cand.id, false)
			}

			// a group created
			q.state.GroupCreated++
		}
	}

	return groups, nil
}

func (q *queue) procCreateImpl(start, count int) (results [][]*party) {
	matched := map[PlayerID]struct{}{}

	// cut party list to [start:start+count]
	parties := q.partiesSorted[start : start+count]

	for baseIdx, baseP := range parties {
		if _, ok := matched[baseP.id]; ok {
			// already matched
			continue
		}

		var (
			candidates []*party
			playerCnt  int
		)

		for idx := baseIdx; idx < len(parties); idx++ {
			p := parties[idx]
			if len(p.players) == 0 {
				continue
			}
			if _, ok := matched[p.id]; ok {
				// already matched
				continue
			}

			if playerCnt+len(p.players) > int(q.config.MaxNumToCreateGroup) {
				continue
			}

			// check if the base party can match with current one
			ok := playerCnt == 0 || baseP.CanMatch(p)
			if !ok {
				continue
			}

			candidates = append(candidates, p)
			playerCnt += len(p.players)

			if playerCnt >= q.config.MaxNumToCreateGroup && len(candidates) >= 2 {
				// enough players are gathered
				break
			}
		}

		// at least 2 candidates are required (number of team = 2)
		if playerCnt >= q.config.MinNumToCreateGroup && len(candidates) >= 2 {
			for _, cand := range candidates {
				matched[cand.id] = struct{}{}
			}
			results = append(results, candidates)
		}
	}

	return
}

func (q *queue) newGroup(candidates []*party) *Group {
	g := &Group{ID: q.idPool + 1, CreatedRound: q.state.Round}

	// sort candidates by number of players, descending
	sort.Slice(candidates, func(i, j int) bool {
		return len(candidates[i].players) > len(candidates[j].players)
	})

	// add candidates to the team of lesser players
	team := 0
	for _, cand := range candidates {
		// if number of players in the current team is bigger than the opposite, change team
		if len(g.Players[team]) > len(g.Players[1-team]) {
			team = 1 - team
		}
		g.Players[team] = append(g.Players[team], cand.players...)
	}

	// if number of players of each team differ more than 1, fail to create a group
	if math.Abs(float64(len(g.Players[0])-len(g.Players[1]))) > 1 {
		return nil
	}

	q.idPool++
	return g
}

// adjustMatchWindow adjusts the queue's match window.
// It affects all match windows of matching parties.
// The match window is adjusted as below:
//
//	rate = num_matched_party / num_total_party
//	if rate < min_rate, match_window += step
//	if rate > max_rate, match_window -= step
func (q *queue) adjustMatchWindow(oldCnt int) {
	if oldCnt == 0 || q.config.WindowAdjustStep <= 0.0 {
		return
	}

	oldWindow := q.matchWindow

	// rate = percentage of matched parties
	rate := 1.0 - float64(len(q.parties))/float64(oldCnt)

	// if the rate is out of keep range, increase queue's match window
	if rate < q.config.MinRateToKeepWindow || rate > q.config.MaxRateToKeepWindow {
		q.matchWindow = min(q.matchWindow+q.config.WindowAdjustStep, q.config.MaxMatchWindow)
	} else {
		q.matchWindow = max(q.matchWindow-q.config.WindowAdjustStep, q.config.MinMatchWindow)
	}

	if q.matchWindow != oldWindow {
		// match window changed; update all parties' window
		for _, p := range q.parties {
			p.UpdateWindowSize(q.matchWindow)
		}
	}
}
