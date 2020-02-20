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

//Get all retention.
func (c *Commander) Getallretentions(w http.ResponseWriter, r *http.Request) {
	db := database.DbConn()
	defer db.Close()
	SetupResponse(&w, r)
	w.Header().Set("Content-Type", "application/json")
	offsets, ok := r.URL.Query()["pages"]
	if !ok || offsets[0] == "" {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Url parameter pages is missing")
		return
	}
	pages := offsets[0]
	i, _ := strconv.Atoi(pages)
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
	offset := i * Limit
	if (strings.Contains(Role, "program manager")) == true {
		var posts []models.Retention
		var totalretention models.Totalretentions
		var Pag models.Pagination
		// offsets, ok := r.URL.Query()["pages"]
		// if !ok || offsets[0] == "" {
		// 	w.WriteHeader(http.StatusNotFound)
		// 	fmt.Fprintf(w, "Url parameter pages is missing")
		// 	return
		// }
		// pages := offsets[0]
		// i, _ := strconv.Atoi(pages)
		// limit, ok := r.URL.Query()["limit"]
		// if !ok || limit[0] == "" {
		// 	w.WriteHeader(http.StatusNotFound)
		// 	fmt.Fprintf(w, "Url parameter limit is missing")
		// 	return
		// }
		// limits := limit[0]
		// Limit, _ := strconv.Atoi(limits)
		// if Limit == 0 {
		// 	w.WriteHeader(http.StatusBadGateway)
		// 	fmt.Fprintf(w, "Limit can't be 0")
		// 	return
		// }
		// offset := i * Limit
		result, err := db.Query("select retention.id, sub_project.sub_project_name, project_manager.project_manager_name, retention_initiated, retained from retention left join sub_project_manager on retention.sub_project_manager_id = sub_project_manager.id left join sub_project on  sub_project_manager.sub_project_id = sub_project.id left join project on sub_project.project_id = project.id left join program_manager on project.program_manager_id=program_manager.id left join project_manager on sub_project_manager.manager_id=project_manager.id WHERE program_manager.program_manager_email = ? AND retention.is_active = 1 LIMIT ?, ?", UserName, offset, Limit)
		catch(err)
		for result.Next() {
			var post models.Retention
			err := result.Scan(&post.ID, &post.ProjectName, &post.ProjectManagerName, &post.RetentionInitiated, &post.Retained)
			catch(err)
			posts = append(posts, post)
		}
		defer result.Close()
		count1, err := db.Query("SELECT count(id) from retention WHERE is_active = 1 AND sub_project_manager_id in (SELECT id from sub_project_manager where sub_project_id in (select id from sub_project where project_id in (select id from project where program_manager_id in (select id from program_manager where program_manager_email= ?))))", UserName)
		catch(err)
		defer count1.Close()
		for count1.Next() {
			err := count1.Scan(&Pag.TotalData)
			catch(err)
		}
		count, err := db.Query("select sum(retention_initiated), sum(retained) from retention left join sub_project_manager on retention.sub_project_manager_id = sub_project_manager.id left join sub_project on  sub_project_manager.sub_project_id = sub_project.id left join project on sub_project.project_id = project.id left join program_manager on project.program_manager_id=program_manager.id left join project_manager on sub_project_manager.manager_id=project_manager.id WHERE program_manager.program_manager_email = ? AND retention.is_active = 1", UserName)
		catch(err)
		defer count.Close()
		for count.Next() {
			err := count.Scan(&totalretention.TotalRetentionInitiated, &totalretention.TotalRetained)
			catch(err)
		}
		totalretention.Data = posts
		Pag.Data = totalretention
		Pag.Limit = Limit
		Pag.TotalPages = (Pag.TotalData / Pag.Limit) + 1
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
		projectManagerID, err := db.Query("select id from sub_project_manager where manager_id in (select id from project_manager where project_manager_email= ? AND sub_project_manager.is_active=1)", UserName)
		catch(err)
		defer projectManagerID.Close()
		for projectManagerID.Next() {
			projectManagerID.Scan(&ManagerID)
			stmt2, err := db.Query("SELECT id from retention WHERE sub_project_manager_id = ? ", ManagerID)
			catch(err)
			if stmt2.Next() == false {

				stmt, err := db.Prepare("INSERT INTO retention(sub_project_manager_id,created_at, updated_at) VALUES(?, now(), now())")
				catch(err)
				defer stmt.Close()
				// var Total int
				// var Sno int
				_, err = stmt.Exec(ManagerID)
				//fmt.Println(ManagerID)
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
			stmt3, err := db.Query("SELECT id from retention WHERE sub_project_manager_id = ? AND retention.is_active=1", deletedManagerID)
			catch(err)
			if stmt3.Next() == true {
				stmt3.Scan(&deleteIDfromtable)
				stmt4, err := db.Prepare("Update retention set is_active = 0 where id = ?")
				catch(err)
				_, err = stmt4.Exec(deleteIDfromtable)
				catch(err)
				defer stmt4.Close()

			}
			defer stmt3.Close()
		}
		//	if duplicateID != 0 {
		// var posts []models.Retention
		// var totalretention models.Totalretentions
		// var Pag models.Pagination
		// offsets, ok := r.URL.Query()["pages"]
		// if !ok || offsets[0] == "" {
		// 	w.WriteHeader(http.StatusNotFound)
		// 	fmt.Fprintf(w, "Url parameter pages is missing")
		// 	return
		// }
		// pages := offsets[0]
		// i, _ := strconv.Atoi(pages)
		// limit, ok := r.URL.Query()["limit"]
		// if !ok || limit[0] == "" {
		// 	w.WriteHeader(http.StatusNotFound)
		// 	fmt.Fprintf(w, "Url parameter limit is missing")
		// 	return
		// }
		// limits := limit[0]
		// Limit, _ := strconv.Atoi(limits)
		// if Limit == 0 {
		// 	w.WriteHeader(http.StatusBadGateway)
		// 	fmt.Fprintf(w, "Limit can't be 0")
		// 	return
		// }
		// offset := i * Limit
		// result, err := db.Query("select retention.id, sub_project.sub_project_name, project_manager.project_manager_name, retention_initiated, retained from retention left join sub_project_manager on retention.sub_project_manager_id = sub_project_manager.id left join sub_project on  sub_project_manager.sub_project_id = sub_project.id left join project_manager on sub_project_manager.manager_id = project_manager.id WHERE project_manager.project_manager_email = ? AND retention.is_active = 1 LIMIT ?, ?", UserName, offset, Limit)
		// catch(err)
		// for result.Next() {
		// 	var post models.Retention
		// 	err := result.Scan(&post.ID, &post.ProjectName, &post.ProjectManagerName, &post.RetentionInitiated, &post.Retained)
		// 	catch(err)
		// 	posts = append(posts, post)
		// }
		// defer result.Close()
		// count, err := db.Query("select sum(retention_initiated), sum(retained) from retention left join sub_project_manager on retention.sub_project_manager_id = sub_project_manager.id left join sub_project on  sub_project_manager.sub_project_id = sub_project.id left join project_manager on sub_project_manager.manager_id = project_manager.id WHERE project_manager.project_manager_email = ? AND retention.is_active = 1", UserName)
		// catch(err)
		// defer count.Close()
		// for count.Next() {
		// 	err := count.Scan(&totalretention.TotalRetentionInitiated, &totalretention.TotalRetained)
		// 	catch(err)
		// }
		// count1, err := db.Query("SELECT count(id) from retention WHERE (is_active = 1) AND (sub_project_manager_id in (SELECT id from sub_project_manager where sub_project_id in (select id from sub_project) AND (manager_id in (select id from project_manager where project_manager_email =?))))", UserName)
		// catch(err)
		// defer count1.Close()
		// for count1.Next() {
		// 	err := count1.Scan(&Pag.TotalData)
		// 	catch(err)
		// }
		// totalretention.Data = posts
		// Pag.Data = totalretention
		// Pag.Limit = Limit
		// Pag.TotalPages = (Pag.TotalData / Pag.Limit) + 1
		// x1 := Pag.TotalData / Pag.Limit
		// x := Pag.TotalData % Pag.Limit
		// x2 := x1 + 1
		// if x == 0 {
		// 	Pag.TotalPages = x1
		// } else {
		// 	Pag.TotalPages = x2
		// }
		// x, _ = strconv.Atoi(pages)
		// if Pag.TotalPages != 0 {
		// 	x1 = x + 1
		// }
		// Pag.Page = x1
		// json.NewEncoder(w).Encode(Pag)
		//} else {
		// var ManagerID int
		// //var post models.HeadCount
		// stmt1, err := db.Query("select id from sub_project_manager where manager_id in (select id from project_manager where project_manager_email= ?)", UserName)
		// catch(err)
		// defer stmt1.Close()
		// for stmt1.Next() != false {
		// 	err := stmt1.Scan(&ManagerID)
		// 	catch(err)

		// 	stmt, err := db.Prepare("INSERT INTO retention(sub_project_manager_id,created_at, updated_at) VALUES(?, now(), now())")
		// 	catch(err)
		// 	defer stmt.Close()
		// 	// var Total int
		// 	// var Sno int
		///// 	_, err = stmt.Exec(ManagerID)
		// 	//fmt.Println(ManagerID)
		// 	catch(err)
		// }
		var posts []models.Retention
		var totalretention models.Totalretentions
		var Pag models.Pagination
		// offsets, ok := r.URL.Query()["pages"]
		// if !ok || offsets[0] == "" {
		// 	w.WriteHeader(http.StatusNotFound)
		// 	fmt.Fprintf(w, "Url parameter pages is missing")
		// 	return
		// }
		// pages := offsets[0]
		// i, _ := strconv.Atoi(pages)
		// limit, ok := r.URL.Query()["limit"]
		// if !ok || limit[0] == "" {
		// 	w.WriteHeader(http.StatusNotFound)
		// 	fmt.Fprintf(w, "Url parameter limit is missing")
		// 	return
		// }
		// limits := limit[0]
		// Limit, _ := strconv.Atoi(limits)
		// if Limit == 0 {
		// 	w.WriteHeader(http.StatusBadGateway)
		// 	fmt.Fprintf(w, "Limit can't be 0")
		// 	return
		// }
		// offset := i * Limit
		result, err := db.Query("select retention.id, sub_project.sub_project_name, retention_initiated, retained from retention left join sub_project_manager on retention.sub_project_manager_id = sub_project_manager.id left join sub_project on  sub_project_manager.sub_project_id = sub_project.id left join project_manager on sub_project_manager.manager_id = project_manager.id WHERE project_manager.project_manager_email = ? AND retention.is_active = 1 LIMIT ?, ?", UserName, offset, Limit)
		catch(err)
		for result.Next() {
			var post models.Retention
			err := result.Scan(&post.ID, &post.ProjectName, &post.RetentionInitiated, &post.Retained)
			catch(err)
			posts = append(posts, post)
		}
		defer result.Close()
		count, err := db.Query("select sum(retention_initiated), sum(retained) from retention left join sub_project_manager on retention.sub_project_manager_id = sub_project_manager.id left join sub_project on  sub_project_manager.sub_project_id = sub_project.id left join project_manager on sub_project_manager.manager_id = project_manager.id WHERE project_manager.project_manager_email = ? AND retention.is_active = 1", UserName)
		catch(err)
		defer count.Close()
		for count.Next() {
			err := count.Scan(&totalretention.TotalRetentionInitiated, &totalretention.TotalRetained)
			catch(err)
		}
		count1, err := db.Query("SELECT count(id) from retention WHERE (is_active = 1) AND (sub_project_manager_id in (SELECT id from sub_project_manager where sub_project_id in (select id from sub_project) AND (manager_id in (select id from project_manager where project_manager_email =?))))", UserName)
		catch(err)
		defer count1.Close()
		for count1.Next() {
			err := count1.Scan(&Pag.TotalData)
			catch(err)
		}
		totalretention.Data = posts
		Pag.Data = totalretention
		Pag.Limit = Limit
		Pag.TotalPages = (Pag.TotalData / Pag.Limit) + 1
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

//To create retention
func (c *Commander) Createretentions(w http.ResponseWriter, r *http.Request) {
	db := database.DbConn()
	defer db.Close()
	SetupResponse(&w, r)
	w.Header().Set("Content-Type", "application/json")
	var ManagerID int
	if (strings.Contains(Role, "project manager")) == true {
		var post models.Retention
		json.NewDecoder(r.Body).Decode(&post)
		stmt1, err1 := db.Query("select id from sub_project_manager where sub_project_id in (select id from sub_project where sub_project_name= ?) and manager_id in (select id from project_manager where project_manager_email= ?)", post.ProjectName, UserName)
		if err1 != nil {
			panic(err1.Error())
		}
		defer stmt1.Close()
		if stmt1.Next() != false {
			err2 := stmt1.Scan(&ManagerID)
			if err2 != nil {
				panic(err2.Error())
			}
		}
		if ManagerID != 0 {
			var dublicateID int
			query := db.QueryRow("SELECT id from retention where (sub_project_manager_id = ? AND retention_initiated = ? AND retained = ?) AND is_active=1", ManagerID, post.RetentionInitiated, post.Retained)
			query.Scan(&dublicateID)
			if dublicateID == 0 {
				stmt, err := db.Prepare("INSERT INTO retention(sub_project_manager_id, retention_initiated, retained, created_at, updated_at) VALUES(?, ?, ?, now(), now())")
				catch(err)
				stmt.Exec(ManagerID, post.RetentionInitiated, post.Retained)
				defer stmt.Close()
				w.WriteHeader(http.StatusCreated)
				fmt.Fprintf(w, "New post was created")
			} else {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "Duplicates record found")
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Project not under you")
		}
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Unauthorised access")
	}

}

// Soft delete
func (c *Commander) Deleteretentions(w http.ResponseWriter, r *http.Request) {
	db := database.DbConn()
	defer db.Close()
	w.Header().Set("Content-Type", "application/json")
	SetupResponse(&w, r)
	var ManagerID int
	if (strings.Contains(Role, "project manager")) == true {
		var post models.Retention
		var email string
		json.NewDecoder(r.Body).Decode(&post)
		stmt1, err := db.Query("select sub_project_manager_id from retention where id=?", post.ID)
		catch(err)
		defer stmt1.Close()
		if stmt1.Next() != false {
			err := stmt1.Scan(&ManagerID)
			catch(err)
		}
		if ManagerID != 0 {
			stmt2, err := db.Query("select project_manager.project_manager_email from retention left join sub_project_manager on retention.sub_project_manager_id=sub_project_manager.id left join sub_project on sub_project_manager.sub_project_id=sub_project.id left join project_manager on sub_project_manager.manager_id=project_manager.id where sub_project_manager.id=? and retention.is_active=1", ManagerID)
			catch(err)
			defer stmt2.Close()
			if stmt2.Next() != false {
				err := stmt2.Scan(&email)
				catch(err)
			}
			if UserName == email {
				var dublicateID int
				query := db.QueryRow("SELECT id from retention where id = ? AND is_active=0", post.ID)
				query.Scan(&dublicateID)
				if dublicateID == 0 {
					stmt, err := db.Prepare("Update retention set is_active = 0 where id = ?")
					if err != nil {
						panic(err.Error())
					}
					_, err = stmt.Exec(post.ID)
					if err != nil {
						panic(err.Error())
					}
					defer stmt.Close()
					respondwithJSON(w, http.StatusOK, map[string]string{"message": "deleted successfully"})
				} else {
					w.WriteHeader(http.StatusBadRequest)
					fmt.Fprintf(w, "Record does not exists")
				}
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

//To update retention table.
func (c *Commander) Updateretentions(w http.ResponseWriter, r *http.Request) {
	db := database.DbConn()
	defer db.Close()
	SetupResponse(&w, r)
	w.Header().Set("Content-Type", "application/json")
	if (strings.Contains(Role, "project manager")) == true {
		var post models.Retention
		var ManagerID int
		json.NewDecoder(r.Body).Decode(&post)
		stmt1, err1 := db.Query("select id from sub_project_manager where sub_project_id in (select id from sub_project where sub_project_name= ?) and manager_id in (select id from project_manager where project_manager_email= ?)", post.ProjectName, UserName)
		if err1 != nil {
			panic(err1.Error())
		}
		defer stmt1.Close()
		if stmt1.Next() != false {
			err2 := stmt1.Scan(&ManagerID)
			if err2 != nil {
				panic(err2.Error())
			}
		}
		if ManagerID != 0 {
			query, err := db.Prepare("Update retention set retention_initiated = ?, retained = ?, updated_at = ? where id = ?")
			catch(err)
			update := time.Now()
			fmt.Println(update)
			_, er := query.Exec(post.RetentionInitiated, post.Retained, update, post.ID)
			catch(er)
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

//Generic search on retention.
func (c *Commander) Getretentionbyprojectname(w http.ResponseWriter, r *http.Request) {
	db := database.DbConn()
	defer db.Close()
	SetupResponse(&w, r)
	w.Header().Set("Content-Type", "application/json")
	if (strings.Contains(Role, "program manager")) == true {
		var posts []models.Retention
		var totalretention models.Totalretentions
		var Pag models.Pagination
		params := mux.Vars(r)
		key := params["projects.project_name"]
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
		count1, _ := db.Query("SELECT count(retention.id) from retention left join sub_project_manager on retention.sub_project_manager_id = sub_project_manager.id left join sub_project on  sub_project_manager.sub_project_id = sub_project.id left join project on sub_project.project_id = project.id left join program_manager on project.program_manager_id=program_manager.id left join project_manager on sub_project_manager.manager_id=project_manager.id WHERE program_manager.program_manager_email = ? AND retention.is_active = 1 AND ((sub_project.sub_project_name LIKE "+searchKey+") OR (project_manager.project_manager_name LIKE "+searchKey+"))", UserName)
		defer count1.Close()
		result, err := db.Query("select retention.id, sub_project.sub_project_name, project_manager.project_manager_name, retention_initiated, retained from retention left join sub_project_manager on retention.sub_project_manager_id = sub_project_manager.id left join sub_project on  sub_project_manager.sub_project_id = sub_project.id left join project on sub_project.project_id = project.id left join program_manager on project.program_manager_id=program_manager.id left join project_manager on sub_project_manager.manager_id=project_manager.id WHERE program_manager.program_manager_email = ? AND retention.is_active = 1 AND ((sub_project.sub_project_name LIKE "+searchKey+") OR (project_manager.project_manager_name LIKE "+searchKey+")) LIMIT ?,?", UserName, Offset, Limit)
		catch(err)
		for result.Next() {
			var post models.Retention
			result.Scan(&post.ID, &post.ProjectName, &post.ProjectManagerName, &post.RetentionInitiated, &post.Retained)
			posts = append(posts, post)
		}
		defer result.Close()
		count, err := db.Query("select ifnull(sum(retention_initiated), 0), ifnull(sum(retained), 0) from retention left join sub_project_manager on retention.sub_project_manager_id = sub_project_manager.id left join sub_project on  sub_project_manager.sub_project_id = sub_project.id left join project on sub_project.project_id = project.id left join program_manager on project.program_manager_id=program_manager.id left join project_manager on sub_project_manager.manager_id=project_manager.id WHERE program_manager.program_manager_email = ? AND retention.is_active = 1 AND ((sub_project.sub_project_name LIKE "+searchKey+") OR (project_manager.project_manager_name LIKE "+searchKey+"))", UserName)
		catch(err)
		for count.Next() {
			err := count.Scan(&totalretention.TotalRetentionInitiated, &totalretention.TotalRetained)
			catch(err)
		}
		defer count.Close()
		var co int //Total data
		if count1.Next() != false {
			count1.Scan(&co)
		} else {
			co = 0
		}
		totalretention.Data = posts
		Pag.Data = totalretention
		Pag.TotalData = co
		Pag.Limit = Limit
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
	} else if (strings.Contains(Role, "project manager")) == true {
		var posts []models.Retention
		var totalretention models.Totalretentions
		var Pag models.Pagination
		params := mux.Vars(r)
		key := params["projects.project_name"]
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
		count1, _ := db.Query("SELECT count(retention.id) from retention left join sub_project_manager on retention.sub_project_manager_id = sub_project_manager.id left join sub_project on  sub_project_manager.sub_project_id = sub_project.id left join project_manager on sub_project_manager.manager_id = project_manager.id WHERE project_manager.project_manager_email = ? AND retention.is_active = 1 AND ((sub_project.sub_project_name LIKE "+searchKey+") OR (project_manager.project_manager_name LIKE "+searchKey+"))", UserName)
		defer count1.Close()
		result, err := db.Query("select retention.id, sub_project.sub_project_name, retention_initiated, retained from retention left join sub_project_manager on retention.sub_project_manager_id = sub_project_manager.id left join sub_project on  sub_project_manager.sub_project_id = sub_project.id left join project_manager on sub_project_manager.manager_id = project_manager.id WHERE project_manager.project_manager_email = ? AND retention.is_active = 1 AND ((sub_project.sub_project_name LIKE "+searchKey+") OR (project_manager.project_manager_name LIKE "+searchKey+")) LIMIT ?,?", UserName, Offset, Limit)
		catch(err)
		for result.Next() {
			var post models.Retention
			result.Scan(&post.ID, &post.ProjectName, &post.RetentionInitiated, &post.Retained)
			posts = append(posts, post)
		}
		defer result.Close()
		count, err := db.Query("select ifnull(sum(retention_initiated), 0), ifnull(sum(retained), 0) from retention left join sub_project_manager on retention.sub_project_manager_id = sub_project_manager.id left join sub_project on  sub_project_manager.sub_project_id = sub_project.id left join project_manager on sub_project_manager.manager_id = project_manager.id WHERE project_manager.project_manager_email = ? AND retention.is_active = 1 AND ((sub_project.sub_project_name LIKE "+searchKey+") OR (project_manager.project_manager_name LIKE "+searchKey+"))", UserName)
		catch(err)
		for count.Next() {
			err := count.Scan(&totalretention.TotalRetentionInitiated, &totalretention.TotalRetained)
			catch(err)
		}
		defer count.Close()
		var co int // total count
		if count1.Next() != false {
			count1.Scan(&co)
		} else {
			co = 0
		}
		totalretention.Data = posts
		Pag.Data = totalretention
		Pag.TotalData = co
		Pag.Limit = Limit
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
	}
}
