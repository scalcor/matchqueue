package matchqueue

import (
	"math"
	"time"
)

type party struct {
	// basic info
	q         *queue
	id        PlayerID
	players   []*Player
	createdAt time.Time

	// matching factors
	avgScore, avgScoreBound, avgScoreMod float64
	matchWindow                          float64

	// state
	waitCnt int
}

// newParty create a new party of the given players.
// Based on the given player list, the party's match score are calculated.
func newParty(q *queue, players []*Player) *party {
	if q == nil || len(players) == 0 {
		return nil
	}

	p := &party{
		q:         q,
		createdAt: NowFunc(),
	}

	p.setPlayers(players)

	return p
}

// setPlayers reset party's match scores and calculates them again based on the given player list.
func (p *party) setPlayers(players []*Player) {
	p.players = players
	p.avgScore = 0.0
	p.avgScoreBound = 0.0
	p.avgScoreMod = 0.0

	if len(players) == 0 {
		return
	}

	p.id = players[0].ID

	// adjust score using party filter
	sumScore := 0.0
	for _, pl := range p.players {
		sumScore += pl.Score
	}
	p.avgScore = p.q.filter.AdjustPartyScore(len(players), sumScore/float64(len(players)))

	// adjust matching factors
	p.AdjustMatchingFactor(0.0)
}

func (p *party) playerIDs() []PlayerID {
	ids := []PlayerID{}
	for _, pl := range p.players {
		ids = append(ids, pl.ID)
	}
	return ids
}

// AdjustMatchingFactor updates the party's matching score and window size.
func (p *party) AdjustMatchingFactor(matchWindow float64) {
	oldScore := p.avgScoreMod

	// make sure that the party's match score is bound within the queue's score boundary
	if p.avgScoreBound == 0.0 {
		p.avgScoreBound = p.q.filter.BoundScore(p.avgScore)
	}

	// modify the party's match score; its score is changed while the party stays longer in the queue
	p.avgScoreMod = p.q.filter.ModifyScore(p.waitCnt, p.avgScoreBound)

	// if the party's match score updated, also update match window size
	if oldScore != p.avgScoreMod && matchWindow > 0.0 {
		p.UpdateWindowSize(matchWindow)
	}
}

// UpdateWindowSize updates its matching window size base on the given window. (queue's)
func (p *party) UpdateWindowSize(matchWindow float64) {
	p.matchWindow = clamp(
		p.q.filter.AdjustWindow(p.avgScoreMod, matchWindow)+float64(p.waitCnt)*p.q.config.WindowAdjustPerRetry,
		p.q.config.MinMatchWindow, p.q.config.MaxMatchWindow,
	)
}

// CanMatch checks if the party can match with the target party.
func (p *party) CanMatch(t *party) bool {
	scoreDiff := math.Abs(t.avgScoreMod - p.avgScoreMod)
	return scoreDiff <= p.matchWindow
}

// HasPriorityTo checks if p has higher priority than t.
// All parties in a queue are sorted by this priority and it affects the order of matching.
func (p *party) HasPriorityTo(t *party) bool {
	// larger party is prior to the smaller one
	if len(p.players) != len(t.players) {
		return len(p.players) > len(t.players)
	}

	// party with higher score has priority
	return p.avgScoreMod > t.avgScoreMod
}
