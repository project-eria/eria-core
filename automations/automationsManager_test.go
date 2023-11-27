package automations

import (
	"errors"
	"testing"
	"time"

	"github.com/project-eria/go-wot/consumer"
	"github.com/project-eria/go-wot/interaction"
	"github.com/project-eria/go-wot/mocks"
	"github.com/project-eria/go-wot/producer"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func Test_getJob(t *testing.T) {
	zerolog.SetGlobalLevel(zerolog.Disabled)
	now := time.Date(2000, time.January, 1, 12, 0, 0, 0, time.UTC)
	exposedMock := &mocks.ExposedThing{}
	aAction := interaction.NewAction(
		"a",
		"No Input, No Output",
		"",
		nil,
		nil,
	)
	exposedAction := producer.NewExposedAction(aAction)
	exposedMock.On("ExposedAction", "on").Return(exposedAction, nil)
	exposedMock.On("ExposedAction", "off").Return(nil, errors.New("exposed action not found"))

	thingMock := &mocks.ConsumedThing{}
	thingMock.On("ReadProperty", "timeProperty").Return("2023-11-02T14:50:10Z", nil)
	thingMock.On("ReadProperty", "otherProperty").Return("", errors.New("property otherProperty not found"))
	consumedMocks := map[string]consumer.ConsumedThing{
		"myThing": thingMock,
	}

	contextsMock := &mocks.ConsumedThing{}
	contextsMock.On("ReadProperty", "away").Return(false, nil)
	contextsMock.On("ReadProperty", "holiday").Return(true, nil)

	onAction := Action{
		Ref:        "on",
		Handler:    nil,
		Parameters: make(map[string]interface{}),
	}

	job := Job{
		Name:          "automation name",
		Action:        onAction,
		ScheduledType: "immediate",
		Scheduled:     "",
	}
	type args struct {
		action string
		groups []Group
	}
	tests := []struct {
		name string
		args args
		want Job
	}{
		{
			name: "No groups/No action",
			args: args{
				action: "",
				groups: []Group{},
			},
			want: Job{},
		},
		{
			name: "No groups/With action",
			args: args{
				action: "on",
				groups: []Group{},
			},
			want: Job{},
		},
		{
			name: "With groups/No action",
			args: args{
				action: "",
				groups: []Group{
					{
						Conditions: []string{"time|before=13:00"},
						Scheduled:  "immediate",
					},
				},
			},
			want: Job{},
		},
		{
			name: "Not existing action",
			args: args{
				action: "off",
				groups: []Group{
					{
						Conditions: []string{"time|before=13:00"},
						Scheduled:  "immediate",
					},
				},
			},
			want: Job{},
		},
		{
			name: "No matching time condition",
			args: args{
				action: "on",
				groups: []Group{
					{
						Conditions: []string{"time|after=13:00"},
						Scheduled:  "immediate",
					},
				},
			},
			want: Job{},
		},
		{
			name: "Matching time condition",
			args: args{
				action: "on",
				groups: []Group{
					{
						Conditions: []string{"time|before=13:00"},
						Scheduled:  "immediate",
					},
				},
			},
			want: job,
		},
		{
			name: "No Matching multiple conditions",
			args: args{
				action: "on",
				groups: []Group{
					{
						Conditions: []string{"time|before=13:00", "context|away"},
						Scheduled:  "immediate",
					},
				},
			},
			want: Job{},
		},
		{
			name: "Matching multiple conditions",
			args: args{
				action: "on",
				groups: []Group{
					{
						Conditions: []string{"time|before=13:00", "context|!away"},
						Scheduled:  "immediate",
					},
				},
			},
			want: job,
		},
		{
			name: "No Matching multiple groups",
			args: args{
				action: "on",
				groups: []Group{
					{
						Conditions: []string{"context|away"},
						Scheduled:  "immediate",
					},
					{
						Conditions: []string{"context|!holiday"},
						Scheduled:  "immediate",
					},
				},
			},
			want: Job{},
		},
		{
			name: "Matching multiple groups #1",
			args: args{
				action: "on",
				groups: []Group{
					{
						Conditions: []string{"context|!away"},
						Scheduled:  "immediate",
					},
					{
						Conditions: []string{"context|!holiday"},
						Scheduled:  "immediate",
					},
				},
			},
			want: job,
		},
		{
			name: "Matching multiple groups #2",
			args: args{
				action: "on",
				groups: []Group{
					{
						Conditions: []string{"context|away"},
						Scheduled:  "immediate",
					},
					{
						Conditions: []string{"context|holiday"},
						Scheduled:  "immediate",
					},
				},
			},
			want: job,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getJob(now, "automation name", tt.args.action, tt.args.groups, contextsMock, exposedMock, consumedMocks)
			assert.Equal(t, tt.want, got)
		})
	}
}
