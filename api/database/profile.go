package database

import (
	"fmt"

	"github.com/jackc/pgx"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

// Profile .
type Profile struct {
	ID           uint64 `json:"id"`
	Email        string `json:"email,omitempty"`
	Username     string `json:"username"`
	PasswordHash string `json:"-"`
	Avatar       string `json:"avatar"`
}

// Authenticate .
func (p *Profile) Authenticate(rawPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(p.PasswordHash), []byte(rawPassword))
	return err == nil
}

// ChangeAvatar .
func (p *Profile) ChangeAvatar(pictureID string) error {
	_, err := pool.Exec("UPDATE profile SET avatar = $2 WHERE id = $1", p.ID, pictureID)
	if err != nil {
		return err
	}
	return nil
}

// ChangeGoogleID .
func (p *Profile) ChangeGoogleID(gid string) error {
	_, err := pool.Exec("UPDATE profile SET google_id = $2 WHERE id = $1", p.ID, gid)
	if err != nil {
		return err
	}
	return nil
}

// ProfileCreateFromGitbub .
func ProfileCreateFromGitbub(username, githubID string) (*Profile, error) {
	p := &Profile{}
	err := pool.QueryRow(
		"INSERT INTO profile (username, github_id, email, google_id) VALUES ($1, $2, $3, $4) RETURNING id",
		username, githubID, uuid.NewV4().String(), uuid.NewV4().String(),
	).Scan(&p.ID)
	if err != nil {
		pgE := err.(pgx.PgError)
		if pgE.Code == "23505" {
			return nil, ErrDuplicated
		}
		return nil, err
	}
	p.Username = username
	return p, nil
}

// ProfileCreateFromGoogle .
func ProfileCreateFromGoogle(username, googleID string) (*Profile, error) {
	p := &Profile{}
	err := pool.QueryRow(
		"INSERT INTO profile (username, google_id, email, github_id) VALUES ($1, $2, $3, $4) RETURNING id",
		username, googleID, uuid.NewV4().String(), uuid.NewV4().String(),
	).Scan(&p.ID)
	if err != nil {
		pgE := err.(pgx.PgError)
		if pgE.Code == "23505" {
			return nil, ErrDuplicated
		}
		return nil, err
	}
	p.Username = username
	return p, nil
}

// ProfileCreate .
func ProfileCreate(email, username, passwordRaw string) (*Profile, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(passwordRaw), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	p := &Profile{}
	err = pool.QueryRow(
		"INSERT INTO profile (email, username, password, google_id, github_id) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		email, username, string(hashedPassword), uuid.NewV4().String(), uuid.NewV4().String(),
	).Scan(&p.ID)
	if err != nil {
		pgE := err.(pgx.PgError)
		if pgE.Code == "23505" {
			return nil, ErrDuplicated
		}
		return nil, err
	}
	p.Email = email
	p.Username = username
	return p, nil
}

// ProfileGetByEmail .
func ProfileGetByEmail(email string) (*Profile, error) {
	p := &Profile{}
	err := pool.QueryRow("SELECT  id, email, username, password, avatar FROM profile WHERE email = $1", email).Scan(&p.ID, &p.Email, &p.Username, &p.PasswordHash, &p.Avatar)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return p, nil
}

// ProfileGetByGithubID .
func ProfileGetByGithubID(gid string) (*Profile, error) {
	p := &Profile{}
	err := pool.QueryRow("SELECT  id,username FROM profile WHERE github_id = $1", gid).Scan(&p.ID, &p.Username)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return p, nil
}

// ProfileGetByGoogleID .
func ProfileGetByGoogleID(gid string) (*Profile, error) {
	p := &Profile{}
	err := pool.QueryRow("SELECT  id, username FROM profile WHERE google_id = $1", gid).Scan(&p.ID, &p.Username)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return p, nil
}

// ProfileGetByID .
func ProfileGetByID(id uint64) (*Profile, error) {
	p := &Profile{}
	err := pool.QueryRow("SELECT id, email, username, password, avatar FROM profile WHERE id = $1", id).Scan(&p.ID, &p.Email, &p.Username, &p.PasswordHash, &p.Avatar)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return p, nil
}

// ProfileGetByUsername .
func ProfileGetByUsername(unsername string) (*Profile, error) {
	p := &Profile{}
	err := pool.QueryRow("SELECT id, email, username, password, avatar FROM profile WHERE username = $1", unsername).Scan(&p.ID, &p.Email, &p.Username, &p.PasswordHash, &p.Avatar)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return p, nil
}

// ProfileSearchByName .
func ProfileSearchByName(p *Profile, name string) ([]Profile, error) {
	out := []Profile{}
	name = "%" + name + "%"
	row, err := pool.Query("SELECT id, username, avatar FROM profile WHERE username LIKE $1 AND id != $2 ORDER BY username LIMIT 5", name, p.ID)
	if err != nil {
		return out, nil
	}
	defer row.Close()
	for row.Next() {
		tmp := Profile{}
		row.Scan(&tmp.ID, &tmp.Username, &tmp.Avatar)
		out = append(out, tmp)
	}
	return out, nil
}

// ProfileAddFriendByID .
func ProfileAddFriendByID(p *Profile, friendID uint64) error {
	_, err := pool.Exec("INSERT INTO friendship (profile_id_owner, profile_id_friend) VALUES ($1, $2) ON CONFLICT DO NOTHING", p.ID, friendID)
	if err != nil {
		return err
	}
	return nil
}

// ProfileRemoveFriendByID .
func ProfileRemoveFriendByID(p *Profile, friendID uint64) error {
	_, err := pool.Exec("DELETE FROM friendship WHERE profile_id_owner = $1 AND profile_id_friend = $2", p.ID, friendID)
	if err != nil {
		return err
	}
	return nil
}

// ProfileListFriends .
func ProfileListFriends(p *Profile) ([]Profile, error) {
	out := []Profile{}
	row, err := pool.Query(`
		SELECT profile.id, profile.username, profile.avatar
		FROM friendship
		LEFT JOIN profile ON friendship.profile_id_friend = profile.id
		WHERE friendship.profile_id_owner = $1
		ORDER BY username
	`, p.ID)
	if err != nil {
		fmt.Println(err)
		return out, nil
	}
	defer row.Close()
	for row.Next() {
		tmp := Profile{}
		row.Scan(&tmp.ID, &tmp.Username, &tmp.Avatar)
		out = append(out, tmp)
	}
	return out, nil
}

// ProfileUpdatePassword .
func ProfileUpdatePassword(p *Profile, newPassword string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = pool.Exec("UPDATE profile SET password = $2 WHERE id = $1", p.ID, string(hashedPassword))
	if err != nil {
		return err
	}
	return nil
}

// ProfileUpdateUsername .
func ProfileUpdateUsername(p *Profile, username string) error {
	_, err := pool.Exec("UPDATE profile SET username = $2 WHERE id = $1", p.ID, username)
	if err != nil {
		pgerr := err.(pgx.PgError)
		if pgerr.Code == "23505" {
			return ErrDuplicated
		}
		return err
	}
	return nil
}
