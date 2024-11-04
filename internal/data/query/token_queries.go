package query

const (
	InsertToken = `
	INSERT INTO tokens
		SET token=?,
			user_id=?
		ON DUPLICATE KEY UPDATE token=?, updated_at=CURRENT_TIMESTAMP
	;`

	SelectTokenByUserId = `
	SELECT * FROM tokens
		WHERE user_id=?
	;`

	DeleteTokenByUserId = `
	SELECT * FROM tokens
		WHERE user_id=?
	;`
)
