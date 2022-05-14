package user

import (
	"context"
	"go-microservices/common"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserDocument struct {
	Id            string `bson:"id"`
	PreferredName string `bson:"preferred_name"`
	Email         string `bson:"email"`
}

type MongoUserStore struct {
	collection *mongo.Collection
}

func NewMongoUserStore() MongoUserStore {
	ctx := context.Background()
	credentials := options.Credential{
		Username: "user",
		Password: "password",
	}
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"), options.Client().SetAuth(credentials))
	common.PanicOnError(err)
	collection := client.Database("human_resources").Collection("users")
	return MongoUserStore{
		collection: collection,
	}
}

func (store *MongoUserStore) Save(user User) error {
	document := UserDocument{
		Id:            user.id.String(),
		PreferredName: user.preferredName,
		Email:         user.email,
	}
	_, err := store.collection.ReplaceOne(context.Background(), bson.D{{"id", user.id.String()}}, document, options.Replace().SetUpsert(true))
	return err
}

// TODO: Missing index
func (store *MongoUserStore) Get(userId UserID) (User, error) {
	var document UserDocument
	err := store.collection.FindOne(context.Background(), bson.D{{"id", userId.String()}}).Decode(&document)
	if err != nil {
		return User{}, convertMongoError(err)
	}
	return parseMongoUserDocument(document)
}

// TODO: Missing index
func (store *MongoUserStore) GetByEmail(email string) (User, error) {
	var document UserDocument
	err := store.collection.FindOne(context.Background(), bson.D{{"email", email}}).Decode(&document)
	if err != nil {
		return User{}, err
	}
	return parseMongoUserDocument(document)
}

func (store *MongoUserStore) GetMany() ([]User, error) {
	var documents []UserDocument
	cursor, err := store.collection.Find(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	err = cursor.All(context.Background(), &documents)
	if err != nil {
		return nil, err
	}

	users := make([]User, len(documents))
	for i, document := range documents {
		user, err := parseMongoUserDocument(document)
		if err != nil {
			return nil, err
		}
		users[i] = user
	}

	return users, nil
}

func (store *MongoUserStore) EmailExists(email string) (bool, error) {
	docs, err := store.collection.CountDocuments(context.Background(), bson.D{{"email", email}})
	if err != nil {
		return false, err
	}
	return docs > 0, nil
}

func parseMongoUserDocument(document UserDocument) (User, error) {
	id, err := common.ParseID(document.Id)
	if err != nil {
		return User{}, err
	}
	return User{
		id:            UserID{id},
		preferredName: document.PreferredName,
		email:         document.Email,
	}, nil
}

func convertMongoError(err error) error {
	if err == mongo.ErrNoDocuments {
		return ErrUserNotFound
	}
	return err
}

func (store *MongoUserStore) Clear() {
	_, err := store.collection.DeleteMany(context.Background(), bson.D{})
	common.PanicOnError(err)
}
