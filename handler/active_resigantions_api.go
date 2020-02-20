package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	database "projectname_projectmanager/driver"
	model "projectname_projectmanager/model"
	models "projectname_projectmanager/model"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

// Get all resigantion under project and program manager.
func (c *Commander) GetAllresignations(w http.ResponseWriter, r *http.Request) {
	db := database.DbConn()
	defer db.Close()
	w.Header().Set("Content-Type", "application/json")
	SetupResponse(&w, r)
	if (strings.Contains(Role, "program manager")) == true {
		var posts []models.Resignations
		var Pag models.Pagination
		offsets, ok := r.URL.Query()["pages"]
		if !ok || offsets[0] == "" {
			w.WriteHeader(http.StatusNotFound)
			respondwithJSON(w, http.StatusOK, map[string]string{"message": "url parameter pages is missing"})
			return
		}
		pages := offsets[0]
		i, _ := strconv.Atoi(pages)

		limit, ok := r.URL.Query()["limit"]
		if !ok || limit[0] == "" {
			w.WriteHeader(http.StatusNotFound)
			respondwithJSON(w, http.StatusOK, map[string]string{"message": "url limit parameter is missing"})
			return
		}
		limits := limit[0]
		Limit, _ := strconv.Atoi(limits)
		if Limit == 0 {
			w.WriteHeader(http.StatusBadGateway)
			respondwithJSON(w, http.StatusOK, map[string]string{"message": "limit can't be 0"})
			return
		}
		offset := i * Limit
		result, err := db.Query("call active_resign_get_all_program(?, ?, ?)", UserName, offset, Limit)
		catch(err)
		defer result.Close()
		Pag.Limit = Limit
		count, err1 := db.Query("SELECT count(id) from active_resignations WHERE is_active = 1 AND sub_project_manager_id in (SELECT id from sub_project_manager where sub_project_id in (select id from sub_project where project_id in (select id from project where program_manager_id in (select id from program_manager where program_manager_email= ?))))", UserName)
		if err1 != nil {
			panic(err1.Error())
		}
		defer count.Close()
		for count.Next() {
			err2 := count.Scan(&Pag.TotalData)
			if err2 != nil {
				panic(err2.Error())
			}
		}
		Pag.TotalPages = (Pag.TotalData / Pag.Limit) + 1
		for result.Next() {
			var post models.Resignations
			err := result.Scan(&post.ID, &post.Empname, &post.Project, &post.Manager, &post.Backfillrequired, &post.Regrenonregre, &post.Status, &post.Dateofresignation, &post.Dateofleaving)
			if err != nil {
				panic(err.Error())
			}
			posts = append(posts, post)
		}
		Pag.Data = posts
		x1 := Pag.TotalData / Pag.Limit
		x := Pag.TotalData % Pag.Limit
		x2 := x1 + 1

		if x == 0 {
			Pag.TotalPages = x1
		} else {
			Pag.TotalPages = x2
		}
		x, _ = strconv.Atoi(pages)
		if Pag.TotalPages != 0 {
			x1 = x + 1
		}
		Pag.Page = x1
		json.NewEncoder(w).Encode(Pag)
	} else if (strings.Contains(Role, "project manager")) == true {
		var posts []models.Resignations
		var Pag models.Pagination
		offsets, ok := r.URL.Query()["pages"]
		if !ok || offsets[0] == "" {
			w.WriteHeader(http.StatusNotFound)
			respondwithJSON(w, http.StatusOK, map[string]string{"message": "url pages parameter is missing"})
			return
		}
		pages := offsets[0]
		i, _ := strconv.Atoi(pages)
		limit, ok := r.URL.Query()["limit"]
		if !ok || limit[0] == "" {
			w.WriteHeader(http.StatusNotFound)
			respondwithJSON(w, http.StatusOK, map[string]string{"message": "url parameter is missing"})
			return
		}
		limits := limit[0]
		Limit, _ := strconv.Atoi(limits)
		if Limit == 0 {
			w.WriteHeader(http.StatusBadGateway)
			respondwithJSON(w, http.StatusOK, map[string]string{"message": "limit can't be 0"})
			return
		}
		offset := i * Limit
		result, err := db.Query("call active_resign_get_all_project(?, ?, ?)", UserName, offset, Limit)
		catch(err)
		defer result.Close()
		Pag.Limit = Limit
		count, err1 := db.Query("SELECT count(id) from active_resignations WHERE (is_active = 1) AND (sub_project_manager_id in (SELECT id from sub_project_manager where sub_project_id in (select id from sub_project) AND (manager_id in (select id from project_manager where project_manager_email =?))))", UserName)
		if err1 != nil {
			panic(err1.Error())
		}
		defer count.Close()
		for count.Next() {
			err2 := count.Scan(&Pag.TotalData)
			if err2 != nil {
				panic(err2.Error())
			}
		}
		Pag.TotalPages = (Pag.TotalData / Pag.Limit) + 1
		for result.Next() {
			var post models.Resignations
			err := result.Scan(&post.ID, &post.Empname, &post.Project, &post.Backfillrequired, &post.Regrenonregre, &post.Status, &post.Dateofresignation, &post.Dateofleaving)
			if err != nil {
				panic(err.Error())
			}
			posts = append(posts, post)
		}
		defer db.Close()
		Pag.Data = posts
		x1 := Pag.TotalData / Pag.Limit
		x := Pag.TotalData % Pag.Limit
		x2 := x1 + 1

		if x == 0 {
			Pag.TotalPages = x1
		} else {
			Pag.TotalPages = x2
		}
		x, _ = strconv.Atoi(pages)
		if Pag.TotalPages != 0 {
			x1 = x + 1
		}
		Pag.Page = x1
		json.NewEncoder(w).Encode(Pag)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

// To cretate Resigantion
func (c *Commander) CreateResignations(w http.ResponseWriter, r *http.Request) {
	db := database.DbConn()
	defer db.Close()
	w.Header().Set("Content-Type", "application/json")
	var ManagerID int
	SetupResponse(&w, r)
	if (strings.Contains(Role, "project manager")) == true {
		var error model.Error
		var post models.Resignations
		err := json.NewDecoder(r.Body).Decode(&post)
		BadRequest(w, err)
		if post.Status != "inprogress" && post.Status != "retained" && post.Status != "exit" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			//error.Code = "1265"
			error.Message = "Data truncated for column 'status'. Invalid entry for column 'status'"
			json.NewEncoder(w).Encode(error)
			return
		}
		stmt1, err := db.Query("select id from sub_project_manager where sub_project_id in (select id from sub_project where sub_project_name= ?) and manager_id in (select id from project_manager where project_manager_email= ?)", post.Project, UserName)
		catch(err)
		defer stmt1.Close()
		if stmt1.Next() != false {
			err := stmt1.Scan(&ManagerID)
			catch(err)
		}
		if ManagerID != 0 {
			stmt, err := db.Prepare("INSERT INTO active_resignations(emp_name, sub_project_manager_id, backfill_required, regre_non_regre, status, date_of_resignation, date_of_leaving, created_at, updated_at) VALUES(?, ?, ?, ?, ?, ?, ?, now(), now())")
			catch(err)
			_, err = stmt.Exec(post.Empname, ManagerID, post.Backfillrequired, post.Regrenonregre, post.Status, post.Dateofresignation, post.Dateofleaving)
			catch(err)
			defer stmt.Close()
			w.WriteHeader(http.StatusCreated)
			respondwithJSON(w, http.StatusOK, map[string]string{"message": "New post was created"})
		} else {
			w.WriteHeader(http.StatusBadRequest)
			respondwithJSON(w, http.StatusOK, map[string]string{"message": "project not under you"})
		}
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		respondwithJSON(w, http.StatusOK, map[string]string{"message": "unauthorised access"})
	}
}

// Get data by Search result
func (c *Commander) GetResignationsbyName(w http.ResponseWriter, r *http.Request) {
	db := database.DbConn()
	SetupResponse(&w, r)
	if (strings.Contains(Role, "project manager")) == true {
		var post models.Resignations
		params := mux.Vars(r)
		key := params["emp_name"]
		var searchKey string = key + "%"
		var Offset int
		var co int //for number of records
		Pages := r.FormValue("pages")
		if Pages == "" {
			w.WriteHeader(http.StatusNotFound)
			respondwithJSON(w, http.StatusOK, map[string]string{"message": "pages url parameter is missing"})
			return
		}
		i1, _ := strconv.Atoi(Pages)
		Offset = 10 * i1
		limit, ok := r.URL.Query()["limit"]
		if !ok || limit[0] == "" {
			w.WriteHeader(http.StatusNotFound)
			respondwithJSON(w, http.StatusOK, map[string]string{"message": "Limit url parameter is missing"})
			return
		}
		limits := limit[0]
		Limit, _ := strconv.Atoi(limits)
		if Limit == 0 {
			w.WriteHeader(http.StatusBadGateway)
			respondwithJSON(w, http.StatusOK, map[string]string{"message": "limit can't be 0"})
			return
		}
		count, err := db.Query("call active_resign_get_by_name_project_count(?, ?)", UserName, searchKey)
		catch(err)
		defer count.Close()
		result, err := db.Query("call active_resign_get_by_name_project_result(?, ?, ?, ?)", UserName, searchKey, Offset, Limit)
		catch(err)
		defer result.Close()
		var posts []models.Resignations
		for result.Next() {
			err := result.Scan(&post.ID, &post.Empname, &post.Project, &post.Backfillrequired, &post.Regrenonregre, &post.Status, &post.Dateofresignation, &post.Dateofleaving)
			catch(err)
			posts = append(posts, post)
		}
		if count.Next() != false {
			count.Scan(&co)
		} else {
			co = 0
		}
		defer db.Close()
		//Pagination
		var Pag models.Pagination
		Pag.TotalData = co
		fmt.Println(co)
		Pag.Limit = Limit
		Pag.Data = posts
		x1 := co / Pag.Limit
		x := co % Pag.Limit
		x2 := x1 + 1

		if x == 0 {
			Pag.TotalPages = x1
		} else {
			Pag.TotalPages = x2
		}
		x, _ = strconv.Atoi(Pages)
		if Pag.TotalPages != 0 {
			x1 = x + 1
		}
		Pag.Page = x1
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Pag)
	} else if (strings.Contains(Role, "program manager")) == true {
		var post models.Resignations
		params := mux.Vars(r)
		key := params["emp_name"]
		var searchKey string = "'" + key + "%'"
		var Offset int
		Pages := r.FormValue("pages")
		if Pages == "" {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "Url parameter pages is missing")
			return
		}
		i1, _ := strconv.Atoi(Pages)
		limit, ok := r.URL.Query()["limit"]
		if !ok || limit[0] == "" {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "Url parameter limit is missing")
			return
		}
		limits := limit[0]
		Limit, _ := strconv.Atoi(limits)
		if Limit == 0 {
			w.WriteHeader(http.StatusBadGateway)
			fmt.Fprintf(w, "Limit can't be 0")
			return
		}
		Offset = Limit * i1
		count, err := db.Query("SELECT count(active_resignations.id) from active_resignations left join sub_project_manager on active_resignations.sub_project_manager_id = sub_project_manager.id left join sub_project on  sub_project_manager.sub_project_id = sub_project.id left join project on sub_project.project_id = project.id left join program_manager on project.program_manager_id=program_manager.id left join project_manager on sub_project_manager.manager_id=project_manager.id WHERE program_manager.program_manager_email = ? AND active_resignations.is_active = 1 AND ((emp_name LIKE "+searchKey+") OR (sub_project.sub_project_name LIKE "+searchKey+") OR (project_manager.project_manager_name LIKE "+searchKey+") or (date_of_resignation LIKE "+searchKey+") OR (date_of_leaving LIKE "+searchKey+") OR (status LIKE "+searchKey+"))", UserName)
		catch(err)
		defer count.Close()
		result, err := db.Query("select active_resignations.id, emp_name, sub_project.sub_project_name, project_manager.project_manager_name, backfill_required, regre_non_regre, status, date_of_resignation, date_of_leaving from active_resignations left join sub_project_manager on active_resignations.sub_project_manager_id = sub_project_manager.id left join sub_project on  sub_project_manager.sub_project_id = sub_project.id left join project on sub_project.project_id = project.id left join program_manager on project.program_manager_id=program_manager.id left join project_manager on sub_project_manager.manager_id=project_manager.id WHERE program_manager.program_manager_email = ? AND active_resignations.is_active = 1 AND ((emp_name LIKE "+searchKey+") OR (sub_project.sub_project_name LIKE "+searchKey+") OR (project_manager.project_manager_name LIKE "+searchKey+") OR (date_of_resignation LIKE "+searchKey+") OR (date_of_leaving LIKE "+searchKey+") OR (status LIKE "+searchKey+")) LIMIT ?, ?", UserName, Offset, Limit)
		catch(err)
		defer result.Close()
		var co int
		var posts []models.Resignations
		for result.Next() {
			err := result.Scan(&post.ID, &post.Empname, &post.Project, &post.Manager, &post.Backfillrequired, &post.Regrenonregre, &post.Status, &post.Dateofresignation, &post.Dateofleaving)
			catch(err)
			posts = append(posts, post)
		}
		if count.Next() != false {
			count.Scan(&co)
		} else {
			co = 0
		}
		defer db.Close()
		//Pagination
		var Pag models.Pagination
		Pag.TotalData = co
		Pag.Limit = Limit
		Pag.Data = posts
		x1 := co / Pag.Limit
		x := int(co) % Pag.Limit
		x2 := x1 + 1
		if x == 0 {
			Pag.TotalPages = x1
		} else {
			Pag.TotalPages = x2
		}
		x, _ = strconv.Atoi(Pages)
		if Pag.TotalPages != 0 {
			x1 = (x + 1)
		}
		Pag.Page = x1
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Pag)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}
}

// Soft delete
func (c *Commander) DeleteResignations(w http.ResponseWriter, r *http.Request) {
	db := database.DbConn()
	defer db.Close()
	var ManagerID int
	SetupResponse(&w, r)
	if (strings.Contains(Role, "project manager")) == true {
		var post models.Resignations
		var email string
		err := json.NewDecoder(r.Body).Decode(&post)
		catch(err)
		stmt1, err := db.Query("select sub_project_manager_id from active_resignations where id=?", post.ID)
		catch(err)
		defer stmt1.Close()
		if stmt1.Next() != false {
			err := stmt1.Scan(&ManagerID)
			catch(err)
		}
		if ManagerID != 0 {
			stmt2, err := db.Query("select project_manager.project_manager_email from head_count left join sub_project_manager on head_count.sub_project_manager_id=sub_project_manager.id left join sub_project on sub_project_manager.sub_project_id=sub_project.id left join project_manager on sub_project_manager.manager_id=project_manager.id where sub_project_manager.id=? and head_count.is_active=1", ManagerID)
			catch(err)
			defer stmt2.Close()
			if stmt2.Next() != false {
				err := stmt2.Scan(&email)
				catch(err)
			}
			if UserName == email {
				stmt, err := db.Prepare("Update active_resignations set is_active = 0 where id = ?")
				catch(err)
				_, err = stmt.Exec(post.ID)
				catch(err)
				defer stmt.Close()
				respondwithJSON(w, http.StatusOK, map[string]string{"message": "deleted successfully"})
			} else {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "Project not under you")
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "ID mis matching")
		}
	}
}

//Update The existing record
func (c *Commander) UpdateResignations(w http.ResponseWriter, r *http.Request) {
	db := database.DbConn()
	defer db.Close()
	SetupResponse(&w, r)
	if (strings.Contains(Role, "project manager")) == true {
		var post models.Resignations
		err := json.NewDecoder(r.Body).Decode(&post)
		catch(err)
		var ManagerID int
		stmt1, err := db.Query("select id from sub_project_manager where sub_project_id in (select id from sub_project where sub_project_name= ?) and manager_id in (select id from project_manager where project_manager_email= ?)", post.Project, UserName)
		catch(err)
		defer stmt1.Close()
		if stmt1.Next() != false {
			err := stmt1.Scan(&ManagerID)
			catch(err)
		}
		if ManagerID != 0 {
			query, err := db.Prepare("Update active_resignations set emp_name = ?, backfill_required = ?, regre_non_regre = ?, status = ?, date_of_resignation = ?, date_of_leaving = ?, updated_at = ? where id = ?")
			catch(err)
			update := time.Now()
			_, err = query.Exec(post.Empname, post.Backfillrequired, post.Regrenonregre, post.Status, post.Dateofresignation, post.Dateofleaving, update, post.ID)
			catch(err)
			defer query.Close()
			respondwithJSON(w, http.StatusOK, map[string]string{"message": "update successfully"})
		} else {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Project is not under you")
		}
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Unauthorised access")
	}
}
func catch(err error) {
	if err != nil {
		WriteLogFile(err)
		panic(err)
	}
}
func respondwithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	fmt.Println(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
