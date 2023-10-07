package dao

import (
	"context"
	"encoding/json"
	"github.com/go-redsync/redsync/v4"
	"github.com/shixiaocaia/tiktok/cmd/commentsvr/config"
	"github.com/shixiaocaia/tiktok/cmd/commentsvr/constant"
	"github.com/shixiaocaia/tiktok/cmd/commentsvr/log"
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

// SetCommentCacheInfo 添加videoId对应的评论信息到redis hset中
func SetCommentCacheInfo(vid int64, comment []*Comment) error {
	redisKey := constant.CommentInfoPrefix + strconv.FormatInt(vid, 10)
	redisCli := GetRedisCli()

	for _, comm := range comment {
		commentBytes, err := json.Marshal(comm)
		if err != nil {
			log.Errorf("json marshal comment err:%v", err)
			return err
		}
		commentIDStr := strconv.FormatInt(comm.Id, 10)
		err = redisCli.HSet(context.Background(), redisKey, commentIDStr, string(commentBytes)).Err()
		if err != nil {
			log.Errorf("redis hset comment err:%v", err)
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

// GetCommentCacheList 获取某个video的评论列表
func GetCommentCacheList(vid int64) ([]*Comment, error) {
	redisKey := constant.CommentInfoPrefix + strconv.FormatInt(vid, 10)
	redisCli := GetRedisCli()
	commentList := make([]*Comment, 0)
	commentMap, err := redisCli.HGetAll(context.Background(), redisKey).Result()
	if err != nil {
		log.Errorf("redis hgetall err:%v", err)
		return nil, err
	}
	for _, v := range commentMap {
		comment := &Comment{}
		err := json.Unmarshal([]byte(v), comment)
		if err != nil {
			log.Errorf("json unmarshal err:%v", err)
			return nil, err
		}
		commentList = append(commentList, comment)
	}

	// HSet存放是无序的，按照评论的ID进行排序
	sort.Slice(commentList, func(i, j int) bool {
		return commentList[i].Id > commentList[j].Id
	})

	return commentList, nil
}

// DelCommentCacheInfo 删除某个video的评论列表
func DelCommentCacheInfo(vid int64) error {
	redisKey := constant.CommentInfoPrefix + strconv.FormatInt(vid, 10)
	redisCli := GetRedisCli()
	if err := redisCli.Del(context.Background(), redisKey).Err(); err != nil {
		log.Errorf("del redis hset err:%v", err)
		return err
	}
	return nil
}

//func CommentCacheAdd(comment *Comment) error {
//	mutex, err := acquireLock(constant.CommentLock)
//	if err != nil {
//		log.Errorf("acquire lock err: %v", err)
//		releaseLock(mutex)
//		return err
//	}
//	defer releaseLock(mutex)
//
//	redisKey := constant.CommentInfoPrefix + strconv.FormatInt(comment.VideoId, 10)
//	redisCli := GetRedisCli()
//	commentBytes, err := json.Marshal(comment)
//	if err != nil {
//		log.Errorf("json marshal comment err: %v", err)
//		return err
//	}
//	commentIDStr := strconv.FormatInt(comment.Id, 10)
//	err = redisCli.HSet(redisCli.Context(), redisKey, commentIDStr, string(commentBytes)).Err()
//	if err != nil {
//		log.Errorf("redis hset comment err: %v", err)
//		return err
//	}
//	log.Info("redis hset comment success")
//	return nil
//
//}
//
//func CommentCacheDel(commentId, vid int64) error {
//	mutex, err := acquireLock(constant.CommentLock)
//	if err != nil {
//		log.Errorf("acquire lock err: %v", err)
//		releaseLock(mutex)
//		return err
//	}
//	defer releaseLock(mutex)
//
//	redisKey := constant.CommentInfoPrefix + strconv.FormatInt(vid, 10)
//	redisCli := GetRedisCli()
//	commentIDStr := strconv.FormatInt(commentId, 10)
//	err = redisCli.HDel(context.Background(), redisKey, commentIDStr).Err()
//	if err != nil {
//		log.Errorf("redis hdel comment err:%v", err)
//		return err
//	}
//	return nil
//}
