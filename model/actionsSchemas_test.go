package model

import (
	"testing"
)

func TestThing_AddActionFromSchema(t *testing.T) {
	// type args struct {
	// 	id     string
	// 	schema string
	// }
	// tests := []struct {
	// 	title string
	// 	args  args
	// 	want  *interaction.Action
	// }{
	// 	{
	// 		title: "ToggleAction schema",
	// 		args:  args{id: "toggle", schema: "ToggleAction"},
	// 		want: &interaction.Action{
	// 			Input:  nil,
	// 			Output: nil,
	// 			Interaction: interaction.Interaction{
	// 				Key:         "toggle",
	// 				Title:       "Toggle",
	// 				Description: "Toggles a boolean state on and off",
	// 				Forms:       []form.Form{},
	// 			},
	// 		},
	// 	},
	// }
	// for _, tt := range tests {
	// 	t.Run(tt.title, func(t *testing.T) {
	// 		testThing := thing.Thing{Actions: make(map[string]*interaction.Action)}
	// 		got, err := AddActionFromSchema(&testThing, tt.args.id, tt.args.schema)
	// 		assert.NoError(t, err, "should not return error")
	// 		assert.Equal(t, tt.want, got, "they should be equal")
	// 	})
	// }
}
