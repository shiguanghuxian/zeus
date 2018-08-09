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

 Date: 10/21/2016 17:08:10 PM
*/

SET NAMES utf8;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
--  Table structure for `zn_appname_field_type`
-- ----------------------------
DROP TABLE IF EXISTS `zn_appname_field_type`;
CREATE TABLE `zn_appname_field_type` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `app_name` varchar(60) DEFAULT 'zn_raw_data' COMMENT '对不同app_name也就是表，设置字段格式',
  `field` varchar(60) DEFAULT NULL COMMENT '字段名',
  `type` varchar(30) DEFAULT 'string' COMMENT '数据类型string,int,int64,float64',
  `unit` varchar(30) DEFAULT NULL COMMENT '字段数据单位',
  `index` int(11) DEFAULT '0' COMMENT '字段是否创建索引,0没有，1有',
  `is_delete` int(11) DEFAULT '0' COMMENT '是否删除，0未删除，1已删除',
  PRIMARY KEY (`id`),
  KEY `field` (`field`),
  KEY `app_name` (`app_name`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8 COMMENT='字段格式配置';

-- ----------------------------
--  Records of `zn_appname_field_type`
-- ----------------------------
BEGIN;
INSERT INTO `zn_appname_field_type` VALUES ('1', 'zn_raw_data', 'code', 'string', null, '0', '0'), ('2', 'zn_raw_data', 'value', 'int', null, '0', '0'), ('3', 'ceshi', 'code', 'int', 'Mb', '0', '1'), ('4', '1', '2', 'string', '3', '1', '1'), ('5', '32', '32', 'float64', 'kb', '0', '1'), ('6', 'zn_hehe', 'value', 'int64', 'Mb', '0', '0');
COMMIT;

-- ----------------------------
--  Table structure for `zn_device`
-- ----------------------------
DROP TABLE IF EXISTS `zn_device`;
CREATE TABLE `zn_device` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `hostname` varchar(200) DEFAULT '' COMMENT '设备主机名',
  `ip` varchar(60) DEFAULT '' COMMENT '设备ip',
  `device_type` varchar(60) DEFAULT '' COMMENT '设备类型,windows,linux,unix或其它自定义',
  `group_name` varchar(60) DEFAULT '' COMMENT '分组信息，上报数据',
  `description` text COMMENT '描述',
  `is_delete` tinyint(4) DEFAULT '0' COMMENT '0未删除，1已删除',
  PRIMARY KEY (`id`),
  KEY `group` (`group_name`),
  KEY `device_type` (`device_type`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8 COMMENT='设备列表本地表';

-- ----------------------------
--  Records of `zn_device`
-- ----------------------------
BEGIN;
INSERT INTO `zn_device` VALUES ('1', 'hostname', '192.168.1.2', 'windows', 'def', '这是描述信息12', '0');
COMMIT;

-- ----------------------------
--  Table structure for `zn_device_group`
-- ----------------------------
DROP TABLE IF EXISTS `zn_device_group`;
CREATE TABLE `zn_device_group` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(60) DEFAULT '' COMMENT '分组名',
  `description` text COMMENT '描述',
  `type` tinyint(4) DEFAULT '0' COMMENT '分组类型，0物理视图，1其它暂定',
  PRIMARY KEY (`id`),
  KEY `type` (`type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='设备分组';

-- ----------------------------
--  Table structure for `zn_device_group_contrast`
-- ----------------------------
DROP TABLE IF EXISTS `zn_device_group_contrast`;
CREATE TABLE `zn_device_group_contrast` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `device_id` int(11) DEFAULT NULL COMMENT '设备id',
  `device_group_id` int(11) DEFAULT NULL COMMENT '分组id',
  PRIMARY KEY (`id`),
  KEY `device_id` (`device_id`),
  KEY `device_group_id` (`device_group_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='设备分组对照';

-- ----------------------------
--  Table structure for `zn_event_level`
-- ----------------------------
DROP TABLE IF EXISTS `zn_event_level`;
CREATE TABLE `zn_event_level` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(30) DEFAULT NULL COMMENT '告警级别名',
  `level` int(11) DEFAULT NULL COMMENT '级别，数值越大，告警越危险',
  PRIMARY KEY (`id`),
  KEY `level` (`level`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8 COMMENT='告警级别';

-- ----------------------------
--  Records of `zn_event_level`
-- ----------------------------
BEGIN;
INSERT INTO `zn_event_level` VALUES ('1', '正常', '0'), ('2', '告警', '1'), ('3', '一般', '2'), ('4', '严重', '3'), ('6', '致命', '4');
COMMIT;

-- ----------------------------
--  Table structure for `zn_event_push`
-- ----------------------------
DROP TABLE IF EXISTS `zn_event_push`;
CREATE TABLE `zn_event_push` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `event_seting_id` int(11) DEFAULT NULL COMMENT '告警设置id',
  `url` varchar(300) DEFAULT NULL COMMENT '告警推送url地址',
  `name` varchar(60) DEFAULT NULL COMMENT '告警推送名',
  `data_type` int(11) DEFAULT '0' COMMENT '0:msg信息，1原始信息',
  PRIMARY KEY (`id`),
  KEY `event_seting_id` (`event_seting_id`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8 COMMENT='告警推送通知地址配置';

-- ----------------------------
--  Records of `zn_event_push`
-- ----------------------------
BEGIN;
INSERT INTO `zn_event_push` VALUES ('1', '1', 'http://www.baidu2.com', '测试推送2', '1'), ('5', '1', 'http://www.sina.com', 'ceshi 1', '0'), ('6', '2', 'http://www.baidu.com', 'ces45', '1');
COMMIT;

-- ----------------------------
--  Table structure for `zn_event_rule`
-- ----------------------------
DROP TABLE IF EXISTS `zn_event_rule`;
CREATE TABLE `zn_event_rule` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `event_level_id` int(11) DEFAULT NULL COMMENT '级别id',
  `event_seting_id` int(11) DEFAULT NULL COMMENT '告警规则id',
  `value` varchar(60) DEFAULT '' COMMENT '比较值',
  `expression` enum('=','>','<','>=','<=','!=') DEFAULT '=' COMMENT '关系，=,>,<,>=,<=,!=',
  `sort` int(11) DEFAULT '0' COMMENT '告警检测顺序',
  `unit` varchar(30) DEFAULT '' COMMENT '单位',
  PRIMARY KEY (`id`),
  KEY `event_level_id` (`event_level_id`),
  KEY `event_seting_id` (`event_seting_id`),
  KEY `value` (`value`),
  KEY `expression` (`expression`),
  KEY `sort` (`sort`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8 COMMENT='告警规则';

-- ----------------------------
--  Records of `zn_event_rule`
-- ----------------------------
BEGIN;
INSERT INTO `zn_event_rule` VALUES ('3', '4', '1', '90', '>', '1', '%'), ('6', '3', '2', '10', '<=', '1', '%');
COMMIT;

-- ----------------------------
--  Table structure for `zn_event_seting`
-- ----------------------------
DROP TABLE IF EXISTS `zn_event_seting`;
CREATE TABLE `zn_event_seting` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(60) DEFAULT NULL COMMENT '告警标题',
  `app_name` varchar(60) DEFAULT NULL COMMENT '应用名',
  `field` varchar(30) DEFAULT NULL COMMENT '判断字段',
  `value_type` enum('当前值','统计值','平均值','最大值','最小值','中位值') DEFAULT NULL COMMENT '告警值类型,当前值，count，平均值，最大值,中位值',
  `describe` text COMMENT '描述信息',
  `continued_count` int(11) DEFAULT '1' COMMENT '出现次数，默认1次，对当前值无效，用redis计数统计',
  `continued_time` int(11) DEFAULT '60' COMMENT '持续时间秒数，步长，对当前值无效',
  `cycle_time` varchar(60) DEFAULT '* */5 * * * *' COMMENT '执行周期，每个多少秒执行一次策略,linux 定时执行格式',
  `event_template_id` int(11) DEFAULT '0' COMMENT '告警模板',
  `enable` int(11) DEFAULT '1' COMMENT '1启用，0不启用',
  `is_delete` int(11) DEFAULT '0' COMMENT '是否删除，0未删除，1已删除',
  PRIMARY KEY (`id`),
  KEY `value_type` (`value_type`),
  KEY `continued_count` (`continued_count`),
  KEY `continued_time` (`continued_time`),
  KEY `cycle_time` (`cycle_time`),
  KEY `event_template_id` (`event_template_id`),
  KEY `app_name` (`app_name`),
  KEY `field` (`field`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8 COMMENT='告警设置表';

-- ----------------------------
--  Records of `zn_event_seting`
-- ----------------------------
BEGIN;
INSERT INTO `zn_event_seting` VALUES ('1', 'name1', 'appname1', 'field1', '统计值', 'desc11', '21', '31', '*/30 * * * * *', '1', '1', '0'), ('2', '测试标题', '应用名', '字段', '统计值', '描述', '1', '60', '* */5 * * * *', '3', '1', '0');
COMMIT;

-- ----------------------------
--  Table structure for `zn_event_template`
-- ----------------------------
DROP TABLE IF EXISTS `zn_event_template`;
CREATE TABLE `zn_event_template` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(30) DEFAULT NULL COMMENT '模板名',
  `content` text COMMENT '告警内容，中间可使用占位符',
  PRIMARY KEY (`id`),
  KEY `name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8 COMMENT='告警模版，用于通知用户，和显示';

-- ----------------------------
--  Records of `zn_event_template`
-- ----------------------------
BEGIN;
INSERT INTO `zn_event_template` VALUES ('1', 'name', '告警触发值{current_value}、告警时间{date}、告警级别{event_level}、主机名{hostname}、ip地址{ip}、分组{group}等'), ('2', 'name1', 'content1'), ('3', 'ee22', 'ewewe22');
COMMIT;

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
) ENGINE=InnoDB AUTO_INCREMENT=35 DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT COMMENT='系统菜单';

-- ----------------------------
--  Records of `zn_menu`
-- ----------------------------
BEGIN;
INSERT INTO `zn_menu` VALUES ('3', '告警管理', 'gaojingguanli', '0', '0', '', 'icon-energy', null, '0', '1'), ('5', '设备管理', 'shebeiguanli', '0', '0', '', 'icon-screen-desktop', null, '0', '1'), ('6', '系统设置', 'canshuguanli', '0', '0', '', 'icon-settings', null, '0', '1'), ('7', '阈值管理', 'yuzhiguanli', '0', '0', '', 'icon-wrench', null, '0', '1'), ('9', '个人设置', 'gerenshezhi', '6', '0', '', 'fa fa-user', null, '0', '1'), ('12', '个人资料', 'gerenziliao', '9', '0', '/user/myuserinfo', 'fa fa-file-o', null, '0', '1'), ('25', '采集规则', 'caijiguize', '6', '0', '', 'fa fa-recycle', null, '0', '1'), ('26', '话题设置', 'huotishezhi', '25', '0', '/topics/index', 'fa fa-cubes', null, '0', '1'), ('27', 'NSQ话题列表', 'nsq huatiliebiao', '25', '0', '/topics/nsqtopics', 'fa fa-deviantart', null, '0', '0'), ('28', '设置字段格式', 'shezhiziduangeshi', '25', '0', '/appname/list', 'fa fa-stumbleupon', null, '0', '1'), ('29', '告警管理', 'gaojingguanli', '6', '0', '', 'fa fa-exclamation-triangle', null, '0', '1'), ('30', '告警设置', 'gaojingshezhi', '29', '0', '/setings/event_seting_list', 'fa fa-bolt', null, '0', '1'), ('31', '告警级别', 'gaojingjibie', '29', '0', '/setings/event_level_list', 'fa fa-warning', null, '0', '1'), ('32', 'ZQL查询', 'zql_search', '0', '0', '/search', 'fa fa-search', null, '0', '1'), ('33', '设备管理', 'shebeiguanli', '5', '0', '/device/list', 'fa fa-desktop', null, '0', '1'), ('34', '设备分组', 'shebeifenzu', '5', '0', '/device/group_list', 'fa fa-folder-open', null, '0', '1');
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
  `is_delete` int(11) DEFAULT '0' COMMENT '是否删除，0未删除，1已删除',
  PRIMARY KEY (`id`),
  KEY `topics` (`topics`),
  KEY `enable` (`enable`),
  KEY `is_delete` (`is_delete`),
  KEY `channel` (`channel`),
  KEY `data_type` (`data_type`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8 COMMENT='配置nsq消息处理规则';

-- ----------------------------
--  Records of `zn_topics_config`
-- ----------------------------
BEGIN;
INSERT INTO `zn_topics_config` VALUES ('1', 'test', 'abc', '2', '1', 'text', '0'), ('2', 'test1', 'abc1', '2', '1', 'text', '0'), ('3', 'test0', 'abc0', '2', '1', 'json', '0');
COMMIT;

-- ----------------------------
--  Table structure for `zn_topics_config_rule`
-- ----------------------------
DROP TABLE IF EXISTS `zn_topics_config_rule`;
CREATE TABLE `zn_topics_config_rule` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `mapped` varchar(255) DEFAULT '' COMMENT 'agent上传数据字段与本地数据库字段映射，入库使用',
  `text_un_type` enum('regular','char') DEFAULT 'char' COMMENT '文本类型数据解析规则，有正则regular，特殊字符char（一个字符）',
  `text_un_rule` varchar(255) DEFAULT NULL COMMENT '文本解析规则',
  `date_format` varchar(60) DEFAULT '2006-01-02 15:04:05' COMMENT '日期格式化规则',
  `topics_config_id` int(11) DEFAULT NULL COMMENT '话题配置id',
  `app_name` varchar(60) DEFAULT 'zn_raw_data' COMMENT '应用名，实际是存储表名,查询时要选择appname',
  `tag` varchar(30) DEFAULT '' COMMENT '标签，用于选择appname也就是那条解析规则',
  `enable` int(11) DEFAULT '1' COMMENT '0未启用，1启用',
  `sort` int(11) DEFAULT '0' COMMENT '排序',
  PRIMARY KEY (`id`),
  KEY `topics_config_id` (`topics_config_id`),
  KEY `enable` (`enable`),
  KEY `app_name` (`app_name`),
  KEY `sort` (`sort`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8 COMMENT='话题规则配置';

-- ----------------------------
--  Records of `zn_topics_config_rule`
-- ----------------------------
BEGIN;
INSERT INTO `zn_topics_config_rule` VALUES ('1', '0:kpiid|1:value|2:date|3:instance|4:code', 'char', '|', '2006-01-02 15:04:05', '1', 'zn_raw_data', '', '1', '0'), ('2', 'kpiid:kpiid|value:value|date:date|instance:instance', 'regular', 'kpiid&(?P<kpiid>\\d*)&jjsd\\|value:(?P<value>.*)\\|date%(?P<date>.*)\\|hehe\\*(?P<instance>.*)@dsdsd', '2006-01-02 15:04:05', '2', 'zn_raw_data', '', '1', '0'), ('3', 'kpiid:kpiid|value:value|date:date|instance:instance', 'char', '', '2006-01-02 15:04:05', '3', 'zn_raw_data', '', '1', '0'), ('8', 'kpiid:kpiid|value:value|date:date|instance:instance|code:code', 'char', null, '2006-01-02 15:04:05', '3', 'zn_raw_data', '', '1', '1');
COMMIT;

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
  `language` varchar(30) DEFAULT 'zh-CN' COMMENT '语言',
  PRIMARY KEY (`id`),
  UNIQUE KEY `username` (`username`),
  KEY `password` (`password`),
  KEY `phone` (`phone`),
  KEY `state` (`state`),
  KEY `group_id` (`group_id`),
  KEY `level` (`level`),
  KEY `is_delete` (`is_delete`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='用户表';

-- ----------------------------
--  Records of `zn_user`
-- ----------------------------
BEGIN;
INSERT INTO `zn_user` VALUES ('1', 'admin', '69b1db5f2c43e65c2b359bdfd7caed2d502cc43d', '左秀朋', 'zuoxiupeng', '/static/img/avatar.png', '1', '13241812313', '1361653339@qq.com', null, '1440485116', '1', '0', '简介', '1', '1468831669', '0', 'zh-CN');
COMMIT;

SET FOREIGN_KEY_CHECKS = 1;
