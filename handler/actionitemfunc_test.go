package handler

import (
	"net/http"
	model "projectname_projectmanager/model"
	"testing"

	_ "github.com/go-sql-driver/mysql" //blank import
)

func TestUpdateData(t *testing.T) {
	actionItemTest := []model.ActionItemClosed{
		{ID: 1,
			ActionItem:  "action item1",
			ProjectName: "Pepperstone",
			Owner:       "Akashnidhi Bindeshwar Prasad",
			MeetingDate: "2020-01-10",
			TargetDate:  "2020-01-15",
			Status:      "closed",
			Comment:     "it is closed",
		},
		{ID: 1,
			ActionItem:  "action item2",
			ProjectName: "Pepperstone",
			Owner:       "Akashnidhi Bindeshwar Prasad",
			MeetingDate: "2020-01-10",
			TargetDate:  "2020-01-15",
			Status:      "open",
			ClosedDate:  "2020-01-17",
			Comment:     "it is open",
		},
		{ID: 1,
			ActionItem:  "action item3",
			ProjectName: "Pepperstone",
			Owner:       "Akashnidhi Bindeshwar Prasad",
			MeetingDate: "2020-01-10",
			TargetDate:  "2020-01-09",
			Status:      "closed",
			ClosedDate:  "2020-01-17",
			Comment:     "it is closed",
		},
		{ID: 1,
			ActionItem:  "action item4",
			ProjectName: "Pepperstone",
			Owner:       "Akashnidhi Bindeshwar Prasad",
			MeetingDate: "2020-01-10",
			TargetDate:  "2020-01-15",
			Status:      "closed",
			ClosedDate:  "2020-01-09",
			Comment:     "it is closed",
		},
		{ID: 1,
			ActionItem:  "action item5",
			ProjectName: "Pepperstone",
			Owner:       "Akashnidhi Bindeshwar Prasad",
			MeetingDate: "2020-01-10",
			TargetDate:  "2020-01-21",
			Status:      "closed",
			ClosedDate:  "2020-01-30",
			Comment:     "it is closed",
		},
	}
	testData := []struct {
		x model.ActionItemClosed
		y int
		z string
	}{
		{actionItemTest[0], http.StatusBadRequest, "Closed date missing"},
		{actionItemTest[1], http.StatusBadRequest, "Closed date cannot be filled if the status is not closed"},
		{actionItemTest[2], http.StatusBadRequest, "Incorrect Target Date Value"},
		{actionItemTest[3], http.StatusBadRequest, "Incorrect Closed Date Value"},
		{actionItemTest[4], http.StatusCreated, ""},
	}
	for _, testDatum := range testData {
		updateOk, updateComment, err := UpdateData(testDatum.x, manager)
		if err != nil {
			WriteLogFile(err)
		} else {
			if updateOk != testDatum.y && updateComment != testDatum.z {
				t.Errorf("Output for %s was incorrect, got: %d and %s, want: %d and %s", testDatum.x.ActionItem, updateOk, updateComment, testDatum.y, testDatum.z)
			}

		}
	}

}

func TestDeleteData(t *testing.T) {
	testData := []struct {
		x int
		y int
	}{
		{deletetest1, http.StatusForbidden},
		{deletetest2, http.StatusOK},
		{deletetest2, http.StatusBadRequest},
	}
	for _, testDatum := range testData {
		deleteOk, err := DeleteData(testDatum.x, manager)
		if err != nil {
			WriteLogFile(err)
		} else {
			if deleteOk != testDatum.y {
				t.Errorf("Output for %d was incorrect, got: %d, want: %d", testDatum.x, deleteOk, testDatum.y)
			}

		}
	}

}
