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

	var (
		id		int64
		stamp	string
		title	string
		body	string
	)

	err = row.Scan(&id, &stamp, &title, &body)
	if err != nil {
		g.HTTPError(http.StatusServiceUnavailable)
		return
	}

	g.Render("blog/onepost", &Post{id, stamp, title, body})
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
		posts[i] = &Post{id, stamp, title, body}
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

	posts := make([]*Post, 10)
	for i := 0; rows.Next(); i++ {
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
		posts[i] = &Post{id, stamp, title, body}
	}

	g.Render("blog/index", posts)
}

