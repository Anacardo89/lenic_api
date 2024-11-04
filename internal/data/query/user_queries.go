package query

const (
	InsertUser = `
	INSERT INTO users
		SET username=?,
			email=?,
			hashpass=?,
			active=?
	;`

	SelectUserById = `
	SELECT * FROM users
		WHERE id=?
	;`

	SelectUserByName = `
	SELECT * FROM users
		WHERE username = ?
	;`

	SelectSearchUsers = `
	SELECT * FROM users
		WHERE username LIKE ? COLLATE utf8mb4_general_ci;
	;`

	SelectUserByEmail = `
	SELECT * FROM users
		WHERE email = ?
	;`

	UpdateUserActive = `
	UPDATE users
		SET active=1
		WHERE username=?
	;`

	UpdatePassword = `
	UPDATE users
		SET hashpass = ?
		WHERE username = ?
	;`

	UpdateProfilePic = `
	UPDATE users
		SET profile_pic=?,
			profile_pic_ext=?
		WHERE username=?
	;`

	SelectUserFollows = `
	SELECT * FROM follows
		WHERE follower_id=? AND followed_id=?
	;`

	SelectUserFollowers = `
	SELECT * FROM follows
		WHERE followed_id=? AND follow_status=1
	;`

	SelectUserFollowing = `
	SELECT * FROM follows
		WHERE follower_id=? AND follow_status=1
	;`

	FollowUser = `
	INSERT INTO follows
		SET follower_id=?,
			followed_id=?,
			follow_status=0
	;`

	AcceptFollow = `
	UPDATE follows
		SET follow_status=1
		WHERE follower_id=? AND followed_id=?
	;`

	UnfollowUser = `
	DELETE FROM follows
		WHERE follower_id=? AND followed_id=?
	;`
)
