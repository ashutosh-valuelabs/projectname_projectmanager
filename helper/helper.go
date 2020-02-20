package helper

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

//FileUpload : this function returns the filename(to save in database) of the saved file or an error if it occurs
func FileUpload(r *http.Request) (string, error) {
	//ParseMultipartForm parses a request body as multipart/form-data
	r.ParseMultipartForm(32 << 20)
	//retrieve the file from form data
	file, handler, err := r.FormFile("file")
	if err != nil {
		return "", err
	}
	defer file.Close() //close the file when we finish
	//this is path which  we want to store the file
	fmt.Println(handler.Filename)
	FilePath := "/home/local/SLS/akashnidhi.p/image/" + handler.Filename
	f, err := os.OpenFile(FilePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return "", err
	}
	defer f.Close()
	io.Copy(f, file)
	return FilePath, nil
}
