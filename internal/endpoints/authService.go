package endpoints

import (
	"context"
	"errors"
	"fmt"

	"github.com/Anacardo89/lenic_api/internal/data/orm"
	"github.com/Anacardo89/lenic_api/internal/pb"
	"github.com/Anacardo89/lenic_api/pkg/auth"
)

type AuthService struct {
	pb.UnimplementedAuthServiceServer
}

func (s *AuthService) Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginResponse, error) {
	u, err := orm.Da.GetUserByName(in.Username)
	if err != nil {
		return nil, fmt.Errorf("could not get user: %v", err)
	}

	if !auth.CheckPasswordHash(in.Password, u.HashPass) {
		return nil, errors.New("invalid credentials")
	}

	following, err := orm.Da.GetFollowing(u.Id)

	var followers []string

	for _, f := range *following {
		user, err := orm.Da.GetUserByID(f.FollowedId)
		if err != nil {
			return nil, fmt.Errorf("could not get user following: %v", err)
		}
		followers = append(followers, user.UserName)
	}

	token, err := auth.GenerateJWT(in.Username, followers)
	if err != nil {
		return nil, fmt.Errorf("could not create token: %v", err)
	}

	res := &pb.LoginResponse{
		Token: token,
	}

	return res, nil
}
