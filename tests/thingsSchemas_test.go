package model

import (
	"testing"
)

func TestNewFromSchemas(t *testing.T) {
	// type args struct {
	// 	urn          string
	// 	ref          string
	// 	title        string
	// 	description  string
	// 	capabilities []string
	// }
	// tests := []struct {
	// 	title   string
	// 	args    args
	// 	want    *thing.Thing
	// 	wantErr bool
	// }{
	// 	{
	// 		title: "LightBasic schema",
	// 		args:  args{urn: "dev:light", title: "Light", description: "My Light", capabilities: []string{"LightBasic"}},
	// 		want: &thing.Thing{
	// 			AtContext:           "https://www.w3.org/2022/wot/td/v1.1",
	// 			AtTypes:             []string{},
	// 			ID:                  "urn:dev:light",
	// 			Title:               "Light",
	// 			Description:         "My Light",
	// 			Security:            []string{},
	// 			SecurityDefinitions: map[string]securityScheme.SecurityScheme{},
	// 		},
	// 	},
	// }
	// for _, tt := range tests {
	// 	t.Run(tt.title, func(t *testing.T) {
	// 		got, err := NewFromSchemas(tt.args.urn, tt.args.title, tt.args.description, tt.args.capabilities)
	// 		assert.NoError(t, err, "should not return error")
	// 		tt.want.Properties = got.Properties
	// 		tt.want.Actions = got.Actions
	// 		assert.Equal(t, tt.want, got, "they should be equal")
	// 	})
	// }
}
