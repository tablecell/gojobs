package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"

	"html/template"
	"log"
	"net/http"
	"time"
)

type (
	Jobs struct {
		Id         int    `sql:"auto_increment primary-key"`
		Title      string `db:"title"`
		Salarytop  int    `db:"salarytop"`
		Salarybase int    `db:"salarybase"`
		Mail       string `db:"mail"`
		Gender     int    `db:"gender"`
		Is996      int    `db:"is996"`
		Level      int    `db:"level"`
		Description string    `db:"description"`
		Remote      string    `db:"remote"`
		Published   time.Time `db:"published"`
	}
	TplData struct {
		Jobs []Jobs
	}
)

var (
 
	htmlTplEngine    *template.Template
	htmlTplEngineErr error
)

 
func main() {
 
	fs := http.FileServer(http.Dir("assets/")) //real dir
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	htmlTplEngine = template.New("htmlTplEngine")

	db, err := sql.Open("sqlite3", "jobs.db")
	fmt.Println(err)
	_, htmlTplEngineErr = htmlTplEngine.ParseGlob("views/*.html")
	if nil != htmlTplEngineErr {
		log.Panic(htmlTplEngineErr.Error())
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		//mux.GET("/", func(w http.ResponseWriter, r *http.Request) {
		jobs := []Jobs{}

		result, err := db.Query("SELECT  title,description, salarybase,salarytop,published FROM jobs WHERE id > ?", 0)

		for result.Next() {
			var j Jobs
			err := result.Scan( &j.Title,&j.Description, &j.Salarybase, &j.Salarytop,    &j.Published) // check err
			fmt.Println(err)
			jobs = append(jobs, j)
		}

		fmt.Println(err)
		data := TplData{Jobs: jobs}
		_ = htmlTplEngine.ExecuteTemplate(w, "home.html", data)

	})

	http.HandleFunc("/new", func(w http.ResponseWriter, r *http.Request) {
	
	
		if "POST" == r.Method {
                fmt.Println(r.FormValue("is996"))
			result, err := db.Exec(`INSERT INTO jobs(title, mail,description, salarybase,salarytop,is996,published) VALUES (?, ?, ?,?,?,?,?)`, r.FormValue("title"), r.FormValue("email"),r.FormValue("jd"), r.FormValue("salarybase"), r.FormValue("salarytop"), r.FormValue("is996"),time.Now())

			fmt.Println(result, err)
			http.Redirect(w, r, "/", 302)
		}

		_ = htmlTplEngine.ExecuteTemplate(w, "new.html", nil)
	})

	fmt.Println("ListenAndServe:  http://localhost:8181")
	log.Fatal(http.ListenAndServe(":8181", nil))
}
