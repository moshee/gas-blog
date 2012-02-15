package blog

import (
	_ "github.com/jbarham/gopgsqldriver"
	"database/sql"
	"net/http"
	"gas"
	"log"
	"strconv"
	"time"
)

type Post struct {
	Id		int64
	Time	time.Time
	Title	string
	Body	string
}

var (
	DB *sql.DB
	QuerySinglePost, QueryPage, QueryNewPost *sql.Stmt
)

func init() {
	var err error
	DB, err = sql.Open("postgres", "user=postgres dbname=postgres")
	if err != nil {
		panic(err)
	}
	QuerySinglePost, _ = DB.Prepare("SELECT * FROM posts WHERE id = $1")
	QueryPage, _ =  DB.Prepare("SELECT * FROM posts ORDER BY id DESC LIMIT 10 OFFSET $1")
	QueryNewPost, _ = DB.Prepare("INSERT INTO posts (timestamp, title, body) VALUES ($1, $2, $3)")
}

func SinglePost(g *gas.Gas, postId string) {
	row := QuerySinglePost.QueryRow(postId)

	var (
		id		int64
		stamp	string
		title	string
		body	string
	)

	err := row.Scan(&id, &stamp, &title, &body)
	if err != nil {
		log.Printf("blog.SinglePost(): %v", err)
		g.ErrorPage(http.StatusServiceUnavailable)
		return
	}

	timestamp, _ := time.Parse("2006-01-02 15:04:05-07", stamp)
	g.Render("blog/onepost", &Post{id, timestamp, title, body})
}

func NewPost(g *gas.Gas) {
	switch g.Method {
	case "GET":
		g.Render("blog/newpost", nil)
	case "POST":
		now := time.Now().Format("2006-01-02 15:04:05-07")
		_, err := QueryNewPost.Exec(now, g.FormValue("title"), g.FormValue("body"))
		if err != nil {
			log.Print(err)
			return
		}
		http.Redirect(g.ResponseWriter, g.Request, "/blog/", http.StatusFound)
		return
	}
}

func Page(g *gas.Gas, page string) {
	pageId, _ := strconv.Atoi(page)
	rows, err := QueryPage.Query(pageId * 10)

	if err != nil {
		log.Printf("blog.Page(): %v", err)
		g.ErrorPage(http.StatusServiceUnavailable)
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
			g.ErrorPage(http.StatusServiceUnavailable)
			return
		}
		timestamp, _ := time.Parse("2006-01-02 15:04:05-07", stamp)
		posts = append(posts, &Post{id, timestamp, title, body})
	}
	if len(posts) < 1 {
		g.ErrorPage(http.StatusNotFound)
	}

	g.Render("blog/index", posts)
}

// TODO: add pagination
func Index(g *gas.Gas) {
	// Note to self: use `SELECT setval('posts_id_seq', currval('posts_id_seq')-1);`
	// when deleting a row so that the index gets set properly
	rows, err := QueryPage.Query("0")

	if err != nil {
		log.Printf("blog.Index(): %v", err)
		g.ErrorPage(http.StatusServiceUnavailable)
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
			g.ErrorPage(http.StatusServiceUnavailable)
			// TODO: make it so you don't have to return on error
			// (add panic() and recover()?)
			return
		}
		timestamp, _ := time.Parse("2006-01-02 15:04:05-07", stamp)
		posts = append(posts, &Post{id, timestamp, title, body})
	}

	g.Render("blog/index", posts)
}

