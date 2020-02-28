package handler

import (
	"encoding/json"
	"net/http"
	database "projectname_projectmanager/driver"
	models "projectname_projectmanager/model"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

//Get all head count
func (c *Commander) Getallheadcount(w http.ResponseWriter, r *http.Request) {
	db := database.DbConn()
	defer db.Close()
	SetupResponse(&w, r)
	w.Header().Set("Content-Type", "application/json")
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
		respondwithJSON(w, http.StatusOK, map[string]string{"message": "url parameter limit is missing"})
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
	if (strings.Contains(Role, "program manager")) == true {
		var posts []models.HeadCount
		var Pag models.Pagination
		result, err := db.Query("select head_count.id, sub_project.sub_project_name, project_manager.project_manager_name, billables_count, billing_on_hold, vt_count, pi_i_count, others, net from head_count left join sub_project_manager on head_count.sub_project_manager_id = sub_project_manager.id left join sub_project on  sub_project_manager.sub_project_id = sub_project.id left join project on sub_project.project_id = project.id left join program_manager on project.program_manager_id=program_manager.id left join project_manager on sub_project_manager.manager_id=project_manager.id WHERE program_manager.program_manager_email = ? AND head_count.is_active = 1 LIMIT ?, ?", UserName, offset, Limit)
		catch(err)
		defer result.Close()
		count, err := db.Query("SELECT count(id) from head_count WHERE is_active = 1 AND sub_project_manager_id in (SELECT id from sub_project_manager where sub_project_id in (select id from sub_project where project_id in (select id from project where program_manager_id in (select id from program_manager where program_manager_email= ?))))", UserName)
		catch(err)
		defer count.Close()
		for count.Next() {
			err := count.Scan(&Pag.TotalData)
			catch(err)
		}
		Pag.Limit = Limit
		Pag.TotalPages = (Pag.TotalData / Pag.Limit) + 1
		for result.Next() {
			var post models.HeadCount
			err := result.Scan(&post.ID, &post.ProjectName, &post.ManagerName, &post.BillablesCount, &post.BillingOnHold, &post.VtCount, &post.PiICount, &post.Others, &post.Net)
			catch(err)
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
	} else if (strings.Contains(Role, "project manager")) == true {
		var ManagerID int
		projectManagerID, err := db.Query("select id from sub_project_manager where manager_id in (select id from project_manager where project_manager_email= ?) AND sub_project_manager.is_active=1", UserName)
		catch(err)
		defer projectManagerID.Close()
		for projectManagerID.Next() {
			projectManagerID.Scan(&ManagerID)
			//var duplicateID int
			stmt2, err := db.Query("SELECT id from head_count WHERE sub_project_manager_id = ? ", ManagerID)
			catch(err)
			if stmt2.Next() == false {

				stmt, err := db.Prepare("INSERT INTO head_count(sub_project_manager_id,created_at, updated_at) VALUES(?, now(), now())")
				catch(err)
				defer stmt.Close()
				// var Total int
				// var Sno int
				_, err = stmt.Exec(ManagerID)
				catch(err)

			}
			stmt2.Close()

		}
		var deletedManagerID int
		var deleteIDfromtable int
		projectManagerID1, err := db.Query("select id from sub_project_manager where manager_id in (select id from project_manager where project_manager_email= ?) AND sub_project_manager.is_active=0", UserName)
		catch(err)
		defer projectManagerID1.Close()
		for projectManagerID1.Next() {
			projectManagerID1.Scan(&deletedManagerID)
			stmt3, err := db.Query("SELECT id from head_count WHERE sub_project_manager_id = ? AND head_count.is_active=1", deletedManagerID)
			catch(err)
			if stmt3.Next() == true {
				stmt3.Scan(&deleteIDfromtable)
				stmt4, err := db.Prepare("delete from head_count where id = ?")
				catch(err)
				_, err = stmt4.Exec(deleteIDfromtable)
				catch(err)
				defer stmt4.Close()

			}
			defer stmt3.Close()
		}

		var posts []models.HeadCount
		var Pag models.Pagination

		Pag.Limit = Limit
		// offset := i * Limit
		result, err := db.Query("select head_count.id, sub_project.sub_project_name, billables_count, billing_on_hold, vt_count, pi_i_count, others, net from head_count left join sub_project_manager on head_count.sub_project_manager_id = sub_project_manager.id left join sub_project on  sub_project_manager.sub_project_id = sub_project.id left join project_manager on sub_project_manager.manager_id = project_manager.id WHERE project_manager.project_manager_email = ? AND head_count.is_active = 1 LIMIT ?, ?", UserName, offset, Limit)
		catch(err)
		defer result.Close()
		count, err := db.Query("SELECT count(id) from head_count WHERE (is_active = 1) AND (sub_project_manager_id in (SELECT id from sub_project_manager where sub_project_id in (select id from sub_project) AND (manager_id in (select id from project_manager where project_manager_email =?))))", UserName)
		catch(err)
		defer count.Close()
		for count.Next() {
			err := count.Scan(&Pag.TotalData)
			catch(err)
		}
		Pag.TotalPages = (Pag.TotalData / Pag.Limit) + 1
		for result.Next() {
			var post models.HeadCount
			err := result.Scan(&post.ID, &post.ProjectName, &post.BillablesCount, &post.BillingOnHold, &post.VtCount, &post.PiICount, &post.Others, &post.Net)
			catch(err)
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

//creating head count for a project
func (c *Commander) Createheadcount(w http.ResponseWriter, r *http.Request) {
	db := database.DbConn()
	defer db.Close()
	// SetupResponse(&w, r)
	w.Header().Set("Content-Type", "application/json")
	var ManagerID int
	if (strings.Contains(Role, "project manager")) == true {
		var post models.HeadCount
		err := json.NewDecoder(r.Body).Decode(&post)
		catch(err)
		stmt1, err := db.Query("select id from sub_project_manager where sub_project_id in (select id from sub_project where sub_project_name= ?) and manager_id in (select id from project_manager where project_manager_email= ?)", post.ProjectName, UserName)
		catch(err)
		defer stmt1.Close()
		if stmt1.Next() != false {
			err := stmt1.Scan(&ManagerID)
			catch(err)
		}
		if ManagerID != 0 {
			var duplicateID int
			query := db.QueryRow("SELECT id from head_count where (sub_project_manager_id = ? OR billables_count = ? OR billing_on_hold = ? OR vt_count = ? OR pi_i_count = ? OR others = ?) AND is_active=1", ManagerID, post.BillablesCount, post.BillingOnHold, post.VtCount, post.PiICount, post.Others)
			query.Scan(&duplicateID)
			if duplicateID == 0 {
				stmt, err := db.Prepare("INSERT INTO head_count(sub_project_manager_id, billables_count, billing_on_hold, vt_count, pi_i_count, others, created_at, updated_at) VALUES(?, ?, ?, ?, ?, ?, now(), now())")
				catch(err)
				defer stmt.Close()
				var Total int
				var Sno int
				_, err = stmt.Exec(ManagerID, post.BillablesCount, post.BillingOnHold, post.VtCount, post.PiICount, post.Others)
				catch(err)
				rows, err := db.Query("select id, ifnull(billables_count, 0) + ifnull(billing_on_hold, 0) + ifnull(vt_count, 0) + ifnull(pi_i_count, 0) + ifnull(others, 0) as total from head_count")
				defer rows.Close()
				catch(err)
				for rows.Next() {
					rows.Scan(&Sno, &Total)
				}
				stmt, err = db.Prepare("update head_count set net = ? where id = ?")
				catch(err)
				_, err = stmt.Exec(Total, Sno)

				catch(err)
				w.WriteHeader(http.StatusCreated)
				respondwithJSON(w, http.StatusOK, map[string]string{"message": "new post was created"})
			} else {
				w.WriteHeader(http.StatusBadRequest)
				respondwithJSON(w, http.StatusOK, map[string]string{"message": "buplicate record was found"})
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
			respondwithJSON(w, http.StatusOK, map[string]string{"message": "project not under you"})
		}
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		respondwithJSON(w, http.StatusOK, map[string]string{"message": "unauthorised access"})
	}

}

//Generic search
func (c *Commander) Getheadcountbyname(w http.ResponseWriter, r *http.Request) {
	db := database.DbConn()
	defer db.Close()
	SetupResponse(&w, r)
	w.Header().Set("Content-Type", "application/json")
	if (strings.Contains(Role, "project manager")) == true {
		var post models.HeadCount
		params := mux.Vars(r)
		key := params["projects.project_name"]
		var searchKey string = "'" + key + "%'"
		var Offset int
		Pages := r.FormValue("pages")
		if len(Pages) < 1 {
			w.WriteHeader(http.StatusNotFound)
			respondwithJSON(w, http.StatusOK, map[string]string{"message": "url pages parameter is missing"})
			return
		}
		i1, _ := strconv.Atoi(Pages)
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
			respondwithJSON(w, http.StatusOK, map[string]string{"message": "Limit can't be 0"})
			return
		}
		Offset = i1 * Limit
		count, _ := db.Query("SELECT count(head_count.id) from head_count left join sub_project_manager on head_count.sub_project_manager_id = sub_project_manager.id left join sub_project on  sub_project_manager.sub_project_id = sub_project.id left join project_manager on sub_project_manager.manager_id = project_manager.id WHERE project_manager.project_manager_email = ? AND head_count.is_active = 1 AND ((sub_project.sub_project_name LIKE "+searchKey+") OR (project_manager.project_manager_name LIKE "+searchKey+"))", UserName)
		defer count.Close()
		result, err := db.Query("select head_count.id, sub_project.sub_project_name, billables_count, billing_on_hold, vt_count, pi_i_count, others, net from head_count left join sub_project_manager on head_count.sub_project_manager_id = sub_project_manager.id left join sub_project on  sub_project_manager.sub_project_id = sub_project.id left join project_manager on sub_project_manager.manager_id = project_manager.id WHERE project_manager.project_manager_email = ? AND head_count.is_active = 1 AND ((sub_project.sub_project_name LIKE "+searchKey+") OR (project_manager.project_manager_name LIKE "+searchKey+")) LIMIT ?, ?", UserName, Offset, Limit)
		catch(err)
		defer result.Close()
		var posts []models.HeadCount
		for result.Next() {
			err := result.Scan(&post.ID, &post.ProjectName, &post.BillablesCount, &post.BillingOnHold, &post.VtCount, &post.PiICount, &post.Others, &post.Net)
			catch(err)
			posts = append(posts, post)
		}
		var co int
		if count.Next() != false {
			count.Scan(&co)
		} else {
			co = 0
		}
		var Pag models.Pagination
		Pag.TotalData = co
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
			x1 = (x + 1)
		}
		Pag.Page = x1
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Pag)

	} else if (strings.Contains(Role, "program manager")) == true {
		var post models.HeadCount
		params := mux.Vars(r)
		key := params["projects.project_name"]
		var searchKey string = "'" + key + "%'"
		var Offset int
		Pages := r.FormValue("pages")
		if len(Pages) < 1 {
			w.WriteHeader(http.StatusNotFound)
			respondwithJSON(w, http.StatusOK, map[string]string{"message": "url pages parameter is missing"})
			return
		}
		i1, _ := strconv.Atoi(Pages)
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
		Offset = i1 * Limit
		count, _ := db.Query("SELECT count(head_count.id) from head_count left join sub_project_manager on head_count.sub_project_manager_id = sub_project_manager.id left join sub_project on  sub_project_manager.sub_project_id = sub_project.id left join project on sub_project.project_id = project.id left join program_manager on project.program_manager_id=program_manager.id left join project_manager on sub_project_manager.manager_id=project_manager.id WHERE program_manager.program_manager_email = ? AND head_count.is_active = 1 AND ((sub_project.sub_project_name LIKE "+searchKey+") OR (project_manager.project_manager_name LIKE "+searchKey+"))", UserName)
		defer count.Close()
		result, err := db.Query("select head_count.id, sub_project.sub_project_name, project_manager.project_manager_name, billables_count, billing_on_hold, vt_count, pi_i_count, others, net from head_count left join sub_project_manager on head_count.sub_project_manager_id = sub_project_manager.id left join sub_project on  sub_project_manager.sub_project_id = sub_project.id left join project on sub_project.project_id = project.id left join program_manager on project.program_manager_id=program_manager.id left join project_manager on sub_project_manager.manager_id=project_manager.id WHERE program_manager.program_manager_email = ? AND head_count.is_active = 1 AND ((sub_project.sub_project_name LIKE "+searchKey+") OR (project_manager.project_manager_name LIKE "+searchKey+")) LIMIT ?,?", UserName, Offset, Limit)
		catch(err)
		defer result.Close()
		var posts []models.HeadCount
		for result.Next() {
			err := result.Scan(&post.ID, &post.ProjectName, &post.ManagerName, &post.BillablesCount, &post.BillingOnHold, &post.VtCount, &post.PiICount, &post.Others, &post.Net)
			catch(err)
			posts = append(posts, post)
		}
		var co int
		if count.Next() != false {
			count.Scan(&co)
		} else {
			co = 0
		}
		var Pag models.Pagination
		Pag.TotalData = co
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
			x1 = (x + 1)
		}
		Pag.Page = x1
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Pag)

	} else {
		w.WriteHeader(http.StatusUnauthorized)
		respondwithJSON(w, http.StatusOK, map[string]string{"message": "unauthorised access"})
	}

}

// Soft delete head count
func (c *Commander) Deleteheadcount(w http.ResponseWriter, r *http.Request) {
	db := database.DbConn()
	defer db.Close()
	w.Header().Set("Content-Type", "application/json")
	SetupResponse(&w, r)
	var ManagerID int
	if (strings.Contains(Role, "project manager")) == true {
		var post models.HeadCount
		var email string
		json.NewDecoder(r.Body).Decode(&post)
		stmt1, err := db.Query("select sub_project_manager_id from head_count where id=?", post.ID)
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
				var dublicateID int
				query := db.QueryRow("SELECT id from head_count where id = ? AND is_active=0", post.ID)
				query.Scan(&dublicateID)
				if dublicateID == 0 {
					stmt, err := db.Prepare("Update head_count set is_active = 0 where id = ?")
					catch(err)
					_, err = stmt.Exec(post.ID)
					catch(err)
					defer stmt.Close()
					respondwithJSON(w, http.StatusOK, map[string]string{"message": "deleted successfully"})
				} else {
					w.WriteHeader(http.StatusBadRequest)
					respondwithJSON(w, http.StatusOK, map[string]string{"message": "record doesn't exists"})
				}
			} else {
				w.WriteHeader(http.StatusBadRequest)
				respondwithJSON(w, http.StatusOK, map[string]string{"message": "project not under you"})
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
			respondwithJSON(w, http.StatusOK, map[string]string{"message": "ID is mis-matching"})
		}

	}
}

//For updating head count record
func (c *Commander) Updateheadcount(w http.ResponseWriter, r *http.Request) {
	db := database.DbConn()
	defer db.Close()
	w.Header().Set("Content-Type", "application/json")
	SetupResponse(&w, r)

	if (strings.Contains(Role, "project manager")) == true {
		var post models.HeadCount
		var ManagerID int
		json.NewDecoder(r.Body).Decode(&post)
		stmt1, err := db.Query("select id from sub_project_manager where sub_project_id in (select id from sub_project where sub_project_name= ?) and manager_id in (select id from project_manager where project_manager_email= ?)", post.ProjectName, UserName)
		catch(err)
		defer stmt1.Close()
		if stmt1.Next() != false {
			err := stmt1.Scan(&ManagerID)
			catch(err)
		}
		if ManagerID != 0 {
			query, err := db.Prepare("Update head_count set billables_count = ?, billing_on_hold = ?, vt_count = ?, pi_i_count = ?, others = ?, updated_at = ? where id = ?")
			catch(err)
			update := time.Now()
			_, er := query.Exec(post.BillablesCount, post.BillingOnHold, post.VtCount, post.PiICount, post.Others, update, post.ID)
			catch(er)
			defer query.Close()
			rows, err := db.Query("select ifnull(billables_count, 0) + ifnull(billing_on_hold, 0) + ifnull(vt_count, 0) + ifnull(pi_i_count, 0) + ifnull(others, 0) as total from head_count where id = ?", post.ID)
			if err != nil {
				WriteLogFile(err)
				panic(err.Error())
			} else {
				var Total int
				for rows.Next() {
					rows.Scan(&Total)
				}
				stmt, err := db.Query("update head_count set net = ? where id = ?", Total, post.ID)

				if err != nil {
					panic(err.Error())
				}

				for stmt.Next() {
					stmt.Scan(&post.Net, &post.ID)
				}
			}

			respondwithJSON(w, http.StatusOK, map[string]string{"message": "update successfully"})
		} else {
			w.WriteHeader(http.StatusBadRequest)
			respondwithJSON(w, http.StatusOK, map[string]string{"message": "project is not under you"})
		}
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		respondwithJSON(w, http.StatusOK, map[string]string{"message": "unauthorised access"})
	}

}
