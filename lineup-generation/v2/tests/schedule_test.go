package tests

import (
	"testing"
	d "lineup-generation/v2/data"
)

func TestSchedule(t *testing.T) {
	// Init the schedule
	d.InitSchedule("/Users/jameskendrick/Code/cv/stopz/src/static/schedule.json")

	if d.ScheduleMap.IsPlaying("1", 0, "SAS") != false {
		t.Errorf("IsPlaying is incorrect for SAS on day 0")
	}
	if d.ScheduleMap.IsPlaying("1", 1, "SAS") != true {
		t.Errorf("IsPlaying is incorrect for SAS on day 1")
	}
	if d.ScheduleMap.IsPlaying("17", 4, "MIN") != false {
		t.Errorf("IsPlaying is incorrect for MIN on day 4")
	}
	if d.ScheduleMap.IsPlaying("17", 12, "MIN") != true {
		t.Errorf("IsPlaying is incorrect for MIN on day 12")
	}

	// Get the schedule for week 1
	week := d.ScheduleMap.GetWeekSchedule("1")
	if week.GetStartDate() != "10/24/2023" {
		t.Errorf("StartDate is incorrect")
	}
	if week.GetEndDate() != "10/29/2023" {
		t.Errorf("EndDate is incorrect")
	}
	if week.GetGameSpan() != 5 {
		t.Errorf("GameSpan is incorrect")
	}
	if d.ScheduleMap.GetGameSpan("1") != 5 {
		t.Errorf("GameSpan is incorrect")
	}

}