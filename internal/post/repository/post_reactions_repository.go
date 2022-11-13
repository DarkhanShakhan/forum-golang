package repository

import (
	"database/sql"
	"errors"
	"forum/internal/entity"
)

const (
	POST_REACTIONS = " post_reactions"
	AND_LIKE       = " AND like = ?"
)

type PostReactionsRepository struct {
	db *sql.DB
}

func NewPostReactionsRepository(db *sql.DB) *PostReactionsRepository {
	return &PostReactionsRepository{db}
}

func (rr *PostReactionsRepository) FetchByPostId(id int, like bool) ([]entity.Reaction, error) {
	reactions := []entity.Reaction{}
	tx, err := rr.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	stmt, err := tx.Prepare("SELECT user_id, date, like FROM" + POST_REACTIONS + " WHERE post_id = ?" + AND_LIKE)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(id, like)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		reaction := entity.Reaction{}
		rows.Scan(&reaction.User.Id, &reaction.Date, &reaction.Like)
		reactions = append(reactions, reaction)
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return reactions, nil
}

func (rr *PostReactionsRepository) StoreReaction(postReaction entity.PostReaction) error {
	tx, err := rr.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	stmt, err := tx.Prepare(`INSERT INTO post_reactions(post_id, user_id, date, like) VALUES(?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	if _, err = stmt.Exec(postReaction.Post.Id, postReaction.Reaction.User.Id, postReaction.Reaction.Date, postReaction.Reaction.Like); err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (rr *PostReactionsRepository) UpdateReaction(postReaction entity.PostReaction) error {
	tx, err := rr.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	stmt, err := tx.Prepare(`UPDATE post_reactions SET like = ? WHERE post_id = ? AND user_id = ?;`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Exec(postReaction.Like, postReaction.Post.Id, postReaction.Reaction.User.Id)
	if err != nil {
		return err
	}
	rAffected, err := res.RowsAffected()
	if rAffected > 1 {
		return errors.New("more than one row has been affected")
	}
	if err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (rr *PostReactionsRepository) DeleteReaction(postReaction entity.PostReaction) error {
	tx, err := rr.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	stmt, err := tx.Prepare(`DELETE FROM post_reactions WHERE post_id = ? AND user_id = ?;`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Exec(postReaction.Post.Id, postReaction.Reaction.User.Id)
	if err != nil {
		return err
	}
	rAffected, err := res.RowsAffected()
	if rAffected > 1 {
		return errors.New("more than one row has been affected")
	}
	if err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (rr *PostReactionsRepository) FetchByUserId(userId int, like bool) ([]entity.PostReaction, error) {
	postReactions := []entity.PostReaction{}
	tx, err := rr.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	stmt, err := tx.Prepare("SELECT post_id, date, like FROM" + POST_REACTIONS + BY_USER_ID + AND_LIKE)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(userId, like)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		postReaction := entity.PostReaction{}
		rows.Scan(&postReaction.Post.Id, &postReaction.Reaction.Date, &postReaction.Reaction.Like)
		postReactions = append(postReactions, postReaction)
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return postReactions, nil
}
