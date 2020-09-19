package todo

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/phalpin/GoGoGadgetRESTAPI/todo/helpers"
	"github.com/phalpin/GoGoGadgetRESTAPI/todo/models"
	"github.com/phalpin/GoGoGadgetRESTAPI/todo/pherr"
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

func (c *ApiController) Initialize(router *mux.Router) {
	router.Name("todo.create").Path("").HandlerFunc(c.CreateTodo).Methods("POST")
	router.Name("todo.read").Path("/{id}").HandlerFunc(c.ReadToDo).Methods("GET")
	router.Name("todo.update").Path("/{id}").HandlerFunc(c.UpdateTodo).Methods("PUT")
	router.Name("todo.delete").Path("/{id}").HandlerFunc(c.DeleteToDo).Methods("DELETE")
}

func (c *ApiController) CreateTodo(w http.ResponseWriter, req *http.Request) {
	var obj *models.ToDo
	err := json.NewDecoder(req.Body).Decode(&obj)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	insertErr := c.service.InsertOne(req.Context(), obj)
	if insertErr != nil {
		http.Error(w, insertErr.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	res, marshalErr := json.Marshal(obj)
	if marshalErr != nil {
		http.Error(w, marshalErr.Error(), http.StatusInternalServerError)
		return
	}

	var _, _ = w.Write(res)
}

func (c *ApiController) ReadToDo(w http.ResponseWriter, req *http.Request) {
	id, err := helpers.GetRouteVariable(req, "id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	obj, getErr := c.service.GetOne(req.Context(), id)
	if getErr != nil {
		casted, ok := getErr.(*pherr.KnownError)
		if ok {
			casted.WriteHttpResponse(w)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res, marshalErr := json.Marshal(obj)
	if marshalErr != nil {
		http.Error(w, marshalErr.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	var _, _ = w.Write(res)
}

func (c *ApiController) UpdateTodo(w http.ResponseWriter, req *http.Request) {
	id, err := helpers.GetRouteVariable(req, "id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var obj *models.ToDo
	decodeErr := json.NewDecoder(req.Body).Decode(&obj)
	if decodeErr != nil {
		http.Error(w, decodeErr.Error(), http.StatusBadRequest)
		return
	}

	obj.Id = id

	updateErr := c.service.UpdateOne(req.Context(), obj)
	if updateErr != nil {
		http.Error(w, updateErr.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (c *ApiController) DeleteToDo(w http.ResponseWriter, req *http.Request) {
	id, err := helpers.GetRouteVariable(req, "id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	deleteErr := c.service.DeleteOne(req.Context(), id)
	if deleteErr != nil {
		http.Error(w, deleteErr.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

}
