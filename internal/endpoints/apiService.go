package endpoints

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Anacardo89/lenic_api/internal/data/model"
	"github.com/Anacardo89/lenic_api/internal/data/orm"
	"github.com/Anacardo89/lenic_api/internal/pb"
	"github.com/Anacardo89/lenic_api/pkg/auth"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var (
	layout = "2006-01-02 15:04:05"
)

type ApiService struct {
	pb.UnimplementedLenicServer
}

func (s *ApiService) CreateUser(ctx context.Context, in *pb.User) (*wrapperspb.Int32Value, error) {

	_, err := orm.Da.GetUserByName(in.Username)
	if err != sql.ErrNoRows {
		return nil, fmt.Errorf("User already exists")
	}

	_, err = orm.Da.GetUserByEmail(in.Email)
	if err != sql.ErrNoRows {
		return nil, fmt.Errorf("Email already exists")
	}

	hashPass, err := auth.HashPassword(in.Pass)
	if err != nil {
		return nil, fmt.Errorf("Could not hash password", err)
	}

	createdAt, err := time.Parse(layout, string(in.CreatedAt))
	if err != nil {
		return nil, fmt.Errorf("Could not parse created at", err)
	}

	updatedAt, err := time.Parse(layout, string(in.UpdatedAt))
	if err != nil {
		return nil, fmt.Errorf("Could not parse updated at", err)
	}

	u := &model.User{
		Id:            int(in.Id),
		UserName:      string(in.Username),
		Email:         string(in.Email),
		HashPass:      hashPass,
		ProfilePic:    string(in.ProfilePic),
		ProfilePicExt: string(in.ProfilePicExt),
		Followers:     int(in.UserFollowers),
		Following:     int(in.UserFollowing),
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
		Active:        int(in.Active),
	}

	res, err := orm.Da.CreateUser(u)
	if err != nil {
		return nil, fmt.Errorf("Error creating user", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("Error getting id", err)
	}

	id32 := int32(id)

	return wrapperspb.Int32(id32), nil
}
