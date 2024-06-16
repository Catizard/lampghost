CREATE TABLE difftable_header (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    alias TEXT,
    last_update TEXT,
    symbol TEXT NOT NULL,
    data_location TEXT NOT NULL,
    data_url TEXT NOT NULL,
    original_url TEXT 
);

CREATE TABLE course_info (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    source TEXT NOT NULL,
    md5s TEXT NOT NULL
);

CREATE TABLE rival_info (
    name TEXT(255) NOT NULL,
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    score_log_path TEXT(255),
    song_data_path TEXT(255),
    lr2_user_data_path TEXT(255) 
);

CREATE TABLE rival_tag (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    rival_id INTEGER NOT NULL,
    tag_name TEXT(255) NOT NULL,
    tag_source TEXT(255) NOT NULL,
    'generated' INTEGER DEFAULT (0) NOT NULL,
    'timestamp' TEXT NOT NULL,
)