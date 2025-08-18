// Mills Mess
// Licensed under the Mills Mess License Agreement
// See LICENSE.md in the root of this repository.

package socket

import (
	"context"
	"github.com/staringfun/millsmess/libs/base"
	"github.com/staringfun/millsmess/libs/internal-core-api"
	"github.com/staringfun/millsmess/libs/public-core-api"
	"github.com/staringfun/millsmess/libs/test-utils"
	"github.com/staringfun/millsmess/libs/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestPreconnectTxNoAuth(t *testing.T) {
	baseCtx := context.Background()
	coreEngine := &test_utils.MockedCoreEngine{}
	var resHeaders map[string]string
	tx := &TxPreconnect{
		BaseTx: &base.BaseTx{
			Ctx: t.Context(),
		},
		Args: PreconnectArgs{
			Token:         "token",
			XForwardedFor: "0.0.0.0",
			UserAgent:     "chrome",
			JoinRoomID:    "r-id",
		},
		CoreInternal: &internal_core_api.InternalCoreAPI{
			PublicCoreAPI: public_core_api.PublicCoreAPI{
				Engine: coreEngine,
			},
		},
		baseCtx: baseCtx,
	}
	coreEngine.On("Fetch", mock.Anything, mock.Anything, mock.Anything, mock.Anything, t.Context()).
		Run(func(args mock.Arguments) {
			headers := args.Get(2).(map[string]string)
			assert.Empty(t, headers["Cookie"])
			assert.NotEmpty(t, headers["Authorization"])
		}).
		Return(200, resHeaders, nil)

	defer tx.DeferPrepare()
	assert.Nil(t, tx.Prepare())

	finished, err := tx.LoadData()
	assert.True(t, finished)
	assert.Nil(t, err)

	assert.Nil(t, tx.ClientData)
}

func TestPreconnectTxSuccessUserAgentXForwardedForCookie(t *testing.T) {
	clock := &test_utils.MockedClock{
		Value: time.Now(),
	}
	baseCtx := context.Background()
	coreEngine := &test_utils.MockedCoreEngine{}
	user := &types.BaseUser{
		ID:       "u-id",
		Username: "u-username",
		Role:     types.UserRoleGuest,
		BannedAt: nil,
	}
	var resHeaders map[string]string
	tx := &TxPreconnect{
		BaseTx: &base.BaseTx{
			Clock: clock,
			Ctx:   t.Context(),
		},
		Args: PreconnectArgs{
			Token:         "",
			Cookie:        "wefew",
			XForwardedFor: "0.0.0.0",
			UserAgent:     "chrome",
			JoinRoomID:    "r-id",
		},
		CoreInternal: &internal_core_api.InternalCoreAPI{
			PublicCoreAPI: public_core_api.PublicCoreAPI{
				Engine: coreEngine,
			},
		},
		baseCtx: baseCtx,
	}
	coreEngine.On("Fetch", mock.Anything, mock.Anything, mock.Anything, mock.Anything, t.Context()).
		Run(func(args mock.Arguments) {
			headers := args.Get(2).(map[string]string)
			assert.NotEmpty(t, headers["Cookie"])
			assert.Empty(t, headers["Authorization"])
			assert.Equal(t, tx.Args.XForwardedFor, headers["X-Forwarded-For"])
			assert.Equal(t, tx.Args.UserAgent, headers["User-Agent"])

			data := args.Get(3).(*public_core_api.FetchMeResponse)
			data.Data = user
		}).
		Return(200, resHeaders, nil)

	defer tx.DeferPrepare()
	assert.Nil(t, tx.Prepare())

	finished, err := tx.LoadData()
	assert.True(t, finished)
	assert.Nil(t, err)

	assert.NotNil(t, tx.ClientData)

	assert.Equal(t, user, tx.ClientData.User)
	assert.NotEmpty(t, tx.ClientData.PlayerID)
}

func TestPreconnectTxSuccessMatchArgsCookie(t *testing.T) {
	clock := &test_utils.MockedClock{
		Value: time.Now(),
	}
	baseCtx := context.Background()
	coreEngine := &test_utils.MockedCoreEngine{}
	user := &types.BaseUser{
		ID:       "u-id",
		Username: "u-username",
		Role:     types.UserRoleGuest,
		BannedAt: nil,
	}
	var resHeaders map[string]string
	tx := &TxPreconnect{
		BaseTx: &base.BaseTx{
			Clock: clock,
			Ctx:   t.Context(),
		},
		Args: PreconnectArgs{
			Token:         "",
			Cookie:        "wefew",
			XForwardedFor: "0.0.0.0",
			UserAgent:     "chrome",
			JoinRoomID:    "r-id",
		},
		CoreInternal: &internal_core_api.InternalCoreAPI{
			PublicCoreAPI: public_core_api.PublicCoreAPI{
				Engine: coreEngine,
			},
		},
		baseCtx: baseCtx,
	}
	coreEngine.On("Fetch", mock.Anything, mock.Anything, mock.Anything, mock.Anything, t.Context()).
		Run(func(args mock.Arguments) {
			data := args.Get(3).(*public_core_api.FetchMeResponse)
			data.Data = user
		}).
		Return(200, resHeaders, nil)

	defer tx.DeferPrepare()
	assert.Nil(t, tx.Prepare())

	finished, err := tx.LoadData()
	assert.True(t, finished)
	assert.Nil(t, err)

	assert.NotNil(t, tx.ClientData)

	assert.Equal(t, user, tx.ClientData.User)
	assert.NotEmpty(t, tx.ClientData.PlayerID)

	assert.NotNil(t, tx.ClientData.joinRoomID)
	assert.Equal(t, tx.Args.JoinRoomID, *tx.ClientData.joinRoomID)
}

func TestPreconnectTxBannedMatchArgsToken(t *testing.T) {
	clock := &test_utils.MockedClock{
		Value: time.Now(),
	}
	baseCtx := context.Background()
	coreEngine := &test_utils.MockedCoreEngine{}
	user := &types.BaseUser{
		ID:       "u-id",
		Username: "u-username",
		Role:     types.UserRoleGuest,
		BannedAt: &time.Time{},
	}
	var resHeaders map[string]string
	tx := &TxPreconnect{
		BaseTx: &base.BaseTx{
			Clock: clock,
			Ctx:   t.Context(),
		},
		Args: PreconnectArgs{
			Token:         "token",
			XForwardedFor: "0.0.0.0",
			UserAgent:     "chrome",
			JoinRoomID:    "r-id",
		},
		CoreInternal: &internal_core_api.InternalCoreAPI{
			PublicCoreAPI: public_core_api.PublicCoreAPI{
				Engine: coreEngine,
			},
		},
		baseCtx: baseCtx,
	}
	coreEngine.On("Fetch", mock.Anything, mock.Anything, mock.Anything, mock.Anything, t.Context()).
		Run(func(args mock.Arguments) {
			headers := args.Get(2).(map[string]string)
			assert.Empty(t, headers["Cookie"])
			assert.NotEmpty(t, headers["Authorization"])

			data := args.Get(3).(*public_core_api.FetchMeResponse)
			data.Data = user
		}).
		Return(200, resHeaders, nil)

	defer tx.DeferPrepare()
	assert.Nil(t, tx.Prepare())

	finished, err := tx.LoadData()
	assert.True(t, finished)
	assert.Nil(t, err)

	assert.Nil(t, tx.ClientData)
}

func TestPreconnectTxBannedBeforeMatchArgsToken(t *testing.T) {
	clock := &test_utils.MockedClock{
		Value: time.Now(),
	}
	bannedAt := time.Now()
	bannedAt.Add(time.Hour * 48)
	baseCtx := context.Background()
	coreEngine := &test_utils.MockedCoreEngine{}
	user := &types.BaseUser{
		ID:       "u-id",
		Username: "u-username",
		Role:     types.UserRoleGuest,
		BannedAt: &bannedAt,
	}
	var resHeaders map[string]string
	tx := &TxPreconnect{
		BaseTx: &base.BaseTx{
			Clock: clock,
			Ctx:   t.Context(),
		},
		Args: PreconnectArgs{
			Token:         "token",
			XForwardedFor: "0.0.0.0",
			UserAgent:     "chrome",
			JoinRoomID:    "r-id",
		},
		CoreInternal: &internal_core_api.InternalCoreAPI{
			PublicCoreAPI: public_core_api.PublicCoreAPI{
				Engine: coreEngine,
			},
		},
		baseCtx: baseCtx,
	}
	coreEngine.On("Fetch", mock.Anything, mock.Anything, mock.Anything, mock.Anything, t.Context()).
		Run(func(args mock.Arguments) {
			headers := args.Get(2).(map[string]string)
			assert.Empty(t, headers["Cookie"])
			assert.NotEmpty(t, headers["Authorization"])

			data := args.Get(3).(*public_core_api.FetchMeResponse)
			data.Data = user
		}).
		Return(200, resHeaders, nil)

	defer tx.DeferPrepare()
	assert.Nil(t, tx.Prepare())

	finished, err := tx.LoadData()
	assert.True(t, finished)
	assert.Nil(t, err)

	assert.NotNil(t, tx.ClientData)

	assert.Equal(t, user, tx.ClientData.User)

	assert.Equal(t, &tx.Args.JoinRoomID, tx.ClientData.joinRoomID)
}
