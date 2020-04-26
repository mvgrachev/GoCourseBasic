package models

import (
	"database/sql"
	"html/template"
	"strings"
	"io/ioutil"
	"github.com/russross/blackfriday"
)

// PostItem - объект поста
type PostItem struct {
	Id	string `json:"id"`
	Title   string `json:"title"` 
	Date    string `json:"date"`
	Summary string `json:"summary"`
	Body    template.HTML `json:"body"`
	File    string `json:"file"`
}

// PostItemSlice - массив постов
type PostItemSlice []PostItem

func (post *PostItem) Insert(db *sql.DB) error {
	_, err := db.Exec(
		"INSERT INTO posts (id, file, status) VALUES (?, ?, ?)",
		post.Id, post.File, 1,
	)
	return err
}

func (post *PostItem) Delete(db *sql.DB) error {
	_, err := db.Exec(
		"UPDATE posts SET status = ? WHERE id = ?",
		2, post.Id,
	)
	return err
}

func (post *PostItem) Update(db *sql.DB) error {
	_, err := db.Exec(
		"UPDATE posts SET file = ?, WHERE id = ?",
		post.File, post.Id,
	)
	return err
}

func GetAllPostItems(db *sql.DB) (PostItemSlice, error) {
	rows, err := db.Query("SELECT id, file FROM posts WHERE status = ? ORDER BY id DESC", 1)
	if err != nil {
		return nil, err
	}
	posts := make(PostItemSlice, 0, 8)
	for rows.Next() {
		post := PostItem{}
		var id string
		var file string
		if err = rows.Scan(&id, &file); err != nil {
			return nil, err
		}
		post = getPostFromFile(id, file)
		posts = append(posts, post)
	}
	return posts, err
}

func GetPost(db *sql.DB, id string) (PostItem, error) {
	var file string

	row := db.QueryRow("SELECT file FROM posts WHERE status = ? AND id = ? ORDER BY id DESC", 1, id)

	post := PostItem{}
	err := row.Scan(&file)	
	switch {
	case err == sql.ErrNoRows:
		return post, err
	case err != nil:
		return post, err
	}

	post = getPostFromFile( id, file) 
	
	return post, err
}

func getPostFromFile( id string, fname string) PostItem {
	fpath := []string{"./posts", fname}
	f := strings.Join(fpath, "/")    
	file := strings.Replace(f, "./posts/", "", -1)
	file = strings.Replace(file, ".md", "", -1)
	fileread, _ := ioutil.ReadFile(f)
	lines := strings.Split(string(fileread), "\n")
	title := string(lines[0])
	date := string(lines[1])
	summary := string(lines[2])
	body := strings.Join(lines[3:len(lines)], "\n")
	postHtml := template.HTML(blackfriday.MarkdownCommon([]byte(body)))
	
	post := PostItem{id, title, date, summary, postHtml, file}

	return post
}
