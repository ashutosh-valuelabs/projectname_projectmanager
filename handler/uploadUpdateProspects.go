package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	database "projectname_projectmanager/driver"
	"strconv"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
)

//Putallupdateprospects : Func to put all updateprospects
func (c *Commander) Putallupdateprospects(w http.ResponseWriter, r *http.Request) {
	// DATABASE CONNECTION
	db := database.DbConn()
	defer db.Close()
	// QUERY TO EXTRACT MANAGER NAME FROM USERNAME
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
	columnName, err := db.Query("SELECT column_name FROM information_schema.columns WHERE table_schema = 'weekly_update' AND table_name = 'update_prospects'")
	var columnNameArray []string
	catch(err)
	for columnName.Next() {
		var columnNameAttributes string
		columnName.Scan(&columnNameAttributes)
		columnNameArray = append(columnNameArray, columnNameAttributes)
	}

	columnNameArray = append(columnNameArray[:0], columnNameArray[2:]...)
	columnNameArray = append(columnNameArray[:4], columnNameArray[4])

	rows, err := xlsx.GetRows("Sheet1")
	catch(err)

	lengthInt := len(rows)
	length := strconv.Itoa(lengthInt)

	style, _ := xlsx.NewStyle(`{"number_format":22}`)
	xlsx.SetCellStyle("Sheet1", "G2", "G"+length, style)

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

			if managerNameExcel != managerName {
				var rowError string
				rowError = rowError + " You are not the manager of this project or watch out the spelling of name "
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
				prospect := inputAttributes[2]
				status := inputAttributes[3]
				comments := inputAttributes[4]
				challenges := inputAttributes[5]

				createdAt := inputAttributes[6]

				createdAt = createdAt + ":00"

				createdAtFormat, err := time.Parse("1/2/06 15:04:05", createdAt)
				catch(err)

				var managerProjectID int
				queryManagerProjectID := db.QueryRow("select id from sub_project_manager where sub_project_id in ( select id from sub_project where sub_project_name = ? ) and manager_id in ( select id from project_manager where project_manager_name = ? )", projectName, managerName)
				queryManagerProjectID.Scan(&managerProjectID)

				if managerProjectID != 0 {
					queryIDIsactive, err := db.Query("SELECT id,is_active FROM update_prospects where sub_project_manager_id = ? ", managerProjectID)
					catch(err)
					var id, flag int
					for queryIDIsactive.Next() {
						queryIDIsactive.Scan(&id, &flag)
					}

					if id != 0 {
						if flag != 0 {
							queryID := db.QueryRow("SELECT id FROM update_prospects where sub_project_manager_id = ? and prospect = ? and status = ? and comments = ? and challenges = ? ", managerProjectID, prospect, status, comments, challenges)
							var updateID int
							queryID.Scan(&updateID)
							if updateID == 0 {
								queryUpdate, err := db.Prepare("update update_prospects set prospect = ?, status = ?, comments = ?, challenges = ?, updated_at = now() where id = ?  ")
								catch(err)
								_, err = queryUpdate.Exec(prospect, status, comments, challenges, id)
							} else {
								w.WriteHeader(http.StatusConflict)
								//fmt.Fprintln(w, "already this data is present")
							}
						} else {
							w.WriteHeader(http.StatusNotFound)
							//fmt.Fprintln(w, "Your project is deleted")
						}
					} else {
						queryInsert, err := db.Prepare("insert into update_prospects (sub_project_manager_id,prospect,status,comments,challenges,created_at,updated_at) values (?,?,?,?,?,?,?) ")
						catch(err)
						_, err = queryInsert.Exec(managerProjectID, prospect, status, comments, challenges, createdAtFormat, createdAtFormat)
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
