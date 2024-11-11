package interceptor

import (
	"context"
	"errors"
	"strings"

	"github.com/Anacardo89/lenic_api/pkg/auth"
	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("missing metadata")
	}

	token := extractToken(md)
	if token == "" {
		return nil, errors.New("missing token")
	}

	claims, err := parseJWT(token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token")
	}

	method := info.FullMethod
	if isUserOnlyAccess(method) {
		if !isSelfRequest(userID, req) {
			return nil, status.Errorf(codes.PermissionDenied, "access denied")
		}
	} else if isFollowerAccess(method) {
		if !isFollowerRequest(userID, req) {
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
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

func isUserOnlyAccess(method string) bool {
	switch method {
	case "/Lenic/ActivateUser":
		return true
	case "/Lenic/FollowUser":
		return true
	case "/Lenic/AcceptFollowUser":
		return true
	case "/Lenic/UnfollowUser":
		return true
	case "/Lenic/UpdateUserPass":
		return true
	case "/Lenic/DeleteUser":
		return true
	case "/Lenic/StartConversation":
		return true
	case "/Lenic/GetUserConversations":
		return true
	case "/Lenic/ReadConversation":
		return true
	case "/Lenic/SendDM":
		return true
	case "/Lenic/GetConversationDMs":
		return true
	case "/Lenic/CreatePost":
		return true
	case "/Lenic/GetFeed":
		return true
	case "/Lenic/RatePostUp":
		return true
	case "/Lenic/RatePostDown":
		return true
	case "/Lenic/UpdatePost":
		return true
	case "/Lenic/DeletePost":
		return true
	case "/Lenic/CreateComment":
		return true
	case "/Lenic/RateCommentUp":
		return true
	case "/Lenic/RateCommentDown":
		return true
	case "/Lenic/UpdateComment":
		return true
	case "/Lenic/DeleteComment":
		return true
	}
}
