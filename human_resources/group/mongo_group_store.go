package group

import (
	"context"
	"go-microservices/common"
	"go-microservices/human_resources/user"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoGroupStore struct {
	collection *mongo.Collection
}

func NewMongoGroupStore() MongoGroupStore {
	ctx := context.Background()
	credentials := options.Credential{
		Username: "user",
		Password: "password",
	}
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"), options.Client().SetAuth(credentials))
	common.PanicOnError(err)
	collection := client.Database("human_resources").Collection("groups")
	return MongoGroupStore{
		collection: collection,
	}
}

func (store *MongoGroupStore) Save(group Group) error {
	document := GroupDocument{
		Id:          group.id.String(),
		Name:        group.name,
		Description: group.description,
		Members:     make([]MembershipDocument, len(group.members)),
	}
	for i, membership := range group.members {
		document.Members[i] = MembershipDocument{
			UserID:   membership.userID.String(),
			JoinedAt: membership.joinedAt,
		}
	}

	_, err := store.collection.ReplaceOne(context.Background(), bson.D{{"id", group.id.String()}}, document, options.Replace().SetUpsert(true))
	return err
}

// TODO: Missing indexes
func (store *MongoGroupStore) Get(groupId GroupID) (Group, error) {
	var document GroupDocument
	err := store.collection.FindOne(context.Background(), bson.D{{"id", groupId.String()}}).Decode(&document)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			err = ErrGroupNotFound
		}
		return Group{}, err
	}
	return parseMongoGroupDocument(document)
}

// TODO: Missing performance (index?)
func (store *MongoGroupStore) FindForUser(userId user.UserID) ([]GroupDto, error) {
	pipeline := mongo.Pipeline{
		bson.D{{"$match", bson.D{
			{"members", bson.D{{"$elemMatch", bson.D{{"user_id", userId.String()}}}}},
		}}},
		bson.D{{"$project", bson.D{
			{"members_count", bson.D{{"$size", "$members"}}},
			{"id", true},
			{"name", true},
			{"description", true},
		}}},
		bson.D{{"$unset", "members"}},
	}
	var documents []GroupDto
	cur, err := store.collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())

	err = cur.All(context.Background(), &documents)

	return documents, err
}

type GroupDocument struct {
	Id          string               `bson:"id"`
	Name        string               `bson:"name"`
	Description string               `bson:"description"`
	Members     []MembershipDocument `bson:"members"`
}

func parseMongoGroupDocument(document GroupDocument) (Group, error) {
	groupID, err := common.ParseID(document.Id)
	if err != nil {
		return Group{}, err
	}
	group := Group{
		id:          GroupID{groupID},
		name:        document.Name,
		description: document.Description,
		members:     make([]Membership, len(document.Members)),
	}
	for i, membershipDocument := range document.Members {
		membership, err := parseMongoMembershipDocument(membershipDocument)
		if err != nil {
			return Group{}, err
		}
		group.members[i] = membership
	}
	return group, nil
}

type MembershipDocument struct {
	UserID   string    `bson:"user_id"`
	JoinedAt time.Time `bson:"joined_at"`
}

func parseMongoMembershipDocument(document MembershipDocument) (Membership, error) {
	userID, err := common.ParseID(document.UserID)
	if err != nil {
		return Membership{}, err
	}
	return Membership{
		userID:   user.UserID{userID},
		joinedAt: document.JoinedAt,
	}, nil
}
