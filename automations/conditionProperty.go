package automations

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/project-eria/go-wot/producer"
	zlog "github.com/rs/zerolog/log"
)

var _conditionPropertyRe = regexp.MustCompile(`^(\w+)(==|!=|<=|>=|<|>)(.+)$`)

// conditionProperty compares a property of the automation's things against a
// reference value. It matches if at least one thing matches (OR, like the
// context condition).
type conditionProperty struct {
	property string
	operator string
	ref      string
	things   map[string]producer.ExposedThing
}

/**
 * `property|<name><operator><value>`
 *	- operators: ==, !=, <, <=, >, >=
 *	- e.g. `property|triggered==true`, `property|temperature>25`
 */
func NewConditionProperty(conditionArray []string, thingRefs []string) (*conditionProperty, error) {
	if len(conditionArray) != 2 {
		return nil, errors.New("invalid condition length")
	}
	match := _conditionPropertyRe.FindStringSubmatch(conditionArray[1])
	if match == nil {
		return nil, errors.New("invalid property condition syntax")
	}
	property, operator, ref := match[1], match[2], match[3]

	things := map[string]producer.ExposedThing{}
	for _, thingRef := range thingRefs {
		if exposedThing, ok := _exposedThings[thingRef]; ok && exposedThing != nil {
			if _, err := exposedThing.ExposedProperty(property); err != nil {
				zlog.Error().Err(err).Str("thing", thingRef).Str("property", property).Msg("[automations:NewConditionProperty] Property not found")
				continue // jump to next thing
			}
			things[thingRef] = exposedThing
		} else {
			zlog.Error().Str("thing", thingRef).Str("property", property).Msg("[automations:NewConditionProperty] Thing not found")
		}
	}
	if len(things) == 0 {
		return nil, errors.New("requested things property not found")
	}

	return &conditionProperty{
		property: property,
		operator: operator,
		ref:      ref,
		things:   things,
	}, nil
}

/**
 * Check if at least one thing's property matches the reference value
 */
func (c *conditionProperty) check(time.Time) (bool, error) {
	for thingRef, exposedThing := range c.things {
		exposedProperty, err := exposedThing.ExposedProperty(c.property)
		if err != nil {
			continue
		}
		value, err := exposedProperty.Read(exposedThing, c.property, nil)
		if err != nil {
			zlog.Warn().Err(err).Str("thing", thingRef).Str("property", c.property).Msg("[automations:conditionProperty] Read failed")
			continue
		}
		ok, err := compareValues(value, c.operator, c.ref)
		if err != nil {
			return false, err
		}
		if ok {
			return true, nil
		}
	}
	return false, nil
}

/**
 * Get the properties to observe, by thing ref
 */
func (c *conditionProperty) list() map[string][]string {
	properties := map[string][]string{}
	for thingRef := range c.things {
		properties[thingRef] = []string{c.property}
	}
	return properties
}

func compareValues(value interface{}, operator string, ref string) (bool, error) {
	// Numeric comparison, when both sides parse as numbers
	if floatValue, ok := toFloat(value); ok {
		if floatRef, err := strconv.ParseFloat(ref, 64); err == nil {
			switch operator {
			case "==":
				return floatValue == floatRef, nil
			case "!=":
				return floatValue != floatRef, nil
			case "<":
				return floatValue < floatRef, nil
			case "<=":
				return floatValue <= floatRef, nil
			case ">":
				return floatValue > floatRef, nil
			case ">=":
				return floatValue >= floatRef, nil
			}
		}
		return false, fmt.Errorf("can't compare number with %q", ref)
	}

	if boolValue, ok := value.(bool); ok {
		boolRef, err := strconv.ParseBool(ref)
		if err != nil {
			return false, fmt.Errorf("can't compare boolean with %q", ref)
		}
		switch operator {
		case "==":
			return boolValue == boolRef, nil
		case "!=":
			return boolValue != boolRef, nil
		}
		return false, fmt.Errorf("unsupported operator %q for boolean", operator)
	}

	stringValue := fmt.Sprintf("%v", value)
	switch operator {
	case "==":
		return stringValue == ref, nil
	case "!=":
		return stringValue != ref, nil
	}
	return false, fmt.Errorf("unsupported operator %q for string", operator)
}

func toFloat(value interface{}) (float64, bool) {
	switch number := value.(type) {
	case float64:
		return number, true
	case float32:
		return float64(number), true
	case int:
		return float64(number), true
	case int8:
		return float64(number), true
	case int16:
		return float64(number), true
	case int32:
		return float64(number), true
	case int64:
		return float64(number), true
	case uint:
		return float64(number), true
	case uint8:
		return float64(number), true
	case uint16:
		return float64(number), true
	case uint32:
		return float64(number), true
	case uint64:
		return float64(number), true
	}
	return 0, false
}
