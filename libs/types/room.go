// Mills Mess
// Licensed under the Mills Mess License Agreement
// See LICENSE.md in the root of this repository.

package types

type RoomID string

func (id RoomID) IsEmpty() bool {
	return id == ""
}

func (id RoomID) String() string {
	return string(id)
}
