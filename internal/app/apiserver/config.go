package apiserver

import "github.com/sadmadrus/chessBox/internal/app/store"

type ServiceCfg struct {
	BindAddr string
	LogLevel string
	Store    *store.StoreConfig
}
