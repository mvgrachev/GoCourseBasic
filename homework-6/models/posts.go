package models

import (
	"html/template"
	"github.com/russross/blackfriday"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const ACTIVE_STATUS = 1
const INACTIVE_STATUS = 2

// PostItem - объект поста
type Post struct {
	Mongo `inline`
	Title   string `bson:"title" json:"title"` 
	Date    string `bson:"datetime" json:"date"`
	Summary string `bson:"summary" json:"summary"`
	Body    interface{} `bson:"body" json:"body"`
	Status  int    `bson:"status" json:"status"`
}

// PostItemSlice - массив постов
type PostItemSlice []Post

func (p *Post) GetMongoCollectionName() string {
	return "posts"
}

func (post *Post) Create(ctx context.Context, db *mongo.Database) error {
	post.Status = ACTIVE_STATUS

	post.Id = primitive.NewObjectID()
	coll := db.Collection(post.GetMongoCollectionName())
	_, err := coll.InsertOne(ctx, post)
	if err != nil {
		return err
	}
	return nil
}

func (post *Post) Delete(ctx context.Context, db *mongo.Database) error {
	coll := db.Collection(post.GetMongoCollectionName())
	_, err := coll.UpdateOne(ctx, bson.M{"_id": post.Id}, bson.M{"": INACTIVE_STATUS})
	return err
}

func (post *Post) Update(ctx context.Context, db *mongo.Database) error {
	coll := db.Collection(post.GetMongoCollectionName())
	_, err := coll.UpdateOne(ctx, bson.M{"_id": post.Id}, post)
	return err
}

func GetAllPostItems(ctx context.Context, db *mongo.Database) (PostItemSlice, error) {
	p := Post{}
	coll := db.Collection(p.GetMongoCollectionName())

	cur, err := coll.Find(ctx, bson.M{})

	if err != nil {
		return nil, err
	}
	
	posts := make(PostItemSlice, 0, 8)
	for cur.Next(ctx) {
		post := Post{}
		if err = cur.Decode(&post); err != nil {
			return nil, err
		}
		
		post.Body = template.HTML(blackfriday.MarkdownCommon([]byte(post.Body.(string))))
		posts = append(posts, post)
	}
	cur.Close(ctx)
	return posts, err
}

func GetPost(ctx context.Context, db *mongo.Database, id string) (Post, error) {
	post := Post{}
	coll := db.Collection(post.GetMongoCollectionName())

	postId, _ := primitive.ObjectIDFromHex(id)	
	res := coll.FindOne(ctx, bson.M{"_id": postId})
	//var body string
	//if err := res.Decode(&post.Id, &post.Title, &post.Date, &post.Summary, &body, &post.Status); err != nil {
	if err := res.Decode(&post); err != nil {
		return post, err
	}
	post.Body = template.HTML(blackfriday.MarkdownCommon([]byte(post.Body.(string))))
	return post, nil
}
