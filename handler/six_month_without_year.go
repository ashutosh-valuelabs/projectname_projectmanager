package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	database "projectname_projectmanager/driver"
	models "projectname_projectmanager/model"
	"strings"
	"time"
)

//Get last six month active resignations
func (c *Commander) Getcompleteresignations(w http.ResponseWriter, r *http.Request) {
	db := database.DbConn()
	defer db.Close()
	w.Header().Set("Content-Type", "application/json")
	SetupResponse(&w, r)

	var post models.Activeresignationforyear
	var posts []models.Activeresignationforyear
	var inf []models.Info
	var infs models.Info
	var count [6]int
	var mname string

	var testname string = "abc"
	var ca int = 0
	var date string

	if (strings.Contains(Role, "program manager")) == true {

		wek, err := db.Query("select MONTHNAME(date_of_resignation), date_format(date_of_resignation, '%Y-%m-%d')  from active_resignations left join sub_project_manager on active_resignations.sub_project_manager_id = sub_project_manager.id left join sub_project on  sub_project_manager.sub_project_id = sub_project.id left join project on sub_project.project_id = project.id left join program_manager on project.program_manager_id=program_manager.id left join project_manager on sub_project_manager.manager_id=project_manager.id WHERE program_manager.program_manager_email = ? AND active_resignations.is_active = 1 AND date_of_resignation>=date_sub(now(), interval 06 MONTH) order by MONTH(date_of_resignation)", UserName)
		catch(err)
		defer wek.Close()
		for wek.Next() {
			wek.Scan(&mname, &date)

			layout := "2006-01-02"
			str, _ := time.Parse(layout, date)
			fmt.Println(mname, date)
			weekNo := NumberOfTheWeekInMonth(str)
			if testname == mname || testname == "abc" {
				count[weekNo]++
				testname = mname
				ca++
			} else {
				for j := 1; j < 6; j++ {
					post.CountNo = count[j]
					post.Week = j
					posts = append(posts, post)

				}
				infs.Month = testname
				infs.Total = ca
				infs.Data = posts
				inf = append(inf, infs)
				ca = 1
				posts = nil
				count[1] = 0
				count[2] = 0
				count[3] = 0
				count[4] = 0
				count[5] = 0
				count[weekNo]++
				testname = mname
			}
		}
		for j := 1; j < 6; j++ {
			post.CountNo = count[j]
			post.Week = j
			posts = append(posts, post)

		}
		infs.Month = testname
		infs.Total = ca
		ca = 1
		infs.Data = posts
		inf = append(inf, infs)
		posts = nil

		if infs.Month != "abc" {

			json.NewEncoder(w).Encode(inf)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	} else if (strings.Contains(Role, "project manager")) == true {
		wek, err := db.Query("select MONTHNAME(date_of_resignation), date_format(date_of_resignation, '%Y-%m-%d')  from active_resignations left join sub_project_manager on active_resignations.sub_project_manager_id = sub_project_manager.id left join sub_project on  sub_project_manager.sub_project_id = sub_project.id left join project_manager on sub_project_manager.manager_id = project_manager.id WHERE project_manager.project_manager_email = ? AND active_resignations.is_active = 1 AND date_of_resignation>=date_sub(now(), interval 06 MONTH) order by MONTH(date_of_resignation)", UserName)
		catch(err)
		defer wek.Close()
		for wek.Next() {
			wek.Scan(&mname, &date)

			layout := "2006-01-02"
			str, _ := time.Parse(layout, date)
			fmt.Println(mname, date)
			weekNo := NumberOfTheWeekInMonth(str)
			if testname == mname || testname == "abc" {
				count[weekNo]++
				testname = mname
				ca++
			} else {
				for j := 1; j < 6; j++ {
					post.CountNo = count[j]
					post.Week = j
					posts = append(posts, post)

				}
				infs.Month = testname
				infs.Total = ca
				infs.Data = posts
				inf = append(inf, infs)
				ca = 1
				posts = nil
				count[1] = 0
				count[2] = 0
				count[3] = 0
				count[4] = 0
				count[5] = 0
				count[weekNo]++
				testname = mname
			}
			//}
		}
		for j := 1; j < 6; j++ {
			post.CountNo = count[j]
			post.Week = j
			posts = append(posts, post)

		}
		infs.Month = testname
		infs.Total = ca
		ca = 1
		infs.Data = posts
		inf = append(inf, infs)
		posts = nil

		if infs.Month != "abc" {

			json.NewEncoder(w).Encode(inf)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
	}

}
func NumberOfTheWeekInMonth(now time.Time) int {
	beginningOfTheMonth := time.Date(now.Year(), now.Month(), 1, 1, 1, 1, 1, time.UTC)
	_, thisWeek := now.ISOWeek()
	_, beginningWeek := beginningOfTheMonth.ISOWeek()
	return 1 + thisWeek - beginningWeek
}
