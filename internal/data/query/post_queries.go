package query

const (
	InsertPost = `
	INSERT INTO posts
		SET post_guid=?,
			author_id=?,
			title=?,
			content=?,
			post_image=?,
			image_ext=?,
			is_public=?,
			rating=?,
			active=?
	;`

	SelectFeed = `
	SELECT p.* FROM posts p
	LEFT JOIN follows f ON p.author_id = f.followed_id AND f.follower_id=?
	WHERE (p.is_public = TRUE AND p.active=1) OR (f.follower_id=? AND f.follow_status = 1 AND p.active=1) OR (p.author_id=? AND p.active=1)
	ORDER BY 
		CASE 
			WHEN p.created_at >= NOW() - INTERVAL 24 HOUR THEN 1 
			ELSE 2 
    	END ASC,
		p.rating DESC,
		p.created_at DESC
	;`

	SelectActivePosts = `
	SELECT * FROM posts
		WHERE active=1
		ORDER BY created_at DESC
	;`

	SelectUserActivePosts = `
	SELECT * FROM posts
		WHERE author_id=? AND active=1
		ORDER BY created_at DESC
	;`

	SelectUserPublicPosts = `
	SELECT * FROM posts
		WHERE author_id=? AND is_public=TRUE AND active=1
		ORDER BY created_at DESC
	;`

	SelectPostByGUID = `
	SELECT * FROM posts
		WHERE post_guid=?
	;`

	UpdatePost = `
	UPDATE posts
		SET title=?,
			content=?,
			is_public=?
		WHERE post_guid=?
	;`

	SetPostAsInactive = `
	UPDATE posts
		SET active=0
		WHERE post_guid=?
	;`

	RatePostUp = `
	INSERT INTO post_ratings
		SET post_id=?,
		user_id=?,
		rating_value=1
		ON DUPLICATE KEY UPDATE rating_value = CASE
        	WHEN rating_value = 1
				THEN 0
        	ELSE 1
    	END
	;`

	RatePostDown = `
	INSERT INTO post_ratings
		SET post_id=?,
		user_id=?,
		rating_value=-1
		ON DUPLICATE KEY UPDATE rating_value = CASE
        	WHEN rating_value = -1
				THEN 0
        	ELSE -1
    	END
	;`

	SelectPostUserRating = `
	SELECT * FROM post_ratings
		WHERE post_id=? AND user_id=?
	;`
)
