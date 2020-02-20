package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	database "projectname_projectmanager/driver"
	model "projectname_projectmanager/model"
)

// GetOpenPositionByAging : get the open position of projects according to their age
func (C *Commander) GetOpenPositionByAging(writer http.ResponseWriter, request *http.Request) {
	db := database.DbConn()
	defer func() {
		err := db.Close()
		if err != nil {
			panic(err.Error())
		}
	}()
	if strings.Contains(Role, "program manager") == true {
		GetProjectName, err := db.Query("SELECT  sub_project_name from sub_project WHERE project_id in (SELECT id FROM project WHERE program_manager_id in (SELECT id from program_manager where program_manager_email = ?))", UserName) //getting all the project name
		if err != nil {
			fmt.Println("error in running query")
			log.Fatal(err)
		}
		defer func() {
			err := GetProjectName.Close()
			if err != nil {
				panic(err.Error())
			}
		}()
		// array of the structure
		var openPositionAging model.OpenPositionAging
		var OpenPositionsData []model.Position
		// instance of the structure
		var OpenPositionData model.Position
		// array of all the project names
		var Names []string
		// instance of the project name
		var Name string
		var Str string
		var Position int
		var Positions []int
		var aging int
		var Aging []int
		for GetProjectName.Next() {
			GetProjectName.Scan(&Name)
			Names = append(Names, Name)

		}
		for i := 0; i < len(Names); i++ {
			N := Names[i]
			//fmt.Println(Names[i], len(Names))
			GetPosition, err := db.Query("SELECT created_at, position FROM open_positions WHERE sub_project_manager_id in (select id from sub_project_manager WHERE sub_project_id IN (SELECT id FROM sub_project WHERE project_id IN (SELECT id FROM project WHERE project_name = ? AND program_manager_id IN (SELECT id FROM program_manager WHERE program_manager_email = ?))))", N, UserName)
			if err != nil {
				fmt.Println("error in running query")
				log.Fatal(err)
			}
			defer func() {
				err := GetPosition.Close()
				if err != nil {
					panic(err.Error())
				}
			}()
			t1 := time.Now()
			t := t1.Format("2006-01-02")
			//fmt.Println(t)
			for GetPosition.Next() {
				GetPosition.Scan(&Str, &Position)
				//fmt.Println(Str, Position)
				DataDiff, err := db.Query("SELECT DATEDIFF(?, ?)", t, Str)
				if err != nil {
					fmt.Println("error in running query")
					log.Fatal(err)
				}
				defer func() {
					err := DataDiff.Close()
					if err != nil {
						panic(err.Error())
					}
				}()
				for DataDiff.Next() {
					DataDiff.Scan(&aging)
				}
				//fmt.Println(aging)
				Aging = append(Aging, aging)
				Positions = append(Positions, Position)

			}
			count1 := 0
			count2 := 0
			count3 := 0
			count4 := 0
			count5 := 0
			count6 := 0
			for j := 0; j < len(Aging); j++ {

				//fmt.Println(Aging[j], Positions[j])

				if Aging[j] < 15 {
					count1 = count1 + Positions[j]

				} else if Aging[j] > 15 && Aging[j] < 30 {
					count2 = count2 + Positions[j]

				} else if Aging[j] > 30 && Aging[j] < 60 {
					count3 = count3 + Positions[j]

				} else if Aging[j] > 60 && Aging[j] < 90 {
					count4 = count4 + Positions[j]

				} else if Aging[j] > 90 && Aging[j] < 120 {
					count5 = count5 + Positions[j]

				} else {
					count6 = count6 + Positions[j]

				}
			}
			Aging = nil
			Positions = nil
			OpenPositionData.ProjectName = N
			OpenPositionData.Between0to15 = count1
			OpenPositionData.Between15to30 = count2
			OpenPositionData.Between30to60 = count3
			OpenPositionData.Between60to90 = count4
			OpenPositionData.Between90to120 = count5
			OpenPositionData.Greaterthen120 = count6
			OpenPositionsData = append(OpenPositionsData, OpenPositionData)
		}
		openPositionAging.Data = OpenPositionsData
		setupResponse(&writer, request)
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		json.NewEncoder(writer).Encode(openPositionAging)
	} else {
		writer.WriteHeader(http.StatusNotFound)
	}
}
