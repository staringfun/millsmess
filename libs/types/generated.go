// Mills Mess
// Licensed under the Mills Mess License Agreement
// See LICENSE.md in the root of this repository.

// Package types: Code generated; DO NOT EDIT;
package types

import "time"

type LockKey string

func NewLockKey(l LockKey) *LockKey {
	return &l
}
func (l LockKey) String() string {
	return string(l)
}
func (l LockKey) IsEmpty() bool {
	return l == ""
}
func (l LockKey) IsValid() bool {
	return l != ""
}

type Lock struct {
	Key       LockKey   `json:"key,omitempty"`
	ExpiresAt time.Time `json:"expiresAt,omitempty"`
}

func (l *Lock) IsValid() bool {
	if l == nil {
		return false
	}
	if !l.Key.IsValid() {
		return false
	}
	return true
}

type JoinResult string

func NewJoinResult(j JoinResult) *JoinResult {
	return &j
}
func (j JoinResult) String() string {
	return string(j)
}
func (j JoinResult) IsEmpty() bool {
	return j == ""
}

const (
	JoinResultFull     JoinResult = "full"
	JoinResultNotFound JoinResult = "not_found"
)

var AllJoinResult = []JoinResult{
	JoinResultFull,
	JoinResultNotFound,
}

func (j JoinResult) IsValid() bool {
	for _, jj := range AllJoinResult {
		if jj == j {
			return true
		}
	}
	return false
}

type LeaveReason string

func NewLeaveReason(l LeaveReason) *LeaveReason {
	return &l
}
func (l LeaveReason) String() string {
	return string(l)
}
func (l LeaveReason) IsEmpty() bool {
	return l == ""
}

const (
	LeaveReasonKicked LeaveReason = "kicked"
	LeaveReasonBanned LeaveReason = "banned"
)

var AllLeaveReason = []LeaveReason{
	LeaveReasonKicked,
	LeaveReasonBanned,
}

func (l LeaveReason) IsValid() bool {
	for _, ll := range AllLeaveReason {
		if ll == l {
			return true
		}
	}
	return false
}

type MessageKey int

func NewMessageKey(m MessageKey) *MessageKey {
	return &m
}
func (m MessageKey) Int() int {
	return int(m)
}
func (m MessageKey) IsEmpty() bool {
	return m == 0
}
func (m MessageKey) IsValid() bool {
	return m != 0
}

type Email string

func NewEmail(e Email) *Email {
	return &e
}
func (e Email) String() string {
	return string(e)
}
func (e Email) IsEmpty() bool {
	return e == ""
}
func (e Email) IsValid() bool {
	return e != ""
}

type Points int64

func NewPoints(p Points) *Points {
	return &p
}
func (p Points) Int64() int64 {
	return int64(p)
}
func (p Points) IsEmpty() bool {
	return p == 0
}
func (p Points) IsValid() bool {
	return p != 0
}

type UserID string

func NewUserID(u UserID) *UserID {
	return &u
}
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

func NewUsername(u Username) *Username {
	return &u
}
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

func NewUserRole(u UserRole) *UserRole {
	return &u
}
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

var AllUserRole = []UserRole{
	UserRoleGuest,
	UserRoleMember,
	UserRoleModerator,
	UserRoleAdmin,
}

func (u UserRole) IsValid() bool {
	for _, uu := range AllUserRole {
		if uu == u {
			return true
		}
	}
	return false
}

type User interface {
	IsValid() bool
	GetID() UserID
	SetID(UserID)
	GetUsername() Username
	SetUsername(Username)
	GetRole() UserRole
	SetRole(UserRole)
	GetPoints() Points
	SetPoints(Points)
	GetRooms() []Room
	SetRooms([]Room)
	GetCreatedAt() time.Time
	SetCreatedAt(time.Time)
	GetBannedAt() *time.Time
	SetBannedAt(*time.Time)
}
type BaseUser struct {
	ID        UserID     `json:"id,omitempty"`
	Username  Username   `json:"username,omitempty"`
	Role      UserRole   `json:"role,omitempty"`
	Points    Points     `json:"points,omitempty"`
	Rooms     []Room     `json:"rooms,omitempty"`
	CreatedAt time.Time  `json:"createdAt,omitempty"`
	BannedAt  *time.Time `json:"bannedAt,omitempty"`
}

func (u *BaseUser) IsValid() bool {
	if u == nil {
		return false
	}
	if !u.ID.IsValid() {
		return false
	}
	if !u.Username.IsValid() {
		return false
	}
	if !u.Role.IsValid() {
		return false
	}
	if !u.Points.IsValid() {
		return false
	}
	for _, uu := range u.Rooms {
		if !uu.IsValid() {
			return false
		}
	}
	return true
}
func (u *BaseUser) GetID() UserID {
	if u == nil {
		var uu UserID
		return uu
	}
	return u.ID
}
func (u *BaseUser) SetID(uu UserID) {
	if u == nil {
		return
	}
	u.ID = uu
}
func (u *BaseUser) GetUsername() Username {
	if u == nil {
		var uu Username
		return uu
	}
	return u.Username
}
func (u *BaseUser) SetUsername(uu Username) {
	if u == nil {
		return
	}
	u.Username = uu
}
func (u *BaseUser) GetRole() UserRole {
	if u == nil {
		var uu UserRole
		return uu
	}
	return u.Role
}
func (u *BaseUser) SetRole(uu UserRole) {
	if u == nil {
		return
	}
	u.Role = uu
}
func (u *BaseUser) GetPoints() Points {
	if u == nil {
		var uu Points
		return uu
	}
	return u.Points
}
func (u *BaseUser) SetPoints(uu Points) {
	if u == nil {
		return
	}
	u.Points = uu
}
func (u *BaseUser) GetRooms() []Room {
	if u == nil {
		var uu []Room
		return uu
	}
	return u.Rooms
}
func (u *BaseUser) SetRooms(uu []Room) {
	if u == nil {
		return
	}
	u.Rooms = uu
}
func (u *BaseUser) GetCreatedAt() time.Time {
	if u == nil {
		var uu time.Time
		return uu
	}
	return u.CreatedAt
}
func (u *BaseUser) SetCreatedAt(uu time.Time) {
	if u == nil {
		return
	}
	u.CreatedAt = uu
}
func (u *BaseUser) GetBannedAt() *time.Time {
	if u == nil {
		var uu *time.Time
		return uu
	}
	return u.BannedAt
}
func (u *BaseUser) SetBannedAt(uu *time.Time) {
	if u == nil {
		return
	}
	u.BannedAt = uu
}

type UserGuest struct {
	*BaseUser
}

func (u *UserGuest) IsValid() bool {
	if u == nil {
		return false
	}
	if !u.BaseUser.IsValid() {
		return false
	}
	return true
}

type UserMember struct {
	Email Email `json:"email,omitempty"`
	*BaseUser
}

func (u *UserMember) IsValid() bool {
	if u == nil {
		return false
	}
	if !u.BaseUser.IsValid() {
		return false
	}
	if !u.Email.IsValid() {
		return false
	}
	return true
}

type UserModerator struct {
	Email Email `json:"email,omitempty"`
	*BaseUser
}

func (u *UserModerator) IsValid() bool {
	if u == nil {
		return false
	}
	if !u.BaseUser.IsValid() {
		return false
	}
	if !u.Email.IsValid() {
		return false
	}
	return true
}

type UserAdmin struct {
	Email Email `json:"email,omitempty"`
	*BaseUser
}

func (u *UserAdmin) IsValid() bool {
	if u == nil {
		return false
	}
	if !u.BaseUser.IsValid() {
		return false
	}
	if !u.Email.IsValid() {
		return false
	}
	return true
}

type RoomID string

func NewRoomID(r RoomID) *RoomID {
	return &r
}
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
	Sessions  []Session `json:"sessions,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
}

func (r *Room) IsValid() bool {
	if r == nil {
		return false
	}
	if !r.ID.IsValid() {
		return false
	}
	if !r.OwnerID.IsValid() {
		return false
	}
	if !r.Owner.IsValid() {
		return false
	}
	for _, rr := range r.Sessions {
		if !rr.IsValid() {
			return false
		}
	}
	return true
}

type SocketRoomPlayerData struct {
	IsMicrophoneMuted bool `json:"isMicrophoneMuted,omitempty"`
	IsVideoMuted      bool `json:"isVideoMuted,omitempty"`
}

func (s *SocketRoomPlayerData) IsValid() bool {
	if s == nil {
		return false
	}
	return true
}

type SocketRoomPlayer struct {
	ID      PlayerID             `json:"id,omitempty"`
	Profile User                 `json:"profile,omitempty"`
	Data    SocketRoomPlayerData `json:"data,omitempty"`
}

func (s *SocketRoomPlayer) IsValid() bool {
	if s == nil {
		return false
	}
	if !s.ID.IsValid() {
		return false
	}
	if !s.Profile.IsValid() {
		return false
	}
	if !s.Data.IsValid() {
		return false
	}
	return true
}

type SocketRoom struct {
	ID      RoomID             `json:"id,omitempty"`
	Profile *Room              `json:"profile,omitempty"`
	Session SocketSession      `json:"session,omitempty"`
	Players []SocketRoomPlayer `json:"players,omitempty"`
}

func (s *SocketRoom) IsValid() bool {
	if s == nil {
		return false
	}
	if !s.ID.IsValid() {
		return false
	}
	if s.Profile != nil && !s.Profile.IsValid() {
		return false
	}
	if !s.Session.IsValid() {
		return false
	}
	for _, ss := range s.Players {
		if !ss.IsValid() {
			return false
		}
	}
	return true
}

type PlayerID string

func NewPlayerID(p PlayerID) *PlayerID {
	return &p
}
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

func (p *Player) IsValid() bool {
	if p == nil {
		return false
	}
	if !p.ID.IsValid() {
		return false
	}
	if p.RoomID != nil && !p.RoomID.IsValid() {
		return false
	}
	for _, pp := range p.Rooms {
		if !pp.IsValid() {
			return false
		}
	}
	return true
}

type Game string

func NewGame(g Game) *Game {
	return &g
}
func (g Game) String() string {
	return string(g)
}
func (g Game) IsEmpty() bool {
	return g == ""
}
func (g Game) IsValid() bool {
	return g != ""
}

type SessionID string

func NewSessionID(s SessionID) *SessionID {
	return &s
}
func (s SessionID) String() string {
	return string(s)
}
func (s SessionID) IsEmpty() bool {
	return s == ""
}
func (s SessionID) IsValid() bool {
	return s != ""
}

type SocketSessionPlayer struct {
	IsWin  bool   `json:"isWin,omitempty"`
	Points Points `json:"points,omitempty"`
	Data   any    `json:"data,omitempty"`
}

func (s *SocketSessionPlayer) IsValid() bool {
	if s == nil {
		return false
	}
	if !s.Points.IsValid() {
		return false
	}
	return true
}

type SocketSession interface {
	IsValid() bool
	GetID() SessionID
	SetID(SessionID)
	GetGame() Game
	SetGame(Game)
	GetPlayers() map[UserID]SocketSessionPlayer
	SetPlayers(map[UserID]SocketSessionPlayer)
	GetData() any
	SetData(any)
	GetCreatedAt() time.Time
	SetCreatedAt(time.Time)
	GetStartedAt() *time.Time
	SetStartedAt(*time.Time)
	GetFinishedAt() *time.Time
	SetFinishedAt(*time.Time)
}
type BaseSocketSession struct {
	ID         SessionID                      `json:"id,omitempty"`
	Game       Game                           `json:"game,omitempty"`
	Players    map[UserID]SocketSessionPlayer `json:"players,omitempty"`
	Data       any                            `json:"data,omitempty"`
	CreatedAt  time.Time                      `json:"createdAt,omitempty"`
	StartedAt  *time.Time                     `json:"startedAt,omitempty"`
	FinishedAt *time.Time                     `json:"finishedAt,omitempty"`
}

func (s *BaseSocketSession) IsValid() bool {
	if s == nil {
		return false
	}
	if !s.ID.IsValid() {
		return false
	}
	if !s.Game.IsValid() {
		return false
	}
	for _, ss := range s.Players {
		if !ss.IsValid() {
			return false
		}
	}
	return true
}
func (s *BaseSocketSession) GetID() SessionID {
	if s == nil {
		var ss SessionID
		return ss
	}
	return s.ID
}
func (s *BaseSocketSession) SetID(ss SessionID) {
	if s == nil {
		return
	}
	s.ID = ss
}
func (s *BaseSocketSession) GetGame() Game {
	if s == nil {
		var ss Game
		return ss
	}
	return s.Game
}
func (s *BaseSocketSession) SetGame(ss Game) {
	if s == nil {
		return
	}
	s.Game = ss
}
func (s *BaseSocketSession) GetPlayers() map[UserID]SocketSessionPlayer {
	if s == nil {
		var ss map[UserID]SocketSessionPlayer
		return ss
	}
	return s.Players
}
func (s *BaseSocketSession) SetPlayers(ss map[UserID]SocketSessionPlayer) {
	if s == nil {
		return
	}
	s.Players = ss
}
func (s *BaseSocketSession) GetData() any {
	if s == nil {
		var ss any
		return ss
	}
	return s.Data
}
func (s *BaseSocketSession) SetData(ss any) {
	if s == nil {
		return
	}
	s.Data = ss
}
func (s *BaseSocketSession) GetCreatedAt() time.Time {
	if s == nil {
		var ss time.Time
		return ss
	}
	return s.CreatedAt
}
func (s *BaseSocketSession) SetCreatedAt(ss time.Time) {
	if s == nil {
		return
	}
	s.CreatedAt = ss
}
func (s *BaseSocketSession) GetStartedAt() *time.Time {
	if s == nil {
		var ss *time.Time
		return ss
	}
	return s.StartedAt
}
func (s *BaseSocketSession) SetStartedAt(ss *time.Time) {
	if s == nil {
		return
	}
	s.StartedAt = ss
}
func (s *BaseSocketSession) GetFinishedAt() *time.Time {
	if s == nil {
		var ss *time.Time
		return ss
	}
	return s.FinishedAt
}
func (s *BaseSocketSession) SetFinishedAt(ss *time.Time) {
	if s == nil {
		return
	}
	s.FinishedAt = ss
}

type MV1PlayerMove struct {
	PlayerID       PlayerID     `json:"playerID,omitempty"`
	UserID         UserID       `json:"userID,omitempty"`
	JoinedRoom     *RoomID      `json:"joinedRoom,omitempty"`
	LeftRoom       *RoomID      `json:"leftRoom,omitempty"`
	LeftRoomReason *LeaveReason `json:"leftRoomReason,omitempty"`
}

func (m *MV1PlayerMove) IsValid() bool {
	if m == nil {
		return false
	}
	if !m.PlayerID.IsValid() {
		return false
	}
	if !m.UserID.IsValid() {
		return false
	}
	if m.JoinedRoom != nil && !m.JoinedRoom.IsValid() {
		return false
	}
	if m.LeftRoom != nil && !m.LeftRoom.IsValid() {
		return false
	}
	if m.LeftRoomReason != nil && !m.LeftRoomReason.IsValid() {
		return false
	}
	return true
}

type Region string

func NewRegion(r Region) *Region {
	return &r
}
func (r Region) String() string {
	return string(r)
}
func (r Region) IsEmpty() bool {
	return r == ""
}
func (r Region) IsValid() bool {
	return r != ""
}

type UserRoom struct {
	PlayerID PlayerID `json:"playerID,omitempty"`
	RoomID   RoomID   `json:"roomID,omitempty"`
	Region   Region   `json:"region,omitempty"`
}

func (u *UserRoom) IsValid() bool {
	if u == nil {
		return false
	}
	if !u.PlayerID.IsValid() {
		return false
	}
	if !u.RoomID.IsValid() {
		return false
	}
	if !u.Region.IsValid() {
		return false
	}
	return true
}

type Session struct {
	ID         SessionID       `json:"id,omitempty"`
	RoomID     *RoomID         `json:"roomID,omitempty"`
	Game       Game            `json:"game,omitempty"`
	Players    []SessionPlayer `json:"players,omitempty"`
	CreatedAt  time.Time       `json:"createdAt,omitempty"`
	StartedAt  *time.Time      `json:"startedAt,omitempty"`
	FinishedAt *time.Time      `json:"finishedAt,omitempty"`
}

func (s *Session) IsValid() bool {
	if s == nil {
		return false
	}
	if !s.ID.IsValid() {
		return false
	}
	if s.RoomID != nil && !s.RoomID.IsValid() {
		return false
	}
	if !s.Game.IsValid() {
		return false
	}
	for _, ss := range s.Players {
		if !ss.IsValid() {
			return false
		}
	}
	return true
}

type SessionPlayer struct {
	SessionID  SessionID  `json:"sessionID,omitempty"`
	UserID     UserID     `json:"userID,omitempty"`
	IsWin      bool       `json:"isWin,omitempty"`
	Points     Points     `json:"points,omitempty"`
	FinishedAt *time.Time `json:"finishedAt,omitempty"`
}

func (s *SessionPlayer) IsValid() bool {
	if s == nil {
		return false
	}
	if !s.SessionID.IsValid() {
		return false
	}
	if !s.UserID.IsValid() {
		return false
	}
	if !s.Points.IsValid() {
		return false
	}
	return true
}

type DisconnectReason string

func NewDisconnectReason(d DisconnectReason) *DisconnectReason {
	return &d
}
func (d DisconnectReason) String() string {
	return string(d)
}
func (d DisconnectReason) IsEmpty() bool {
	return d == ""
}

const (
	DisconnectReasonBanned DisconnectReason = "banned"
)

var AllDisconnectReason = []DisconnectReason{
	DisconnectReasonBanned,
}

func (d DisconnectReason) IsValid() bool {
	for _, dd := range AllDisconnectReason {
		if dd == d {
			return true
		}
	}
	return false
}

type MV1Disconnect struct {
	Reason DisconnectReason `json:"reason,omitempty"`
}

func (m *MV1Disconnect) IsValid() bool {
	if m == nil {
		return false
	}
	if !m.Reason.IsValid() {
		return false
	}
	return true
}

type MV1PlayerUpdate struct {
	Player         *Player      `json:"player,omitempty"`
	JoinResult     *JoinResult  `json:"joinResult,omitempty"`
	JoinedRoom     *SocketRoom  `json:"joinedRoom,omitempty"`
	LeftRoom       *RoomID      `json:"leftRoom,omitempty"`
	LeftRoomReason *LeaveReason `json:"leftRoomReason,omitempty"`
	Key            *MessageKey  `json:"key,omitempty"`
}

func (m *MV1PlayerUpdate) IsValid() bool {
	if m == nil {
		return false
	}
	if m.Player != nil && !m.Player.IsValid() {
		return false
	}
	if m.JoinResult != nil && !m.JoinResult.IsValid() {
		return false
	}
	if m.JoinedRoom != nil && !m.JoinedRoom.IsValid() {
		return false
	}
	if m.LeftRoom != nil && !m.LeftRoom.IsValid() {
		return false
	}
	if m.LeftRoomReason != nil && !m.LeftRoomReason.IsValid() {
		return false
	}
	if m.Key != nil && !m.Key.IsValid() {
		return false
	}
	return true
}

type MV1RoomJoin struct {
	RoomID RoomID      `json:"roomID,omitempty"`
	Key    *MessageKey `json:"key,omitempty"`
}

func (m *MV1RoomJoin) IsValid() bool {
	if m == nil {
		return false
	}
	if !m.RoomID.IsValid() {
		return false
	}
	if m.Key != nil && !m.Key.IsValid() {
		return false
	}
	return true
}

type MV1RoomLeave struct {
	RoomID RoomID      `json:"roomID,omitempty"`
	Key    *MessageKey `json:"key,omitempty"`
}

func (m *MV1RoomLeave) IsValid() bool {
	if m == nil {
		return false
	}
	if !m.RoomID.IsValid() {
		return false
	}
	if m.Key != nil && !m.Key.IsValid() {
		return false
	}
	return true
}

type MV1ProfilesUpdate struct {
	Rooms []RoomID `json:"rooms,omitempty"`
	Users []User   `json:"users,omitempty"`
}

func (m *MV1ProfilesUpdate) IsValid() bool {
	if m == nil {
		return false
	}
	for _, mm := range m.Rooms {
		if !mm.IsValid() {
			return false
		}
	}
	for _, mm := range m.Users {
		if !mm.IsValid() {
			return false
		}
	}
	return true
}

type MV1SessionUpdate struct {
	Session   SocketSession `json:"session,omitempty"`
	Timestamp time.Time     `json:"timestamp,omitempty"`
}

func (m *MV1SessionUpdate) IsValid() bool {
	if m == nil {
		return false
	}
	if !m.Session.IsValid() {
		return false
	}
	return true
}

type SocketMessageTypeCommand string

func NewSocketMessageTypeCommand(s SocketMessageTypeCommand) *SocketMessageTypeCommand {
	return &s
}
func (s SocketMessageTypeCommand) String() string {
	return string(s)
}
func (s SocketMessageTypeCommand) IsEmpty() bool {
	return s == ""
}

const (
	SocketMessageTypeCommandV1RoomLeave        SocketMessageTypeCommand = "v1:room.leave"
	SocketMessageTypeCommandV1RoomJoin         SocketMessageTypeCommand = "v1:room.join"
	SocketMessageTypeCommandV1Connect          SocketMessageTypeCommand = "v1:connect"
	SocketMessageTypeCommandV1ClientDisconnect SocketMessageTypeCommand = "v1:client.disconnect"
)

var AllSocketMessageTypeCommand = []SocketMessageTypeCommand{
	SocketMessageTypeCommandV1RoomLeave,
	SocketMessageTypeCommandV1RoomJoin,
	SocketMessageTypeCommandV1Connect,
	SocketMessageTypeCommandV1ClientDisconnect,
}

func (s SocketMessageTypeCommand) IsValid() bool {
	for _, ss := range AllSocketMessageTypeCommand {
		if ss == s {
			return true
		}
	}
	return false
}

type SocketMessageTypeEvent string

func NewSocketMessageTypeEvent(s SocketMessageTypeEvent) *SocketMessageTypeEvent {
	return &s
}
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

var AllSocketMessageTypeEvent = []SocketMessageTypeEvent{
	SocketMessageTypeEventV1PlayerUpdate,
	SocketMessageTypeEventV1Disconnect,
}

func (s SocketMessageTypeEvent) IsValid() bool {
	for _, ss := range AllSocketMessageTypeEvent {
		if ss == s {
			return true
		}
	}
	return false
}

type TopicName string

func NewTopicName(t TopicName) *TopicName {
	return &t
}
func (t TopicName) String() string {
	return string(t)
}
func (t TopicName) IsEmpty() bool {
	return t == ""
}

const (
	TopicNameProfilesUpdate TopicName = "profiles.update"
	TopicNameSessionUpdate  TopicName = "session.update"
)

var AllTopicName = []TopicName{
	TopicNameProfilesUpdate,
	TopicNameSessionUpdate,
}

func (t TopicName) IsValid() bool {
	for _, tt := range AllTopicName {
		if tt == t {
			return true
		}
	}
	return false
}
