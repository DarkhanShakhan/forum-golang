package repository

import (
	"context"
	"database/sql"
	"forum_auth/internal/entity"
	"log"
	"time"
)

const sessionExpiry = 10 * time.Minute

type SessionsRepository struct {
	db       *sql.DB
	errorLog *log.Logger
}

func NewSessionsRepository(db *sql.DB, errorLog *log.Logger) *SessionsRepository {
	return &SessionsRepository{db, errorLog}
}

func (sr *SessionsRepository) Fetch(ctx context.Context, token string) (entity.Session, error) {
	session := entity.Session{}
	tx, err := sr.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		sr.errorLog.Println(err)
		return entity.Session{}, err
	}
	defer tx.Rollback()
	stmt, err := tx.PrepareContext(ctx, "SELECT * FROM sessions WHERE token=?")
	if err != nil {
		sr.errorLog.Println(err)
		return entity.Session{}, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, token)
	if err != nil {
		sr.errorLog.Println(err)
		return entity.Session{}, err
	}
	temp := ""
	if rows.Next() {
		rows.Scan(&session.UserId, &session.Token, &temp)
	}
	session.ExpiryTime, err = time.Parse(time.Layout, temp)
	if err != nil {
		sr.errorLog.Println(err)
	}
	if err = tx.Commit(); err != nil {
		sr.errorLog.Println(err)
		return entity.Session{}, err
	}
	return session, nil
}

func (sr *SessionsRepository) FetchByUserId(ctx context.Context, id int64) (entity.Session, error) {
	session := entity.Session{}
	tx, err := sr.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		sr.errorLog.Println(err)
		return entity.Session{}, err
	}
	defer tx.Rollback()
	stmt, err := tx.PrepareContext(ctx, "SELECT * FROM sessions WHERE user_id=?")
	if err != nil {
		sr.errorLog.Println(err)
		return entity.Session{}, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, id)
	if err != nil {
		sr.errorLog.Println(err)
		return entity.Session{}, err
	}
	temp := ""
	if rows.Next() {
		rows.Scan(&session.UserId, &session.Token, &session.ExpiryTime)
	}
	session.ExpiryTime, err = time.Parse(time.Layout, temp)
	if err != nil {
		sr.errorLog.Println(err)
	}
	if err = tx.Commit(); err != nil {
		sr.errorLog.Println(err)
		return entity.Session{}, err
	}
	return session, nil
}

func (sr *SessionsRepository) Store(ctx context.Context, session entity.Session) (entity.Session, error) {
	tx, err := sr.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		sr.errorLog.Println(err)
		return entity.Session{}, err
	}
	defer tx.Rollback()
	stmt, err := tx.PrepareContext(ctx, "INSERT OR REPLACE INTO sessions(user_id, token, expiry_date) VALUES (?, ?, ?);")
	if err != nil {
		sr.errorLog.Println(err)
		return entity.Session{}, err
	}
	defer stmt.Close()

	expiryTime := time.Now().Add(sessionExpiry).Format(time.Layout)
	_, err = stmt.ExecContext(ctx, session.UserId, session.Token, expiryTime)
	if err != nil {
		sr.errorLog.Println(err)
		return entity.Session{}, err
	}
	session.ExpiryTime, err = time.Parse(time.Layout, expiryTime)
	if err != nil {
		sr.errorLog.Println(err)
		return entity.Session{}, err
	}
	if err = tx.Commit(); err != nil {
		sr.errorLog.Println(err)
		return entity.Session{}, err
	}
	return session, err
}

func (sr *SessionsRepository) Update(ctx context.Context, session entity.Session) (entity.Session, error) {
	tx, err := sr.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		sr.errorLog.Println(err)
		return entity.Session{}, err
	}
	defer tx.Rollback()
	stmt, err := tx.PrepareContext(ctx, "UPDATE sessions SET expiry_date=? WHERE token =?;")
	if err != nil {
		sr.errorLog.Println(err)
		return entity.Session{}, err
	}
	defer stmt.Close()

	expiryTime := time.Now().Add(sessionExpiry).Format(time.Layout)
	_, err = stmt.ExecContext(ctx, expiryTime, session.Token)
	if err != nil {
		sr.errorLog.Println(err)
		return entity.Session{}, err
	}
	session.ExpiryTime, err = time.Parse(time.Layout, expiryTime)
	if err != nil {
		sr.errorLog.Println(err)
		return entity.Session{}, err
	}
	if err = tx.Commit(); err != nil {
		sr.errorLog.Println(err)
		return entity.Session{}, err
	}
	return session, nil
}

func (sr *SessionsRepository) Delete(ctx context.Context, session entity.Session) error {
	tx, err := sr.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		sr.errorLog.Println(err)
		return err
	}
	defer tx.Rollback()
	stmt, err := tx.PrepareContext(ctx, "DELETE FROM sessions WHERE token=?;")
	if err != nil {
		sr.errorLog.Println(err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, session.Token)
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
