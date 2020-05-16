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

type Inquirer interface {
	GetMongoCollectionName() string
	Create(context.Context, *mongo.Database ) error
	Read(context.Context, *mongo.Database ) error
	Update(context.Context, *mongo.Database ) error
	Delete(context.Context, *mongo.Database ) error
}

// PostItem - объект поста
type Post struct {
	Mongo `inline`
	Title   string `bson:"title" json:"title"` 
	Date    string `bson:"datetime" json:"date"`
	Summary string `bson:"summary" json:"summary"`
	Body    interface{} `bson:"body" json:"body"`
	Status  int    `bson:"status" json:"status"`
}

func CreatePost(ctx context.Context, i Inquirer, db *mongo.Database) error {
	i.Create(ctx, db)
	return nil
}

func ReadPost(ctx context.Context, i Inquirer, db *mongo.Database) error {
	i.Read(ctx, db)
	return nil
}

func UpdatePost(ctx context.Context, i Inquirer, db *mongo.Database) error {
	i.Update(ctx, db)
	return nil
}

func DeletePost(ctx context.Context, i Inquirer, db *mongo.Database) error {
	i.Delete(ctx, db)
	return nil
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

func (post *Post) Read(ctx context.Context, db *mongo.Database) error {
	post.Id, _ = primitive.ObjectIDFromHex(post.Id.(string))
	coll := db.Collection(post.GetMongoCollectionName())
	res := coll.FindOne(ctx, bson.M{"_id": post.Id})
	if err := res.Decode(post); err != nil {
		return err
	}
	post.Body = template.HTML(blackfriday.MarkdownCommon([]byte(post.Body.(string))))
	post.Id = post.Id.(primitive.ObjectID).Hex()
	return nil
}

func (post *Post) Update(ctx context.Context, db *mongo.Database) error {
	post.Id, _ = primitive.ObjectIDFromHex(post.Id.(string))
	coll := db.Collection(post.GetMongoCollectionName())
	_, err := coll.ReplaceOne(ctx, bson.M{"_id": post.Id}, bson.M{"title":post.Title,"summary":post.Summary,"body":post.Body},)
	return err
}

func (post *Post) Delete(ctx context.Context, db *mongo.Database) error {
	post.Id, _ = primitive.ObjectIDFromHex(post.Id.(string))
	coll := db.Collection(post.GetMongoCollectionName())
	_, err := coll.UpdateOne(ctx, bson.M{"_id": post.Id}, bson.D{{"$set", bson.D{{"status", INACTIVE_STATUS}}},},)
	return err
}

func GetAllPostItems(ctx context.Context, db *mongo.Database) (PostItemSlice, error) {
	p := Post{}
	coll := db.Collection(p.GetMongoCollectionName())

	cur, err := coll.Find(ctx, bson.M{"status": ACTIVE_STATUS})

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
		post.Id = post.Id.(primitive.ObjectID).Hex()
		posts = append(posts, post)
	}
	cur.Close(ctx)
	return posts, err
}
