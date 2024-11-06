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

var (
	layout = "2006-01-02 15:04:05"
)

type ApiService struct {
	pb.UnimplementedLenicServer
}

func (s *ApiService) CreateUser(ctx context.Context, in *pb.User) (*wrapperspb.Int32Value, error) {

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

func (s *ApiService) ActivateUser(ctx context.Context, in *wrapperspb.StringValue) (*wrapperspb.StringValue, error) {
	err := orm.Da.SetUserAsActive(in.Value)
	if err != nil {
		return nil, fmt.Errorf("could not activate user: %v", err)
	}
	return in, nil
}

func (s *ApiService) GetUser(ctx context.Context, in *wrapperspb.StringValue) (*pb.User, error) {
	u, err := orm.Da.GetUserByName(in.Value)
	if err != nil {
		return nil, fmt.Errorf("could not get user: %v", err)
	}
	user := pb.User{
		Id:            int32(u.Id),
		Username:      u.UserName,
		Email:         u.Email,
		UserFollowers: int32(u.Followers),
		UserFollowing: int32(u.Following),
		CreatedAt:     u.CreatedAt.Format(layout),
		UpdatedAt:     u.UpdatedAt.Format(layout),
		Active:        int32(u.Active),
	}
	return &user, nil
}

func (s *ApiService) SearchUsers(in *wrapperspb.StringValue, stream pb.Lenic_SearchUsersServer) error {
	users, err := orm.Da.GetSearchUsers(in.Value)
	if err != nil {
		return fmt.Errorf("could not get users: %v", err)
	}
	for _, u := range *users {
		user := pb.User{
			Id:            int32(u.Id),
			Username:      u.UserName,
			Email:         u.Email,
			UserFollowers: int32(u.Followers),
			UserFollowing: int32(u.Following),
			CreatedAt:     u.CreatedAt.Format(layout),
			UpdatedAt:     u.UpdatedAt.Format(layout),
			Active:        int32(u.Active),
		}
		err := stream.Send(&user)
		if err != nil {
			return fmt.Errorf("error sending message to stream: %v", err)
		}
	}
	return nil
}

func (s *ApiService) GetUserFollowers(in *wrapperspb.StringValue, stream pb.Lenic_GetUserFollowersServer) error {

	user, err := orm.Da.GetUserByName(in.Value)
	if err != nil {
		return fmt.Errorf("could not get user by name: %v", err)
	}

	follows, err := orm.Da.GetFollowers(user.Id)
	if err != nil {
		return fmt.Errorf("could not get users: %v", err)
	}

	for _, f := range *follows {
		u, err := orm.Da.GetUserByID(f.FollowerId)
		if err != nil {
			return fmt.Errorf("could not get Id from follower: %v", err)
		}
		u_out := pb.User{
			Id:            int32(u.Id),
			Username:      u.UserName,
			Email:         u.Email,
			UserFollowers: int32(u.Followers),
			UserFollowing: int32(u.Following),
			CreatedAt:     u.CreatedAt.Format(layout),
			UpdatedAt:     u.UpdatedAt.Format(layout),
			Active:        int32(u.Active),
		}
		err = stream.Send(&u_out)
		if err != nil {
			return fmt.Errorf("error sending message to stream: %v", err)
		}
	}
	return nil
}

func (s *ApiService) GetUserFollowing(in *wrapperspb.StringValue, stream pb.Lenic_GetUserFollowingServer) error {

	user, err := orm.Da.GetUserByName(in.Value)
	if err != nil {
		return fmt.Errorf("could not get user by name: %v", err)
	}

	follows, err := orm.Da.GetFollowing(user.Id)
	if err != nil {
		return fmt.Errorf("could not get users: %v", err)
	}

	for _, f := range *follows {
		u, err := orm.Da.GetUserByID(f.FollowedId)
		if err != nil {
			return fmt.Errorf("could not get Id from follower: %v", err)
		}
		u_out := pb.User{
			Id:            int32(u.Id),
			Username:      u.UserName,
			Email:         u.Email,
			UserFollowers: int32(u.Followers),
			UserFollowing: int32(u.Following),
			CreatedAt:     u.CreatedAt.Format(layout),
			UpdatedAt:     u.UpdatedAt.Format(layout),
			Active:        int32(u.Active),
		}
		err = stream.Send(&u_out)
		if err != nil {
			return fmt.Errorf("error sending message to stream: %v", err)
		}
	}
	return nil
}

func (s *ApiService) FollowUser(ctx context.Context, in *pb.Follow) (*wrapperspb.StringValue, error) {

	_, err := orm.Da.FollowUser(int(in.FollowerId), int(in.FollowedId))
	if err != nil {
		return nil, fmt.Errorf("error following user: %v", err)
	}

	return &wrapperspb.StringValue{Value: "OK"}, nil
}

func (s *ApiService) AcceptFollow(ctx context.Context, in *pb.Follow) (*wrapperspb.StringValue, error) {

	err := orm.Da.AcceptFollow(int(in.FollowerId), int(in.FollowedId))
	if err != nil {
		return nil, fmt.Errorf("error accepting follow: %v", err)
	}

	return &wrapperspb.StringValue{Value: "OK"}, nil
}

func (s *ApiService) UnfollowUser(ctx context.Context, in *pb.Follow) (*wrapperspb.StringValue, error) {

	err := orm.Da.UnfollowUser(int(in.FollowerId), int(in.FollowedId))
	if err != nil {
		return nil, fmt.Errorf("error unfollowing: %v", err)
	}

	return &wrapperspb.StringValue{Value: "OK"}, nil
}

func (s *ApiService) UpdateUserPass(ctx context.Context, in *pb.User) (*wrapperspb.StringValue, error) {
	err := orm.Da.SetNewPassword(in.Username, in.Pass)
	if err != nil {
		return nil, fmt.Errorf("could not update password: %v", err)
	}
	return wrapperspb.String("OK"), nil
}

func (s *ApiService) DeleteUser(ctx context.Context, in *wrapperspb.StringValue) (*wrapperspb.StringValue, error) {
	err := orm.Da.DeleteUser(in.Value)
	if err != nil {
		return nil, fmt.Errorf("could not delete user: %v", err)
	}
	return wrapperspb.String("OK"), nil

}
