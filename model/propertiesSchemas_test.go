package model

import (
	"testing"
)

func TestAddPropertyFromSchema(t *testing.T) {
	// type args struct {
	// 	id           string
	// 	defaultValue interface{}
	// 	schema       string
	// }
	// tests := []struct {
	// 	title string
	// 	args  args
	// 	want  *interaction.Property
	// }{
	// 	{
	// 		title: "OnOffProperty schema",
	// 		args:  args{id: "onoff", defaultValue: false, schema: "OnOffProperty"},
	// 		want: &interaction.Property{
	// 			Interaction: interaction.Interaction{
	// 				Key:         "onoff",
	// 				Title:       "On/Off",
	// 				Description: "Whether the device is turned on",
	// 				Forms:       []form.Form{},
	// 			},
	// 			Data: dataSchema.Data{
	// 				Default:    false,
	// 				ReadOnly:   false,
	// 				WriteOnly:  false,
	// 				Type:       "boolean",
	// 				DataSchema: dataSchema.Boolean{},
	// 			},
	// 		},
	// 	},
	// }
	// for _, tt := range tests {
	// 	t.Run(tt.title, func(t *testing.T) {
	// 		testThing := thing.Thing{Properties: make(map[string]*interaction.Property)}
	// 		got, err := AddPropertyFromSchema(&testThing, tt.args.id, tt.args.defaultValue, tt.args.schema)
	// 		assert.NoError(t, err, "should not return error")
	// 		assert.Equal(t, tt.want, got, "they should be equal")
	// 	})
	// }
}
