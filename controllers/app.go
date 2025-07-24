package controllers

import (
	sql "github.com/jmoiron/sqlx"
)

type App struct {
	db *sql.DB
}

func NewApp(db *sql.DB) *App {
	return &App{
		db: db,
	}
}
