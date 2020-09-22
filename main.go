package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/phalpin/GoGoGadgetRESTAPI/middleware"
	"github.com/phalpin/GoGoGadgetRESTAPI/todo"
	"github.com/phalpin/libapi"
	"net/http"
	"os"
)

func main() {

	connStr := os.Getenv("MongoConnStr")
	dbName := "GoGoGadgetRestfulAPI"

	var Services = map[string]IServiceController{
		"todo": todo.NewController(
			todo.NewService(
				todo.NewRepo(
					connStr,
					dbName,
				))),
	}

	var Middleware = []mux.MiddlewareFunc{
		middleware.ContextStorage,
		middleware.CorrelationManager,
		middleware.MetricsCollector,
		middleware.SimpleLogging,
	}

	router := mux.NewRouter()

	//Register Services
	for key, elem := range Services {
		subRouter := router.Name(fmt.Sprintf("[Controller] %v", key)).PathPrefix(fmt.Sprintf("/%v", key)).Subrouter()
		handlers := elem.GetHandlers()
		for _, handler := range handlers {
			subRouter.Name(handler.Name).HandlerFunc(handler.ServeHTTP).Path(handler.Path).Methods(handler.Methods...)
		}
		router.Handle(fmt.Sprintf("/%v", key), subRouter)
	}

	//Register Middleware
	for _, elem := range Middleware {
		router.Use(elem)
	}

	//Handle Home
	router.HandleFunc("/", homeLink)

	//Begin
	err := http.ListenAndServe(":8080", router)

	if err != nil {
		panic(err)
	}

}

func homeLink(w http.ResponseWriter, r *http.Request) {
	var _, _ = fmt.Fprintf(w, "Hello, world!")
	var _ = r
}

type IServiceController interface {
	GetHandlers() []*libapi.HandlerPackage
}
