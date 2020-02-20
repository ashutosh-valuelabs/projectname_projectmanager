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

//Putallactiveresignations : Func to put all activeresignations
func (c *Commander) Putallactiveresignations(w http.ResponseWriter, r *http.Request) {
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
	//cell, err := xlsx.GetCellValue("Sheet1", "A1")
	//catch(err)
	//fmt.Println(cell)

	//// Get all the rows in the Sheet1.
	//	if cell == "Open Positions" {
	columnName, err := db.Query("SELECT column_name FROM information_schema.columns WHERE table_schema = 'weekly_update' AND table_name = 'active_resignations'")
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

	style1, _ := xlsx.NewStyle(`{"number_format":14}`)
	xlsx.SetCellStyle("Sheet1", "H2", "H"+length, style1)
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
				w.WriteHeader(http.StatusBadRequest)
				//fmt.Fprintln(w, "please correct your column ordering")
				break
			}

		} else {
			var inputAttributes []string

			for j < len(row) {
				inputAttributes = append(inputAttributes, row[j])
				j++
			}

			managerNameExcel := inputAttributes[2]
			backfillRequired := inputAttributes[3]
			regreNonRegre := inputAttributes[4]
			dateOfResignation := inputAttributes[6]
			dateOfLeaving := inputAttributes[7]

			//fmt.Println(managerNameExcel)

			if managerNameExcel != managerName {
				var rowError string
				rowError = rowError + " You are not the manager of this project or watch out the spelling of name "
				fmt.Println(rowError)
				iString := strconv.Itoa(i + 1)
				excelErrors = append(excelErrors, "error on line: "+iString+" "+rowError)
			}
			if backfillRequired != "0" && backfillRequired != "1" {
				var rowError string
				rowError = rowError + " Backfill Required can have only 0/1  "
				fmt.Println(rowError)
				iString := strconv.Itoa(i + 1)
				excelErrors = append(excelErrors, "error on line: "+iString+" "+rowError)
			}
			if regreNonRegre != "0" && regreNonRegre != "1" {
				var rowError string
				rowError = rowError + " Regre_non_regre can have only 0/1  "
				fmt.Println(rowError)
				iString := strconv.Itoa(i + 1)
				excelErrors = append(excelErrors, "error on line: "+iString+" "+rowError)
			}
			if dateOfLeaving < dateOfResignation {
				var rowError string
				rowError = rowError + " Date of leaving cannot be smaller than date of resignation  "
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
				employeeName := inputAttributes[0]
				projectName := inputAttributes[1]
				backfillRequired := inputAttributes[3]
				regreNonRegre := inputAttributes[4]
				status := inputAttributes[5]
				dateOfResignation := inputAttributes[6]
				dateOfLeaving := inputAttributes[7]
				createdAt := inputAttributes[8]

				createdAt = createdAt + ":00"

				dateOfLeavingFormat, err := time.Parse("1-2-06", dateOfLeaving)
				catch(err)
				dateOfResignationFormat, err := time.Parse("1-2-06", dateOfResignation)
				catch(err)

				createdAtFormat, err := time.Parse("1/2/06 15:04:05", createdAt)
				catch(err)

				var managerProjectID int
				queryManagerProjectID := db.QueryRow("select id from sub_project_manager where sub_project_id in ( select id from sub_project where sub_project_name = ? ) and manager_id in ( select id from project_manager where project_manager_name = ? )", projectName, managerName)
				queryManagerProjectID.Scan(&managerProjectID)

				if managerProjectID != 0 {
					queryIDIsactive, err := db.Query("SELECT id,is_active FROM active_resignations where sub_project_manager_id = ? and emp_name = ? ", managerProjectID, employeeName)
					catch(err)
					var id, flag int
					for queryIDIsactive.Next() {
						queryIDIsactive.Scan(&id, &flag)
					}

					if id != 0 {
						if flag != 0 {
							queryID := db.QueryRow("SELECT id FROM active_resignations where sub_project_manager_id = ? and emp_name = ? and backfill_required = ? and regre_non_regre = ? and status = ? and date_of_resignation = ?  and date_of_leaving = ?  ", managerProjectID, backfillRequired, regreNonRegre, status, dateOfResignationFormat, dateOfLeavingFormat)
							var updateID int
							queryID.Scan(&updateID)
							if updateID == 0 {
								queryUpdate, err := db.Prepare("update active_resignations set backfill_required = ?, regre_non_regre = ?, status = ?, date_of_resignation = ?, date_of_leaving = ?, updated_at = now() where id = ? ")
								catch(err)
								_, err = queryUpdate.Exec(backfillRequired, regreNonRegre, status, dateOfResignationFormat, dateOfLeavingFormat, id)
							} else {
								w.WriteHeader(http.StatusConflict)
								//fmt.Fprintln(w, "already this data is present")
							}
						} else {
							w.WriteHeader(http.StatusNotFound)
							//fmt.Fprintln(w, "Your project is deleted")
						}
					} else {
						queryInsert, err := db.Prepare("insert into active_resignations (sub_project_manager_id,emp_name,backfill_required,regre_non_regre,status,date_of_resignation,date_of_leaving,created_at,updated_at) values (?,?,?,?,?,?,?,?,?)  ")
						catch(err)
						_, err = queryInsert.Exec(managerProjectID, employeeName, backfillRequired, regreNonRegre, status, dateOfResignationFormat, dateOfLeavingFormat, createdAtFormat, createdAtFormat)
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
