package orm

import (
	"database/sql"
	"time"

	"github.com/Anacardo89/lenic_api/internal/data/model"
	"github.com/Anacardo89/lenic_api/internal/data/query"
	"github.com/Anacardo89/lenic_api/pkg/db"
)

func (da *DataAccess) CreateNotification(n *model.Notification) (sql.Result, error) {
	result, err := da.Db.Exec(query.InsertNotification,
		n.UserID,
		n.FromUserId,
		n.NotifType,
		n.NotifMsg,
		n.ResourceId,
		n.ParentId)
	return result, err
}

func (da *DataAccess) GetFollowNotification(user_id int, from_user_id int) (*model.Notification, error) {
	var (
		createdAt []byte
		updatedAt []byte
	)
	n := model.Notification{}
	row := da.Db.QueryRow(query.SelectFollowNotification, user_id, from_user_id)
	err := row.Scan(
		&n.Id,
		&n.UserID,
		&n.FromUserId,
		&n.NotifType,
		&n.NotifMsg,
		&n.ResourceId,
		&n.ParentId,
		&n.IsRead,
		&createdAt,
		&updatedAt)
	if err != nil {
		return nil, err
	}
	n.CreatedAt, err = time.Parse(db.DateLayout, string(createdAt))
	if err != nil {
		return nil, err
	}
	n.UpdatedAt, err = time.Parse(db.DateLayout, string(updatedAt))
	if err != nil {
		return nil, err
	}

	return &n, nil
}

func (da *DataAccess) GetNotificationById(id int) (*model.Notification, error) {
	var (
		createdAt []byte
		updatedAt []byte
	)
	n := model.Notification{}
	row := da.Db.QueryRow(query.SelectNotificationById, id)
	err := row.Scan(
		&n.Id,
		&n.UserID,
		&n.FromUserId,
		&n.NotifType,
		&n.NotifMsg,
		&n.ResourceId,
		&n.ParentId,
		&n.IsRead,
		&createdAt,
		&updatedAt)
	if err != nil {
		return nil, err
	}
	n.CreatedAt, err = time.Parse(db.DateLayout, string(createdAt))
	if err != nil {
		return nil, err
	}
	n.UpdatedAt, err = time.Parse(db.DateLayout, string(updatedAt))
	if err != nil {
		return nil, err
	}

	return &n, nil
}

func (da *DataAccess) GetNotificationsByUser(user_id int, limit int, offset int) ([]*model.Notification, error) {
	notifs := []*model.Notification{}
	rows, err := da.Db.Query(query.SelectNotificationsByUser, user_id, limit, offset)
	if err != nil {
		if err == sql.ErrNoRows {
			return notifs, nil
		}
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			createdAt []byte
			updatedAt []byte
		)
		n := model.Notification{}
		err = rows.Scan(
			&n.Id,
			&n.UserID,
			&n.FromUserId,
			&n.NotifType,
			&n.NotifMsg,
			&n.ResourceId,
			&n.ParentId,
			&n.IsRead,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, err
		}
		n.CreatedAt, err = time.Parse(db.DateLayout, string(createdAt))
		if err != nil {
			return nil, err
		}
		n.UpdatedAt, err = time.Parse(db.DateLayout, string(updatedAt))
		if err != nil {
			return nil, err
		}
		notifs = append(notifs, &n)
	}
	return notifs, nil
}

func (da *DataAccess) UpdateNotificationRead(id int) error {
	_, err := da.Db.Exec(query.UpdateNotificationRead, id)
	if err != nil {
		return err
	}
	return nil
}

func (da *DataAccess) DeleteNotificationByID(id int) error {
	_, err := da.Db.Exec(query.DeleteNotificationByID, id)
	if err != nil {
		return err
	}
	return nil
}
