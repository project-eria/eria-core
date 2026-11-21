package propertyModel

import (
	"github.com/project-eria/go-wot/dataSchema"
	"github.com/project-eria/go-wot/interaction"
)

// A JPEG image: raw over HTTP, base64 in JSON and over WebSocket.
var jpeg, _ = dataSchema.NewString(
	dataSchema.StringContentMediaType("image/jpeg"),
	dataSchema.StringContentEncoding("base64"),
)

// Latest is a camera's latest small image (sub stream), pushed often while observed.
var Latest = Meta{
	Title:       "Latest image",
	Description: "Latest image of the sub stream (small), pushed every second while observed",
	Data:        jpeg,
	Options:     []interaction.PropertyOption{interaction.PropertyReadOnly(true)},
}

// LatestMain is a camera's latest full image (main stream), pushed while observed only.
var LatestMain = Meta{
	Title:       "Latest full image",
	Description: "Latest image of the main stream (full size), pushed every second while observed only",
	Data:        jpeg,
	Options:     []interaction.PropertyOption{interaction.PropertyReadOnly(true)},
}
