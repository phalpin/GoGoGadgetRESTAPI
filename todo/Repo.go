package todo

import (
	"context"
	"errors"
	"github.com/phalpin/GoGoGadgetRESTAPI/todo/models"
	"github.com/phalpin/GoGoGadgetRESTAPI/todo/pherr"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const collectionName = "Todos"
const insertInvalidDueToIdExistingAlready = "unable to insert a record that is already saved"
const mustPassIdToGetOneMethod = "you must provide an id when retrieving a saved record"
const mustPassIdToUpdateOneMethod = "you must provide an id when updating a saved record"
const mustPassIdToDeleteOneMethod = "you must provide an id when deleting a saved record"

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
			return nil, err
		}
		t.client = client
		//defer client.Disconnect(ctx)

		return client, nil
	}
	return t.client, nil
}

func (t *Repo) getCollection(collName string) (*mongo.Collection, error) {
	cl, err := t.getClient()
	if err != nil {
		return nil, pherr.NewKnown(err, "The Database is experience problems. Please try again...", pherr.WithErrorType(pherr.Unknown))
	}

	collection := cl.Database(t.mongoDbName).Collection(collName)

	return collection, nil
}

func (t *Repo) InsertOne(ctx context.Context, todo *models.ToDo) error {
	if todo.Id != "" {
		return errors.New(insertInvalidDueToIdExistingAlready)
	}

	coll, err := t.getCollection(collectionName)
	if err != nil {
		return err
	}

	result, insertErr := coll.InsertOne(ctx, todo)
	if insertErr != nil {
		return insertErr
		//Note: Need to either use a proper logging package or PrintF and come up with some kind of formatting.
		//log.Fatal("Failure to insert record:", insertErr.Error())
	}

	todo.Id = result.InsertedID.(primitive.ObjectID).String()

	return nil
}

func (t *Repo) GetOne(ctx context.Context, id string) (*models.ToDo, error) {
	if id == "" {
		return nil, errors.New(mustPassIdToGetOneMethod)
	}

	coll, err := t.getCollection(collectionName)
	if err != nil {
		return nil, err
	}

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, pherr.NewKnown(err, "Malformed ID", pherr.WithErrorType(pherr.InvalidArgument))
	}

	filter := bson.D{{"_id", objId}}

	retVal := &models.ToDo{}
	findErr := coll.FindOne(ctx, filter).Decode(&retVal)
	defer ctx.Done()
	if findErr != nil {
		return nil, err
	}

	return retVal, nil
}

func (t *Repo) UpdateOne(ctx context.Context, todo *models.ToDo) error {
	if todo.Id == "" {
		return errors.New(mustPassIdToUpdateOneMethod)
	}

	coll, err := t.getCollection(collectionName)
	if err != nil {
		return err
	}

	objectId, objIdErr := primitive.ObjectIDFromHex(todo.Id)
	if objIdErr != nil {
		return objIdErr
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
		return updateErr
	}

	return nil
}

func (t *Repo) DeleteOne(ctx context.Context, id string) error {
	if id == "" {
		return errors.New(mustPassIdToDeleteOneMethod)
	}

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
