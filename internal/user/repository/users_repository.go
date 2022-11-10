package repository

import (
	"database/sql"
	"forum/internal/entity"
)

const (
	SELECT_QUERY = "SELECT * FROM users"
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

// FIXME:
func (ur *UsersRepository) Update(user entity.User) error {
	return nil
}

// FIXME:
func (ur *UsersRepository) Delete(id int) error {
	return nil
}
