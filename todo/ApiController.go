package todo

import (
	"encoding/json"
	"github.com/phalpin/GoGoGadgetRESTAPI/todo/models"
	"github.com/phalpin/libapi"
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

func (c *ApiController) GetHandlers() []*libapi.HandlerPackage {
	return []*libapi.HandlerPackage{
		libapi.HandlerRecord("todo.create", "", c.CreateTodo, "POST"),
		libapi.HandlerRecord("todo.read", "/{id}", c.ReadToDo, "GET"),
		libapi.HandlerRecord("todo.update", "/{id}", c.UpdateTodo, "PUT"),
		libapi.HandlerRecord("todo.delete", "/{id}", c.DeleteToDo, "DELETE"),
	}
}

func (c *ApiController) CreateTodo(req *http.Request) (*libapi.ActionResult, error) {
	var obj *models.ToDo
	err := json.NewDecoder(req.Body).Decode(&obj)
	if err != nil {
		return libapi.ErrResult(err)
	}

	insertErr := c.service.InsertOne(req.Context(), obj)
	if insertErr != nil {
		return libapi.ErrResult(insertErr)
	}

	return libapi.ObjResult(obj)
}

func (c *ApiController) ReadToDo(req *http.Request) (*libapi.ActionResult, error) {
	id, err := libapi.GetRouteVariable(req, "id")
	if err != nil {
		return libapi.ErrResult(err)
	}

	obj, getErr := c.service.GetOne(req.Context(), id)
	if getErr != nil {
		return libapi.ErrResult(getErr)
	}

	return libapi.ObjResult(obj)
}

func (c *ApiController) UpdateTodo(req *http.Request) (*libapi.ActionResult, error) {
	id, err := libapi.GetRouteVariable(req, "id")
	if err != nil {
		return libapi.ErrResult(err)
	}

	var obj *models.ToDo
	decodeErr := json.NewDecoder(req.Body).Decode(&obj)
	if decodeErr != nil {
		return libapi.ErrResult(decodeErr)
	}

	obj.Id = id

	updateErr := c.service.UpdateOne(req.Context(), obj)
	if updateErr != nil {
		return libapi.ErrResult(updateErr)
	}

	return libapi.NilResult()
}

func (c *ApiController) DeleteToDo(req *http.Request) (*libapi.ActionResult, error) {
	id, err := libapi.GetRouteVariable(req, "id")
	if err != nil {
		return libapi.ErrResult(err)
	}

	deleteErr := c.service.DeleteOne(req.Context(), id)
	if deleteErr != nil {
		return libapi.ErrResult(deleteErr)
	}

	return libapi.NilResult()
}
