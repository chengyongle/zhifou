create database zhifou_comment;
use zhifou_comment;

CREATE TABLE `comment` (
     `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
     `biz_id` varchar(64) NOT NULL DEFAULT '' COMMENT '业务ID',
     `obj_id` bigint(20) UNSIGNED NOT NULL DEFAULT '0' COMMENT '评论对象id',
     `comment_user_id` bigint(20) UNSIGNED NOT NULL DEFAULT '0' COMMENT '发表评论的用户ID',
     `be_comment_user_id` bigint(20) UNSIGNED NOT NULL DEFAULT '0' COMMENT '被评论用户ID',
     `parent_id` bigint(20) UNSIGNED NOT NULL DEFAULT '0' COMMENT '父评论ID，为0表示这是根评论',
     `content` text COLLATE utf8_unicode_ci NOT NULL COMMENT '内容',
     `status` tinyint(4) NOT NULL DEFAULT '0' COMMENT '状态 0:正常 1:删除',
     `like_num` int(11) NOT NULL DEFAULT '0' COMMENT '点赞数',
     `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
     `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
     PRIMARY KEY (`id`),
     KEY `ix_update_time` (`update_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='评论表';


CREATE TABLE `comment_count` (
    `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `biz_id` varchar(64) NOT NULL DEFAULT '' COMMENT '业务ID',
    `obj_id` bigint(20) UNSIGNED NOT NULL DEFAULT '0' COMMENT '评论对象id',
    `comment_num` int(11) NOT NULL DEFAULT '0' COMMENT '评论总数',
    `comment_root_num` int(11) NOT NULL DEFAULT '0' COMMENT '根评论总数',
    `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
    PRIMARY KEY (`id`),
    KEY `ix_update_time` (`update_time`),
    UNIQUE KEY `uk_biz_obj` (`biz_id`,`obj_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='评论计数';

