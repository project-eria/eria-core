package propertyModel

import (
	"github.com/project-eria/go-wot/dataSchema"
	"github.com/project-eria/go-wot/interaction"
)

type Meta struct {
	Title       string
	Description string
	Data        dataSchema.Data
	Options     []interaction.PropertyOption
}
