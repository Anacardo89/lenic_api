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
	"github.com/Anacardo89/lenic_api/pkg/logger"
	"github.com/google/uuid"
)

var (
	layout = "2006-01-02 15:04:05"
)

type ApiService struct {
	pb.UnimplementedLenicServer
}

func (s *ApiService) Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginResponse, error) {
	u, err := orm.Da.GetUserByName(in.Username)
	if err != nil {
		logger.Error.Println("could not get user: ", err)
		return nil, fmt.Errorf("could not get user: %v", err)
	}

	if !auth.CheckPasswordHash(in.Password, u.HashPass) {
		logger.Error.Println("invalid credentials")
		return nil, errors.New("invalid credentials")
	}

	following, err := orm.Da.GetFollowing(u.Id)
	if err != nil {
		logger.Error.Println("could not get following: ", err)
		return nil, fmt.Errorf("could not get following: %v", err)
	}

	var followers []string

	for _, f := range *following {
		user, err := orm.Da.GetUserByID(f.FollowedId)
		if err != nil {
			logger.Error.Println("could not get user following: ", err)
			return nil, fmt.Errorf("could not get user following: %v", err)
		}
		followers = append(followers, user.UserName)
	}

	token, err := auth.GenerateJWT(in.Username, followers)
	if err != nil {
		logger.Error.Println("could not create token: ", err)
		return nil, fmt.Errorf("could not create token: %v", err)
	}

	res := &pb.LoginResponse{
		Token: token,
	}

	return res, nil
}

func (s *ApiService) CreateUser(ctx context.Context, in *pb.User) (*pb.CreateUserResponse, error) {

	_, err := orm.Da.GetUserByName(in.Username)
	if err != sql.ErrNoRows {
		logger.Error.Println("user already exists")
		return nil, errors.New("user already exists")
	}

	_, err = orm.Da.GetUserByEmail(in.Email)
	if err != sql.ErrNoRows {
		logger.Error.Println("email already exists")
		return nil, errors.New("email already exists")
	}

	hashPass, err := auth.HashPassword(in.Pass)
	if err != nil {
		logger.Error.Println("could not hash password: ", err)
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
		logger.Error.Println("error creating user: ", err)
		return nil, fmt.Errorf("error creating user: %v", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		logger.Error.Println("error getting id: ", err)
		return nil, fmt.Errorf("error getting id: %v", err)
	}

	resp := &pb.CreateUserResponse{
		Id: int32(id),
	}

	return resp, nil
}

func (s *ApiService) GetUser(ctx context.Context, in *pb.GetUserRequest) (*pb.User, error) {
	u, err := orm.Da.GetUserByName(in.Username)
	if err != nil {
		logger.Error.Println("could not get user: ", err)
		return nil, fmt.Errorf("could not get user: %v", err)
	}
	user := pb.User{
		Id:            int32(u.Id),
		Username:      u.UserName,
		UserFollowers: int32(u.Followers),
		UserFollowing: int32(u.Following),
		CreatedAt:     u.CreatedAt.Format(layout),
		UpdatedAt:     u.UpdatedAt.Format(layout),
		Active:        int32(u.Active),
	}
	return &user, nil
}

func (s *ApiService) SearchUsers(in *pb.SearchUsersRequest, stream pb.Lenic_SearchUsersServer) error {
	users, err := orm.Da.GetSearchUsers(in.Username)
	if err != nil {
		logger.Error.Println("could not get users: ", err)
		return fmt.Errorf("could not get users: %v", err)
	}
	for _, u := range *users {
		user := pb.User{
			Id:            int32(u.Id),
			Username:      u.UserName,
			UserFollowers: int32(u.Followers),
			UserFollowing: int32(u.Following),
			CreatedAt:     u.CreatedAt.Format(layout),
			UpdatedAt:     u.UpdatedAt.Format(layout),
			Active:        int32(u.Active),
		}
		err := stream.Send(&user)
		if err != nil {
			logger.Error.Println("error sending message to stream:  ", err)
			return fmt.Errorf("error sending message to stream: %v", err)
		}
	}
	return nil
}

func (s *ApiService) GetUserFollowers(in *pb.GetUserFollowersRequest, stream pb.Lenic_GetUserFollowersServer) error {

	user, err := orm.Da.GetUserByName(in.Username)
	if err != nil {
		logger.Error.Println("could not get user by name: ", err)
		return fmt.Errorf("could not get user by name: %v", err)
	}

	follows, err := orm.Da.GetFollowers(user.Id)
	if err != nil {
		logger.Error.Println("could not get users: ", err)
		return fmt.Errorf("could not get users: %v", err)
	}

	for _, f := range *follows {
		u, err := orm.Da.GetUserByID(f.FollowerId)
		if err != nil {
			logger.Error.Println("could not get Id from follower: ", err)
			return fmt.Errorf("could not get Id from follower: %v", err)
		}
		u_out := pb.User{
			Id:            int32(u.Id),
			Username:      u.UserName,
			UserFollowers: int32(u.Followers),
			UserFollowing: int32(u.Following),
			CreatedAt:     u.CreatedAt.Format(layout),
			UpdatedAt:     u.UpdatedAt.Format(layout),
			Active:        int32(u.Active),
		}
		err = stream.Send(&u_out)
		if err != nil {
			logger.Error.Println("error sending message to stream: ", err)
			return fmt.Errorf("error sending message to stream: %v", err)
		}
	}
	return nil
}

func (s *ApiService) GetUserFollowing(in *pb.GetUserFollowingRequest, stream pb.Lenic_GetUserFollowingServer) error {

	user, err := orm.Da.GetUserByName(in.Username)
	if err != nil {
		logger.Error.Println("could not get user by name: ", err)
		return fmt.Errorf("could not get user by name: %v", err)
	}

	follows, err := orm.Da.GetFollowing(user.Id)
	if err != nil {
		logger.Error.Println("could not get users: ", err)
		return fmt.Errorf("could not get users: %v", err)
	}

	for _, f := range *follows {
		u, err := orm.Da.GetUserByID(f.FollowedId)
		if err != nil {
			logger.Error.Println("could not get Id from follower: ", err)
			return fmt.Errorf("could not get Id from follower: %v", err)
		}
		u_out := pb.User{
			Id:            int32(u.Id),
			Username:      u.UserName,
			UserFollowers: int32(u.Followers),
			UserFollowing: int32(u.Following),
			CreatedAt:     u.CreatedAt.Format(layout),
			UpdatedAt:     u.UpdatedAt.Format(layout),
			Active:        int32(u.Active),
		}
		err = stream.Send(&u_out)
		if err != nil {
			logger.Error.Println("error sending message to stream: ", err)
			return fmt.Errorf("error sending message to stream: %v", err)
		}
	}
	return nil
}

func (s *ApiService) FollowUser(ctx context.Context, in *pb.FollowUserRequest) (*pb.FollowUserResponse, error) {

	res := &pb.FollowUserResponse{
		Response: "NOK",
	}

	_, err := orm.Da.FollowUser(int(in.FollowerId), int(in.FollowedId))
	if err != nil {
		logger.Error.Println("error following user: ", err)
		return res, fmt.Errorf("error following user: %v", err)
	}

	res.Response = "OK"

	return res, nil
}

func (s *ApiService) AcceptFollow(ctx context.Context, in *pb.AcceptFollowRequest) (*pb.AcceptFollowResponse, error) {

	res := &pb.AcceptFollowResponse{
		Response: "NOK",
	}

	err := orm.Da.AcceptFollow(int(in.FollowerId), int(in.FollowedId))
	if err != nil {
		logger.Error.Println("error accepting follow: ", err)
		return res, fmt.Errorf("error accepting follow: %v", err)
	}

	res.Response = "OK"

	return res, nil
}

func (s *ApiService) UnfollowUser(ctx context.Context, in *pb.UnfollowRequest) (*pb.UnfollowUserResponse, error) {

	res := &pb.UnfollowUserResponse{
		Response: "NOK",
	}

	err := orm.Da.UnfollowUser(int(in.FollowerId), int(in.FollowedId))
	if err != nil {
		logger.Error.Println("error unfollowing: ", err)
		return res, fmt.Errorf("error unfollowing: %v", err)
	}

	res.Response = "OK"

	return res, nil
}

func (s *ApiService) UpdateUserPass(ctx context.Context, in *pb.User) (*pb.UpdateUserPassResponse, error) {

	res := &pb.UpdateUserPassResponse{
		Response: "NOK",
	}

	hash, err := auth.HashPassword(in.Pass)
	if err != nil {
		logger.Error.Println("could not hash password: ", err)
		return res, fmt.Errorf("could not hash password: %v", err)
	}

	err = orm.Da.SetNewPassword(in.Username, hash)
	if err != nil {
		logger.Error.Println("could not update password: ", err)
		return res, fmt.Errorf("could not update password: %v", err)
	}

	res.Response = "OK"

	return res, nil
}

func (s *ApiService) DeleteUser(ctx context.Context, in *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {

	res := &pb.DeleteUserResponse{
		Response: "NOK",
	}

	err := orm.Da.DeleteUser(in.Username)
	if err != nil {
		logger.Error.Println("could not delete user: ", err)
		return res, fmt.Errorf("could not delete user: %v", err)
	}

	res.Response = "OK"

	return res, nil
}

func (s *ApiService) StartConversation(ctx context.Context, in *pb.Conversation) (*pb.StartConversationResponse, error) {

	c := model.Conversation{
		User1Id: int(in.User1Id),
		User2Id: int(in.User2Id),
	}

	res, err := orm.Da.CreateConversation(&c)
	if err != nil {
		logger.Error.Println("could not create conversation: ", err)
		return nil, fmt.Errorf("could not create conversation: %v", err)
	}

	id64, err := res.LastInsertId()
	if err != nil {
		logger.Error.Println("could not get get conversation id: ", err)
		return nil, fmt.Errorf("could not get conversation id: %v", err)
	}

	id := int32(id64)

	resp := &pb.StartConversationResponse{
		Id: id,
	}

	return resp, nil
}

// NEEDS LOGGING

func (s *ApiService) GetUserConversations(in *pb.GetUserConversationsRequest, stream pb.Lenic_GetUserConversationsServer) error {

	u, err := orm.Da.GetUserByName(in.Username)
	if err != nil {
		return fmt.Errorf("could not get user id: %v", err)
	}

	convos, err := orm.Da.GetConversationsByUserId(u.Id)
	if err != nil {
		return fmt.Errorf("could not get user convos: %v", err)
	}

	for _, c := range convos {
		convo := pb.Conversation{
			Id:        int32(c.Id),
			User1Id:   int32(c.User1Id),
			User2Id:   int32(c.User2Id),
			CreatedAt: c.CreatedAt.Format(layout),
			UpdatedAt: c.UpdatedAt.Format(layout),
		}
		err = stream.Send(&convo)
		if err != nil {
			return fmt.Errorf("error sending message to stream: %v", err)
		}
	}
	return nil
}

func (s *ApiService) ReadConversation(ctx context.Context, in *pb.ReadConversationRequest) (*pb.ReadConversationResponse, error) {

	res := &pb.ReadConversationResponse{
		Response: "NOK",
	}

	dms, err := orm.Da.GetDMsByConversationId(int(in.Id))
	if err != nil {
		return res, fmt.Errorf("could not gt dms: %v", err)
	}

	for _, dm := range dms {
		err := orm.Da.UpdateDMReadById(dm.Id)
		if err != nil {
			return res, fmt.Errorf("could not mark dm %v as read: %v", dm.Id, err)
		}
	}

	res.Response = "OK"

	return res, nil
}

func (s *ApiService) SendDM(ctx context.Context, in *pb.DM) (*pb.SendDMResponse, error) {

	dm := model.DMessage{
		ConversationId: int(in.ConversationId),
		SenderId:       int(in.SenderId),
		Content:        in.Content,
		IsRead:         false,
	}

	res, err := orm.Da.CreateDMessage(&dm)
	if err != nil {
		return nil, fmt.Errorf("could not send DM: %v", err)
	}

	id64, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("could not get conversation id: %v", err)
	}

	id := int32(id64)

	resp := &pb.SendDMResponse{
		Id: id,
	}

	return resp, nil
}

func (s *ApiService) GetConversationDMs(in *pb.GetConversationDMsRequest, stream pb.Lenic_GetConversationDMsServer) error {
	dms, err := orm.Da.GetDMsByConversationId(int(in.Id))
	if err != nil {
		return fmt.Errorf("could not get DMs: %v", err)
	}

	for _, d := range dms {
		dm := pb.DM{
			Id:             int32(d.Id),
			ConversationId: int32(d.ConversationId),
			SenderId:       int32(d.SenderId),
			Content:        d.Content,
			IsRead:         d.IsRead,
			CreatedAt:      d.CreatedAt.Format(layout),
		}
		err = stream.Send(&dm)
		if err != nil {
			return fmt.Errorf("error sending message to stream: %v", err)
		}
	}
	return nil
}

func (s *ApiService) CreatePost(ctx context.Context, in *pb.Post) (*pb.CreatePostResponse, error) {

	guid := uuid.New().String()
	p := model.Post{
		GUID:     guid,
		AuthorId: int(in.AuthorId),
		Title:    in.Title,
		Content:  in.Content,
		Image:    "",
		ImageExt: "",
		IsPublic: in.IsPublic,
		Rating:   0,
		Active:   1,
	}

	_, err := orm.Da.CreatePost(&p)
	if err != nil {
		return nil, fmt.Errorf("could not create post: %v", err)
	}

	res := &pb.CreatePostResponse{
		Uuid: guid,
	}

	return res, nil
}

func (s *ApiService) GetPost(ctx context.Context, in *pb.GetPostRequest) (*pb.Post, error) {
	p, err := orm.Da.GetPostByGUID(in.Uuid)
	if err != nil {
		return nil, fmt.Errorf("could not get post: %v", err)
	}

	active := false
	if p.Active > 0 {
		active = true
	}

	post := pb.Post{
		Id:        int32(p.Id),
		PostGuid:  p.GUID,
		AuthorId:  int32(p.AuthorId),
		Title:     p.Title,
		Content:   p.Content,
		CreatedAt: p.CreatedAt.Format(layout),
		UpdatedAt: p.UpdatedAt.Format(layout),
		IsPublic:  p.IsPublic,
		Rating:    int32(p.Rating),
		Active:    active,
	}

	return &post, nil
}

func (s *ApiService) GetUserPosts(in *pb.GetUserPostsRequest, stream pb.Lenic_GetUserPostsServer) error {
	u, err := orm.Da.GetUserByName(in.Username)
	if err != nil {
		return fmt.Errorf("could not get user: %v", err)
	}

	posts, err := orm.Da.GetUserPosts(u.Id)
	if err != nil {
		return fmt.Errorf("could not get posts: %v", err)
	}

	for _, p := range *posts {
		active := false
		if p.Active > 0 {
			active = true
		}
		post := pb.Post{
			Id:        int32(p.Id),
			PostGuid:  p.GUID,
			AuthorId:  int32(p.AuthorId),
			Title:     p.Title,
			Content:   p.Content,
			CreatedAt: p.CreatedAt.Format(layout),
			UpdatedAt: p.UpdatedAt.Format(layout),
			IsPublic:  p.IsPublic,
			Rating:    int32(p.Rating),
			Active:    active,
		}
		err = stream.Send(&post)
		if err != nil {
			return fmt.Errorf("error sending message to stream: %v", err)
		}
	}
	return nil
}

func (s *ApiService) GetUserPublicPosts(in *pb.GetUserPublicPostsRequest, stream pb.Lenic_GetUserPostsServer) error {
	u, err := orm.Da.GetUserByName(in.Username)
	if err != nil {
		return fmt.Errorf("could not get user: %v", err)
	}

	posts, err := orm.Da.GetUserPublicPosts(u.Id)
	if err != nil {
		return fmt.Errorf("could not get posts: %v", err)
	}

	for _, p := range *posts {
		active := false
		if p.Active > 0 {
			active = true
		}
		post := pb.Post{
			Id:        int32(p.Id),
			PostGuid:  p.GUID,
			AuthorId:  int32(p.AuthorId),
			Title:     p.Title,
			Content:   p.Content,
			CreatedAt: p.CreatedAt.Format(layout),
			UpdatedAt: p.UpdatedAt.Format(layout),
			IsPublic:  p.IsPublic,
			Rating:    int32(p.Rating),
			Active:    active,
		}
		err = stream.Send(&post)
		if err != nil {
			return fmt.Errorf("error sending message to stream: %v", err)
		}
	}
	return nil
}

func (s *ApiService) GetFeed(in *pb.GetFeedRequest, stream pb.Lenic_GetFeedServer) error {
	u, err := orm.Da.GetUserByName(in.Username)
	if err != nil {
		return fmt.Errorf("could not get user: %v", err)
	}

	posts, err := orm.Da.GetFeed(u.Id)
	if err != nil {
		return fmt.Errorf("could not get posts: %v", err)
	}

	for _, p := range *posts {
		active := false
		if p.Active > 0 {
			active = true
		}
		post := pb.Post{
			Id:        int32(p.Id),
			PostGuid:  p.GUID,
			AuthorId:  int32(p.AuthorId),
			Title:     p.Title,
			Content:   p.Content,
			CreatedAt: p.CreatedAt.Format(layout),
			UpdatedAt: p.UpdatedAt.Format(layout),
			IsPublic:  p.IsPublic,
			Rating:    int32(p.Rating),
			Active:    active,
		}
		err = stream.Send(&post)
		if err != nil {
			return fmt.Errorf("error sending message to stream: %v", err)
		}
	}
	return nil
}

func (s *ApiService) RatePostUp(ctx context.Context, in *pb.PostRating) (*pb.RatePostUpResponse, error) {

	res := &pb.RatePostUpResponse{
		Response: "NOK",
	}

	err := orm.Da.RatePostUp(int(in.PostId), int(in.UserId))
	if err != nil {
		return res, fmt.Errorf("could not rate post up: %v", err)
	}

	res.Response = "OK"

	return res, nil
}

func (s *ApiService) RatePostDown(ctx context.Context, in *pb.PostRating) (*pb.RatePostDownResponse, error) {

	res := &pb.RatePostDownResponse{
		Response: "NOK",
	}

	err := orm.Da.RatePostDown(int(in.PostId), int(in.UserId))
	if err != nil {
		return res, fmt.Errorf("could not rate post down: %v", err)
	}

	res.Response = "OK"

	return res, nil
}

func (s *ApiService) UpdatePost(ctx context.Context, in *pb.Post) (*pb.UpdatePostResponse, error) {

	res := &pb.UpdatePostResponse{
		Response: "NOK",
	}

	p := model.Post{
		GUID:     in.PostGuid,
		Title:    in.Title,
		Content:  in.Content,
		IsPublic: in.IsPublic,
	}

	err := orm.Da.UpdatePost(p)
	if err != nil {
		return res, fmt.Errorf("could not update post: %v", err)
	}

	res.Response = "OK"

	return res, nil
}

func (s *ApiService) DeletePost(ctx context.Context, in *pb.DeletePostRequest) (*pb.DeletePostResponse, error) {

	res := &pb.DeletePostResponse{
		Response: "NOK",
	}

	err := orm.Da.DisablePost(in.Uuid)
	if err != nil {
		return nil, fmt.Errorf("could not delete post: %v", err)
	}

	res.Response = "OK"

	return res, nil
}

func (s *ApiService) CreateComment(ctx context.Context, in *pb.Comment) (*pb.CreateCommentResponse, error) {

	c := model.Comment{
		PostGUID: in.PostGuid,
		AuthorId: int(in.AuthorId),
		Content:  in.Content,
		Rating:   int(in.Rating),
		Active:   1,
	}

	res, err := orm.Da.CreateComment(&c)
	if err != nil {
		return nil, fmt.Errorf("could not create comment: %v", err)
	}

	id64, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("could not get conversation id: %v", err)
	}

	id := int32(id64)

	resp := &pb.CreateCommentResponse{
		Id: id,
	}

	return resp, nil
}

func (s *ApiService) GetComment(ctx context.Context, in *pb.GetCommentRequest) (*pb.Comment, error) {
	c, err := orm.Da.GetCommentById(int(in.Id))
	if err != nil {
		return nil, fmt.Errorf("could not get comment: %v", err)
	}

	active := false
	if c.Active > 0 {
		active = true
	}
	comment := pb.Comment{
		Id:        int32(c.Id),
		PostGuid:  c.PostGUID,
		AuthorId:  int32(c.AuthorId),
		Content:   c.Content,
		CreatedAt: c.CreatedAt.Format(layout),
		UpdatedAt: c.UpdatedAt.Format(layout),
		Rating:    int32(c.Rating),
		Active:    active,
	}

	return &comment, nil
}

func (s *ApiService) GetCommentsFromPost(in *pb.GetCommentsFromPostRequest, stream pb.Lenic_GetCommentsFromPostServer) error {
	comments, err := orm.Da.GetCommentsByPost(in.Uuid)
	if err != nil {
		return fmt.Errorf("could not get commentss: %v", err)
	}

	for _, c := range *comments {
		active := false
		if c.Active > 0 {
			active = true
		}
		comment := pb.Comment{
			Id:        int32(c.Id),
			PostGuid:  c.PostGUID,
			AuthorId:  int32(c.AuthorId),
			Content:   c.Content,
			CreatedAt: c.CreatedAt.Format(layout),
			UpdatedAt: c.UpdatedAt.Format(layout),
			Rating:    int32(c.Rating),
			Active:    active,
		}
		err = stream.Send(&comment)
		if err != nil {
			return fmt.Errorf("error sending message to stream: %v", err)
		}
	}
	return nil
}

func (s *ApiService) RateCommentUp(ctx context.Context, in *pb.CommentRating) (*pb.RateCommentUpResponse, error) {

	res := &pb.RateCommentUpResponse{
		Response: "NOK",
	}

	err := orm.Da.RateCommentUp(int(in.CommentId), int(in.UserId))
	if err != nil {
		return res, fmt.Errorf("could not rate comment up: %v", err)
	}

	res.Response = "OK"

	return res, nil
}

func (s *ApiService) RateCommentDown(ctx context.Context, in *pb.CommentRating) (*pb.RateCommentDownResponse, error) {

	res := &pb.RateCommentDownResponse{
		Response: "NOK",
	}

	err := orm.Da.RateCommentDown(int(in.CommentId), int(in.UserId))
	if err != nil {
		return res, fmt.Errorf("could not rate comment down: %v", err)
	}

	res.Response = "OK"

	return res, nil
}

func (s *ApiService) UpdateComment(ctx context.Context, in *pb.Comment) (*pb.UpdateCommentResponse, error) {

	res := &pb.UpdateCommentResponse{
		Response: "NOK",
	}

	err := orm.Da.UpdateCommentText(int(in.Id), in.Content)
	if err != nil {
		return res, fmt.Errorf("could not update comment: %v", err)
	}

	res.Response = "OK"

	return res, nil
}

func (s *ApiService) DeleteComment(ctx context.Context, in *pb.DeleteCommentRequest) (*pb.DeleteCommentResponse, error) {

	res := &pb.DeleteCommentResponse{
		Response: "NOK",
	}

	err := orm.Da.DisableComment(int(in.Id))
	if err != nil {
		return nil, fmt.Errorf("could not delete comment: %v", err)
	}

	res.Response = "OK"

	return res, nil
}
