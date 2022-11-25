-- +goose Up
CREATE TABLE leaderboard (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    game_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    coins INTEGER NOT NULL,
    FOREIGN KEY (user_id) REFERENCES user(id),
    FOREIGN KEY (game_id) REFERENCES game(id)
);

CREATE INDEX idx_leaderboard_coins ON leaderboard(coins);

-- +goose Down
DROP TABLE leaderboard;

