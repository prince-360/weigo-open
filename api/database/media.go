package database

import (
    "time"
)

// Media .
type Media struct {
    ID string `json:"id"`
    CreatedAt time.Time `json:"created_at"`
    Profile Profile `json:"profile"`
}

// MediaCreated .
func MediaCreated(p *Profile, id string) (*Media, error) {
    m := &Media{}
    m.Profile = *p
    m.Profile.Email = ""
    m.Profile.PasswordHash = ""
    err := pool.QueryRow("INSERT INTO media (id, profile_id) VALUES ($1, $2) ON CONFLICT DO NOTHING RETURNING id, created_at", id, p.ID).
    Scan(&m.ID, &m.CreatedAt)
    if err != nil {
        return nil, err
    }
    return m, nil
}
