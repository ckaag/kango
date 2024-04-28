CREATE TABLE IF NOT EXISTS authors (
    id INTEGER PRIMARY KEY,
    name text NOT NULL,
    bio text
);
CREATE TABLE IF NOT EXISTS counters (
    "name" text PRIMARY key,
    "counter" integer NOT NULL
);