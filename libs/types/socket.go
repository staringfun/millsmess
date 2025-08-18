// Mills Mess
// Licensed under the Mills Mess License Agreement
// See LICENSE.md in the root of this repository.

package types

type SocketEventV1PlayerUpdateData struct {
	Player *Player `json:"player,omitempty"`
}

type SocketDisconnectReason string

const (
	SocketDisconnectReasonBan = SocketDisconnectReason("ban")
)

type SocketV1DisconnectDTO struct {
	Reason SocketDisconnectReason `json:"reason"`
}

type SocketData = any

type SocketDataConnect struct {
}
