package blog

import (
	"github.com/russross/blackfriday"
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
	Tag		string
}

var (
	DB *sql.DB
	QuerySinglePost, QueryPage, QueryNewPost, QueryEditPost *sql.Stmt
)

func init() {
	var err error
	DB, err = sql.Open("postgres", "user=postgres dbname=postgres")
	if err != nil {
		panic(err)
	}
	QuerySinglePost, _ = DB.Prepare("SELECT * FROM posts WHERE id = $1")
	QueryPage, _ =  DB.Prepare("SELECT * FROM posts ORDER BY id DESC LIMIT 10 OFFSET $1")
	QueryNewPost, _ = DB.Prepare("INSERT INTO posts (timestamp, title, body, tag) VALUES ($1, $2, $3, $4)")
	QueryEditPost, _ = DB.Prepare("UPDATE posts SET body = $1 WHERE id = $2")
}

/////////////////////////////
// Data-getters
/////////////////////////////

//func getPost(g *gas.Gas) *Post {

//}

func getPosts(g *gas.Gas, postId interface{}) []*Post {
	rows, err := QueryPage.Query(postId)

	if err != nil {
		g.ErrorPage(http.StatusServiceUnavailable)
		panic(err)
	}

	posts := []*Post{}
	for rows.Next() {
		// TODO: better way to do this?
		var (
			id		int64
			stamp	string
			title	string
			body	[]byte
			tag		string
		)

		err = rows.Scan(&id, &stamp, &title, &body, &tag)
		if err != nil {
			g.ErrorPage(http.StatusServiceUnavailable)
			// TODO: make it so you don't have to return on error
			// (add panic() and recover()?)
//			panic(err)
		}
		timestamp, _ := time.Parse("2006-01-02 15:04:05-07", stamp)
		posts = append(posts, &Post{id, timestamp, title, string(blackfriday.MarkdownCommon(body)), tag})
	}

	return posts
}

/////////////////////////////
// Actions
/////////////////////////////

func NewPost(g *gas.Gas) {
	switch g.Method {
	case "GET":
		g.Render("blog/newpost", nil)
	case "POST":
		now := time.Now().Format("2006-01-02 15:04:05-07")
		_, err := QueryNewPost.Exec(now, g.FormValue("title"), g.FormValue("body"), g.FormValue("tag"))
		if err != nil {
			log.Print(err)
			return
		}
		http.Redirect(g.ResponseWriter, g.Request, "/blog/", http.StatusFound)
	}
}

//func PreviewPost(g *gas.Gas) {


func EditPost(g *gas.Gas, postId string) {
	switch g.Method {
	case "GET":
		row := QuerySinglePost.QueryRow(postId)
		var (
			id int64
			tag string
			stamp, title string
			body string
		)
		err := row.Scan(&id, &stamp, &title, &body, &tag)
		if err != nil {
			log.Printf("blog.EditPost(): %v", err)
			g.ErrorPage(http.StatusServiceUnavailable)
		}

		timestamp, _ := time.Parse("2006-01-02 15:04:05-07", stamp)
		g.Render("blog/editpost", &Post{id, timestamp, title, body, tag})

	case "POST":
		QueryEditPost.Exec(g.FormValue("body"), postId)
		http.Redirect(g.ResponseWriter, g.Request, "/blog/", http.StatusFound)
	}

}

func SinglePost(g *gas.Gas, postId string) {
	row := QuerySinglePost.QueryRow(postId)

	var (
		id		int64
		stamp	string
		title	string
		body	[]byte
		tag		string
	)

	err := row.Scan(&id, &stamp, &title, &body, &tag)
	if err != nil {
		log.Printf("blog.SinglePost(): %v", err)
		g.ErrorPage(http.StatusServiceUnavailable)
	}

	timestamp, _ := time.Parse("2006-01-02 15:04:05-07", stamp)
	g.Render("blog/onepost", &Post{id, timestamp, title, string(blackfriday.MarkdownCommon(body)), tag})
}

func Page(g *gas.Gas, page string) {
	pageId, _ := strconv.Atoi(page)
	posts := getPosts(g, pageId)
	if len(posts) < 1 {
		g.ErrorPage(http.StatusNotFound)
	}

	g.Render("blog/index", posts)
}

// TODO: add pagination
func Index(g *gas.Gas) {
	// Note to self: use `SELECT setval('posts_id_seq', currval('posts_id_seq')-1);`
	// when deleting a row so that the index gets set properly

	g.Render("blog/index", getPosts(g, "0"))
}

func RSS(g *gas.Gas) {
	g.ResponseWriter.Header().Set("Content-Type", "application/rss+xml")
	g.Render("blog/rss", getPosts(g, "0"))
}
