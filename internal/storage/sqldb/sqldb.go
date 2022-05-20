package sqldb

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"

	"github.com/nikitapozdeev/feed-repost-bot/internal/storage"
)

const (
	schemaSQL = `
		CREATE TABLE IF NOT EXISTS subscriptions(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			clientId NUMBER,
			feedLink VARCHAR(200),
			updated NUMBER
		);
	`
	selectSQL = `
		SELECT * 
			FROM subscriptions 
		 WHERE clientId = ?
	`

	insertSQL = `
		INSERT INTO subscriptions (
			clientId, feedLink
		) VALUES (
			?, ?
		)
	`

	updateSQL = `
		UPDATE subscriptions
			 SET clientId = ?,
			 		 feedLink = ?,
			 		 updated = ?
		 WHERE id = ?
	`

	deleteSQL = `
		DELETE FROM subscriptions
			WHERE id = ?
	`
)

// SqlDB implements Storage interface
type SqlDB struct {
	sql        *sql.DB
	stmtInsert *sql.Stmt
	stmtUpdate *sql.Stmt
	stmtDelete *sql.Stmt
}

// NewDB creates new database
func NewDB(dbFile string) (*SqlDB, error) {
	sqlDB, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, err
	}

	if _, err = sqlDB.Exec(schemaSQL); err != nil {
		return nil, err
	}

	stmtInsert, err := sqlDB.Prepare(insertSQL)
	if err != nil {
		return nil, err
	}

	stmtUpdate, err := sqlDB.Prepare(updateSQL)
	if err != nil {
		return nil, err
	}

	stmtDelete, err := sqlDB.Prepare(deleteSQL)
	if err != nil {
		return nil, err
	}

	return &SqlDB{
		sql:        sqlDB,
		stmtInsert: stmtInsert,
		stmtUpdate: stmtUpdate,
		stmtDelete: stmtDelete,
	}, nil
}

// Get gets all client subscriptions
func (db *SqlDB) Get(clientId int64) ([]storage.Subscription, error) {
	rows, err := db.sql.Query(selectSQL, clientId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	subscriptions := make([]storage.Subscription, 0)

	for rows.Next() {
		subscription := storage.Subscription{}
		err := rows.Scan(&subscription)
		if err != nil {
			return nil, err
		}

		subscriptions = append(subscriptions, subscription)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return subscriptions, nil
}

// Add inserts new subscription
func (db *SqlDB) Add(subscription storage.Subscription) error {
	tx, err := db.sql.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Stmt(db.stmtInsert).Exec(subscription.ClientID, subscription.FeedLink)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

// Update updates subscription
func (db *SqlDB) Update(subscription storage.Subscription) error {
	tx, err := db.sql.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Stmt(db.stmtUpdate).Exec(
		subscription.ClientID,
		subscription.FeedLink,
		subscription.Updated,
		subscription.ID,
	)

	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// Delete removes subscription
func (db *SqlDB) Delete(id int64) error {
	tx, err := db.sql.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Stmt(db.stmtDelete).Exec(id)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// Close closes database
func (db *SqlDB) Close() error {
	defer func() {
		db.stmtInsert.Close()
		db.stmtUpdate.Close()
		db.stmtDelete.Close()
		db.sql.Close()
	}()

	return nil
}
