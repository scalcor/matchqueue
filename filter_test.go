package matchqueue

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_matchFilter_AdjustPartyScore(t *testing.T) {
	type (
		args struct {
			cnt   int
			score float64
		}
		test struct {
			name string
			args args
			want float64
		}
	)

	f := &matchFilter{}

	tests := []test{
		{"big party", args{5, 10.0}, 15.0},
		{"single party", args{1, 10.0}, 10.0},
		{"empty party", args{0, 10.0}, 10.0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := f.AdjustPartyScore(tt.args.cnt, tt.args.score)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_matchFilter_BoundScore(t *testing.T) {
	type (
		args struct {
			score float64
		}
		test struct {
			name string
			args args
			want float64
		}
	)

	t.Run("simple", func(t *testing.T) {
		f := newMatchFilter("simple", "", nil)

		tests := []test{
			{"in range", args{10.0}, 20.0},
			{"oob lower", args{-10.0}, 0.0},
			{"oob upper", args{110.0}, 100.0},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got := f.BoundScore(tt.args.score)
				assert.Equal(t, tt.want, got)
			})
		}
	})

	t.Run("curve", func(t *testing.T) {
		f := newMatchFilter("curve", "", nil)

		tests := []test{
			{"in range", args{10.0}, 16.141446721709528},
			{"oob lower", args{-10.0}, 7.44027652986172},
			{"oob upper", args{110.0}, 96.88925592415524},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got := f.BoundScore(tt.args.score)
				assert.Equal(t, tt.want, got)
			})
		}
	})
}

func Test_matchFilter_ModifyScore(t *testing.T) {
	f := newMatchFilter("", "", defaultModRatio)

	type args struct {
		retry int
		score float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{"first", args{0, 10.0}, 10.0},
		{"seventh lower", args{6, 10.0}, 26.0},
		{"tenth lower", args{9, 10.0}, 42.0},
		{"seventh upper", args{6, 85.0}, 71.0},
		{"tenth upper", args{9, 85.0}, 57.0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := f.ModifyScore(tt.args.retry, tt.args.score)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_matchFilter_AdjustWindow(t *testing.T) {
	type (
		args struct {
			score  float64
			window float64
		}
		test struct {
			name string
			args args
			want float64
		}
	)

	t.Run("simple", func(t *testing.T) {
		f := newMatchFilter("", "simple", nil)

		tests := []test{
			{"normal", args{25.0, 10.0}, 10.0},
			{"retry", args{25.0, 10.0}, 10.0},
			{"zero", args{25.0, 0.0}, 0.0},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got := f.AdjustWindow(tt.args.score, tt.args.window)
				assert.Equal(t, tt.want, got)
			})
		}
	})

	t.Run("calculated", func(t *testing.T) {
		f := newMatchFilter("", "calculated", nil)

		tests := []test{
			{"normal 1", args{25.0, 10.0}, 4.99963172016425},
			{"normal 2", args{72.5, 23.0}, 13.298016797524529},
			{"oob lower", args{-5.0, 10.0}, 0.24469939362609208},
			{"oob upper", args{115.0, 10.0}, 2.0609219281580544},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got := f.AdjustWindow(tt.args.score, tt.args.window)
				assert.Equal(t, tt.want, got)
			})
		}
	})
}
