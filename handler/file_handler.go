package handler

import (
	"bytes"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/memgo_server/database"
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
	name         string
	desc         string
	class        string
	img          string
	lecturer     int
	rawpPrice    float64
	salePrice    float64
	videoCnt     int
	period       string
	audience     string
	target       string
	recommend    int
	parts        Parts
	discountType int
	duration     int
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
	duration int
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
			db := database.Db
			err := db.QueryRow("select id from class where name = ? and type = 0", row[2]).Scan(&c.current().class)
			if err != nil {
				log.Panic(err)
			}
			//c.current().class = row[2]
		}
		file, raw, err := f.GetPicture("Sheet1", "D"+strconv.Itoa(i+1))
		if err != nil {
			log.Panicln(err)
		}
		if len(raw) != 0 && len(file) != 0 {
			p := SendPostRequest2(`http://localhost:8443/courseTrain/imgUpload`, file, bytes.NewReader(raw))
			//ioutil.WriteFile(filepath.Join("D:/tmp", file), raw, 0644)
			c.current().img = string(p)
		}
		if len(row[4]) != 0 {
			db := database.Db
			err := db.QueryRow("select user_id from user where real_name = ? limit 1;", row[4]).Scan(&c.current().lecturer)
			if err != nil {
				log.Panic(err)
			}
			//c.current().lecturer = row[4]
		}
		if len(row[5]) != 0 {
			c.current().rawpPrice, err = strconv.ParseFloat(row[5], 32)
			if err != nil {
				log.Panic(err)
			}
		}
		if len(row[6]) != 0 {
			c.current().salePrice, err = strconv.ParseFloat(row[6], 32)
			if err != nil {
				log.Panic(err)
			}
		}
		if len(row[7]) != 0 {
			c.current().recommend, err = strconv.Atoi(row[7])
			if err != nil {
				log.Panic(err)
			}
		}
		if len(row[8]) != 0 {
			c.current().discountType, err = strconv.Atoi(row[8])
			if err != nil {
				log.Panic(err)
			}
		}
		if len(row[9]) != 0 {
			c.current().audience = row[9]
		}
		if len(row[10]) != 0 {
			c.current().target = row[10]
		}
		c.current().parts.addName(row[11])
		if len(row[12]) != 0 {
			c.currentPart().desc = row[12]
		}
		c.currentPart().videos.addName(row[13])
		if len(row[14]) != 0 {
			c.currentVideo().duration, err = strconv.Atoi(row[14])
			if err != nil {
				log.Panic(err)
			}
		}
		if len(row[15]) != 0 {
			c.currentVideo().url = row[15]
		}
	}
	SaveCourse(c)
}

func SaveCourse(c *Courses) {
	db := database.Db
	tx, err := db.Begin()
	if err != nil {
		log.Panic(err)
	}
	for _, a := range c.courses {
		r, err := tx.Exec(`insert into course(name, description, class_id,image_path, coin_cnt, lecturer, is_on_line,
                   raw_price, sale_price, creator, video_cnt, recommend, discount_type,
                   audience, target, duration)
				values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, a.name, a.desc, a.class, a.img, 0, a.lecturer, 1, a.rawpPrice, a.salePrice, 1, a.videoCnt, a.recommend, a.discountType, a.audience, a.target, a.duration)
		if err != nil {
			log.Panic(err)
		}
		cid, err := r.LastInsertId()
		if err != nil {
			log.Panic(err)
		}
		for _, p := range a.parts.parts {
			r, err := tx.Exec("insert into part(name, description, course_id) values (?,?,?) ", p.name, p.desc, cid)
			if err != nil {
				log.Panic(err)
			}
			pid, err := r.LastInsertId()
			if err != nil {
				log.Panic(err)
			}
			for _, v := range p.videos.videos {
				_, err := tx.Exec("insert into video(pid, name, type, link, description, duration) values (?,?,?,?,?,?)", pid, v.name, 2, v.url, "", v.duration)
				if err != nil {
					log.Panic(err)
				}
				a.duration += v.duration
			}
			//a.videoCnt += len(p.videos.videos)
		}
		// todo 更新总时长和课时
		_, err = tx.Exec("update course set video_cnt = ?, duration = ? where id = ?", a.videoCnt, a.duration, cid)
		if err != nil {
			log.Panic(err)
		}
	}
	tx.Commit()
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
