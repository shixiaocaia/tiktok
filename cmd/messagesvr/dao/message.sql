USE camps_tiktok;

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