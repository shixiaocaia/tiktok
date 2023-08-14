package dao

import (
	"context"
	"encoding/json"
	"github.com/go-redsync/redsync/v4"
	"github.com/shixiaocaia/tiktok/cmd/videosvr/config"
	"github.com/shixiaocaia/tiktok/cmd/videosvr/constant"
	"github.com/shixiaocaia/tiktok/cmd/videosvr/log"
	"sort"
	"strconv"
	"time"
)

// GetRedisLock 获取redis锁
func acquireLock(key string) (*redsync.Mutex, error) {
	redisSync := GetRedSync()
	mutex := redisSync.NewMutex(key)
	if err := mutex.Lock(); err != nil {
		log.Errorf("redis lock err: %v", err)
		return nil, err
	}
	log.Info("redis lock success")
	return mutex, nil
}

// DelRedisLock 释放redis锁
func releaseLock(mutex *redsync.Mutex) error {
	if ok, err := mutex.Unlock(); err != nil || !ok {
		log.Errorf("redis unlock err: %v", err)
		return err
	}
	log.Info("redis unlock success")
	return nil
}

// SetVideoCache 添加video信息到redis hset中
func SetVideoCache(videoList []Video) error {
	redisKey := constant.VideosInfoPrefix
	redisCli := GetRedisCli()

	for _, video := range videoList {
		videoBytes, err := json.Marshal(video)
		if err != nil {
			log.Errorf("json marshal video err:%v", err)
			return err
		}
		videoIDStr := strconv.FormatInt(video.Id, 10)
		err = redisCli.HSet(context.Background(), redisKey, videoIDStr, string(videoBytes)).Err()
		if err != nil {
			log.Errorf("redis hset video err:%v", err)
			return err
		}
	}

	// 设置过期时间
	expired := time.Second * time.Duration(config.GetGlobalConfig().RedisConfig.Expired)
	err := redisCli.Expire(context.Background(), redisKey, expired).Err()
	if err != nil {
		log.Errorf("redis expire err:%v", err)
		return err
	}

	return nil
}

// GetVideoCache 获取单个video信息
func GetVideoCache(videoID int64) (*Video, error) {
	redisKey := constant.VideosInfoPrefix
	redisCli := GetRedisCli()

	videoInfo, err := redisCli.HGet(context.Background(), redisKey, strconv.FormatInt(videoID, 10)).Result()
	if err != nil {
		log.Errorf("redis hget video err:%v", err)
		return nil, err
	}

	var video Video
	err = json.Unmarshal([]byte(videoInfo), &video)
	if err != nil {
		log.Errorf("json unmarshal video err:%v", err)
		return nil, err
	}

	return &video, nil
}

// GetVideoListCache 获取一组video信息
func GetVideoListCache(lastTime int64) ([]Video, error) {
	redisKey := constant.VideosInfoPrefix
	redisCli := GetRedisCli()

	// 获取所有的video信息
	videoMap, err := redisCli.HGetAll(context.Background(), redisKey).Result()
	if err != nil {
		log.Errorf("redis hgetall video err:%v", err)
		return nil, err
	}

	// 将map转换为slice
	var videoList []Video
	for _, videoStr := range videoMap {
		var video Video
		err := json.Unmarshal([]byte(videoStr), &video)
		if err != nil {
			log.Errorf("json unmarshal video err:%v", err)
			return nil, err
		}
		videoList = append(videoList, video)
	}

	// 根据lastTime过滤
	var filterVideoList []Video
	for _, video := range videoList {
		if video.PublishTime < lastTime {
			filterVideoList = append(filterVideoList, video)
		}
		if len(filterVideoList) >= 30 {
			break
		}
	}
	// 根据publishTime排序，从大到小
	sort.Slice(filterVideoList, func(i, j int) bool {
		return filterVideoList[i].PublishTime > filterVideoList[j].PublishTime
	})

	return filterVideoList, nil
}

// UpdateVideoCache 清空缓存
func UpdateVideoCache() error {
	redisKey := constant.VideosInfoPrefix
	redisCli := GetRedisCli()

	if err := redisCli.Del(context.Background(), redisKey).Err(); err != nil {
		log.Errorf("del videoInfo hset err:%v", err)
		return err
	}
	return nil
}

// UpdateFavoriteCountCache 更新点赞信息缓存
func UpdateFavoriteCountCache(videoID int64, actionType int64) error {
	// 1. 获取分布式锁
	videoLock := constant.VideosInfoPrefix + "_" + strconv.FormatInt(videoID, 10)
	mutex, err := acquireLock(videoLock)
	if err != nil {
		log.Errorf("acquire lock err: %v", err)
		releaseLock(mutex)
		return err
	}
	defer releaseLock(mutex)

	// 2. 更新信息
	redisKey := constant.VideosInfoPrefix
	redisCli := GetRedisCli()

	// 更新videoID缓存
	realVideoInfo, err := GetVideoCache(videoID)
	if err != nil {
		log.Errorf("get videocache err: %v", err)
		return err
	}

	count := 1
	if actionType == 2 {
		count = -1
	}

	realVideoInfo.FavoriteCount += int64(count)
	videoBytes, err := json.Marshal(realVideoInfo)
	if err != nil {
		log.Errorf("json marshal video err: %v", err)
		return err
	}
	err = redisCli.HSet(redisCli.Context(), redisKey, strconv.FormatInt(videoID, 10), string(videoBytes)).Err()
	if err != nil {
		log.Errorf("redis hset video err: %v", err)
		return err
	}

	return nil
}

// UpdateCommentCountCache 更新评论信息缓存
func UpdateCommentCountCache(videoID, actionType int64) error {
	// 1. 获取分布式锁
	videoLock := constant.VideosInfoPrefix + "_" + strconv.FormatInt(videoID, 10)
	mutex, err := acquireLock(videoLock)
	if err != nil {
		log.Errorf("acquire lock err: %v", err)
		releaseLock(mutex)
		return err
	}
	defer releaseLock(mutex)

	// 2. 更新信息
	redisKey := constant.VideosInfoPrefix
	redisCli := GetRedisCli()

	// 更新videoID缓存
	realVideoInfo, err := GetVideoCache(videoID)
	if err != nil {
		log.Errorf("get videocache err: %v", err)
		return err
	}

	count := 1
	if actionType == 2 {
		count = -1
	}
	realVideoInfo.CommentCount += int64(count)
	videoBytes, err := json.Marshal(realVideoInfo)
	if err != nil {
		log.Errorf("json marshal video err: %v", err)
		return err
	}
	err = redisCli.HSet(redisCli.Context(), redisKey, strconv.FormatInt(videoID, 10), string(videoBytes)).Err()
	if err != nil {
		log.Errorf("redis hset video err: %v", err)
		return err
	}

	return nil
}
