package handler

import (
	"fmt"
	"time"

	"io"
	"net/http"
	"os"
	database "projectname_projectmanager/driver"
	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize"
)

//Putallactionitems : Func to put all actionitems
func (c *Commander) Putallactionitems(w http.ResponseWriter, r *http.Request) {
	// DATABASE CONNECTION
	db := database.DbConn()
	defer db.Close()
	// QUERY TO EXTRACT MANAGER NAME FROM UserName
	user := db.QueryRow("select manager_name from manager where manager_email_id = ?", UserName)
	user.Scan(&managerName)
	// SETTING CONTENT TYPE TO FORM DATA
	w.Header().Set("Content-Type", "multipart/form-data")
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("File")
	catch(err)
	defer file.Close()
	f, err := os.OpenFile("/home/shivangivarshney/temp/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	catch(err)
	defer f.Close()
	io.Copy(f, file)
	xlsx, err := excelize.OpenFile("/home/shivangivarshney/temp/" + handler.Filename)
	catch(err)
	// Get value from cell by given worksheet name and axis.
	/*cell, err := xlsx.GetCellValue("Sheet1", "A1")
	catch(err)
	fmt.Println(cell)*/

	//// Get all the rows in the Sheet1.
	//	if cell == "Open Positions" {
	columnName, err := db.Query("SELECT column_name FROM information_schema.columns WHERE table_schema = 'weekly_update' AND table_name = 'action_items'")
	var columnNameArray []string
	catch(err)
	for columnName.Next() {
		var columnNameAttributes string
		columnName.Scan(&columnNameAttributes)
		columnNameArray = append(columnNameArray, columnNameAttributes)
	}

	columnNameArray = append(columnNameArray[:0], columnNameArray[1:]...)
	columnNameArray = append(columnNameArray[:1], columnNameArray[2:]...)
	columnNameArray = append(columnNameArray[:6], columnNameArray[6])
	//fmt.Println(columnNameArray)
	rows, err := xlsx.GetRows("Sheet1")
	catch(err)

	lengthInt := len(rows)
	length := strconv.Itoa(lengthInt)

	style, _ := xlsx.NewStyle(`{"number_format":14}`)
	xlsx.SetCellStyle("Sheet1", "D2", "D"+length, style)
	style1, _ := xlsx.NewStyle(`{"number_format":14}`)
	xlsx.SetCellStyle("Sheet1", "E2", "E"+length, style1)
	style2, _ := xlsx.NewStyle(`{"number_format":14}`)
	xlsx.SetCellStyle("Sheet1", "G2", "G"+length, style2)
	style3, _ := xlsx.NewStyle(`{"number_format":22}`)
	xlsx.SetCellStyle("Sheet1", "I2", "I"+length, style3)

	rows, err = xlsx.GetRows("Sheet1")
	catch(err)

	var excelColumnName []string
	var excelErrors []string

	for i, row := range rows {
		j := 0
		if i == 0 {
			continue
		}
		if i == 1 {
			for j < len(row) {
				excelColumnName = append(excelColumnName, row[j])
				j++
			}

			excelColumnName = append(excelColumnName[:1], excelColumnName[3:]...)

			if IsEqual(excelColumnName, columnNameArray) {
				continue
			} else {
				fmt.Fprintln(w, "please correct your column ordering")
				break
			}

		} else {
			var inputAttributes []string

			for j < len(row) {
				inputAttributes = append(inputAttributes, row[j])
				j++
			}

			managerNameExcel := inputAttributes[2]
			status := inputAttributes[5]
			targetDate := inputAttributes[4]
			meetingDate := inputAttributes[3]
			//fmt.Println(managerNameExcel)

			if managerNameExcel != managerName {
				var rowError string
				rowError = rowError + " You are not the manager of this project or watch out the spelling of name "
				fmt.Println(rowError)
				iString := strconv.Itoa(i + 1)
				excelErrors = append(excelErrors, "error on line: "+iString+" "+rowError)
			}
			if status != "open" && status != "Open" && status != "closed" && status != "Closed" && status != "onhold" && status != "Onhold" && status != "inprogress" && status != "Inprogress" {
				var rowError string
				rowError = rowError + " status must be Open/Inprogress/Onhold/Closed "
				fmt.Println(rowError)
				iString := strconv.Itoa(i + 1)
				excelErrors = append(excelErrors, "error on line: "+iString+" "+rowError)
			}
			if targetDate < meetingDate {
				var rowError string
				rowError = rowError + " target_date should be greater than or equal to meetingdate "
				fmt.Println(rowError)
				iString := strconv.Itoa(i + 1)
				excelErrors = append(excelErrors, "error on line: "+iString+" "+rowError)
			}
		}
	}

	if len(excelErrors) == 0 {

		for i, row := range rows {
			j := 0
			if i == 0 || i == 1 {
				continue
			} else {
				var inputAttributes []string
				for j < len(row) {
					inputAttributes = append(inputAttributes, row[j])
					j++
				}
				actionItem := inputAttributes[0]
				projectName := inputAttributes[1]
				meetingDate := inputAttributes[3]
				targetDate := inputAttributes[4]
				status := inputAttributes[5]
				closedDate := inputAttributes[6]
				comment := inputAttributes[7]
				createdAt := inputAttributes[8]

				createdAt = createdAt + ":00"

				meetingDateFormat, err := time.Parse("1-2-06", meetingDate)
				catch(err)
				targetDateFormat, err := time.Parse("1-2-06", targetDate)
				catch(err)
				closedDateFormat, err := time.Parse("1-2-06", closedDate)
				catch(err)
				createdAtFormat, err := time.Parse("1/2/06 15:04:05", createdAt)
				catch(err)

				var closedInTime int
				if closedDate <= targetDate {
					closedInTime = 1
				} else if closedDate > targetDate {
					closedInTime = 0
				}
				var managerProjectID int
				queryManagerProjectID := db.QueryRow("select id from sub_project_manager where sub_project_id in ( select id from sub_project where sub_project_name = ? ) and manager_id in ( select id from project_manager where project_manager_name = ? )", projectName, managerName)
				queryManagerProjectID.Scan(&managerProjectID)

				if managerProjectID != 0 {
					queryIDIsactive, err := db.Query("SELECT id,is_active FROM action_items where sub_project_manager_id = ? and action_item = ? ", managerProjectID, actionItem)
					catch(err)
					var id, flag int
					for queryIDIsactive.Next() {
						queryIDIsactive.Scan(&id, &flag)
					}

					if id != 0 {
						if flag != 0 {
							queryID := db.QueryRow("SELECT id FROM action_items where sub_project_manager_id = ? and meeting_date = ? and target_date = ? and status = ? and closed_date = ? and comment = ? ", managerProjectID, meetingDateFormat, targetDateFormat, status, closedDateFormat, comment)
							var updateID int
							queryID.Scan(&updateID)
							if updateID == 0 {
								queryUpdate, err := db.Prepare("update action_items set meeting_date = ?, target_date = ?, status = ?, closed_date = ?, comment = ?, closed_in_time = ?, updated_at = now() where id = ? ")
								catch(err)
								_, err = queryUpdate.Exec(meetingDateFormat, targetDateFormat, status, closedDateFormat, comment, closedInTime, id)
							} else {
								w.WriteHeader(http.StatusConflict)
								//fmt.Fprintln(w, "already this data is present")
							}
						} else {
							w.WriteHeader(http.StatusNotFound)
							//fmt.Fprintln(w, "Your project is deleted")
						}
					} else {
						queryInsert, err := db.Prepare("insert into action_items (sub_project_manager_id,action_item,meeting_date,target_date,status,closed_date,comment,created_at,closed_in_time,updated_at) values (?,?,?,?,?,?,?,?,?,?) ")
						catch(err)
						_, err = queryInsert.Exec(managerProjectID, actionItem, meetingDateFormat, targetDateFormat, status, closedDateFormat, comment, createdAtFormat, closedInTime, createdAtFormat)
						catch(err)
					}
				} else {
					w.WriteHeader(http.StatusForbidden)
					//fmt.Fprintln(w, "project not under this manager")
				}
			}
		}
	} else {
		fmt.Fprintln(w, excelErrors)
	}
}
