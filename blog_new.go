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

var DB *sql.DB

func init() {
	var err error
	DB, err = sql.Open("postgres", "user=postgres dbname=postgres")
	if err != nil {
		panic(err)
	}
}

func sqlEsc(in string) string {
	return strings.Replace(in, "'", "''", -1)
}

func SinglePost(g *gas.Gas, postId string) {
	row := DB.QueryRow("SELECT * FROM posts WHERE id = $1", postId)

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
	case "GET":
		g.Render("blog/newpost", nil)
	case "POST":
		now := time.Now().Format("2006-01-02 15:04:05")
		_, err := DB.Exec("INSERT INTO posts (timestamp, title, body) VALUES ('$1', '$2', '$3')",
			now, g.FormValue("title"), g.FormValue("body"))
		if err != nil {
			log.Print(err)
		}
		http.Redirect(g.ResponseWriter, g.Request, "/blog/", http.StatusFound)
	}
}

func Page(g *gas.Gas, page string) {
	pageId, _ := strconv.Atoi(page)
	rows, err := DB.Query("SELECT * FROM posts ORDER BY id DESC OFFSET $1 LIMIT 10", pageId*10)

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
	rows, err := DB.Query("SELECT * FROM posts ORDER BY id DESC LIMIT 10")

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

