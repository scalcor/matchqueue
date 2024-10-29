package matchqueue

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
)

type queueTestSuite struct {
	suite.Suite

	q       *queue
	players [][]*Player
}

func (ts *queueTestSuite) SetupSuite() {
	ts.q = New(DefaultConfig()).(*queue)

	// add initial data
	// #1: [1]10.0, [2]11.5
	// #2: [3]32.4, [4]9.0, [5]88.0, [6]22.1
	// #3: [7]45.0
	// #4: [8]0.004
	// #5: [9]67.8
	// #6: [10]33.3, [11]33.3, [12]22.2, [13]14.5
	data := [][]float64{{10.0, 11.5}, {32.4, 9.0, 88.0, 22.1}, {45.0}, {0.004}, {67.8}, {33.3, 33.3, 22.2, 14.5}}

	var id PlayerID
	cnt := 0
	for _, vs := range data {
		pl := []*Player{}
		for _, v := range vs {
			id++
			pl = append(pl, &Player{ID: id, Score: v})
		}
		ts.players = append(ts.players, pl)
	}

	for i, pl := range ts.players {
		ts.q.AddPlayer(pl)
		cnt += len(pl)

		// check state
		ts.Assert().Equal(cnt, ts.q.playerCnt)
		ts.Assert().Equal(i+1, len(ts.q.parties))
	}
}

func Test_queueTestSuite(t *testing.T) {
	suite.Run(t, new(queueTestSuite))
}

func (ts *queueTestSuite) TestProcMatching() {
	ts.Run("no match", func() {
		// nothing happens; group is created every 2 rounds
		for round := 1; round <= 7; round++ {
			ts.Run(fmt.Sprintf("%v", round), func() {
				mw := ts.q.matchWindow

				groups, err := ts.q.ProcMatching()
				ts.Assert().NoError(err)
				ts.Assert().Empty(groups)

				// match windows increases a little
				ts.Assert().Greater(ts.q.matchWindow, mw)
			})
		}
	})

	ts.Run("match", func() {
		// groups are created
		// #2 + #3 vs #6, #5
		var expected [2][]*Player
		for _, n := range []int{1, 2} {
			expected[0] = append(expected[0], ts.players[n]...)
		}
		for _, n := range []int{5, 4} {
			expected[1] = append(expected[1], ts.players[n]...)
		}

		groups, err := ts.q.ProcMatching()
		ts.Require().NoError(err)
		ts.Require().Equal(1, len(groups))

		got := groups[0]
		ts.Assert().EqualValues(1, got.ID)
		ts.Assert().EqualValues(expected, got.Players)
		ts.Assert().EqualValues(8, got.CreatedRound)
	})
}
