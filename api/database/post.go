package database

import (
	"fmt"
	"time"

	"github.com/jackc/pgx"
)

// Post .
type Post struct {
	ID           uint64    `json:"id"`
	Profile      Profile   `json:"profile"`
	Content      string    `json:"content"`
	CreatedAt    time.Time `json:"created_at"`
	LikeCount    int       `json:"like_count"`
	CommentCount int       `json:"comment_count"`
	IsLiked      bool      `json:"is_liked"`
	Medias       []string  `json:"medias"`
}

// Hydrate .
func (p *Post) Hydrate(profileID uint64) {
	p.IsLiked = p.IsLikedByProfileID(profileID)
}

// IsLikedByProfileID .
func (p *Post) IsLikedByProfileID(profileID uint64) bool {
	err := pool.QueryRow("SELECT id FROM post_like WHERE profile_id = $1 AND post_id = $2", profileID, p.ID).Scan()
	if err == pgx.ErrNoRows {
		return false
	}
	return true
}

// PostCreate .
func PostCreate(profile *Profile, content string, medias []string) (*Post, error) {
	if medias == nil {
		medias = []string{}
	}
	post := &Post{}
	post.Profile = *profile
	post.Profile.Email = ""
	post.Profile.PasswordHash = ""
	post.Content = content
	post.Medias = medias
	err := pool.QueryRow("INSERT INTO post (profile_id,content,medias) VALUES ($1, $2, $3) RETURNING id, created_at;", profile.ID, content, medias).
		Scan(&post.ID, &post.CreatedAt)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func postLikedByProfileIDRange(start, end, profileID uint64) map[uint64]bool {
	out := map[uint64]bool{}
	query := `SELECT post_id FROM post_like WHERE post_id >= $1 AND post_id <= $2 AND profile_id = $3;`
	res, err := pool.Query(query, start, end, profileID)
	if err != nil {
		return out
	}
	defer res.Close()
	for res.Next() {
		var i = uint64(0)
		res.Scan(&i)
		out[i] = true
	}
	return out
}

// PostListIDs .
func PostListIDs(profile *Profile, from *uint64) ([]*Post, error) {
	out := []*Post{}
	var res *pgx.Rows
	var err error
	if from == nil {
		res, err = pool.Query(`
			SELECT
				post.id, post.content, post.created_at, post.like_count, post.comment_count, post_like.id, post.medias,
				profile.id, profile.username, profile.avatar
			FROM post
			LEFT JOIN post_like ON post.id = post_like.post_id AND post_like.profile_id = $1
			LEFT JOIN profile ON post.profile_id = profile.id
			WHERE post.profile_id = $1 OR post.profile_id IN (SELECT friendship.profile_id_friend FROM friendship WHERE friendship.profile_id_owner = $1)
			ORDER BY post.id DESC LIMIT 10;
		`, profile.ID)
	} else {
		res, err = pool.Query(`
			SELECT
				post.id, post.content, post.created_at, post.like_count, post.comment_count, post_like.id, post.medias,
				profile.id, profile.username, profile.avatar
			FROM post
			LEFT JOIN post_like ON post.id = post_like.post_id AND post_like.profile_id = $1
			LEFT JOIN profile ON post.profile_id = profile.id
			WHERE post.id < $2 AND (post.profile_id = $1 OR post.profile_id IN (SELECT friendship.profile_id_friend FROM friendship WHERE friendship.profile_id_owner = $1) )
			ORDER BY post.id DESC LIMIT 10;
		`, profile.ID, *from)
	}
	if err == pgx.ErrNoRows {
		return out, nil
	} else if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer res.Close()
	for res.Next() {
		tmp := &Post{}
		var postLikeID *uint64
		res.Scan(
			&tmp.ID, &tmp.Content, &tmp.CreatedAt, &tmp.LikeCount, &tmp.CommentCount, &postLikeID, &tmp.Medias,
			&tmp.Profile.ID, &tmp.Profile.Username, &tmp.Profile.Avatar,
		)
		if tmp.Medias == nil {
			tmp.Medias = []string{}
		}
		tmp.IsLiked = postLikeID != nil
		out = append(out, tmp)
	}

	return out, nil
}

// PostGetByID .
func PostGetByID(id uint64) (*Post, error) {
	p := &Post{}
	var postLikeID *uint64
	err := pool.QueryRow(`
		SELECT
			post.id, post.content, post.created_at, post.like_count, post_like.id, post.medias,
			profile.id, profile.username, profile.avatar
		FROM post
		LEFT JOIN post_like ON post.id = post_like.post_id AND post_like.profile_id = $1
		LEFT JOIN profile ON post.profile_id = profile.id
		WHERE post.id = $1
	`, id).Scan(
		&p.ID, &p.Content, &p.CreatedAt, &p.LikeCount, &postLikeID, &p.Medias,
		&p.Profile.ID, &p.Profile.Username, &p.Profile.Avatar,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	if p.Medias == nil {
		p.Medias = []string{}
	}
	p.IsLiked = postLikeID != nil
	return p, nil
}

// PostLikeByID .
func PostLikeByID(id, userID uint64) error {
	_, err := pool.Exec("INSERT INTO post_like (post_id, profile_id) VALUES ($1, $2) ON CONFLICT DO NOTHING", id, userID)
	return err
}

// PostUnlikeByID .
func PostUnlikeByID(id, userID uint64) error {
	_, err := pool.Exec("DELETE FROM post_like WHERE post_id = $1 AND profile_id = $2", id, userID)
	return err
}
