package main

import "go.uber.org/zap"

type application struct {
	config config
	// store
	logger *zap.SugaredLogger
}

type config struct {
	addr string
	db dbConfig
	env string
	apiURL string
	auth authConfig
}

type dbConfig struct {
	host string
	port string
	name string
}

type authConfig struct {

}
























