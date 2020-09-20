package todo

import (
	"context"
	"github.com/phalpin/GoGoGadgetRESTAPI/todo/models"
	"github.com/phalpin/liberr"
	"github.com/phalpin/liberr/errortype"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const collectionName = "Todos"

type IRepo interface {
	InsertOne(ctx context.Context, todo *models.ToDo) error
	GetOne(ctx context.Context, id string) (*models.ToDo, error)
	UpdateOne(ctx context.Context, todo *models.ToDo) error
	DeleteOne(ctx context.Context, id string) error
}

type Repo struct {
	mongoConnStr string
	mongoDbName  string
	client       *mongo.Client
}

func NewRepo(connStr string, dbName string) *Repo {
	repo := &Repo{
		mongoConnStr: connStr,
		mongoDbName:  dbName,
	}

	client, _ := repo.getClient()
	repo.client = client

	return repo
}

/* Private Methods */
func (t *Repo) getClient() (*mongo.Client, error) {
	if t.client == nil {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(t.mongoConnStr))
		if err != nil {
			return nil, liberr.NewBaseFromError(err)
		}
		t.client = client

		return client, nil
	}
	return t.client, nil
}

func (t *Repo) getCollection(collName string) (*mongo.Collection, error) {
	cl, err := t.getClient()
	if err != nil {
		return nil, liberr.NewKnownFromErr(err, "The Database is experiencing problems. Please try again...", liberr.WithErrorType(errortype.Unknown))
	}

	collection := cl.Database(t.mongoDbName).Collection(collName)

	return collection, nil
}

func (t *Repo) InsertOne(ctx context.Context, todo *models.ToDo) error {

	coll, err := t.getCollection(collectionName)
	if err != nil {
		return liberr.NewBaseFromError(err)
	}

	result, insertErr := coll.InsertOne(ctx, todo)
	if insertErr != nil {
		return liberr.NewBaseFromError(insertErr)
	}

	todo.Id = result.InsertedID.(primitive.ObjectID).String()

	return nil
}

func (t *Repo) GetOne(ctx context.Context, id string) (*models.ToDo, error) {
	coll, err := t.getCollection(collectionName)
	if err != nil {
		return nil, liberr.NewBaseFromError(err)
	}

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, liberr.NewKnownFromErr(err, "Malformed ID", liberr.WithErrorType(errortype.InvalidArgument))
	}

	filter := bson.D{{"_id", objId}}

	retVal := &models.ToDo{}
	findResult := coll.FindOne(ctx, filter)
	if findResult.Err() == nil {
		decodeErr := findResult.Decode(&retVal)
		if decodeErr != nil {
			return nil, liberr.NewBaseFromError(decodeErr)
		}

		defer ctx.Done()
		return retVal, nil
	}

	return nil, nil
}

func (t *Repo) UpdateOne(ctx context.Context, todo *models.ToDo) error {
	coll, err := t.getCollection(collectionName)
	if err != nil {
		return liberr.NewBaseFromError(err)
	}

	objectId, objIdErr := primitive.ObjectIDFromHex(todo.Id)
	if objIdErr != nil {
		return liberr.NewKnown("invalid id passed", "An invalid ID was passed. Please check the ID and try again.", liberr.WithErrorType(errortype.InvalidArgument))
	}

	filter := bson.D{{"_id", objectId}}
	update := bson.M{
		"$set": bson.M{
			"Title":       todo.Title,
			"Description": todo.Description,
			"Completed":   todo.Completed,
		},
	}

	_, updateErr := coll.UpdateOne(ctx, filter, update)
	if updateErr != nil {
		return liberr.NewBaseFromError(updateErr)
	}

	return nil
}

func (t *Repo) DeleteOne(ctx context.Context, id string) error {
	coll, err := t.getCollection(collectionName)
	if err != nil {
		return err
	}

	objectId, objIdErr := primitive.ObjectIDFromHex(id)
	if objIdErr != nil {
		return objIdErr
	}

	filter := bson.D{{"_id", objectId}}
	_, deleteErr := coll.DeleteOne(ctx, filter, nil)
	if deleteErr != nil {
		return deleteErr
	}

	return nil
}
