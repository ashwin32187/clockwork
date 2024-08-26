package clockwork

import (
	"fmt"
	"sync"
	"time"
)

// add all tasks here for execution using AddTasks(),
// further calls to AddTasks() will append to existing tasks list,
// tasks can also be scheduled by using AddTask(),
// stop and remove tasks from execution list using ClearTasks(), RemoveTask(), RemoveTaskByName()
// view status of current tasks list using Report(), keep track of tasks using FetchTasks()
// stop specific tasks by calling Stop()
// see stop field for detail changes to accessible fields on any task after being picked up
// for execution are not supported by package internal copies of accessible fields are made
// when AddTasks (), AddTask () or Reinstate ( is called and these will be used for processing

var tasksForExecution []*Task

type Task struct {
	// task name
	Name string // keeping task names unique helps with tracking
	name string
	// function should take string task name, int run count as arguments
	// this can be used to build more granular controls like having an hourly schedule
	// but running desired operations within task only once every 2 hours by skipping for odd or even run counts
	// first run will always be initiated as close to the AddTasks() or run() call as possible and run count will be passed as 1
	Function func(string, int)
	function func(string, int)
	// for default supported schedule types use SupportedScheduleTypes function call
	// custom schedules can be implemented by implementing Schedule interface
	Schedule Schedule
	schedule Schedule
	// Cron can be used like cron expression but using time Time, time fields will be parsed Cron time.Time used as per schedule, see schedules for details
	Cron time.Time
	cron time.Time
	// for internal use
	// stop specific tasks using Stop(), once stop is initiated ongoing execution will still complete but no further runs will occur
	// the task can be rescheduled by creating a new task with the same fields to the execution list and calling AddTask()
	// the task can also be rescheduled by calling Reinstate() for the stopped task, see Reinstate() for details
	stopped   bool // true only when Stop() or RemoveTask(), RemoveTaskByName() is called for the task, ClearTasks() sets this to true for all tasks
	running   bool // true only when being picked up for execution, false when not started or after Stop() is called
	firstRun  time.Time
	latestRun time.Time
	nextRun   time.Time
	once      sync.Once
	// to be_processed internally and sent to task function which can be used to skip occurrences
	runCount int
}

// returns execution list currently being tracked by package
func FetchTasks() []*Task {
	if len(tasksForExecution) == 0 {
		return nil
	}
	destination := make([]*Task, len(tasksForExecution))
	copy(destination, tasksForExecution)
	return destination
}

// add, clear and remove methods use mutex lock as execution list is modified
var lock sync.RWMutex

// run tasks and them to execution list
func AddTasks(tasks []*Task) int {
	lock.Lock()
	defer lock.Unlock()
	if len(tasksForExecution) == 0 {
		tasksForExecution = make([]*Task, 0)
	}
	for index, task := range tasks {
		if task.isValid() && !task.running && !task.stopped {
			// make immutable copies for run reference
			task.name = task.Name
			task.function = task.Function
			task.schedule = task.Schedule
			task.cron = task.Cron

			func(taskIndex int) {
				go (task).run()
				tasksForExecution = append(tasksForExecution, task)
			}(index)
		}
	}
	return len(tasksForExecution)
}

// run a new task and add it to the execution list for reporting
// will return true if task execution has been started and added to the list
func (t *Task) AddTask() (int, bool) {
	lock.Lock()
	defer lock.Unlock()
	if t.isValid() && !t.running && !t.stopped {
		if len(tasksForExecution) == 0 {
			tasksForExecution = make([]*Task, 0)
		}
		// make immutable copies for run reference
		t.name = t.Name
		t.function = t.Function
		t.schedule = t.Schedule
		t.cron = t.Cron
		go t.run()
		tasksForExecution = append(tasksForExecution, t)
		return len(tasksForExecution), true
	}
	return len(tasksForExecution), false // invalid task, return length of existing tasks list, false to indicate no action has been taken
}

// core task runner always needs to be invoked in separate goroutine as it does not return until task is stopped
func (t *Task) run() {
	if t == nil || t.function == nil || t.schedule == nil {
		return
	}
	runTime := t.computeFirstRunTime()
	for {
		// release go routine before wait if task is nil or stopped
		if t.stopped {
			t.running = false
			return
		}
		t.running = true
		time.Sleep(time.Until(runTime))
		// double checking after wait
		if t.stopped {
			t.running = false
			return
		}
		func(name string, count int) {
			go t.function(name, count)
		}(t.name, t.runCount+1)
		t.runCount = t.runCount + 1
		t.latestRun = runTime
		runTime = t.schedule.NextRun(runTime)
		t.nextRun = runTime
	}
}

func (t *Task) computeFirstRunTime() time.Time {
	runTime := t.schedule.FirstRun(t.cron)
	for {
		if runTime.After(time.Now()) {
			break
		}
		runTime = t.schedule.NextRun(runTime)
	}
	if t.firstRun.Year() == 1 {
		t.firstRun = runTime
	}
	return runTime
}

// can be used to stopped a task, once called, the task cannot  be run again, ongoing execution will still complete
// the stopped task can be rescheduled by calling Reinstate() on a stopped task
func (t *Task) Stop() {
	if t != nil {
		t.stopped = true
	}
}

// can be used to recreate an existing stopped task (only once per stopped task) with fresh stats
// if called on an running task, no action is taken and existing task is returned with false
// a new task is created and execution is started when Reinstate() is called on a stopped task
// pointer to the new task and true are returned
func (t *Task) Reinstate() (*Task, bool) {
	if t.isValid() && t.stopped {
		var restarted Task
		t.once.Do(func() {
			t.reinstateOnce(&restarted)
		})
		if restarted.schedule != nil && restarted.function != nil {
			// do once invoked reinstated once
			return &restarted, true
		}
	}
	return t, false
}

// to ensure a stopped thread can only be truly reinstated once
func (t *Task) reinstateOnce(restarted *Task) {
	restarted.Name = t.name
	restarted.Function = t.function
	restarted.Schedule = t.schedule
	restarted.Cron = t.cron
	restarted.AddTask()
}

// stop and clear current execution list, Report() will not track tasks after ClearTasks() is called
func ClearTasks() {
	lock.Lock()
	defer lock.Unlock()
	for index := range tasksForExecution {
		tasksForExecution[index].stopped = true
	}
	tasksForExecution = nil
}

// stop and remove task while preserving order, removed tasks will no longer be covered by Report()
func (t *Task) RemoveTask() bool {
	lock.Lock()
	defer lock.Unlock()
	for index, task := range tasksForExecution {
		if task == t {
			task.stopped = true
			copy(tasksForExecution[index:], tasksForExecution[index+1:]) // shift a[i+l:] left one index
			tasksForExecution[len(tasksForExecution)-1] = nil
			tasksForExecution = tasksForExecution[:len(tasksForExecution)-1] // erase last element (write zero value)
			return true
		}
	}
	return false
}

// stop and remove task by name (first match) while preserving order, removed tasks will no longer be covered by Report()
func RemoveTaskByName(name string) bool {
	lock.Lock()
	defer lock.Unlock()
	for index, task := range tasksForExecution {
		if task.name == name {
			task.stopped = true
			copy(tasksForExecution[index:], tasksForExecution[index+1:]) // shift a[i+l:] left one index
			tasksForExecution[len(tasksForExecution)-1] = nil
			tasksForExecution = tasksForExecution[:len(tasksForExecution)-1] // erase last element (write zero value)
			return true
		}
	}
	return false
}

// check task for mandatory values
func (t *Task) isValid() bool {
	return t != nil && t.Schedule != nil && t.Function != nil
}

// can be used to create custom schedules that can be used with clockwork
type Schedule interface {
	// calculate first run time for the schedule variants
	FirstRun(time.Time) time.Time
	// calculate subsequent run time for schedule variants based on last run time
	NextRun(time.Time) time.Time
}

// supported schedule implementations
type ScheduleType struct {
	string
}

var EverySecond = ScheduleType{"SECOND"}
var EveryMinute = ScheduleType{"MINUTE"}
var Hourly = ScheduleType{"HOURLY"}
var Daily = ScheduleType{"DAILY"}
var Weekly = ScheduleType{"WEEKLY"}
var Monthly = ScheduleType{"Monthly"}
var Yearly = ScheduleType{"YEARLY"}

// used to calculate subsequent run times for tasks as per schedule based on last run time after the first run
func (s *ScheduleType) NextRun(currentRun time.Time) time.Time {
	var timeOfRun time.Time
	switch *s {
	case EverySecond:
		timeOfRun = currentRun.Add(time.Second)
	case EveryMinute:
		timeOfRun = currentRun.Add(time.Minute)
	case Hourly:
		timeOfRun = currentRun.Add(time.Hour)
	case Daily:
		timeOfRun = currentRun.AddDate(0, 0, 1)
	case Weekly:
		timeOfRun = currentRun.AddDate(0, 0, 7)
	case Monthly:
		timeOfRun = currentRun.AddDate(0, 1, 0)
	case Yearly:
		timeOfRun = currentRun.AddDate(1, 0, 0)
	default:
		panic("unsupported schedule type")
	}
	return timeOfRun
}

// used to calculate subsequent run times for tasks as per schedule based on last run time after the first run
func (s *ScheduleType) FirstRun(cron time.Time) time.Time {
	timeOfRun := time.Now().In(cron.Location())
	switch *s {
	case EverySecond: // for scheduling every second, cron will not be sed and can be zero value
		timeOfRun = timeOfRun.Add(time.Second)

	case EveryMinute: // for scheduling every minute, second will be used
		timeOfRun = time.Date(timeOfRun.Year(), timeOfRun.Month(), timeOfRun.Day(), timeOfRun.Hour(), timeOfRun.Minute(), cron.Second(), 0, timeOfRun.Location())

	case Hourly: // for hourly schedule, location, minute and second will be used
		timeOfRun = time.Date(timeOfRun.Year(), timeOfRun.Month(), timeOfRun.Day(), timeOfRun.Hour(), cron.Minute(), cron.Second(), 0, timeOfRun.Location())

	case Daily: // for daily schedule, location, hour, minute and second will be used
		timeOfRun = time.Date(timeOfRun.Year(), timeOfRun.Month(), timeOfRun.Day(), cron.Hour(), cron.Minute(), cron.Second(), 0, timeOfRun.Location())

	case Weekly: // for weekly schedule, day of week, location, hour, minute and second will be used
		year, month, day := GetNearestDate(cron.Weekday(), cron.Location())
		timeOfRun = time.Date(year, month, day, cron.Hour(), cron.Minute(), cron.Second(), 0, timeOfRun.Location())

	case Monthly: // for monthly schedule, date, location, hour, minute and second will be used
		timeOfRun = time.Date(timeOfRun.Year(), timeOfRun.Month(), cron.Day(), cron.Hour(), cron.Minute(), cron.Second(), 0, timeOfRun.Location())

	case Yearly: // for yearly schedule, month, date, location, hour, minute and second will be used
		timeOfRun = time.Date(timeOfRun.Year(), cron.Month(), cron.Day(), cron.Hour(), cron.Minute(), cron.Second(), 0, timeOfRun.Location())

	default:
		panic("unsupported schedule type")
	}
	for { // ensure future time for first run
		if timeOfRun.After(time.Now()) {
			break
		}
		timeOfRun = s.NextRun(timeOfRun)
	}
	return timeOfRun
}

// get the nearest future date for any weekday (Sunday to Saturday)
// when creating a weekly schedule, the year, month and day returned can be used with time.Date(...)
// along with the right timestamp to feed into the Cron param
func GetNearestDate(dayOfWeek time.Weekday, location *time.Location) (int, time.Month, int) {
	if location == nil {
		location = time.Local
	}
	nearestDate := time.Now().In(location)
	today := nearestDate.Weekday()
	if today <= dayOfWeek {
		nearestDate = nearestDate.AddDate(0, 0, int(dayOfWeek-today))
	} else {
		nearestDate = nearestDate.AddDate(0, 0, int(dayOfWeek+7-today))
	}
	return nearestDate.Date()
}

func SupportedScheduleTypes() []ScheduleType {
	return []ScheduleType{EverySecond, EveryMinute, Hourly, Daily, Weekly, Monthly, Yearly}
}

// used only to create a report on current status of tasks being tracked
type TaskReport struct {
	Name      string
	Running   bool
	RunCount  int
	FirstRun  time.Time
	LatestRun time.Time
	NextRun   time.Time
}

func (t TaskReport) String() string {
	var first, latest, next string
	if t.FirstRun.Year() != 1 {
		first = t.FirstRun.String()
	}
	if t.LatestRun.Year() != 1 {
		latest = t.LatestRun.String()
	}
	if t.Running && t.NextRun.Year() != 1 {
		next = t.NextRun.String()
	}
	return fmt.Sprintf("{\"name\": \"%s\", \"running\": %v, \"runCount\": %d, \"firstRun\": \"%s\", \"latestRun\": \"%s\", \"nextRun\": \"%s\"}", t.Name, t.Running, t.RunCount, first, latest, next)
}

// used to get current stats for task execution on demand
func Report() []TaskReport {
	report := make([]TaskReport, 0, len(tasksForExecution))
	for _, task := range tasksForExecution {
		if task != nil {
			taskReport := TaskReport{
				Name:      task.name,
				Running:   task.running,
				RunCount:  task.runCount,
				FirstRun:  task.firstRun,
				LatestRun: task.latestRun,
				NextRun:   task.nextRun,
			}
			report = append(report, taskReport)
		}
	}
	return report
}

// scheduler to handle monthly schedules on a set date of the month
// for month end dates past 28th of the month schedule will run at the end of the month
// without rolling over to next month
// use weekday adjustment and postpone flags to run on the nearest weekday
// in case day of month to run the schedule falls on weekend
// this is also an illustration of implementing custom schedules as necessary
type DayOfMonthSchedule struct {
	// will default to last day of month if invalid date is set
	DayOfMonth        int
	WeekdayAdjustment bool
	Postpone          bool
	// will default to Saturday and Sunday for the weekend if incorrect configuration is set
	weekend []time.Weekday
}

// allows for configurable weekend days, returns status if weekend config was updated correctly
func (s *DayOfMonthSchedule) SetWeekend(weekend []time.Weekday) bool {
	// ensure scheduler can run at least one day of the week
	// will not handle duplicate days in array leading to size >= 7
	if len(weekend) >= 7 {
		return false
	}
	s.weekend = weekend
	return true
}

func (s *DayOfMonthSchedule) FirstRun(cron time.Time) time.Time {
	currentDate := time.Now().In(cron.Location())
	// for this schedule, location, hour, minute and second will be used from task cron
	timeOfRun := time.Date(currentDate.Year(), currentDate.Month(), 1, cron.Hour(), cron.Minute(), cron.Second(), 0, cron.Location())
	timeOfRun = s.adjust(timeOfRun)
	for {
		if timeOfRun.Before(currentDate) {
			// cannot run for this month as time of run has already crossed
			// reset day to 1
			timeOfRun = time.Date(timeOfRun.Year(), timeOfRun.Month(), 1, cron.Hour(), cron.Minute(), cron.Second(), 0, timeOfRun.Location())
			// add a month
			timeOfRun = timeOfRun.AddDate(0, 1, 0)
			timeOfRun = s.adjust(timeOfRun)
		} else {
			break
		}
	}
	return timeOfRun
}

func (s *DayOfMonthSchedule) adjust(timeOfRun time.Time) time.Time {
	// go to last day of month
	timeOfRun = timeOfRun.AddDate(0, 1, -1)
	if s.DayOfMonth > 0 && s.DayOfMonth < timeOfRun.Day() {
		// run on either the last day or the day of month in schedule as per current month
		timeOfRun = time.Date(timeOfRun.Year(), timeOfRun.Month(), s.DayOfMonth, timeOfRun.Hour(), timeOfRun.Minute(), timeOfRun.Second(), 0, timeOfRun.Location())
	}
	// adjust to run on a weekday as per configuration in schedule
	if s.WeekdayAdjustment {
		if s.weekend == nil {
			s.weekend = []time.Weekday{time.Saturday, time.Sunday}
		}
		dateAdjustment := -1
		if s.Postpone {
			dateAdjustment = 1
		}
		for {
			isWeekend := false
			for _, weekday := range s.weekend {
				if timeOfRun.Weekday() == weekday {
					isWeekend = true
					break
				}
			}
			if isWeekend {
				timeOfRun = timeOfRun.AddDate(0, 0, dateAdjustment)
			} else {
				break
			}
		}
	}
	return timeOfRun
}

func (s *DayOfMonthSchedule) NextRun(cron time.Time) time.Time {
	previousRunDate := cron
	// reset day to 1
	cron = time.Date(cron.Year(), cron.Month(), 1, cron.Hour(), cron.Minute(), cron.Second(), 0, cron.Location())
	monthsToAdd := 1
	if s.WeekdayAdjustment {
		if s.Postpone && s.DayOfMonth > 20 && previousRunDate.Day() < 7 {
			// adjust for when weekday adjusment postpones run into next month
			monthsToAdd = 0
		} else if !s.Postpone && previousRunDate.Day() > s.DayOfMonth {
			// adjust for when weekday adjusment takes previous run date to the month before
			monthsToAdd = 2
		}
	}
	cron = cron.AddDate(0, monthsToAdd, 0)
	return s.adjust(cron)
}
