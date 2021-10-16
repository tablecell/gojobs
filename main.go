package main

import (
	"fmt"
	"github.com/azer/crud"
	_ "github.com/mattn/go-sqlite3"
	"github.com/xujiajun/gorouter"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)
type (
	Jobs struct {
		Id          int       `sql:"auto_increment primary-key"`
		Title       string    `db:"title"`
		Salarytop   int       `db:"salarytop"`
		Salarybase  int       `db:"salarybase"`
		Description string    `db:"description"`
		Published   time.Time `db:"published"`
	}
	TplData struct {
		Jobs []*Jobs
	}
)

var (
	DB               *crud.DB
	htmlTplEngine    *template.Template
	htmlTplEngineErr error
)
func ServeFiles(w http.ResponseWriter, r *http.Request) {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	dir := wd + "/assets"
	fmt.Println(dir)
	http.StripPrefix("/assets/", http.FileServer(http.Dir(dir))).ServeHTTP(w, r)
}

func main() {

	var err error
	DB, err := crud.Connect("sqlite3", "jobs.db")
	//	err = DB.Ping()
	fmt.Println(err)

	htmlTplEngine = template.New("htmlTplEngine")

	// 模板根目录下的模板文件 一些公共文件
	_, htmlTplEngineErr = htmlTplEngine.ParseGlob("views/*.html")
	if nil != htmlTplEngineErr {
		log.Panic(htmlTplEngineErr.Error())
	}

	mux := gorouter.New()

	mux2 := mux.Group("/assets")
	mux2.GET("/{filename:[0-9a-zA-Z_.-]+}", func(w http.ResponseWriter, r *http.Request) {
		ServeFiles(w, r)
	})

	mux.GET("/", func(w http.ResponseWriter, r *http.Request) {
		jobs := []*Jobs{}

		err := DB.Read(&jobs, "SELECT  * FROM jobs WHERE id > ?", 0)
		fmt.Println(err)
		data := TplData{Jobs: jobs}

		_ = htmlTplEngine.ExecuteTemplate(w, "home.html", data)
 
	})

	mux.POST("/new", func(w http.ResponseWriter, r *http.Request) {
		// title
		// fmt. r.FormValue("title")
		base, err := strconv.Atoi(r.FormValue("salarybase"))
		top, err := strconv.Atoi(r.FormValue("salarytop"))
		fmt.Println(err)
		var book = Jobs{
			Title:       r.FormValue("title"),
			Description: r.FormValue("jd"),
			Published:   time.Now(),
			Salarybase:  base,
			Salarytop:   top,
		}

		err = DB.Create(book)
		fmt.Println(err)
 http.Redirect(w, r, "/", 302)
//		tmpl := template.Must(template.ParseFiles("./test.html"))
//		tmpl.Execute(w, struct{ Success bool }{true})
	})

	mux.GET("/new", func(w http.ResponseWriter, r *http.Request) {
		_ = htmlTplEngine.ExecuteTemplate(w, "new.html", nil)
	})


    fmt.Println("ListenAndServe:  http://localhost:8181")
	log.Fatal(http.ListenAndServe(":8181", mux))
}
