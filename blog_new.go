package blog

import (
	_ "github.com/jbarham/gopgsqldriver"
	"database/sql"
	"net/http"
	"gas"
	"log"
	"strconv"
	"strings"
	"time"
)

type Post struct {
	Id		int64
	Time	time.Time
	Title	string
	Body	string
}

func db() *sql.DB {
	database, err := sql.Open("postgres", "user=postgres dbname=postgres")
	if err != nil {
		panic(err)
	}
	return database
}

func sqlEsc(in string) string {
	return strings.Replace(in, "'", "''", -1)
}

func SinglePost(g *gas.Gas, postId string) {
	row := db().QueryRow("SELECT * FROM posts WHERE id = $1", postId)

	var (
		id		int64
		stamp	string
		title	string
		body	string
	)

	err := row.Scan(&id, &stamp, &title, &body)
	if err != nil {
		log.Printf("blog.SinglePost(): %v", err)
		g.HTTPError(http.StatusServiceUnavailable)
		return
	}

	timestamp, _ := time.Parse("2006-01-02 15:04:05", stamp)
	g.Render("blog/onepost", &Post{id, timestamp, title, body})
}

func NewPost(g *gas.Gas) {
	switch g.Method {
	case "POST":
		g.Render("blog/onepost", &Post{0, time.Now(), g.FormValue("title"), g.FormValue("body")})
	case "GET":
		g.Render("blog/newpost", nil)
	}
}

func Page(g *gas.Gas, page string) {
	pageId, _ := strconv.Atoi(page)
	rows, err := db().Query("SELECT * FROM posts ORDER BY id DESC OFFSET $1 LIMIT 10", pageId*10)

	if err != nil {
		log.Printf("blog.Page(): %v", err)
		g.HTTPError(http.StatusServiceUnavailable)
	}

	posts := []*Post{}
	for rows.Next() {
		var (
			id		int64
			stamp	string
			title	string
			body	string
		)

		err = rows.Scan(&id, &stamp, &title, &body)
		if err != nil {
			g.HTTPError(http.StatusServiceUnavailable)
			return
		}
		timestamp, _ := time.Parse("2006-01-02 15:04:05", stamp)
		posts = append(posts, &Post{id, timestamp, title, body})
	}

	g.Render("blog/index", posts)
}

// TODO: add pagination
func Index(g *gas.Gas) {
	rows, err := db().Query("SELECT * FROM posts ORDER BY id DESC LIMIT 10")

	if err != nil {
		log.Printf("blog.Index(): %v", err)
		g.HTTPError(http.StatusServiceUnavailable)
		return
	}

	posts := []*Post{}
	for rows.Next() {
		// TODO: better way to do this?
		var (
			id		int64
			stamp	string
			title	string
			body	string
		)

		err = rows.Scan(&id, &stamp, &title, &body)
		if err != nil {
			log.Printf("rows.Scan(): %v", err)
			g.HTTPError(http.StatusServiceUnavailable)
			// TODO: make it so you don't have to return on error
			// (add panic() and recover()?)
			return
		}
		timestamp, _ := time.Parse("2006-01-02 15:04:05", stamp)
		posts = append(posts, &Post{id, timestamp, title, body})
	}

	g.Render("blog/index", posts)
}

