package blog

import (
	_ "github.com/jbarham/gopgsqldriver"
	"database/sql"
	"net/http"
	"gas"
	"log"
	"strconv"
)

type Post struct {
	Id		int64
	Time	string
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

func SinglePost(g *gas.Gas, id string) {
	row := db().QueryRow("SELECT * FROM posts WHERE id = ?", id)

	post := &Post{}
	if err := row.Scan(&post.Id, &post.Time, &post.Title, &post.Body); err != nil {
		g.HTTPError(http.StatusServiceUnavailable)
	}

	g.Render("blog/onepost", post)
}

func Page(g *gas.Gas, page string) {
	pageId, _ := strconv.Atoi(page)
	rows, err := db().Query("SELECT * FROM posts ORDER BY id DESC OFFSET ? LIMIT 10", pageId*10)

	if err != nil {
		log.Printf("blog.Page(): %v", err)
		g.HTTPError(http.StatusServiceUnavailable)
	}
	posts := make([]*Post, 10)
	for i := 0; rows.Next(); i++ {
		err = rows.Scan(&posts[i].Id, &posts[i].Time, &posts[i].Title, &posts[i].Body)
		if err != nil {
			g.HTTPError(http.StatusServiceUnavailable)
		}
	}

	g.Render("blog/index", posts)
}

func Index(g *gas.Gas) {
	rows, err := db().Query("SELECT * FROM posts ORDER BY id DESC LIMIT 10")

	if err != nil {
		log.Printf("blog.Index(): %v", err)
		g.HTTPError(http.StatusServiceUnavailable)
		return
	}

	posts := make([]*Post, 10)
	for i := 0; rows.Next(); i++ {
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
			return
		}
		posts[i] = &Post{id, stamp, title, body}
	}

	g.Render("blog/index", posts)
}

