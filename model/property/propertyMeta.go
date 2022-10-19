package propertyModel

type Meta struct {
	Title       string
	Description string
	Type        string
	ReadOnly    bool
	Unit        string
	Minimum     int16
	Maximum     int16
	MinLength   uint16
	MaxLength   uint16
	Pattern     string
}
