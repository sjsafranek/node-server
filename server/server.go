package server

import (
	"fmt"
	"net/http"

	"github.com/sjsafranek/logger"
	"github.com/sjsafranek/node-server/config"
	"github.com/sjsafranek/node-server/database"
	"github.com/sjsafranek/node-server/httprouter"
	"github.com/sjsafranek/node-server/models"
)

func New(conf config.Configuration) (*Server, error) {
	db, err := database.New(conf.Database)
	if nil != err {
		return &Server{}, err
	}

	router := httprouter.New()
	server := &Server{db: db, config: conf, router: router}
	server.Init()
	return server, err
}

type Server struct {
	db     models.Database
	config config.Configuration
	router httprouter.Router
}

func (self *Server) GetDatabase() models.Database {
	return self.db
}

func (self *Server) Init() {
	self.router.AttachHandlerFunc("/api/v1/{bucketId}/get/{key}", self.HttpGetKeyHandler, []string{"GET"})
	self.router.AttachHandlerFunc("/api/v1/{bucketId}/set/{key}", self.HttpSetKeyHandler, []string{"POST"})
}

func (self *Server) ListenAndServe(address string) error {
	logger.Info(fmt.Sprintf("Magic happens on port %v...", address))
	return http.ListenAndServe(address, self.router)
}
