package grpc

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	ssov1 "github.com/MOONLAYT400/Proto_sso/gen/go/sso"
	grpcLog "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpcRetry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	api ssov1.AuthClient
	log *slog.Logger
}

func New(ctx context.Context, log *slog.Logger, addr string, timeout time.Duration, retriesCount int) (*Client,error) {
	const op = "clients.sso.grpc.New"

	retryOpts := []grpcRetry.CallOption{
		grpcRetry.WithCodes(codes.NotFound, codes.Aborted,codes.DeadlineExceeded),
		grpcRetry.WithMax(uint(retriesCount)),
		grpcRetry.WithPerRetryTimeout(timeout),
	}

	logOpts := []grpcLog.Option{
		grpcLog.WithLogOnEvents(grpcLog.PayloadReceived, grpcLog.PayloadSent),
	}

	cc,err := grpc.DialContext(ctx,addr,grpc.WithTransportCredentials(insecure.NewCredentials()),
	grpc.WithChainUnaryInterceptor(
		grpcRetry.UnaryClientInterceptor(retryOpts...),
		grpcLog.UnaryClientInterceptor(InterceptorLogger(log),logOpts...),
	),
)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Client{
		api: ssov1.NewAuthClient(cc),
		log: log,
	}, nil
}

func (c *Client) IsAdmin(ctx context.Context, userId int64) (bool, error) {
	const op = "clients.sso.grpc.IsAdmin"

	response, err := c.api.IsAdmin(ctx, &ssov1.IsAdminRequest{UserId: userId})
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return response.IsAdmin, nil
}

func InterceptorLogger(log *slog.Logger) grpcLog.Logger {
	return grpcLog.LoggerFunc(
		func(ctx context.Context, level grpcLog.Level, msg string, fields ...any) {
			log.Log(ctx,slog.Level(level), msg, fields...)
		},	
	)
}