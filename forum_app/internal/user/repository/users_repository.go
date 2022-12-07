package repository

import (
	"context"
	"database/sql"
	"forum_app/internal/entity"
	"log"
	"time"
)

type UsersRepository struct {
	db       *sql.DB
	errorLog *log.Logger
}

func NewUsersRepository(db *sql.DB, errorLog *log.Logger) *UsersRepository {
	return &UsersRepository{db, errorLog}
}

func (ur *UsersRepository) FetchById(ctx context.Context, id int) (entity.User, error) {
	user := entity.User{}
	tx, err := ur.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		ur.errorLog.Println(err)
		return user, err
	}
	defer tx.Rollback()
	stmt, err := tx.PrepareContext(ctx, "SELECT id, name, email, registration_date FROM users WHERE id = ?;")
	if err != nil {
		ur.errorLog.Println(err)
		return user, err
	}
	defer stmt.Close()
	rows, err := stmt.QueryContext(ctx, id)
	if err != nil {
		ur.errorLog.Println(err)
		return user, err
	}
	if rows.Next() {
		rows.Scan(&user.Id, &user.Name, &user.Email, &user.RegDate)
	}
	if err = tx.Commit(); err != nil {
		ur.errorLog.Println(err)
		return entity.User{}, err
	}
	return user, nil
}

func (ur *UsersRepository) FetchAll(ctx context.Context) ([]entity.User, error) {
	users := []entity.User{}
	tx, err := ur.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		ur.errorLog.Println(err)
		return users, err
	}
	defer tx.Rollback()
	stmt, err := tx.PrepareContext(ctx, "SELECT id, name, email, registration_date FROM users;")
	if err != nil {
		ur.errorLog.Println(err)
		return users, err
	}
	defer stmt.Close()
	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		ur.errorLog.Println(err)
		return users, err
	}
	for rows.Next() {
		tempUser := entity.User{}
		rows.Scan(&tempUser.Id, &tempUser.Name, &tempUser.Email, &tempUser.RegDate)
		users = append(users, tempUser)
	}
	if err = tx.Commit(); err != nil {
		ur.errorLog.Println(err)
		return nil, err
	}
	return users, nil
}
func (ur *UsersRepository) FetchByEmail(ctx context.Context, email string) (entity.User, error) {
	user := entity.User{}
	tx, err := ur.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		ur.errorLog.Println(err)
		return user, err
	}
	defer tx.Rollback()
	stmt, err := tx.PrepareContext(ctx, "SELECT * FROM users WHERE email = ?")
	if err != nil {
		ur.errorLog.Println(err)
		return user, err
	}
	defer stmt.Close()
	rows, err := stmt.QueryContext(ctx, email)
	if err != nil {
		ur.errorLog.Println(err)
		return user, err
	}
	if rows.Next() {
		rows.Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.RegDate)
	}
	if err = tx.Commit(); err != nil {
		ur.errorLog.Println(err)
		return entity.User{}, err
	}
	return user, nil
}

func (ur *UsersRepository) Store(ctx context.Context, user entity.User) (int64, error) {
	tx, err := ur.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		ur.errorLog.Println(err)
		return 0, err
	}
	defer tx.Rollback()
	stmt, err := tx.PrepareContext(ctx, "INSERT INTO users(name, email, password, registration_date) VALUES (?, ?, ?, ?);")
	if err != nil {
		ur.errorLog.Println(err)
		return 0, err
	}
	defer stmt.Close()
	user.RegDate = time.Now().Format("2006-01-02")
	res, err := stmt.ExecContext(ctx, user.Name, user.Email, user.Password, user.RegDate)
	if err != nil {
		ur.errorLog.Println(err)
		return 0, err
	}
	if err = tx.Commit(); err != nil {
		ur.errorLog.Println(err)
		return 0, err
	}
	return res.LastInsertId()
}
