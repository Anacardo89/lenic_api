package orm

import (
	"database/sql"
	"time"

	"github.com/Anacardo89/lenic_api/internal/data/model"
	"github.com/Anacardo89/lenic_api/internal/data/query"
	"github.com/Anacardo89/lenic_api/pkg/db"
)

func (da *DataAccess) CreateConversation(c *model.Conversation) (sql.Result, error) {
	result, err := da.Db.Exec(query.InsertConversation,
		c.User1Id,
		c.User2Id,
	)
	return result, err
}

func (da *DataAccess) CreateDMessage(d *model.DMessage) (sql.Result, error) {
	result, err := da.Db.Exec(query.InsertDMessage,
		d.ConversationId,
		d.SenderId,
		d.Content,
	)
	return result, err
}

func (da *DataAccess) GetConversationById(id int) (*model.Conversation, error) {
	var (
		createdAt []byte
		updatedAt []byte
	)
	c := model.Conversation{}
	row := da.Db.QueryRow(query.SelectConversationById, id)
	err := row.Scan(
		&c.Id,
		&c.User1Id,
		&c.User2Id,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return nil, err
	}
	c.CreatedAt, err = time.Parse(db.DateLayout, string(createdAt))
	if err != nil {
		return nil, err
	}
	c.UpdatedAt, err = time.Parse(db.DateLayout, string(updatedAt))
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (da *DataAccess) GetConversationByUserIds(user1_id int, user2_id int) (*model.Conversation, error) {
	min := user1_id
	max := user2_id
	if user1_id > user2_id {
		min = user2_id
		max = user1_id
	}
	var (
		createdAt []byte
		updatedAt []byte
	)
	c := model.Conversation{}
	row := da.Db.QueryRow(query.SelectConversationByUserIds, min, max)
	err := row.Scan(
		&c.Id,
		&c.User1Id,
		&c.User2Id,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return nil, err
	}
	c.CreatedAt, err = time.Parse(db.DateLayout, string(createdAt))
	if err != nil {
		return nil, err
	}
	c.UpdatedAt, err = time.Parse(db.DateLayout, string(updatedAt))
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (da *DataAccess) GetConversationsByUserId(user_id int, limit int, offset int) ([]*model.Conversation, error) {
	conversations := []*model.Conversation{}
	rows, err := da.Db.Query(query.SelectConversationsByUserId, user_id, user_id, limit, offset)
	if err != nil {
		if err == sql.ErrNoRows {
			return conversations, nil
		}
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			createdAt []byte
			updatedAt []byte
		)
		c := model.Conversation{}
		err = rows.Scan(
			&c.Id,
			&c.User1Id,
			&c.User2Id,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, err
		}
		c.CreatedAt, err = time.Parse(db.DateLayout, string(createdAt))
		if err != nil {
			return nil, err
		}
		c.UpdatedAt, err = time.Parse(db.DateLayout, string(updatedAt))
		if err != nil {
			return nil, err
		}
		conversations = append(conversations, &c)
	}
	return conversations, nil
}

func (da *DataAccess) GetDMById(id int) (*model.DMessage, error) {
	var createdAt []byte
	m := model.DMessage{}
	row := da.Db.QueryRow(query.SelectDMById, id)
	err := row.Scan(
		&m.Id,
		&m.ConversationId,
		&m.SenderId,
		&m.Content,
		&m.IsRead,
		&createdAt,
	)
	if err != nil {
		return nil, err
	}
	m.CreatedAt, err = time.Parse(db.DateLayout, string(createdAt))
	if err != nil {
		return nil, err
	}

	return &m, nil
}

func (da *DataAccess) GetLastDMBySenderInConversation(converrsation_id int, sender_id int) (*model.DMessage, error) {
	var createdAt []byte
	m := model.DMessage{}
	row := da.Db.QueryRow(query.SelectLastDMBySenderInConversation, converrsation_id, sender_id)
	err := row.Scan(
		&m.Id,
		&m.ConversationId,
		&m.SenderId,
		&m.Content,
		&m.IsRead,
		&createdAt,
	)
	if err != nil {
		return nil, err
	}
	m.CreatedAt, err = time.Parse(db.DateLayout, string(createdAt))
	if err != nil {
		return nil, err
	}

	return &m, nil
}

func (da *DataAccess) GetDMsByConversationId(conversation_id int, limit int, offset int) ([]*model.DMessage, error) {
	dms := []*model.DMessage{}
	rows, err := da.Db.Query(query.SelectDMsByConversationId, conversation_id, limit, offset)
	if err != nil {
		if err == sql.ErrNoRows {
			return dms, nil
		}
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var createdAt []byte
		m := model.DMessage{}
		err = rows.Scan(
			&m.Id,
			&m.ConversationId,
			&m.SenderId,
			&m.Content,
			&m.IsRead,
			&createdAt,
		)
		if err != nil {
			return nil, err
		}
		m.CreatedAt, err = time.Parse(db.DateLayout, string(createdAt))
		if err != nil {
			return nil, err
		}
		dms = append(dms, &m)
	}
	return dms, nil
}

func (da *DataAccess) UpdateConversationById(id int) error {
	_, err := da.Db.Exec(query.UpdateConversationById, id)
	if err != nil {
		return err
	}
	return nil
}

func (da *DataAccess) UpdateDMReadById(id int) error {
	_, err := da.Db.Exec(query.UpdateDMReadById, id)
	if err != nil {
		return err
	}
	return nil
}
