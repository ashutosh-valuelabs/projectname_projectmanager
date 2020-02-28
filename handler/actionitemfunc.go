package handler

import (
	"net/http"
	database "projectname_projectmanager/driver"
	models "projectname_projectmanager/model"
	"time"

	_ "github.com/go-sql-driver/mysql" //blank import
)

//UpdateData :
func UpdateData(actionitem models.ActionItemClosed, User string) (int, string, error) {
	type data map[string]interface{}
	var managerDetailsID string
	//var managerEmail string
	var closedInTime int
	db := database.DbConn()
	if actionitem.ID != 0 {
		if actionitem.Status != "open" && actionitem.Status != "inprogress" && actionitem.Status != "closed" && actionitem.Status != "onhold" {
			// w.Header().Set("Content-Type", "application/json")
			// w.WriteHeader(http.StatusUnprocessableEntity)
			// //error.Code = "1265"
			// error.Message = "Data truncated for column 'status'. Invalid entry for column 'status'"
			// json.NewEncoder(w).Encode(error)
			return http.StatusBadRequest, "Data truncated for column 'status'. Invalid entry for column 'status'", nil
		}
		myDateString := "2006-01-02"
		myMeetingDate, err := time.Parse(myDateString, actionitem.MeetingDate)
		if err != nil {
			return 0, "", err
			// WriteLogFile(err)
			// panic(err)
		}
		myTargetDate, err := time.Parse(myDateString, actionitem.TargetDate)
		if err != nil {
			return 0, "", err
			// WriteLogFile(err)
			// panic(err)
		}
		isTargetDateBeforeMeetingDate := myTargetDate.Before(myMeetingDate)
		if actionitem.ClosedDate != "" {
			if actionitem.Status == "closed" {
				myClosedDate, err := time.Parse(myDateString, actionitem.ClosedDate)
				if err != nil {
					return 0, "", err
					// WriteLogFile(err)
					// panic(err)
				}
				isClosedDateBeforeMeetingDate := myClosedDate.Before(myMeetingDate)
				if isClosedDateBeforeMeetingDate == false && isTargetDateBeforeMeetingDate == false {
					isClosedDateAfterTargetDate := myClosedDate.After(myTargetDate)
					if isClosedDateAfterTargetDate == false {
						closedInTime = 1
					} else {
						closedInTime = 0
					}
					selDB, err := db.Query("SELECT id FROM sub_project_manager WHERE sub_project_id = (SELECT id FROM sub_project WHERE sub_project_name=? )and manager_id = (SELECT id FROM project_manager WHERE project_manager_email =?)", actionitem.ProjectName, User)
					defer selDB.Close()
					if err != nil {
						return http.StatusUnprocessableEntity, "", err
						//WriteLogFile(err)
						// If the structure of the body is wrong, return an HTTP error
						// w.WriteHeader(http.StatusUnprocessableEntity)
						// return
					}
					if selDB.Next() != false {
						err := selDB.Scan(&managerDetailsID)
						if err != nil {
							return http.StatusUnprocessableEntity, "", err
							// WriteLogFile(err)
							// // If the structure of the body is wrong, return an HTTP error
							// w.WriteHeader(http.StatusUnprocessableEntity)
							// return
						}
					} else {
						return http.StatusForbidden, "", nil
						//w.WriteHeader(http.StatusBadRequest)
					}

					updatedAt := time.Now()
					updForm, err := db.Prepare("UPDATE action_items SET action_item=? , sub_project_manager_id=? , meeting_date=? , target_date=? , status=? , closed_date=? , comment=? , updated_at=? , closed_in_time=? WHERE id=?")
					if err != nil {
						return 0, "", err
						// WriteLogFile(err)
						// panic(err.Error())
					}
					defer updForm.Close()
					updForm.Exec(actionitem.ActionItem,
						managerDetailsID,
						actionitem.MeetingDate,
						actionitem.TargetDate,
						actionitem.Status,
						actionitem.ClosedDate,
						actionitem.Comment,
						updatedAt,
						closedInTime,
						actionitem.ID)
					defer updForm.Close()
					// sii, _ := solr.NewSolrInterface("http://localhost:8983/solr", "action_items")
					// params := &url.Values{}
					// params.Add("commitWithin", "500")
					// solrID := fmt.Sprintf("ID_i:%d", int(actionitem.ID))
					// res, _ := sii.Delete(data{"query": solrID}, params)
					// fmt.Println(res)
					return http.StatusCreated, "", nil
					// w.Header().Set("Content-Type", "application/json")
					// w.WriteHeader(http.StatusCreated)
				}
				// return
				// w.Header().Set("Content-Type", "application/json")
				// w.WriteHeader(http.StatusUnprocessableEntity)
				// //error.Code = "1292"
				if isClosedDateBeforeMeetingDate == true && isTargetDateBeforeMeetingDate == true {
					return http.StatusBadRequest, "Incorrect Closed Date and Target Date Values", nil
				} else if isTargetDateBeforeMeetingDate == true {
					return http.StatusBadRequest, "Incorrect Target Date Value", nil
				} else if isClosedDateBeforeMeetingDate == true {
					return http.StatusBadRequest, "Incorrect Closed Date Value", nil
				}
				//json.NewEncoder(w).Encode(error)

			} else {
				return http.StatusBadRequest, "Closed date cannot be filled if the status is not closed", nil
				// w.Header().Set("Content-Type", "application/json")
				// w.WriteHeader(http.StatusBadRequest)
				// //error.Code = "1293"
				// error.Message = "Closed date cannot be filled if the status is not closed"
				// json.NewEncoder(w).Encode(error)
			}
		} else {
			if actionitem.Status != "closed" {
				if isTargetDateBeforeMeetingDate == false {
					selDB, err := db.Query("SELECT id FROM sub_project_manager WHERE sub_project_id = (SELECT id FROM sub_project WHERE sub_project_name=? )and manager_id = (SELECT id FROM project_manager WHERE project_manager_email =?)", actionitem.ProjectName, User)
					defer selDB.Close()
					if err != nil {
						return http.StatusUnprocessableEntity, "", err
						// WriteLogFile(err)
						// // If the structure of the body is wrong, return an HTTP error
						// w.WriteHeader(http.StatusUnprocessableEntity)
						// return
					}
					if selDB.Next() != false {
						err := selDB.Scan(&managerDetailsID)
						if err != nil {
							return http.StatusUnprocessableEntity, "", err
							// WriteLogFile(err)
							// // If the structure of the body is wrong, return an HTTP error
							// w.WriteHeader(http.StatusUnprocessableEntity)
							// return
						}
					} else {
						return http.StatusBadRequest, "", nil
						//w.WriteHeader(http.StatusBadRequest)
					}
					updatedAt := time.Now()
					updForm, err := db.Prepare("UPDATE action_items SET action_item=? , sub_project_manager_id=? , meeting_date=? , target_date=? , status=?, closed_date=NULL , comment=? , updated_at=? , closed_in_time = NULL WHERE id=?")
					if err != nil {
						return 0, "", err
						// WriteLogFile(err)
						// panic(err.Error())
					}
					updForm.Exec(actionitem.ActionItem,
						managerDetailsID,
						actionitem.MeetingDate,
						actionitem.TargetDate,
						actionitem.Status,
						actionitem.Comment,
						updatedAt,
						actionitem.ID)
					defer updForm.Close()
					// sii, _ := solr.NewSolrInterface("http://localhost:8983/solr", "action_items")
					// params := &url.Values{}
					// params.Add("commitWithin", "500")
					// solrID := fmt.Sprintf("ID_i:%d", int(actionitem.ID))
					// res, _ := sii.Delete(data{"query": solrID}, params)
					// fmt.Println(res)
					return http.StatusCreated, "", nil
					// w.Header().Set("Content-Type", "application/json")
					// w.WriteHeader(http.StatusCreated)
				}
				// w.Header().Set("Content-Type", "application/json")
				// w.WriteHeader(http.StatusUnprocessableEntity)
				//error.Code = "1292"
				if isTargetDateBeforeMeetingDate == true {
					return http.StatusBadRequest, "Incorrect Target Date Value", nil
					//json.NewEncoder(w).Encode(error)
				}

			} else {
				return http.StatusBadRequest, "Closed date missing", nil
				// w.Header().Set("Content-Type", "application/json")
				// w.WriteHeader(http.StatusUnprocessableEntity)
				// error.Message = "Closed date missing"
				// json.NewEncoder(w).Encode(error)

			}
		}
	} else {
		return http.StatusBadRequest, "", nil
	}
	return 0, "", nil
}

//DeleteData :
func DeleteData(x int, User string) (int, error) {
	var managerDetailsID string
	var managerEmail string
	db := database.DbConn()
	selDB, err := db.Query("SELECT sub_project_manager_id from action_items WHERE id =? and is_active='1'", x)
	defer selDB.Close()
	if err != nil {
		return 0, err
	}
	if selDB.Next() != false {
		err := selDB.Scan(&managerDetailsID)
		if err != nil {
			return 0, err
		}
	} else {
		return http.StatusBadRequest, nil
	}
	selDB, err = db.Query("SELECT project_manager.project_manager_email FROM project_manager JOIN sub_project_manager ON project_manager.id = sub_project_manager.manager_id WHERE sub_project_manager.id =?", managerDetailsID)
	defer selDB.Close()
	if err != nil {
		return 0, err
	}
	if selDB.Next() != false {
		err := selDB.Scan(&managerEmail)
		if err != nil {
			return 0, err
		}
	} else {
		return http.StatusBadRequest, nil //w.WriteHeader(http.StatusBadRequest)
	}
	if managerEmail == User {
		updatedAt := time.Now()
		updForm, err := db.Prepare("UPDATE action_items SET is_active = '0', updated_at = ? WHERE id=?")
		if err != nil {
			return 0, err
			//WriteLogFile(err)
			//panic(err.Error())
		}
		defer updForm.Close()
		updForm.Exec(updatedAt, x)
		return http.StatusOK, nil
		//w.WriteHeader(http.StatusOK)
	}
	return http.StatusForbidden, nil
	//w.WriteHeader(http.StatusForbidden)

}
