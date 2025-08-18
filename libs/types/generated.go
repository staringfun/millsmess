// Mills Mess
// Licensed under the Mills Mess License Agreement
// See LICENSE.md in the root of this repository.

// Package types: Code generated; DO NOT EDIT;
package types

import "time"

type Email string

func (e Email) String() string {
	return string(e)
}
func (e Email) IsEmpty() bool {
	return e == ""
}
func (e Email) IsValid() bool {
	return e != ""
}

type UserID string

func (u UserID) String() string {
	return string(u)
}
func (u UserID) IsEmpty() bool {
	return u == ""
}
func (u UserID) IsValid() bool {
	return u != ""
}

type Username string

func (u Username) String() string {
	return string(u)
}
func (u Username) IsEmpty() bool {
	return u == ""
}

var UsernameRunes = []rune{
	'A',
	'B',
	'C',
	'D',
	'E',
	'F',
	'G',
	'H',
	'I',
	'J',
	'K',
	'L',
	'M',
	'N',
	'O',
	'P',
	'Q',
	'R',
	'S',
	'T',
	'U',
	'V',
	'W',
	'X',
	'Y',
	'Z',
	'a',
	'b',
	'c',
	'd',
	'e',
	'f',
	'g',
	'h',
	'i',
	'j',
	'k',
	'l',
	'm',
	'n',
	'o',
	'p',
	'q',
	'r',
	's',
	't',
	'u',
	'v',
	'w',
	'x',
	'y',
	'z',
	'0',
	'1',
	'2',
	'3',
	'4',
	'5',
	'6',
	'7',
	'8',
	'9',
	'.',
	'_',
	'-',
}

const UsernameMinLength = 2
const UsernameMaxLength = 32

func (u Username) IsValid() bool {
	if len(u) < UsernameMinLength {
		return false
	}
	if len(u) > UsernameMaxLength {
		return false
	}
outer:
	for _, uu := range []rune(u) {
		for _, uuu := range UsernameRunes {
			if uu == uuu {
				continue outer
			}
		}
		return false
	}
	return u != ""
}

type UserRole string

func (u UserRole) String() string {
	return string(u)
}
func (u UserRole) IsEmpty() bool {
	return u == ""
}

const (
	UserRoleGuest     UserRole = "guest"
	UserRoleMember    UserRole = "member"
	UserRoleModerator UserRole = "moderator"
	UserRoleAdmin     UserRole = "admin"
)

var AllUserRole = []UserRole{UserRoleGuest, UserRoleMember, UserRoleModerator, UserRoleAdmin}

func (u UserRole) IsValid() bool {
	for _, uu := range AllUserRole {
		if uu == u {
			return true
		}
	}
	return false
}

type User interface {
	GetID() UserID
	SetID(UserID)
	GetUsername() Username
	SetUsername(Username)
	GetRole() UserRole
	SetRole(UserRole)
	GetCreatedAt() time.Time
	SetCreatedAt(time.Time)
	GetBannedAt() *time.Time
	SetBannedAt(*time.Time)
}
type BaseUser struct {
	ID        UserID     `json:"id,omitempty"`
	Username  Username   `json:"username,omitempty"`
	Role      UserRole   `json:"role,omitempty"`
	CreatedAt time.Time  `json:"createdAt,omitempty"`
	BannedAt  *time.Time `json:"bannedAt,omitempty"`
}

func (u *BaseUser) GetID() UserID {
	return u.ID
}
func (u *BaseUser) SetID(v UserID) {
	u.ID = v
}
func (u *BaseUser) GetUsername() Username {
	return u.Username
}
func (u *BaseUser) SetUsername(v Username) {
	u.Username = v
}
func (u *BaseUser) GetRole() UserRole {
	return u.Role
}
func (u *BaseUser) SetRole(v UserRole) {
	u.Role = v
}
func (u *BaseUser) GetCreatedAt() time.Time {
	return u.CreatedAt
}
func (u *BaseUser) SetCreatedAt(v time.Time) {
	u.CreatedAt = v
}
func (u *BaseUser) GetBannedAt() *time.Time {
	return u.BannedAt
}
func (u *BaseUser) SetBannedAt(v *time.Time) {
	u.BannedAt = v
}

type UserGuest struct {
	*BaseUser
}

type UserMember struct {
	Email Email `json:"email,omitempty"`
	*BaseUser
}

type UserModerator struct {
	Email Email `json:"email,omitempty"`
	*BaseUser
}

type UserAdmin struct {
	Email Email `json:"email,omitempty"`
	*BaseUser
}

type RoomID string

func (r RoomID) String() string {
	return string(r)
}
func (r RoomID) IsEmpty() bool {
	return r == ""
}
func (r RoomID) IsValid() bool {
	return r != ""
}

type Room struct {
	ID        RoomID    `json:"id,omitempty"`
	OwnerID   UserID    `json:"ownerID,omitempty"`
	Owner     User      `json:"owner,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
}

type SocketRoomPlayer struct {
	ID      PlayerID `json:"id,omitempty"`
	Profile User     `json:"profile,omitempty"`
}

type SocketRoom struct {
	ID      RoomID             `json:"id,omitempty"`
	Profile *Room              `json:"profile,omitempty"`
	Players []SocketRoomPlayer `json:"players,omitempty"`
}

type PlayerID string

func (p PlayerID) String() string {
	return string(p)
}
func (p PlayerID) IsEmpty() bool {
	return p == ""
}
func (p PlayerID) IsValid() bool {
	return p != ""
}

type Player struct {
	ID     PlayerID     `json:"id,omitempty"`
	RoomID *RoomID      `json:"roomID,omitempty"`
	Rooms  []SocketRoom `json:"rooms,omitempty"`
}

type DisconnectReason string

func (d DisconnectReason) String() string {
	return string(d)
}
func (d DisconnectReason) IsEmpty() bool {
	return d == ""
}

const (
	DisconnectReasonBanned DisconnectReason = "banned"
)

var AllDisconnectReason = []DisconnectReason{DisconnectReasonBanned}

func (d DisconnectReason) IsValid() bool {
	for _, dd := range AllDisconnectReason {
		if dd == d {
			return true
		}
	}
	return false
}

type MV1Connect struct {
	RoomID *RoomID `json:"roomID,omitempty"`
}

type MV1Disconnect struct {
	Reason DisconnectReason `json:"reason,omitempty"`
}

type MV1PlayerUpdate struct {
	Player *Player `json:"player,omitempty"`
}

type SocketMessageTypeCommand string

func (s SocketMessageTypeCommand) String() string {
	return string(s)
}
func (s SocketMessageTypeCommand) IsEmpty() bool {
	return s == ""
}

const (
	SocketMessageTypeCommandV1Connect SocketMessageTypeCommand = "v1:connect"
)

var AllSocketMessageTypeCommand = []SocketMessageTypeCommand{SocketMessageTypeCommandV1Connect}

func (s SocketMessageTypeCommand) IsValid() bool {
	for _, ss := range AllSocketMessageTypeCommand {
		if ss == s {
			return true
		}
	}
	return false
}

type SocketMessageTypeEvent string

func (s SocketMessageTypeEvent) String() string {
	return string(s)
}
func (s SocketMessageTypeEvent) IsEmpty() bool {
	return s == ""
}

const (
	SocketMessageTypeEventV1PlayerUpdate SocketMessageTypeEvent = "v1:player.update"
	SocketMessageTypeEventV1Disconnect   SocketMessageTypeEvent = "v1:disconnect"
)

var AllSocketMessageTypeEvent = []SocketMessageTypeEvent{SocketMessageTypeEventV1PlayerUpdate, SocketMessageTypeEventV1Disconnect}

func (s SocketMessageTypeEvent) IsValid() bool {
	for _, ss := range AllSocketMessageTypeEvent {
		if ss == s {
			return true
		}
	}
	return false
}
