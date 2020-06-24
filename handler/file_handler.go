package handler

import (
	"bytes"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/memgo_server/database"
	"github.com/shopspring/decimal"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"unsafe"
)

func init() {
	genSetFieldFuncs(&course{})
	genSetFieldFuncs(&part{})
	genSetFieldFuncs(&video{})

}

type course struct {
	name         string `cellNum:"0"`
	desc         string `cellNum:"1"`
	class        string `cellNum:"2"`
	img          string `cellNum:"3"`
	lecturer     int    `cellNum:"4"`
	rawpPrice    int    `cellNum:"5"`
	salePrice    int    `cellNum:"6"`
	videoCnt     int
	period       string
	recommend    int    `cellNum:"7"`
	discountType int    `cellNum:"8"`
	audience     string `cellNum:"9"`
	target       string `cellNum:"10"`
	parts        Parts
	duration     int
}

type receiveField func(interface{}, string)

var fieldFuncs = make(map[string]receiveField)

func genSetFieldFuncs(c interface{}) {
	v := reflect.ValueOf(c).Elem() // the struct variable
	for i := 0; i < v.NumField(); i++ {
		fieldInfo := v.Type().Field(i) // a reflect.StructField
		tag := fieldInfo.Tag           // a reflect.StructTag
		cellNum := tag.Get("cellNum")
		fieldFuncs[cellNum] = genSetFieldFunc(fieldInfo)
	}
}

func genSetFieldFunc(field reflect.StructField) receiveField {
	if field.Type.Kind() == reflect.Int {
		return func(c interface{}, v string) {
			var vv int
			var err error
			if strings.Contains(v, ".") {
				d, err := decimal.NewFromString(v)
				if err != nil {
					log.Panic(err)
				}
				f, _ := d.Float64()
				vv = int(f * 100)
			} else {
				vv, err = strconv.Atoi(v)
				if err != nil {
					log.Panicln(err)
				}
			}
			rv := reflect.ValueOf(c).Elem()
			rv.FieldByName(field.Name).Set(reflect.ValueOf(vv))
		}
	}
	if field.Type.Kind() == reflect.String {
		return func(c interface{}, v string) {
			rv := reflect.ValueOf(c).Elem()
			rv.FieldByName(field.Name).Set(reflect.ValueOf(v))
		}
	}
	return nil

}

type Courses struct {
	courses []course
	index   int
}

type part struct {
	name   string `cellNum:"11"`
	desc   string `cellNum:"12"`
	videos Videos
}

type Parts struct {
	parts []part
	index int
}

type video struct {
	name     string `cellNum:"13"`
	duration int    `cellNum:"14"`
	url      string `cellNum:"15"`
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
		name:   name,
	}
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
		if len(row) < 16 || i == 0 {
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
		}
		file, raw, err := f.GetPicture("Sheet1", "D"+strconv.Itoa(i+1))
		if err != nil {
			log.Panicln(err)
		}
		if len(raw) != 0 && len(file) != 0 {
			p := SendPostRequest2(`http://139.199.3.134:8443/courseTrain/imgUpload`, file, bytes.NewReader(raw))
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
			d, err := decimal.NewFromString(row[5])
			if err != nil {
				log.Panic(err)
			}
			rp, _ := d.Float64()
			c.current().rawpPrice = int(rp * 100)
		}
		if len(row[6]) != 0 {
			//sp, err := strconv.ParseFloat(row[6], 32)
			d, err := decimal.NewFromString(row[6])
			if err != nil {
				log.Panic(err)
			}
			sp, _ := d.Float64()
			c.current().salePrice = int(sp * 100)
		}
		if len(row[7]) != 0 {
			c.current().recommend, err = strconv.Atoi(row[7])
			if err != nil {
				log.Panic(err)
			}
		}
		if len(row[8]) != 0 {
			switch row[8] {
			case "免费":
				c.current().discountType = 0
			case "会员免费":
				c.current().discountType = 1
			case "折扣":
				c.current().discountType = 2
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

func SaveCourseFromExcelPlus(r io.Reader) {
	f, err := excelize.OpenReader(r)
	if err != nil {
		log.Panicln(err)
	}
	c := NewCourses()
	rows, err := f.GetRows("Sheet1")
	for i, row := range rows {
		if len(row) < 16 || i == 0 {
			continue
		}
		for ci, cell := range row {
			if len(cell) == 0 {
				if ci != 3 {
					continue
				}
				file, raw, err := f.GetPicture("Sheet1", "D"+strconv.Itoa(i+1))
				if err != nil {
					log.Panicln(err)
				}
				if len(raw) != 0 && len(file) != 0 {
					p := SendPostRequest2(`http://139.199.3.134:8443/courseTrain/imgUpload`, file, bytes.NewReader(raw))
					cell = string(p)
				}

			}
			if ci == 0 {
				c.addName(cell)
				continue
			}
			if ci == 11 {
				c.current().parts.addName(cell)
				continue
			}
			if ci == 13 {
				c.currentPart().videos.addName(cell)
				continue
			}
			if ci > 13 {
				fieldFuncs[strconv.Itoa(ci)](*c.currentVideo(), cell)
				continue
			}
			if ci > 11 {
				fieldFuncs[strconv.Itoa(ci)](*c.currentPart(), cell)
				continue
			}
			fieldFuncs[strconv.Itoa(ci)](*c.currentPart(), cell)
		}
	}
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
			a.videoCnt += len(p.videos.videos)
		}
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

type courseTrain struct {
	name      string `cellNum:"0"`
	class     string `cellNum:"1"`
	img       string `cellNum:"2"`
	lecturer  int    `cellNum:"3"`
	desc      string `cellNum:"4"`
	length    int    `cellNum:"5"`
	audience  string `cellNum:"6"`
	target    string `cellNum:"7"`
	directory string `cellNum:"8"`
}

func SaveCourTrainFromExcel(r io.Reader) {
	f, err := excelize.OpenReader(r)
	if err != nil {
		log.Panicln(err)
	}
	rows, err := f.GetRows("Sheet1")
	cs := make([]*courseTrain, 0)
	var c *courseTrain
	for ri, row := range rows {
		if len(row) < 9 || ri == 0 {
			continue
		}
		c = &courseTrain{}
		for ci, cell := range row {
			if len(cell) == 0 && ci != 2 {
				continue
			}
			switch ci {
			case 0:
				c.name = cell
			case 1:
				db := database.Db
				err := db.QueryRow("select id from class where name = ? and type = 1", cell).Scan(&c.class)
				if err != nil {
					log.Panic(err)
				}
			case 2:
				file, raw, err := f.GetPicture("Sheet1", "C"+strconv.Itoa(ri+1))
				if err != nil {
					log.Panicln(err)
				}
				if len(raw) != 0 && len(file) != 0 {
					p := SendPostRequest2(`http://139.199.3.134:8443/courseTrain/imgUpload`, file, bytes.NewReader(raw))
					c.img = string(p)
				}
			case 3:
				db := database.Db
				err := db.QueryRow("select user_id from user where real_name = ? limit 1;", cell).Scan(&c.lecturer)
				if err != nil {
					log.Panic(err)
				}
			case 4:
				c.desc = cell
			case 5:
				c.length, err = strconv.Atoi(cell)
				if err != nil {
					log.Panicln(err)
				}
			case 6:
				c.audience = cell
			case 7:
				c.target = cell
			case 8:
				c.directory = cell
			}
		}
		if unsafe.Sizeof(*c) != 0 {
			cs = append(cs, c)
		}
	}
	db := database.Db
	for _, a := range cs {
		_, err := db.Exec(`INSERT INTO fdl.course_train(name,description,class_id,image_path,lecturer,length,creator,audience,target,directory)
                                  VALUES(?,?,?,?,?,?,?,?,?,?);`, a.name, a.desc, a.class, a.img, a.lecturer, a.length, 1, a.audience, a.target, a.directory)
		if err != nil {
			log.Panic(err)
		}
	}
}
