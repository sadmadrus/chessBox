package store

type StoreConfig struct {
	DatabaseURL string
}

func CreateNewConfig() *StoreConfig {
	return &StoreConfig{}
}
