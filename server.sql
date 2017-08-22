-- --------------------------------------------------------
-- Хост:                         127.0.0.1
-- Версия сервера:               10.2.7-MariaDB - mariadb.org binary distribution
-- Операционная система:         Win64
-- HeidiSQL Версия:              9.4.0.5125
-- --------------------------------------------------------

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET NAMES utf8 */;
/*!50503 SET NAMES utf8mb4 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;


-- Дамп структуры базы данных ficbook
CREATE DATABASE IF NOT EXISTS `ficbook` /*!40100 DEFAULT CHARACTER SET utf8 */;
USE `ficbook`;

-- Дамп структуры для таблица ficbook.bans_list
CREATE TABLE IF NOT EXISTS `bans_list` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `login_banned` varbinary(100) NOT NULL,
  `login_banning` varbinary(100) NOT NULL,
  `reason` text NOT NULL,
  `time_ban` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  `time_expired` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
  PRIMARY KEY (`id`),
  KEY `bans_list_login_banned_index` (`login_banned`),
  KEY `bans_list_login_banning_index` (`login_banning`)
) ENGINE=InnoDB AUTO_INCREMENT=171 DEFAULT CHARSET=utf8mb4;

-- Экспортируемые данные не выделены.
-- Дамп структуры для таблица ficbook.chat_message_all
CREATE TABLE IF NOT EXISTS `chat_message_all` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `login` varchar(50) NOT NULL,
  `message` text NOT NULL,
  `timestamp` timestamp NOT NULL DEFAULT current_timestamp(),
  `room_uuid` varchar(50) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `chat_message_all_login_index` (`login`),
  KEY `chat_message_all_room_id_index` (`room_uuid`),
  KEY `i_timestamp` (`timestamp`),
  KEY `i_room_id` (`room_uuid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Экспортируемые данные не выделены.
-- Дамп структуры для таблица ficbook.chat_rooms
CREATE TABLE IF NOT EXISTS `chat_rooms` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(50) NOT NULL,
  `topic` varchar(50) NOT NULL DEFAULT 'unknown',
  `about` text NOT NULL DEFAULT 'unknown',
  `type` varchar(50) NOT NULL DEFAULT 'public',
  `uuid` varchar(50) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=76 DEFAULT CHARSET=utf8mb4;

-- Экспортируемые данные не выделены.
-- Дамп структуры для таблица ficbook.subscriptions
CREATE TABLE IF NOT EXISTS `subscriptions` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `login` varbinary(100) NOT NULL,
  `room_id` int(11) NOT NULL,
  `timestamp` timestamp NOT NULL DEFAULT current_timestamp(),
  PRIMARY KEY (`id`),
  KEY `subscriptions_login_index` (`login`),
  KEY `subscriptions_room_id_index` (`room_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Экспортируемые данные не выделены.
-- Дамп структуры для таблица ficbook.users
CREATE TABLE IF NOT EXISTS `users` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `login` varbinary(50) NOT NULL DEFAULT 'unknown_login',
  `password` blob NOT NULL DEFAULT 'unknown_password',
  `power` int(11) DEFAULT -1,
  `date_reg` timestamp NOT NULL DEFAULT current_timestamp(),
  `date_visit` timestamp NOT NULL DEFAULT current_timestamp(),
  PRIMARY KEY (`id`),
  UNIQUE KEY `users_login_unique` (`login`)
) ENGINE=InnoDB AUTO_INCREMENT=447 DEFAULT CHARSET=utf8mb4;

-- Экспортируемые данные не выделены.
/*!40101 SET SQL_MODE=IFNULL(@OLD_SQL_MODE, '') */;
/*!40014 SET FOREIGN_KEY_CHECKS=IF(@OLD_FOREIGN_KEY_CHECKS IS NULL, 1, @OLD_FOREIGN_KEY_CHECKS) */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
