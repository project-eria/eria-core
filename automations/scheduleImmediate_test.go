package automations

import (
	"errors"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func Test_immediateSchedule(t *testing.T) {
	zerolog.SetGlobalLevel(zerolog.Disabled)

	type args struct {
		schedulingArray []string
	}
	tests := []struct {
		name string
		args args
		want string
		err  error
	}{
		{
			name: "No params",
			args: args{
				schedulingArray: []string{"immediate"},
			},
			want: "",
			err:  nil,
		},
		{
			name: "Incorrect number of params",
			args: args{
				schedulingArray: []string{"immediate", "x"},
			},
			want: "",
			err:  errors.New("invalid immediate schedule length"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := immediateSchedule(tt.args.schedulingArray)
			assert.Empty(t, got)
			assert.Equal(t, tt.err, err)
		})
	}
}
