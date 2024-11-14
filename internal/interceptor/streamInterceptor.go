package interceptor

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/Anacardo89/lenic_api/internal/pb"
	"github.com/Anacardo89/lenic_api/pkg/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

func AuthStreamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {

	ctx := ss.Context()

	claims, err := extractClaimsFromContext(ctx)
	if err != nil {
		return status.Errorf(codes.Unauthenticated, "failed to extract claims: %v", err)
	}

	method := info.FullMethod

	bs := &BufferedStream{ServerStream: ss}

	req, err := extractRequestFromStream(bs, method)
	if err != nil {
		return fmt.Errorf("coul not get request from stream: %v", err)
	}

	if canBePublic(method) {
		uuid := getUUIDFromRequest(req)

		if getIsPublic(uuid) {
			return handler(srv, bs)
		} else {
			if !isSelfRequest(claims.Username, req) && !isFollowerRequest(claims.FollowingIDs, req) {
				return status.Errorf(codes.PermissionDenied, "access denied for private post")
			}
		}
	}

	if isUserOnlyAccess(method) {
		if !isSelfRequest(claims.Username, req) {
			return status.Errorf(codes.PermissionDenied, "access denied")
		}
	} else if isFollowerAccess(method) {
		if !isSelfRequest(claims.Username, req) && !isFollowerRequest(claims.FollowingIDs, req) {
			return status.Errorf(codes.PermissionDenied, "access restricted to followers")
		}
	}
	return handler(srv, bs)
}

func extractClaimsFromContext(ctx context.Context) (*auth.Claims, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("metadata not found in context")
	}

	authHeaders, exists := md["authorization"]
	if !exists || len(authHeaders) == 0 {
		return nil, errors.New("authorization header not provided")
	}

	token := extractToken(md)
	if token == "" {
		return nil, errors.New("missing token")
	}

	claims, err := parseJWT(token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token")
	}

	return claims, nil

}

func extractRequestFromStream(bs *BufferedStream, method string) (proto.Message, error) {
	req, err := extractRequestType(method)
	if err != nil {
		return nil, err
	}

	err = bs.RecvMsg(req)
	if err != nil {
		return nil, fmt.Errorf("could not read message from stream: %v", err)
	}

	return req, nil
}

func extractRequestType(method string) (proto.Message, error) {
	switch method {
	case "/lenic.Lenic/SearchUsers":
		return &pb.SearchUsersRequest{}, nil
	case "/lenic.Lenic/GetUserFollowers":
		return &pb.GetUserFollowersRequest{}, nil
	case "/lenic.Lenic/GetUserFollowing":
		return &pb.GetUserFollowingRequest{}, nil
	case "/lenic.Lenic/GetUserConversations":
		return &pb.GetUserConversationsRequest{}, nil
	case "/lenic.Lenic/GetConversationDMs":
		return &pb.GetConversationDMsRequest{}, nil
	case "/lenic.Lenic/GetUserPosts":
		return &pb.GetUserPostsRequest{}, nil
	case "/lenic.Lenic/GetUserPublicPosts":
		return &pb.GetUserPublicPostsRequest{}, nil
	case "/lenic.Lenic/GetFeed":
		return &pb.GetFeedRequest{}, nil
	case "/lenic.Lenic/GetCommentsFromPost":
		return &pb.GetCommentsFromPostRequest{}, nil
	default:
		return nil, fmt.Errorf("unknown method: %s", method)
	}
}

type BufferedStream struct {
	grpc.ServerStream
	bufferedRequest proto.Message
}

func (bs *BufferedStream) RecvMsg(m interface{}) error {
	// If the message is already buffered, return it without consuming the stream
	if bs.bufferedRequest != nil {
		// Clone the buffered request into the provided message.
		reflect.ValueOf(m).Elem().Set(reflect.ValueOf(bs.bufferedRequest).Elem())
		return nil
	}

	// Otherwise, read the message from the stream normally
	err := bs.ServerStream.RecvMsg(m)
	if err != nil {
		return err
	}

	// Buffer the received message for future access without consuming the stream
	bs.bufferedRequest = m.(proto.Message)
	return nil
}
