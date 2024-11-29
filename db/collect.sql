create database zhifou_collect;
use zhifou_collect;

CREATE TABLE `collect_record` (
   `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
   `biz_id` varchar(64) NOT NULL DEFAULT '' COMMENT '业务ID',
   `obj_id` bigint(20) UNSIGNED NOT NULL DEFAULT '0' COMMENT '收藏对象id',
   `user_id` bigint(20) UNSIGNED NOT NULL DEFAULT '0' COMMENT '用户ID',
   `collect_status` tinyint(4) NOT NULL DEFAULT '1' COMMENT '状态 1:收藏 2:取消收藏',
   `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
   `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
   PRIMARY KEY (`id`),
   KEY `ix_update_time` (`update_time`),
   UNIQUE KEY `uk_biz_obj_uid` (`biz_id`,`obj_id`,`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='收藏记录表';

CREATE TABLE `collect_count` (
      `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
      `biz_id` varchar(64) NOT NULL DEFAULT '' COMMENT '业务ID',
      `obj_id` bigint(20) UNSIGNED NOT NULL DEFAULT '0' COMMENT '收藏对象id',
      `collect_num` int(11) NOT NULL DEFAULT '0' COMMENT '收藏数',
      `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
      `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
      PRIMARY KEY (`id`),
      KEY `ix_update_time` (`update_time`),
      UNIQUE KEY `uk_biz_obj` (`biz_id`,`obj_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='收藏计数表';