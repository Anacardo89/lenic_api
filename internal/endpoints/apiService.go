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
	"github.com/google/uuid"
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

func (s *ApiService) StartConversation(ctx context.Context, in *pb.Conversation) (*wrapperspb.Int32Value, error) {
	c := model.Conversation{
		User1Id: int(in.User1Id),
		User2Id: int(in.User2Id),
	}
	res, err := orm.Da.CreateConversation(&c)
	if err != nil {
		return nil, fmt.Errorf("could not create conversation: %v", err)
	}

	id64, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("could not get conversation id: %v", err)
	}

	id := int32(id64)

	return wrapperspb.Int32(id), nil
}

func (s *ApiService) GetUserConversations(in *wrapperspb.StringValue, stream pb.Lenic_GetUserConversationsServer) error {
	u, err := orm.Da.GetUserByName(in.Value)
	if err != nil {
		return fmt.Errorf("could not get user id: %v", err)
	}

	convos, err := orm.Da.GetConversationsByUserId(u.Id)
	if err != nil {
		fmt.Errorf("could not get user convos: %v", err)
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

func (s *ApiService) ReadConversation(ctx context.Context, in *wrapperspb.Int32Value) (*wrapperspb.StringValue, error) {
	dms, err := orm.Da.GetDMsByConversationId(int(in.Value))
	if err != nil {
		return nil, fmt.Errorf("could not gt dms: %v", err)
	}

	for _, dm := range dms {
		err := orm.Da.UpdateDMReadById(dm.Id)
		if err != nil {
			return nil, fmt.Errorf("could not mark dm %v as read: %v", dm.Id, err)
		}
	}
	return &wrapperspb.StringValue{Value: "OK"}, nil
}

func (s *ApiService) SendDM(ctx context.Context, in *pb.DM) (*wrapperspb.Int32Value, error) {
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

	return wrapperspb.Int32(id), nil
}

func (s *ApiService) GetConversationDMs(in *wrapperspb.Int32Value, stream pb.Lenic_GetConversationDMsServer) error {
	dms, err := orm.Da.GetDMsByConversationId(int(in.Value))
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

func (s *ApiService) CreatePost(ctx context.Context, in *pb.Post) (*wrapperspb.StringValue, error) {

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

	return &wrapperspb.StringValue{Value: guid}, nil
}

func (s *ApiService) GetPost(ctx context.Context, in *wrapperspb.StringValue) (*pb.Post, error) {
	p, err := orm.Da.GetPostByGUID(in.Value)
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

func (s *ApiService) GetUserPosts(in *wrapperspb.StringValue, stream pb.Lenic_GetUserPostsServer) error {
	u, err := orm.Da.GetUserByName(in.Value)
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

func (s *ApiService) GetUserPublicPosts(in *wrapperspb.StringValue, stream pb.Lenic_GetUserPostsServer) error {
	u, err := orm.Da.GetUserByName(in.Value)
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

func (s *ApiService) GetFeed(in *wrapperspb.StringValue, stream pb.Lenic_GetFeedServer) error {
	u, err := orm.Da.GetUserByName(in.Value)
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

func (s *ApiService) RatePostUp(ctx context.Context, in *pb.PostRating) (*wrapperspb.StringValue, error) {
	err := orm.Da.RatePostUp(int(in.PostId), int(in.UserId))
	if err != nil {
		return nil, fmt.Errorf("could not rate post up: %v", err)
	}

	return &wrapperspb.StringValue{Value: "OK"}, nil
}

func (s *ApiService) RatePostDown(ctx context.Context, in *pb.PostRating) (*wrapperspb.StringValue, error) {
	err := orm.Da.RatePostDown(int(in.PostId), int(in.UserId))
	if err != nil {
		return nil, fmt.Errorf("could not rate post down: %v", err)
	}

	return &wrapperspb.StringValue{Value: "OK"}, nil
}

func (s *ApiService) UpdatePost(ctx context.Context, in *pb.Post) (*wrapperspb.StringValue, error) {
	p := model.Post{
		GUID:     in.PostGuid,
		Title:    in.Title,
		Content:  in.Content,
		IsPublic: in.IsPublic,
	}

	err := orm.Da.UpdatePost(p)
	if err != nil {
		return nil, fmt.Errorf("could not update post: %v", err)
	}

	return &wrapperspb.StringValue{Value: "OK"}, nil
}

func (s *ApiService) DeletePost(ctx context.Context, in *wrapperspb.StringValue) (*wrapperspb.StringValue, error) {
	err := orm.Da.DisablePost(in.Value)
	if err != nil {
		return nil, fmt.Errorf("could not delete post: %v", err)
	}

	return &wrapperspb.StringValue{Value: "OK"}, nil
}

func (s *ApiService) CreateComment(ctx context.Context, in *pb.Comment) (*wrapperspb.Int32Value, error) {

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

	return wrapperspb.Int32(id), nil
}

func (s *ApiService) GetComment(ctx context.Context, in *wrapperspb.Int32Value) (*pb.Comment, error) {
	c, err := orm.Da.GetCommentById(int(in.Value))
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

func (s *ApiService) GetCommentsFromPost(in *wrapperspb.StringValue, stream pb.Lenic_GetCommentsFromPostServer) error {
	comments, err := orm.Da.GetCommentsByPost(in.Value)
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

func (s *ApiService) RateCommentUp(ctx context.Context, in *pb.CommentRating) (*wrapperspb.StringValue, error) {
	err := orm.Da.RateCommentUp(int(in.CommentId), int(in.UserId))
	if err != nil {
		return nil, fmt.Errorf("could not rate comment up: %v", err)
	}

	return &wrapperspb.StringValue{Value: "OK"}, nil
}

func (s *ApiService) RateCommentDown(ctx context.Context, in *pb.CommentRating) (*wrapperspb.StringValue, error) {
	err := orm.Da.RateCommentDown(int(in.CommentId), int(in.UserId))
	if err != nil {
		return nil, fmt.Errorf("could not rate comment down: %v", err)
	}

	return &wrapperspb.StringValue{Value: "OK"}, nil
}

func (s *ApiService) UpdateComment(ctx context.Context, in *pb.Comment) (*wrapperspb.StringValue, error) {
	err := orm.Da.UpdateCommentText(int(in.Id), in.Content)
	if err != nil {
		return nil, fmt.Errorf("could not update comment: %v", err)
	}

	return &wrapperspb.StringValue{Value: "OK"}, nil
}

func (s *ApiService) DeleteComment(ctx context.Context, in *wrapperspb.Int32Value) (*wrapperspb.StringValue, error) {
	err := orm.Da.DisableComment(int(in.Value))
	if err != nil {
		return nil, fmt.Errorf("could not delete comment: %v", err)
	}

	return &wrapperspb.StringValue{Value: "OK"}, nil
}
