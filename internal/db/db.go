package db

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/nikitapozdeev/feed-repost-bot/internal/model"
)

var schemaSQL = `
CREATE TABLE IF NOT EXISTS subscriptions (
	clientId NUMBER,
	feedLink VARCHAR(100)
);
`

var insertSQL = `
	INSERT INTO subscriptions (
		clientId, feedLink
	) VALUES (
		?, ?
	)
`

type DB struct {
	sql *sql.DB
	stmt *sql.Stmt
	buffer []model.Subscription
}

func NewDB(dbFile string) (*DB, error) {
	sqlDB, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, err
	}

	if _, err = sqlDB.Exec(schemaSQL); err != nil {
		return nil, err
	}

	stmt, err := sqlDB.Prepare(insertSQL)
	if err != nil {
		return nil, err
	}

	return &DB{
		sql: sqlDB,
		stmt: stmt,
		buffer: make([]model.Subscription, 0, 1024),
	}, nil
}

func (db *DB) Add(subscription model.Subscription) error {
	if len(db.buffer) == cap(db.buffer) {
		return errors.New("Subscription buffer is full")
	}

	db.buffer = append(db.buffer, subscription)
	if len(db.buffer) == cap(db.buffer) {
		if err := db.Flush(); err != nil {
			return fmt.Errorf("Unable to flush subscriptions: %w", err)
		}
	}

	return nil
}

func (db *DB) Flush() error {
	tx, err := db.sql.Begin()
	if err != nil {
		return err
	}

	for _, subscription := range db.buffer {
		_, err := tx.Stmt(db.stmt).Exec(subscription.ClientID, subscription.FeedLink)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	db.buffer = db.buffer[:0]
	return tx.Commit()
}

func (db *DB) Close() error {
	defer func() {
		db.stmt.Close()
		db.sql.Close()
	}()

	if err := db.Flush(); err != nil {
		return err
	}

	return nil
}