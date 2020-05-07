package models

import (
	"database/sql"
	"html/template"
	"github.com/russross/blackfriday"
)

const ACTIVE_STATUS = 1
const INACTIVE_STATUS = 2

// PostItem - объект поста
type PostItem struct {
	Id	string `json:"id"`
	Title   string `json:"title"` 
	Date    string `json:"date"`
	Summary string `json:"summary"`
	Body    interface{} `json:"body"`
	Status  int    `json:"status"`
}

// PostItemSlice - массив постов
type PostItemSlice []PostItem

func (post *PostItem) Insert(db *sql.DB) error {
	_, err := db.Exec(
		"INSERT INTO posts (title, summary, body, status) VALUES (?, ?, ?, ?)",
		post.Title, post.Summary, post.Body, ACTIVE_STATUS,
	)
	return err
}

func (post *PostItem) Delete(db *sql.DB) error {
	_, err := db.Exec(
		"UPDATE posts SET status = ? WHERE id = ?",
		INACTIVE_STATUS, post.Id,
	)
	return err
}

func (post *PostItem) Update(db *sql.DB) error {
	_, err := db.Exec(
		"UPDATE posts SET body = ? WHERE id = ?",
		post.Body, post.Id,
	)
	return err
}

func GetAllPostItems(db *sql.DB) (PostItemSlice, error) {
	rows, err := db.Query("SELECT id, title, created_at, summary, body, status FROM posts WHERE status = ? ORDER BY id DESC", ACTIVE_STATUS)
	if err != nil {
		return nil, err
	}
	posts := make(PostItemSlice, 0, 8)
	for rows.Next() {
		post := PostItem{}
		var body string
		if err = rows.Scan(&post.Id, &post.Title, &post.Date, &post.Summary, &body, &post.Status); err != nil {
			return nil, err
		}
		
		post.Body = template.HTML(blackfriday.MarkdownCommon([]byte(body)))
		posts = append(posts, post)
	}
	return posts, err
}

func GetPost(db *sql.DB, id string) (PostItem, error) {
	row := db.QueryRow("SELECT id, title, created_at, summary, body, status FROM posts WHERE status = ? AND id = ? ORDER BY id DESC", ACTIVE_STATUS, id)

	post := PostItem{}
	var body string
	err := row.Scan(&post.Id, &post.Title, &post.Date, &post.Summary, &body, &post.Status)	
	switch {
	case err == sql.ErrNoRows:
		return post, err
	case err != nil:
		return post, err
	}

	post.Body = template.HTML(blackfriday.MarkdownCommon([]byte(body)))
	
	return post, err
}
