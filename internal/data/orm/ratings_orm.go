package orm

import (
	"database/sql"

	"github.com/Anacardo89/lenic_api/internal/data/model"
	"github.com/Anacardo89/lenic_api/internal/data/query"
)

func (da *DataAccess) GetPostUserRating(post_id int, user_id int) (*model.PostRatings, error) {
	pr := model.PostRatings{}
	row := da.Db.QueryRow(query.SelectPostUserRating, post_id, user_id)
	err := row.Scan(
		&pr.PostId,
		&pr.UserId,
		&pr.RatingValue,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return &pr, err
		} else {
			return nil, err
		}
	}
	return &pr, nil
}

func (da *DataAccess) GetCommentUserRating(comment_id int, user_id int) (*model.CommentRatings, error) {
	cr := model.CommentRatings{}
	row := da.Db.QueryRow(query.SelectCommentUserRating, comment_id, user_id)
	err := row.Scan(
		&cr.CommentId,
		&cr.UserId,
		&cr.RatingValue,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return &cr, err
		} else {
			return nil, err
		}
	}
	return &cr, nil
}
