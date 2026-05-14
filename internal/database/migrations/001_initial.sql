CREATE TABLE IF NOT EXISTS jokes (
    id INTEGER PRIMARY KEY,
    joke text NOT NULL,
    nsfw bool NOT NULL,
    created_at date NOT NULL
);

