package httprouter

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/schollz/httpfileserver"
	"github.com/sjsafranek/lemur/middleware"
	"github.com/sjsafranek/logger"
)

func Vars(r *http.Request) map[string]string {
	return mux.Vars(r)
}

type Router struct {
	*mux.Router
}

func (self *Router) AttachFileServer(route, directory string) {
	self.PathPrefix("/static/").Handler(httpfileserver.New("/static/", directory))
}

func (self *Router) AttachHandlerFunc(pattern string, handler http.HandlerFunc, methods []string) {
	logger.Info("Attaching HTTP handler for route: ", methods, " ", pattern)
	self.Methods(methods...).Path(pattern).Handler(handler)
}

func New() Router {
	// create http router
	router := mux.NewRouter().StrictSlash(true)

	// attach middleware
	router.Use(middleware.LoggingMiddleWare, middleware.SetHeadersMiddleWare, middleware.CORSMiddleWare)
	handlers.CompressHandler(router)

	return Router{router}
}
