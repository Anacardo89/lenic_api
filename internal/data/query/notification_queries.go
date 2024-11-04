package query

const (
	InsertNotification = `
	INSERT INTO notifications
		SET user_id=?,
			from_user_id=?,
			notif_type=?,
			notif_message=?,
			resource_id=?,
			parent_id=?
	;`

	SelectFollowNotification = `
	SELECT * FROM notifications
		WHERE user_id=? AND from_user_id=? AND notif_type='follow_request'
	;`

	SelectNotificationById = `
	SELECT * FROM notifications
		WHERE id=?
	;`

	SelectNotificationsByUser = `
	SELECT * FROM notifications
		WHERE user_id=?
			ORDER BY created_at DESC
			LIMIT ? OFFSET ?
	;`

	UpdateNotificationRead = `
	UPDATE notifications
		SET is_read=TRUE
		WHERE id=?
	;`

	DeleteNotificationByID = `
	DELETE FROM notifications
		WHERE id=?
	;`
)
