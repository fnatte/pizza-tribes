-- +goose Up
CREATE TABLE user_item (
    user_id INTEGER NOT NULL,
    item_id TEXT NOT NULL,
    PRIMARY KEY(user_id, item_id),
    FOREIGN KEY (user_id) REFERENCES user(id)
);

-- +goose Down
DROP TABLE user_item;

