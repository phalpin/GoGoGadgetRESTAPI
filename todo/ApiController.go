package todo

import (
	"encoding/json"
	"github.com/phalpin/GoGoGadgetRESTAPI/todo/api"
	"github.com/phalpin/GoGoGadgetRESTAPI/todo/helpers"
	"github.com/phalpin/GoGoGadgetRESTAPI/todo/models"
	"net/http"
)

type ApiController struct {
	service IService
}

func NewController(svc IService) *ApiController {
	retVal := &ApiController{
		service: svc,
	}

	return retVal
}

func (c *ApiController) GetHandlers() []*api.HandlerPackage {
	return []*api.HandlerPackage{
		api.HandlerRecord("todo.create", "", c.CreateTodo, "POST"),
		api.HandlerRecord("todo.read", "/{id}", c.ReadToDo, "GET"),
		api.HandlerRecord("todo.update", "/{id}", c.UpdateTodo, "PUT"),
		api.HandlerRecord("todo.delete", "/{id}", c.DeleteToDo, "DELETE"),
	}
}

func (c *ApiController) CreateTodo(req *http.Request) (*api.ActionResult, error) {
	var obj *models.ToDo
	err := json.NewDecoder(req.Body).Decode(&obj)
	if err != nil {
		return api.ErrResult(err)
	}

	insertErr := c.service.InsertOne(req.Context(), obj)
	if insertErr != nil {
		return api.ErrResult(insertErr)
	}

	return api.ObjResult(obj)
}

func (c *ApiController) ReadToDo(req *http.Request) (*api.ActionResult, error) {
	id, err := helpers.GetRouteVariable(req, "id")
	if err != nil {
		return api.ErrResult(err)
	}

	obj, getErr := c.service.GetOne(req.Context(), id)
	if getErr != nil {
		return api.ErrResult(getErr)
	}

	return api.ObjResult(obj)
}

func (c *ApiController) UpdateTodo(req *http.Request) (*api.ActionResult, error) {
	id, err := helpers.GetRouteVariable(req, "id")
	if err != nil {
		return api.ErrResult(err)
	}

	var obj *models.ToDo
	decodeErr := json.NewDecoder(req.Body).Decode(&obj)
	if decodeErr != nil {
		return api.ErrResult(decodeErr)
	}

	obj.Id = id

	updateErr := c.service.UpdateOne(req.Context(), obj)
	if updateErr != nil {
		return api.ErrResult(updateErr)
	}

	return api.NilResult()
}

func (c *ApiController) DeleteToDo(req *http.Request) (*api.ActionResult, error) {
	id, err := helpers.GetRouteVariable(req, "id")
	if err != nil {
		return api.ErrResult(err)
	}

	deleteErr := c.service.DeleteOne(req.Context(), id)
	if deleteErr != nil {
		return api.ErrResult(deleteErr)
	}

	return api.NilResult()
}
