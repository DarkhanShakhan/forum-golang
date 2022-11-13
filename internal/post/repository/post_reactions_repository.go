package repository

import (
	"database/sql"
	"errors"
	"fmt"
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
	rows, err := rr.db.Query("SELECT user_id, date, like FROM"+POST_REACTIONS+BY_ID+AND_LIKE, id, like)
	if err != nil {
		return reactions, err
	}
	for rows.Next() {
		reaction := entity.Reaction{}
		rows.Scan(&reaction.User.Id, &reaction.Date, &reaction.Like)
		reactions = append(reactions, reaction)
	}
	return reactions, nil
}

func (rr *PostReactionsRepository) StoreReaction(postReaction entity.PostReaction) error {
	_, err := rr.db.Exec(`INSERT INTO post_reactions(post_id, user_id, date, like) VALUES(?,?,?,?)`, postReaction.Post.Id, postReaction.Reaction.User.Id, postReaction.Reaction.Date, 1)
	fmt.Println(err)
	if err != nil {
		return err
	}
	return nil
}

func (rr *PostReactionsRepository) UpdateReaction(postReaction entity.PostReaction) error {
	res, err := rr.db.Exec(`UPDATE post_reactions SET like = ? WHERE post_id = ? AND user_id = ?;`, postReaction.Like, postReaction.Post.Id, postReaction.Reaction.User.Id)
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
	return nil
}

func (rr *PostReactionsRepository) DeleteReaction(postReaction entity.PostReaction) error {
	res, err := rr.db.Exec(`DELETE FROM post_reactions WHERE post_id = ? AND user_id = ?`, postReaction.Post.Id, postReaction.Reaction.User.Id)
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
	return nil
}

func (rr *PostReactionsRepository) FetchByUserId(userId int, like bool) ([]entity.PostReaction, error) {
	postReactions := []entity.PostReaction{}
	rows, err := rr.db.Query("SELECT post_id, date, like FROM", POST_REACTIONS, BY_USER_ID+AND_LIKE, userId, like)
	if err != nil {
		return postReactions, err
	}
	for rows.Next() {
		postReaction := entity.PostReaction{}
		rows.Scan(&postReaction.Post.Id, &postReaction.Reaction.Date, &postReaction.Reaction.Like)
		postReactions = append(postReactions, postReaction)
	}

	return postReactions, nil
}
