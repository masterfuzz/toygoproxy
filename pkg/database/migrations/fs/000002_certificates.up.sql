CREATE TABLE IF NOT EXISTS certificates (
    hostname TEXT UNIQUE PRIMARY KEY,
    certificate TEXT NOT NULL,
    private_key TEXT NOT NULL
);
