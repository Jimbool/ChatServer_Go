# ************************************************************
# Sequel Pro SQL dump
# Version 4499
#
# http://www.sequelpro.com/
# https://github.com/sequelpro/sequelpro
#
# Host: 192.168.1.226 (MySQL 5.5.44)
# Database: chatserver
# Generation Time: 2016-01-25 02:24:13 +0000
# ************************************************************


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;


# Dump of table config
# ------------------------------------------------------------

DROP TABLE IF EXISTS `config`;

CREATE TABLE `config` (
  `AppId` varchar(32) NOT NULL COMMENT '应用Id',
  `AppName` varchar(32) NOT NULL COMMENT '应用名称',
  `AppKey` varchar(64) NOT NULL COMMENT '应用Key，用于加密',
  `SocketServerConfig` varchar(1024) NOT NULL COMMENT 'Socket服务器配置：如IP、Port等',
  `WebServerConfig` varchar(1024) NOT NULL COMMENT 'Web服务器配置，如：端口，映射关系'
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

LOCK TABLES `config` WRITE;
/*!40000 ALTER TABLE `config` DISABLE KEYS */;

INSERT INTO `config` (`AppId`, `AppName`, `AppKey`, `SocketServerConfig`, `WebServerConfig`)
VALUES
	('1','皮卡丘','7DE9DAA1-E87C-FF51-BCE3-15946DBE9462','{\n	\"ServerHost\": \"192.168.1.68\",\n	\"ServerPort\": 8001,\n	\"CheckExpireInterval\": 300,\n	\"ClientExpireSeconds\": 60,\n	\"MaxMsgLength\": 100\n}','{\n	\"ServerHost\":\"192.168.1.68\",\n	\"ServerPort\": 8002\n}');

/*!40000 ALTER TABLE `config` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table player
# ------------------------------------------------------------

DROP TABLE IF EXISTS `player`;

CREATE TABLE `player` (
  `Id` varchar(64) NOT NULL COMMENT '玩家Id',
  `Name` varchar(64) NOT NULL COMMENT '玩家名称',
  `UnionId` varchar(64) DEFAULT NULL COMMENT '公会Id',
  `ExtraMsg` varchar(1024) DEFAULT NULL COMMENT '额外透传信息',
  `RegisterTime` datetime NOT NULL COMMENT '注册时间',
  `LoginTime` datetime NOT NULL COMMENT '登录时间',
  `IsForbidden` tinyint(1) NOT NULL COMMENT '是否封号',
  `SilentEndTime` datetime NOT NULL COMMENT '禁言结束时间',
  PRIMARY KEY (`Id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


# Dump of table request_log
# ------------------------------------------------------------

DROP TABLE IF EXISTS `request_log`;

CREATE TABLE `request_log` (
  `Id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'Id',
  `APIName` varchar(32) NOT NULL COMMENT 'API名称',
  `Content` varchar(1024) NOT NULL COMMENT '请求内容',
  `Crdate` datetime NOT NULL COMMENT '请求时间',
  PRIMARY KEY (`Id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


# Dump of table sensitivewords
# ------------------------------------------------------------

DROP TABLE IF EXISTS `sensitivewords`;

CREATE TABLE `sensitivewords` (
  `Text` varchar(32) NOT NULL COMMENT '内容',
  PRIMARY KEY (`Text`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
