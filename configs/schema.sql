DROP DATABASE IF EXISTS `simpleblog`;

CREATE DATABASE IF NOT EXISTS `simpleblog` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE `simpleblog`;

DROP TABLE IF EXISTS `casbin_rule`;

CREATE TABLE `casbin_rule` (
	`id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
	`ptype` varchar(100) DEFAULT NULL,
	`v0` varchar(100) DEFAULT NULL,
	`v1` varchar(100) DEFAULT NULL,
	`v2` varchar(100) DEFAULT NULL,
	`v3` varchar(100) DEFAULT NULL,
	`v4` varchar(100) DEFAULT NULL,
	`v5` varchar(100) DEFAULT NULL,
	PRIMARY KEY (`id`),
	UNIQUE KEY `idx_casbin_rule` (`ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`)
) ENGINE = InnoDB AUTO_INCREMENT = 1 DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

INSERT INTO
	`casbin_rule`
VALUES
	(18, 'g', 'user-000000', 'role::admin', NULL, NULL, '', ''),
	(21, 'p', 'role::admin', '*', '*', 'allow', '', ''),
	(7, 'p', 'role::user', '/v1.SimpleBlog/DeleteUser', 'CALL', 'deny', '', ''),
	(8, 'p', 'role::user', '/v1.SimpleBlog/ListUser', 'CALL', 'deny', '', ''),
	(9, 'p', 'role::user', '/v1/users', 'GET', 'deny', '', ''),
	(10, 'p', 'role::user', '/v1/users/*', 'DELETE', 'deny', '', '');

DROP TABLE IF EXISTS `post`;

CREATE TABLE `post` (
	`id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
	`userID` varchar(36) NOT NULL DEFAULT '' COMMENT 'User unique ID',
	`postID` varchar(35) NOT NULL DEFAULT '' COMMENT 'Post unique ID',
	`title` varchar(256) NOT NULL DEFAULT '' COMMENT 'Post title',
	`content` longtext NOT NULL COMMENT 'Post content',
	`createdAt` datetime NOT NULL DEFAULT current_timestamp() COMMENT 'Post created time',
	`updatedAt` datetime NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp() COMMENT 'Post last updated time',
	PRIMARY KEY (`id`),
	UNIQUE KEY `post.postID` (`postID`),
	KEY `idx.post.userID` (`userID`)
) ENGINE = InnoDB AUTO_INCREMENT = 1 DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = 'Posts table';

DROP TABLE IF EXISTS `user`;

CREATE TABLE `user` (
	`id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
	`userID` varchar(36) NOT NULL DEFAULT '' COMMENT 'User unique ID',
	`username` varchar(255) NOT NULL DEFAULT '' COMMENT 'Username (unique)',
	`password` varchar(255) NOT NULL DEFAULT '' COMMENT 'User password (hashed)',
	`nickname` varchar(30) NOT NULL DEFAULT '' COMMENT 'User nickname',
	`email` varchar(256) NOT NULL DEFAULT '' COMMENT 'User email address',
	`phone` varchar(16) NOT NULL DEFAULT '' COMMENT 'User phone number',
	`createdAt` datetime NOT NULL DEFAULT current_timestamp() COMMENT 'User created time',
	`updatedAt` datetime NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp() COMMENT 'User last updated time',
	PRIMARY KEY (`id`),
	UNIQUE KEY `user.userID` (`userID`),
	UNIQUE KEY `user.username` (`username`),
	UNIQUE KEY `user.phone` (`phone`)
) ENGINE = InnoDB AUTO_INCREMENT = 1 DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = 'Users table';

INSERT INTO
	`user`
VALUES
	(96, 'user-000000', 'root', '$2a$10$ctsFXEUAMd7rXXpmccNlO.ZRiYGYz0eOfj8EicPGWqiz64YBBgR1y', 'colin404', 'colin404@foxmail.com', '18110000000', '2024-12-12 03:55:25', '2024-12-12 03:55:25');