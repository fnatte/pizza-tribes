package sql

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"strconv"

	"github.com/fnatte/pizza-tribes/internal/game/persist"
	"golang.org/x/crypto/bcrypt"
)

type userDb struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *userDb {
	return &userDb{ db: db }
}

var validUsername = regexp.MustCompile(`^[a-zA-Z]+[a-zA-Z0-9_]*$`)

func IsValidUsername(username string) bool {
	return validUsername.MatchString(username) && len(username) >= 3 && len(username) <= 30
}

func (db *userDb) CreateUser(ctx context.Context, username string, password string) (string, error) {
	if !IsValidUsername(username) {
		return "", persist.ErrInvalidUsername
	}

	// Check for existing user with this username
	res, err := db.db.QueryContext(ctx, `SELECT id FROM user WHERE username = ?`, username)
	if err != nil {
		return "", err
	}
	if res.Next() {
		return "", persist.ErrUsernameTaken
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	res2, err := db.db.ExecContext(ctx,
		`INSERT INTO user (username, hashed_password) VALUES (?, ?)`,
		username, hash)
	if err != nil {
		return "", err
	}
	id, err := res2.LastInsertId()
	if err != nil {
		return "", err
	}


	return strconv.FormatInt(id, 10), err
}

func (db *userDb) GetAllUsers(ctx context.Context) ([]*persist.User, error) {
	rows, err := db.db.QueryContext(ctx, `SELECT id, username, hashed_password FROM user`)
	if err != nil {
		return nil, err
	}

	users := []*persist.User{}

	for rows.Next() {
		u := persist.User{}
		err := rows.Scan(&u.Id, &u.Username, &u.HashedPassword)
		if err != nil {
			return nil, err
		}
		users = append(users, &u)
	}

	return users, nil
}

func (db *userDb) GetUser(ctx context.Context, userId string) (*persist.User, error) {
	row := db.db.QueryRowContext(ctx, `SELECT id, username, hashed_password FROM user WHERE id = ?`, userId)

	u := persist.User{}
	err := row.Scan(&u.Id, &u.Username, &u.HashedPassword)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (db *userDb) GetUserByUsername(ctx context.Context, username string) (*persist.User, error) {
	row := db.db.QueryRowContext(ctx, `SELECT id, username, hashed_password FROM user WHERE username = ?`, username)

	u := persist.User{}
	err := row.Scan(&u.Id, &u.Username, &u.HashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, persist.ErrUserNotFound
		}
		return nil, err
	}

	return &u, nil
}

func (db *userDb) GetUserCount(ctx context.Context) (int64, error) {
	row := db.db.QueryRowContext(ctx, `SELECT COUNT(1) FROM user`)

	var count int64
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (db *userDb) DeleteUser(ctx context.Context, userId string) error {
	res, err := db.db.ExecContext(ctx, `DELETE FROM user WHERE id = ?`, userId)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return persist.ErrUserNotFound
	}

	return nil
}

