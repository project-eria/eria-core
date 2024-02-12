package automations

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/gookit/goutil/arrutil"
)

type conditionContexts []conditionContext

type conditionContext struct {
	context string
	invert  bool
}

/**
 * Matches at least one context (OR) condition
 * `context|<context name>|<context name>`
 *	- not: `context|!<context name>|!<context name>`
 */
func NewConditionContext(conditionArray []string) (*conditionContexts, error) {
	contextsThing := _consumer.ThingFromTag("contexts")
	if contextsThing == nil {
		return nil, errors.New("contexts thing not configured")
	}

	// Check if the condition has the correct number of parameters
	if len(conditionArray) < 2 {
		return nil, errors.New("invalid condition length")
	}

	var contexts conditionContexts
	re, _ := regexp.Compile(`^!?\w*$`)
	for i := 1; i < len(conditionArray); i++ {
		c := conditionArray[i]

		// Check if the context name is valid
		match := re.MatchString(c)
		if !match {
			return nil, errors.New("invalid context name")
		}

		if strings.HasPrefix(c, "!") {
			// If the context is NOT active
			contexts = append(contexts, conditionContext{
				context: c[1:],
				invert:  true,
			})
		} else {
			contexts = append(contexts, conditionContext{
				context: c,
				invert:  false,
			})
		}
	}
	return &contexts, nil
}

/**
 * Check if one of the contexts matches the active ones
 */
func (conditions *conditionContexts) check(time.Time) (bool, error) {
	for _, condition := range *conditions {
		// Check if the context is on the actives list
		active := arrutil.InStrings(condition.context, _activeContexts)
		if condition.invert != active {
			// The condition matches
			return true, nil
		}
	}
	// No match found
	return false, nil
}

/**
 * Get the list of contexts, to observe
 */
func (conditions *conditionContexts) list() []string {
	return arrutil.Map(*conditions, func(c conditionContext) (string, bool) { return c.context, true })
}
