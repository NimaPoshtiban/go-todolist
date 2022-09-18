package models

import "database/sql"

type Models struct {
	DB DBModel
}

type DBModel struct {
	DB *sql.DB
}

func NewModels(db *sql.DB) Models {
	return Models{DB: DBModel{
		DB: db,
	}}
}
