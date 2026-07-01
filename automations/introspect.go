package automations

import (
	"sort"
	"time"
)

// ScheduledAutomation describes one currently-planned automation job, as
// returned by the EriaAppController "automations" endpoint. There is one entry
// per active group of each automation (a group with no matching/enabled
// schedule produces no entry).
type ScheduledAutomation struct {
	// Automation is the name of the owning automation (from the YAML config).
	Automation string `json:"automation"`
	// Devices are the thing references the action is run against.
	Devices []string `json:"devices"`
	// Action is the action key invoked on the devices (e.g. "on", "toggle").
	Action string `json:"action"`
	// Value is the optional action value, when the action carries one.
	Value interface{} `json:"value,omitempty"`
	// Schedule is a human-readable description of the schedule
	// (e.g. "every day at 07:30:00", "immediately").
	Schedule string `json:"schedule"`
	// NextRun is the next execution date/time, for time-based schedules. It is
	// nil for schedules that don't have a future run (immediate/disabled).
	NextRun *time.Time `json:"nextRun,omitempty"`
}

// Scheduled returns the list of automation jobs currently planned. It reflects
// the live scheduling state: entries change as contexts change and hour-based
// schedules are re-planned.
func Scheduled() []ScheduledAutomation {
	_mu.Lock()
	defer _mu.Unlock()

	result := make([]ScheduledAutomation, 0)
	for _, automation := range _automations {
		for _, job := range automation.jobs {
			entry := ScheduledAutomation{
				Automation: automation.name,
				Devices:    []string{},
				Schedule:   job.schedule.string(),
			}
			if a, ok := job.action.(*Action); ok {
				entry.Action = a.Ref
				entry.Value = a.Value
				entry.Devices = deviceRefs(a)
			}
			if next, ok := scheduleNextRun(job.schedule); ok {
				entry.NextRun = &next
			}
			result = append(result, entry)
		}
	}
	return result
}

// deviceRefs returns the (sorted, for stable output) thing references an action
// targets.
func deviceRefs(a *Action) []string {
	refs := make([]string, 0, len(a.ExposedThings))
	for ref := range a.ExposedThings {
		refs = append(refs, ref)
	}
	sort.Strings(refs)
	return refs
}

// scheduleNextRun returns the next execution time of a schedule when it has one
// (hour-based schedules backed by a cron job). Immediate and disabled schedules
// return false.
func scheduleNextRun(s Schedule) (time.Time, bool) {
	if at, ok := s.(*scheduleAtHour); ok && at.cronJob != nil {
		if next, err := at.cronJob.NextRun(); err == nil {
			return next, true
		}
	}
	return time.Time{}, false
}
