package repository

import (
	"database/sql"
	"errors"
	"forum/internal/entity"
)

const (
	SELECT_QUERY = "SELECT * FROM users"
	UPDATE_QUERY = "UPDATE users SET"
	DELETE_QUERY = "DELETE FROM users"
	NAME         = " name = "
	PASSWORD     = " password = "
	BY_ID        = " WHERE id = ?"
)

type UsersRepository struct {
	db *sql.DB
}

func NewUsersRepository(db *sql.DB) *UsersRepository {
	return &UsersRepository{db}
}

func (ur *UsersRepository) FetchById(id int) (entity.User, error) {
	user := entity.User{}
	rows, err := ur.db.Query(SELECT_QUERY+BY_ID, id)
	if err != nil {
		return user, err
	}
	for rows.Next() {
		rows.Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.RegDate)
	}
	return user, nil
}

func (ur *UsersRepository) FetchAll() ([]entity.User, error) {
	users := []entity.User{}
	rows, err := ur.db.Query(SELECT_QUERY)
	if err != nil {
		return users, err
	}
	tempUser := entity.User{}
	for rows.Next() {
		rows.Scan(&tempUser.Id, &tempUser.Name, &tempUser.Email, &tempUser.Password, &tempUser.RegDate)
	}
	return users, nil
}

func (ur *UsersRepository) Update(user entity.User) error {
	query := UPDATE_QUERY
	if user.Id == 0 {
		return errors.New("user id not provided")
	}
	if user.Name != "" {
		query += NAME + user.Name
	}
	if user.Password != "" {
		query += PASSWORD + user.Password
	}
	if query == UPDATE_QUERY {
		return errors.New("no attributes to update")
	}
	query += BY_ID
	result, err := ur.db.Exec(query, user.Id)
	if err != nil {
		return err
	}
	nbr, err := result.RowsAffected()
	if nbr > 1 {
		return errors.New("more than one row has been affected")
	}
	if err != nil {
		return err
	}
	return nil
}

func (ur *UsersRepository) Delete(id int) error {
	query := DELETE_QUERY
	query += BY_ID
	result, err := ur.db.Exec(query, id)
	if err != nil {
		return err
	}
	nbr, err := result.RowsAffected()
	if nbr > 1 {
		return errors.New("more than one row has been affected")
	}
	if err != nil {
		return err
	}
	return nil
}
