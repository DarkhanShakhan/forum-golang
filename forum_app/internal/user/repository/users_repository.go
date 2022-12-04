package repository

import (
	"context"
	"database/sql"
	"forum_app/internal/entity"
)

type UsersRepository struct {
	db *sql.DB
}

func NewUsersRepository(db *sql.DB) *UsersRepository {
	return &UsersRepository{db}
}

func (ur *UsersRepository) FetchById(ctx context.Context, id int) (entity.User, error) {
	user := entity.User{}
	tx, err := ur.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return user, err
	}
	defer tx.Rollback()
	stmt, err := tx.PrepareContext(ctx, "SELECT * FROM users WHERE id = ?;")
	if err != nil {
		return user, err
	}
	defer stmt.Close()
	rows, err := stmt.QueryContext(ctx, id)
	if err != nil {
		return user, err
	}
	if rows.Next() {
		rows.Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.RegDate)
	}
	if err = tx.Commit(); err != nil {
		return entity.User{}, err
	}
	return user, nil
}

func (ur *UsersRepository) FetchAll(ctx context.Context) ([]entity.User, error) {
	users := []entity.User{}
	tx, err := ur.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return users, err
	}
	defer tx.Rollback()
	stmt, err := tx.PrepareContext(ctx, "SELECT * FROM users;")
	if err != nil {
		return users, err
	}
	defer stmt.Close()
	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return users, err
	}
	for rows.Next() {
		tempUser := entity.User{}
		rows.Scan(&tempUser.Id, &tempUser.Name, &tempUser.Email, &tempUser.Password, &tempUser.RegDate)
		users = append(users, tempUser)
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return users, nil
}
func (ur *UsersRepository) FetchByEmail(ctx context.Context, email string) (entity.User, error) {
	user := entity.User{}
	tx, err := ur.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return user, err
	}
	defer tx.Rollback()
	stmt, err := tx.PrepareContext(ctx, "SELECT * FROM users WHERE email = ?")
	if err != nil {
		return user, err
	}
	defer stmt.Close()
	rows, err := stmt.QueryContext(ctx, email)
	if err != nil {
		return user, err
	}
	if rows.Next() {
		rows.Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.RegDate)
	}
	if err = tx.Commit(); err != nil {
		return entity.User{}, err
	}
	return user, nil
}

//for future use
// func (ur *UsersRepository) Update(user entity.User) error {
// 	query := UPDATE_QUERY
// 	if user.Id == 0 {
// 		return errors.New("user id not provided")
// 	}
// 	if user.Name != "" {
// 		query += NAME + `"` + user.Name + `"`
// 	}
// 	if user.Password != "" {
// 		if user.Name != "" {
// 			query += `,`
// 		}
// 		query += PASSWORD + `"` + user.Password + `"`
// 	}
// 	if query == UPDATE_QUERY {
// 		return errors.New("no attributes to update")
// 	}
// 	query += BY_ID
// 	fmt.Println(query)
// 	result, err := ur.db.Exec(query, user.Id)
// 	if err != nil {
// 		return err
// 	}
// 	nbr, err := result.RowsAffected()
// 	if nbr > 1 {
// 		return errors.New("more than one row has been affected")
// 	}
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (ur *UsersRepository) Delete(id int) error {
// 	result, err := ur.db.Exec(DELETE_QUERY+BY_ID, id)
// 	if err != nil {
// 		return err
// 	}
// 	nbr, err := result.RowsAffected()
// 	if nbr > 1 {
// 		return errors.New("more than one row has been affected")
// 	}
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
