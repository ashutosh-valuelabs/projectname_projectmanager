package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	database "projectname_projectmanager/driver"
	models "projectname_projectmanager/model"
)

//ActionItemGetOverview : To get the overview of Action Items
func (C *Commander) ActionItemGetOverview(writer http.ResponseWriter, response *http.Request) {
	db := database.DbConn()
	defer db.Close()
	writer.Header().Set("Content-Type", "application/json")
	var overview models.Overview
	t := time.Now().Format("2006-01-02")
	count, err := db.Query("SELECT COUNT(id) FROM action_items WHERE closed_in_time = 1")
	catch(err)
	defer count.Close()
	for count.Next() {
		err := count.Scan(&overview.ClosedInTime)
		catch(err)
	}
	count1, err := db.Query("SELECT COUNT(id) FROM action_items WHERE closed_in_time = 0")
	catch(err)
	defer count1.Close()
	for count1.Next() {
		err := count1.Scan(&overview.ClosedOutTime)
		catch(err)
	}
	count2, err := db.Query("SELECT COUNT(id) FROM action_items WHERE ? < target_date AND status LIKE 'inprogress'", t)
	catch(err)
	defer count2.Close()
	for count2.Next() {
		err := count2.Scan(&overview.InProgressInTime)
		catch(err)
	}
	count3, err := db.Query("SELECT COUNT(id) FROM action_items WHERE ? > target_date AND status LIKE 'inprogress'", t)
	catch(err)
	defer count3.Close()
	for count3.Next() {
		err := count3.Scan(&overview.InProgressOutTime)
		catch(err)
	}
	count4, err := db.Query("SELECT COUNT(id) FROM action_items WHERE status LIKE 'onhold'")
	catch(err)
	defer count4.Close()
	for count4.Next() {
		err := count4.Scan(&overview.OnHold)
		catch(err)
	}
	overview.GrandTotal = 0
	overview.GrandTotal = overview.GrandTotal + overview.ClosedInTime + overview.ClosedOutTime + overview.InProgressInTime + overview.InProgressOutTime + overview.OnHold
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(overview)
}

//ActionItemGetSayDo : to get the SayDo Ratio of Action Items
func (C *Commander) ActionItemGetSayDo(writer http.ResponseWriter, response *http.Request) {
	db := database.DbConn()
	defer db.Close()
	writer.Header().Set("Content-Type", "application/json")
	var total float64
	var closedintime float64
	var overviews []models.Manager
	var overview models.Manager
	var data models.Data
	count, err := db.Query("SELECT project_manager_name, count(project_manager_name) FROM action_items LEFT JOIN sub_project_manager ON action_items.sub_project_manager_id = sub_project_manager.id LEFT JOIN project_manager ON sub_project_manager.manager_id = project_manager.id LEFT JOIN sub_project ON sub_project_manager.sub_project_id = sub_project.id WHERE action_items.is_active = 1 AND WEEK(CURDATE())-WEEK(action_items.updated_at) < 1 GROUP BY project_manager_name")
	catch(err)
	defer count.Close()
	for count.Next() {
		err := count.Scan(&overview.Name, &total)
		fmt.Println(total)
		catch(err)
		Owner := overview.Name
		count1, err := db.Query("SELECT count(project_manager_name) FROM action_items LEFT JOIN sub_project_manager ON action_items.sub_project_manager_id = sub_project_manager.id LEFT JOIN project_manager ON sub_project_manager.manager_id = project_manager.id LEFT JOIN sub_project ON sub_project_manager.sub_project_id = sub_project.id WHERE closed_in_time = 1 AND action_items.is_active = 1 AND project_manager.project_manager_email = ? GROUP BY project_manager_name", UserName)
		catch(err)
		defer count1.Close()
		for count1.Next() {
			err := count1.Scan(&closedintime)
			catch(err)
		}
		fmt.Println(closedintime)
		overview.SayDo = 0.0
		overview.SayDo = ((closedintime / total) * 100)
		manager, err := db.Query("SELECT manager FROM saydo WHERE manager = ?", Owner)
		catch(err)
		if manager.Next() != false {
			update, err := db.Prepare("UPDATE saydo SET saydoratio = ? WHERE manager = ?")
			catch(err)
			update.Exec(overview.SayDo, Owner)
		} else {
			insert, err := db.Prepare("INSERT INTO saydo(manager,saydoratio) VALUES(?,?)")
			catch(err)
			Say := overview.SayDo
			insert.Exec(Owner, Say)
		}
		saydo, err := db.Query("SELECT saydoratio FROM saydo WHERE manager = ?", Owner)
		catch(err)
		defer saydo.Close()
		for saydo.Next() {
			err := saydo.Scan(&overview.SayDo)
			catch(err)
		}
		overviews = append(overviews, overview)
		data.Data = overviews
	}
	if count.Next() == false {
		count, err := db.Query("SELECT project_manager_name,count(project_manager_name) FROM action_items LEFT JOIN sub_project_manager ON action_items.sub_project_manager_id = sub_project_manager.id LEFT JOIN project_manager ON sub_project_manager.manager_id = project_manager.id LEFT JOIN sub_project ON sub_project_manager.sub_project_id = sub_project.id WHERE action_items.is_active = 1 GROUP BY project_manager_name ")
		catch(err)
		defer count.Close()
		for count.Next() {
			count.Scan(&overview.Name, &total)
			overview.SayDo = 0
			overviews = append(overviews, overview)
		}
		data.Data = overviews
	}

	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(data)
}

//ActionItemGetTrend : to get the trend of Action Items
func (C *Commander) ActionItemGetTrend(writer http.ResponseWriter, response *http.Request) {
	db := database.DbConn()
	defer db.Close()
	writer.Header().Set("Content-Type", "application/json")
	var wk, s int
	var ct, saydo float32
	var counts, counts1 [53]float32
	var week, week1 [53]int
	var SayDOTrend models.Saydo
	var SayDoTrends []models.Saydo
	var data models.SayDoData

	offsets, ok := response.URL.Query()["year"]
	if !ok || len(offsets[0]) < 1 {
		respondwithJSON(writer, http.StatusBadRequest, map[string]string{"Message": "Bad Request ; URL Parameter 'year' is missing"})
		return
	}
	year := offsets[0]
	yr, err := strconv.Atoi(year)
	catch(err)

	count, err := db.Query("SELECT COUNT(WEEK(updated_at)),WEEK(updated_at) FROM action WHERE YEAR(updated_at) = ? GROUP BY WEEK(updated_at)", yr)
	defer count.Close()
	catch(err)
	i := 0
	for count.Next() {
		count.Scan(&ct, &wk)
		counts[wk] = ct
		week[i] = wk
		i++
		s++
	}
	wk = 0
	ct = 0
	closed, err := db.Query("SELECT COUNT(WEEK(updated_at)),WEEK(updated_at) FROM action WHERE closed_in_time = 1 AND YEAR(updated_at) = ? GROUP BY WEEK(updated_at) ", yr)
	defer closed.Close()
	catch(err)
	i = 0
	for closed.Next() {
		closed.Scan(&ct, &wk)
		counts1[wk] = ct
		week1[i] = wk
		i++

	}
	for i = 0; i < s; i++ {
		saydo = counts1[week[i]] / counts[week[i]] * 100
		Week := week[i]
		SayDOTrend.Saydo = saydo
		x, y := WeekRange(yr, Week)
		x1 := x.Format("02-01-2006")
		y1 := y.Format("02-01-2006")
		SayDOTrend.WeekStart = x1
		SayDOTrend.WeekEnd = y1
		SayDoTrends = append(SayDoTrends, SayDOTrend)
		data.Data = SayDoTrends
	}
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(data)
}

// WeekRange : To get the range of week using start date
func WeekRange(year, week int) (start, end time.Time) {
	start = WeekStart(year, week)
	end = start.AddDate(0, 0, 6)
	return
}

// WeekStart : To get the start date using week number
func WeekStart(year, week int) time.Time {
	// Start from the middle of the year:
	t := time.Date(year, 7, 1, 0, 0, 0, 0, time.UTC)
	// Roll back to Monday:
	if wd := t.Weekday(); wd == time.Sunday {
		t = t.AddDate(0, 0, -6)
	} else {
		t = t.AddDate(0, 0, -int(wd)+1)
	}
	_, writer := t.ISOWeek()
	t = t.AddDate(0, 0, (week-writer)*7)
	return t
}
