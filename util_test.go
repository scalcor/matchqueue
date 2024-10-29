package matchqueue

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_clamp(t *testing.T) {
	type args struct {
		v     float64
		lower float64
		upper float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{"in range", args{0.00001, 0.000001, 0.0001}, 0.00001},
		{"lower", args{0, 11.30123, 21.54546}, 11.30123},
		{"upper", args{0.001, 0.000001, 0.0001}, 0.0001},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := clamp(tt.args.v, tt.args.lower, tt.args.upper)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_insertSortedSlice(t *testing.T) {
	type (
		ty struct {
			b float64
		}
		args struct {
			s []ty
			e ty
		}
	)

	f := func(s []ty, e ty) func(int) bool {
		return func(i int) bool {
			return s[i].b >= e.b
		}
	}

	tests := []struct {
		name string
		args args
		want []ty
	}{
		{"insert middle", args{[]ty{{0.45}}, ty{1.55}}, []ty{{0.45}, {1.55}}},
		{"append back", args{[]ty{{1.5}, {1.505}}, ty{1.51}}, []ty{{1.5}, {1.505}, {1.51}}},
		{"append front", args{[]ty{{1.5}, {1.505}}, ty{1.499}}, []ty{{1.499}, {1.5}, {1.505}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := insertSortedSlice(tt.args.s, tt.args.e, f(tt.args.s, tt.args.e))
			assert.Equal(t, tt.want, got)
		})
	}
}
