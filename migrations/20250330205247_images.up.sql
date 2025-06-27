
CREATE TYPE STATUS_ENUM as
    ENUM(
        'week_top',
        'month_top',
        'year_top',
        ''
        );
CREATE TABLE IF NOT EXISTS images(
    id TEXT PRIMARY KEY ,
    text TEXT NOT NULL ,
    content bytea NOT NULL ,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    status STATUS_ENUM DEFAULT '',
    prompt TEXT NOT NULL
);
