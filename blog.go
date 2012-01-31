// This file is deprecated. It's only here so I can pick and choose the good ideas from the bad before I forget them.
package blog

import (
	"net/http"
	"pgsql"
	"text/template"
	//	"io"
	"errors"
	"strings"
	//	"os"
	"fmt"
	"regexp"
	"strconv"
	//"time"
)

var templates map[string]*template.Template

// A blog post
type Post struct {
	Id        int64
	Timestamp string
	Title     string
	Body      string
}

// Advance to the next row and return true if successful
func nextRow(rs *pgsql.ResultSet) bool {
	gotRow, err := rs.FetchNext()
	if gotRow && (err == nil) {
		return true
	}
	println(err)
	return false
}

// Get a single post with given id from database
func postFromDb(pg *pgsql.Conn, id int64) *Post {
	rs, _ := pg.Query(fmt.Sprintf("SELECT * FROM posts WHERE id = %d;", id))
	defer rs.Close()

	if nextRow(rs) {
		timestamp, _, _ := rs.Time(1)
		title, _, _ := rs.String(2)
		body, _, _ := rs.String(3)
		return &Post{id, timestamp.Format("02 Jan 2006 15:04"), title, body}
	}

	return nil
}

// Get posts bewteen specified ids from database
func postsFromDb(pg *pgsql.Conn, fromId, toId int64) []*Post {
	println("blog.postsFromDb()")
	rs, _ := pg.Query(fmt.Sprintf("SELECT * FROM posts WHERE id >= %d AND id <= %d ORDER BY id DESC", fromId, toId))
	defer rs.Close()

	p := make([]*Post, (toId - fromId + 1))
	for i := 0; nextRow(rs); i++ {
		id, _, _ := rs.Int64(0)
		timestamp, _, _ := rs.Time(1)
		title, _, _ := rs.String(2)
		body, _, _ := rs.String(3)
		p[i] = &Post{id, timestamp.Format("02 Jan 2006 15:04"), title, body}
	}
	return p
}

// shorter regexp matcher
func match(s, re string) bool {
	return regexp.MustCompile(re).MatchString(s)
}

// Render a template to w, can either be a preparsed template or a string
func render(w http.ResponseWriter, tmpl, data interface{}) error {
	print("blog.render()...")
	switch tmpl.(type) {
	case *template.Template:
		println("it's a template")
		return tmpl.(*template.Template).Execute(w, data)
	case string:
		println("it's a string")
		t, _ := template.New("string").Parse(tmpl.(string))
		return t.Execute(w, data)
	}
	println("wat")
	return errors.New("Erroneous template type")
}

func onoez(w http.ResponseWriter, errorCode, path string) {
	render(w, templates[errorCode], path)
}

// main handler
func Gas(w http.ResponseWriter, r *http.Request) {
	println("requested:", r.URL.Path)

	pg, err := pgsql.Connect("user=postgres timeout=5", pgsql.LogFatal)
	if err != nil {
		println("pgsql.Connect():", err.Error())
		return
	}
	defer pg.Close()

	path := strings.Split(r.URL.Path, "/")
	fmt.Printf("%#v\n", path)
	switch len(path) {
	case 1:
		if path[0] != "" {
			http.Redirect(w, r, "..", http.StatusSeeOther)
		} else {
			println("trying to render index")
			render(w, templates["index"], postsFromDb(pg, 1, 2))
		}

	case 2:
		switch path[0] {
		case "post":
			println("trying to render post", path[1])
			if match(path[1], `\d+`) {
				postId, err := strconv.ParseInt(path[1], 10, 64)
				if err != nil {
					println("ERROR:", err.Error())
					onoez(w, "503", r.URL.Path)
				} else {
					render(w, templates["onepost"], postFromDb(pg, postId))
				}
			} else {
				println("regex match failed for post")
				onoez(w, "503", r.URL.Path)
			}
		case "page":
			println("trying to render page", path[1])
			if match(path[1], `\d+`) {
				page, err := strconv.ParseInt(path[1], 10, 64)
				fromId, toId := (page+1)*10, fromId+10
			} else {
				println("regex match failed for page")
				onoez(w, "503", r.URL.Path)
			}

	default:
		println("sending a 503")
		onoez(w, "503", r.URL.Path)
	}
}

func init() {
	templates = make(map[string]*template.Template)
	var err error
	for _, t := range []string{"index", "onepost", "503"} {
		templates[t], err = template.ParseFiles("./blog/"+t+".html.got")
		if err != nil {
			panic(err)
		}
	}
}
