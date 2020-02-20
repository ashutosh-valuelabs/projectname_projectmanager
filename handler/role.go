package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	database "projectname_projectmanager/driver"
	model "projectname_projectmanager/model"
)

func (c *Commander) GetRole(writer http.ResponseWriter, request *http.Request) {
	var Project model.UserRole
	Email := request.URL.Query()["role"]
	fmt.Println(Email[0])

	//json.NewDecoder(request.Body).Decode(&Data)
	db := database.DbConn()
	defer db.Close()
	getProjectManagerID, err := db.Query("SELECT id from project_manager where project_manager_email = ? ", Email[0])
	defer getProjectManagerID.Close()
	if err != nil {
		panic(err)
	}
	if getProjectManagerID.Next() == true {
		role := "project manager"
		Project.Role = append(Project.Role, role)
	}
	getProgramManagerID, err := db.Query("SELECT id from program_manager where program_manager_email = ? ", Email[0])
	defer getProgramManagerID.Close()
	if err != nil {
		panic(err)
	}
	if getProgramManagerID.Next() == true {
		role := "program manager"
		Project.Role = append(Project.Role, role)
	}
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(Project)

}
