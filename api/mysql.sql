CREATE DATABASE IF NOT EXISTS resume DEFAULT CHARACTER SET utf8mb4 DEFAULT COLLATE utf8mb4_0900_ai_ci;

USE resume;

CREATE TABLE IF NOT EXISTS `jobs` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(128) NOT NULL DEFAULT '' COMMENT '职位名称',
  `demand` text COMMENT '职位要求(详细描述对应聘者的技能、经验等要求)',
  `desc` text COMMENT '职位描述(详细描述职位的工作内容、职责等)',
  `create_uid` int unsigned NOT NULL DEFAULT '0' COMMENT '创建人',
  `create_time` int unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_time` int unsigned NOT NULL DEFAULT '0' COMMENT '修改时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='招聘职位表';

CREATE TABLE IF NOT EXISTS `resumes` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `job_id` int unsigned NOT NULL DEFAULT '0' COMMENT '关联的招聘职位ID',
  `filename` varchar(128) NOT NULL DEFAULT '' COMMENT '上传简历的文件名',
  `size` int unsigned NOT NULL DEFAULT '0' COMMENT '简历文件大小',
  `create_uid` int unsigned NOT NULL DEFAULT '0' COMMENT '上传人',
  `create_time` int unsigned NOT NULL DEFAULT '0' COMMENT '创建时间即上传时间',
  `update_time` int unsigned NOT NULL DEFAULT '0' COMMENT '修改时间',
  PRIMARY KEY (`id`),
  KEY `job_id` (`job_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='上传的简历表';

CREATE TABLE IF NOT EXISTS `resume_analyze_records` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `resume_id` int unsigned NOT NULL DEFAULT '0' COMMENT '关联的简历ID',
  `status` varchar(32) NOT NULL DEFAULT '' COMMENT '分析状态 pending分析中 completed已完成 failed失败',
  `result` text COMMENT 'AI分析结果',
  `create_uid` int unsigned NOT NULL DEFAULT '0' COMMENT '创建分析的用户',
  `create_time` int unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `resume_id` (`resume_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='简历的AI分析历史记录表';

CREATE TABLE IF NOT EXISTS `users` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(32) NOT NULL DEFAULT '' COMMENT '用户名',
  `password` varchar(32) NOT NULL DEFAULT '' COMMENT '密码',
  `email` varchar(255)  NOT NULL DEFAULT '' COMMENT '邮箱地址',
  `create_uid` int unsigned NOT NULL DEFAULT '0' COMMENT '创建人',
  `create_time` int unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_time` int unsigned NOT NULL DEFAULT '0' COMMENT '修改时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

CREATE TABLE IF NOT EXISTS `user_roles` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `uid` int unsigned NOT NULL DEFAULT '0' COMMENT '用户ID',
  `role_id` int unsigned NOT NULL DEFAULT '0' COMMENT '关联的角色ID',
  `create_time` int unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `uid` (`uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户关联的角色';

CREATE TABLE IF NOT EXISTS `roles` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(32) NOT NULL DEFAULT '' COMMENT '角色名',
  `create_uid` int unsigned NOT NULL DEFAULT '0' COMMENT '创建人',
  `create_time` int unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_time` int unsigned NOT NULL DEFAULT '0' COMMENT '修改时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色表';

CREATE TABLE IF NOT EXISTS `role_perms` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `role_id` int unsigned NOT NULL DEFAULT '0' COMMENT '关联的角色ID',
  `value` varchar(128) NOT NULL DEFAULT '' COMMENT '权限名称, 如rest.job.get',
  `create_uid` int unsigned NOT NULL DEFAULT '0' COMMENT '创建人',
  `create_time` int unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `role_id` (`role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色权限表';
