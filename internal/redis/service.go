package redis

import (
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

const EmptyLogin = "not saved login"

func SetPath(logger *zap.Logger, rdb *redis.Client, userID int64, val string) {
	id := strconv.FormatInt(userID, 10)
	res := rdb.Set(id, val, 7*24*time.Hour)
	if res.Err() != nil {
		logger.Error("set path", zap.Error(res.Err()))
	}
}

func SetTaskUserID(logger *zap.Logger, rdb *redis.Client, userID int64, toUserID int64) {
	id := strconv.FormatInt(userID, 10)
	tId := strconv.FormatInt(toUserID, 10)
	res := rdb.Set("task_"+id, tId, 7*24*time.Hour)
	if res.Err() != nil {
		logger.Error("set path", zap.Error(res.Err()))
	}
}

func GetTaskUserID(logger *zap.Logger, rdb *redis.Client, userID int64) string {
	id := strconv.FormatInt(userID, 10)
	have, err := rdb.Exists("task_" + id).Result()
	if err != nil {
		logger.Error("path exists", zap.Error(err))
	}
	if have == 0 {
		return EmptyLogin
	}

	value, err := rdb.Get("task_" + id).Result()
	if err != nil {
		logger.Error("get path", zap.Error(err))
	}

	return value
}

func GetPath(logger *zap.Logger, rdb *redis.Client, userID int64) string {
	id := strconv.FormatInt(userID, 10)
	have, err := rdb.Exists(id).Result()
	if err != nil {
		logger.Error("path exists", zap.Error(err))
	}
	if have == 0 {
		return EmptyLogin
	}

	value, err := rdb.Get(id).Result()
	if err != nil {
		logger.Error("get path", zap.Error(err))
	}

	return value
}

func SetLogin(logger *zap.Logger, rdb *redis.Client, userID int64, login string) {
	id := strconv.FormatInt(userID, 10)
	res := rdb.Set("login"+id, login, 7*24*time.Hour)
	if res.Err() != nil {
		logger.Error("set login", zap.Error(res.Err()))
	}
}

func GetLogin(logger *zap.Logger, rdb *redis.Client, userID int64) string {
	id := strconv.FormatInt(userID, 10)
	have, err := rdb.Exists("login" + id).Result()
	if err != nil {
		logger.Error("login exists", zap.Error(err))
	}
	if have == 0 {
		return EmptyLogin
	}

	value, err := rdb.Get("login" + id).Result()
	if err != nil {
		logger.Error("get login", zap.Error(err))
	}
	return value
}
