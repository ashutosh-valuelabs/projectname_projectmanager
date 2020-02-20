package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	database "projectname_projectmanager/driver"
	models "projectname_projectmanager/model"
	"strings"
)

//Get total active head count under a program manager.
func (C *Commander) Getactiveheadcount(w http.ResponseWriter, r *http.Request) {
	db := database.DbConn()
	defer db.Close()
	w.Header().Set("Content-Type", "application/json")
	SetupResponse(&w, r)
	if (strings.Contains(Role, "program manager")) == true {
		result, err := db.Query("SELECT sum(net), sum(billables_count), sum(billing_on_hold), sum(vt_count), sum(pi_i_count), sum(others) from head_count left join sub_project_manager on head_count.sub_project_manager_id = sub_project_manager.id left join sub_project on  sub_project_manager.sub_project_id = sub_project.id left join project on sub_project.project_id = project.id left join program_manager on project.program_manager_id=program_manager.id left join project_manager on sub_project_manager.manager_id=project_manager.id WHERE program_manager.program_manager_email = ? AND head_count.is_active = 1", UserName)
		catch(err)
		defer result.Close()
		var post models.Active
		for result.Next() {
			err := result.Scan(&post.ActiveHeadCount, &post.Billable, &post.BillingOnHold, &post.ValueTrade, &post.ProjectInvestment, &post.Others)
			catch(err)
		}
		json.NewEncoder(w).Encode(post)
	} else if (strings.Contains(Role, "project manager")) == true {
		result, err := db.Query("SELECT sum(net), sum(billables_count), sum(billing_on_hold), sum(vt_count), sum(pi_i_count), sum(others) from head_count left join sub_project_manager on head_count.sub_project_manager_id = sub_project_manager.id left join sub_project on  sub_project_manager.sub_project_id = sub_project.id left join project_manager on sub_project_manager.manager_id = project_manager.id WHERE project_manager.project_manager_email = ? AND head_count.is_active = 1", UserName)
		catch(err)
		defer result.Close()
		var post models.Active
		for result.Next() {
			err := result.Scan(&post.ActiveHeadCount, &post.Billable, &post.BillingOnHold, &post.ValueTrade, &post.ProjectInvestment, &post.Others)
			catch(err)
		}
		json.NewEncoder(w).Encode(post)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Unauthorised access")
	}
}
