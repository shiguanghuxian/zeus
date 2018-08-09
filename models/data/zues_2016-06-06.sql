/*
 Navicat Premium Data Transfer

 Source Server         : localhost
 Source Server Type    : MySQL
 Source Server Version : 50624
 Source Host           : localhost
 Source Database       : zues

 Target Server Type    : MySQL
 Target Server Version : 50624
 File Encoding         : utf-8

 Date: 06/06/2016 11:15:59 AM
*/

SET NAMES utf8;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
--  Table structure for `zn_topics_config`
-- ----------------------------
DROP TABLE IF EXISTS `zn_topics_config`;
CREATE TABLE `zn_topics_config` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `topics` char(8) NOT NULL COMMENT '话题名，使用不超过8个字符',
  `channel` varchar(20) NOT NULL COMMENT '读取会话的通道',
  `channel_count` int(11) DEFAULT '1' COMMENT '该会话使用多少个同名通道取取数据',
  `enable` tinyint(4) DEFAULT '1' COMMENT '1启用，0不启用',
  `data_type` char(4) DEFAULT 'json' COMMENT '这个会话使用的数据传输类型(json,text)，默认json',
  `mapped` varchar(255) DEFAULT '' COMMENT 'agent上传数据字段与本地数据库字段映射，入库使用',
  `text_un_type` enum('regular','char') DEFAULT 'char' COMMENT '文本类型数据解析规则，有正则regular，特殊字符char（一个字符）',
  `text_un_rule` varchar(255) DEFAULT '' COMMENT '文本解析规则',
  PRIMARY KEY (`id`),
  KEY `topics` (`topics`),
  KEY `enable` (`enable`)
) ENGINE=MyISAM AUTO_INCREMENT=4 DEFAULT CHARSET=utf8 COMMENT='配置nsq消息处理规则';

-- ----------------------------
--  Records of `zn_topics_config`
-- ----------------------------
BEGIN;
INSERT INTO `zn_topics_config` VALUES ('1', 'test', 'abc', '2', '1', 'text', '0:kpiid|1:value|2:date|3:instance', 'char', '|'), ('2', 'test1', 'abc1', '2', '1', 'text', 'kpiid:kpiid|value:value|date:date|instance:instance', 'regular', 'kpiid&(?P<kpiid>\\d*)&jjsd\\|value:(?P<value>.*)\\|date%(?P<date>.*)\\|hehe\\*(?P<instance>.*)@dsdsd'), ('3', 'test0', 'abc0', '2', '1', 'json', 'kpiid:kpiid|value:value|date:date|instance:instance', 'char', '');
COMMIT;

SET FOREIGN_KEY_CHECKS = 1;
