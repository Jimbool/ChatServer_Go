# ************************************************************
# Sequel Pro SQL dump
# Version 4499
#
# http://www.sequelpro.com/
# https://github.com/sequelpro/sequelpro
#
# Host: 192.168.1.226 (MySQL 5.5.44)
# Database: chatserver
# Generation Time: 2016-02-02 07:10:03 +0000
# ************************************************************


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;


# Dump of table appconfig
# ------------------------------------------------------------

DROP TABLE IF EXISTS `appconfig`;

CREATE TABLE `appconfig` (
  `AppId` varchar(32) NOT NULL COMMENT '应用Id',
  `AppName` varchar(32) NOT NULL COMMENT '应用名称',
  `AppKey` varchar(64) NOT NULL COMMENT '应用Key，用于加密'
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

LOCK TABLES `appconfig` WRITE;
/*!40000 ALTER TABLE `appconfig` DISABLE KEYS */;

INSERT INTO `appconfig` (`AppId`, `AppName`, `AppKey`)
VALUES
	('PKQ','皮卡丘','7DE9DAA1-E87C-FF51-BCE3-15946DBE9462');

/*!40000 ALTER TABLE `appconfig` ENABLE KEYS */;
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
  `ServerGroupId` int(11) NOT NULL COMMENT '服务器组Id',
  `APIName` varchar(32) NOT NULL COMMENT 'API名称',
  `Content` varchar(1024) NOT NULL COMMENT '请求内容',
  `Crdate` datetime NOT NULL COMMENT '请求时间',
  PRIMARY KEY (`Id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


# Dump of table serverconfig
# ------------------------------------------------------------

DROP TABLE IF EXISTS `serverconfig`;

CREATE TABLE `serverconfig` (
  `ServerGroupId` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `SocketServerConfig` varchar(1024) DEFAULT NULL COMMENT 'Socket服务器配置：如IP、Port等',
  `WebServerConfig` varchar(1024) DEFAULT NULL COMMENT 'Web服务器配置，如：端口，映射关系',
  PRIMARY KEY (`ServerGroupId`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

LOCK TABLES `serverconfig` WRITE;
/*!40000 ALTER TABLE `serverconfig` DISABLE KEYS */;

INSERT INTO `serverconfig` (`ServerGroupId`, `SocketServerConfig`, `WebServerConfig`)
VALUES
	(1,'{\n	\"ServerHost\": \"192.168.1.68\",\n	\"ServerPort\": 8001,\n	\"CheckExpireInterval\": 300,\n	\"ClientExpireSeconds\": 60,\n	\"MaxMsgLength\": 100,\n	\"MaxHistoryCount\":10\n}\n','{\n	\"ServerHost\":\"192.168.1.68\",\n	\"ServerPort\": 8002\n}'),
	(2,'{\n	\"ServerHost\": \"192.168.1.68\",\n	\"ServerPort\": 9001,\n	\"CheckExpireInterval\": 300,\n	\"ClientExpireSeconds\": 60,\n	\"MaxMsgLength\": 100,\n	\"MaxHistoryCount\":10\n}\n','{\n	\"ServerHost\":\"192.168.1.68\",\n	\"ServerPort\": 9002\n}');

/*!40000 ALTER TABLE `serverconfig` ENABLE KEYS */;
UNLOCK TABLES;



/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
