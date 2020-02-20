package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	database "projectname_projectmanager/driver"
	models "projectname_projectmanager/model"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

// Getallopenpositions :Func for fetching all entries
func (c *Commander) Getallopenpositions(w http.ResponseWriter, r *http.Request) {

	db := database.DbConn()
	defer db.Close()

	w.Header().Set("Content-Type", "application/json")

	if strings.Contains(Role, "Project Manager") || strings.Contains(Role, "project manager") {
		statusGet, ok := r.URL.Query()["status"]

		if !ok || statusGet[0] == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		status := statusGet[0]

		if status == "daily" || status == "Daily" {
			offsetGet, ok := r.URL.Query()["pages"]

			if !ok || offsetGet[0] == "" {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			off := offsetGet[0]
			offInt, _ := strconv.Atoi(off)
			offset := offInt * 10

			limitGet, ok := r.URL.Query()["limit"]

			if !ok || limitGet[0] == "" {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			limit, _ := strconv.Atoi(limitGet[0])

			var pagination models.Pagination
			var dailyArray []models.Daily

			queryDaily, err := db.Query("select o.id,p.sub_project_name, o.designation,o.type_position,o.position, o.priority, o.additional_comment, o.l1_due,o.l2_due, o.client_due, o.created_at from open_positions o inner join sub_project_manager s on o.sub_project_manager_id = s.id inner join sub_project p on s.sub_project_id = p.id where s.manager_id in (select id from project_manager where project_manager_email = ?)  and o.is_active = 1 LIMIT ?, ? ", UserName, offset, limit)
			catch(err)
			defer queryDaily.Close()

			count, err := db.Query("select count(id) from open_positions where is_active = 1 and sub_project_manager_id in ( select id from sub_project_manager where manager_id in ( select id from project_manager where project_manager_email = ?))", UserName)
			catch(err)
			defer count.Close()

			for count.Next() {
				err2 := count.Scan(&pagination.TotalData)
				catch(err2)

			}
			pagination.Limit = limit
			x := pagination.TotalData / pagination.Limit
			x1 := pagination.TotalData % pagination.Limit
			if x1 == 0 {
				pagination.TotalPages = x
			} else {
				pagination.TotalPages = x + 1
			}
			pagination.Page = offInt + 1

			for queryDaily.Next() {
				var createdAt string
				var daily models.Daily
				err := queryDaily.Scan(&daily.Id, &daily.Project_name, &daily.Designation, &daily.Type_position, &daily.Position, &daily.Priority, &daily.Additonal_comment, &daily.L1_due, &daily.L2_due, &daily.Client_due, &createdAt)
				catch(err)

				location, err := time.LoadLocation("Asia/Kolkata")
				catch(err)

				layout := "2006-01-02 15:04:05"
				t, err := time.Parse(layout, createdAt)
				to := time.Now()
				to = to.In(location)
				days := to.Sub(t) / (24 * time.Hour)
				daily.Ageing = int(days)

				dailyArray = append(dailyArray, daily)
				pagination.Data = dailyArray

			}
			json.NewEncoder(w).Encode(pagination)

		} else if status == "Weekly" || status == "weekly" {

			offsetGet, ok := r.URL.Query()["pages"]
			if !ok || offsetGet[0] == "" {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			off := offsetGet[0]
			offInt, _ := strconv.Atoi(off)
			offset := offInt * 10

			limitGet, ok := r.URL.Query()["limit"]

			if !ok || limitGet[0] == "" {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			limit := limitGet[0]
			limitFloat, _ := strconv.Atoi(limit)

			var pagination models.Pagination
			var weeklyArray []models.Weekly

			queryWeekly, err := db.Query("select o.id,p.sub_project_name, o.designation,o.type_position,o.position, o.priority, o.additional_comment, o.l1_happened,o.l2_happened, o.client_happened, o.created_at from open_positions o inner join sub_project_manager s on o.sub_project_manager_id = s.id inner join sub_project p on s.sub_project_id = p.id where s.manager_id in (select id from project_manager where project_manager_email = ?)  and o.is_active = 1 LIMIT ?, ?  ", UserName, offset, limitFloat)
			catch(err)
			defer queryWeekly.Close()

			count, err1 := db.Query("select count(id) from open_positions where is_active = 1 and sub_project_manager_id in ( select id from sub_project_manager where manager_id in ( select id from project_manager where project_manager_email = ?))", UserName)
			catch(err1)
			defer count.Close()

			for count.Next() {
				err := count.Scan(&pagination.TotalData)
				catch(err)
			}

			pagination.Limit = limitFloat
			x := pagination.TotalData / pagination.Limit
			x1 := pagination.TotalData % pagination.Limit
			if x1 == 0 {
				pagination.TotalPages = x
			} else {
				pagination.TotalPages = x + 1
			}
			pagination.Page = offInt + 1
			for queryWeekly.Next() {
				var weekly models.Weekly
				var createdAt string

				err := queryWeekly.Scan(&weekly.Id, &weekly.Project_name, &weekly.Designation, &weekly.Type_position, &weekly.Position, &weekly.Priority, &weekly.Additonal_comment, &weekly.L1_Happened, &weekly.L2_Happened, &weekly.Client_Happened, &createdAt)
				catch(err)

				location, err := time.LoadLocation("Asia/Kolkata")
				catch(err)

				layout := "2006-01-02 15:04:05"
				t, err := time.Parse(layout, createdAt)
				to := time.Now()
				to = to.In(location)
				days := to.Sub(t) / (24 * time.Hour)
				weekly.Ageing = int(days)

				weeklyArray = append(weeklyArray, weekly)
				pagination.Data = weeklyArray
			}
			json.NewEncoder(w).Encode(pagination)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}

	} else if strings.Contains(Role, "Program Manager") || strings.Contains(Role, "program manager") {
		statusGet, ok := r.URL.Query()["status"]

		if !ok || statusGet[0] == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		status := statusGet[0]

		if status == "daily" || status == "Daily" {
			offsetGet, ok := r.URL.Query()["pages"]

			if !ok || offsetGet[0] == "" {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			off := offsetGet[0]
			offInt, _ := strconv.Atoi(off)
			offset := offInt * 10

			limitGet, ok := r.URL.Query()["limit"]

			if !ok || limitGet[0] == "" {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			limit := limitGet[0]
			limitFloat, _ := strconv.Atoi(limit)

			var pagination models.Pagination
			var dailyarray []models.Daily
			queryDaily, err := db.Query("select o.id,p.project_name, o.designation,o.type_position,o.position, o.priority, o.additional_comment, o.l1_due,o.l2_due, o.client_due, o.created_at from open_positions o inner join sub_project_manager s on o.sub_project_manager_id = s.id inner join sub_project sp on s.sub_project_id = sp.id inner join project p on sp.project_id = p.id where p.program_manager_id in (select id from program_manager where program_manager_email = ?)  and o.is_active = 1 LIMIT ?, ?", UserName, offset, limitFloat)
			catch(err)
			defer queryDaily.Close()
			count, err1 := db.Query("select count(id) from open_positions where is_active = 1 and sub_project_manager_id in (select id from sub_project_manager where sub_project_id in ( select id from sub_project where project_id in ( select id from project where program_manager_id in ( select id from program_manager where program_manager_email = ? ))))", UserName)
			catch(err1)

			defer count.Close()
			for count.Next() {
				err = count.Scan(&pagination.TotalData)
				catch(err)

			}
			pagination.Limit = limitFloat
			x := pagination.TotalData / pagination.Limit
			x1 := pagination.TotalData % pagination.Limit
			if x1 == 0 {
				pagination.TotalPages = x
			} else {
				pagination.TotalPages = x + 1
			}
			pagination.Page = offInt + 1
			for queryDaily.Next() {
				var daily models.Daily
				var createdAt string
				err := queryDaily.Scan(&daily.Id, &daily.Project_name, &daily.Designation, &daily.Type_position, &daily.Position, &daily.Priority, &daily.Additonal_comment, &daily.L1_due, &daily.L2_due, &daily.Client_due, &createdAt)
				catch(err)
				location, err := time.LoadLocation("Asia/Kolkata")
				catch(err)
				layout := "2006-01-02 15:04:05"
				t, err := time.Parse(layout, createdAt)
				to := time.Now()
				to = to.In(location)
				days := to.Sub(t) / (24 * time.Hour)
				daily.Ageing = int(days)

				dailyarray = append(dailyarray, daily)
				pagination.Data = dailyarray

			}
			json.NewEncoder(w).Encode(pagination)
		} else if status == "Weekly" || status == "weekly" {
			offsetGet, ok := r.URL.Query()["pages"]

			if !ok || offsetGet[0] == "" {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			off := offsetGet[0]
			offInt, _ := strconv.Atoi(off)
			offset := offInt * 10

			limitGet, ok := r.URL.Query()["limit"]

			if !ok || limitGet[0] == "" {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			limit := limitGet[0]
			limitFloat, _ := strconv.Atoi(limit)

			var pagination models.Pagination
			var weeklyArray []models.Weekly
			queryWeekly, err := db.Query("select o.id,p.project_name, o.designation,o.type_position,o.position, o.priority, o.additional_comment, o.l1_happened,o.l2_happened, o.client_happened,  o.created_at from open_positions o inner join sub_project_manager s on o.sub_project_manager_id = s.id inner join sub_project sp on s.sub_project_id = sp.id inner join project p on sp.project_id = p.id where p.program_manager_id in (select id from program_manager where program_manager_email= ?)  and o.is_active = 1 LIMIT ?, ? ", UserName, offset, limitFloat)
			catch(err)
			defer queryWeekly.Close()
			count, err := db.Query("select count(id) from open_positions where is_active = 1 and sub_project_manager_id in (select id from sub_project_manager where sub_project_id in ( select id from sub_project where project_id in ( select id from project where program_manager_id in ( select id from program_manager where program_manager_email = ? ))))", UserName)
			catch(err)

			defer count.Close()
			for count.Next() {
				err = count.Scan(&pagination.TotalData)
				catch(err)
			}
			pagination.Limit = limitFloat
			x := pagination.TotalData / pagination.Limit
			x1 := pagination.TotalData % pagination.Limit
			if x1 == 0 {
				pagination.TotalPages = x
			} else {
				pagination.TotalPages = x + 1
			}
			pagination.Page = offInt + 1
			for queryWeekly.Next() {
				var weekly models.Weekly
				var createdAt string

				err := queryWeekly.Scan(&weekly.Id, &weekly.Project_name, &weekly.Designation, &weekly.Type_position, &weekly.Position, &weekly.Priority, &weekly.Additonal_comment, &weekly.L1_Happened, &weekly.L2_Happened, &weekly.Client_Happened, &createdAt)
				catch(err)
				location, err := time.LoadLocation("Asia/Kolkata")
				catch(err)
				layout := "2006-01-02 15:04:05"
				t, err := time.Parse(layout, createdAt)
				to := time.Now()
				to = to.In(location)
				days := to.Sub(t) / (24 * time.Hour)

				weekly.Ageing = int(days)

				weeklyArray = append(weeklyArray, weekly)
				pagination.Data = weeklyArray

			}
			json.NewEncoder(w).Encode(pagination)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}
}

// Createopenpositions :Func for creating new entries
func (c *Commander) Createopenpositions(w http.ResponseWriter, r *http.Request) {
	db := database.DbConn()
	defer db.Close()

	w.Header().Set("Content-Type", "application/json")

	if strings.Contains(Role, "Project Manager") || strings.Contains(Role, "project manager") {
		statusGet, ok := r.URL.Query()["status"]

		if !ok || statusGet[0] == "" {
			w.WriteHeader(http.StatusCreated)
			return
		}
		status := statusGet[0]

		if status == "daily" || status == "Daily" {
			var daily models.Daily
			json.NewDecoder(r.Body).Decode(&daily)

			querySubProjectManagerID := db.QueryRow("select id from sub_project_manager where sub_project_id in ( select id from sub_project where sub_project_name = ? ) and manager_id in ( select id from project_manager where project_manager_email = ? )", daily.Project_name, UserName)
			var managerid int
			querySubProjectManagerID.Scan(&managerid)

			if managerid != 0 {
				var id int
				queryExisting := db.QueryRow("select id from open_positions where sub_project_manager_id = ? and designation = ? and type_position = ?  and is_active = 1 ", managerid, daily.Designation, daily.Type_position)
				queryExisting.Scan(&id)

				if id == 0 {
					queryInsert, err := db.Prepare("INSERT INTO open_positions( sub_project_manager_id, designation, type_position, position, priority, additional_comment, l1_due, l2_due, client_due, created_at, updated_at) VALUES( ?, ?, ?, ?, ?, ?, ?, ?, ?, now(), now())")
					catch(err)
					_, err = queryInsert.Exec(managerid, daily.Designation, daily.Type_position, daily.Position, daily.Priority, daily.Additonal_comment, daily.L1_due, daily.L2_due, daily.Client_due)
					catch(err)
					w.WriteHeader(http.StatusCreated)
					defer queryInsert.Close()
				} else {
					w.WriteHeader(http.StatusConflict)
				}
			} else {
				w.WriteHeader(http.StatusForbidden)
			}
		} else if status == "Weekly" || status == "weekly" {
			var weekly models.Weekly
			json.NewDecoder(r.Body).Decode(&weekly)

			querySubProjectManagerID := db.QueryRow("select id from sub_project_manager where sub_project_id in ( select id from sub_project where sub_project_name = ? ) and manager_id in ( select id from project_manager where project_manager_email = ? )", weekly.Project_name, UserName)
			var managerid int
			querySubProjectManagerID.Scan(&managerid)

			if managerid != 0 {
				var id int
				queryExisting := db.QueryRow("select id from open_positions where sub_project_manager_id = ? and designation = ? and type_position = ? and is_active = 1 ", managerid, weekly.Designation, weekly.Type_position)
				queryExisting.Scan(&id)
				if id != 0 {
					queryUpdate, err := db.Prepare("UPDATE open_positions set l1_happened = ?, l2_happened=?, client_happened= ?,  updated_at =  now() where id = ? ")
					catch(err)
					_, err = queryUpdate.Exec(weekly.L1_Happened, weekly.L2_Happened, weekly.Client_Happened, id)
					catch(err)
					w.WriteHeader(http.StatusCreated)
					defer queryUpdate.Close()
				} else {
					w.WriteHeader(http.StatusNotFound)
				}
			} else {
				w.WriteHeader(http.StatusForbidden)
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}

}

// Getopenpositions :Func for fetching specific entries
func (c *Commander) Getopenpositions(w http.ResponseWriter, r *http.Request) {
	db := database.DbConn()
	defer db.Close()

	w.Header().Set("Content-Type", "application/json")

	if strings.Contains(Role, "Project Manager") || strings.Contains(Role, "project manager") {
		params := mux.Vars(r)
		statusGet, ok := r.URL.Query()["status"]

		if !ok || statusGet[0] == "" {
			w.WriteHeader(http.StatusBadRequest)
			//fmt.Fprintf(w, "Url Param status is missing")
			return
		}
		status := statusGet[0]

		if status == "daily" || status == "Daily" {

			offsetGet, ok := r.URL.Query()["pages"]

			if !ok || offsetGet[0] == "" {
				w.WriteHeader(http.StatusBadRequest)
				//fmt.Fprintf(w, "Url Param offset is missing")
				return
			}

			off := offsetGet[0]
			offsetInt, _ := strconv.Atoi(off)
			offset := offsetInt * 10

			limitGet, ok := r.URL.Query()["limit"]

			if !ok || limitGet[0] == "" {
				w.WriteHeader(http.StatusBadRequest)
				//fmt.Fprintf(w, "Url Param limit is missing")
				return
			}

			limit := limitGet[0]
			limitFloat, _ := strconv.Atoi(limit)

			var pagination models.Pagination
			var dailyArray []models.Daily
			var search string
			search = "'" + params["id"] + "%'"
			queryDaily, err := db.Query("select o.id,p.sub_project_name, o.designation,o.type_position,o.position, o.priority, o.additional_comment, o.l1_due,o.l2_due, o.client_due, o.created_at from open_positions o inner join sub_project_manager s on o.sub_project_manager_id = s.id inner join sub_project p on s.sub_project_id = p.id where s.manager_id in (select id from project_manager where project_manager_email = ?) and o.designation like "+search+" OR o.type_position like "+search+" OR o.position like "+search+" OR o.priority like "+search+" OR o.additional_comment like "+search+" OR o.l1_due like "+search+" OR o.l2_due like "+search+" OR o.client_due like "+search+"  and o.is_active = 1 LIMIT ?, ? ", UserName, offset, limitFloat)
			catch(err)
			defer queryDaily.Close()
			count, err := db.Query("select count(id) from open_positions where is_active = 1 and sub_project_manager_id in ( select id from sub_project_manager where manager_id in ( select id from project_manager where project_manager_email = ?)) and designation like "+search+" OR type_position like "+search+" OR position like "+search+" OR priority like "+search+" OR additional_comment like "+search+" OR l1_due like "+search+" OR l2_due like "+search+" OR client_due like "+search+"", UserName)
			catch(err)

			defer count.Close()
			for count.Next() {
				err = count.Scan(&pagination.TotalData)
				catch(err)

			}
			pagination.Limit = limitFloat
			x := pagination.TotalData / pagination.Limit
			x1 := pagination.TotalData % pagination.Limit
			if x1 == 0 {
				pagination.TotalPages = x
			} else {
				pagination.TotalPages = x + 1
			}
			pagination.Page = offsetInt + 1
			for queryDaily.Next() {
				var daily models.Daily
				var createdAt string
				err := queryDaily.Scan(&daily.Id, &daily.Project_name, &daily.Designation, &daily.Type_position, &daily.Position, &daily.Priority, &daily.Additonal_comment, &daily.L1_due, &daily.L2_due, &daily.Client_due, &createdAt)
				catch(err)
				location, err := time.LoadLocation("Asia/Kolkata")
				catch(err)
				layout := "2006-01-02 15:04:05"
				t, err := time.Parse(layout, createdAt)
				to := time.Now()
				to = to.In(location)
				days := to.Sub(t) / (24 * time.Hour)

				daily.Ageing = int(days)

				//if daily.Is_active == "1" {
				dailyArray = append(dailyArray, daily)
				//}
				pagination.Data = dailyArray
			}
			json.NewEncoder(w).Encode(pagination)
		} else if status == "Weekly" || status == "weekly" {
			offsetGet, ok := r.URL.Query()["pages"]

			if !ok || offsetGet[0] == "" {
				fmt.Fprintf(w, "Url Param offset is missing")
				return
			}

			off := offsetGet[0]
			offsetInt, _ := strconv.Atoi(off)
			offset := offsetInt * 10

			limitGet, ok := r.URL.Query()["limit"]

			if !ok || limitGet[0] == "" {
				w.WriteHeader(http.StatusBadRequest)
				//fmt.Fprintf(w, "Url Param limit is missing")
				return
			}

			limit := limitGet[0]
			limitFloat, _ := strconv.Atoi(limit)

			var pagination models.Pagination
			var weeklyArray []models.Weekly
			var search string
			search = "'" + params["id"] + "%'"
			queryWeekly, err := db.Query("select o.id,p.sub_project_name, o.designation,o.type_position,o.position, o.priority, o.additional_comment, o.l1_happened,o.l2_happened, o.client_happened, o.created_at from open_positions o inner join sub_project_manager s on o.sub_project_manager_id = s.id inner join sub_project p on s.sub_project_id = p.id where s.manager_id in (select id from project_manager where project_manager_email = ?) and o.designation like "+search+" OR o.type_position like "+search+" OR o.position like "+search+" OR o.priority like "+search+" OR o.additional_comment like "+search+" OR o.l1_happened like "+search+" OR o.l2_happened like "+search+" OR o.client_happened like "+search+" LIMIT ?, ? ", UserName, offset, limitFloat)
			catch(err)
			defer queryWeekly.Close()
			count, err := db.Query("select count(id) from open_positions where is_active = 1 and sub_project_manager_id in ( select id from sub_project_manager where manager_id in ( select id from project_manager where project_manager_email = ?)) and designation like "+search+" OR type_position like "+search+" OR position like "+search+" OR priority like "+search+" OR additional_comment like "+search+" OR l1_happened like "+search+" OR l2_happened like "+search+" OR client_happened like "+search+"", UserName)
			catch(err)

			defer count.Close()
			for count.Next() {
				err = count.Scan(&pagination.TotalData)
				catch(err)
			}
			pagination.Limit = limitFloat
			x := pagination.TotalData / pagination.Limit
			x1 := pagination.TotalData % pagination.Limit
			if x1 == 0 {
				pagination.TotalPages = x
			} else {
				pagination.TotalPages = x + 1
			}
			pagination.Page = offsetInt + 1
			for queryWeekly.Next() {
				var weekly models.Weekly
				var createdAt string

				err := queryWeekly.Scan(&weekly.Id, &weekly.Project_name, &weekly.Designation, &weekly.Type_position, &weekly.Position, &weekly.Priority, &weekly.Additonal_comment, &weekly.L1_Happened, &weekly.L2_Happened, &weekly.Client_Happened, &createdAt)
				catch(err)
				location, err := time.LoadLocation("Asia/Kolkata")
				catch(err)
				layout := "2006-01-02 15:04:05"
				t, err := time.Parse(layout, createdAt)
				to := time.Now()
				to = to.In(location)
				days := to.Sub(t) / (24 * time.Hour)

				weekly.Ageing = int(days)

				//if weekly.Is_active == "1" {
				weeklyArray = append(weeklyArray, weekly)
				//	}
				pagination.Data = weeklyArray

			}
			json.NewEncoder(w).Encode(pagination)
		} else {
			w.WriteHeader(http.StatusBadRequest)
			//fmt.Fprintf(w, "No status")
		}
	} else if strings.Contains(Role, "Program Manager") || strings.Contains(Role, "program manager") {
		params := mux.Vars(r)
		st, ok1 := r.URL.Query()["status"]

		if !ok1 || st[0] == "" {
			fmt.Fprintf(w, "status is missing")
			return
		}
		status := st[0]

		if status == "daily" || status == "Daily" {
			offsetGet, ok := r.URL.Query()["pages"]

			if !ok || offsetGet[0] == "" {
				fmt.Fprintf(w, "Url Param offset is missing")
				return
			}

			off := offsetGet[0]
			offsetInt, _ := strconv.Atoi(off)
			offset := offsetInt * 10

			limitGet, ok := r.URL.Query()["limit"]

			if !ok || limitGet[0] == "" {
				w.WriteHeader(http.StatusBadRequest)
				//fmt.Fprintf(w, "Url Param limit is missing")
				return
			}

			limit := limitGet[0]
			limitFloat, _ := strconv.Atoi(limit)

			var pagination models.Pagination
			var dailyArray []models.Daily
			var search string
			search = "'" + params["id"] + "%'"
			queryDaily, err := db.Query("select o.id,p.project_name, o.designation,o.type_position,o.position, o.priority, o.additional_comment, o.l1_due,o.l2_due, o.client_due, o.created_at from open_positions o inner join sub_project_manager s on o.sub_project_manager_id = s.id inner join sub_project sp on s.sub_project_id = sp.id inner join project p on sp.project_id = p.id where p.program_manager_id in (select id from program_manager where program_manager_email = ?)  and o.is_active = 1 and o.designation like "+search+" OR o.type_position like "+search+" OR o.position like "+search+" OR o.priority like "+search+" OR o.additional_comment like "+search+" OR o.l1_due like "+search+" OR o.l2_due like "+search+" OR o.client_due like "+search+" LIMIT ?, ? ", UserName, offset, limitFloat)
			catch(err)
			defer queryDaily.Close()
			count, err := db.Query("select count(id) from open_positions where is_active = 1 and sub_project_manager_id in ( select id from sub_project_manager where sub_project_id in ( select id from sub_project where project_id in ( select id from project where program_manager_id in ( select id from program_manager where program_manager_email = ? )))) and designation like "+search+" OR type_position like "+search+" OR position like "+search+" OR priority like "+search+" OR additional_comment like "+search+" OR l1_due like "+search+" OR l2_due like "+search+" OR client_due like "+search+"", UserName)
			catch(err)
			defer count.Close()
			for count.Next() {
				err := count.Scan(&pagination.TotalData)
				catch(err)
			}
			pagination.Limit = limitFloat
			x := pagination.TotalData / pagination.Limit
			x1 := pagination.TotalData % pagination.Limit
			if x1 == 0 {
				pagination.TotalPages = x
			} else {
				pagination.TotalPages = x + 1
			}
			pagination.Page = offsetInt + 1
			for queryDaily.Next() {
				var daily models.Daily
				var createdAt string
				err := queryDaily.Scan(&daily.Id, &daily.Project_name, &daily.Designation, &daily.Type_position, &daily.Position, &daily.Priority, &daily.Additonal_comment, &daily.L1_due, &daily.L2_due, &daily.Client_due, &createdAt)
				catch(err)
				location, err := time.LoadLocation("Asia/Kolkata")
				catch(err)
				layout := "2006-01-02 15:04:05"
				t, err := time.Parse(layout, createdAt)
				to := time.Now()
				to = to.In(location)
				days := to.Sub(t) / (24 * time.Hour)

				daily.Ageing = int(days)

				//	if daily.Is_active == "1" {
				dailyArray = append(dailyArray, daily)
				//	}
				pagination.Data = dailyArray

			}
			json.NewEncoder(w).Encode(pagination)
		} else if status == "Weekly" || status == "weekly" {
			offsetGet, ok := r.URL.Query()["pages"]

			if !ok || offsetGet[0] == "" {
				fmt.Fprintf(w, "Url Param offset is missing")
				return
			}

			off := offsetGet[0]
			offsetInt, _ := strconv.Atoi(off)
			offset := offsetInt * 10

			limitGet, ok := r.URL.Query()["limit"]

			if !ok || limitGet[0] == "" {
				w.WriteHeader(http.StatusBadRequest)
				//fmt.Fprintf(w, "Url Param limit is missing")
				return
			}

			limit := limitGet[0]
			limitFloat, _ := strconv.Atoi(limit)

			var pagination models.Pagination
			var weeklyArray []models.Weekly
			var search string
			search = "'" + params["id"] + "%'"
			queryWeekly, err := db.Query("select o.id,p.project_name, o.designation,o.type_position,o.position, o.priority, o.additional_comment, o.l1_happened,o.l2_happened, o.client_happened, o.created_at from open_positions o inner join sub_project_manager s on o.sub_project_manager_id = s.id inner join sub_project sp on s.sub_project_id = sp.id inner join project p on sp.project_id = p.id where p.program_manager_id in (select id from program_manager where program_manager_email = ?)  and o.is_active = 1 and designation like "+search+" OR type_position like "+search+" OR position like "+search+" OR priority like "+search+" OR additional_comment like "+search+" OR l1_happened like "+search+" OR l2_happened like "+search+" OR client_happened like "+search+"  LIMIT ?, ? ", UserName, offset, limitFloat)
			catch(err)
			defer queryWeekly.Close()
			count, err := db.Query("select count(id) from open_positions where is_active = 1 and sub_project_manager_id in ( select id from sub_project_manager where sub_project_id in ( select id from sub_project where project_id in ( select id from project where program_manager_id in ( select id from program_manager where program_manager_email = ? )))) and designation like "+search+" OR type_position like "+search+" OR position like "+search+" OR priority like "+search+" OR additional_comment like "+search+" OR l1_happened like "+search+" OR l2_happened like "+search+" OR client_happened like "+search+"", UserName)
			catch(err)
			defer count.Close()
			for count.Next() {
				err := count.Scan(&pagination.TotalData)
				catch(err)
			}
			pagination.Limit = limitFloat
			x := pagination.TotalData / pagination.Limit
			x1 := pagination.TotalData % pagination.Limit
			if x1 == 0 {
				pagination.TotalPages = x
			} else {
				pagination.TotalPages = x + 1
			}
			pagination.Page = offsetInt + 1
			for queryWeekly.Next() {
				var weekly models.Weekly
				var createdAt string

				err := queryWeekly.Scan(&weekly.Id, &weekly.Project_name, &weekly.Designation, &weekly.Type_position, &weekly.Position, &weekly.Priority, &weekly.Additonal_comment, &weekly.L1_Happened, &weekly.L2_Happened, &weekly.Client_Happened, &createdAt)
				catch(err)
				location, err := time.LoadLocation("Asia/Kolkata")
				catch(err)
				layout := "2006-01-02 15:04:05"
				t, err := time.Parse(layout, createdAt)
				to := time.Now()
				to = to.In(location)
				days := to.Sub(t) / (24 * time.Hour)

				weekly.Ageing = int(days)

				//	if weekly.Is_active == "1" {
				weeklyArray = append(weeklyArray, weekly)
				//	}
				pagination.Data = weeklyArray

			}
			json.NewEncoder(w).Encode(pagination)
		}
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		//fmt.Fprintf(w, "No authentication")
	}

}

// Deleteopenpositions :Func for deleting entries
func (c *Commander) Deleteopenpositions(w http.ResponseWriter, r *http.Request) {
	db := database.DbConn()
	defer db.Close()
	w.Header().Set("Content-Type", "application/json")

	if strings.Contains(Role, "Project Manager") || strings.Contains(Role, "project manager") {
		var daily models.Daily
		json.NewDecoder(r.Body).Decode(&daily)
		querySubProjectManagerEmail := db.QueryRow("select project_manager_email from project_manager where id in ( select manager_id from sub_project_manager where id in ( select sub_project_manager_id from open_positions where id = ? ))", daily.Id)
		var managerEmail string
		querySubProjectManagerEmail.Scan(&managerEmail)
		//	fmt.Println(manager_email)

		if managerEmail == UserName {
			var id int
			queryExisting := db.QueryRow("select id from open_positions where id = ? and is_active = 1", daily.Id)
			queryExisting.Scan(&id)

			if id != 0 {
				queryUpdate, err := db.Prepare("Update open_positions set is_active = 0 where id = ? ")
				catch(err)
				_, err = queryUpdate.Exec(id)
				catch(err)
				w.WriteHeader(http.StatusOK)
				defer queryUpdate.Close()
				//fmt.Fprintf(w, "deleted successfully")
			} else {
				w.WriteHeader(http.StatusNotFound)
				//fmt.Fprintf(w, "Project doesnot exist")
			}
		} else {
			w.WriteHeader(http.StatusForbidden)
			//fmt.Fprintf(w, "project not under you")
		}
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		//fmt.Fprintf(w, "Unauthorised access")
	}
}

// Updateopenpositions :Func for updating entries
func (c *Commander) Updateopenpositions(w http.ResponseWriter, r *http.Request) {
	db := database.DbConn()
	defer db.Close()

	w.Header().Set("Content-Type", "application/json")

	if strings.Contains(Role, "Project Manager") || strings.Contains(Role, "project manager") {
		statusGet, ok := r.URL.Query()["status"]

		if !ok || statusGet[0] == "" {
			w.WriteHeader(http.StatusBadRequest)
			//fmt.Fprintf(w, "status is missing")
			return
		}
		status := statusGet[0]

		if status == "daily" || status == "Daily" {
			var daily models.Daily
			json.NewDecoder(r.Body).Decode(&daily)
			querySubProjectManagerID := db.QueryRow("select id from sub_project_manager where sub_project_id in ( select id from sub_project where sub_project_name = ? ) and manager_id in ( select id from project_manager where project_manager_email = ? )", daily.Project_name, UserName)
			var managerid int
			querySubProjectManagerID.Scan(&managerid)
			if managerid != 0 {
				var id int
				queryExisting := db.QueryRow("select id from open_positions where sub_project_manager_id = ? and designation = ? and is_active = 1 ", managerid, daily.Designation)
				queryExisting.Scan(&id)

				if id != 0 {
					queryUpdate, err := db.Prepare("Update open_positions set type_position = ?, position = ?, priority = ?, additional_comment = ?, l1_due = ?, l2_due = ?, client_due = ?, updated_at = now() where id = ?")
					catch(err)
					_, err = queryUpdate.Exec(daily.Type_position, daily.Position, daily.Priority, daily.Additonal_comment, daily.L1_due, daily.L2_due, daily.Client_due, id)
					catch(err)
					w.WriteHeader(http.StatusOK)
					defer queryUpdate.Close()
					//fmt.Fprintf(w, "update successfully")
				} else {
					w.WriteHeader(http.StatusNotFound)
					//fmt.Fprintf(w, "Project doesnot exist")
				}
			} else {
				w.WriteHeader(http.StatusForbidden)
				//fmt.Fprintf(w, "Project not under you")
			}
		} else if status == "Weekly" || status == "weekly" {
			var weekly models.Weekly
			json.NewDecoder(r.Body).Decode(&weekly)
			querySubProjectManagerID := db.QueryRow("select id from sub_project_manager where sub_project_id in ( select id from sub_project where sub_project_name = ? ) and manager_id in ( select id from project_manager where project_manager_email = ? )", weekly.Project_name, UserName)
			var managerid int
			querySubProjectManagerID.Scan(&managerid)

			if managerid != 0 {
				var id int
				queryExisting := db.QueryRow("select id from open_positions where sub_project_manager_id = ? and designation = ? and is_active = 1 ", managerid, weekly.Designation)
				queryExisting.Scan(&id)
				if id != 0 {
					queryUpdate, err := db.Prepare("Update open_positions set type_position = ?, position = ?, priority = ?, additional_comment = ?, l1_happened = ?, l2_happened = ?, client_happened = ?, updated_at = now() where id = ?")
					catch(err)
					_, err = queryUpdate.Exec(weekly.Type_position, weekly.Position, weekly.Priority, weekly.Additonal_comment, weekly.L1_Happened, weekly.L2_Happened, weekly.Client_Happened, id)
					catch(err)
					w.WriteHeader(http.StatusOK)
					defer queryUpdate.Close()
				} else {
					w.WriteHeader(http.StatusNotFound)
				}
			} else {
				w.WriteHeader(http.StatusForbidden)
			}
		}
	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}
}
