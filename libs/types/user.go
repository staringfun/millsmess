// Mills Mess
// Licensed under the Mills Mess License Agreement
// See LICENSE.md in the root of this repository.

package types

import "time"

type User interface {
	GetID() UserID
	SetID(UserID)

	GetUsername() Username
	SetUsername(Username)

	GetRole() UserRole
	SetRole(UserRole)

	GetBannedAt() *time.Time
	SetBannedAt(*time.Time)
}

type BaseUser struct {
	ID       UserID     `json:"id"`
	Username Username   `json:"username"`
	Role     UserRole   `json:"role"`
	BannedAt *time.Time `json:"bannedAt"`
}

func (u *BaseUser) GetID() UserID {
	return u.ID
}

func (u *BaseUser) SetID(id UserID) {
	u.ID = id
}

func (u *BaseUser) GetUsername() Username {
	return u.Username
}

func (u *BaseUser) SetUsername(username Username) {
	u.Username = username
}

func (u *BaseUser) GetRole() UserRole {
	return u.Role
}

func (u *BaseUser) SetRole(role UserRole) {
	u.Role = role
}

func (u *BaseUser) GetBannedAt() *time.Time {
	return u.BannedAt
}

func (u *BaseUser) SetBannedAt(bannedAt *time.Time) {
	u.BannedAt = bannedAt
}
