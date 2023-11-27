package automations

import (
	"errors"
	"testing"

	"github.com/project-eria/go-wot/consumer"
	"github.com/project-eria/go-wot/mocks"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func Test_atHourSchedule(t *testing.T) {
	zerolog.SetGlobalLevel(zerolog.Disabled)

	consumedThingMock := &mocks.ConsumedThing{}
	consumedThingMock.On("ReadProperty", "timeProperty").Return("2023-11-02T14:50:10Z", nil)
	consumedThingMock.On("ReadProperty", "otherProperty").Return("", errors.New("property otherProperty not found"))
	consumedThings := map[string]consumer.ConsumedThing{
		"astral": consumedThingMock,
	}

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
				schedulingArray: []string{"at", "hour"},
			},
			want: "",
			err:  errors.New("invalid at scheduling length"),
		},
		{
			name: "Fixed hour",
			args: args{
				schedulingArray: []string{"at", "hour", "12:20"},
			},
			want: "12:20:00",
			err:  nil,
		},
		{
			name: "Fixed hour, with seconds",
			args: args{
				schedulingArray: []string{"at", "hour", "12:20:40"},
			},
			want: "12:20:40",
			err:  nil,
		},
		{
			name: "Fixed hour, extra params",
			args: args{
				schedulingArray: []string{"at", "hour", "12:00", "x"},
			},
			want: "12:00:00",
			err:  nil,
		},
		{
			name: "Fixed hour, incorrect hour",
			args: args{
				schedulingArray: []string{"at", "hour", "1200"},
			},
			want: "",
			err:  errors.New("invalid time"),
		},
		{
			name: "Thing hour, not existing thing",
			args: args{
				schedulingArray: []string{"at", "hour", "mything:myproperty"},
			},
			want: "",
			err:  errors.New("time thing doesn't exist"),
		},
		{
			name: "Thing hour, not existing property",
			args: args{
				schedulingArray: []string{"at", "hour", "astral:otherProperty"},
			},
			want: "",
			err:  errors.New("time thing property not available: property otherProperty not found"),
		},
		{
			name: "Thing hour, existing property",
			args: args{
				schedulingArray: []string{"at", "hour", "astral:timeProperty"},
			},
			want: "14:50:10",
			err:  nil,
		},
		{
			name: "Thing hour, after min",
			args: args{
				schedulingArray: []string{"at", "hour", "astral:timeProperty", "min=14:50:00"},
			},
			want: "14:50:10",
			err:  nil,
		},
		{
			name: "Thing hour, before min",
			args: args{
				schedulingArray: []string{"at", "hour", "astral:timeProperty", "min=14:50:20"},
			},
			want: "14:50:20",
			err:  nil,
		},
		{
			name: "Thing hour, before max",
			args: args{
				schedulingArray: []string{"at", "hour", "astral:timeProperty", "max=14:50:20"},
			},
			want: "14:50:10",
			err:  nil,
		},
		{
			name: "Thing hour, after max",
			args: args{
				schedulingArray: []string{"at", "hour", "astral:timeProperty", "max=14:50:00"},
			},
			want: "14:50:00",
			err:  nil,
		},
		{
			name: "Thing hour, between min and max",
			args: args{
				schedulingArray: []string{"at", "hour", "astral:timeProperty", "min=14:50:00", "max=14:50:20"},
			},
			want: "14:50:10",
			err:  nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := atHourSchedule(tt.args.schedulingArray, consumedThings)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.err, err)
		})
	}
}
