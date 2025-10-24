-- phpMyAdmin SQL Dump
-- version 5.2.1
-- https://www.phpmyadmin.net/
--
-- 主机： localhost
-- 生成日期： 2025-10-24 09:52:43
-- 服务器版本： 8.0.33
-- PHP 版本： 8.3.6

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- 数据库： `geekai_agent`
--
CREATE DATABASE IF NOT EXISTS `geekai_agent` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
USE `geekai_agent`;

-- --------------------------------------------------------

--
-- 表的结构 `geekai_admin_users`
--

DROP TABLE IF EXISTS `geekai_admin_users`;
CREATE TABLE `geekai_admin_users` (
  `id` int NOT NULL,
  `username` varchar(30) NOT NULL COMMENT '用户名',
  `password` char(64) NOT NULL COMMENT '密码',
  `salt` char(12) NOT NULL COMMENT '密码盐',
  `status` tinyint(1) NOT NULL COMMENT '当前状态',
  `last_login_at` bigint NOT NULL COMMENT '最后登录时间',
  `last_login_ip` char(16) NOT NULL COMMENT '最后登录 IP',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime NOT NULL COMMENT '更新时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='系统用户' ROW_FORMAT=DYNAMIC;

--
-- 转存表中的数据 `geekai_admin_users`
--

INSERT INTO `geekai_admin_users` (`id`, `username`, `password`, `salt`, `status`, `last_login_at`, `last_login_ip`, `created_at`, `updated_at`) VALUES
(1, 'admin', '6d17e80c87d209efb84ca4b2e0824f549d09fac8b2e1cc698de5bb5e1d75dfd0', 'mmrql75o', 1, 1761299213, '127.0.0.1', '2024-03-11 16:30:20', '2025-10-24 17:46:53');

-- --------------------------------------------------------

--
-- 表的结构 `geekai_apps`
--

DROP TABLE IF EXISTS `geekai_apps`;
CREATE TABLE `geekai_apps` (
  `id` int NOT NULL,
  `name` varchar(30) DEFAULT NULL COMMENT '名称',
  `type` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT 'openai' COMMENT 'openai,dify,coze',
  `bot_id` varchar(30) NOT NULL COMMENT '机器人ID（coze 专用）',
  `enabled` tinyint(1) DEFAULT NULL COMMENT '是否启用',
  `configs` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '智能体配置参数',
  `score` bigint NOT NULL DEFAULT '0' COMMENT '单次对话消耗积分',
  `summary` varchar(512) DEFAULT NULL COMMENT '应用简介',
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `cid` int NOT NULL COMMENT '分类ID',
  `check` tinyint NOT NULL DEFAULT '0' COMMENT '审核状态 0:未审核 1:审核通过 -1:审核不通过',
  `check_note` varchar(255) DEFAULT NULL COMMENT '审核备注',
  `creator_id` int NOT NULL DEFAULT '0' COMMENT '创作者ID',
  `is_hot` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否热门',
  `use_count` int NOT NULL DEFAULT '0' COMMENT '使用次数',
  `icon` varchar(255) DEFAULT NULL COMMENT '应用图标'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='OpenAI API ';

-- --------------------------------------------------------

--
-- 表的结构 `geekai_app_categories`
--

DROP TABLE IF EXISTS `geekai_app_categories`;
CREATE TABLE `geekai_app_categories` (
  `id` bigint UNSIGNED NOT NULL,
  `name` varchar(30) NOT NULL COMMENT '分类名称',
  `enabled` tinyint NOT NULL DEFAULT '0' COMMENT '状态',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime NOT NULL COMMENT '更新时间',
  `creator_id` bigint NOT NULL COMMENT '创作者ID'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- 转存表中的数据 `geekai_app_categories`
--

INSERT INTO `geekai_app_categories` (`id`, `name`, `enabled`, `created_at`, `updated_at`, `creator_id`) VALUES
(5, '绘图', 1, '2025-05-25 18:05:15', '2025-05-25 18:19:30', 0),
(6, '教育', 1, '2025-05-25 18:13:45', '2025-05-25 18:19:23', 0),
(7, '视频', 1, '2025-05-25 18:13:48', '2025-05-25 18:19:40', 0),
(8, '育儿', 1, '2025-06-15 18:17:04', '2025-06-15 18:17:04', 0),
(11, '视频', 1, '2025-07-07 15:58:03', '2025-07-07 16:11:33', 1),
(12, '教育', 1, '2025-07-07 15:59:21', '2025-07-30 11:16:33', 1),
(16, '数字人', 1, '2025-07-07 16:16:57', '2025-07-30 11:17:31', 1);

-- --------------------------------------------------------

--
-- 表的结构 `geekai_chat_items`
--

DROP TABLE IF EXISTS `geekai_chat_items`;
CREATE TABLE `geekai_chat_items` (
  `id` int NOT NULL,
  `chat_id` char(40) NOT NULL COMMENT '会话 ID',
  `conversation_id` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '会话ID(coze/dify)',
  `user_id` bigint NOT NULL COMMENT '用户 ID',
  `app_id` bigint NOT NULL COMMENT '智能体ID',
  `title` varchar(100) NOT NULL COMMENT '会话标题',
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `icon` varchar(255) NOT NULL COMMENT '图标地址'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户会话列表';

-- --------------------------------------------------------

--
-- 表的结构 `geekai_chat_messages`
--

DROP TABLE IF EXISTS `geekai_chat_messages`;
CREATE TABLE `geekai_chat_messages` (
  `id` bigint NOT NULL,
  `user_id` bigint NOT NULL COMMENT '用户 ID',
  `chat_id` char(40) NOT NULL COMMENT '会话 ID',
  `role` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT 'user or ai',
  `app_id` bigint NOT NULL COMMENT '智能体ID',
  `content` text NOT NULL COMMENT '聊天内容',
  `tokens` smallint NOT NULL COMMENT '耗费 token 数量',
  `use_context` tinyint(1) NOT NULL COMMENT '是否允许作为上下文语料',
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='聊天信息';

-- --------------------------------------------------------

--
-- 表的结构 `geekai_configs`
--

DROP TABLE IF EXISTS `geekai_configs`;
CREATE TABLE `geekai_configs` (
  `id` int NOT NULL,
  `name` varchar(20) NOT NULL COMMENT '配置名称',
  `value` text NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- 转存表中的数据 `geekai_configs`
--

INSERT INTO `geekai_configs` (`id`, `name`, `value`) VALUES
(1, 'system', '{\"title\":\"GeekAI 智能体\",\"slogan\":\"让每一个人都能用好AI\",\"admin_title\":\"GeekAgent 控制台\",\"logo\":\"/images/logo.png\",\"copyright\":\"极客学长\",\"init_score\":100,\"daily_score\":10,\"invite_score\":100,\"enabled_register\":true,\"wechat_card_url\":\"/images/wx.png\",\"email_white_list\":[\"qq.com\",\"163.com\",\"gmail.com\",\"hotmail.com\",\"126.com\",\"outlook.com\",\"foxmail.com\",\"yahoo.com\",\"pvc123.com\"],\"app_id\":61}'),
(3, 'notice', '{\"content\":\"## Geek-Agent v1.0.4 更新日志\\n- Bug 修复：管理后台的超级管理员不能被删除和禁用。\\n- 功能优化：登录组件增加用户协议和隐私政策提示。\\n- 功能新增：支持创作者模式，创作者可以创建自己的智能体应用，并设置价格。🎉🎉🎉。\\n- 功能新增：增加创作者提现功能，创作者可以提现自己的收益。🎉🎉🎉。\\n- 功能新增：完成对话页面重新生成功能。\\n- 功能新增：新增对话分享功能页面。\\n\\n\\n项目介绍：[https://docs.geekai.me](https://docs.geekai.me/agent/)\\n部署教程：[https://docs.geekai.me/agent/install.html](https://docs.geekai.me/agent/install.html)\"}');

-- --------------------------------------------------------

--
-- 表的结构 `geekai_creators`
--

DROP TABLE IF EXISTS `geekai_creators`;
CREATE TABLE `geekai_creators` (
  `id` bigint UNSIGNED NOT NULL COMMENT '主键ID',
  `user_id` int NOT NULL COMMENT '关联用户ID',
  `name` varchar(100) NOT NULL COMMENT '创作者名称',
  `description` varchar(512) DEFAULT NULL COMMENT '创作者简介',
  `logo` varchar(255) DEFAULT NULL COMMENT '创作者Logo',
  `enabled` tinyint(1) NOT NULL DEFAULT '1' COMMENT '是否启用',
  `scores` bigint DEFAULT '0' COMMENT '积分',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime NOT NULL COMMENT '更新时间',
  `check` tinyint(1) NOT NULL DEFAULT '0' COMMENT '审核状态 0:未审核 1:审核通过 2:审核不通过	',
  `check_note` varchar(255) DEFAULT NULL COMMENT '审核备注',
  `withdraw_configs` text COMMENT '提现配置',
  `fee` smallint DEFAULT '0' COMMENT '提现费率(0-100)',
  `username` varchar(30) DEFAULT NULL COMMENT '用户名'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- 转存表中的数据 `geekai_creators`
--

INSERT INTO `geekai_creators` (`id`, `user_id`, `name`, `description`, `logo`, `enabled`, `scores`, `created_at`, `updated_at`, `check`, `check_note`, `withdraw_configs`, `fee`, `username`) VALUES
(1, 4, '极客学长AIGC', '智能体创客不只是技术操作者，更是规则重构者：他们优化智能体的创作逻辑，定义人机协作的新范式，在艺术、设计、教育等领域催生 “人机共创” 的新业态，让 AI 的效率与人类的温度碰撞出更具想象力的火花。', '/static/upload/2025/7/1752476354833413.png', 1, 1052, '2025-06-25 10:54:36', '2025-07-31 09:33:56', 1, '审核费用支付失败', '{\"account\":\"yangjian\",\"method\":\"alipay\",\"mobile\":\"18575670125\",\"name\":\"阳建\",\"qrcode\":\"/static/upload/2025/7/1752044716865779.jpg\",\"score_to_rmb_ratio\":1000}', 15, 'yangjian'),
(3, 57, '暴躁的叶孤城@661308', '我是创作者凑满10个字', '/images/avatar/user23.png', 1, 0, '2025-07-09 18:57:33', '2025-07-10 15:16:37', 1, '测试审核通过', '', 0, '3'),
(5, 55, '快乐的周芷若@983067', '18591927365测试', '/images/avatar/user38.png', 1, 0, '2025-07-31 16:28:03', '2025-07-31 16:28:03', 0, '', '', 0, 'rubz8ruqtnr7');

-- --------------------------------------------------------

--
-- 表的结构 `geekai_creator_score_logs`
--

DROP TABLE IF EXISTS `geekai_creator_score_logs`;
CREATE TABLE `geekai_creator_score_logs` (
  `id` bigint UNSIGNED NOT NULL,
  `user_id` bigint NOT NULL COMMENT '用户ID',
  `creator_id` bigint NOT NULL COMMENT '创作者ID',
  `app_id` bigint NOT NULL DEFAULT '0' COMMENT '应用ID',
  `type` char(20) NOT NULL COMMENT '类型（income：收入，withdraw：提现）',
  `score` int NOT NULL COMMENT '积分数值',
  `balance` bigint NOT NULL COMMENT '余额',
  `subject` varchar(50) NOT NULL COMMENT '主题',
  `remark` varchar(512) NOT NULL COMMENT '备注',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `mark` tinyint(1) NOT NULL COMMENT '资金类型（0：支出，1：收入）'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- --------------------------------------------------------

--
-- 表的结构 `geekai_creator_withdraws`
--

DROP TABLE IF EXISTS `geekai_creator_withdraws`;
CREATE TABLE `geekai_creator_withdraws` (
  `id` int NOT NULL COMMENT '主键ID',
  `creator_id` int NOT NULL COMMENT '创作者ID',
  `account` varchar(100) NOT NULL COMMENT '收款账号',
  `status` varchar(20) NOT NULL COMMENT '状态(pending/success/reject)',
  `note` varchar(255) DEFAULT NULL COMMENT '备注',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime NOT NULL COMMENT '更新时间',
  `qr_code` varchar(255) NOT NULL COMMENT '收款二维码',
  `method` varchar(20) NOT NULL COMMENT '收款方式(alipay/wxpay)',
  `scores` int NOT NULL COMMENT '提现积分',
  `fee` decimal(10,2) NOT NULL COMMENT '提现手续费',
  `total_money` decimal(10,2) NOT NULL COMMENT '提现总金额',
  `real_money` decimal(10,2) NOT NULL COMMENT '提现到账金额',
  `account_name` varchar(100) NOT NULL COMMENT '收款人姓名'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- --------------------------------------------------------

--
-- 表的结构 `geekai_files`
--

DROP TABLE IF EXISTS `geekai_files`;
CREATE TABLE `geekai_files` (
  `id` int NOT NULL,
  `user_id` bigint NOT NULL COMMENT '用户 ID',
  `name` varchar(100) NOT NULL COMMENT '文件名',
  `obj_key` varchar(100) DEFAULT NULL COMMENT '文件标识',
  `url` varchar(255) NOT NULL COMMENT '文件地址',
  `ext` varchar(10) NOT NULL COMMENT '文件后缀',
  `size` bigint NOT NULL DEFAULT '0' COMMENT '文件大小',
  `created_at` datetime NOT NULL COMMENT '创建时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户文件表';

-- --------------------------------------------------------

--
-- 表的结构 `geekai_orders`
--

DROP TABLE IF EXISTS `geekai_orders`;
CREATE TABLE `geekai_orders` (
  `id` int NOT NULL,
  `user_id` bigint NOT NULL COMMENT '用户ID',
  `pid` bigint NOT NULL COMMENT '产品ID',
  `username` varchar(30) NOT NULL COMMENT '用户名',
  `order_no` varchar(30) NOT NULL COMMENT '订单ID',
  `trade_no` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '支付平台交易流水号',
  `subject` varchar(100) NOT NULL COMMENT '订单产品',
  `amount` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '订单金额',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '订单状态（0：待支付，1：已扫码，2：支付成功）',
  `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '备注',
  `pay_time` bigint DEFAULT NULL COMMENT '支付时间',
  `pay_way` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '支付方式:alipay,wxpay',
  `channel` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '支付类型渠道：支付宝，微信，聚合支付',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime NOT NULL COMMENT '更新时间',
  `checked` tinyint NOT NULL DEFAULT '0' COMMENT '是否已检查'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='充值订单表';

-- --------------------------------------------------------

--
-- 表的结构 `geekai_products`
--

DROP TABLE IF EXISTS `geekai_products`;
CREATE TABLE `geekai_products` (
  `id` int NOT NULL,
  `name` varchar(100) NOT NULL COMMENT '产品名称',
  `price` decimal(10,2) NOT NULL COMMENT '产品价格',
  `credit` bigint NOT NULL COMMENT '积分额度',
  `enabled` tinyint(1) NOT NULL DEFAULT '1' COMMENT '启用状态',
  `sales` bigint NOT NULL DEFAULT '0' COMMENT '销量',
  `sort_num` tinyint NOT NULL DEFAULT '0' COMMENT '排序',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime NOT NULL COMMENT '更新时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='会员套餐表';

--
-- 转存表中的数据 `geekai_products`
--

INSERT INTO `geekai_products` (`id`, `name`, `price`, `credit`, `enabled`, `sales`, `sort_num`, `created_at`, `updated_at`) VALUES
(5, '100积分', 0.10, 100, 1, 0, 1, '2023-08-28 10:55:08', '2025-04-24 16:02:25'),
(6, '200积分', 19.90, 200, 1, 0, 2, '1970-01-01 08:00:00', '2025-04-24 08:38:40'),
(7, '300积分', 29.90, 300, 1, 0, 3, '2025-04-24 08:26:14', '2025-04-24 08:38:40'),
(8, '400积分', 29.90, 400, 1, 0, 4, '2025-04-24 08:35:27', '2025-04-24 08:38:40'),
(9, '500积分', 39.00, 500, 1, 0, 5, '2025-04-24 08:35:45', '2025-04-24 08:38:40');

-- --------------------------------------------------------

--
-- 表的结构 `geekai_redeems`
--

DROP TABLE IF EXISTS `geekai_redeems`;
CREATE TABLE `geekai_redeems` (
  `id` int NOT NULL,
  `user_id` bigint NOT NULL COMMENT '用户 ID',
  `name` varchar(30) NOT NULL COMMENT '兑换码名称',
  `amount` bigint NOT NULL COMMENT '额度',
  `code` varchar(100) NOT NULL COMMENT '兑换码',
  `enabled` tinyint(1) NOT NULL COMMENT '是否启用',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `redeemed_at` bigint NOT NULL COMMENT '兑换时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='兑换码';

-- --------------------------------------------------------

--
-- 表的结构 `geekai_score_logs`
--

DROP TABLE IF EXISTS `geekai_score_logs`;
CREATE TABLE `geekai_score_logs` (
  `id` int NOT NULL,
  `user_id` bigint NOT NULL COMMENT '用户ID',
  `username` varchar(30) NOT NULL COMMENT '用户名',
  `type` tinyint(1) NOT NULL COMMENT '类型（1：充值，2：消费，3：退费）',
  `amount` smallint NOT NULL COMMENT '算力数值',
  `balance` bigint NOT NULL COMMENT '余额',
  `subject` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '主题',
  `remark` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '备注',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `mark` tinyint(1) NOT NULL COMMENT '资金类型（0：支出，1：收入）'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户积分消费日志';

-- --------------------------------------------------------

--
-- 表的结构 `geekai_users`
--

DROP TABLE IF EXISTS `geekai_users`;
CREATE TABLE `geekai_users` (
  `id` int NOT NULL,
  `username` varchar(30) NOT NULL,
  `nickname` varchar(30) NOT NULL,
  `password` char(64) NOT NULL,
  `avatar` varchar(255) NOT NULL,
  `salt` char(12) NOT NULL,
  `scores` bigint DEFAULT '0',
  `expired_time` bigint NOT NULL,
  `last_login_at` bigint NOT NULL,
  `vip` tinyint(1) DEFAULT '0',
  `last_login_ip` char(16) NOT NULL,
  `openid` varchar(100) DEFAULT NULL,
  `platform` varchar(30) DEFAULT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `enabled` tinyint(1) NOT NULL DEFAULT '1',
  `invitor` bigint DEFAULT '0',
  `invite_code` varchar(100) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户表';

--
-- 转存表中的数据 `geekai_users`
--

INSERT INTO `geekai_users` (`id`, `username`, `nickname`, `password`, `avatar`, `salt`, `scores`, `expired_time`, `last_login_at`, `vip`, `last_login_ip`, `openid`, `platform`, `created_at`, `updated_at`, `enabled`, `invitor`, `invite_code`) VALUES
(4, '18888888888', '极客学长', 'c4a1188178c48afab68735ded98b71eb05b14f2990edf3bb98c4d94e5fd470b4', '/static/upload/2025/6/1749805607059551.png', 'nnubmo94', 6637, 0, 1755685086, 0, '127.0.0.1', '', '', '2024-11-22 12:12:53', '2025-08-20 18:18:06', 1, 0, NULL),
(54, 'abc@qq.com', '无敌的陆小凤@445124', '7fabf731cd818c8b5b84236e994bd7edc9007c961ba918be36b28905f65daba2', '/images/avatar/user.png', 't1z9n4u1', 123, 0, 0, 0, '', '', '', '2025-06-13 10:07:17', '2025-06-13 11:02:22', 1, 0, NULL),
(55, '18591927365', '冷酷的姜泥@646841', '0e814d937bb94051d5562be62c72860bcf983a79978ad8f436bc3897ca7297f1', '/images/avatar/user.png', '4olcqcku', 199, 0, 1753950536, 0, '127.0.0.1', '', '', '2025-06-13 10:58:50', '2025-07-31 16:28:56', 1, 0, NULL),
(57, 'yangjian@pvc123.com', '闪亮的徐凤年@200297', 'eb8d0fe5d92d0b08650a6872ec8bde5d4fad548c09bf7a7797f791699a00442c', '/images/avatar/user18.png', '2jnpi31n', 99, 0, 1752651962, 0, '127.0.0.1', '', '', '2025-06-14 20:37:52', '2025-07-16 15:46:02', 1, 0, NULL),
(58, 'user@16867438', '逍遥的小东邪@533486', 'b87ffd84df54f293fd695b2e7b97d596fbfc5c6254cd317354cd65da8cb6d5ac', '/images/avatar/user10.png', 'kkxn26om', 100, 0, 1749952584, 0, '127.0.0.1', 'oPyyL6iIjHa--j75ddSwjq2xKG_s', 'wechat', '2025-06-15 09:56:25', '2025-06-15 09:56:25', 1, 0, NULL),
(60, 'user@18992158', '呆萌的陆小凤@242059', '7d90f95bdb731297c1c493a2a09be3adb15c8edc2647c3491cbe07ec7facd9f3', '/images/avatar/user33.png', '6sb5osfe', 100, 1753718400, 1752652009, 0, '127.0.0.1', 'oPyyL6v9UKmvWsk7W9GfzVlIuZiY', 'wechat', '2025-06-15 09:59:26', '2025-07-30 19:50:39', 1, 0, NULL);

-- --------------------------------------------------------

--
-- 表的结构 `geekai_user_login_logs`
--

DROP TABLE IF EXISTS `geekai_user_login_logs`;
CREATE TABLE `geekai_user_login_logs` (
  `id` int NOT NULL,
  `user_id` bigint NOT NULL COMMENT '用户ID',
  `username` varchar(30) NOT NULL COMMENT '用户名',
  `login_ip` char(16) NOT NULL COMMENT '登录IP',
  `login_address` varchar(30) NOT NULL COMMENT '登录地址',
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户登录日志';

--
-- 转储表的索引
--

--
-- 表的索引 `geekai_admin_users`
--
ALTER TABLE `geekai_admin_users`
  ADD PRIMARY KEY (`id`) USING BTREE,
  ADD UNIQUE KEY `username` (`username`) USING BTREE,
  ADD UNIQUE KEY `username_2` (`username`),
  ADD UNIQUE KEY `username_3` (`username`);

--
-- 表的索引 `geekai_apps`
--
ALTER TABLE `geekai_apps`
  ADD PRIMARY KEY (`id`);

--
-- 表的索引 `geekai_app_categories`
--
ALTER TABLE `geekai_app_categories`
  ADD PRIMARY KEY (`id`);

--
-- 表的索引 `geekai_chat_items`
--
ALTER TABLE `geekai_chat_items`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `chat_id` (`chat_id`),
  ADD KEY `idx_geekai_chat_items_user_id` (`user_id`),
  ADD KEY `idx_geekai_chat_items_app_id` (`app_id`);

--
-- 表的索引 `geekai_chat_messages`
--
ALTER TABLE `geekai_chat_messages`
  ADD PRIMARY KEY (`id`),
  ADD KEY `chat_id` (`chat_id`),
  ADD KEY `idx_geekai_chat_messages_app_id` (`app_id`),
  ADD KEY `idx_geekai_chat_messages_chat_id` (`chat_id`),
  ADD KEY `idx_geekai_chat_messages_user_id` (`user_id`);

--
-- 表的索引 `geekai_configs`
--
ALTER TABLE `geekai_configs`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `name` (`name`);

--
-- 表的索引 `geekai_creators`
--
ALTER TABLE `geekai_creators`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `user_id` (`user_id`),
  ADD UNIQUE KEY `username` (`username`),
  ADD UNIQUE KEY `username_2` (`username`);

--
-- 表的索引 `geekai_creator_score_logs`
--
ALTER TABLE `geekai_creator_score_logs`
  ADD PRIMARY KEY (`id`),
  ADD KEY `idx_geekai_creator_score_logs_user_id` (`user_id`),
  ADD KEY `idx_geekai_creator_score_logs_creator_id` (`creator_id`);

--
-- 表的索引 `geekai_creator_withdraws`
--
ALTER TABLE `geekai_creator_withdraws`
  ADD PRIMARY KEY (`id`);

--
-- 表的索引 `geekai_files`
--
ALTER TABLE `geekai_files`
  ADD PRIMARY KEY (`id`),
  ADD KEY `idx_geekai_files_user_id` (`user_id`);

--
-- 表的索引 `geekai_orders`
--
ALTER TABLE `geekai_orders`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `order_no` (`order_no`);

--
-- 表的索引 `geekai_products`
--
ALTER TABLE `geekai_products`
  ADD PRIMARY KEY (`id`);

--
-- 表的索引 `geekai_redeems`
--
ALTER TABLE `geekai_redeems`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `code` (`code`),
  ADD KEY `idx_geekai_redeems_user_id` (`user_id`);

--
-- 表的索引 `geekai_score_logs`
--
ALTER TABLE `geekai_score_logs`
  ADD PRIMARY KEY (`id`),
  ADD KEY `idx_geekai_score_logs_user_id` (`user_id`);

--
-- 表的索引 `geekai_users`
--
ALTER TABLE `geekai_users`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `username` (`username`);

--
-- 表的索引 `geekai_user_login_logs`
--
ALTER TABLE `geekai_user_login_logs`
  ADD PRIMARY KEY (`id`),
  ADD KEY `idx_geekai_user_login_logs_user_id` (`user_id`);

--
-- 在导出的表使用AUTO_INCREMENT
--

--
-- 使用表AUTO_INCREMENT `geekai_admin_users`
--
ALTER TABLE `geekai_admin_users`
  MODIFY `id` int NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=137;

--
-- 使用表AUTO_INCREMENT `geekai_apps`
--
ALTER TABLE `geekai_apps`
  MODIFY `id` int NOT NULL AUTO_INCREMENT;

--
-- 使用表AUTO_INCREMENT `geekai_app_categories`
--
ALTER TABLE `geekai_app_categories`
  MODIFY `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=17;

--
-- 使用表AUTO_INCREMENT `geekai_chat_items`
--
ALTER TABLE `geekai_chat_items`
  MODIFY `id` int NOT NULL AUTO_INCREMENT;

--
-- 使用表AUTO_INCREMENT `geekai_chat_messages`
--
ALTER TABLE `geekai_chat_messages`
  MODIFY `id` bigint NOT NULL AUTO_INCREMENT;

--
-- 使用表AUTO_INCREMENT `geekai_configs`
--
ALTER TABLE `geekai_configs`
  MODIFY `id` int NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=15;

--
-- 使用表AUTO_INCREMENT `geekai_creators`
--
ALTER TABLE `geekai_creators`
  MODIFY `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID', AUTO_INCREMENT=6;

--
-- 使用表AUTO_INCREMENT `geekai_creator_score_logs`
--
ALTER TABLE `geekai_creator_score_logs`
  MODIFY `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- 使用表AUTO_INCREMENT `geekai_creator_withdraws`
--
ALTER TABLE `geekai_creator_withdraws`
  MODIFY `id` int NOT NULL AUTO_INCREMENT COMMENT '主键ID';

--
-- 使用表AUTO_INCREMENT `geekai_files`
--
ALTER TABLE `geekai_files`
  MODIFY `id` int NOT NULL AUTO_INCREMENT;

--
-- 使用表AUTO_INCREMENT `geekai_orders`
--
ALTER TABLE `geekai_orders`
  MODIFY `id` int NOT NULL AUTO_INCREMENT;

--
-- 使用表AUTO_INCREMENT `geekai_products`
--
ALTER TABLE `geekai_products`
  MODIFY `id` int NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=12;

--
-- 使用表AUTO_INCREMENT `geekai_redeems`
--
ALTER TABLE `geekai_redeems`
  MODIFY `id` int NOT NULL AUTO_INCREMENT;

--
-- 使用表AUTO_INCREMENT `geekai_score_logs`
--
ALTER TABLE `geekai_score_logs`
  MODIFY `id` int NOT NULL AUTO_INCREMENT;

--
-- 使用表AUTO_INCREMENT `geekai_users`
--
ALTER TABLE `geekai_users`
  MODIFY `id` int NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=64;

--
-- 使用表AUTO_INCREMENT `geekai_user_login_logs`
--
ALTER TABLE `geekai_user_login_logs`
  MODIFY `id` int NOT NULL AUTO_INCREMENT;
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
