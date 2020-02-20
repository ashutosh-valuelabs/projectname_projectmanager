package handler

import (
	"fmt"
	"time"

	//_"github.com/go-sql-driver/mysql"
	"io"
	"net/http"
	"os"
	database "projectname_projectmanager/driver"
	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize"
)

var managerName string

//Putallopenpositions : Func to put all openpositions
func (c *Commander) Putallopenpositions(w http.ResponseWriter, r *http.Request) {
	// DATABASE CONNECTION
	db := database.DbConn()
	defer db.Close()
	// QUERY TO EXTRACT MANAGER NAME FROM UserName
	user := db.QueryRow("select project_manager_name from project_manager where project_manager_email = ?", UserName)
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
	columnName, err := db.Query("SELECT column_name FROM information_schema.columns WHERE table_schema = 'weekly_update' AND table_name = 'open_positions'")
	var columnNameArray []string
	catch(err)
	for columnName.Next() {
		var columnNameAttributes string
		columnName.Scan(&columnNameAttributes)
		columnNameArray = append(columnNameArray, columnNameAttributes)
	}

	columnNameArray = append(columnNameArray[:0], columnNameArray[2:]...)
	columnNameArray = append(columnNameArray[:11], columnNameArray[12])

	rows, err := xlsx.GetRows("Sheet1")
	catch(err)

	lengthInt := len(rows)
	length := strconv.Itoa(lengthInt)

	style, _ := xlsx.NewStyle(`{"number_format":22}`)
	xlsx.SetCellStyle("Sheet1", "N2", "N"+length, style)

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

			excelColumnName = append(excelColumnName[:0], excelColumnName[2:]...)

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

			managerNameExcel := inputAttributes[1]
			typePosition := inputAttributes[3]
			priority := inputAttributes[4]
			fmt.Println(managerNameExcel)

			if managerNameExcel != managerName {
				var rowError string
				rowError = rowError + " You are not the manager of this project or watch out the spelling of name "
				fmt.Println(rowError)
				iString := strconv.Itoa(i + 1)
				excelErrors = append(excelErrors, "error on line: "+iString+" "+rowError)
			}
			if typePosition != "new" && typePosition != "New" && typePosition != "resignation" && typePosition != "Resignation" && typePosition != "Replacement" && typePosition != "replacement" {
				var rowError string
				rowError = rowError + " type_position must be New/Replacement/resignation "
				fmt.Println(rowError)
				iString := strconv.Itoa(i + 1)
				excelErrors = append(excelErrors, "error on line: "+iString+" "+rowError)
			}
			if priority != "high" && priority != "High" && priority != "low" && priority != "Low" && priority != "Medium" && priority != "medium" {
				var rowError string
				rowError = rowError + " priority must be High/Low/Medium "
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

				projectName := inputAttributes[0]
				designation := inputAttributes[2]
				typePosition := inputAttributes[3]
				priority := inputAttributes[4]
				position := inputAttributes[5]
				additionalComment := inputAttributes[6]
				l1Due := inputAttributes[7]
				l2Due := inputAttributes[8]
				clientDue := inputAttributes[9]
				l1Happened := inputAttributes[10]
				l2Happened := inputAttributes[11]
				clientHappened := inputAttributes[12]
				createdAt := inputAttributes[13]

				createdAt = createdAt + ":00"

				createdAtFormat, err := time.Parse("1/2/06 15:04:05", createdAt)
				catch(err)
				fmt.Println(projectName)
				var managerProjectID int
				queryManagerProjectID := db.QueryRow("select id from sub_project_manager where sub_project_id in ( select id from sub_project where sub_project_name = ? ) and manager_id in ( select id from project_manager where project_manager_name = ? ) ", projectName, managerName)
				queryManagerProjectID.Scan(&managerProjectID)

				if managerProjectID != 0 {
					queryIDIsactive, err := db.Query("SELECT id,is_active FROM open_positions where sub_project_manager_id = ? and designation =? and type_position = ? ", managerProjectID, designation, typePosition)
					catch(err)
					var id, flag int
					for queryIDIsactive.Next() {
						queryIDIsactive.Scan(&id, &flag)
					}

					if id != 0 {
						if flag != 0 {
							queryID := db.QueryRow("SELECT id FROM open_positions where sub_project_manager_id = ? and designation =? and type_position = ? and priority = ? and position = ? and additional_comment = ? and l1_due = ? and l2_due = ? and client_due = ? and l1_happened = ? and l2_happened = ? and client_happened = ? ", managerProjectID, designation, typePosition, priority, position, additionalComment, l1Due, l2Due, clientDue, l1Happened, l2Happened, clientHappened)
							var updateID int
							queryID.Scan(&updateID)
							if updateID == 0 {
								queryUpdate, err := db.Prepare("update open_positions set priority = ?, position = ?, additional_comment = ?, l1_due = ?, l2_due = ?, client_due = ?,l1_happened = ?, l2_happened = ?, client_happened = ?, updated_at = now() where id = ? ")
								catch(err)
								_, err = queryUpdate.Exec(priority, position, additionalComment, l1Due, l2Due, clientDue, l1Happened, l2Happened, clientHappened, id)
							} else {
								w.WriteHeader(http.StatusConflict)
								//fmt.Fprintln(w, "already this data is present")
							}
						} else {
							w.WriteHeader(http.StatusNotFound)
							//fmt.Fprintln(w, "Your project is deleted")
						}
					} else {
						queryInsert, err := db.Prepare("insert into open_positions (sub_project_manager_id,designation,type_position,priority,position,additional_comment,l1_due,l2_due,client_due,l1_happened,l2_happened,client_happened,created_at,updated_at) values (?,?,?,?,?,?,?,?,?,?,?,?,?,?) ")
						catch(err)
						_, err = queryInsert.Exec(managerProjectID, designation, typePosition, priority, position, additionalComment, l1Due, l2Due, clientDue, l1Happened, l2Happened, clientHappened, createdAtFormat, createdAtFormat)
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

//IsEqual : Func to check equality of strings
func IsEqual(a1 []string, a2 []string) bool {

	if len(a1) == len(a2) {
		for i, v := range a1 {
			if v != a2[i] {
				return false
			}
		}
	} else {
		return false
	}
	return true
}
