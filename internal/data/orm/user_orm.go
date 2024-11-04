package orm

import (
	"database/sql"
	"time"

	"github.com/Anacardo89/lenic_api/internal/data/model"
	"github.com/Anacardo89/lenic_api/internal/data/query"
	"github.com/Anacardo89/lenic_api/pkg/db"
)

func (da *DataAccess) CreateUser(u *model.User) (sql.Result, error) {
	res, err := da.Db.Exec(query.InsertUser,
		u.UserName,
		u.Email,
		u.HashPass,
		u.Active)
	return res, err
}

func (da *DataAccess) GetUserByID(id int) (*model.User, error) {
	var (
		createdAt []byte
		updatedAt []byte
	)
	u := model.User{}
	row := da.Db.QueryRow(query.SelectUserById, id)
	err := row.Scan(
		&u.Id,
		&u.UserName,
		&u.Email,
		&u.HashPass,
		&u.ProfilePic,
		&u.ProfilePicExt,
		&u.Followers,
		&u.Following,
		&createdAt,
		&updatedAt,
		&u.Active)
	if err != nil {
		return nil, err
	}
	u.CreatedAt, err = time.Parse(db.DateLayout, string(createdAt))
	if err != nil {
		return nil, err
	}
	u.UpdatedAt, err = time.Parse(db.DateLayout, string(updatedAt))
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (da *DataAccess) GetUserByName(name string) (*model.User, error) {
	var (
		createdAt []byte
		updatedAt []byte
	)
	u := model.User{}
	row := da.Db.QueryRow(query.SelectUserByName, name)
	err := row.Scan(
		&u.Id,
		&u.UserName,
		&u.Email,
		&u.HashPass,
		&u.ProfilePic,
		&u.ProfilePicExt,
		&u.Followers,
		&u.Following,
		&createdAt,
		&updatedAt,
		&u.Active)
	if err != nil {
		return nil, err
	}
	u.CreatedAt, err = time.Parse(db.DateLayout, string(createdAt))
	if err != nil {
		return nil, err
	}
	u.UpdatedAt, err = time.Parse(db.DateLayout, string(updatedAt))
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (da *DataAccess) GetSearchUsers(username string) (*[]model.User, error) {
	users := []model.User{}
	likeuser := "%" + username + "%"
	rows, err := da.Db.Query(query.SelectSearchUsers, likeuser)
	if err != nil {
		if err == sql.ErrNoRows {
			return &users, nil
		}
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			createdAt []byte
			updatedAt []byte
		)
		u := model.User{}
		err = rows.Scan(
			&u.Id,
			&u.UserName,
			&u.Email,
			&u.HashPass,
			&u.ProfilePic,
			&u.ProfilePicExt,
			&u.Followers,
			&u.Following,
			&createdAt,
			&updatedAt,
			&u.Active,
		)
		if err != nil {
			return nil, err
		}
		u.CreatedAt, err = time.Parse(db.DateLayout, string(createdAt))
		if err != nil {
			return nil, err
		}
		u.UpdatedAt, err = time.Parse(db.DateLayout, string(updatedAt))
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return &users, nil
}

func (da *DataAccess) GetUserByEmail(email string) (*model.User, error) {
	var (
		createdAt []byte
		updatedAt []byte
	)
	u := model.User{}
	row := da.Db.QueryRow(query.SelectUserByEmail, email)
	err := row.Scan(
		&u.Id,
		&u.UserName,
		&u.Email,
		&u.HashPass,
		&u.ProfilePic,
		&u.ProfilePicExt,
		&u.Followers,
		&u.Following,
		&createdAt,
		&updatedAt,
		&u.Active)
	if err != nil {
		return nil, err
	}
	u.CreatedAt, err = time.Parse(db.DateLayout, string(createdAt))
	if err != nil {
		return nil, err
	}
	u.UpdatedAt, err = time.Parse(db.DateLayout, string(updatedAt))
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (da *DataAccess) SetUserAsActive(name string) error {
	_, err := da.Db.Exec(query.UpdateUserActive, name)
	if err != nil {
		return err
	}
	return nil
}

func (da *DataAccess) SetNewPassword(user string, pass string) error {
	_, err := da.Db.Exec(query.UpdatePassword, pass, user)
	if err != nil {
		return err
	}
	return nil
}

func (da *DataAccess) GetUserFollows(follower_id int, followed_id int) (*model.Follows, error) {
	f := model.Follows{}
	row := da.Db.QueryRow(query.SelectUserFollows, follower_id, followed_id)
	err := row.Scan(
		&f.FollowerId,
		&f.FollowedId,
		&f.Status)
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func (da *DataAccess) GetFollowers(followed_id int) (*[]model.Follows, error) {
	follows := []model.Follows{}
	rows, err := da.Db.Query(query.SelectUserFollowers, followed_id)
	if err != nil {
		if err == sql.ErrNoRows {
			return &follows, nil
		}
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		f := model.Follows{}
		err = rows.Scan(
			&f.FollowerId,
			&f.FollowedId,
			&f.Status,
		)
		if err != nil {
			return nil, err
		}
		if f.Status == 1 {
			follows = append(follows, f)
		}
	}
	return &follows, nil
}

func (da *DataAccess) GetFollowing(follower_id int) (*[]model.Follows, error) {
	follows := []model.Follows{}
	rows, err := da.Db.Query(query.SelectUserFollowing, follower_id)
	if err != nil {
		if err == sql.ErrNoRows {
			return &follows, nil
		}
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		f := model.Follows{}
		err = rows.Scan(
			&f.FollowerId,
			&f.FollowedId,
			&f.Status,
		)
		if err != nil {
			return nil, err
		}
		if f.Status == 1 {
			follows = append(follows, f)
		}
	}
	return &follows, nil
}

func (da *DataAccess) FollowUser(follower_id int, followed_id int) error {
	_, err := da.Db.Exec(query.FollowUser, follower_id, followed_id)
	if err != nil {
		return err
	}
	return nil
}

func (da *DataAccess) AcceptFollow(follower_id int, followed_id int) error {
	_, err := da.Db.Exec(query.AcceptFollow, follower_id, followed_id)
	if err != nil {
		return err
	}
	return nil

}

func (da *DataAccess) UnfollowUser(follower_id int, followed_id int) error {
	_, err := da.Db.Exec(query.UnfollowUser, follower_id, followed_id)
	if err != nil {
		return err
	}
	return nil
}

func (da *DataAccess) UpdateProfilePic(profile_pic string, profile_pic_ext string, username string) error {
	_, err := da.Db.Exec(query.UpdateProfilePic, profile_pic, profile_pic_ext, username)
	if err != nil {
		return err
	}
	return nil
}
