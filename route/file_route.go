package route

import (
	"github.com/julienschmidt/httprouter"
	"github.com/memgo_server/handler"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const BaseUploadPath = "D:/tmp"

func init() {
	router := HttpRouter()
	router.POST("/file/upload", fileupload)
	router.POST("/excel/save", saveCourseFromExcel)
	router.POST("/excel/course", saveCourseFromExcel2)
	router.POST("/excel/course/plus", saveCourseFromExcelPlus)
	router.POST("/excel/courseTrain", saveCourTrainFromExcel2)
}

func fileupload(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	//从表单中读取文件
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		log.Panicln(err)
		_, _ = io.WriteString(w, "Read file error")
		return
	}
	defer file.Close()
	log.Println("filename:" + fileHeader.Filename)
	//创建文件
	newFile, err := os.Create(BaseUploadPath + "/" + fileHeader.Filename)
	if err != nil {
		_, _ = io.WriteString(w, "Create file error")
		return
	}
	//defer 结束时关闭文件
	defer newFile.Close()

	//将文件写到本地
	_, err = io.Copy(newFile, file)
	if err != nil {
		_, _ = io.WriteString(w, "Write file error")
		return
	}
	io.WriteString(w, "Upload success")
}

func saveCourseFromExcel(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	//创建文件
	fileType := r.Header.Get("Content-Type")
	fileType = strings.Split(fileType, "/")[1]
	newFile, err := os.Create(BaseUploadPath + "/" + strconv.FormatInt(time.Now().Unix(), 10) + "." + fileType)
	if err != nil {
		_, _ = io.WriteString(w, err.Error())
		return
	}
	//defer 结束时关闭文件
	defer newFile.Close()
	//将文件写到本地
	_, err = io.Copy(newFile, r.Body)
	if err != nil {
		_, _ = io.WriteString(w, "Write file error")
		return
	}
	io.WriteString(w, "Upload success")
}

func saveCourseFromExcel2(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	handler.SaveCourseFromExcel(r.Body)
	//handler.SavePictureFromExcel()
	r.Body.Close()
	io.WriteString(w, "Upload success")
}

func saveCourseFromExcelPlus(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	handler.SaveCourseFromExcelPlus(r.Body)
	//handler.SavePictureFromExcel()
	r.Body.Close()
	io.WriteString(w, "Upload success")
}

func saveCourTrainFromExcel2(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	handler.SaveCourTrainFromExcel(r.Body)
	//handler.SavePictureFromExcel()
	r.Body.Close()
	io.WriteString(w, "Upload success")
}
