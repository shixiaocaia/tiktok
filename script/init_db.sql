SET NAMES utf8mb4;
SET
FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for comments
-- ----------------------------
USE mini_tiktok;

DROP TABLE IF EXISTS `t_comment`;
CREATE TABLE `t_comment`
(
    `id`           bigint(20) NOT NULL AUTO_INCREMENT COMMENT '评论id，自增主键',
    `user_id`      bigint(20) NOT NULL COMMENT '评论发布用户id',
    `video_id`     bigint(20) NOT NULL COMMENT '评论视频id',
    `comment_text` varchar(255) NOT NULL COMMENT '评论内容',
    `create_time`  datetime     NOT NULL COMMENT '评论发布时间',
    PRIMARY KEY (`id`),
    KEY            `videoIdIdx` (`video_id`) USING BTREE COMMENT '评论列表使用视频id作为索引-方便查看视频下的评论列表'
) ENGINE=InnoDB AUTO_INCREMENT=1206 DEFAULT CHARSET=utf8 COMMENT='评论表';

-- ----------------------------
-- Table structure for follows
-- ----------------------------
DROP TABLE IF EXISTS `t_relation`;
CREATE TABLE `t_relation`
(
    `id`          bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增主键',
    `follow_id`   bigint(20) NOT NULL COMMENT '用户id',
    `follower_id` bigint(20) NOT NULL COMMENT '关注的用户',
    PRIMARY KEY (`id`),
    UNIQUE KEY `followIdtoFollowerIdIdx` (`follow_id`,`follower_id`) USING BTREE,
    KEY           `FollowIdIdx` (`follow_id`) USING BTREE,
    KEY           `FollowerIdIdx` (`follower_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1096 DEFAULT CHARSET=utf8 COMMENT='关注表';

-- ----------------------------
-- Table structure for likes
-- ----------------------------
DROP TABLE IF EXISTS `t_favorite`;
CREATE TABLE `t_favorite`
(
    `id`       bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增主键',
    `user_id`  bigint(20) NOT NULL COMMENT '点赞用户id',
    `video_id` bigint(20) NOT NULL COMMENT '被点赞的视频id',
    PRIMARY KEY (`id`),
    UNIQUE KEY `userIdtoVideoIdIdx` (`user_id`,`video_id`) USING BTREE,
    KEY        `userIdIdx` (`user_id`) USING BTREE,
    KEY        `videoIdx` (`video_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1229 DEFAULT CHARSET=utf8 COMMENT='点赞表';

-- ----------------------------
-- Table structure for users
-- ----------------------------
DROP TABLE IF EXISTS `t_user`;
CREATE TABLE `t_user`
(
    `id`               bigint(20) NOT NULL AUTO_INCREMENT COMMENT '用户id，自增主键',
    `user_name`             varchar(255) NOT NULL COMMENT '用户名',
    `password`         varchar(255) NOT NULL COMMENT '用户密码',
    `follow_count`     bigint(20) NOT NULL DEFAULT 0 COMMENT '该用户关注其他用户个数',
    `follower_count`   bigint(20) NOT NULL DEFAULT 0 COMMENT '该用户粉丝个数',
    `total_favorited`  bigint(20) NOT NULL DEFAULT 0 COMMENT '该用户被喜欢的视频数量',
    `favorite_count`   bigint(20) NOT NULL DEFAULT 0 COMMENT '该用户喜欢的视频数量',
    `signature` varchar(1024) COMMENT '签名',
    `avatar`           varchar(1024) COMMENT '用户头像',
    `background_image` varchar(1024) COMMENT '主页背景',
    `work_count`      bigint(20) NOT NULL DEFAULT 0 COMMENT '作品数量',
    PRIMARY KEY (`id`),
    KEY                `name_password_idx` (`user_name`,`password`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=20044 DEFAULT CHARSET=utf8 COMMENT='用户表';

-- ----------------------------
-- Table structure for videos
-- ----------------------------
DROP TABLE IF EXISTS `t_video`;
CREATE TABLE `t_video`
(
    `id`           bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增主键，视频唯一id',
    `author_id`    bigint(20) NOT NULL COMMENT '视频作者id',
    `play_url`     varchar(255) NOT NULL COMMENT '播放url',
    `cover_url`    varchar(255) NOT NULL COMMENT '封面url',
    `favorite_count` bigint(20) NOT NULL DEFAULT 0 COMMENT '视频的点赞数量',
    `comment_count` bigint(20) NOT NULL DEFAULT 0 COMMENT '视频的评论数量',
    `publish_time` bigint(20)     NOT NULL COMMENT '发布时间戳',
    `title`        varchar(255) DEFAULT NULL COMMENT '视频名称',
    PRIMARY KEY (`id`),
    KEY            `time` (`publish_time`) USING BTREE,
    KEY            `author` (`author_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=115 DEFAULT CHARSET=utf8 COMMENT='视频表';

-- ----------------------------
-- Table structure for messages
-- ----------------------------
DROP TABLE IF EXISTS `t_message`;
CREATE TABLE `t_message`
(
    `id`            int(11) NOT NULL AUTO_INCREMENT COMMENT '主键',
    `to_user_id`    int(11) NOT NULL COMMENT '接收用户id',
    `from_user_id`  int(11) NOT NULL COMMENT '发送用户id',
    `content`       varchar(255) NOT NULL COMMENT '消息内容',
    `create_time`   datetime NOT NULL COMMENT '创建时间',
    PRIMARY KEY (`id`),
    KEY            `to_user_id` (`to_user_id`),
    KEY           `from_user_id` (`from_user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1406 DEFAULT CHARSET=utf8 COMMENT='消息表';
