# Weigo

Weigo is a testing project using Golang and Vuejs-2 that copycats twitter. You can test it directly [here](https://weigo.tuxlinuxien.com/).

Don't use it in production, however you are free to fork it and modify it as much
as you need.

## Install

```
cd $GOPATH/src
git clone https://github.com/tuxlinuxien/weigo-open.git ./weigo
```

## SQL

Weigo is using Postgres 9.5+.

```sql
DROP SCHEMA public CASCADE;
CREATE SCHEMA public;

CREATE TABLE oauth_profile (
    key TEXT PRIMARY KEY,
    type TEXT,
    email TEXT,
    username TEXT,
    account_id TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now(),
    used BOOLEAN DEFAULT False
);

CREATE TABLE profile (
    id SERIAL PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    email TEXT UNIQUE NOT NULL,
    google_id TEXT UNIQUE NOT NULL,
    github_id TEXT UNIQUE NOT NULL,
    password TEXT default '',
    avatar TEXT DEFAULT '',
    description TEXT DEFAULT ''
);

CREATE TABLE post (
    id SERIAL PRIMARY KEY,
    profile_id INTEGER REFERENCES profile(id) ON DELETE CASCADE,
    content TEXT DEFAULT '',
    created_at TIMESTAMPTZ DEFAULT now(),
    like_count INTEGER DEFAULT 0,
    medias TEXT[] DEFAULT ARRAY[]::TEXT[],
    comment_count INTEGER DEFAULT 0
);

CREATE TABLE post_comment (
    id SERIAL PRIMARY KEY,
    post_id INTEGER REFERENCES post(id) ON DELETE CASCADE,
    profile_id INTEGER REFERENCES profile(id) ON DELETE CASCADE,
    content TEXT DEFAULT '',
    created_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE post_like (
    id SERIAL PRIMARY KEY,
    profile_id INTEGER REFERENCES profile(id) ON DELETE CASCADE,
    post_id INTEGER REFERENCES post(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT now(),
    UNIQUE (profile_id, post_id)
);

CREATE TABLE friendship (
    id SERIAL PRIMARY KEY,
    profile_id_owner INTEGER REFERENCES profile(id) ON DELETE CASCADE,
    profile_id_friend INTEGER REFERENCES profile(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT now(),
    UNIQUE (profile_id_owner, profile_id_friend)
);

CREATE TABLE media (
    id TEXT PRIMARY KEY UNIQUE,
    profile_id INTEGER REFERENCES profile(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT now()
);

CREATE OR REPLACE FUNCTION post_like_count_procedure() RETURNS TRIGGER AS $$
DECLARE
    post_like_count INTEGER;
BEGIN
    IF ( TG_OP = 'INSERT' ) THEN
        SELECT INTO post_like_count COUNT(*) FROM post_like WHERE post_id = NEW.post_id;
        UPDATE post SET like_count = post_like_count WHERE id = NEW.post_id;
        RETURN NEW;
    ELSIF ( TG_OP = 'DELETE' ) THEN
        SELECT INTO post_like_count COUNT(*) FROM post_like WHERE post_id = OLD.post_id;
        UPDATE post SET like_count = post_like_count WHERE id = OLD.post_id;
        RETURN OLD;
    END IF;
    RETURN NULL;
END; $$
LANGUAGE plpgsql;

CREATE TRIGGER post_like_count_trigger AFTER INSERT OR DELETE ON post_like
FOR EACH ROW EXECUTE PROCEDURE post_like_count_procedure();

CREATE OR REPLACE FUNCTION post_comment_count_procedure() RETURNS TRIGGER AS $$
DECLARE
    post_comment_count INTEGER;
BEGIN
    IF ( TG_OP = 'INSERT' ) THEN
        SELECT INTO post_comment_count COUNT(*) FROM post_comment WHERE post_id = NEW.post_id;
        UPDATE post SET comment_count = post_comment_count WHERE id = NEW.post_id;
        RETURN NEW;
    ELSIF ( TG_OP = 'DELETE' ) THEN
        SELECT INTO post_comment_count COUNT(*) FROM post_comment WHERE post_id = OLD.post_id;
        UPDATE post SET comment_count = post_comment_count WHERE id = OLD.post_id;
        RETURN OLD;
    END IF;
    RETURN NULL;
END; $$
LANGUAGE plpgsql;

CREATE TRIGGER post_comment_count_trigger AFTER INSERT OR DELETE ON post_comment
FOR EACH ROW EXECUTE PROCEDURE post_comment_count_procedure();
```

## NGINX conf

```conf
server {
    listen listen 443 ssl http2;
    server_name weigo.tuxlinuxien.com;

    ssl_certificate     /etc/letsencrypt/live/weigo.tuxlinuxien.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/weigo.tuxlinuxien.com/privkey.pem;
    ssl_protocols TLSv1 TLSv1.1 TLSv1.2; # Dropping SSLv3, ref: POODLE
    ssl_prefer_server_ciphers on;
    ssl_ciphers EECDH+CHACHA20:EECDH+AES128:RSA+AES128:EECDH+AES256:RSA+AES256:EECDH+3DES:RSA+3DES:!MD5;
    ssl_session_cache shared:SSL:10m;
    ssl_session_timeout 10m;
    ssl_stapling on;
    ssl_stapling_verify on;
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains; preload" always;

    keepalive_timeout 5;
    charset utf8;

    location /dist/ {
        expires 1d;
        access_log off;
        add_header Cache-Control "public";
        root ;
    }

    location /uploaded/ {
        expires 1M;
        access_log off;
        add_header Cache-Control "public";
        root ;
    }

    location / {
        proxy_set_header Host $host;
        proxy_redirect off;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_pass http://localhost:1314;
    }
}

server {
    listen 80;
    server_name weigo.tuxlinuxien.com;
    return 302 https://$server_name$request_uri;
}
```

## Config

### dev

```
DB=postgres://weigo@127.0.0.1/weigo
UI_DOMAIN=http://localhost:8080/index.dev.html
GOOGLE_CLIENT_KEY=
GOOGLE_CLIENT_SECRET=
GOOGLE_REDIRECT=http://localhost:1314/user/oauth/google/callback
GITHUB_CLIENT_KEY=
GITHUB_CLIENT_SECRET=
GITHUB_REDIRECT=http://localhost:1314/user/oauth/github/callback
```

### prod

```
DB=postgres://weigo@127.0.0.1/weigo
UI_DOMAIN=https://weigo.tuxlinuxien.com/
GOOGLE_CLIENT_KEY=
GOOGLE_CLIENT_SECRET=
GOOGLE_REDIRECT=https://weigo.tuxlinuxien.com/user/oauth/google/callback
GITHUB_CLIENT_KEY=
GITHUB_CLIENT_SECRET=
GITHUB_REDIRECT=https://weigo.tuxlinuxien.com/user/oauth/github/callback
```

## Daemon

Since **systemd** is integrated with any modern linux, it's encouraged to use it
instead of **init.d**.

```
[Unit]
Description=Api weigo
After=network.target

[Service]
EnvironmentFile=
ExecStart=
Restart=always
RestartSec=30

[Install]
WantedBy=multi-user.target
```
