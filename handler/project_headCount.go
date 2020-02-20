package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	database "projectname_projectmanager/driver"
	models "projectname_projectmanager/model"
	"strings"
)

//Get project wise head count
func (c *Commander) GetSingleProjectCount(w http.ResponseWriter, r *http.Request) {
	db := database.DbConn()
	defer db.Close()
	w.Header().Set("Content-Type", "application/json")
	SetupResponse(&w, r)
	offsets, ok := r.URL.Query()["project"]
	if !ok || offsets[0] == "" {
		w.WriteHeader(http.StatusNotFound)
		respondwithJSON(w, http.StatusOK, map[string]string{"message": "url parameter project is missing"})
		return
	}
	project := offsets[0]
	// i, _ := strconv.Atoi(pages)
	// offset := i * Limit
	if (strings.Contains(Role, "program manager")) == true {
		var post models.HeadCount
		// err := json.NewDecoder(r.Body).Decode(&post)
		// catch(err)
		result, err := db.Query("select billables_count, billing_on_hold, vt_count, pi_i_count, others from head_count left join sub_project_manager on head_count.sub_project_manager_id = sub_project_manager.id left join sub_project on  sub_project_manager.sub_project_id = sub_project.id left join project on sub_project.project_id = project.id left join program_manager on project.program_manager_id=program_manager.id left join project_manager on sub_project_manager.manager_id=project_manager.id WHERE sub_project.sub_project_name = ? AND head_count.is_active = 1", project)
		catch(err)
		defer result.Close()
		for result.Next() {
			//var post models.Bottom
			err := result.Scan(&post.BillablesCount, &post.BillingOnHold, &post.VtCount, &post.PiICount, &post.Others)
			catch(err)
			//posts = append(posts, post)
		}
		json.NewEncoder(w).Encode(post)
	} else if (strings.Contains(Role, "project manager")) == true {
		var post models.HeadCount
		// err := json.NewDecoder(r.Body).Decode(&post)
		// catch(err)
		result, err := db.Query("select billables_count, billing_on_hold, vt_count, pi_i_count, others from head_count left join sub_project_manager on head_count.sub_project_manager_id = sub_project_manager.id left join sub_project on  sub_project_manager.sub_project_id = sub_project.id left join project_manager on sub_project_manager.manager_id = project_manager.id WHERE sub_project.sub_project_name = ? AND head_count.is_active = 1", project)
		catch(err)
		defer result.Close()
		for result.Next() {
			//var post models.Bottom
			err := result.Scan(&post.BillablesCount, &post.BillingOnHold, &post.VtCount, &post.PiICount, &post.Others)
			catch(err)
			//posts = append(posts, post)
		}
		json.NewEncoder(w).Encode(post)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Unauthorised access")
	}
}
