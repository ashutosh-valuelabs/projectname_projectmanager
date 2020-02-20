package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	database "projectname_projectmanager/driver"
	models "projectname_projectmanager/model"

	"github.com/gorilla/mux"
)

// UpdateProspectGetData : Get all the Prospects of a particular Program Manager or Project Manager
func (C *Commander) UpdateProspectGetData(writer http.ResponseWriter, response *http.Request) {
	db := database.DbConn()
	defer db.Close()
	writer.Header().Set("Content-Type", "application/json")
	if (strings.Contains(Role, "project manager")) == true {
		offsets, ok := response.URL.Query()["pages"]
		if !ok || len(offsets[0]) < 1 {
			respondwithJSON(writer, http.StatusBadRequest, map[string]string{"Message": "Bad Request ; URL Parameter 'pages' is missing"})
			return
		}
		var prospects []models.Prospects
		var pagination models.Pagination
		pages := offsets[0]
		i, _ := strconv.Atoi(pages)
		offset := i * 10
		limit, ok := response.URL.Query()["limit"]
		if !ok || len(limit[0]) < 1 {
			respondwithJSON(writer, http.StatusBadRequest, map[string]string{"Message": "Bad Request ; URL Parameter 'limit' is missing"})
			return
		}
		limits := limit[0]
		lim, _ := strconv.Atoi(limits)
		pagination.Limit = lim
		result, error := db.Query("call getallupdateprospect_Project(?,?,?)", UserName, offset, lim)
		catch(error)
		defer result.Close()

		count, error := db.Query("call getallupdateprospect_Projectcount(?)", UserName)
		catch(error)
		defer count.Close()

		for result.Next() {
			var prospect models.Prospects
			error := result.Scan(&prospect.ID, &prospect.Project, &prospect.Prospect, &prospect.Status, &prospect.Comments, &prospect.Challenges)
			catch(error)
			prospects = append(prospects, prospect)
		}
		var pagecount int
		if count.Next() != false {
			count.Scan(&pagecount)
		} else {
			pagecount = 0
		}
		pagination.TotalData = pagecount
		pagination.Data = prospects
		x1 := pagecount / 10
		x := pagecount % 10
		x2 := x1 + 1

		if x == 0 {
			pagination.TotalPages = x1
		} else {
			pagination.TotalPages = x2
		}
		x, error = strconv.Atoi(pages)
		catch(error)
		if pagination.TotalPages != 0 {
			x1 = (x + 1)
		}
		pagination.Page = x1
		writer.WriteHeader(http.StatusOK)
		json.NewEncoder(writer).Encode(pagination)

	} else if (strings.Contains(Role, "program manager")) == true {
		offsets, ok := response.URL.Query()["pages"]
		if !ok || len(offsets[0]) < 1 {
			respondwithJSON(writer, http.StatusBadRequest, map[string]string{"Message": "Bad Request ; URL Parameter 'pages' is missing"})
			return
		}
		pages := offsets[0]
		i, _ := strconv.Atoi(pages)
		offset := i * 10

		limit, ok := response.URL.Query()["limit"]
		if !ok || len(limit[0]) < 1 {
			respondwithJSON(writer, http.StatusBadRequest, map[string]string{"Message": "Bad Request ; URL Parameter 'limit' is missing"})
			return
		}
		limits := limit[0]
		lim, _ := strconv.Atoi(limits)
		var prospects []models.Prospects
		var pagination models.Pagination
		result, error := db.Query("call getallupdateprospect_Program(?,?,?)", UserName, offset, lim)
		catch(error)
		defer result.Close()
		pagination.Limit = lim
		count, error := db.Query("call getallupdateprospect_Programcount(?)", UserName)
		catch(error)
		defer count.Close()

		for result.Next() {
			var prospect models.Prospects
			error := result.Scan(&prospect.ID, &prospect.Project, &prospect.Prospect, &prospect.Status, &prospect.Comments, &prospect.Challenges)
			catch(error)
			prospects = append(prospects, prospect)
		}
		var pagecount int
		if count.Next() != false {
			count.Scan(&pagecount)
		} else {
			pagecount = 0
		}
		pagination.TotalData = pagecount
		pagination.Data = prospects
		x1 := pagecount / 10
		x := pagecount % 10
		x2 := x1 + 1

		if x == 0 {
			pagination.TotalPages = x1
		} else {
			pagination.TotalPages = x2
		}
		x, _ = strconv.Atoi(pages)
		if pagination.TotalPages != 0 {
			x1 = (x + 1)
		}
		pagination.Page = x1
		writer.WriteHeader(http.StatusOK)

		json.NewEncoder(writer).Encode(pagination)
	} else {
		fmt.Println("Not Found")
		writer.WriteHeader(http.StatusNotFound)
	}
}

// UpdateProspectCreateData : to insert the data
func (C *Commander) UpdateProspectCreateData(writer http.ResponseWriter, response *http.Request) {
	db := database.DbConn()
	defer db.Close()
	writer.Header().Set("Content-Type", "application/json")
	if (strings.Contains(Role, "project manager")) == true {
		var prospect models.Prospects
		var managerID int
		json.NewDecoder(response.Body).Decode(&prospect)
		mid, error := db.Query("call updateprospects_getmanagerID(?,?)", prospect.Project, UserName)
		catch(error)
		defer mid.Close()
		if mid.Next() != false {
			error := mid.Scan(&managerID)
			catch(error)
		}
		if managerID != 0 {
			var id int
			var flag bool
			upid, error := db.Query("SELECT id FROM update_prospects WHERE sub_project_manager_id = ? AND is_active = 1", managerID)
			catch(error)
			defer upid.Close()
			if upid.Next() != false {
				error := upid.Scan(&id)
				catch(error)
			}
			if prospect.Status == "active" || prospect.Status == "onhold" || prospect.Status == "inactive" || prospect.Status == "lost" || prospect.Status == "won" {
				pros, error := db.Query("select exists (select prospect from update_prospects where sub_project_manager_id = ? and update_prospects.is_active = 1 and update_prospects.prospect = ?)", managerID, prospect.Prospect)
				catch(error)
				defer pros.Close()
				if pros.Next() != false {
					error := pros.Scan(&flag)
					catch(error)
				}
				if flag == false {
					insert, error := db.Prepare("INSERT INTO update_prospects(sub_project_manager_id,prospect,status,comments,challenges,created_at,updated_at) VALUES(?,?,?,?,?,now(),now())")
					catch(error)
					_, error = insert.Exec(managerID, prospect.Prospect, prospect.Status, prospect.Comments, prospect.Challenges)
					catch(error)
					writer.WriteHeader(http.StatusCreated)
					respondwithJSON(writer, http.StatusOK, map[string]string{"message": "Inserted Successfully"})
				} else {
					respondwithJSON(writer, http.StatusConflict, map[string]string{"message": "Duplicate entry for prospect"})
				}
			} else {
				respondwithJSON(writer, http.StatusExpectationFailed, map[string]string{"message": "Wrong Status"})
			}
		} else {
			respondwithJSON(writer, http.StatusBadRequest, map[string]string{"message": "Project not under you"})
		}
	} else {
		respondwithJSON(writer, http.StatusUnauthorized, map[string]string{"message": "Not Authorized"})
	}
}

// UpdateProspectGetDataID : to get the data according to search
func (C *Commander) UpdateProspectGetDataID(writer http.ResponseWriter, response *http.Request) { // Get all the Prospects of a particular Program Manager or Project Manager
	db := database.DbConn()
	defer db.Close()
	writer.Header().Set("Content-Type", "application/json")

	if (strings.Contains(Role, "program manager")) == true {
		offsets, ok := response.URL.Query()["pages"]
		if !ok || len(offsets[0]) < 1 {
			respondwithJSON(writer, http.StatusBadRequest, map[string]string{"Message": "Bad Request ; URL Parameter 'pages' is missing"})
			return
		}
		pages := offsets[0]
		i, _ := strconv.Atoi(pages)
		offset := i * 10
		limit, ok := response.URL.Query()["limit"]
		if !ok || len(limit[0]) < 1 {
			respondwithJSON(writer, http.StatusBadRequest, map[string]string{"Message": "Bad Request ; URL Parameter 'limit' is missing"})
			return
		}
		limits := limit[0]
		lim, _ := strconv.Atoi(limits)
		var prospects []models.Prospects
		var pagination models.Pagination
		parameter := mux.Vars(response)
		key := parameter["id"]
		var search string = key + "%"
		pagination.Limit = lim

		result, error := db.Query("call getallupdateprospect_ProgramID(?,?,?,?)", UserName, search, offset, lim)
		catch(error)
		defer result.Close()
		for result.Next() {
			var prospect models.Prospects
			error := result.Scan(&prospect.ID, &prospect.Project, &prospect.Prospect, &prospect.Status, &prospect.Comments, &prospect.Challenges)
			catch(error)
			prospects = append(prospects, prospect)
		}
		count, error := db.Query("call getallupdateprospect_ProgramIDcount(?,?)", UserName, search)
		catch(error)
		defer count.Close()
		var pagecount int
		if count.Next() != false {
			count.Scan(&pagecount)
		} else {
			pagecount = 0
		}
		pagination.TotalData = pagecount
		pagination.Data = prospects
		x1 := pagecount / 10
		x := pagecount % 10
		x2 := x1 + 1

		if x == 0 {
			pagination.TotalPages = x1
		} else {
			pagination.TotalPages = x2
		}
		x, _ = strconv.Atoi(pages)
		if pagination.TotalPages != 0 {
			x1 = (x + 1)
		}
		pagination.Page = x1
		writer.WriteHeader(http.StatusOK)
		json.NewEncoder(writer).Encode(pagination)

	} else if (strings.Contains(Role, "project manager")) == true {

		offsets, ok := response.URL.Query()["pages"]
		if !ok || len(offsets[0]) < 1 {
			respondwithJSON(writer, http.StatusBadRequest, map[string]string{"Message": "Bad Request ; URL Parameter 'pages' is missing"})
			return
		}
		pages := offsets[0]
		i, _ := strconv.Atoi(pages)
		offset := i * 10

		limit, ok := response.URL.Query()["limit"]
		if !ok || len(limit[0]) < 1 {
			respondwithJSON(writer, http.StatusBadRequest, map[string]string{"Message": "Bad Request ; URL Parameter 'limit' is missing"})
			return
		}
		limits := limit[0]
		lim, _ := strconv.Atoi(limits)

		var prospects []models.Prospects
		var pagination models.Pagination
		parameter := mux.Vars(response)
		key := parameter["id"]
		var search string = key + "%"
		pagination.Limit = lim
		result, error := db.Query("call getallupdateprospects_ProjectID(?,?,?,?)", UserName, search, offset, lim)
		catch(error)
		defer result.Close()
		for result.Next() {
			var prospect models.Prospects
			error := result.Scan(&prospect.ID, &prospect.Project, &prospect.Prospect, &prospect.Status, &prospect.Comments, &prospect.Challenges)
			catch(error)
			prospects = append(prospects, prospect)
		}
		count, error := db.Query("call getallupdateprospect_ProjectcountID(?,?)", UserName, search)
		catch(error)
		defer count.Close()
		var pagecount int

		if count.Next() != false {
			count.Scan(&pagecount)
		} else {
			pagecount = 0
		}
		pagination.TotalData = pagecount
		pagination.Data = prospects
		x1 := pagecount / 10
		x := pagecount % 10
		x2 := x1 + 1

		if x == 0 {
			pagination.TotalPages = x1
		} else {
			pagination.TotalPages = x2
		}
		x, _ = strconv.Atoi(pages)
		if pagination.TotalPages != 0 {
			x1 = (x + 1)
		}
		pagination.Page = x1
		writer.WriteHeader(http.StatusOK)
		json.NewEncoder(writer).Encode(pagination)
	} else {
		writer.WriteHeader(http.StatusNotFound)
	}
}

// UpdateProspectUpdateData : to update the data
func (C *Commander) UpdateProspectUpdateData(writer http.ResponseWriter, response *http.Request) {
	db := database.DbConn()
	defer db.Close()
	writer.Header().Set("Content-Type", "application/json")
	if (strings.Contains(Role, "project manager")) == true {
		var prospect models.Prospects
		var managerID int
		json.NewDecoder(response.Body).Decode(&prospect)
		mid, error := db.Query("call updateprospects_getmanagerID(?,?)", prospect.Project, UserName)
		catch(error)
		defer mid.Close()
		if mid.Next() != false {
			error := mid.Scan(&managerID)
			catch(error)
		}
		fmt.Println(managerID)
		if managerID != 0 {
			var id int
			upid, error := db.Query("SELECT id FROM update_prospects WHERE sub_project_manager_id = ? AND is_active = 1", managerID)
			catch(error)
			defer upid.Close()
			if upid.Next() != false {
				error := upid.Scan(&id)
				catch(error)
			}

			if id != 0 {
				query, error := db.Prepare("UPDATE update_prospects SET prospect = ?, status = ?, comments = ?, challenges = ?, updated_at = ? WHERE sub_project_manager_id = ? AND update_prospects.id = ? AND is_active = 1")
				catch(error)

				_, error = query.Exec(prospect.Prospect, prospect.Status, prospect.Comments, prospect.Challenges, time.Now().Format("2006-01-02 15:04:05"), managerID, prospect.ID)
				catch(error)
				defer query.Close()
				respondwithJSON(writer, http.StatusOK, map[string]string{"message": "Updated Successfully"})
			} else {
				respondwithJSON(writer, http.StatusBadRequest, map[string]string{"Message": "Data Not Found"})
			}
		} else {
			respondwithJSON(writer, http.StatusBadRequest, map[string]string{"Message": "Project not under you"})
		}
	} else {
		respondwithJSON(writer, http.StatusUnauthorized, map[string]string{"Message": "Not Authorized"})
	}
}

// UpdateProspectDeleteData : to delete the data
func (C *Commander) UpdateProspectDeleteData(writer http.ResponseWriter, response *http.Request) {
	db := database.DbConn()
	defer db.Close()
	writer.Header().Set("Content-Type", "application/json")
	if (strings.Contains(Role, "project manager")) == true {
		var prospect models.Prospects
		var managerID int
		json.NewDecoder(response.Body).Decode(&prospect)
		result, error := db.Query("SELECT sub_project_manager_id FROM update_prospects WHERE id = ?", prospect.ID)
		catch(error)
		defer result.Close()
		if result.Next() != false {
			error := result.Scan(&managerID)
			catch(error)
		}
		if managerID != 0 {
			var email string
			manageremail, error := db.Query("call deleteupdateprospect(?)", managerID)
			catch(error)
			defer manageremail.Close()
			if manageremail.Next() != false {
				error := manageremail.Scan(&email)
				catch(error)
			}
			if email == UserName {
				stmt, error := db.Prepare("UPDATE update_prospects SET is_active = 0 WHERE sub_project_manager_id = ? and update_prospects.id = ?")
				catch(error)
				_, error = stmt.Exec(managerID, prospect.ID)
				catch(error)
				defer stmt.Close()
				respondwithJSON(writer, http.StatusOK, map[string]string{"message": "Deleted Successfully"})
			} else {
				respondwithJSON(writer, http.StatusUnauthorized, map[string]string{"Message": "Not Authorized to Delete"})
			}
		} else {
			respondwithJSON(writer, http.StatusBadRequest, map[string]string{"Message": "Error from Front End"})
		}
	} else {
		respondwithJSON(writer, http.StatusUnauthorized, map[string]string{"Message": "Not Authorized"})
	}
}
