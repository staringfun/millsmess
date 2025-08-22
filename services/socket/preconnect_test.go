// Mills Mess
// Licensed under the Mills Mess License Agreement
// See LICENSE.md in the root of this repository.

package socket

import (
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
	var resHeaders map[string]string
	coreEngine := &test_utils.MockedCoreEngine{}
	preconnect := &Preconnect{
		Base: &base.Base{
			Ctx: t.Context(),
		},
		CoreInternal: &internal_core_api.InternalCoreAPI{
			PublicCoreAPI: &public_core_api.PublicCoreAPI{
				Engine: coreEngine,
			},
		},
	}
	coreEngine.On("Fetch", mock.Anything, mock.Anything, mock.Anything, mock.Anything, t.Context()).
		Run(func(args mock.Arguments) {
			headers := args.Get(2).(map[string]string)
			assert.Empty(t, headers["Cookie"])
			assert.NotEmpty(t, headers["Authorization"])
		}).
		Return(200, resHeaders, nil)

	clientData, err := preconnect.Run(PreconnectArgs{
		Token:         "token",
		XForwardedFor: "0.0.0.0",
		UserAgent:     "chrome",
		JoinArgs: &types.MV1RoomJoin{
			RoomID: "r-id",
		},
	}, t.Context())

	assert.Nil(t, clientData)
	assert.Nil(t, err)
}

func TestPreconnectTxSuccessUserAgentXForwardedForCookie(t *testing.T) {
	clock := &test_utils.MockedClock{
		Value: time.Now(),
	}
	user := &types.BaseUser{
		ID:       "u-id",
		Username: "u-username",
		Role:     types.UserRoleGuest,
	}
	coreEngine := &test_utils.MockedCoreEngine{}
	preconnect := &Preconnect{
		Base: &base.Base{
			Ctx:   t.Context(),
			Clock: clock,
		},
		CoreInternal: &internal_core_api.InternalCoreAPI{
			PublicCoreAPI: &public_core_api.PublicCoreAPI{
				Engine: coreEngine,
			},
		},
	}
	var resHeaders map[string]string
	coreEngine.On("Fetch", mock.Anything, mock.Anything, mock.Anything, mock.Anything, t.Context()).
		Run(func(args mock.Arguments) {
			headers := args.Get(2).(map[string]string)
			assert.NotEmpty(t, headers["Cookie"])
			assert.Empty(t, headers["Authorization"])

			data := args.Get(3).(*public_core_api.FetchMeResponse)
			data.Data = user
		}).
		Return(200, resHeaders, nil)

	clientData, err := preconnect.Run(PreconnectArgs{
		Token:         "",
		Cookie:        "wefew",
		XForwardedFor: "0.0.0.0",
		UserAgent:     "chrome",
		JoinArgs: &types.MV1RoomJoin{
			RoomID: "r-id",
		},
	}, t.Context())
	assert.Nil(t, err)

	assert.NotNil(t, clientData)

	assert.Equal(t, user, clientData.User)
	assert.NotEmpty(t, clientData.PlayerID)
}

func TestPreconnectTxSuccessMatchArgsCookie(t *testing.T) {
	clock := &test_utils.MockedClock{
		Value: time.Now(),
	}
	user := &types.BaseUser{
		ID:       "u-id",
		Username: "u-username",
		Role:     types.UserRoleGuest,
	}
	coreEngine := &test_utils.MockedCoreEngine{}
	preconnect := &Preconnect{
		Base: &base.Base{
			Ctx:   t.Context(),
			Clock: clock,
		},
		CoreInternal: &internal_core_api.InternalCoreAPI{
			PublicCoreAPI: &public_core_api.PublicCoreAPI{
				Engine: coreEngine,
			},
		},
	}
	var resHeaders map[string]string
	coreEngine.On("Fetch", mock.Anything, mock.Anything, mock.Anything, mock.Anything, t.Context()).
		Run(func(args mock.Arguments) {
			data := args.Get(3).(*public_core_api.FetchMeResponse)
			data.Data = user
		}).
		Return(200, resHeaders, nil)
	args := PreconnectArgs{
		Token:         "",
		Cookie:        "wefew",
		XForwardedFor: "0.0.0.0",
		UserAgent:     "chrome",
		JoinArgs: &types.MV1RoomJoin{
			RoomID: "r-id",
		},
	}
	clientData, err := preconnect.Run(args, t.Context())
	assert.Nil(t, err)

	assert.NotNil(t, clientData)

	assert.Equal(t, user, clientData.User)
	assert.NotEmpty(t, clientData.PlayerID)

	assert.NotNil(t, clientData.joinArgs)
	assert.Equal(t, args.JoinArgs.RoomID, clientData.joinArgs.RoomID)
}

func TestPreconnectTxBannedMatchArgsToken(t *testing.T) {
	clock := &test_utils.MockedClock{
		Value: time.Now(),
	}
	user := &types.BaseUser{
		ID:       "u-id",
		Username: "u-username",
		Role:     types.UserRoleGuest,
		BannedAt: &time.Time{},
	}
	coreEngine := &test_utils.MockedCoreEngine{}
	preconnect := &Preconnect{
		Base: &base.Base{
			Ctx:   t.Context(),
			Clock: clock,
		},
		CoreInternal: &internal_core_api.InternalCoreAPI{
			PublicCoreAPI: &public_core_api.PublicCoreAPI{
				Engine: coreEngine,
			},
		},
	}
	var resHeaders map[string]string
	coreEngine.On("Fetch", mock.Anything, mock.Anything, mock.Anything, mock.Anything, t.Context()).
		Run(func(args mock.Arguments) {
			headers := args.Get(2).(map[string]string)
			assert.Empty(t, headers["Cookie"])
			assert.NotEmpty(t, headers["Authorization"])

			data := args.Get(3).(*public_core_api.FetchMeResponse)
			data.Data = user
		}).
		Return(200, resHeaders, nil)

	args := PreconnectArgs{
		Token:         "token",
		XForwardedFor: "0.0.0.0",
		UserAgent:     "chrome",
	}
	clientData, err := preconnect.Run(args, t.Context())
	assert.Nil(t, err)

	assert.Nil(t, clientData)
}

func TestPreconnectTxBannedBeforeMatchArgsToken(t *testing.T) {
	clock := &test_utils.MockedClock{
		Value: time.Now(),
	}
	bannedAt := time.Now()
	bannedAt.Add(time.Hour * 48)
	user := &types.BaseUser{
		ID:       "u-id",
		Username: "u-username",
		Role:     types.UserRoleGuest,
		BannedAt: &bannedAt,
	}
	coreEngine := &test_utils.MockedCoreEngine{}
	preconnect := &Preconnect{
		Base: &base.Base{
			Ctx:   t.Context(),
			Clock: clock,
		},
		CoreInternal: &internal_core_api.InternalCoreAPI{
			PublicCoreAPI: &public_core_api.PublicCoreAPI{
				Engine: coreEngine,
			},
		},
	}
	var resHeaders map[string]string
	coreEngine.On("Fetch", mock.Anything, mock.Anything, mock.Anything, mock.Anything, t.Context()).
		Run(func(args mock.Arguments) {
			headers := args.Get(2).(map[string]string)
			assert.Empty(t, headers["Cookie"])
			assert.NotEmpty(t, headers["Authorization"])

			data := args.Get(3).(*public_core_api.FetchMeResponse)
			data.Data = user
		}).
		Return(200, resHeaders, nil)

	args := PreconnectArgs{
		Token:         "token",
		XForwardedFor: "0.0.0.0",
		UserAgent:     "chrome",
		JoinArgs: &types.MV1RoomJoin{
			RoomID: "r-id",
		},
	}
	clientData, err := preconnect.Run(args, t.Context())
	assert.Nil(t, err)

	assert.NotNil(t, clientData)

	assert.Equal(t, user, clientData.User)

	assert.Equal(t, clientData.joinArgs.RoomID, args.JoinArgs.RoomID)
}
