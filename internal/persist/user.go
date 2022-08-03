package persist

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/fnatte/pizza-tribes/internal/redis"
	"github.com/rs/xid"
	"golang.org/x/crypto/bcrypt"
)

type UserDbo struct {
	Id             string `redis:"id"`
	Username       string `redis:"username"`
	HashedPassword string `redis:"hashed_password"`
}

type userRepo struct {
	rdb redis.RedisClient
}

var ErrUsernameTaken = errors.New("username is taken")
var ErrInvalidUsername = errors.New("username is invalid")

func NewUserRepository(rdb redis.RedisClient) *userRepo {
	return &userRepo{
		rdb: rdb,
	}
}

func (r *userRepo) GetUserLatestActivity(ctx context.Context, userId string) (int64, error) {
	key := fmt.Sprintf("user:%s:latest_activity", userId)

	val, err := r.rdb.Get(ctx, key).Int64()
	if err != nil && err != redis.Nil {
		return 0, fmt.Errorf("failed to get latest user activity: %w", err)
	}

	return val, nil
}

func (r *userRepo) SetUserLatestActivity(ctx context.Context, userId string, val int64) error {
	key := fmt.Sprintf("user:%s:latest_activity", userId)

	err := r.rdb.Set(ctx, key, val, 0).Err()
	if err != nil {
		return fmt.Errorf("failed to set latest user activity: %w", err)
	}

	return nil
}

func (r *userRepo) GetAllUsers(ctx context.Context) ([]string, error) {
	return r.rdb.SMembers(ctx, "users").Result()
}

func (r *userRepo) DeleteUser(ctx context.Context, userId string) error {
	key := fmt.Sprintf("user:%s", userId)

	u, err := r.GetUser(ctx, userId)
	if err != nil {
		return err
	}

	pipe := r.rdb.Pipeline()
	pipe.Del(ctx, fmt.Sprintf("%s:latest_activity", key))
	pipe.Del(ctx, fmt.Sprintf("%s:reports", key))
	pipe.Del(ctx, fmt.Sprintf("%s:reportsByDate", key))
	pipe.Del(ctx, fmt.Sprintf("%s:gamestate", key))
	pipe.Del(ctx, fmt.Sprintf("%s:fcm_tokens", key))
	pipe.Del(ctx, fmt.Sprintf("%s:ts_pizzas", key))
	pipe.Del(ctx, fmt.Sprintf("%s:ts_coins", key))
	pipe.Del(ctx, key)
	pipe.Del(ctx, fmt.Sprintf("username:%s", u.Username))
	pipe.SRem(ctx, "users", userId)
	pipe.ZRem(ctx, "user_updates", userId)
	pipe.ZRem(ctx, "leaderboard", userId)

	_, err = pipe.Exec(ctx)
	return err
}

func (r *userRepo) GetUser(ctx context.Context, userId string) (*UserDbo, error) {
	userKey := fmt.Sprintf("user:%s", userId)
	user := UserDbo{}
	err := r.rdb.HGetAll(ctx, userKey).Scan(&user)
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepo) FindUser(ctx context.Context, username string) (string, error) {
	usernameKey := fmt.Sprintf("username:%s", strings.ToLower(username))
	userId, err := r.rdb.Get(ctx, usernameKey).Result()
	if err != nil && err != redis.Nil {
		return "", err
	}

	return userId, nil
}

var validUsername = regexp.MustCompile(`^[a-zA-Z]+[a-zA-Z0-9_]*$`)

func IsValidUsername(username string) bool {
	return validUsername.MatchString(username) && len(username) >= 3 && len(username) <= 20
}

func (r *userRepo) CreateUser(ctx context.Context, username string, password string) (string, error) {
	if !IsValidUsername(username) {
		return "", ErrInvalidUsername
	}

	id := xid.New().String()
	usernameKey := fmt.Sprintf("username:%s", strings.ToLower(username))
	userKey := fmt.Sprintf("user:%s", id)

	// Check for existing user with this username
	res, err := r.rdb.Exists(ctx, usernameKey).Result()
	if err != nil && err != redis.Nil {
		return "", err
	}
	if res != 0 {
		return "", ErrUsernameTaken
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	txf := func(tx *redis.Tx) error {
		_, err := tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			pipe.Set(ctx, usernameKey, id, 0)
			pipe.HSet(ctx, userKey, "id", id, "username", username, "hashed_password", hash)
			pipe.SAdd(ctx, "users", id)
			return nil
		})
		return err
	}

	err = r.rdb.Watch(ctx, txf, usernameKey, userKey)

	return id, err
}
