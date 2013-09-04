# ************************************************************
# Sequel Pro SQL dump
# Version 4096
#
# http://www.sequelpro.com/
# http://code.google.com/p/sequel-pro/
#
# Host: 127.0.0.1 (MySQL 5.5.32-0ubuntu0.12.04.1)
# Database: match_db
# Generation Time: 2013-09-03 08:56:03 +0000
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



# Dump of table matches
# ------------------------------------------------------------

DROP TABLE IF EXISTS `matches`;

CREATE TABLE `matches` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(60) NOT NULL,
  `appid` int(11) NOT NULL,
  `gameid` int(11) NOT NULL,
  `sort` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0=ASC; 1=DESC',
  `begin` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' ON UPDATE CURRENT_TIMESTAMP,
  `end` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;




/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
