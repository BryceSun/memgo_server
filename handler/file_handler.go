package handler

import (
	"bytes"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

type course struct {
	name      string
	desc      string
	class     string
	img       string
	lecturer  string
	rawpPrice string
	salePrice string
	period    string
	audience  string
	target    string
	parts     Parts
}

type Courses struct {
	courses []course
	index   int
}

type part struct {
	name   string
	desc   string
	videos Videos
}

type Parts struct {
	parts []part
	index int
}

type video struct {
	name     string
	url      string
	duration string
}

type Videos struct {
	videos []video
	index  int
}

func NewCourses() *Courses {
	return &Courses{
		courses: make([]course, 0),
		index:   -1,
	}
}

func NewParts() *Parts {
	return &Parts{
		parts: make([]part, 0),
		index: -1,
	}
}

func NewVideos() *Videos {
	return &Videos{
		videos: make([]video, 0),
		index:  -1,
	}
}

func (c *Courses) current() *course {
	return &c.courses[c.index]
}

func (c *Courses) currentPart() *part {
	return c.current().parts.current()
}

func (c *Courses) currentVideo() *video {
	return c.current().parts.current().videos.current()
}

func (p *Parts) current() *part {
	return &p.parts[p.index]
}

func (v *Videos) current() *video {
	return &v.videos[v.index]
}
func (c *Courses) addName(name string) {

	if len(name) == 0 {
		return
	}
	for _, v := range c.courses {
		if v.name == name {
			return
		}
	}
	cc := course{
		parts: *NewParts(),
		name:  name,
	}
	c.courses = append(c.courses, cc)
	c.index++
}

func (p *Parts) addName(name string) {
	if len(name) == 0 {
		return
	}
	for _, v := range p.parts {
		if v.name == name {
			return
		}
	}
	pp := part{
		videos: *NewVideos(),
	}
	pp.name = name
	p.parts = append(p.parts, pp)
	p.index++
}

func (v *Videos) addName(name string) {
	if len(name) == 0 {
		return
	}
	for _, v := range v.videos {
		if v.name == name {
			return
		}
	}
	vv := video{}
	vv.name = name
	v.videos = append(v.videos, vv)
	v.index++
}

func SendPostRequest(url string, filename string, filetype string) []byte {
	file, err := os.Open(filename)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(filetype, filepath.Base(file.Name()))

	if err != nil {
		log.Fatal(err)
	}

	io.Copy(part, file)
	writer.Close()
	request, err := http.NewRequest("POST", url, body)

	if err != nil {
		log.Fatal(err)
	}

	request.Header.Add("Content-Type", writer.FormDataContentType())
	client := &http.Client{}

	response, err := client.Do(request)

	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	content, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Fatal(err)
	}

	return content
}

func SendPostRequest2(url, fileName string, reader io.Reader) []byte {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", fileName)

	if err != nil {
		log.Fatal(err)
	}

	io.Copy(part, reader)
	writer.Close()
	request, err := http.NewRequest("POST", url, body)

	if err != nil {
		log.Fatal(err)
	}

	request.Header.Add("Content-Type", writer.FormDataContentType())
	client := &http.Client{}

	response, err := client.Do(request)

	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	content, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Fatal(err)
	}

	return content
}

func SaveCourseFromExcel(r io.Reader) {
	f, err := excelize.OpenReader(r)
	if err != nil {
		log.Panicln(err)
	}
	c := NewCourses()
	rows, err := f.GetRows("Sheet1")
	for i, row := range rows {
		if len(row) < 13 || i == 0 {
			continue
		}
		c.addName(row[0])
		if len(row[1]) != 0 {
			c.current().desc = row[1]
		}
		if len(row[2]) != 0 {
			c.current().class = row[2]
		}
		file, raw, err := f.GetPicture("Sheet1", "D"+strconv.Itoa(i+1))
		if err != nil {
			log.Panicln(err)
		}
		p := SendPostRequest2(`http://localhost:8443/courseTrain/imgUpload`, file, bytes.NewReader(raw))
		//ioutil.WriteFile(filepath.Join("D:/tmp", file), raw, 0644)
		if len(row[3]) != 0 {
			c.current().img = string(p)
		}
		if len(row[4]) != 0 {
			c.current().lecturer = row[4]
		}
		if len(row[5]) != 0 {
			c.current().rawpPrice = row[5]
		}
		if len(row[6]) != 0 {
			c.current().salePrice = row[6]
		}
		if len(row[7]) != 0 {
			c.current().period = row[7]
		}
		c.current().parts.addName(row[8])
		if len(row[9]) != 0 {
			c.currentPart().desc = row[9]
		}
		c.currentPart().videos.addName(row[10])
		//if len(row[11]) != 0 {
		//c.currentVideo().url=row[11]

		//}
		if len(row[12]) != 0 {
			c.currentVideo().duration = row[12]
		}
	}
}

func SavePictureFromExcel() {
	f, err := excelize.OpenFile("C:\\Users\\wjw\\Desktop\\test.xlsx")
	if err != nil {
		log.Panicln(err)
	}
	_, err = f.GetRows("Sheet1")

	file, raw, err := f.GetPicture("Sheet1", "A1")
	if err != nil {
		log.Panicln(err)
	}
	ioutil.WriteFile(filepath.Join("D:/tmp", file), raw, 0644)
}
