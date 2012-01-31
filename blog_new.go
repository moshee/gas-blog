package blog

import (
	_ "github.com/jbarham/pgsqldriver"
	"database/sql"
	"net/http"
	//"strings"
	//"regexp"
	"text/template"
	"time"
	"../gas/gas"
)

func init() {
	gas.RegisterTemplates("blog/index", "blog/onepost", "blog/503")
}

type Post struct {
	Id		int64
	Time	time.Time
	Title	string
	Body	string
}

func db() *sql.DB {
	database, err := sql.Open("postgres", "dbname=postgres")
	if err != nil {
		panic(err)
	}
	return database
}

func Post(g gas.Request, id int64) {
	row := db().QueryRow("SELECT * FROM posts WHERE id = ?", id)
	post := &Post{}
	if err = row.Scan(&post.Id, &post.Time, &post.Title, &post.Body); err != nil {
		gas.HTTPError(http.ErrorServiceUnavailable)
	}

	gas.Render("onepost", post)
}

func Page(g gas.Request, page int) {
	rows, err := db().Query("SELECT * FROM posts ORDER BY id DESC OFFSET ? LIMIT 10", page*10)
	if err != nil {
		println("blog.Page():", err.Error)
		gas.HTTPError(http.ErrorServiceUnavailable)
	}

	posts := make([]*Post, toId-fromId+1)
	for i := 0; rows.Next(); i++ {
		err = rows.Scan(&posts[i].Id, &posts[i].Time, &posts[i].Title, &posts[i].Body)
		if err != nil {
			gas.HTTPError(http.ErrorServiceUnavailable)
		}
	}

	gas.Render("blog/index", posts)
}

func Index(g gas.Request) {
	rows, err := db().Query("SELECT * FROM posts ORDER BY id DESC LIMIT 10")
	if err != nil {
		println("blog.Index():", err.Error)
		gas.HTTPError(http.ErrorServiceUnavailable)
	}

	posts := make([]*Post, 10)
	for i := 0; rows.Next(); i++ {
		err = rows.Scan(&posts[i].Id, &posts[i].Time, &posts[i].Title, &posts[i].Body)
		if err != nil {
			gas.HTTPError(http.ErrorServiceUnavailable)
		}
	}

	gas.Render("blog/index", posts)
}

