package propertyModel

var LogLevel = Meta{
	Title:       "Log level",
	Description: "Application log level [error, warn, info, debug, trace]",
	Type:        "string",
	Enum:        []interface{}{"error", "warn", "info", "debug", "trace"},
}
