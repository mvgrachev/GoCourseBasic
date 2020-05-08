package models

import (
	"database/sql"
	"html/template"
	"github.com/russross/blackfriday"
	"github.com/jinzhu/gorm"
)

const ACTIVE_STATUS = 1
const INACTIVE_STATUS = 2

// PostItem - объект поста
type Post struct {
	Id	string `gorm:"id" json:"id"`
	Title   string `gorm:"type:varchar(255); not null" json:"title"` 
	Date    string `gorm:"-" json:"date"`
	Summary string `gorm:"type:varchar(255); not null" json:"summary"`
	Body    interface{} `json:"body"`
	Status  int    `json:"status"`
}

// PostItemSlice - массив постов
type PostItemSlice []Post

func (post *Post) Create(db *gorm.DB) error {
	post.Status = ACTIVE_STATUS
	result := db.Create(post)
	return result.Error
}

func (post *Post) Delete(db *gorm.DB) error {
	chPost := &Post{}
	result := db.First(chPost, post.Id)
	if result.Error != nil {
		return result.Error
	}
	result = db.Model(&chPost).Update("status", INACTIVE_STATUS)
	return result.Error
}

func (post *Post) Update(db *gorm.DB) error {
	
	chPost := &Post{}
	result := db.First(chPost, post.Id)
	if result.Error != nil {
		return result.Error
	}
	result = db.Model(&chPost).Update("body", post.Body)
	return result.Error
	
}

func GetAllPostItems(db *gorm.DB) (PostItemSlice, error) {
	rows, err := db.Raw("SELECT id, title, created_at, summary, body, status FROM posts WHERE status = ? ORDER BY id DESC", ACTIVE_STATUS).Rows()
	//rows, err := db.Model(&PostItem{}).Where("status = ?", ACTIVE_STATUS).Select("id, title, created_at, summary, body, status").Rows()
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	posts := make(PostItemSlice, 0, 8)
	for rows.Next() {
		post := Post{}
		var body string
		if err = rows.Scan(&post.Id, &post.Title, &post.Date, &post.Summary, &body, &post.Status); err != nil {
			return nil, err
		}
		
		post.Body = template.HTML(blackfriday.MarkdownCommon([]byte(body)))
		posts = append(posts, post)
	}
	return posts, err
}

func GetPost(db *gorm.DB, id string) (Post, error) {
	row := db.Raw("SELECT id, title, created_at, summary, body, status FROM posts WHERE status = ? AND id = ? ORDER BY id DESC", ACTIVE_STATUS, id).Row()
	post := Post{}
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
