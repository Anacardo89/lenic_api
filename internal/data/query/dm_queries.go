package query

const (
	InsertConversation = `
	INSERT INTO conversations
		SET user1_id=?,
			user2_id=?
	;`

	InsertDMessage = `
	INSERT INTO dmessages
		SET conversation_id=?,
			sender_id=?,
			content=?,
			is_read=FALSE
	;`

	SelectConversationById = `
	SELECT * FROM conversations
		WHERE id=?
	;`

	SelectConversationByUserIds = `
	SELECT * FROM conversations
		WHERE user1_id=? AND user2_id=?
	;`

	SelectConversationsByUserId = `
	SELECT * FROM conversations
		WHERE user1_id=? OR user2_id=?
			ORDER BY updated_at DESC
			LIMIT ? OFFSET ?
	;`

	SelectDMById = `
	SELECT * FROM dmessages
		WHERE id=?
	;`

	SelectLastDMBySenderInConversation = `
	SELECT * FROM dmessages
		WHERE conversation_id = ? AND sender_id = ?
		ORDER BY created_at DESC
		LIMIT 1;`

	SelectDMsByConversationId = `
	SELECT * FROM dmessages
		WHERE conversation_id=?
			ORDER BY created_at
			LIMIT ? OFFSET ?
	;`

	UpdateConversationById = `
	UPDATE conversations
		SET updated_at=CURRENT_TIMESTAMP
		WHERE id=?
	;`

	UpdateDMReadById = `
	UPDATE dmessages
		SET is_read=TRUE
		WHERE id=?
	;`
)
