package propertyModel

type Meta struct {
	Title       string
	Description string
	Type        string
	ReadOnly    bool
	Unit        string
	Enum        []interface{}
	Minimum     int
	Maximum     int
	MinLength   uint16
	MaxLength   uint16
	Pattern     string
}
