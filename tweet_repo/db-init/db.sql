CREATE DATABASE art;

USE art;

CREATE TABLE IF NOT EXISTS tweets (
    ID BIGINT PRIMARY KEY,
    AUTHOR_ID BIGINT NOT NULL,
    CONTENT VARCHAR(500),
    CREATED_AT DATETIME NOT NULL,
    LANG VARCHAR(10),
    RETWEET_COUNT INT NOT NULL,
    LIKE_COUNT INT NOT NULL,
    HASHTAGS TEXT,
    SENTIMENT FLOAT
)
