package query

const (
	InsertSession = `
	INSERT INTO sessions
		SET session_id=?,
			user_id=?,
			active=?
		ON DUPLICATE KEY UPDATE user_id=?, updated_at=CURRENT_TIMESTAMP
	;`

	SelectSessionById = `
	SELECT * FROM sessions
		WHERE id=?
	;`

	SelectSessionBySessionId = `
	SELECT * FROM sessions
		WHERE session_id=?
	;`
)
