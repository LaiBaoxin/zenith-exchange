/*
 Navicat Premium Data Transfer

 Source Server         : localhostMySQL8
 Source Server Type    : MySQL
 Source Server Version : 50729 (5.7.29)
 Source Host           : localhost:3306
 Source Schema         : cex

 Target Server Type    : MySQL
 Target Server Version : 50729 (5.7.29)
 File Encoding         : 65001

 Date: 30/04/2026 20:20:58
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for accounts
-- ----------------------------
DROP TABLE IF EXISTS `accounts`;
CREATE TABLE `accounts` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint(20) unsigned DEFAULT NULL,
  `currency` varchar(20) DEFAULT NULL,
  `available` decimal(36,18) DEFAULT '0.000000000000000000',
  `frozen` decimal(36,18) DEFAULT '0.000000000000000000',
  `version` int(10) unsigned NOT NULL DEFAULT '0',
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_currency` (`user_id`,`currency`),
  CONSTRAINT `fk_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`),
  CONSTRAINT `fk_users_accounts` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for balance_logs
-- ----------------------------
DROP TABLE IF EXISTS `balance_logs`;
CREATE TABLE `balance_logs` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint(20) unsigned NOT NULL COMMENT '用户ID',
  `currency` varchar(20) NOT NULL COMMENT '币种',
  `change_type` varchar(50) NOT NULL COMMENT '变更类型: freeze, unfreeze, trade等',
  `amount` decimal(40,18) NOT NULL COMMENT '变动金额',
  `balance` decimal(40,18) NOT NULL COMMENT '变动后的可用余额',
  `log_time` datetime(3) NOT NULL COMMENT '记录时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_currency` (`user_id`,`currency`)
) ENGINE=InnoDB AUTO_INCREMENT=49 DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for orders
-- ----------------------------
DROP TABLE IF EXISTS `orders`;
CREATE TABLE `orders` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint(20) unsigned NOT NULL,
  `symbol` varchar(20) NOT NULL,
  `side` enum('buy','sell') NOT NULL,
  `type` enum('limit','market') NOT NULL DEFAULT 'limit',
  `price` decimal(36,18) NOT NULL,
  `amount` decimal(36,18) NOT NULL,
  `filled_amount` decimal(36,18) NOT NULL DEFAULT '0.000000000000000000',
  `status` tinyint(4) NOT NULL DEFAULT '0',
  `msg_hash` char(66) DEFAULT NULL,
  `signature` text,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `is_mock` tinyint(4) DEFAULT '0',
  `updated_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_user_status` (`user_id`,`status`),
  KEY `idx_symbol_status` (`symbol`,`status`)
) ENGINE=InnoDB AUTO_INCREMENT=12 DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for users
-- ----------------------------
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `wallet_address` char(42) NOT NULL,
  `api_key` varchar(64) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uni_users_api_key` (`api_key`),
  KEY `idx_wallet` (`wallet_address`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4;

SET FOREIGN_KEY_CHECKS = 1;
