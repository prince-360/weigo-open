package database

import (
	"strings"
	"time"

	"github.com/jackc/pgx"
	"github.com/satori/go.uuid"
)

// OauthProfile .
type OauthProfile struct {
	Key       string    `json:"key"`
	Type      string    `json:"-"`
	AccountID string    `json:"-"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_ad"`
}

// Insert .
func (OauthProfile) Insert(op *OauthProfile) error {
	key := uuid.NewV4().String() + uuid.NewV4().String()
	key = strings.Replace(key, "-", "", -1)
	op.Key = key
	err := pool.QueryRow(
		`INSERT INTO oauth_profile (key, type, account_id, username, email) VALUES ($1,$2,$3,$4,$5) RETURNING created_at`,
		op.Key, op.Type, op.AccountID, op.Username, op.Email,
	).Scan(&op.CreatedAt)
	return err
}

// GetByKey .
func (OauthProfile) GetByKey(key string) (*OauthProfile, error) {
	op := &OauthProfile{}
	op.Key = key
	err := pool.QueryRow(
		`SELECT type, account_id, username, email, created_at FROM oauth_profile WHERE key = $1`,
		key,
	).Scan(&op.Type, &op.AccountID, &op.Username, &op.Email, &op.CreatedAt)
	if err != nil && err == pgx.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return op, nil
}
