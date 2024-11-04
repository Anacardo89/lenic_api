package orm

import (
	"github.com/Anacardo89/lenic_api/internal/data/model"
	"github.com/Anacardo89/lenic_api/internal/data/query"
)

func (da *DataAccess) CreateTag(t *model.Tag) error {
	_, err := da.Db.Exec(query.InsertTag,
		t.TagName,
		t.TagType,
	)
	return err
}

func (da *DataAccess) GetTagByName(tag_name string) (*model.Tag, error) {
	t := model.Tag{}
	row := da.Db.QueryRow(query.SelectTagByName, tag_name)
	err := row.Scan(
		&t.Id,
		&t.TagName,
		&t.TagType,
	)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (da *DataAccess) CreateUserTag(t *model.UserTag) error {
	_, err := da.Db.Exec(query.InsertUserTag,
		t.TagId,
		t.PostId,
		t.CommentId,
		t.TagPlace,
	)
	return err
}

func (da *DataAccess) GetUserTagById(tag_id int) (*model.UserTag, error) {
	t := model.UserTag{}
	row := da.Db.QueryRow(query.SelectUserTagById, tag_id)
	err := row.Scan(
		&t.Id,
		&t.TagId,
		&t.PostId,
		&t.CommentId,
		&t.TagPlace,
	)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (da *DataAccess) DeleteUserTagByID(id int) error {
	_, err := da.Db.Exec(query.DeleteUserTagById, id)
	return err
}

func (da *DataAccess) CreateReferenceTag(t *model.ReferenceTag) error {
	_, err := da.Db.Exec(query.InsertUserTag,
		t.TagId,
		t.PostId,
		t.CommentId,
		t.TagPlace,
	)
	return err
}

func (da *DataAccess) GetReferenceTagById(tag_id int) (*model.ReferenceTag, error) {
	t := model.ReferenceTag{}
	row := da.Db.QueryRow(query.SelectReferenceTagById, tag_id)
	err := row.Scan(
		&t.Id,
		&t.TagId,
		&t.PostId,
		&t.CommentId,
		&t.TagPlace,
	)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (da *DataAccess) DeleteReferenceTagByID(id int) error {
	_, err := da.Db.Exec(query.DeleteReferenceTagById, id)
	return err
}
