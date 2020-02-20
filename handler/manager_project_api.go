package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	database "projectname_projectmanager/driver"
	model "projectname_projectmanager/model"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

//setupResponse : To set all the CORS request
func setupResponse(writer *http.ResponseWriter, request *http.Request) {
	(*writer).Header().Set("Access-Control-Allow-Origin", "*")
	(*writer).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*writer).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

// GetProjectManagerSearchResult : send all the data with te requested manager name
func (C *Commander) GetProjectManagerSearchResult(writer http.ResponseWriter, request *http.Request) {
	var error model.Error
	db := database.DbConn()
	defer func() {
		err := db.Close()
		if err != nil {
			panic(err.Error())
		}
	}()
	if strings.Contains(Role, "program manager") == true {
		p := mux.Vars(request)
		key1 := p["id"]
		key := strings.TrimSpace(key1)
		SearchString := key + "%"
		Offset := 0
		Pages := request.URL.Query()["Pages"]
		limit, err := strconv.Atoi(request.URL.Query().Get("limit"))
		if limit != 10 && limit != 20 && limit != 50 {
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusBadRequest)
			error.Message = "Incorrect Limit Value"
			json.NewEncoder(writer).Encode(error)
			return
		}
		if err != nil {
			WriteLogFile(err)
			return
		}
		i1, err := strconv.Atoi(Pages[0])
		if err != nil {
			WriteLogFile(err)
			return
		}
		Offset = 10 * i1
		Total := 0
		SearchResults, err := db.Query("call GetManagerDetailsByManagerName(?, ?, ?, ?)", SearchString, UserName, Offset, limit)
		if err != nil {
			WriteLogFile(err)
			return
		}
		defer func() {
			err := SearchResults.Close()
			if err != nil {
				panic(err.Error())
			}
		}()
		var ManagerDetailData model.Project
		var ManagerDetailsData []model.Project
		for SearchResults.Next() {
			SearchResults.Scan(&ManagerDetailData.ProjectName, &ManagerDetailData.SubProjectName, &ManagerDetailData.ManagerName, &ManagerDetailData.ManagerEmailID, &ManagerDetailData.Id)
			ManagerDetailsData = append(ManagerDetailsData, ManagerDetailData)
			Total++
		}

		var PaginationFormat model.Pagination
		PaginationFormat.TotalData = Total
		PaginationFormat.Limit = limit
		PaginationFormat.Data = ManagerDetailsData
		x1 := Total / limit
		x := Total % limit
		if x == 0 {
			PaginationFormat.TotalPages = x1
		} else {
			PaginationFormat.TotalPages = x1 + 1
		}
		x, _ = strconv.Atoi(Pages[0])
		if PaginationFormat.TotalPages != 0 {
			x1 = x + 1
		}
		PaginationFormat.Page = x1
		setupResponse(&writer, request)
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		json.NewEncoder(writer).Encode(PaginationFormat)
	} else {
		writer.WriteHeader(http.StatusNotFound)
	}
}

// GetProjectName : send all the project name
func (C *Commander) GetProjectName(writer http.ResponseWriter, request *http.Request) {
	db := database.DbConn()
	defer func() {
		err := db.Close()
		if err != nil {
			panic(err.Error())
		}
	}()
	if strings.Contains(Role, "program manager") == true {
		GetProjectName, err := db.Query("SELECT  sub_project_name FROM sub_project WHERE project_id IN (SELECT id  FROM project WHERE program_manager_id in (SELECT id FROM program_manager WHERE program_manager_email = ?))", UserName)
		if err != nil {
			WriteLogFile(err)
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
		defer func() {
			err := GetProjectName.Close()
			if err != nil {
				panic(err.Error())
			}
		}()
		var ProjectNames []string
		var ProjectName string
		for GetProjectName.Next() {
			err := GetProjectName.Scan(&ProjectName)
			if err != nil {
				WriteLogFile(err)
				writer.WriteHeader(http.StatusBadRequest)
				return
			}
			ProjectNames = append(ProjectNames, ProjectName)
		}
		setupResponse(&writer, request)
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		json.NewEncoder(writer).Encode(ProjectNames)
	}
	if strings.Contains(Role, "project manager") == true {
		GetProjectName, err := db.Query("SELECT  sub_project.sub_project_name FROM sub_project_manager JOIN sub_project ON sub_project_manager.sub_project_id = sub_project.id Join project_manager ON sub_project_manager.manager_id = project_manager.id WHERE project_manager_email = ? ", UserName)
		if err != nil {
			WriteLogFile(err)
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
		defer func() {
			err := GetProjectName.Close()
			if err != nil {
				panic(err.Error())
			}
		}()
		var ProjectNames []string
		var ProjectName string
		for GetProjectName.Next() {
			err := GetProjectName.Scan(&ProjectName)
			if err != nil {
				WriteLogFile(err)
				writer.WriteHeader(http.StatusBadRequest)
				return
			}
			ProjectNames = append(ProjectNames, ProjectName)
		}
		setupResponse(&writer, request)
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		json.NewEncoder(writer).Encode(ProjectNames)

	}
}

// GetData : get all data
func (C *Commander) GetData(writer http.ResponseWriter, request *http.Request) {
	var error model.Error
	db := database.DbConn()
	defer func() {
		err := db.Close()
		if err != nil {
			panic(err.Error())
		}
	}()
	if strings.Contains(Role, "program manager") == true {
		var Offset int
		Pages := request.URL.Query()["Pages"]
		fmt.Println(Pages)
		if Pages[0] != "" {
			limit, err := strconv.Atoi(request.URL.Query().Get("limit"))
			if limit != 10 && limit != 20 && limit != 50 {
				writer.Header().Set("Content-Type", "application/json")
				writer.WriteHeader(http.StatusBadRequest)
				error.Message = "Incorrect Limit Value"
				json.NewEncoder(writer).Encode(error)
				return
			}
			i1, _ := strconv.Atoi(Pages[0])
			fmt.Println(i1)
			Offset = 10 * i1
			count, _ := db.Query("SELECT COUNT(Id) FROM sub_project_manager WHERE sub_project_id in (SELECT id FROM sub_project WHERE project_id in (SELECT id FROM project WHERE program_manager_id in (SELECT id FROM program_manager where program_manager_email = ?)))", UserName)
			defer func() {
				err := count.Close()
				if err != nil {
					panic(err.Error())
				}
			}()
			GetManagerDetails, err := db.Query("call GetAllManagerDetailsData(?, ?, ?)", UserName, Offset, limit)
			if err != nil {
				WriteLogFile(err)
				return
			}
			defer func() {
				err := GetManagerDetails.Close()
				if err != nil {
					panic(err.Error())
				}
			}()
			var Total int
			var ManagerDetailData model.Project
			var ManagerDetailsData []model.Project
			for GetManagerDetails.Next() {
				GetManagerDetails.Scan(&ManagerDetailData.ProjectName, &ManagerDetailData.SubProjectName, &ManagerDetailData.ManagerName, &ManagerDetailData.ManagerEmailID, &ManagerDetailData.Id)
				ManagerDetailsData = append(ManagerDetailsData, ManagerDetailData)
			}
			if count.Next() != false {
				count.Scan(&Total)
			} else {
				Total = 0
			}
			var PaginationFormat model.Pagination
			PaginationFormat.TotalData = Total
			PaginationFormat.Limit = limit
			PaginationFormat.Data = ManagerDetailsData
			x1 := Total / limit
			x := Total % limit
			if x == 0 {
				PaginationFormat.TotalPages = x1
			} else {
				PaginationFormat.TotalPages = x1 + 1
			}
			x, _ = strconv.Atoi(Pages[0])
			if PaginationFormat.TotalPages != 0 {
				x1 = x + 1
			}
			PaginationFormat.Page = x1
			setupResponse(&writer, request)
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusOK)
			json.NewEncoder(writer).Encode(PaginationFormat)
		} else {
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusBadRequest)
			error.Message = "Incorrect Page Value"
			json.NewEncoder(writer).Encode(error)
			return

		}
	} else {
		writer.WriteHeader(http.StatusNotFound)
	}
}
