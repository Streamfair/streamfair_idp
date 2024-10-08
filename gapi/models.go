package gapi


import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserSvcSession struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	RefreshToken string    `json:"refresh_token"`
	UserAgent    string    `json:"user_agent"`
	ClientIp     string    `json:"client_ip"`
	IsBlocked    bool      `json:"is_blocked"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
}

type UserSvcUser struct {
	ID                int64       `json:"id"`
	Username          string      `json:"username"`
	FullName          string      `json:"full_name"`
	Email             string      `json:"email"`
	PasswordHash      string      `json:"password_hash"`
	PasswordSalt      string      `json:"password_salt"`
	CountryCode       string      `json:"country_code"`
	RoleID            pgtype.Int8 `json:"role_id"`
	Status            pgtype.Text `json:"status"`
	LastLoginAt       time.Time   `json:"last_login_at"`
	UsernameChangedAt time.Time   `json:"username_changed_at"`
	EmailChangedAt    time.Time   `json:"email_changed_at"`
	PasswordChangedAt time.Time   `json:"password_changed_at"`
	CreatedAt         time.Time   `json:"created_at"`
	UpdatedAt         time.Time   `json:"updated_at"`
}

type TokenSvcRefreshToken struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Token     string    `json:"token"`
	Revoked   bool      `json:"revoked"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

type TokenSvcToken struct {
	ID        int64       `json:"id"`
	UserID    int64       `json:"user_id"`
	TokenType pgtype.Text `json:"token_type"`
	Token     string      `json:"token"`
	Revoked   bool        `json:"revoked"`
	ExpiresAt time.Time   `json:"expires_at"`
	CreatedAt time.Time   `json:"created_at"`
}

type SessionSvcSession struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	RefreshToken string    `json:"refresh_token"`
	UserAgent    string    `json:"user_agent"`
	ClientIp     string    `json:"client_ip"`
	IsBlocked    bool      `json:"is_blocked"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
}


type UserSvcAccount struct {
	ID          int64       `json:"id"`
	AccountType int32       `json:"account_type"`
	Owner       string      `json:"owner"`
	AvatarUri   pgtype.Text `json:"avatar_uri"`
	Plays       int64       `json:"plays"`
	Likes       int64       `json:"likes"`
	Follows     int64       `json:"follows"`
	Shares      int64       `json:"shares"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

type UserSvcAccountType struct {
	ID          int32     `json:"id"`
	Type        string    `json:"type"`
	Permissions []byte    `json:"permissions"`
	IsArtist    bool      `json:"is_artist"`
	IsProducer  bool      `json:"is_producer"`
	IsWriter    bool      `json:"is_writer"`
	IsLabel     bool      `json:"is_label"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UserSvcAccountTypesAccount struct {
	AccountTypesID int32 `json:"AccountTypes_id"`
	AccountsID     int64 `json:"Accounts_id"`
}

type AccessCtrlSvcRole struct {
	ID          int64  `json:"id"`
	RoleName    string `json:"role_name"`
	Permissions []byte `json:"permissions"`
}

type AccessCtrlSvcUserRole struct {
	UserID    int64     `json:"user_id"`
	RoleID    int64     `json:"role_id"`
	CreatedAt time.Time `json:"created_at"`
}
