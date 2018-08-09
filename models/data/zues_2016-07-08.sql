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

 Date: 07/08/2016 14:52:17 PM

修改话题配置格式

*/

SET NAMES utf8;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
--  Table structure for `zn_appname_field_type`
-- ----------------------------
DROP TABLE IF EXISTS `zn_appname_field_type`;
CREATE TABLE `zn_appname_field_type` (
  `id` int(11) NOT NULL,
  `app_name` varchar(60) DEFAULT 'zn_raw_data' COMMENT '对不同app_name也就是表，设置字段格式',
  `field` varchar(60) DEFAULT NULL COMMENT '字段名',
  `type` varchar(30) DEFAULT 'string' COMMENT '数据类型',
  PRIMARY KEY (`id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8 COMMENT='字段格式配置';

-- ----------------------------
--  Table structure for `zn_menu`
-- ----------------------------
DROP TABLE IF EXISTS `zn_menu`;
CREATE TABLE `zn_menu` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(60) CHARACTER SET utf8mb4 DEFAULT NULL COMMENT '显示名',
  `name_en` varchar(60) CHARACTER SET utf8mb4 DEFAULT NULL COMMENT '英文标题',
  `parent_id` int(11) DEFAULT '0' COMMENT '父级菜单id',
  `sort` int(11) DEFAULT '0' COMMENT '排序，倒叙',
  `url` varchar(180) CHARACTER SET utf8mb4 DEFAULT NULL,
  `icon` varchar(255) CHARACTER SET utf8mb4 DEFAULT NULL COMMENT '图片地址',
  `other` varchar(255) CHARACTER SET utf8mb4 DEFAULT NULL COMMENT '其它信息',
  `role_type_id` int(11) DEFAULT '0' COMMENT '角色类型',
  `show` int(11) DEFAULT '1' COMMENT '1显示，0隐藏',
  PRIMARY KEY (`id`),
  KEY `parent_id` (`parent_id`),
  KEY `sort` (`sort`),
  KEY `role_type_id` (`role_type_id`),
  KEY `url` (`url`),
  KEY `show` (`show`)
) ENGINE=InnoDB AUTO_INCREMENT=28 DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT COMMENT='beego  菜单';

-- ----------------------------
--  Records of `zn_menu`
-- ----------------------------
BEGIN;
INSERT INTO `zn_menu` VALUES ('1', '席位登录', 'xiweidenglu', '0', '0', '', 'icon-users', null, '0', '1'), ('2', '程序化交易', 'chengxuhuajiaoyi', '0', '0', '', 'icon-hourglass', null, '0', '1'), ('3', '告警管理', 'gaojingguanli', '0', '0', '', 'icon-energy', null, '0', '1'), ('4', '统计查询', 'tongjichaxun', '0', '0', '', 'icon-docs', null, '0', '1'), ('5', '设备管理', 'shebeiguanli', '0', '0', '', 'icon-screen-desktop', null, '0', '1'), ('6', '系统设置', 'canshuguanli', '0', '0', '', 'icon-settings', null, '0', '1'), ('7', '阈值管理', 'yuzhiguanli', '0', '0', '', 'icon-wrench', null, '0', '1'), ('8', '图表设置', 'tubiaoshezhi', '6', '0', '', 'fa fa-bar-chart-o', null, '0', '0'), ('9', '个人设置', 'gerenshezhi', '6', '0', '', 'fa fa-user', null, '0', '1'), ('10', '安全管理', 'anquanguanli', '6', '0', '', 'fa fa-shield', null, '0', '0'), ('11', '首页图表设置', 'homesetchart', '8', '0', '/settings/homechart', 'fa fa-file-o', null, '0', '0'), ('12', '个人资料', 'gerenziliao', '9', '0', '/user/myuserinfo', 'fa fa-file-o', null, '0', '1'), ('13', '人员管理', 'renyuanguanli', '10', '0', '/user/userlist', 'fa fa-user', null, '0', '0'), ('14', '角色管理', 'jueseguanli', '10', '0', '', 'fa fa-group', null, '0', '0'), ('15', '组织机构', 'zuzhijigou', '10', '0', '/user/comefrom', 'fa fa-slack', null, '0', '0'), ('16', '会员告警', 'huiyuangaojing', '3', '0', '', 'fa fa-user', null, '0', '0'), ('17', '采集端告警', 'caijiduangaojing', '3', '0', '', null, null, '0', '0'), ('18', '本地告警', 'bendigaojing', '3', '0', '', null, null, '0', '0'), ('19', '席位频繁登录监控', 'xiweipinfandenglu', '1', '0', '/seatlogin/loginchart', 'fa fa-circle-o-notch', null, '0', '0'), ('20', '告警列表', 'gaojingliebiao', '3', '0', '/event/eventlist', 'fa fa-bolt', null, '0', '0'), ('21', '席位报文速率监控', 'xiweibaowensulu', '2', '0', '/transaction/chart', 'fa fa-bolt', null, '0', '0'), ('22', '实时 登录/退出', 'loginandlogout', '1', '1', '/seatlogin/secloginchart', 'fa fa-bar-chart-o', null, '0', '0'), ('23', '席位报文速率', 'xiweibaowensulv', '2', '1', '/transaction/ratechart', 'fa fa-bar-chart-o', null, '0', '0'), ('24', '错单速率', 'cuodansulv', '1', '0', '/seatlogin/showchart', 'fa fa-bar-chart-o', null, '0', '0'), ('25', '采集规则', 'caijiguize', '6', '0', '', 'fa fa-recycle', null, '0', '1'), ('26', '话题设置', 'huotishezhi', '25', '0', '/topics/index', 'fa fa-cubes', null, '0', '1'), ('27', 'NSQ话题列表', 'nsq huatiliebiao', '25', '0', '/topics/nsqtopics', 'fa fa-deviantart', null, '0', '1');
COMMIT;

-- ----------------------------
--  Table structure for `zn_topics_config`
-- ----------------------------
DROP TABLE IF EXISTS `zn_topics_config`;
CREATE TABLE `zn_topics_config` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `topics` char(8) NOT NULL COMMENT '话题名，使用不超过8个字符',
  `channel` varchar(20) NOT NULL COMMENT '读取会话的通道',
  `channel_count` int(11) DEFAULT '1' COMMENT '该会话使用多少个同名通道取取数据',
  `enable` int(11) DEFAULT '1' COMMENT '1启用，0不启用',
  `data_type` char(4) DEFAULT 'json' COMMENT '这个会话使用的数据传输类型(json,text)，默认json',
  `mapped` varchar(255) DEFAULT '' COMMENT 'agent上传数据字段与本地数据库字段映射，入库使用',
  `text_un_type` enum('regular','char') DEFAULT 'char' COMMENT '文本类型数据解析规则，有正则regular，特殊字符char（一个字符）',
  `text_un_rule` varchar(255) DEFAULT '' COMMENT '文本解析规则',
  `is_delete` int(11) DEFAULT '0' COMMENT '是否删除，0未删除，1已删除',
  PRIMARY KEY (`id`),
  KEY `topics` (`topics`),
  KEY `enable` (`enable`)
) ENGINE=MyISAM AUTO_INCREMENT=9 DEFAULT CHARSET=utf8 COMMENT='配置nsq消息处理规则';

-- ----------------------------
--  Records of `zn_topics_config`
-- ----------------------------
BEGIN;
INSERT INTO `zn_topics_config` VALUES ('1', 'test', 'abc', '2', '1', 'text', '0:kpiid|1:value|2:date|3:instance|4:code', 'char', '|', '0'), ('2', 'test1', 'abc1', '2', '1', 'text', 'kpiid:kpiid|value:value|date:date|instance:instance', 'regular', 'kpiid&(?P<kpiid>\\d*)&jjsd\\|value:(?P<value>.*)\\|date%(?P<date>.*)\\|hehe\\*(?P<instance>.*)@dsdsd', '0'), ('3', 'test0', 'abc0', '2', '1', 'json', 'kpiid:kpiid|value:value|date:date|instance:instance|code:code', 'char', '', '0'), ('6', '12', '1221', '2', '1', 'json', '12', 'regular', '', '1'), ('7', '23', '23', '3', '0', 'json', '23', 'regular', '323', '1'), ('8', '343', '323', '33', '1', 'json', '323', 'regular', '1212', '1');
COMMIT;

-- ----------------------------
--  Table structure for `zn_topics_config_rule`
-- ----------------------------
DROP TABLE IF EXISTS `zn_topics_config_rule`;
CREATE TABLE `zn_topics_config_rule` (
  `id` int(11) NOT NULL,
  `mapped` varchar(255) DEFAULT '' COMMENT 'agent上传数据字段与本地数据库字段映射，入库使用',
  `text_un_type` enum('regular','char') DEFAULT 'char' COMMENT '文本类型数据解析规则，有正则regular，特殊字符char（一个字符）',
  `text_un_rule` varchar(255) DEFAULT NULL COMMENT '文本解析规则',
  `date_format` varchar(60) DEFAULT '2006-01-02 15:04:05' COMMENT '日期格式化规则',
  `topics_config_id` int(11) DEFAULT NULL COMMENT '话题配置id',
  `app_name` varchar(60) DEFAULT 'zn_raw_data' COMMENT '应用名，实际是存储表名,查询时要选择appname',
  PRIMARY KEY (`id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8 COMMENT='话题规则配置';

-- ----------------------------
--  Table structure for `zn_user`
-- ----------------------------
DROP TABLE IF EXISTS `zn_user`;
CREATE TABLE `zn_user` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(30) DEFAULT NULL COMMENT '用户名',
  `password` varchar(40) DEFAULT NULL,
  `name` varchar(30) DEFAULT NULL COMMENT '人员名',
  `name_en` varchar(60) DEFAULT NULL COMMENT '用户英文名',
  `icon` varchar(300) DEFAULT NULL COMMENT '头像地址',
  `sex` char(1) DEFAULT '0' COMMENT '性别0位置，1男，2女',
  `phone` char(11) DEFAULT NULL COMMENT '手机号',
  `email` varchar(300) DEFAULT NULL COMMENT '用户邮箱',
  `token` char(40) DEFAULT NULL COMMENT '手机登录标识',
  `addtime` int(11) DEFAULT NULL COMMENT '注册时间',
  `state` enum('0','1') DEFAULT '1' COMMENT '1在职，0离职－－－工作状态',
  `group_id` int(11) DEFAULT '0',
  `info` text COMMENT '简介',
  `level` tinyint(4) DEFAULT NULL COMMENT '用户级别，1一线人员，2二线人员，3经理，4副经理，5总经理',
  `uptime` int(11) DEFAULT NULL COMMENT '修改时间',
  `is_delete` enum('0','1') DEFAULT '0' COMMENT '是否删除，0未删除，1已删除',
  PRIMARY KEY (`id`),
  UNIQUE KEY `username` (`username`),
  KEY `password` (`password`),
  KEY `phone` (`phone`),
  KEY `state` (`state`),
  KEY `group_id` (`group_id`),
  KEY `level` (`level`),
  KEY `is_delete` (`is_delete`)
) ENGINE=MyISAM AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='用户表';

-- ----------------------------
--  Records of `zn_user`
-- ----------------------------
BEGIN;
INSERT INTO `zn_user` VALUES ('1', 'admin', '69b1db5f2c43e65c2b359bdfd7caed2d502cc43d', '左秀朋', 'zuoxiupeng', '/static/img/avatar.png', '1', '13241812313', '1361653339@qq.com', null, '1440485116', '1', '0', '简介', '1', '1465799056', '0');
COMMIT;

SET FOREIGN_KEY_CHECKS = 1;
