package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	database "projectname_projectmanager/driver"
	helper "projectname_projectmanager/helper"
	model "projectname_projectmanager/model"
	"strconv"
	"time"
)

// ChangeProfileImage : uploding the Profile Image to Server.
func (C *Commander) ChangeProfileImage(writer http.ResponseWriter, request *http.Request) {
	Time := time.Now()
	db := database.DbConn()
	defer db.Close()
	var user model.Profile
	//reading the user whose image we want to change from the database
	user.Name = UserName
	user.Role = Role
	//here we call the function we made to get the image and save it
	imageName, err := helper.FileUpload(request)
	if err != nil {
		http.Error(writer, "Invalid Data", http.StatusBadRequest)
		return
	}
	User, _ := db.Query("SELECT id FROM profile WHERE username = ?", UserName)
	defer User.Close()
	if User.Next() != false {
		UpdateProfile, _ := db.Query("UPDATE profile set username = ?, role = ?, image_path = ?, updated_at = ?", user.Name, user.Role, imageName, Time)
		defer UpdateProfile.Close()
	} else {
		InsertProfile, _ := db.Query("INSERT into profile(username, role, image_path, created_at, updated_at)VALUES(?, ?, ?, ?, ?)", user.Name, user.Role, imageName, Time, Time)
		defer InsertProfile.Close()
	}
	setupResponse(&writer, request)
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	json.NewEncoder(writer).Encode(user)
}

// FileServe : Serving file to client.
func (C *Commander) FileServe(writer http.ResponseWriter, request *http.Request) {
	db := database.DbConn()
	defer db.Close()
	var Filename string
	GetFilePath, _ := db.Query("SELECT image_path FROM profile WHERE username = ?", UserName)
	defer GetFilePath.Close()
	if GetFilePath.Next() != false {
		GetFilePath.Scan(&Filename)
	}
	if Filename == "" {
		//Get not set, send a 400 bad request
		http.Error(writer, "Get 'file' not specified in url.", 400)
		return
	}
	fmt.Println("Client requests: " + Filename)

	Openfile, err := os.Open(Filename)
	defer Openfile.Close()
	if err != nil {
		//File not found, send 404
		http.Error(writer, "File not found.", 404)
		return
	}

	FileHeader := make([]byte, 512)
	Openfile.Read(FileHeader)
	FileContentType := http.DetectContentType(FileHeader)
	FileStat, _ := Openfile.Stat()
	FileSize := strconv.FormatInt(FileStat.Size(), 10)
	setupResponse(&writer, request)
	writer.Header().Set("Content-Disposition", "attachment; filename="+Filename)
	writer.Header().Set("Content-Type", FileContentType)
	writer.Header().Set("Content-Length", FileSize)
	//We read 512 bytes from the file already, so we reset the offset back to 0
	Openfile.Seek(0, 0)
	//'Copy' the file to the client
	io.Copy(writer, Openfile)
	return
}
