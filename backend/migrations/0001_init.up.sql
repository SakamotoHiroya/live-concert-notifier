CREATE TABLE users (
    id         uuid PRIMARY KEY,
    email      text NOT NULL UNIQUE,
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE artists (
    id                 uuid PRIMARY KEY,
    name               text NOT NULL,
    official_site_url  text NOT NULL UNIQUE,
    created_at         timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE user_artists (
    user_id     uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    artist_id   uuid NOT NULL REFERENCES artists (id) ON DELETE CASCADE,
    followed_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, artist_id)
);

CREATE TABLE concerts (
    id              uuid PRIMARY KEY,
    artist_id       uuid NOT NULL REFERENCES artists (id) ON DELETE CASCADE,
    title           text NOT NULL DEFAULT '',
    venue_name      text NOT NULL,
    venue_location  text NOT NULL,
    date            date NOT NULL,
    co_performers   text[] NOT NULL DEFAULT '{}',
    is_festival     boolean NOT NULL DEFAULT false,
    source_url      text NOT NULL,
    raw_text        text NOT NULL DEFAULT '',
    discovered_at   timestamptz NOT NULL DEFAULT now(),
    created_at      timestamptz NOT NULL DEFAULT now(),
    UNIQUE (artist_id, date, venue_name)
);

CREATE INDEX concerts_date_idx ON concerts (date);

CREATE TABLE scrape_jobs (
    id            uuid PRIMARY KEY,
    artist_id     uuid NOT NULL REFERENCES artists (id) ON DELETE CASCADE,
    status        text NOT NULL CHECK (status IN ('pending', 'running', 'succeeded', 'failed')),
    started_at    timestamptz,
    finished_at   timestamptz,
    error_message text
);

CREATE INDEX scrape_jobs_artist_id_idx ON scrape_jobs (artist_id);
CREATE INDEX scrape_jobs_status_idx ON scrape_jobs (status);
