package matchqueue

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_newParty(t *testing.T) {
	q := &queue{filter: newMatchFilter("", "", defaultModRatio)}
	players := []*Player{{ID: 1, Score: 34.0}, {ID: 2, Score: 23.25}, {ID: 3, Score: 22.5}, {ID: 4, Score: 19.25}}

	type args struct {
		q       *queue
		players []*Player
	}
	tests := []struct {
		name string
		args args
		want *party
	}{
		{"normal", args{q, players}, &party{avgScore: 37.125, avgScoreMod: 74.25}},
		{"empty", args{q, nil}, nil},
		{"nil", args{nil, players}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newParty(tt.args.q, tt.args.players)
			if tt.want != nil {
				assert.Equal(t, tt.want.avgScore, got.avgScore)
				assert.Equal(t, tt.want.avgScoreMod, got.avgScoreMod)
			} else {
				assert.Nil(t, got)
			}
		})
	}
}

func Test_party_playerIDs(t *testing.T) {
	players := []*Player{{ID: 1}, {ID: 3}, {ID: 6}, {ID: 9}}

	tests := []struct {
		name string
		p    *party
		want []PlayerID
	}{
		{"normal", &party{players: players}, []PlayerID{1, 3, 6, 9}},
		{"empty", &party{}, []PlayerID{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.p.playerIDs()
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_party_AdjustMatchingFactor(t *testing.T) {
	type args struct {
		matchWindow float64
	}
	tests := []struct {
		name       string
		p          *party
		args       args
		wantScore  float64
		wantWindow float64
	}{
		{
			"simple",
			&party{avgScore: 24.4, q: &queue{config: Config{MaxMatchWindow: 100.0}, filter: newMatchFilter("", "", defaultModRatio)}},
			args{15.0},
			48.8, 15.0,
		},
		{"complicated",
			&party{avgScore: 24.4, q: &queue{config: Config{MaxMatchWindow: 100.0}, filter: newMatchFilter("curve", "calculated", defaultModRatio)}},
			args{15.0},
			47.712116831117335, 14.92154188734874,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.p.AdjustMatchingFactor(tt.args.matchWindow)
			assert.Equal(t, tt.wantScore, tt.p.avgScoreMod)
			assert.Equal(t, tt.wantWindow, tt.p.matchWindow)
		})
	}
}

func Test_party_CanMatch(t *testing.T) {
	type args struct {
		t *party
	}
	tests := []struct {
		name string
		p    *party
		args args
		want bool
	}{
		{"match", &party{avgScoreMod: 67.3, matchWindow: 15.2}, args{&party{avgScoreMod: 55.5}}, true},
		{"no match", &party{avgScoreMod: 67.3, matchWindow: 10.2}, args{&party{avgScoreMod: 55.5}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.p.CanMatch(tt.args.t)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_party_HasPriorityTo(t *testing.T) {
	// party with 4 players and score of 50
	p := &party{players: []*Player{{}, {}, {}, {}}, avgScoreMod: 50.0}

	type args struct {
		t *party
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"bigger", args{&party{players: []*Player{{}, {}, {}}}}, true},
		{"smaller", args{&party{players: []*Player{{}, {}, {}, {}, {}}}}, false},
		{"higher", args{&party{players: []*Player{{}, {}, {}, {}}, avgScoreMod: 49.99}}, true},
		{"lower", args{&party{players: []*Player{{}, {}, {}, {}}, avgScoreMod: 50.01}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := p.HasPriorityTo(tt.args.t)
			assert.Equal(t, tt.want, got)
		})
	}
}
