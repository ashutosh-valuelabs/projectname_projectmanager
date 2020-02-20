package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	database "projectname_projectmanager/driver"
	model "projectname_projectmanager/model"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql" //blank import
	"github.com/gorilla/mux"
)

var field string = "id,action_item,manager_details_id,meeting_date,target_date,status,closed_date,comment,is_active"
var fields string = "action_items.id,action_items.action_item,action_items.meeting_date,action_items.target_date,action_items.status,action_items.closed_date,action_items.comment,action_items.is_active,projects.project_name,manager.manager_name"

//ActionItemPostData : to post the data into the database
func (c *Commander) ActionItemPostData(w http.ResponseWriter, r *http.Request) {
	fmt.Println(UserName, Role)
	var error model.Error
	var actionitem model.ActionItemClosed
	var ManagerDetailsID int
	SetupResponse(&w, r)
	if (*r).Method == "OPTIONS" {
		w.Header().Set("Access-Control-Max-Age", "86400")
		w.WriteHeader(http.StatusOK)
		return
	}
	db := database.DbConn()
	defer db.Close()
	err := json.NewDecoder(r.Body).Decode(&actionitem)
	if err != nil {
		WriteLogFile(err)
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	if actionitem.Status != "open" && actionitem.Status != "inprogress" && actionitem.Status != "closed" && actionitem.Status != "onhold" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		//error.Code = "1265"
		error.Message = "Data truncated for column 'status'. Invalid entry for column 'status'"
		json.NewEncoder(w).Encode(error)
		return
	}
	myDateString := "2006-01-02"
	myMeetingDate, err := time.Parse(myDateString, actionitem.MeetingDate)
	if err != nil {
		WriteLogFile(err)
		panic(err)
	}
	myTargetDate, err := time.Parse(myDateString, actionitem.TargetDate)
	if err != nil {
		WriteLogFile(err)
		panic(err)
	}
	isTargetDateBeforeMeetingDate := myTargetDate.Before(myMeetingDate)
	if isTargetDateBeforeMeetingDate == false {
		var ActionItem, MeetingDate, TargetDate, Status, Comment string
		selDB, err := db.Query("SELECT id FROM sub_project_manager WHERE sub_project_id = (SELECT id FROM sub_project WHERE sub_project_name=? )and manager_id = (SELECT id FROM project_manager WHERE project_manager_email =?)", actionitem.ProjectName, UserName)
		defer selDB.Close()
		if err != nil {
			WriteLogFile(err)
			// If the structure of the body is wrong, return an HTTP error
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
		if selDB.Next() != false {
			err := selDB.Scan(&ManagerDetailsID)
			if err != nil {
				WriteLogFile(err)
				// If the structure of the body is wrong, return an HTTP error
				w.WriteHeader(http.StatusUnprocessableEntity)
				return
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
		selDB, err = db.Query("SELECT action_item,meeting_date,target_date,status,comment from action_items WHERE sub_project_manager_id=? AND action_item=? AND is_active='1'", ManagerDetailsID, actionitem.ActionItem)
		defer selDB.Close()
		if err != nil {
			WriteLogFile(err)
			// If the structure of the body is wrong, return an HTTP error
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
		for selDB.Next() {
			err := selDB.Scan(&ActionItem, &MeetingDate, &TargetDate, &Status, &Comment)
			if err != nil {
				WriteLogFile(err)
				// If the structure of the body is wrong, return an HTTP error
				w.WriteHeader(http.StatusUnprocessableEntity)
				return
			}
			if ActionItem == actionitem.ActionItem &&
				MeetingDate == actionitem.MeetingDate &&
				TargetDate == actionitem.TargetDate &&
				Status == actionitem.Status &&
				Comment == actionitem.Comment {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusConflict)
				//error.Code = "2627"
				error.Message = "Violation of Unique constraint: Duplicate row insertion"
				json.NewEncoder(w).Encode(error)
				return
			}
		}
		if err != nil {
			WriteLogFile(err)
			// If the structure of the body is wrong, return an HTTP error
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
		CreatedAt := time.Now()
		insForm, err := db.Prepare("INSERT INTO action_items(action_item,sub_project_manager_id,meeting_date,target_date,status,comment,created_at,updated_at)VALUES(?,?,?,?,?,?,?,?)")
		defer insForm.Close()
		if err != nil {
			WriteLogFile(err)
			return
		}
		insForm.Exec(actionitem.ActionItem,
			ManagerDetailsID,
			actionitem.MeetingDate,
			actionitem.TargetDate,
			actionitem.Status,
			actionitem.Comment,
			CreatedAt,
			CreatedAt)
		if err != nil {
			WriteLogFile(err)
			// If the structure of the body is wrong, return an HTTP error
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		//error.Code = "1292"
		if isTargetDateBeforeMeetingDate == true {
			error.Message = "Incorrect Target Date Value"
		}
		json.NewEncoder(w).Encode(error)
	}
}

//ActionItemUpdateData : to update the details stored in database
func (c *Commander) ActionItemUpdateData(w http.ResponseWriter, r *http.Request) {
	var error model.Error
	var actionitem model.ActionItemClosed
	//var closedInTime, managerDetailsID int
	SetupResponse(&w, r)
	if (*r).Method == "OPTIONS" {
		w.Header().Set("Access-Control-Max-Age", "86400")
		w.WriteHeader(http.StatusOK)
	}
	db := database.DbConn()
	defer func() {
		err := db.Close()
		if err != nil {
			WriteLogFile(err)
			return
		}
	}()
	err := json.NewDecoder(r.Body).Decode(&actionitem)
	if err != nil {
		WriteLogFile(err)
		//If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	updateOk, updateComment, err := UpdateData(actionitem, UserName)
	if err != nil {
		WriteLogFile(err)
	}
	if updateComment == "" {
		if updateOk == http.StatusUnprocessableEntity {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		} else if updateOk == http.StatusBadRequest {
			w.WriteHeader(http.StatusBadRequest)
			return
		} else if updateOk == http.StatusForbidden {
			w.WriteHeader(http.StatusForbidden)
			return
		} else {
			w.WriteHeader(http.StatusCreated)
			return
		}
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		error.Message = updateComment
		json.NewEncoder(w).Encode(error)
		return
	}
	// myDateString := "2006-01-02"
	// myMeetingDate, err := time.Parse(myDateString, actionitem.MeetingDate)
	// if err != nil {
	// 	WriteLogFile(err)
	// 	panic(err)
	// }
	// myTargetDate, err := time.Parse(myDateString, actionitem.TargetDate)
	// if err != nil {
	// 	WriteLogFile(err)
	// 	panic(err)
	// }
	// isTargetDateBeforeMeetingDate := myTargetDate.Before(myMeetingDate)
	// if actionitem.ClosedDate != "" {
	// 	if actionitem.Status == "closed" {
	// 		myClosedDate, err := time.Parse(myDateString, actionitem.ClosedDate)
	// 		if err != nil {
	// 			WriteLogFile(err)
	// 			panic(err)
	// 		}
	// 		isClosedDateBeforeMeetingDate := myClosedDate.Before(myMeetingDate)
	// 		if isClosedDateBeforeMeetingDate == false && isTargetDateBeforeMeetingDate == false {
	// 			isClosedDateAfterTargetDate := myClosedDate.After(myTargetDate)
	// 			if isClosedDateAfterTargetDate == false {
	// 				closedInTime = 1
	// 			} else {
	// 				closedInTime = 0
	// 			}
	// 			selDB, err := db.Query("SELECT id FROM sub_project_manager WHERE sub_project_id = (SELECT id FROM sub_project WHERE sub_project_name=? )and manager_id = (SELECT id FROM project_manager WHERE project_manager_email =?)", actionitem.ProjectName, User)
	// 			defer selDB.Close()
	// 			if err != nil {
	// 				WriteLogFile(err)
	// 				// If the structure of the body is wrong, return an HTTP error
	// 				w.WriteHeader(http.StatusUnprocessableEntity)
	// 				return
	// 			}
	// 			if selDB.Next() != false {
	// 				err := selDB.Scan(&managerDetailsID)
	// 				if err != nil {
	// 					WriteLogFile(err)
	// 					// If the structure of the body is wrong, return an HTTP error
	// 					w.WriteHeader(http.StatusUnprocessableEntity)
	// 					return
	// 				}
	// 			} else {
	// 				w.WriteHeader(http.StatusBadRequest)
	// 			}

	// 			updatedAt := time.Now()
	// 			updForm, err := db.Prepare("UPDATE action_items SET action_item=? , sub_project_manager_id=? , meeting_date=? , target_date=? , status=? , closed_date=? , comment=? , updated_at=? , closed_in_time=? WHERE id=?")
	// 			if err != nil {
	// 				WriteLogFile(err)
	// 				panic(err.Error())
	// 			}
	// 			defer updForm.Close()
	// 			updForm.Exec(actionitem.ActionItem,
	// 				managerDetailsID,
	// 				actionitem.MeetingDate,
	// 				actionitem.TargetDate,
	// 				actionitem.Status,
	// 				actionitem.ClosedDate,
	// 				actionitem.Comment,
	// 				updatedAt,
	// 				closedInTime,
	// 				actionitem.SNo)
	// 			defer updForm.Close()
	// 			w.Header().Set("Content-Type", "application/json")
	// 			w.WriteHeader(http.StatusCreated)
	// 		} else {
	// 			w.Header().Set("Content-Type", "application/json")
	// 			w.WriteHeader(http.StatusUnprocessableEntity)
	// 			//error.Code = "1292"
	// 			if isClosedDateBeforeMeetingDate == true && isTargetDateBeforeMeetingDate == true {
	// 				error.Message = "Incorrect Closed Date and Target Date Values"
	// 			} else if isTargetDateBeforeMeetingDate == true {
	// 				error.Message = "Incorrect Target Date Value"
	// 			} else if isClosedDateBeforeMeetingDate == true {
	// 				error.Message = "Incorrect Closed Date Value"
	// 			}
	// 			json.NewEncoder(w).Encode(error)
	// 		}
	// 	} else {
	// 		w.Header().Set("Content-Type", "application/json")
	// 		w.WriteHeader(http.StatusBadRequest)
	// 		//error.Code = "1293"
	// 		error.Message = "Closed date cannot be filled if the status is not closed"
	// 		json.NewEncoder(w).Encode(error)
	// 	}
	// } else {
	// 	if actionitem.Status != "closed" {
	// 		if isTargetDateBeforeMeetingDate == false {
	// 			selDB, err := db.Query("SELECT id FROM sub_project_manager WHERE sub_project_id = (SELECT id FROM sub_project WHERE sub_project_name=? )and manager_id = (SELECT id FROM project_manager WHERE project_manager_email =?)", actionitem.ProjectName, User)
	// 			defer selDB.Close()
	// 			if err != nil {
	// 				WriteLogFile(err)
	// 				// If the structure of the body is wrong, return an HTTP error
	// 				w.WriteHeader(http.StatusUnprocessableEntity)
	// 				return
	// 			}
	// 			if selDB.Next() != false {
	// 				err := selDB.Scan(&managerDetailsID)
	// 				if err != nil {
	// 					WriteLogFile(err)
	// 					// If the structure of the body is wrong, return an HTTP error
	// 					w.WriteHeader(http.StatusUnprocessableEntity)
	// 					return
	// 				}
	// 			} else {
	// 				w.WriteHeader(http.StatusBadRequest)
	// 			}
	// 			updatedAt := time.Now()
	// 			updForm, err := db.Prepare("UPDATE action_items SET action_item=? , sub_project_manager_id=? , meeting_date=? , target_date=? , status=? , comment=? , updated_at=? , closed_in_time = NULL WHERE id=?")
	// 			if err != nil {
	// 				WriteLogFile(err)
	// 				panic(err.Error())
	// 			}
	// 			updForm.Exec(actionitem.ActionItem,
	// 				managerDetailsID,
	// 				actionitem.MeetingDate,
	// 				actionitem.TargetDate,
	// 				actionitem.Status,
	// 				actionitem.Comment,
	// 				updatedAt,
	// 				actionitem.SNo)
	// 			defer updForm.Close()
	// 			w.Header().Set("Content-Type", "application/json")
	// 			w.WriteHeader(http.StatusCreated)
	// 		} else {
	// 			w.Header().Set("Content-Type", "application/json")
	// 			w.WriteHeader(http.StatusUnprocessableEntity)
	// 			//error.Code = "1292"
	// 			if isTargetDateBeforeMeetingDate == true {
	// 				error.Message = "Incorrect Target Date Value"
	// 				json.NewEncoder(w).Encode(error)
	// 			}
	// 		}
	// 	} else {
	// 		w.Header().Set("Content-Type", "application/json")
	// 		w.WriteHeader(http.StatusUnprocessableEntity)
	// 		error.Message = "Closed date missing"
	// 		json.NewEncoder(w).Encode(error)

	// 	}
	// }
}

//ActionItemDeleteData : to soft delete the details
func (c *Commander) ActionItemDeleteData(w http.ResponseWriter, r *http.Request) {
	var actionitem model.ActionItemClosed
	//var managerDetailsID int
	//var managerEmail string
	SetupResponse(&w, r)
	if (*r).Method == "OPTIONS" {
		w.Header().Set("Access-Control-Max-Age", "86400")
		w.WriteHeader(http.StatusOK)
		return
	}
	//db := database.DbConn()
	//defer db.Close()
	err := json.NewDecoder(r.Body).Decode(&actionitem)

	if err != nil {
		WriteLogFile(err)
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	sno := actionitem.ID
	deleteOk, err := DeleteData(sno, UserName)
	if err != nil {
		if err != nil {
			WriteLogFile(err)
			// If the structure of the body is wrong, return an HTTP error
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
	} else {
		if deleteOk == http.StatusBadRequest {
			w.WriteHeader(http.StatusBadRequest)
		} else if deleteOk == http.StatusOK {
			w.WriteHeader(http.StatusOK)
		} else if deleteOk == http.StatusForbidden {
			w.WriteHeader(http.StatusForbidden)
		}
	}
	// selDB, err := db.Query("SELECT manager_project_id from action_items WHERE id =? and is_active='1'", actionitem.SNo)
	// defer selDB.Close()
	// BadRequest(w, err)
	// if selDB.Next() != false {
	// 	err := selDB.Scan(&managerDetailsID)
	// 	BadRequest(w, err)
	// } else {
	// 	w.WriteHeader(http.StatusBadRequest)
	// }
	// selDB, err = db.Query("SELECT manager.manager_email_id FROM manager JOIN manager_project ON manager.id = manager_project.manager_id WHERE manager_project.id =?", managerDetailsID)
	// defer selDB.Close()
	// BadRequest(w, err)
	// if selDB.Next() != false {
	// 	err := selDB.Scan(&managerEmail)
	// 	BadRequest(w, err)
	// } else {
	// 	w.WriteHeader(http.StatusBadRequest)
	// }
	// if managerEmail == UserName {
	// 	updatedAt := time.Now()
	// 	updForm, err := db.Prepare("UPDATE action_items SET is_active = '0', updated_at = ? WHERE id=?")
	// 	if err != nil {
	// 		WriteLogFile(err)
	// 		panic(err.Error())
	// 	}
	// 	defer updForm.Close()
	// 	updForm.Exec(updatedAt, actionitem.SNo)
	// 	w.WriteHeader(http.StatusOK)
	// } else {
	// 	w.WriteHeader(http.StatusForbidden)
	// }
}

//ActionItemGetData : to the get the action items
func (c *Commander) ActionItemGetData(w http.ResponseWriter, r *http.Request) {
	var data []model.ActionItemClosed
	var Page model.Pagination
	var countt int
	var error model.Error
	SetupResponse(&w, r)
	if (*r).Method == "OPTIONS" {
		w.Header().Set("Access-Control-Max-Age", "86400")
		w.WriteHeader(http.StatusOK)
		return
	}
	// cacheKey := "actionclosed"
	// e := `"` + cacheKey + `"`
	// w.Header().Set("Etag", e)
	// w.Header().Set("Cache-Control", "max-age=2592000") // 30 days

	// if match := r.Header.Get("If-None-Match"); match != "" {
	// 	if strings.Contains(match, e) {
	// 		w.WriteHeader(http.StatusNotModified)
	// 		return
	// 	}
	// }
	db := database.DbConn()
	defer db.Close()
	//pages := r.FormValue("pages")
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		WriteLogFile(err)
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	if limit != 10 && limit != 20 && limit != 50 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		error.Message = "Incorrect Limit Value"
		json.NewEncoder(w).Encode(error)
		return
	}
	pages := r.URL.Query().Get("pages")
	page, err := strconv.Atoi(pages)
	Page.Page = page + 1
	if err != nil {
		WriteLogFile(err)
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	offset := page * limit
	status := r.URL.Query().Get("status")
	if strings.Contains(Role, "project manager") == true {
		if status == "closed" {
			var SNo int
			var ProjectName, ActionItem, MeetingDate, TargetDate, Status, ClosedDate, Comment, Flag string
			count, err := db.Query("SELECT count(*) FROM action_items a JOIN sub_project_manager spm ON spm.id = a.sub_project_manager_id JOIN sub_project s ON spm.sub_project_id = s.id JOIN project_manager pj ON pj.id = spm.manager_id WHERE pj.project_manager_email=? AND a.status= 3 AND a.is_active='1'", UserName)
			if err != nil {
				WriteLogFile(err)
				// If the structure of the body is wrong, return an HTTP error
				w.WriteHeader(http.StatusUnprocessableEntity)
				return
			}
			defer count.Close()
			for count.Next() {
				err := count.Scan(&countt)
				if err != nil {
					WriteLogFile(err)
					// If the structure of the body is wrong, return an HTTP error
					w.WriteHeader(http.StatusUnprocessableEntity)
					return
				}
			}
			selDB, err := db.Query("call getClosedActionsProject(?,?,?)", UserName, offset, limit)
			if err != nil {
				WriteLogFile(err)
				// If the structure of the body is wrong, return an HTTP error
				w.WriteHeader(http.StatusUnprocessableEntity)
				return
			}
			defer selDB.Close()
			for selDB.Next() {
				err := selDB.Scan(&SNo, &ActionItem, &MeetingDate, &TargetDate, &Status, &ClosedDate, &Comment, &Flag, &ProjectName)
				if err != nil {
					WriteLogFile(err)
					// If the structure of the body is wrong, return an HTTP error
					w.WriteHeader(http.StatusUnprocessableEntity)
					return
				}
				data = append(data, model.ActionItemClosed{ID: SNo, ProjectName: ProjectName, ActionItem: ActionItem, MeetingDate: MeetingDate, TargetDate: TargetDate, Status: Status, ClosedDate: ClosedDate, Comment: Comment})
			}
		} else if status == "open" {
			var SNo int
			var ProjectName, ActionItem, MeetingDate, TargetDate, Status, Comment, Flag string
			var ClosedDate sql.NullString
			count, err := db.Query("SELECT count(*) FROM action_items a JOIN sub_project_manager spm ON spm.id = a.sub_project_manager_id JOIN sub_project s ON spm.sub_project_id = s.id JOIN project_manager pj ON pj.id = spm.manager_id WHERE pj.project_manager_email=? AND a.status!= 3 AND a.is_active='1'", UserName)
			if err != nil {
				WriteLogFile(err)
				// If the structure of the body is wrong, return an HTTP error
				w.WriteHeader(http.StatusUnprocessableEntity)
				return
			}
			defer count.Close()
			for count.Next() {
				err := count.Scan(&countt)
				if err != nil {
					WriteLogFile(err)
					// If the structure of the body is wrong, return an HTTP error
					w.WriteHeader(http.StatusUnprocessableEntity)
					return
				}
			}
			selDB, err := db.Query("call getAllActionsProject(?,?,?)", UserName, offset, limit)
			if err != nil {
				WriteLogFile(err)
				// If the structure of the body is wrong, return an HTTP error
				w.WriteHeader(http.StatusUnprocessableEntity)
				return
			}
			defer selDB.Close()
			for selDB.Next() {
				err := selDB.Scan(&SNo, &ActionItem, &MeetingDate, &TargetDate, &Status, &ClosedDate, &Comment, &Flag, &ProjectName)
				if err != nil {
					WriteLogFile(err)
					// If the structure of the body is wrong, return an HTTP error
					w.WriteHeader(http.StatusUnprocessableEntity)
					return
				}
				data = append(data, model.ActionItemClosed{ID: SNo, ProjectName: ProjectName, ActionItem: ActionItem, MeetingDate: MeetingDate, TargetDate: TargetDate, Status: Status, ClosedDate: "NA", Comment: Comment})
			}
		} else if status == "" {
			var SNo int
			var ProjectName, ActionItem, MeetingDate, TargetDate, Status, Comment, Flag string
			var ClosedDate sql.NullString
			count, err := db.Query("SELECT count(*) FROM action_items a JOIN sub_project_manager spm ON spm.id = a.sub_project_manager_id JOIN sub_project s ON spm.sub_project_id = s.id JOIN project_manager pj ON pj.id = spm.manager_id WHERE pj.project_manager_email=? AND a.is_active='1'", UserName)
			if err != nil {
				WriteLogFile(err)
				// If the structure of the body is wrong, return an HTTP error
				w.WriteHeader(http.StatusUnprocessableEntity)
				return
			}
			defer count.Close()
			for count.Next() {
				err := count.Scan(&countt)
				if err != nil {
					WriteLogFile(err)
					// If the structure of the body is wrong, return an HTTP error
					w.WriteHeader(http.StatusUnprocessableEntity)
					return
				}
			}
			selDB, err := db.Query("call getActionsProject(?,?,?)", UserName, offset, limit)
			if err != nil {
				WriteLogFile(err)
				// If the structure of the body is wrong, return an HTTP error
				w.WriteHeader(http.StatusUnprocessableEntity)
				return
			}
			defer selDB.Close()
			for selDB.Next() {
				err := selDB.Scan(&SNo, &ActionItem, &MeetingDate, &TargetDate, &Status, &ClosedDate, &Comment, &Flag, &ProjectName)
				if err != nil {
					WriteLogFile(err)
					// If the structure of the body is wrong, return an HTTP error
					w.WriteHeader(http.StatusUnprocessableEntity)
					return
				}
				if ClosedDate.Valid == false {
					data = append(data, model.ActionItemClosed{ID: SNo, ProjectName: ProjectName, ActionItem: ActionItem, MeetingDate: MeetingDate, TargetDate: TargetDate, Status: Status, ClosedDate: "NA", Comment: Comment})
				} else {
					data = append(data, model.ActionItemClosed{ID: SNo, ProjectName: ProjectName, ActionItem: ActionItem, MeetingDate: MeetingDate, TargetDate: TargetDate, Status: Status, ClosedDate: ClosedDate.String, Comment: Comment})
				}
			}

		} else {
			error.Code = "405"
			error.Message = "Method not allowed"
			json.NewEncoder(w).Encode(error)

		}
	} else {
		if status == "closed" {
			var SNo int
			var ProjectName, ActionItem, Owner, MeetingDate, TargetDate, Status, ClosedDate, Comment, Flag string
			count, err := db.Query("SELECT count(*) FROM action_items a JOIN sub_project_manager spm ON spm.id = a.sub_project_manager_id JOIN sub_project s ON spm.sub_project_id = s.id JOIN project p ON p.id = s.project_id JOIN program_manager pg ON pg.id = p.program_manager_id JOIN project_manager pj ON pj.id = spm.manager_id WHERE pg.program_manager_email=? AND a.status= 3 AND a.is_active='1'", UserName)
			if err != nil {
				WriteLogFile(err)
				// If the structure of the body is wrong, return an HTTP error
				w.WriteHeader(http.StatusUnprocessableEntity)
				return
			}
			defer count.Close()
			for count.Next() {
				err := count.Scan(&countt)
				if err != nil {
					WriteLogFile(err)
					// If the structure of the body is wrong, return an HTTP error
					w.WriteHeader(http.StatusUnprocessableEntity)
					return
				}
			}
			selDB, err := db.Query("call getClosedActionsProgram(?,?,?)", UserName, offset, limit)
			defer selDB.Close()
			if err != nil {
				WriteLogFile(err)
				// If the structure of the body is wrong, return an HTTP error
				w.WriteHeader(http.StatusUnprocessableEntity)
				return
			}
			for selDB.Next() {
				err := selDB.Scan(&SNo, &ActionItem, &MeetingDate, &TargetDate, &Status, &ClosedDate, &Comment, &Flag, &ProjectName, &Owner)
				if err != nil {
					WriteLogFile(err)
					// If the structure of the body is wrong, return an HTTP error
					w.WriteHeader(http.StatusUnprocessableEntity)
					return
				}
				data = append(data, model.ActionItemClosed{ID: SNo, ProjectName: ProjectName, ActionItem: ActionItem, Owner: Owner, MeetingDate: MeetingDate, TargetDate: TargetDate, Status: Status, ClosedDate: ClosedDate, Comment: Comment})
			}
		} else if status == "open" {
			var SNo int
			var ProjectName, ActionItem, Owner, MeetingDate, TargetDate, Status, Comment, Flag string
			var ClosedDate sql.NullString
			count, err := db.Query("SELECT count(*) FROM action_items a JOIN sub_project_manager spm ON spm.id = a.sub_project_manager_id JOIN sub_project s ON spm.sub_project_id = s.id JOIN project p ON p.id = s.project_id JOIN program_manager pg ON pg.id = p.program_manager_id JOIN project_manager pj ON pj.id = spm.manager_id WHERE pg.program_manager_email=? AND a.status!= 3 AND a.is_active='1'", UserName)
			if err != nil {
				WriteLogFile(err)
				// If the structure of the body is wrong, return an HTTP error
				w.WriteHeader(http.StatusUnprocessableEntity)
				return
			}
			defer count.Close()
			for count.Next() {
				err := count.Scan(&countt)
				if err != nil {
					WriteLogFile(err)
					// If the structure of the body is wrong, return an HTTP error
					w.WriteHeader(http.StatusUnprocessableEntity)
					return
				}
			}
			selDB, err := db.Query("call getAllActionsProgram(?,?,?)", UserName, offset, limit)
			if err != nil {
				WriteLogFile(err)
				// If the structure of the body is wrong, return an HTTP error
				w.WriteHeader(http.StatusUnprocessableEntity)
				return
			}
			defer selDB.Close()
			for selDB.Next() {
				err := selDB.Scan(&SNo, &ActionItem, &MeetingDate, &TargetDate, &Status, &ClosedDate, &Comment, &Flag, &ProjectName, &Owner)
				if err != nil {
					WriteLogFile(err)
					// If the structure of the body is wrong, return an HTTP error
					w.WriteHeader(http.StatusUnprocessableEntity)
					return
				}
				data = append(data, model.ActionItemClosed{ID: SNo, ProjectName: ProjectName, ActionItem: ActionItem, Owner: Owner, MeetingDate: MeetingDate, TargetDate: TargetDate, Status: Status, ClosedDate: "NA", Comment: Comment})
			}
		} else if status == "" {
			var SNo int
			var ProjectName, ActionItem, MeetingDate, TargetDate, Owner, Status, Comment, Flag string
			var ClosedDate sql.NullString
			count, err := db.Query("SELECT count(*) FROM action_items a JOIN sub_project_manager spm ON spm.id = a.sub_project_manager_id JOIN sub_project s ON spm.sub_project_id = s.id JOIN project p ON p.id = s.project_id JOIN program_manager pg ON pg.id = p.program_manager_id JOIN project_manager pj ON pj.id = spm.manager_id WHERE pg.program_manager_email=? AND a.is_active='1'", UserName)
			if err != nil {
				WriteLogFile(err)
				// If the structure of the body is wrong, return an HTTP error
				w.WriteHeader(http.StatusUnprocessableEntity)
				return
			}
			defer count.Close()
			for count.Next() {
				err := count.Scan(&countt)
				if err != nil {
					WriteLogFile(err)
					// If the structure of the body is wrong, return an HTTP error
					w.WriteHeader(http.StatusUnprocessableEntity)
					return
				}
			}
			selDB, err := db.Query("call getActionsProgram(?,?,?)", UserName, offset, limit)
			if err != nil {
				WriteLogFile(err)
				// If the structure of the body is wrong, return an HTTP error
				w.WriteHeader(http.StatusUnprocessableEntity)
				return
			}
			defer selDB.Close()
			for selDB.Next() {
				err := selDB.Scan(&SNo, &ActionItem, &MeetingDate, &TargetDate, &Status, &ClosedDate, &Comment, &Flag, &ProjectName, &Owner)
				if err != nil {
					WriteLogFile(err)
					// If the structure of the body is wrong, return an HTTP error
					w.WriteHeader(http.StatusUnprocessableEntity)
					return
				}
				if ClosedDate.Valid == false {
					data = append(data, model.ActionItemClosed{ID: SNo, ProjectName: ProjectName, ActionItem: ActionItem, Owner: Owner, MeetingDate: MeetingDate, TargetDate: TargetDate, Status: Status, ClosedDate: "NA", Comment: Comment})
				} else {
					data = append(data, model.ActionItemClosed{ID: SNo, ProjectName: ProjectName, ActionItem: ActionItem, Owner: Owner, MeetingDate: MeetingDate, TargetDate: TargetDate, Status: Status, ClosedDate: ClosedDate.String, Comment: Comment})
				}
			}

		} else {
			error.Code = "405"
			error.Message = "Method not allowed"
			json.NewEncoder(w).Encode(error)
		}

	}
	w.Header().Set("Content-Type", "application/json")
	Page.TotalData = countt
	Page.Limit = limit
	x := countt
	page = x / limit
	x = x % limit
	if x == 0 {
		Page.TotalPages = page
	} else {
		Page.TotalPages = page + 1
	}
	Page.Data = data
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Page)
}

//ActionItemGetDataID : to search the action items
func (c *Commander) ActionItemGetDataID(w http.ResponseWriter, r *http.Request) {
	var data []model.ActionItemClosed
	var Page model.Pagination
	var countt int
	var error model.Error
	SetupResponse(&w, r)
	if (*r).Method == "OPTIONS" {
		w.Header().Set("Access-Control-Max-Age", "86400")
		w.WriteHeader(http.StatusOK)
		return
	}
	db := database.DbConn()
	defer db.Close()
	// cacheKey := "actionclosedID"
	// e := `"` + cacheKey + `"`
	// w.Header().Set("Etag", e)
	// w.Header().Set("Cache-Control", "max-age=2592000") // 30 days

	// if match := r.Header.Get("If-None-Match"); match != "" {
	// 	if strings.Contains(match, e) {
	// 		w.WriteHeader(http.StatusNotModified)
	// 		return
	// 	}
	// }
	//pages := r.FormValue("pages")
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		WriteLogFile(err)
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	if limit != 10 && limit != 20 && limit != 50 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		error.Message = "Incorrect Limit Value"
		json.NewEncoder(w).Encode(error)
		return
	}
	pages := r.URL.Query().Get("pages")
	page, err := strconv.Atoi(pages)
	Page.Page = page + 1
	if err != nil {
		WriteLogFile(err)
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	offset := page * limit
	status := r.URL.Query().Get("status")
	p := mux.Vars(r)
	removeWhiteSpace := p["id"]
	removeWhiteSpace = strings.TrimSpace(removeWhiteSpace)
	key := removeWhiteSpace + "%"
	fmt.Println(key)
	if strings.Contains(Role, "project manager") == true {
		if status == "closed" {
			var SNo int
			var ProjectName, ActionItem, MeetingDate, TargetDate, Status, ClosedDate, Comment, Flag string
			count, err := db.Query("SELECT count(*) FROM action_items a JOIN sub_project_manager spm ON spm.id = a.sub_project_manager_id JOIN sub_project s ON spm.sub_project_id = s.id JOIN project_manager pj ON pj.id = spm.manager_id WHERE (a.status = 3 AND pj.project_manager_email= ? AND a.is_active='1') AND (a.action_item like '"+key+"' OR a.meeting_date like '"+key+"' OR a.target_date like '"+key+"' OR a.closed_date like '"+key+"' OR a.comment like '"+key+"' OR s.sub_project_name like '"+key+"')", UserName)
			if err != nil {
				WriteLogFile(err)
				// If the structure of the body is wrong, return an HTTP error
				w.WriteHeader(http.StatusUnprocessableEntity)
				return
			}
			defer count.Close()
			for count.Next() {
				err := count.Scan(&countt)
				if err != nil {
					WriteLogFile(err)
					// If the structure of the body is wrong, return an HTTP error
					w.WriteHeader(http.StatusUnprocessableEntity)
					return
				}
			}
			selDB, err := db.Query("call getClosedActionsProjectQuery(?,?,?,?)", key, UserName, offset, limit)
			if err != nil {
				WriteLogFile(err)
				// If the structure of the body is wrong, return an HTTP error
				w.WriteHeader(http.StatusUnprocessableEntity)
				return
			}
			defer selDB.Close()
			for selDB.Next() {
				err := selDB.Scan(&SNo, &ActionItem, &MeetingDate, &TargetDate, &Status, &ClosedDate, &Comment, &Flag, &ProjectName)
				if err != nil {
					WriteLogFile(err)
					// If the structure of the body is wrong, return an HTTP error
					w.WriteHeader(http.StatusUnprocessableEntity)
					return
				}
				data = append(data, model.ActionItemClosed{ID: SNo, ProjectName: ProjectName, ActionItem: ActionItem, MeetingDate: MeetingDate, TargetDate: TargetDate, Status: Status, ClosedDate: ClosedDate, Comment: Comment})
			}
		} else if status == "open" {
			var SNo int
			var ProjectName, ActionItem, MeetingDate, TargetDate, Status, Comment, Flag string
			var ClosedDate sql.NullString
			count, err := db.Query("SELECT count(*) FROM action_items a JOIN sub_project_manager spm ON spm.id = a.sub_project_manager_id JOIN sub_project s ON spm.sub_project_id = s.id JOIN project_manager pj ON pj.id = spm.manager_id WHERE (a.status != 3 AND pj.project_manager_email= ? AND a.is_active='1') AND (a.action_item like '"+key+"' OR a.meeting_date like '"+key+"' OR a.target_date like '"+key+"' OR a.closed_date like '"+key+"' OR a.comment like '"+key+"' OR s.sub_project_name like '"+key+"')", UserName)
			if err != nil {
				WriteLogFile(err)
				// If the structure of the body is wrong, return an HTTP error
				w.WriteHeader(http.StatusUnprocessableEntity)
				return
			}
			defer count.Close()
			for count.Next() {
				err := count.Scan(&countt)
				if err != nil {
					WriteLogFile(err)
					// If the structure of the body is wrong, return an HTTP error
					w.WriteHeader(http.StatusUnprocessableEntity)
					return
				}
			}
			selDB, err := db.Query("call getAllActionsProjectQuery(?,?,?,?)", key, UserName, offset, limit)
			if err != nil {
				WriteLogFile(err)
				// If the structure of the body is wrong, return an HTTP error
				w.WriteHeader(http.StatusUnprocessableEntity)
				return
			}
			defer selDB.Close()
			for selDB.Next() {
				err := selDB.Scan(&SNo, &ActionItem, &MeetingDate, &TargetDate, &Status, &ClosedDate, &Comment, &Flag, &ProjectName)
				if err != nil {
					WriteLogFile(err)
					// If the structure of the body is wrong, return an HTTP error
					w.WriteHeader(http.StatusUnprocessableEntity)
					return
				}
				data = append(data, model.ActionItemClosed{ID: SNo, ProjectName: ProjectName, ActionItem: ActionItem, MeetingDate: MeetingDate, TargetDate: TargetDate, Status: Status, ClosedDate: "NA", Comment: Comment})
			}
		} else if status == "" {
			var SNo int
			var ProjectName, ActionItem, MeetingDate, TargetDate, Status, Comment, Flag string
			var ClosedDate sql.NullString
			count, err := db.Query("SELECT count(*) FROM action_items a JOIN sub_project_manager spm ON spm.id = a.sub_project_manager_id JOIN sub_project s ON spm.sub_project_id = s.id JOIN project_manager pj ON pj.id = spm.manager_id WHERE (pj.project_manager_email= ? AND a.is_active='1') AND (a.action_item like '"+key+"' OR a.meeting_date like '"+key+"' OR a.target_date like '"+key+"' OR a.closed_date like '"+key+"' OR a.comment like '"+key+"' OR s.sub_project_name like '"+key+"')", UserName)
			if err != nil {
				WriteLogFile(err)
				// If the structure of the body is wrong, return an HTTP error
				w.WriteHeader(http.StatusUnprocessableEntity)
				return
			}
			defer count.Close()
			for count.Next() {
				err := count.Scan(&countt)
				if err != nil {
					WriteLogFile(err)
					// If the structure of the body is wrong, return an HTTP error
					w.WriteHeader(http.StatusUnprocessableEntity)
					return
				}
			}
			selDB, err := db.Query("call getActionsProjectQuery(?,?,?,?)", key, UserName, offset, limit)
			if err != nil {
				WriteLogFile(err)
				// If the structure of the body is wrong, return an HTTP error
				w.WriteHeader(http.StatusUnprocessableEntity)
				return
			}
			defer selDB.Close()
			for selDB.Next() {
				err := selDB.Scan(&SNo, &ActionItem, &MeetingDate, &TargetDate, &Status, &ClosedDate, &Comment, &Flag, &ProjectName)
				if err != nil {
					WriteLogFile(err)
					// If the structure of the body is wrong, return an HTTP error
					w.WriteHeader(http.StatusUnprocessableEntity)
					return
				}
				if ClosedDate.Valid == false {
					data = append(data, model.ActionItemClosed{ID: SNo, ProjectName: ProjectName, ActionItem: ActionItem, MeetingDate: MeetingDate, TargetDate: TargetDate, Status: Status, ClosedDate: "NA", Comment: Comment})
				} else {
					data = append(data, model.ActionItemClosed{ID: SNo, ProjectName: ProjectName, ActionItem: ActionItem, MeetingDate: MeetingDate, TargetDate: TargetDate, Status: Status, ClosedDate: ClosedDate.String, Comment: Comment})
				}
			}

		} else {
			error.Code = "405"
			error.Message = "Method not allowed"
			json.NewEncoder(w).Encode(error)
		}
	} else {
		if status == "closed" {
			var SNo int
			var ProjectName, ActionItem, Owner, MeetingDate, TargetDate, Status, ClosedDate, Comment, Flag string
			count, err := db.Query("SELECT count(*) FROM action_items a JOIN sub_project_manager spm ON spm.id = a.sub_project_manager_id JOIN sub_project s ON spm.sub_project_id = s.id JOIN project p ON p.id = s.project_id JOIN program_manager pg ON pg.id = p.program_manager_id JOIN project_manager pj ON pj.id = spm.manager_id WHERE (a.status = 3 AND pg.program_manager_email= ? AND a.is_active='1') AND (a.action_item like '"+key+"' OR a.meeting_date like '"+key+"' OR a.target_date like '"+key+"' OR a.status like '"+key+"' OR a.closed_date like '"+key+"' OR a.comment like '"+key+"' OR s.sub_project_name like '"+key+"' OR pj.project_manager_name like '"+key+"')", UserName)
			if err != nil {
				WriteLogFile(err)
				// If the structure of the body is wrong, return an HTTP error
				w.WriteHeader(http.StatusUnprocessableEntity)
				return
			}
			defer count.Close()
			for count.Next() {
				err := count.Scan(&countt)
				if err != nil {
					WriteLogFile(err)
					// If the structure of the body is wrong, return an HTTP error
					w.WriteHeader(http.StatusUnprocessableEntity)
					return
				}
			}
			selDB, err := db.Query("call getClosedActionsProgramQuery(?,?,?,?)", key, UserName, offset, limit)
			if err != nil {
				WriteLogFile(err)
				// If the structure of the body is wrong, return an HTTP error
				w.WriteHeader(http.StatusUnprocessableEntity)
				return
			}
			defer selDB.Close()
			for selDB.Next() {
				err := selDB.Scan(&SNo, &ActionItem, &MeetingDate, &TargetDate, &Status, &ClosedDate, &Comment, &Flag, &ProjectName, &Owner)
				if err != nil {
					WriteLogFile(err)
					// If the structure of the body is wrong, return an HTTP error
					w.WriteHeader(http.StatusUnprocessableEntity)
					return
				}
				data = append(data, model.ActionItemClosed{ID: SNo, ProjectName: ProjectName, ActionItem: ActionItem, Owner: Owner, MeetingDate: MeetingDate, TargetDate: TargetDate, Status: Status, ClosedDate: ClosedDate, Comment: Comment})
			}
		} else if status == "open" {
			var SNo int
			var ProjectName, ActionItem, Owner, MeetingDate, TargetDate, Status, Comment, Flag string
			var ClosedDate sql.NullString
			count, err := db.Query("SELECT count(*) FROM action_items a JOIN sub_project_manager spm ON spm.id = a.sub_project_manager_id JOIN sub_project s ON spm.sub_project_id = s.id JOIN project p ON p.id = s.project_id JOIN program_manager pg ON pg.id = p.program_manager_id JOIN project_manager pj ON pj.id = spm.manager_id WHERE (a.status != 3 AND pg.program_manager_email= ? AND a.is_active='1') AND (a.action_item like '"+key+"' OR a.meeting_date like '"+key+"' OR a.target_date like '"+key+"' OR a.status like '"+key+"' OR a.closed_date like '"+key+"' OR a.comment like '"+key+"' OR s.sub_project_name like '"+key+"' OR pj.project_manager_name like '"+key+"')", UserName)
			if err != nil {
				WriteLogFile(err)
				// If the structure of the body is wrong, return an HTTP error
				w.WriteHeader(http.StatusUnprocessableEntity)
				return
			}
			defer count.Close()
			for count.Next() {
				err := count.Scan(&countt)
				if err != nil {
					WriteLogFile(err)
					// If the structure of the body is wrong, return an HTTP error
					w.WriteHeader(http.StatusUnprocessableEntity)
					return
				}
			}
			selDB, err := db.Query("call getAllActionsProgramQuery(?,?,?,?)", key, UserName, offset, limit)
			if err != nil {
				WriteLogFile(err)
				// If the structure of the body is wrong, return an HTTP error
				w.WriteHeader(http.StatusUnprocessableEntity)
				return
			}
			defer selDB.Close()
			for selDB.Next() {
				err := selDB.Scan(&SNo, &ActionItem, &MeetingDate, &TargetDate, &Status, &ClosedDate, &Comment, &Flag, &ProjectName, &Owner)
				if err != nil {
					WriteLogFile(err)
					// If the structure of the body is wrong, return an HTTP error
					w.WriteHeader(http.StatusUnprocessableEntity)
					return
				}
				data = append(data, model.ActionItemClosed{ID: SNo, ProjectName: ProjectName, ActionItem: ActionItem, Owner: Owner, MeetingDate: MeetingDate, TargetDate: TargetDate, Status: Status, ClosedDate: "NA", Comment: Comment})
			}
		} else if status == "" {
			var SNo int
			var ProjectName, ActionItem, MeetingDate, TargetDate, Owner, Status, Comment, Flag string
			var ClosedDate sql.NullString
			count, err := db.Query("SELECT count(*) FROM action_items a JOIN sub_project_manager spm ON spm.id = a.sub_project_manager_id JOIN sub_project s ON spm.sub_project_id = s.id JOIN project p ON p.id = s.project_id JOIN program_manager pg ON pg.id = p.program_manager_id JOIN project_manager pj ON pj.id = spm.manager_id WHERE (pg.program_manager_email= ? AND a.is_active='1') AND (a.action_item like '"+key+"' OR a.meeting_date like '"+key+"' OR a.target_date like '"+key+"' OR a.status like '"+key+"' OR a.closed_date like '"+key+"' OR a.comment like '"+key+"' OR s.sub_project_name like '"+key+"' OR pj.project_manager_name like '"+key+"')", UserName)
			if err != nil {
				WriteLogFile(err)
				// If the structure of the body is wrong, return an HTTP error
				w.WriteHeader(http.StatusUnprocessableEntity)
				return
			}
			defer count.Close()
			for count.Next() {
				err := count.Scan(&countt)
				if err != nil {
					WriteLogFile(err)
					// If the structure of the body is wrong, return an HTTP error
					w.WriteHeader(http.StatusUnprocessableEntity)
					return
				}
			}
			selDB, err := db.Query("call getActionsProgramQuery(?,?,?,?)", key, UserName, offset, limit)
			if err != nil {
				WriteLogFile(err)
				// If the structure of the body is wrong, return an HTTP error
				w.WriteHeader(http.StatusUnprocessableEntity)
				return
			}
			defer selDB.Close()
			for selDB.Next() {
				err := selDB.Scan(&SNo, &ActionItem, &MeetingDate, &TargetDate, &Status, &ClosedDate, &Comment, &Flag, &ProjectName, &Owner)
				if err != nil {
					WriteLogFile(err)
					// If the structure of the body is wrong, return an HTTP error
					w.WriteHeader(http.StatusUnprocessableEntity)
					return
				}
				if ClosedDate.Valid == false {
					data = append(data, model.ActionItemClosed{ID: SNo, ProjectName: ProjectName, ActionItem: ActionItem, Owner: Owner, MeetingDate: MeetingDate, TargetDate: TargetDate, Status: Status, ClosedDate: "NA", Comment: Comment})
				} else {
					data = append(data, model.ActionItemClosed{ID: SNo, ProjectName: ProjectName, ActionItem: ActionItem, Owner: Owner, MeetingDate: MeetingDate, TargetDate: TargetDate, Status: Status, ClosedDate: ClosedDate.String, Comment: Comment})
				}
			}

		} else {
			error.Code = "405"
			error.Message = "Method not allowed"
			json.NewEncoder(w).Encode(error)
		}
	}
	Page.TotalData = countt
	Page.Limit = limit
	x := countt
	page = x / limit
	x = x % limit
	if x == 0 {
		Page.TotalPages = page
	} else {
		Page.TotalPages = page + 1
	}
	Page.Data = data
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Page)
}
