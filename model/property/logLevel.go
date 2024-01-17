package propertyModel

import "github.com/project-eria/go-wot/dataSchema"

var level, _ = dataSchema.NewString(
	dataSchema.StringEnum([]string{"error", "warn", "info", "debug", "trace"}),
)

var LogLevel = Meta{
	Title:       "Log level",
	Description: "Application log level [error, warn, info, debug, trace]",
	Data:        level,
}
