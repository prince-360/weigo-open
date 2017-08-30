package database

import (
	"log"
	"time"
)

// Comment .
type Comment struct {
	ID        uint64    `json:"id"`
	PostID    uint64    `json:"post_id"`
	Profile   Profile   `json:"profile"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// CommentCreate .
func CommentCreate(postID, profileID uint64, content string) (uint64, error) {
	var outID = uint64(0)
	err := pool.QueryRow("INSERT INTO post_comment (post_id, profile_id,content) VALUES ($1, $2, $3) RETURNING id;", postID, profileID, content).Scan(&outID)
	if err != nil {
		return 0, err
	}
	return outID, nil
}

// CommentListByPostID .
func CommentListByPostID(postID uint64) ([]Comment, error) {
	var out = []Comment{}
	row, err := pool.Query(`
		SELECT
			post_comment.id, post_comment.post_id, post_comment.content, post_comment.created_at,
			profile.id, profile.username, profile.avatar
		FROM post_comment
		LEFT JOIN profile ON post_comment.profile_id = profile.id
		WHERE post_id = $1
		ORDER BY post_comment.id ASC
	`, postID)
	if err != nil {
		log.Println(err)
		return out, err
	}
	defer row.Close()
	for row.Next() {
		tmp := Comment{}
		err := row.Scan(
			&tmp.ID, &tmp.PostID, &tmp.Content, &tmp.CreatedAt,
			&tmp.Profile.ID, &tmp.Profile.Username, &tmp.Profile.Avatar,
		)
		if err != nil {
			break
		}
		out = append(out, tmp)
	}
	return out, nil
}
