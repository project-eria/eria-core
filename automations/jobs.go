package automations

import (
	"time"

	"github.com/gookit/goutil/arrutil"
	zlog "github.com/rs/zerolog/log"
)

// scheduleJobs builds and schedules every configured automation.
func scheduleJobs(automations []AutomationConfig) {
	now := time.Now().In(_location)

	for _, automationConfig := range automations {
		automation, observables, err := getAutomation(automationConfig)
		if err != nil {
			zlog.Warn().Err(err).Str("automation", automationConfig.Name).Msg("[automations:scheduleJobs] Skipping...")
			continue
		}
		automation.scheduleJob(now)
		// Register the automation in the global registry, for introspection.
		_mu.Lock()
		_automations = append(_automations, automation)
		_mu.Unlock()
		// Register the automation against the contexts it observes, so it can be
		// re-planned when one of them changes.
		for _, context := range observables.contexts {
			if _, found := _contextsAutomations[context]; !found {
				_contextsAutomations[context] = make([]*Automation, 0)
			}
			_contextsAutomations[context] = append(_contextsAutomations[context], automation)
		}
	}
}

// getAutomation prepares the runtime Automation from its configuration, and
// collects the contexts referenced by its conditions (for monitoring).
func getAutomation(c AutomationConfig) (*Automation, *Observables, error) {
	zlog.Info().Str("automation", c.Name).Msg("[automations:getAutomation] Adding automation")

	observables := &Observables{contexts: make([]string, 0)}
	groups := make([]actionGroup, len(c.Groups))
	for gi, groupConfig := range c.Groups {
		// One action per group
		action, err := getAction(c.Things, c.Name, groupConfig.Action)
		if err != nil {
			zlog.Error().Err(err).Str("automation", c.Name).Str("action", groupConfig.Action).Msg("[automations:getAutomation] Can't get action")
			return nil, nil, err // Skip this automation
		}

		// First-match list of schedule branches
		entries := make([]scheduleEntry, len(groupConfig.Schedules))
		for ei, entryConfig := range groupConfig.Schedules {
			conditions, o, err := getConditions(entryConfig.Conditions)
			if err != nil {
				zlog.Error().Err(err).Str("automation", c.Name).Int("group", gi).Int("entry", ei).Msg("[automations:getAutomation] Can't get conditions")
				return nil, observables, err // Skip this automation
			}
			mergeContexts(observables, o)

			var schedule Schedule
			if entryConfig.Schedule != "" {
				schedule, err = getSchedule(entryConfig.Schedule)
				if err != nil {
					zlog.Error().Err(err).Str("automation", c.Name).Int("group", gi).Int("entry", ei).Str("schedule", entryConfig.Schedule).Msg("[automations:getAutomation] Can't get schedule")
					return nil, nil, err // Skip this automation
				}
			}
			entries[ei] = scheduleEntry{conditions: conditions, schedule: schedule}
		}
		groups[gi] = actionGroup{action: action, entries: entries}
	}

	return &Automation{name: c.Name, groups: groups}, observables, nil
}

// plan computes the set of actions to schedule now: for each group, the first
// schedule entry whose conditions all match wins; an entry without a schedule
// disables the group.
func (automation *Automation) plan(now time.Time) []scheduledAction {
	var planned []scheduledAction
	for gi, group := range automation.groups {
		for _, entry := range group.entries {
			ok, err := checkConditions(entry.conditions, now)
			if err != nil {
				zlog.Error().Err(err).Str("automation", automation.name).Int("group", gi).Msg("[automations:plan]")
				continue
			}
			if !ok {
				continue
			}
			// First matching branch wins.
			if entry.schedule != nil {
				if err := entry.schedule.job(); err != nil {
					zlog.Error().Err(err).Str("automation", automation.name).Int("group", gi).Msg("[automations:plan]")
				} else {
					planned = append(planned, scheduledAction{schedule: entry.schedule, action: group.action})
				}
			}
			break
		}
	}
	return planned
}

// scheduleJob (re)plans the automation: it cancels the currently active jobs and
// starts the newly planned ones. Conditions are evaluated here (planning time),
// so the active job set always reflects exactly what will run.
func (automation *Automation) scheduleJob(now time.Time) {
	zlog.Trace().Str("automation", automation.name).Msg("[automations:scheduleJob] scheduling jobs")
	planned := automation.plan(now)

	// Snapshot the currently active jobs, then cancel them outside the lock
	// (cancel() may do network I/O via unobserve).
	_mu.Lock()
	previous := automation.jobs
	_mu.Unlock()
	for _, job := range previous {
		job.schedule.cancel()
	}

	// Build the new job set (start() may do network I/O via observe).
	status := "success"
	newJobs := make([]scheduledAction, 0, len(planned))
	for _, job := range planned {
		if err := job.schedule.start(job.action); err != nil {
			zlog.Error().Err(err).Str("automation", automation.name).Msg("[automations:scheduleJob]")
			status = err.Error()
			continue
		}
		newJobs = append(newJobs, job)
		zlog.Info().Str("automation", automation.name).Msgf("[automations:scheduleJob] Job scheduled: %s", job.schedule.string())
	}

	// Install the new plan atomically for readers (Scheduled()).
	_mu.Lock()
	automation.jobs = newJobs
	automation.status = status
	automation.lastScheduled = now
	_mu.Unlock()
}

// mergeContexts appends the contexts observed by o into observables, skipping
// duplicates.
func mergeContexts(observables *Observables, o *Observables) {
	if o == nil {
		return
	}
	for _, context := range o.contexts {
		if !arrutil.StringsHas(observables.contexts, context) {
			observables.contexts = append(observables.contexts, context)
		}
	}
}
