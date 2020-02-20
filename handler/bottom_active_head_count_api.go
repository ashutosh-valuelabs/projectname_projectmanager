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
func (c *Commander) Getprojectheadcount(w http.ResponseWriter, r *http.Request) {
	db := database.DbConn()
	defer db.Close()
	w.Header().Set("Content-Type", "application/json")
	SetupResponse(&w, r)
	if (strings.Contains(Role, "program manager")) == true {
		var posts []models.Bottom
		result, err := db.Query("select projects.project_name, sum(net) from head_count left join sub_project_manager on head_count.sub_project_manager_id = sub_project_manager.id left join sub_project on  sub_project_manager.sub_project_id = sub_project.id left join project on sub_project.project_id = project.id left join program_manager on project.program_manager_id=program_manager.id left join project_manager on sub_project_manager.manager_id=project_manager.id WHERE program_manager.program_manager_email = ? AND head_count.is_active = 1", UserName)
		catch(err)
		defer result.Close()
		for result.Next() {
			var post models.Bottom
			err := result.Scan(&post.Project, &post.HeadCount)
			catch(err)
			posts = append(posts, post)
		}
		json.NewEncoder(w).Encode(posts)
	} else if (strings.Contains(Role, "project manager")) == true {
		var posts []models.Bottom
		result, err := db.Query("select projects.project_name, sum(net) from head_count left join sub_project_manager on head_count.sub_project_manager_id = sub_project_manager.id left join sub_project on  sub_project_manager.sub_project_id = sub_project.id left join project_manager on sub_project_manager.manager_id = project_manager.id WHERE project_manager.project_manager_email = ? AND head_count.is_active = 1", UserName)
		catch(err)
		defer result.Close()
		for result.Next() {
			var post models.Bottom
			err := result.Scan(&post.Project, &post.HeadCount)
			catch(err)
			posts = append(posts, post)
		}
		json.NewEncoder(w).Encode(posts)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Unauthorised access")
	}
}
