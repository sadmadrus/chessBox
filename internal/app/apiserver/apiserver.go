package apiserver

import (
	"github.com/gorilla/mux"
	"github.com/sadmadrus/chessBox/internal/app/store"
	"github.com/sirupsen/logrus"
)

type APIServer struct {
	config *ServiceCfg
	logger *logrus.Logger
	router *mux.Router
	store  *store.Store
}

func CreateNewAPIServer(config *ServiceCfg) *APIServer {
	return &APIServer{
		config: config,
		logger: logrus.New(),
		router: mux.NewRouter(),
	}
}
