// Mills Mess
// Licensed under the Mills Mess License Agreement
// See LICENSE.md in the root of this repository.

package internal_core_api

import (
	"context"
	"github.com/staringfun/millsmess/libs/public-core-api"
	"github.com/staringfun/millsmess/libs/types"
)

type InternalCoreAPI struct {
	*public_core_api.PublicCoreAPI
}

func (c *InternalCoreAPI) AppendUserHeader(id types.UserID, headers map[string]string) {
	headers["X-User-ID"] = id.String()
}

func (c *InternalCoreAPI) LockUserRoom(id types.UserID, ctx context.Context) (public_core_api.CoreLock, error) {
	headers := map[string]string{}
	c.AppendUserHeader(id, headers)
	return c.LockUserRoomWithHeaders(public_core_api.LockUserRoomArgs{}, headers, ctx)
}

func (c *InternalCoreAPI) FetchUserRooms(id types.UserID, ctx context.Context) ([]types.UserRoom, error) {
	headers := map[string]string{}
	c.AppendUserHeader(id, headers)
	res, _, err := c.FetchUserRoomHeaders(public_core_api.FetchUserRoomsArgs{}, headers, ctx)
	return res, err
}
