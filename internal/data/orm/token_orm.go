package orm

import (
	"time"

	"github.com/Anacardo89/lenic_api/internal/data/model"
	"github.com/Anacardo89/lenic_api/internal/data/query"
	"github.com/Anacardo89/lenic_api/pkg/db"
)

func (da *DataAccess) CreateToken(t *model.Token) error {
	_, err := da.Db.Exec(query.InsertToken,
		t.Token,
		t.UserId,
		t.Token,
	)
	return err
}

func (da *DataAccess) GetTokenByUserId(id int) (*model.Token, error) {
	var (
		createdAt []byte
		updatedAt []byte
	)
	t := model.Token{}
	row := da.Db.QueryRow(query.SelectTokenByUserId, id)
	err := row.Scan(
		&t.Id,
		&t.Token,
		&t.UserId,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return nil, err
	}
	t.CreatedAt, err = time.Parse(db.DateLayout, string(createdAt))
	if err != nil {
		return nil, err
	}
	t.UpdatedAt, err = time.Parse(db.DateLayout, string(updatedAt))
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (da *DataAccess) DeleteTokenByUserId(id int) error {
	_, err := da.Db.Exec(query.DeleteTokenByUserId, id)
	if err != nil {
		return err
	}
	return nil
}
