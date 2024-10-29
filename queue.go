package matchqueue

import (
	"slices"
	"time"
)

type (
	Notification struct {
		Message string
		Data    map[string]any
	}

	Queue interface {
		// AddPlayers adds players to the queue and updates matching factors for the queue.
		AddPlayer([]*Player)

		// Remove player removes player from the queue.
		// All players added together will be removed together.
		RemovePlayer(PlayerID, bool)

		// State returns the queue's current state.
		State() State

		// ProcMatching does a matching process.
		ProcMatching() ([]*Group, error)
	}
)

var NowFunc = time.Now

type (
	queue struct {
		// basic info
		config Config

		// party
		parties       []*party // party list sorted by the join order
		partiesSorted []*party // party list sorted by its priority

		// match filter
		filter *matchFilter

		// match state
		matchWindow       float64
		playerCnt         int
		roundGroupCreated uint64

		state *State

		idPool GroupID
	}
)

var _ Queue = new(queue)

// New creates a new matching queue.
func New(conf *Config) Queue {
	q := &queue{config: *conf, state: &State{}}
	q.Init()
	return q
}

// Init initializes the new queue.
// It sets matching factors using its configuration.
func (q *queue) Init() {
	q.matchWindow = q.config.InitMatchWindow
	q.filter = newMatchFilter(q.config.ScoreBoundFilter, q.config.MatchingWindowFilter, q.config.ScoreModRatio)
}

// implementation of Queue
func (q *queue) AddPlayer(players []*Player) {
	if len(players) == 0 {
		return
	}

	p := newParty(q, players)
	p.UpdateWindowSize(q.matchWindow)

	q.addParty(p)
}

func (q *queue) addParty(p *party) {
	if p == nil {
		return
	}

	q.parties = append(q.parties, p)

	// partiesSorted must be sorted by party's priority, desc
	q.partiesSorted = insertSortedSlice(q.partiesSorted, p, func(i int) bool {
		return p.HasPriorityTo(q.partiesSorted[i])
	})

	// total number of players in the queue
	q.playerCnt += len(p.players)
}

func (q *queue) RemovePlayer(leader PlayerID, updateState bool) {
	if !leader.IsValid() {
		return
	}

	idx, p := q.findParty(leader)
	if p == nil {
		return
	}

	memberCnt := len(p.players)

	q.parties = slices.Delete(q.parties, idx, idx)
	q.partiesSorted = slices.DeleteFunc(q.partiesSorted, func(t *party) bool { return t.id == leader })

	// update queued player count
	q.playerCnt = max(q.playerCnt-memberCnt, 0)

	// update state
	if updateState {
		waitTime := NowFunc().Sub(p.createdAt)
		q.state.AddCanceled(uint64(max(waitTime, 0)), memberCnt)
	}
}

func (q *queue) findParty(id PlayerID) (idx int, p *party) {
	idx = slices.IndexFunc(q.parties, func(t *party) bool { return t.id == id })
	if idx >= 0 {
		p = q.parties[idx]
	}
	return
}

func (q *queue) State() State {
	return *q.state
}
