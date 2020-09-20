package todo

import (
	"context"
	"fmt"
	"github.com/phalpin/GoGoGadgetRESTAPI/todo/models"
	"github.com/phalpin/liberr"
	"github.com/phalpin/liberr/errortype"
)

type IService interface {
	InsertOne(context.Context, *models.ToDo) error
	GetOne(context.Context, string) (*models.ToDo, error)
	UpdateOne(context.Context, *models.ToDo) error
	DeleteOne(context.Context, string) error
}

type Service struct {
	repo IRepo
}

func NewService(repo IRepo) *Service {
	return &Service{
		repo: repo,
	}
}

func (t *Service) InsertOne(ctx context.Context, todo *models.ToDo) error {
	if todo.Id != "" {
		return liberr.NewKnown("todo already exists", "You must provide a todo with a blank id in order to create it.", liberr.WithErrorType(errortype.InvalidArgument))
	}

	return t.repo.InsertOne(ctx, todo)
}

func (t *Service) GetOne(ctx context.Context, id string) (*models.ToDo, error) {
	if id == "" {
		return nil, liberr.NewKnown("no id passed", "You must provide an id when retrieving a saved record.", liberr.WithErrorType(errortype.InvalidArgument))
	}

	retVal, err := t.repo.GetOne(ctx, id)

	if err != nil {
		return nil, liberr.NewBaseFromError(err, liberr.WithErrorType(errortype.Unknown))
	}

	if retVal == nil {
		return nil, liberr.NewKnown("not found", fmt.Sprintf("No todo found with ID: %v", id), liberr.WithErrorType(errortype.NotFound))
	}

	return retVal, nil
}

func (t *Service) UpdateOne(ctx context.Context, todo *models.ToDo) error {
	if todo.Id == "" {
		return liberr.NewKnown("no id passed", "You must provide an id when updating a todo.", liberr.WithErrorType(errortype.InvalidArgument))
	}

	return t.repo.UpdateOne(ctx, todo)
}

func (t *Service) DeleteOne(ctx context.Context, id string) error {
	if id == "" {
		return liberr.NewKnown("no id passed", "You must provide an ID when deleting a Todo.", liberr.WithErrorType(errortype.InvalidArgument))
	}

	return t.repo.DeleteOne(ctx, id)
}
