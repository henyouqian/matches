# ************************************************************
# Sequel Pro SQL dump
# Version 4096
#
# http://www.sequelpro.com/
# http://code.google.com/p/sequel-pro/
#
# Host: 127.0.0.1 (MySQL 5.6.10)
# Database: match_db
# Generation Time: 2013-09-10 06:39:39 +0000
# ************************************************************


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;


# Dump of table gametpls
# ------------------------------------------------------------

DROP TABLE IF EXISTS `gametpls`;

CREATE TABLE `gametpls` (
  `id` int(11) unsigned NOT NULL,
  `appid` int(11) NOT NULL DEFAULT '0',
  `name` varchar(40) DEFAULT NULL,
  `order` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '0=ASC; 1=DESC',
  PRIMARY KEY (`appid`,`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



# Dump of table insertTest
# ------------------------------------------------------------

DROP TABLE IF EXISTS `insertTest`;

CREATE TABLE `insertTest` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `a` int(11) DEFAULT NULL,
  `b` int(11) DEFAULT NULL,
  `c` int(11) DEFAULT NULL,
  `d` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



# Dump of table matches
# ------------------------------------------------------------

DROP TABLE IF EXISTS `matches`;

CREATE TABLE `matches` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(60) NOT NULL,
  `appid` int(11) unsigned NOT NULL,
  `gameid` int(11) unsigned NOT NULL,
  `sort` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '0=ASC; 1=DESC',
  `begin` bigint(60) unsigned NOT NULL,
  `end` bigint(60) unsigned NOT NULL,
  PRIMARY KEY (`id`),
  KEY `appid` (`appid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;




/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
