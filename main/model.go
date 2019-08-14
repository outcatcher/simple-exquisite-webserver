package main

import (
	"database/sql"
	"errors"
	"net"
)

type entity struct {
	ID   int    `json:"id"`
	Data string `json:"Data"`
}

func (e *entity) GetEntity(db *sql.DB) error {
	return db.QueryRow("SELECT Data FROM entities WHERE id=$1", e.ID).Scan(&e.Data)
}

func (e *entity) UpdateEntity(db *sql.DB) error {
	_, err := db.Exec("UPDATE entities SET Data=$1 WHERE id=$2", e.Data, e.ID)
	return err
}

func (e *entity) DeleteEntity(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM entities WHERE id=$1", e.ID)
	return err
}

func (e *entity) CreateEntity(db *sql.DB) error {
	// postgres doesn't return the last inserted ID so this is the workaround
	err := db.QueryRow(
		"INSERT INTO entities(Data) VALUES($1) RETURNING id",
		e.Data).Scan(&e.ID)
	return err
}

func isConnectionError(err error) bool {
	switch err.(type) {
	default:
		return false
	case *net.OpError:
		return true
	}
}

func GetEntities(db *sql.DB, start, count int) ([]entity, error) {
	rows, err := db.Query("SELECT id, Data FROM entities LIMIT $1 OFFSET $2", count, start)
	if err != nil {
		if isConnectionError(err) {
			return nil, errors.New("Can't connect to database.")
		}
		return nil, err
	}

	defer rows.Close()

	var entities []entity

	for rows.Next() {
		var e entity
		if err := rows.Scan(&e.ID, &e.Data); err != nil {
			return nil, err
		}
		entities = append(entities, e)
	}

	return entities, nil
}
