package gapi

import (
	"context"
	"fmt"
	"testing"
	"time"

	db "github.com/spaghetti-lover/simplebank/db/sqlc"
	"github.com/spaghetti-lover/simplebank/token"
	"github.com/spaghetti-lover/simplebank/util"
	"github.com/spaghetti-lover/simplebank/worker"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
)

func newTestServer(t *testing.T, store db.Store, taskDistributor worker.TaskDistributor) *Server {
	config := util.Config{
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := NewServer(config, store, taskDistributor)
	require.NoError(t, err)

	return server
}

func newContextWithBearerToken(t *testing.T, tokenMaker token.Maker, username string, role string, duration time.Duration, tokenType token.TokenType) context.Context {
	accessToken, _, err := tokenMaker.CreateToken(username, role, duration, tokenType)
	require.NoError(t, err)

	bearerToken := fmt.Sprintf("%s %s", authorizationBearer, accessToken)
	md := metadata.MD{
		authorizationHeader: []string{
			bearerToken,
		},
	}

	return metadata.NewIncomingContext(context.Background(), md)
}
