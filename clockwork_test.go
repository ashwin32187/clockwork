package clockwork

import (
	"fmt"
	"testing"
	"time"
)

func Test_clockwork(_ *testing.T) {
	testDayOfMonthScheduler()
	testRunner()
}
func testRunner() {
	yearlyAfter, yearlyBefore, yearlyLocal, yearly, monthlyAfter, monthlyBefore, monthly, weeklyAfter, weeklyBefore, weekly,
		dailyAfter, dailyBefore, daily, hourly, everyMinute, everySecond, dayOfTheMonth := testData()
	fmt.Println("tasks added", AddTasks([]*Task{yearlyAfter, yearlyBefore, yearlyLocal, yearly, monthlyAfter, monthlyBefore, monthly, weeklyAfter, weeklyBefore, weekly,
		dailyAfter, dailyBefore, daily, hourly, everyMinute, everySecond, dayOfTheMonth}))
	time.Sleep(5 * time.Second)
	// cannot change tasks in execution list from outside package
	FetchTasks()[0] = nil
	// still has access to task post nil assigment
	FetchTasks()[0].Stop()
	fmt.Println(FetchTasks()[0].Name)
	time.Sleep(1 * time.Second)
	fmt.Println(Report())
	// remove a task
	FetchTasks()[0].RemoveTask()
	fmt.Println(Report())
	newEverySecond := &Task{
		Name:     "new every second",
		Function: Work,
		Schedule: &EverySecond,
		Cron:     time.Now(),
	}
	// add a single task
	fmt.Println(newEverySecond.AddTask())
	time.Sleep(3 * time.Second)
	// stop a task
	newEverySecond.Stop()
	everySecond.Stop()
	for count := 0; count < 3; count++ {
		fmt.Println(Report())
		time.Sleep(1 * time.Second)
		everyMinute.Stop()
	}
	fmt.Println("test RemoveTask()")
	fmt.Println(newEverySecond.RemoveTask())
	fmt.Println(newEverySecond.RemoveTask())
	fmt.Println(Report())
	fmt.Println("test Reinstate()")
	everyMinute, flag := everyMinute.Reinstate()
	fmt.Println(flag)
	everyMinute, flag = everyMinute.Reinstate()
	fmt.Println(flag)
	time.Sleep(5 * time.Second)
	fmt.Println(Report())
	time.Sleep(1 * time.Minute)
	fmt.Println(Report())
	everyMinute.Stop()
	RemoveTaskByName(everyMinute.Name)
	time.Sleep(5 * time.Second)
	time.Sleep(1 * time.Minute)
	fmt.Println(Report())
	ClearTasks()
	time.Sleep(2 * time.Second)
	everySecond.Reinstate()
	time.Sleep(2 * time.Second)
	fmt.Println(Report())
	ClearTasks()
}

func Work(name string, count int) {
	fmt.Printf("running %s, count %d, at %s \n", name, count, time.Now().String())
}

func testData() (*Task, *Task, *Task, *Task,
	*Task, *Task, *Task, *Task,
	*Task, *Task, *Task, *Task,
	*Task, *Task, *Task, *Task, *Task) {
	refTime := time.Now()
	return &Task{
			Name:     "yearlyAfter",
			Function: Work,
			Schedule: &Yearly,
			Cron:     time.Date(refTime.Year(), refTime.Month(), refTime.Day()+1, 00, 00, 00, 0, time.Local),
		}, &Task{
			Name:     "yearlyBefore",
			Function: Work,
			Schedule: &Yearly,
			Cron:     time.Date(refTime.Year(), refTime.Month(), refTime.Day()-1, 00, 00, 00, 0, time.Local),
		}, &Task{
			Name:     "yearlyLocal",
			Function: Work,
			Schedule: &Yearly,
			Cron:     time.Date(refTime.Year(), refTime.Month(), refTime.Day(), refTime.Hour()-1, 11, 00, 0, time.Local),
		},
		&Task{
			Name:     "yearly",
			Function: Work,
			Schedule: &Yearly,
			Cron:     time.Date(refTime.Year(), refTime.Month(), refTime.Day(), refTime.Hour()-5, refTime.Minute()-28, refTime.Second(), 0, time.UTC),
		},
		&Task{
			Name:     "monthlyAfter",
			Function: Work,
			Schedule: &Monthly,
			Cron:     time.Date(refTime.Year(), refTime.Month(), refTime.Day()+1, 07, 07, 07, 0, time.UTC),
		},
		&Task{
			Name:     "monthlyBefore",
			Function: Work,
			Schedule: &Monthly,
			Cron:     time.Date(refTime.Year(), refTime.Month(), refTime.Day()-1, 17, 30, 00, 0, time.UTC),
		},
		&Task{
			Name:     "monthly",
			Function: Work,
			Schedule: &Monthly,
			Cron:     time.Date(refTime.Year(), refTime.Month(), refTime.Day(), refTime.Hour()-5, refTime.Minute()-28, refTime.Second(), 0, time.UTC),
		},
		&Task{
			Name:     "weeklyAfter",
			Function: Work,
			Schedule: &Weekly,
			Cron:     time.Date(refTime.Year(), refTime.Month(), refTime.Day()+1, 07, 14, 21, 0, time.UTC),
		},
		&Task{
			Name:     "weeklyBefore",
			Function: Work,
			Schedule: &Weekly,
			Cron:     time.Date(refTime.Year(), refTime.Month(), refTime.Day()-2, 01, 11, 48, 0, time.UTC),
		},
		&Task{
			Name:     "weekly",
			Function: Work,
			Schedule: &Weekly,
			Cron:     time.Date(refTime.Year(), refTime.Month(), refTime.Day(), refTime.Hour()-5, refTime.Minute()-28, 00, 0, time.UTC),
		},
		&Task{
			Name:     "dailyAfter",
			Function: Work,
			Schedule: &Daily,
			Cron:     time.Date(refTime.Year(), refTime.Month(), refTime.Day(), 20, 56, 07, 0, time.UTC),
		},
		&Task{
			Name:     "dailyBefore",
			Function: Work,
			Schedule: &Daily,
			Cron:     time.Date(refTime.Year(), refTime.Month(), refTime.Day(), 1, 17, 59, 0, time.UTC),
		},
		&Task{
			Name:     "daily",
			Function: Work,
			Schedule: &Daily,
			Cron:     time.Date(refTime.Year(), refTime.Month(), refTime.Day(), refTime.Hour()-5, refTime.Minute()-28, 00, 0, time.UTC),
		},
		&Task{
			Name:     "hourly",
			Function: Work,
			Schedule: &Hourly,
			Cron:     time.Date(refTime.Year(), refTime.Month(), refTime.Day(), refTime.Hour()-5, refTime.Minute()-28, 00, 0, time.UTC),
		},
		&Task{
			Name:     "everyMinute",
			Function: Work,
			Schedule: &EveryMinute,
			Cron:     time.Now(),
		},
		&Task{
			Name:     "everySecond",
			Function: Work,
			Schedule: &EverySecond,
			Cron:     time.Now(),
		},
		&Task{
			Name:     "dayOfTheMonth",
			Function: Work,
			Schedule: &DayOfMonthSchedule{
				DayOfMonth: 31,
			},
			Cron: time.Now(),
		}
}

func testDayOfMonthScheduler() {
	schedule := DayOfMonthSchedule{
		DayOfMonth:        31,
		WeekdayAdjustment: false,
		Postpone:          false,
	}

	runTime := schedule.FirstRun(time.Now())
	fmt.Println("End of month")
	fmt.Println(runTime)
	for month := 1; month < 13; month++ {
		runTime = schedule.NextRun(runTime)
		fmt.Println(runTime)
	}
	fmt.Println("30th of month")
	schedule.DayOfMonth = 30
	runTime = schedule.FirstRun(time.Now())
	fmt.Println(runTime)
	for month := 1; month < 13; month++ {
		runTime = schedule.NextRun(runTime)
		fmt.Println(runTime)
	}
	fmt.Println("29th of month")
	schedule.DayOfMonth = 29
	runTime = schedule.FirstRun(time.Now())
	fmt.Println(runTime)
	for month := 1; month < 13; month++ {
		runTime = schedule.NextRun(runTime)
		fmt.Println(runTime)
	}
	fmt.Println("28th of month")
	schedule.DayOfMonth = 28
	runTime = schedule.FirstRun(time.Now())
	fmt.Println(runTime)
	for month := 1; month < 13; month++ {
		runTime = schedule.NextRun(runTime)
		fmt.Println(runTime)
	}
	fmt.Println("10th of month")
	schedule.DayOfMonth = 10
	runTime = schedule.FirstRun(time.Now())
	fmt.Println(runTime)
	for month := 0; month < 12; month++ {
		runTime = schedule.NextRun(runTime)
		fmt.Println(runTime)
	}
	fmt.Println("invalid config")
	schedule.DayOfMonth = -1
	runTime = schedule.FirstRun(time.Now())
	fmt.Println(runTime)
	for month := 1; month < 13; month++ {
		runTime = schedule.NextRun(runTime)
		fmt.Println(runTime)
	}
	fmt.Println("adjust by bringing forward")
	schedule.DayOfMonth = 1
	schedule.WeekdayAdjustment = true
	runTime = schedule.FirstRun(time.Now())
	fmt.Println(runTime)
	for month := 1; month < 13; month++ {
		runTime = schedule.NextRun(runTime)
		fmt.Println(runTime)
	}
	fmt.Println("adjust by postponing")
	schedule.DayOfMonth = 31
	schedule.WeekdayAdjustment = true
	schedule.Postpone = true
	runTime = schedule.FirstRun(time.Now())
	fmt.Println(runTime)
	for month := 1; month < 13; month++ {
		runTime = schedule.NextRun(runTime)
		fmt.Println(runTime)
	}
	fmt.Println("leap year test")
	schedule.WeekdayAdjustment = false
	schedule.DayOfMonth = 29
	runTime = schedule.FirstRun(time.Now())
	fmt.Println(runTime)
	for month := 1; month < 48; month++ {
		runTime = schedule.NextRun(runTime)
		if runTime.Year()%4 == 0 || runTime.Year()%4 == 3 {
			fmt.Println(runTime)
		}
	}
}
