package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	database "projectname_projectmanager/driver"
	models "projectname_projectmanager/model"

	"github.com/gorilla/mux"
)

// ProjectUpdatesGetData : To get the project updates
func (C *Commander) ProjectUpdatesGetData(writer http.ResponseWriter, response *http.Request) {
	db := database.DbConn()
	defer db.Close()
	writer.Header().Set("Content-Type", "application/json")

	if (strings.Contains(Role, "program manager")) == true { //Converting the role from database to lower case for comparing
		var updates []models.Updates
		var data models.UpdatesData

		result, error := db.Query("call getallprojectupdates_Program(?)", UserName)
		//catch(error)
		WriteLogFile(error)
		defer result.Close()
		for result.Next() {
			var update models.Updates
			error := result.Scan(&update.ID, &update.Project, &update.Manager, &update.Ups, &update.Downs, &update.ProjectUpdates, &update.GeneralUpdates, &update.Challenges, &update.NeedHelp, &update.ClientVisits, &update.TeamSize, &update.OpenPositions, &update.HighPerformer, &update.LowPerformer)
			//catch(error)
			WriteLogFile(error)
			updates = append(updates, update)
			data.UData = updates
		}
		json.NewEncoder(writer).Encode(data)

	} else if (strings.Contains(Role, "project manager")) == true {

		HeadCount()
		OpenPosition()
		Insertion()
		var updates []models.Updates
		var update models.Updates

		var data models.UpdatesData
		json.NewDecoder(response.Body).Decode(&update)

		result, error := db.Query("call getallprojectupdates_Project(?)", UserName)
		//catch(error)
		WriteLogFile(error)
		defer result.Close()
		for result.Next() {
			var update models.Updates
			error := result.Scan(&update.ID, &update.Project, &update.Ups, &update.Downs, &update.ProjectUpdates, &update.GeneralUpdates, &update.Challenges, &update.NeedHelp, &update.ClientVisits, &update.TeamSize, &update.OpenPositions, &update.HighPerformer, &update.LowPerformer)
			//catch(error)
			WriteLogFile(error)
			updates = append(updates, update)
			data.UData = updates
		}
		json.NewEncoder(writer).Encode(data)

	} else {
		respondwithJSON(writer, http.StatusNotFound, map[string]string{"Message": "Not Found"})
	}
}

// ProjectUpdatesCreateData : To create project updates
// func (C *Commander) ProjectUpdatesCreateData(writer http.ResponseWriter, response *http.Request) {
// 	db := database.Dbconn()
// 	defer db.Close()
// 	writer.Header().Set("Content-Type", "application/json")

// 	if (strings.Contains(Role, "project manager")) == true {
// 		var update models.Updates
// 		var managerID int
// 		json.NewDecoder(response.Body).Decode(&update)
// 		result, error := db.Query("call projectupdate_ManagerID(?,?)", update.Project, UserName)
// 		catch(error)
// 		defer result.Close()
// 		if result.Next() != false {
// 			error := result.Scan(&managerID)
// 			catch(error)
// 		}
// 		fmt.Println(managerID)
// 		if managerID != 0 {
// 			var id int
// 			net, error := db.Query("SELECT net FROM head_count WHERE sub_project_manager_id = ?", managerID) // Getting Team Size from head count table
// 			InternalServerError(writer, error)
// 			defer net.Close()
// 			if net.Next() != false {
// 				error := net.Scan(&update.TeamSize)
// 				catch(error)
// 			}

// 			sum, error := db.Query("SELECT sum(position) FROM open_positions WHERE sub_project_manager_id = ?", managerID) // Getting open position from open position table
// 			catch(error)
// 			defer sum.Close()
// 			if sum.Next() != false {
// 				error := sum.Scan(&update.OpenPositions)
// 				catch(error)
// 			}
// 			puid, error := db.Query("SELECT id FROM project_updates WHERE sub_project_manager_id = ? AND is_active = 1", managerID)
// 			catch(error)
// 			defer puid.Close()
// 			if puid.Next() != false {
// 				error := puid.Scan(&id)
// 				catch(error)
// 			}

// 			if id == 0 {
// 				insert, error := db.Prepare("INSERT INTO project_updates(ups, downs, project_updates, general_updates, challenges, need_help, client_visits, team_size, open_positions, high_performer, low_performer, sub_project_manager_id, created_at, updated_at ) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,now(),now())")
// 				catch(error)
// 				_, error = insert.Exec(update.Ups, update.Downs, update.ProjectUpdates, update.GeneralUpdates, update.Challenges, update.NeedHelp, update.ClientVisits, update.TeamSize, update.OpenPositions, update.HighPerformer, update.LowPerformer, managerID)
// 				catch(error)
// 				respondwithJSON(writer, http.StatusCreated, map[string]string{"Message": "Inserted Successfully"})
// 			} else {
// 				respondwithJSON(writer, http.StatusConflict, map[string]string{"Message": "Duplicates cannot be Created"})
// 			}
// 		} else {
// 			respondwithJSON(writer, http.StatusForbidden, map[string]string{"Message": "Project not under you"})
// 		}

// 	} else {
// 		respondwithJSON(writer, http.StatusUnauthorized, map[string]string{"Message": "Not Authorized"})
// 	}
// }

// ProjectUpdatesGetDataID : to get the data according to search
func (C *Commander) ProjectUpdatesGetDataID(writer http.ResponseWriter, response *http.Request) {
	db := database.DbConn()
	defer db.Close()

	writer.Header().Set("Content-Type", "application/json")
	if (strings.Contains(Role, "program manager")) == true {
		var updates []models.Updates
		var data models.UpdatesData
		parameter := mux.Vars(response)
		key := parameter["id"]
		var search string = key + "%"
		result, error := db.Query("call getallprojectupdates_ProgramID(?,?)", UserName, search)
		catch(error)
		defer result.Close()
		for result.Next() {
			var update models.Updates

			error := result.Scan(&update.ID, &update.Project, &update.Ups, &update.Downs, &update.ProjectUpdates, &update.GeneralUpdates, &update.Challenges, &update.NeedHelp, &update.ClientVisits, &update.TeamSize, &update.OpenPositions, &update.HighPerformer, &update.LowPerformer)
			catch(error)
			updates = append(updates, update)
			data.UData = updates
		}
		json.NewEncoder(writer).Encode(data)

	} else if (strings.Contains(Role, "project manager")) == true {
		var updates []models.Updates
		var data models.UpdatesData
		parameter := mux.Vars(response)
		key := parameter["id"]
		var search string = key + "%"
		result, error := db.Query("call getallprojectupdates_ProjectID(?,?) ", UserName, search)
		catch(error)
		defer result.Close()
		for result.Next() {
			var update models.Updates
			error := result.Scan(&update.ID, &update.Project, &update.Ups, &update.Downs, &update.ProjectUpdates, &update.GeneralUpdates, &update.Challenges, &update.NeedHelp, &update.ClientVisits, &update.TeamSize, &update.OpenPositions, &update.HighPerformer, &update.LowPerformer)
			catch(error)
			updates = append(updates, update)
			data.UData = updates
		}
		json.NewEncoder(writer).Encode(data)

	} else {
		respondwithJSON(writer, http.StatusNotFound, map[string]string{"Message": "Not Found"})
	}
}

// ProjectUpdatesUpdateData : To update a particular prospect
func (C *Commander) ProjectUpdatesUpdateData(writer http.ResponseWriter, response *http.Request) {
	db := database.DbConn()
	defer db.Close()
	writer.Header().Set("Content-Type", "application/json")
	if (strings.Contains(Role, "project manager")) == true {
		var update models.Updates
		var managerID int
		json.NewDecoder(response.Body).Decode(&update)
		result, error := db.Query("call projectupdate_ManagerID(?,?)", update.Project, UserName)
		catch(error)
		defer result.Close()
		if result.Next() != false {
			error := result.Scan(&managerID)
			catch(error)
		}

		if managerID != 0 {
			var id int
			puid, error := db.Query("SELECT id FROM project_updates WHERE sub_project_manager_id = ? and is_active = 1", managerID)
			catch(error)
			defer puid.Close()
			if puid.Next() != false {
				error := puid.Scan(&id)
				catch(error)
			}
			if id != 0 {
				query, error := db.Prepare("UPDATE project_updates SET ups = ?, downs = ?, project_updates = ?, general_updates = ?, challenges = ?, need_help = ?, client_visits = ?, high_performer = ?, low_performer = ?, updated_at = ? WHERE sub_project_manager_id = ? AND is_active = 1")
				catch(error)
				_, error = query.Exec(update.Ups, update.Downs, update.ProjectUpdates, update.GeneralUpdates, update.Challenges, update.NeedHelp, update.ClientVisits, update.HighPerformer, update.LowPerformer, time.Now().Format("2006-01-02 15:04:05"), managerID)
				catch(error)
				defer query.Close()
				respondwithJSON(writer, http.StatusOK, map[string]string{"message": "Updated Successfully"})
			} else {
				respondwithJSON(writer, http.StatusConflict, map[string]string{"Message": "Data Not Found"})
			}
		} else {
			respondwithJSON(writer, http.StatusForbidden, map[string]string{"Message": "Project Not Under You"})
		}
	} else {
		respondwithJSON(writer, http.StatusUnauthorized, map[string]string{"Message": "Not Authorized"})
	}

}

// ProjectUpdatesDeleteData : To delete a particular prospect
func (C *Commander) ProjectUpdatesDeleteData(writer http.ResponseWriter, response *http.Request) {
	db := database.DbConn()
	defer db.Close()
	writer.Header().Set("Content-Type", "application/json")
	if (strings.Contains(Role, "project manager")) == true {
		var update models.Updates
		var managerID int
		json.NewDecoder(response.Body).Decode(&update)
		result, error := db.Query("SELECT sub_project_manager_id FROM project_updates WHERE id = ?", update.ID)
		catch(error)
		defer result.Close()
		if result.Next() != false {
			error := result.Scan(&managerID)
			catch(error)
		}
		if managerID != 0 {
			var email string
			manageremail, error := db.Query("SELECT project_manager_email FROM project_updates LEFT JOIN sub_project_manager ON project_updates.sub_project_manager_id = sub_project_manager.id LEFT JOIN project_manager ON sub_project_manager.manager_id = project_manager.id LEFT JOIN sub_project ON sub_project_manager.sub_project_id = sub_project.id WHERE project_manager.project_manager_email = ? AND project_updates.is_active = 1", UserName)
			catch(error)
			defer manageremail.Close()
			if manageremail.Next() != false {
				error := manageremail.Scan(&email)
				catch(error)
			}
			if email == UserName {
				stmt, error := db.Prepare("UPDATE project_updates SET is_active = 0 WHERE sub_project_manager_id = ? AND project_updates.id = ?")
				catch(error)
				_, error = stmt.Exec(managerID, update.ID)
				catch(error)
				defer stmt.Close()
				respondwithJSON(writer, http.StatusOK, map[string]string{"message": "Deleted Successfully"})
			} else {
				respondwithJSON(writer, http.StatusUnauthorized, map[string]string{"Message": "Not Authorized to Delete"})
			}
		} else {
			respondwithJSON(writer, http.StatusBadRequest, map[string]string{"Message": "Project Not Under You"})
		}
	} else {
		respondwithJSON(writer, http.StatusUnauthorized, map[string]string{"Message": "Not Authorized"})
	}
}

//HeadCount : To set the head count to default value
func HeadCount() {
	db := database.DbConn()
	var ManagerID int
	defer db.Close()
	projectManagerID, err := db.Query("select id from sub_project_manager where manager_id in (select id from project_manager where project_manager_email= ?)", UserName)
	catch(err)
	defer projectManagerID.Close()
	for projectManagerID.Next() {
		projectManagerID.Scan(&ManagerID)

		stmt2, err := db.Query("SELECT id from head_count WHERE sub_project_manager_id = ? ", ManagerID)
		catch(err)
		if stmt2.Next() == false {
			stmt, err := db.Prepare("INSERT INTO head_count(sub_project_manager_id,created_at, updated_at) VALUES(?, now(), now())")
			catch(err)
			defer stmt.Close()
			_, err = stmt.Exec(ManagerID)
			catch(err)
		}
		stmt2.Close()

	}
}

//OpenPosition :
func OpenPosition() {
	db := database.DbConn()
	var ManagerID int
	defer db.Close()
	projectManagerID, err := db.Query("select id from sub_project_manager where manager_id in (select id from project_manager where project_manager_email= ?)", UserName)
	catch(err)
	defer projectManagerID.Close()
	for projectManagerID.Next() {
		projectManagerID.Scan(&ManagerID)

		stmt2, err := db.Query("SELECT id from open_positions WHERE sub_project_manager_id = ? ", ManagerID)
		catch(err)
		if stmt2.Next() == false {
			stmt, err := db.Prepare("INSERT INTO open_positions(sub_project_manager_id,created_at, updated_at) VALUES(?, now(), now())")
			catch(err)
			defer stmt.Close()
			_, err = stmt.Exec(ManagerID)
			catch(err)
		}
		stmt2.Close()

	}
}

//Insertion : To set the default value
func Insertion() {
	db := database.DbConn()
	var update models.Updates
	var ManagerID int
	defer db.Close()
	projectManagerID, err := db.Query("select id from sub_project_manager where manager_id in (select id from project_manager where project_manager_email= ?)", UserName)
	//catch(err)
	WriteLogFile(err)
	defer projectManagerID.Close()
	for projectManagerID.Next() {
		projectManagerID.Scan(&ManagerID)

		net, error := db.Query("SELECT net FROM head_count WHERE sub_project_manager_id = ?", ManagerID) // Getting Team Size from head count table

		defer net.Close()
		if net.Next() != false {
			error := net.Scan(&update.TeamSize)
			//catch(error)
			WriteLogFile(error)
		}

		sum, error := db.Query("SELECT sum(position) FROM open_positions WHERE sub_project_manager_id = ?", ManagerID) // Getting open position from open position table
		//catch(error)
		WriteLogFile(error)
		defer sum.Close()
		if sum.Next() != false {
			error := sum.Scan(&update.OpenPositions)
			//catch(error)
			WriteLogFile(error)
		}
		fmt.Println(update.OpenPositions)
		stmt2, err := db.Query("SELECT id from project_updates WHERE sub_project_manager_id = ? ", ManagerID)
		//catch(err)
		WriteLogFile(err)
		if stmt2.Next() == false {

			stmt, err := db.Prepare("INSERT INTO project_updates(sub_project_manager_id,team_size,open_positions,created_at, updated_at) VALUES(?,?,?,now(), now())")
			//catch(err)
			WriteLogFile(err)
			defer stmt.Close()
			_, err = stmt.Exec(ManagerID, update.TeamSize, update.OpenPositions)
			//catch(err)
			WriteLogFile(error)

		}
		stmt2.Close()
	}
}
