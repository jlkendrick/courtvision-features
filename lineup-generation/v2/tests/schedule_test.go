package tests

import (
	"strconv"
	"testing"
	d "v2/data"
)

func TestSchedule(t *testing.T) {
	// Init the schedule
	d.InitSchedule("/Users/jameskendrick/Code/cv/features/lineup-generation/v2/static/schedule24-25.json")

	// Check that the schedule has been initialized
	for week:=1; week<=22; week++ {
		week_str := strconv.Itoa(week)
		week_schedule := d.ScheduleMap.GetWeekSchedule(week_str)
		if week_schedule.StartDate == "" {
			t.Errorf("Week %v schedule not initialized", week)
		}
	}

}