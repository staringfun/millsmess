// Mills Mess
// Licensed under the Mills Mess License Agreement
// See LICENSE.md in the root of this repository.

package types

type RoomID string

func (r RoomID) IsEmpty() bool {
	return r == ""
}
