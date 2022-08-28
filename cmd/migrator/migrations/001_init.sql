-- +goose Up
CREATE TABLE user (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL,
    email TEXT,
    hashed_password TEXT NOT NULL
);

CREATE TABLE game (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    status TEXT NOT NULL,
    host TEXT NOT NULL
);

CREATE TABLE user_game (
    user_id INTEGER NOT NULL,
    game_id INTEGER NOT NULL,
    PRIMARY KEY(user_id, game_id),
    FOREIGN KEY (user_id) REFERENCES user(id),
    FOREIGN KEY (game_id) REFERENCES game(id)
);

-- +goose Down
DROP TABLE user;
DROP TABLE game;

