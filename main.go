package main

import (
	"clockwork/clockwork"
	"fmt"
	"time"
)

func main() {
	testDayOfMonthScheduler()
	testRunner()
}
func testRunner() {
	yearlyAfter, yearlyBefore, yearlyLocal, yearly, monthlyAfter, monthlyBefore, monthly, weeklyAfter, weeklyBefore, weekly,
		dailyAfter, dailyBefore, daily, hourly, everyMinute, everySecond, dayOfTheMonth := testData()
	fmt.Println("tasks added", clockwork.AddTasks([]*clockwork.Task{yearlyAfter, yearlyBefore, yearlyLocal, yearly, monthlyAfter, monthlyBefore, monthly, weeklyAfter, weeklyBefore, weekly,
		dailyAfter, dailyBefore, daily, hourly, everyMinute, everySecond, dayOfTheMonth}))
	time.Sleep(5 * time.Second)
	// cannot change tasks in execution list from outside package
	clockwork.FetchTasks()[0] = nil
	// still has access to task post nil assigment
	clockwork.FetchTasks()[0].Stop()
	fmt.Println(clockwork.FetchTasks()[0].Name)
	time.Sleep(1 * time.Second)
	fmt.Println(clockwork.Report())
	// remove a task
	clockwork.FetchTasks()[0].RemoveTask()
	fmt.Println(clockwork.Report())
	newEverySecond := &clockwork.Task{
		Name:     "new every second",
		Function: Work,
		Schedule: &clockwork.EverySecond,
		Cron:     time.Now(),
	}
	// add a single task
	fmt.Println(newEverySecond.AddTask())
	time.Sleep(3 * time.Second)
	// stop a task
	newEverySecond.Stop()
	everySecond.Stop()
	for count := 0; count < 3; count++ {
		fmt.Println(clockwork.Report())
		time.Sleep(1 * time.Second)
		everyMinute.Stop()
	}
	fmt.Println("test RemoveTask()")
	fmt.Println(newEverySecond.RemoveTask())
	fmt.Println(newEverySecond.RemoveTask())
	fmt.Println(clockwork.Report())
	fmt.Println("test Reinstate()")
	everyMinute, flag := everyMinute.Reinstate()
	fmt.Println(flag)
	everyMinute, flag = everyMinute.Reinstate()
	fmt.Println(flag)
	time.Sleep(5 * time.Second)
	fmt.Println(clockwork.Report())
	time.Sleep(1 * time.Minute)
	fmt.Println(clockwork.Report())
	everyMinute.Stop()
	clockwork.RemoveTaskByName(everyMinute.Name)
	time.Sleep(5 * time.Second)
	time.Sleep(1 * time.Minute)
	fmt.Println(clockwork.Report())
	clockwork.ClearTasks()
	time.Sleep(2 * time.Second)
	everySecond.Reinstate()
	time.Sleep(2 * time.Second)
	fmt.Println(clockwork.Report())
	clockwork.ClearTasks()
}

func Work(name string, count int) {
	fmt.Printf("running %s, count %d, at %s \n", name, count, time.Now().String())
}

func testData() (*clockwork.Task, *clockwork.Task, *clockwork.Task, *clockwork.Task,
	*clockwork.Task, *clockwork.Task, *clockwork.Task, *clockwork.Task,
	*clockwork.Task, *clockwork.Task, *clockwork.Task, *clockwork.Task,
	*clockwork.Task, *clockwork.Task, *clockwork.Task, *clockwork.Task, *clockwork.Task) {
	refTime := time.Now()
	return &clockwork.Task{
			Name:     "yearlyAfter",
			Function: Work,
			Schedule: &clockwork.Yearly,
			Cron:     time.Date(refTime.Year(), refTime.Month(), refTime.Day()+1, 00, 00, 00, 0, time.Local),
		}, &clockwork.Task{
			Name:     "yearlyBefore",
			Function: Work,
			Schedule: &clockwork.Yearly,
			Cron:     time.Date(refTime.Year(), refTime.Month(), refTime.Day()-1, 00, 00, 00, 0, time.Local),
		}, &clockwork.Task{
			Name:     "yearlyLocal",
			Function: Work,
			Schedule: &clockwork.Yearly,
			Cron:     time.Date(refTime.Year(), refTime.Month(), refTime.Day(), refTime.Hour()-1, 11, 00, 0, time.Local),
		},
		&clockwork.Task{
			Name:     "yearly",
			Function: Work,
			Schedule: &clockwork.Yearly,
			Cron:     time.Date(refTime.Year(), refTime.Month(), refTime.Day(), refTime.Hour()-5, refTime.Minute()-28, refTime.Second(), 0, time.UTC),
		},
		&clockwork.Task{
			Name:     "monthlyAfter",
			Function: Work,
			Schedule: &clockwork.Monthly,
			Cron:     time.Date(refTime.Year(), refTime.Month(), refTime.Day()+1, 07, 07, 07, 0, time.UTC),
		},
		&clockwork.Task{
			Name:     "monthlyBefore",
			Function: Work,
			Schedule: &clockwork.Monthly,
			Cron:     time.Date(refTime.Year(), refTime.Month(), refTime.Day()-1, 17, 30, 00, 0, time.UTC),
		},
		&clockwork.Task{
			Name:     "monthly",
			Function: Work,
			Schedule: &clockwork.Monthly,
			Cron:     time.Date(refTime.Year(), refTime.Month(), refTime.Day(), refTime.Hour()-5, refTime.Minute()-28, refTime.Second(), 0, time.UTC),
		},
		&clockwork.Task{
			Name:     "weeklyAfter",
			Function: Work,
			Schedule: &clockwork.Weekly,
			Cron:     time.Date(refTime.Year(), refTime.Month(), refTime.Day()+1, 07, 14, 21, 0, time.UTC),
		},
		&clockwork.Task{
			Name:     "weeklyBefore",
			Function: Work,
			Schedule: &clockwork.Weekly,
			Cron:     time.Date(refTime.Year(), refTime.Month(), refTime.Day()-2, 01, 11, 48, 0, time.UTC),
		},
		&clockwork.Task{
			Name:     "weekly",
			Function: Work,
			Schedule: &clockwork.Weekly,
			Cron:     time.Date(refTime.Year(), refTime.Month(), refTime.Day(), refTime.Hour()-5, refTime.Minute()-28, 00, 0, time.UTC),
		},
		&clockwork.Task{
			Name:     "dailyAfter",
			Function: Work,
			Schedule: &clockwork.Daily,
			Cron:     time.Date(refTime.Year(), refTime.Month(), refTime.Day(), 20, 56, 07, 0, time.UTC),
		},
		&clockwork.Task{
			Name:     "dailyBefore",
			Function: Work,
			Schedule: &clockwork.Daily,
			Cron:     time.Date(refTime.Year(), refTime.Month(), refTime.Day(), 1, 17, 59, 0, time.UTC),
		},
		&clockwork.Task{
			Name:     "daily",
			Function: Work,
			Schedule: &clockwork.Daily,
			Cron:     time.Date(refTime.Year(), refTime.Month(), refTime.Day(), refTime.Hour()-5, refTime.Minute()-28, 00, 0, time.UTC),
		},
		&clockwork.Task{
			Name:     "hourly",
			Function: Work,
			Schedule: &clockwork.Hourly,
			Cron:     time.Date(refTime.Year(), refTime.Month(), refTime.Day(), refTime.Hour()-5, refTime.Minute()-28, 00, 0, time.UTC),
		},
		&clockwork.Task{
			Name:     "everyMinute",
			Function: Work,
			Schedule: &clockwork.EveryMinute,
			Cron:     time.Now(),
		},
		&clockwork.Task{
			Name:     "everySecond",
			Function: Work,
			Schedule: &clockwork.EverySecond,
			Cron:     time.Now(),
		},
		&clockwork.Task{
			Name:     "dayOfTheMonth",
			Function: Work,
			Schedule: &clockwork.DayOfMonthSchedule{
				DayOfMonth: 31,
			},
			Cron: time.Now(),
		}
}

func testDayOfMonthScheduler() {
	schedule := clockwork.DayOfMonthSchedule{
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
