package endpoints

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Anacardo89/lenic_api/internal/data/model"
	"github.com/Anacardo89/lenic_api/internal/data/orm"
	"github.com/Anacardo89/lenic_api/internal/pb"
	"github.com/Anacardo89/lenic_api/pkg/auth"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type CreateUserService struct {
	pb.UnimplementedCreateUserServiceServer
}

func (s *CreateUserService) CreateUser(ctx context.Context, in *pb.User) (*wrapperspb.Int32Value, error) {

	_, err := orm.Da.GetUserByName(in.Username)
	if err != sql.ErrNoRows {
		return nil, errors.New("user already exists")
	}

	_, err = orm.Da.GetUserByEmail(in.Email)
	if err != sql.ErrNoRows {
		return nil, errors.New("email already exists")
	}

	hashPass, err := auth.HashPassword(in.Pass)
	if err != nil {
		return nil, fmt.Errorf("could not hash password: %v", err)
	}

	u := &model.User{
		UserName: string(in.Username),
		Email:    string(in.Email),
		HashPass: hashPass,
		Active:   0,
	}

	res, err := orm.Da.CreateUser(u)
	if err != nil {
		return nil, fmt.Errorf("error creating user: %v", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error getting id: %v", err)
	}

	id32 := int32(id)

	return wrapperspb.Int32(id32), nil
}
