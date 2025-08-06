// Mills Mess
// Licensed under the Mills Mess License Agreement
// See LICENSE.md in the root of this repository.

package types

type SocketEvent string

const (
	SocketEventConnect    = SocketEvent("connect")
	SocketEventDisconnect = SocketEvent("connect")
)

type SocketData any

type SocketDataConnect struct {
}
