package actionModel

import (
	"github.com/project-eria/go-wot/dataSchema"
)

type Meta struct {
	Title       string
	Description string
	Input       dataSchema.Data
	Output      dataSchema.Data
}
