package auth_repo

import (
	"context"
	"database/sql"
	"forum_auth/internal/entity"
	"log"
	"time"
)

type SessionsRepository struct {
	db       *sql.DB
	errorLog *log.Logger
}

func NewSessionsRepository(db *sql.DB, errorLog *log.Logger) *SessionsRepository {
	return &SessionsRepository{db, errorLog}
}

func (sr *SessionsRepository) Store(ctx context.Context, session entity.Session) error {
	tx, err := sr.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		sr.errorLog.Println(err)
		return err
	}
	defer tx.Rollback()
	stmt, err := tx.PrepareContext(ctx, "INSERT INTO users(user_id, token, expiry_date) VALUES (?, ?, ?);")
	if err != nil {
		sr.errorLog.Println(err)
		return err
	}
	defer stmt.Close()

	session.ExpiryDate = time.Now().AddDate(0, 0, 15).Format("2006-01-02")
	_, err = stmt.ExecContext(ctx, session.UserId, session.Token, session.ExpiryDate)
	if err != nil {
		sr.errorLog.Println(err)
		return err
	}
	if err = tx.Commit(); err != nil {
		sr.errorLog.Println(err)
		return err
	}
	return nil
}
