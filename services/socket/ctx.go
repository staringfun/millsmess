package socket

import (
	"context"
	"github.com/staringfun/millsmess/libs/base"
	"github.com/staringfun/millsmess/libs/types"
)

const PlayerIDKey = "playerID"

func WithPlayerID(ctx context.Context, id types.PlayerID) context.Context {
	return context.WithValue(ctx, PlayerIDKey, id)
}

func GetPlayerID(ctx context.Context) types.PlayerID {
	id := ctx.Value(PlayerIDKey)
	playerID, _ := id.(types.PlayerID)
	return playerID
}

const UserIDKey = "userID"

func WithUserID(ctx context.Context, id types.UserID) context.Context {
	return context.WithValue(ctx, UserIDKey, id)
}

func GetUserID(ctx context.Context) types.UserID {
	id := ctx.Value(UserIDKey)
	playerID, _ := id.(types.UserID)
	return playerID
}

func WithConnectionLoggerFields(ctx context.Context, connectionCtx context.Context) context.Context {
	return base.GetLogger(ctx).With().
		Any("playerID", GetPlayerID(connectionCtx)).
		Any("userID", GetUserID(connectionCtx)).
		WithContext(ctx)
}
