package models

type Mongo struct {
	Id interface{} `bson:"_id"`
}

func (m *Mongo) GetMongoCollectionName() string {
	panic("GetMongoCollectionName not implemented")
	return ""
}
