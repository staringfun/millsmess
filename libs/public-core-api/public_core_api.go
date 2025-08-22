// Mills Mess
// Licensed under the Mills Mess License Agreement
// See LICENSE.md in the root of this repository.

package public_core_api

import (
	"context"
	"errors"
	"github.com/staringfun/millsmess/libs/types"
	"net/http"
	"time"
)

type PublicCoreAPIEngine interface {
	Fetch(query string, variables map[string]any, reqHeaders map[string]string, res any, ctx context.Context) (
		status int,
		resHeaders map[string]string,
		err error,
	)
}

type PublicCoreAPI struct {
	Engine PublicCoreAPIEngine
}

type AuthArgs struct {
	Token  string
	Cookie string
}

var UnauthorizedError = errors.New("unauthorized")

func (c *PublicCoreAPI) AppendAuthHeaders(args AuthArgs, headers map[string]string) {
	if args.Token != "" {
		headers["Authorization"] = "Bearer " + args.Token
	}
	if args.Cookie != "" {
		headers["Cookie"] = args.Cookie
	}
}

type FetchMeArgs struct {
	Auth AuthArgs
}

type FetchMeResponse struct {
	Data types.User `json:"data"`
}

var FetchMeError = errors.New("fetch me")

const MeQuery = `
query me {
	me {
		id
		username
		role
	}
}
`

func (c *PublicCoreAPI) FetchMeHeaders(args FetchMeArgs, headers map[string]string, ctx context.Context) (types.User, map[string]string, error) {
	if headers == nil {
		headers = map[string]string{}
	}
	c.AppendAuthHeaders(args.Auth, headers)
	var res FetchMeResponse
	status, resHeaders, err := c.Engine.Fetch(MeQuery, nil, headers, &res, ctx)
	if err != nil {
		return nil, nil, err
	}
	switch {
	case status == http.StatusUnauthorized:
		return nil, resHeaders, UnauthorizedError
	case status != http.StatusOK || res.Data == nil:
		return nil, resHeaders, FetchMeError
	default:
		return res.Data, resHeaders, nil
	}
}

type FetchUserRoomsArgs struct {
	Auth AuthArgs
}

type FetchUserRoomsResponse struct {
	Data []types.UserRoom `json:"data"`
}

var FetchUserRoomsError = errors.New("fetch me")

const UserRoomQuery = `
query ur {
	userRooms {
		playerID
		userID
		roomID
	}
}
`

func (c *PublicCoreAPI) FetchUserRoomHeaders(args FetchUserRoomsArgs, headers map[string]string, ctx context.Context) ([]types.UserRoom, map[string]string, error) {
	if headers == nil {
		headers = map[string]string{}
	}
	c.AppendAuthHeaders(args.Auth, headers)
	var res FetchUserRoomsResponse
	status, resHeaders, err := c.Engine.Fetch(UserRoomQuery, nil, headers, &res, ctx)
	if err != nil {
		return nil, nil, err
	}
	switch {
	case status == http.StatusUnauthorized:
		return nil, resHeaders, UnauthorizedError
	case status != http.StatusOK || res.Data == nil:
		return nil, resHeaders, FetchUserRoomsError
	default:
		return res.Data, resHeaders, nil
	}
}

func (c *PublicCoreAPI) FetchMe(args FetchMeArgs, ctx context.Context) (types.User, error) {
	res, _, err := c.FetchMeHeaders(args, nil, ctx)
	return res, err
}

var LockUserRoomError = errors.New("lock user room")

const LockUserRoomMutation = `
mutation lock {
	lockUserRoom {
		key
		expiresAt
	}
}
`

type LockUserRoomArgs struct {
}

type LockUserRoomResponse struct {
	Data *types.Lock `json:"data,omitempty"`
}

func (c *PublicCoreAPI) LockUserRoomHeaders(args LockUserRoomArgs, headers map[string]string, ctx context.Context) (*types.Lock, map[string]string, error) {
	if headers == nil {
		headers = map[string]string{}
	}
	var res LockUserRoomResponse
	status, resHeaders, err := c.Engine.Fetch(LockUserRoomMutation, nil, headers, &res, ctx)
	if err != nil {
		return nil, nil, err
	}
	switch {
	case status == http.StatusUnauthorized:
		return nil, resHeaders, UnauthorizedError
	case status != http.StatusOK || res.Data == nil:
		return nil, resHeaders, LockUserRoomError
	default:
		return res.Data, resHeaders, nil
	}
}

var UnlockUserRoomError = errors.New("unlock user room")

const UnlockUserRoomMutation = `
mutation unlock($key: LockKey!) {
	unlockUserRoom({key: $key}) {
		ok
	}
}
`

type UnlockUserRoomArgs struct {
	Key types.LockKey
}

type UserRoomResponse struct {
	Ok bool `json:"ok"`
}

type UnlockUserRoomResponse struct {
	Data *UserRoomResponse `json:"data,omitempty"`
}

func (c *PublicCoreAPI) UnlockUserRoomHeaders(args UnlockUserRoomArgs, headers map[string]string, ctx context.Context) (map[string]string, error) {
	if headers == nil {
		headers = map[string]string{}
	}
	var res UnlockUserRoomResponse
	status, resHeaders, err := c.Engine.Fetch(UnlockUserRoomMutation, map[string]any{
		"key": args.Key,
	}, headers, &res, ctx)
	if err != nil {
		return nil, err
	}
	switch {
	case status == http.StatusUnauthorized:
		return resHeaders, UnauthorizedError
	case status != http.StatusOK || res.Data == nil || !res.Data.Ok:
		return resHeaders, UnlockUserRoomError
	default:
		return resHeaders, nil
	}
}

var ExtendUserRoomError = errors.New("extend user room")

const ExtendUserRoomMutation = `
mutation extend($key: LockKey!) {
	extendUserRoom({key: $key}) {
		expiresAt
	}
}
`

type ExtendUserRoomArgs struct {
	Key types.LockKey
}

type ExtendUserRoomResponseData struct {
	ExpiresAt time.Time `json:"expiresAt,omitempty"`
}

type ExtendUserRoomResponse struct {
	Data *ExtendUserRoomResponseData `json:"data,omitempty"`
}

func (c *PublicCoreAPI) ExtendUserRoomHeaders(args ExtendUserRoomArgs, headers map[string]string, ctx context.Context) (*ExtendUserRoomResponseData, map[string]string, error) {
	if headers == nil {
		headers = map[string]string{}
	}
	var res ExtendUserRoomResponse
	status, resHeaders, err := c.Engine.Fetch(ExtendUserRoomMutation, map[string]any{
		"key": args.Key,
	}, headers, &res, ctx)
	if err != nil {
		return nil, nil, err
	}
	switch {
	case status == http.StatusUnauthorized:
		return nil, resHeaders, UnauthorizedError
	case status != http.StatusOK || res.Data == nil:
		return nil, resHeaders, ExtendUserRoomError
	default:
		return res.Data, resHeaders, nil
	}
}

var AddUserRoomError = errors.New("add user room")

const AddUserRoomMutation = `
mutation add($key: LockKey!, $roomID: RoomID!, $playerID: PlayerID!) {
	addUserRoom({key: $key, roomID: $roomID, playerID: $playerID}) {
		playerID
		userID
		roomID
	}
}
`

type AddUserRoomArgs struct {
	RoomID   types.RoomID   `json:"roomID,omitempty"`
	PlayerID types.PlayerID `json:"playerID,omitempty"`
	Key      types.LockKey  `json:"key,omitempty"`
}

type AddUserRoomResponseData = types.UserRoom

type AddUserRoomResponse struct {
	Data *AddUserRoomResponseData `json:"data,omitempty"`
}

func (c *PublicCoreAPI) AddUserRoomHeaders(args AddUserRoomArgs, headers map[string]string, ctx context.Context) (*AddUserRoomResponseData, map[string]string, error) {
	if headers == nil {
		headers = map[string]string{}
	}
	var res AddUserRoomResponse
	status, resHeaders, err := c.Engine.Fetch(AddUserRoomMutation, map[string]any{
		"key":      args.Key,
		"roomID":   args.RoomID,
		"playerID": args.PlayerID,
	}, headers, &res, ctx)
	if err != nil {
		return nil, nil, err
	}
	switch {
	case status == http.StatusUnauthorized:
		return nil, resHeaders, UnauthorizedError
	case status != http.StatusOK || res.Data == nil:
		return nil, resHeaders, AddUserRoomError
	default:
		return res.Data, resHeaders, nil
	}
}

var DeleteUserRoomError = errors.New("delete user room")

const DeleteUserRoomMutation = `
mutation delete($key: LockKey!, $playerID: PlayerID!) {
	deleteUserRoom({key: $key, playerID: $playerID}) {
		playerID
		userID
		roomID
	}
}
`

type DeleteUserRoomArgs struct {
	PlayerID types.PlayerID `json:"playerID,omitempty"`
	Key      types.LockKey  `json:"key,omitempty"`
}

type DeleteUserRoomResponseData = types.UserRoom

type DeleteUserRoomResponse struct {
	Data *DeleteUserRoomResponseData `json:"data,omitempty"`
}

func (c *PublicCoreAPI) DeleteUserRoomHeaders(args DeleteUserRoomArgs, headers map[string]string, ctx context.Context) (*AddUserRoomResponseData, map[string]string, error) {
	if headers == nil {
		headers = map[string]string{}
	}
	var res AddUserRoomResponse
	status, resHeaders, err := c.Engine.Fetch(DeleteUserRoomMutation, map[string]any{
		"key":      args.Key,
		"playerID": args.PlayerID,
	}, headers, &res, ctx)
	if err != nil {
		return nil, nil, err
	}
	switch {
	case status == http.StatusUnauthorized:
		return nil, resHeaders, UnauthorizedError
	case status != http.StatusOK || res.Data == nil:
		return nil, resHeaders, DeleteUserRoomError
	default:
		return res.Data, resHeaders, nil
	}
}

type CoreLock interface {
	Unlock(ctx context.Context) error
	Extend(ctx context.Context) (*time.Time, error)
	AddUserRoom(roomID types.RoomID, playerID types.PlayerID, ctx context.Context) (*types.UserRoom, error)
	DeleteUserRoom(playerID types.PlayerID, ctx context.Context) (*types.UserRoom, error)
}

type DefaultCoreLock struct {
	unlock         func(ctx context.Context) error
	extend         func(ctx context.Context) (*time.Time, error)
	addUserRoom    func(roomID types.RoomID, playerID types.PlayerID, ctx context.Context) (*types.UserRoom, error)
	deleteUserRoom func(playerID types.PlayerID, ctx context.Context) (*types.UserRoom, error)
}

func (d *DefaultCoreLock) Unlock(ctx context.Context) error {
	return d.unlock(ctx)
}

func (d *DefaultCoreLock) Extend(ctx context.Context) (*time.Time, error) {
	return d.extend(ctx)
}

func (d *DefaultCoreLock) AddUserRoom(roomID types.RoomID, playerID types.PlayerID, ctx context.Context) (*types.UserRoom, error) {
	return d.addUserRoom(roomID, playerID, ctx)
}

func (d *DefaultCoreLock) DeleteUserRoom(playerID types.PlayerID, ctx context.Context) (*types.UserRoom, error) {
	return d.deleteUserRoom(playerID, ctx)
}

func (c *PublicCoreAPI) LockUserRoomWithHeaders(args LockUserRoomArgs, headers map[string]string, ctx context.Context) (CoreLock, error) {
	res, _, err := c.LockUserRoomHeaders(args, headers, ctx)
	if err != nil {
		return nil, err
	}
	return &DefaultCoreLock{
		unlock: func(ctx context.Context) error {
			_, err := c.UnlockUserRoomHeaders(UnlockUserRoomArgs{
				Key: res.Key,
			}, nil, ctx)
			return err
		},
		extend: func(ctx context.Context) (*time.Time, error) {
			t, _, err := c.ExtendUserRoomHeaders(ExtendUserRoomArgs{
				Key: res.Key,
			}, nil, ctx)
			if err != nil {
				return nil, err
			}
			return &t.ExpiresAt, nil
		},
		addUserRoom: func(roomID types.RoomID, playerID types.PlayerID, ctx context.Context) (*types.UserRoom, error) {
			userRoom, _, err := c.AddUserRoomHeaders(AddUserRoomArgs{
				RoomID:   roomID,
				PlayerID: playerID,
				Key:      res.Key,
			}, nil, ctx)
			return userRoom, err
		},
		deleteUserRoom: func(playerID types.PlayerID, ctx context.Context) (*types.UserRoom, error) {
			userRoom, _, err := c.DeleteUserRoomHeaders(DeleteUserRoomArgs{
				PlayerID: playerID,
				Key:      res.Key,
			}, nil, ctx)
			return userRoom, err
		},
	}, nil
}

func (c *PublicCoreAPI) LockUserRoom(args LockUserRoomArgs, ctx context.Context) (CoreLock, error) {
	return c.LockUserRoomWithHeaders(args, nil, ctx)
}
