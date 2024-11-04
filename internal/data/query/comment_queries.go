package query

const (
	InsertComment = `
	INSERT INTO comments
		SET post_guid=?,
			author_id=?,
			content=?,
			rating=?,
			active=?
	;`

	SelectCommentById = `
	SELECT * FROM comments
		WHERE id=?
	;`

	SelectActiveCommentsByPost = `
	SELECT * FROM comments
		WHERE post_guid=? AND active=1 ORDER BY rating DESC
	;`

	UpdateCommentText = `
	UPDATE comments
		SET content=?
		WHERE id=?
	;`

	SetCommentAsInactive = `
	UPDATE comments
		SET active=0
		WHERE id=?
	;`

	RateCommentUp = `
	INSERT INTO comment_ratings
		SET comment_id=?,
		user_id=?,
		rating_value=1
		ON DUPLICATE KEY UPDATE rating_value = CASE
        	WHEN rating_value = 1
				THEN 0
        	ELSE 1
    	END
	;`

	RateCommentDown = `
	INSERT INTO comment_ratings
		SET comment_id=?,
		user_id=?,
		rating_value=-1
		ON DUPLICATE KEY UPDATE rating_value = CASE
        	WHEN rating_value = -1
				THEN 0
        	ELSE -1
    	END
	;`

	SelectCommentUserRating = `
	SELECT * FROM comment_ratings
		WHERE comment_id=? AND user_id=?
	;`
)
