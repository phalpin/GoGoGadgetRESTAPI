package todo

import (
	"context"
	"errors"
	"fmt"
	"github.com/phalpin/GoGoGadgetRESTAPI/todo/models"
	"github.com/phalpin/GoGoGadgetRESTAPI/todo/pherr"
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
	return t.repo.InsertOne(ctx, todo)
}

func (t *Service) GetOne(ctx context.Context, id string) (*models.ToDo, error) {
	retVal, err := t.repo.GetOne(ctx, id)

	if err != nil {
		return nil, err
	}

	if retVal == nil {
		return nil, pherr.NewKnown(errors.New("not found"), fmt.Sprintf("No todo found with ID: %v", id), pherr.WithErrorType(pherr.NotFound))
	}

	return retVal, nil
}

func (t *Service) UpdateOne(ctx context.Context, todo *models.ToDo) error {
	return t.repo.UpdateOne(ctx, todo)
}

func (t *Service) DeleteOne(ctx context.Context, id string) error {
	return t.repo.DeleteOne(ctx, id)
}
