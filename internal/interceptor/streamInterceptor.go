package interceptor

import (
	"context"
	"errors"
	"fmt"

	"github.com/Anacardo89/lenic_api/pkg/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

func AuthStreamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {

	ctx := ss.Context()

	claims, err := extractClaimsFromContext(ctx) // Your function to extract claims
	if err != nil {
		return status.Errorf(codes.Unauthenticated, "failed to extract claims: %v", err)
	}

	method := info.FullMethod

	req, err := extractRequestFromStream(ss)
	if err != nil {
		return fmt.Errorf("coul not get request from stream: %v", err)
	}

	if canBePublic(method) {
		uuid := getUUIDFromRequest(req)

		if getIsPublic(uuid) {
			return handler(srv, ss)
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
	return handler(srv, ss)
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

func extractRequestFromStream(stream grpc.ServerStream) (interface{}, error) {
	var req interface{}

	err := stream.RecvMsg(&req)
	if err != nil {
		return nil, fmt.Errorf("failed to read message from stream: %w", err)
	}

	if p, ok := peer.FromContext(stream.Context()); ok {
		fmt.Printf("Request received from peer: %s\n", p.Addr)
	}

	return req, nil
}
