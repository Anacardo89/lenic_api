package query

const (
	InsertTag = `
	INSERT INTO tags
		SET tag_name=?,
			tag_type=?
	;`

	SelectTagByName = `
	SELECT * FROM tags
		WHERE tag_name=?
	;`

	InsertUserTag = `
	INSERT INTO user_tags
		SET tag_id=?,
			post_id=?,
			comment_id=?,
			tag_place=?
	;`

	SelectUserTagById = `
	SELECT * FROM user_tags
		WHERE tag_id=?
	;`

	SelectUserTagsByPostId = `
	SELECT * FROM user_tags
		WHERE post_id=?
	;`

	SelectUserTagsByCommentId = `
	SELECT * FROM user_tags
		WHERE comment_id=?
	;`

	DeleteUserTagById = `
	DELETE FROM user_tags
		WHERE id=?
	;`

	InsertReferenceTag = `
	INSERT INTO reference_tags
		SET tag_id=?,
			post_id=?,
			comment_id=?,
			tag_place=?
	;`

	SelectReferenceTagById = `
	SELECT * FROM reference_tags
		WHERE tag_id=?
	;`

	DeleteReferenceTagById = `
	DELETE FROM reference_tags
		WHERE id=?
	;`
)
