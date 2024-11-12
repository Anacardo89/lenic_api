package interceptor

import (
	"context"
	"errors"
	"strings"

	"github.com/Anacardo89/lenic_api/internal/data/orm"
	"github.com/Anacardo89/lenic_api/internal/pb"
	"github.com/Anacardo89/lenic_api/pkg/auth"
	"github.com/Anacardo89/lenic_api/pkg/logger"
	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func AuthUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	method := info.FullMethod
	if method == "/lenic.Lenic/Login" || method == "/lenic.Lenic/CreateUser" {
		return handler(ctx, req)
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		logger.Error.Println("missing metadata")
		return nil, errors.New("missing metadata")
	}

	token := extractToken(md)
	if token == "" {
		logger.Error.Println("missing token")
		return nil, errors.New("missing token")
	}

	claims, err := parseJWT(token)
	if err != nil {
		logger.Error.Println("invalid token")
		return nil, status.Errorf(codes.Unauthenticated, "invalid token")
	}

	if canBePublic(method) {
		uuid := getUUIDFromRequest(req)

		if getIsPublic(uuid) {
			return handler(ctx, req)
		} else {
			if !isSelfRequest(claims.Username, req) && !isFollowerRequest(claims.FollowingIDs, req) {
				logger.Error.Println("access denied for private post")
				return nil, status.Errorf(codes.PermissionDenied, "access denied for private post")
			}
		}
	}

	if isUserOnlyAccess(method) {
		if !isSelfRequest(claims.Username, req) {
			logger.Error.Println("access denied")
			return nil, status.Errorf(codes.PermissionDenied, "access denied")
		}
	} else if isFollowerAccess(method) {
		if !isSelfRequest(claims.Username, req) && !isFollowerRequest(claims.FollowingIDs, req) {
			logger.Error.Println("access restricted to followers")
			return nil, status.Errorf(codes.PermissionDenied, "access restricted to followers")
		}
	}
	return handler(ctx, req)
}

func extractToken(md metadata.MD) string {
	values := md["authorization"]
	if len(values) == 0 {
		return ""
	}
	return strings.TrimPrefix(values[0], "Bearer ")
}

func parseJWT(tokenString string) (*auth.Claims, error) {
	claims := &auth.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(token *jwt.Token) (interface{}, error) {
			return auth.JwtKey, nil
		},
	)
	if err != nil || !token.Valid {
		logger.Error.Println("invalid token")
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

func isFollowerAccess(method string) bool {
	switch method {
	case "/lenic.Lenic/GetUserPosts":
		return true
	default:
		return false
	}

}

func canBePublic(method string) bool {
	switch method {
	case "/lenic.Lenic/GetPost":
		return true
	case "/lenic.Lenic/GetComment":
		return true
	case "/lenic.Lenic/GetCommentsFromPost": //Stream
		return true
	case "/lenic.Lenic/RatePostUp":
		return true
	case "/lenic.Lenic/RatePostDown":
		return true
	case "/lenic.Lenic/RateCommentUp":
		return true
	case "/lenic.Lenic/RateCommentDown":
		return true
	default:
		return false
	}
}

func isFollowerRequest(following []string, request interface{}) bool {
	switch req := request.(type) {
	case *pb.GetUserPostsRequest:
		for _, f := range following {
			if f == req.Username {
				return true
			}
		}
		return false
	case *pb.GetPostRequest:
		p, err := orm.Da.GetPostByGUID(req.Uuid)
		if err != nil {
			return false
		}
		u, err := orm.Da.GetUserByID(p.AuthorId)
		if err != nil {
			return false
		}
		for _, f := range following {
			if f == u.UserName {
				return true
			}
		}
		return false
	case *pb.GetCommentRequest:
		c, err := orm.Da.GetCommentById(int(req.Id))
		if err != nil {
			return false
		}
		p, err := orm.Da.GetPostByGUID(c.PostGUID)
		if err != nil {
			return false
		}
		u, err := orm.Da.GetUserByID(p.AuthorId)
		if err != nil {
			return false
		}
		for _, f := range following {
			if f == u.UserName {
				return true
			}
		}
		return false
	case *pb.GetCommentsFromPostRequest:
		p, err := orm.Da.GetPostByGUID(req.Uuid)
		if err != nil {
			return false
		}
		u, err := orm.Da.GetUserByID(p.AuthorId)
		if err != nil {
			return false
		}
		for _, f := range following {
			if f == u.UserName {
				return true
			}
		}
		return false
	case *pb.PostRating:
		p, err := orm.Da.GetPostByID(int(req.PostId))
		if err != nil {
			return false
		}
		u, err := orm.Da.GetUserByID(p.AuthorId)
		if err != nil {
			return false
		}
		for _, f := range following {
			if f == u.UserName {
				return true
			}
		}
		return false
	case *pb.CommentRating:
		c, err := orm.Da.GetCommentById(int(req.CommentId))
		if err != nil {
			return false
		}
		u, err := orm.Da.GetUserByID(c.AuthorId)
		if err != nil {
			return false
		}
		for _, f := range following {
			if f == u.UserName {
				return true
			}
		}
		return false
	default:
		return false
	}
}

func getUUIDFromRequest(request interface{}) string {
	switch req := request.(type) {
	case *pb.GetPostRequest:
		return req.Uuid
	case *pb.GetCommentRequest:
		c, err := orm.Da.GetCommentById(int(req.Id))
		if err != nil {
			return ""
		}
		return c.PostGUID
	case *pb.GetCommentsFromPostRequest:
		return req.Uuid
	case *pb.PostRating:
		p, err := orm.Da.GetPostByID(int(req.PostId))
		if err != nil {
			return ""
		}
		return p.GUID
	case *pb.CommentRating:
		c, err := orm.Da.GetCommentById(int(req.CommentId))
		if err != nil {
			return ""
		}
		return c.PostGUID
	default:
		return ""
	}
}

func getIsPublic(uuid string) bool {
	p, err := orm.Da.GetPostByGUID(uuid)
	if err != nil {
		return false
	}
	return p.IsPublic
}

func isUserOnlyAccess(method string) bool {
	switch method {
	case "/lenic.Lenic/ActivateUser":
		return true
	case "/lenic.Lenic/UpdateUserPass":
		return true
	case "/lenic.Lenic/DeleteUser":
		return true
	case "/lenic.Lenic/FollowUser":
		return true
	case "/lenic.Lenic/AcceptFollowUser":
		return true
	case "/lenic.Lenic/UnfollowUser":
		return true
	case "/lenic.Lenic/StartConversation":
		return true
	case "/lenic.Lenic/GetUserConversations": // stream
		return true
	case "/lenic.Lenic/ReadConversation":
		return true
	case "/lenic.Lenic/SendDM":
		return true
	case "/lenic.Lenic/GetConversationDMs": // stream
		return true
	case "/lenic.Lenic/CreatePost":
		return true
	case "/lenic.Lenic/GetFeed": // stream
		return true
	case "/lenic.Lenic/UpdatePost":
		return true
	case "/lenic.Lenic/DeletePost":
		return true
	case "/lenic.Lenic/CreateComment":
		return true
	case "/lenic.Lenic/UpdateComment":
		return true
	case "/lenic.Lenic/DeleteComment":
		return true
	default:
		return false
	}
}

func isSelfRequest(username string, request interface{}) bool {
	switch req := request.(type) {
	case *pb.ActivateUserRequest:
		return req.Username == username
	case *pb.User:
		return req.Username == username
	case *pb.DeleteUserRequest:
		return req.Username == username
	case *pb.FollowUserRequest:
		u, err := orm.Da.GetUserByID(int(req.FollowerId))
		if err != nil {
			return false
		}
		return u.UserName == username
	case *pb.AcceptFollowRequest:
		u, err := orm.Da.GetUserByID(int(req.FollowedId))
		if err != nil {
			return false
		}
		return u.UserName == username
	case *pb.UnfollowRequest:
		u, err := orm.Da.GetUserByID(int(req.FollowedId))
		if err != nil {
			return false
		}
		return u.UserName == username
	case *pb.Conversation:
		u1, err := orm.Da.GetUserByID(int(req.User1Id))
		if err != nil {
			return false
		}
		u2, err := orm.Da.GetUserByID(int(req.User1Id))
		if err != nil {
			return false
		}
		return u1.UserName == username || u2.UserName == username
	case *pb.GetUserConversationsRequest:
		return req.Username == username
	case *pb.ReadConversationRequest:
		c, err := orm.Da.GetConversationById(int(req.Id))
		if err != nil {
			return false
		}
		u1, err := orm.Da.GetUserByID(int(c.User1Id))
		if err != nil {
			return false
		}
		u2, err := orm.Da.GetUserByID(int(c.User1Id))
		if err != nil {
			return false
		}
		return u1.UserName == username || u2.UserName == username
	case *pb.DM:
		c, err := orm.Da.GetConversationById(int(req.Id))
		if err != nil {
			return false
		}
		u1, err := orm.Da.GetUserByID(int(c.User1Id))
		if err != nil {
			return false
		}
		u2, err := orm.Da.GetUserByID(int(c.User1Id))
		if err != nil {
			return false
		}
		return u1.UserName == username || u2.UserName == username
	case *pb.GetConversationDMsRequest:
		c, err := orm.Da.GetConversationById(int(req.Id))
		if err != nil {
			return false
		}
		u1, err := orm.Da.GetUserByID(int(c.User1Id))
		if err != nil {
			return false
		}
		u2, err := orm.Da.GetUserByID(int(c.User1Id))
		if err != nil {
			return false
		}
		return u1.UserName == username || u2.UserName == username
	case *pb.Post:
		u, err := orm.Da.GetUserByID(int(req.AuthorId))
		if err != nil {
			return false
		}
		return u.UserName == username
	case *pb.DeletePostRequest:
		p, err := orm.Da.GetPostByGUID(req.Uuid)
		if err != nil {
			return false
		}
		u, err := orm.Da.GetUserByID(p.AuthorId)
		if err != nil {
			return false
		}
		return u.UserName == username
	case *pb.GetFeedRequest:
		return req.Username == username
	case *pb.Comment:
		p, err := orm.Da.GetPostByGUID(req.PostGuid)
		if err != nil {
			return false
		}
		u, err := orm.Da.GetUserByID(p.AuthorId)
		if err != nil {
			return false
		}
		return u.UserName == username
	case *pb.DeleteCommentRequest:
		c, err := orm.Da.GetCommentById(int(req.Id))
		if err != nil {
			return false
		}
		u, err := orm.Da.GetUserByID(c.AuthorId)
		if err != nil {
			return false
		}
		return u.UserName == username
	default:
		return false
	}
}

// func getAuthorFromRequest(request interface{}) string {
// 	switch req := request.(type) {
// 	case *pb.GetPostRequest:
// 		p, err := orm.Da.GetPostByGUID(req.Uuid)
// 		if err != nil {
// 			return ""
// 		}
// 		u, err := orm.Da.GetUserByID(p.AuthorId)
// 		if err != nil {
// 			return ""
// 		}
// 		return u.UserName
// 	case *pb.GetCommentRequest:
// 		c, err := orm.Da.GetCommentById(int(req.Id))
// 		if err != nil {
// 			return ""
// 		}
// 		u, err := orm.Da.GetUserByID(c.AuthorId)
// 		if err != nil {
// 			return ""
// 		}
// 		return u.UserName
// 	case *pb.GetCommentsFromPostRequest:
// 		p, err := orm.Da.GetPostByGUID(req.Uuid)
// 		if err != nil {
// 			return ""
// 		}
// 		u, err := orm.Da.GetUserByID(p.AuthorId)
// 		if err != nil {
// 			return ""
// 		}
// 		return u.UserName
// 	case *pb.PostRating:
// 		p, err := orm.Da.GetPostByID(int(req.PostId))
// 		if err != nil {
// 			return ""
// 		}
// 		u, err := orm.Da.GetUserByID(p.AuthorId)
// 		if err != nil {
// 			return ""
// 		}
// 		return u.UserName
// 	case *pb.CommentRating:
// 		c, err := orm.Da.GetCommentById(int(req.CommentId))
// 		if err != nil {
// 			return ""
// 		}
// 		u, err := orm.Da.GetUserByID(c.AuthorId)
// 		if err != nil {
// 			return ""
// 		}
// 		return u.UserName
// 	default:
// 		return ""
// 	}
// }
