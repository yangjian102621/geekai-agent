-- version 5.2.1
-- https://www.phpmyadmin.net/
--
-- 主机： 127.0.0.1
-- 生成日期： 2025-01-18 11:37:09
-- 服务器版本： 8.0.33
-- PHP 版本： 8.1.2-1ubuntu2.20

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
  `username` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '用户名',
  `password` char(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '密码',
  `salt` char(12) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '密码盐',
  `status` tinyint(1) NOT NULL COMMENT '当前状态',
  `last_login_at` int NOT NULL COMMENT '最后登录时间',
  `last_login_ip` char(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '最后登录 IP',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime NOT NULL COMMENT '更新时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='系统用户' ROW_FORMAT=DYNAMIC;

--
-- 转存表中的数据 `geekai_admin_users`
--

INSERT INTO `geekai_admin_users` (`id`, `username`, `password`, `salt`, `status`, `last_login_at`, `last_login_ip`, `created_at`, `updated_at`) VALUES
(1, 'admin', '6d17e80c87d209efb84ca4b2e0824f549d09fac8b2e1cc698de5bb5e1d75dfd0', 'mmrql75o', 1, 1737170930, '127.0.0.1', '2024-03-11 16:30:20', '2025-01-18 11:28:51');

-- --------------------------------------------------------

--
-- 表的结构 `geekai_apps`
--

DROP TABLE IF EXISTS `geekai_apps`;
CREATE TABLE `geekai_apps` (
  `id` int NOT NULL,
  `name` varchar(30) DEFAULT NULL COMMENT '名称',
  `type` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT 'openai' COMMENT 'openai,dify,coze',
  `enabled` tinyint(1) DEFAULT NULL COMMENT '是否启用',
  `configs` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '智能体配置参数',
  `score` int NOT NULL DEFAULT '0' COMMENT '单次对话消耗积分',
  `icon` varchar(255) DEFAULT NULL COMMENT '应用图标',
  `summary` varchar(255) DEFAULT NULL COMMENT '应用简介',
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='OpenAI API ';

-- --------------------------------------------------------

--
-- 表的结构 `geekai_chat_items`
--

DROP TABLE IF EXISTS `geekai_chat_items`;
CREATE TABLE `geekai_chat_items` (
  `id` int NOT NULL,
  `chat_id` char(40) NOT NULL COMMENT '会话 ID',
  `conversation_id` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '会话ID(coze/dify)',
  `user_id` int NOT NULL COMMENT '用户 ID',
  `app_id` int NOT NULL COMMENT '智能体ID',
  `title` varchar(100) NOT NULL COMMENT '会话标题',
  `icon` varchar(255) NOT NULL COMMENT '图标地址',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime NOT NULL COMMENT '更新时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户会话列表';

-- --------------------------------------------------------

--
-- 表的结构 `geekai_chat_messages`
--

DROP TABLE IF EXISTS `geekai_chat_messages`;
CREATE TABLE `geekai_chat_messages` (
  `id` bigint NOT NULL,
  `user_id` int NOT NULL COMMENT '用户 ID',
  `chat_id` char(40) NOT NULL COMMENT '会话 ID',
  `role` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT 'user or ai',
  `icon` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '角色图标',
  `app_id` int NOT NULL COMMENT '智能体ID',
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
  `marker` varchar(20) NOT NULL COMMENT '标识',
  `config_json` text NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- 转存表中的数据 `geekai_configs`
--

INSERT INTO `geekai_configs` (`id`, `marker`, `config_json`) VALUES
(1, 'system', '{\"title\":\"GeekAI 智能体\",\"slogan\":\"让每一个人都能用好AI\",\"admin_title\":\"GeekAgent 控制台\",\"logo\":\"/images/logo.png\",\"init_score\":123,\"daily_score\":456,\"invite_score\":789,\"enabled_register\":true,\"wechat_card_url\":\"/images/wx.png\",\"email_white_list\":[\"qq.com\",\"163.com\",\"gmail.com\",\"hotmail.com\",\"126.com\",\"outlook.com\",\"foxmail.com\",\"yahoo.com\",\"pvc123.com\"],\"app_id\":6,\"bot_avatar\":\"/images/avatar/assistant.png\",\"user_avatar\":\"/images/avatar/user.png\"}'),
(3, 'notice', '{\"sd_neg_prompt\":\"\",\"mj_mode\":\"\",\"index_bg_url\":\"\",\"index_navs\":null,\"copyright\":\"\",\"mark_map_text\":\"\",\"enabled_verify\":false,\"email_white_list\":null,\"content\":\"## v4.1.6 更新日志\\n\\n* 功能新增：**支持OpenAI实时语音对话功能** Beta 版，目前没有做算力计费控制，目前只有 VIP 用户可以使用。\\n* 功能优化：优化MysQL容器配置文档，解决MysQL容器资源占用过高问题\\n* 功能新增：管理后台增加AI绘图任务管理，可在管理后台浏览和删除用户的绘图任务\\n* 功能新增：管理后台增加Suno和Luma任务管理功能\\n* Bug修复：修复管理后台删除兑换码报 404 错误\\n* 功能优化：优化充值产品定价逻辑，可以设置原价和优惠价，**升级当前版本之后请务必要到管理后台去重新设置一下产品价格，以免造成损失！！！\\n\\n注意：当前站点仅为开源项目 \\u003ca style=\\\"color: #F56C6C\\\" href=\\\"https://github.com/yangjian102621/geekai\\\" target=\\\"_blank\\\"\\u003eGeekAI-Plus\\u003c/a\\u003e 的演示项目，本项目单纯就是给大家体验项目功能使用。\\n\\u003cstrong style=\\\"color: #F56C6C\\\"\\u003e体验额度用完之后请不要在当前站点进行任何充值操作！！！\\u003c/strong\\u003e\\n\\u003cstrong style=\\\"color: #F56C6C\\\"\\u003e体验额度用完之后请不要在当前站点进行任何充值操作！！！\\u003c/strong\\u003e\\n\\u003cstrong style=\\\"color: #F56C6C\\\"\\u003e体验额度用完之后请不要在当前站点进行任何充值操作！！！\\u003c/strong\\u003e\\n 如果觉得好用你就花几分钟自己部署一套，没有API KEY 的同学可以去下面几个推荐的中转站购买：\\n1、\\u003ca href=\\\"https://api.chat-plus.net\\\" target=\\\"_blank\\\"\\n   style=\\\"font-size: 20px;color:#F56C6C\\\"\\u003ehttps://api.chat-plus.net\\u003c/a\\u003e\\n2、\\u003ca href=\\\"https://api.geekai.me\\\" target=\\\"_blank\\\"\\n   style=\\\"font-size: 20px;color:#F56C6C\\\"\\u003ehttps://api.geekai.me\\u003c/a\\u003e\\n支持MidJourney，GPT，Claude，Google Gemmi，以及国内各个厂家的大模型，现在有超级优惠，价格远低于 OpenAI 官方。关于中转 API 的优势和劣势请参考 [中转API技术原理](https://docs.geekai.me/config/chat/#%E4%B8%AD%E8%BD%ACapi%E5%B7%A5%E4%BD%9C%E5%8E%9F%E7%90%86)。GPT-3.5，GPT-4，DALL-E3 绘图......你都可以随意使用，无需魔法。\\n接入教程： \\u003ca href=\\\"https://docs.geekai.me\\\" target=\\\"_blank\\\"\\n             style=\\\"font-size: 20px;color:#F56C6C\\\"\\u003ehttps://docs.geekai.me\\u003c/a\\u003e\\n本项目源码地址：\\u003ca href=\\\"https://github.com/yangjian102621/geekai\\\" target=\\\"_blank\\\"\\u003ehttps://github.com/yangjian102621/geekai\\u003c/a\\u003e\",\"updated\":true}');

-- --------------------------------------------------------

--
-- 表的结构 `geekai_files`
--

DROP TABLE IF EXISTS `geekai_files`;
CREATE TABLE `geekai_files` (
  `id` int NOT NULL,
  `user_id` int NOT NULL COMMENT '用户 ID',
  `name` varchar(100) NOT NULL COMMENT '文件名',
  `obj_key` varchar(100) DEFAULT NULL COMMENT '文件标识',
  `url` varchar(255) NOT NULL COMMENT '文件地址',
  `ext` varchar(10) NOT NULL COMMENT '文件后缀',
  `size` bigint NOT NULL DEFAULT '0' COMMENT '文件大小',
  `created_at` datetime NOT NULL COMMENT '创建时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户文件表';

-- --------------------------------------------------------

--
-- 表的结构 `geekai_products`
--

DROP TABLE IF EXISTS `geekai_products`;
CREATE TABLE `geekai_products` (
  `id` int NOT NULL,
  `name` varchar(30) NOT NULL COMMENT '名称',
  `price` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '价格',
  `discount` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '优惠金额',
  `days` smallint NOT NULL DEFAULT '0' COMMENT '延长天数',
  `power` int NOT NULL DEFAULT '0' COMMENT '增加算力值',
  `enabled` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否启动',
  `sales` int NOT NULL DEFAULT '0' COMMENT '销量',
  `sort_num` tinyint NOT NULL DEFAULT '0' COMMENT '排序',
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `app_url` varchar(255) DEFAULT NULL COMMENT 'App跳转地址',
  `url` varchar(255) DEFAULT NULL COMMENT '跳转地址'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='会员套餐表';

-- --------------------------------------------------------

--
-- 表的结构 `geekai_redeems`
--

DROP TABLE IF EXISTS `geekai_redeems`;
CREATE TABLE `geekai_redeems` (
  `id` int NOT NULL,
  `user_id` int NOT NULL COMMENT '用户 ID',
  `name` varchar(30) NOT NULL COMMENT '兑换码名称',
  `amount` int NOT NULL COMMENT '额度',
  `code` varchar(100) NOT NULL COMMENT '兑换码',
  `enabled` tinyint(1) NOT NULL COMMENT '是否启用',
  `created_at` datetime NOT NULL,
  `redeemed_at` int NOT NULL COMMENT '兑换时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='兑换码';

-- --------------------------------------------------------

--
-- 表的结构 `geekai_score_logs`
--

DROP TABLE IF EXISTS `geekai_score_logs`;
CREATE TABLE `geekai_score_logs` (
  `id` int NOT NULL,
  `user_id` int NOT NULL COMMENT '用户ID',
  `username` varchar(30) NOT NULL COMMENT '用户名',
  `type` tinyint(1) NOT NULL COMMENT '类型（1：充值，2：消费，3：退费）',
  `amount` smallint NOT NULL COMMENT '算力数值',
  `balance` int NOT NULL COMMENT '余额',
  `subject` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '主题',
  `remark` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '备注',
  `mark` tinyint(1) NOT NULL COMMENT '资金类型（0：支出，1：收入）',
  `created_at` datetime NOT NULL COMMENT '创建时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户积分消费日志';

-- --------------------------------------------------------

--
-- 表的结构 `geekai_users`
--

DROP TABLE IF EXISTS `geekai_users`;
CREATE TABLE `geekai_users` (
  `id` int NOT NULL,
  `username` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '用户名',
  `mobile` char(11) DEFAULT NULL COMMENT '手机号',
  `email` varchar(50) DEFAULT NULL COMMENT '邮箱地址',
  `nickname` varchar(30) NOT NULL COMMENT '昵称',
  `password` char(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '密码',
  `avatar` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '头像',
  `salt` char(12) NOT NULL COMMENT '密码盐',
  `scores` int NOT NULL DEFAULT '0' COMMENT '剩余积分',
  `expired_time` int NOT NULL COMMENT '用户过期时间',
  `status` tinyint(1) NOT NULL COMMENT '当前状态',
  `last_login_at` int NOT NULL COMMENT '最后登录时间',
  `vip` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否会员',
  `last_login_ip` char(16) NOT NULL COMMENT '最后登录 IP',
  `openid` varchar(100) DEFAULT NULL COMMENT '第三方登录账号ID',
  `platform` varchar(30) DEFAULT NULL COMMENT '登录平台',
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户表';

--
-- 转存表中的数据 `geekai_users`
--

INSERT INTO `geekai_users` (`id`, `username`, `mobile`, `email`, `nickname`, `password`, `avatar`, `salt`, `scores`, `expired_time`, `status`, `last_login_at`, `vip`, `last_login_ip`, `openid`, `platform`, `created_at`, `updated_at`) VALUES
(4, '18888888888', '18888888888', 'abc123@qq.com', '极客学长', 'c4a1188178c48afab68735ded98b71eb05b14f2990edf3bb98c4d94e5fd470b4', '/images/avatar/user.png', 'nnubmo94', 3570, 1734019200, 1, 1737162657, 0, '127.0.0.1', '', '', '2024-11-22 12:12:53', '2025-01-18 09:10:58'),
(50, 'geekmaster', '18575670123', 'yangjian@123.com', '用户@199191', 'e3f5417a818899298f3e9dbba68850a2027c0d3f15f881553921879c12b1e8c5', '/images/avatar/user.png', '6w691hmg', 1234, 1734019200, 1, 0, 0, '', '', '', '2024-12-05 09:04:17', '2024-12-05 09:04:17'),
(51, 'yangjian@pvc123.com', '', 'yangjian@pvc123.com', '用户@469083', '632d359bf3b1a29360c87292d708a215b722447f444485198615f8a5ebef695d', '/images/avatar/user.png', 'j6zzv3gq', 123, 0, 1, 0, 0, '', '', '', '2024-12-05 12:00:40', '2024-12-05 12:00:40'),
(52, 'yangjian102621@gmail.com', '', 'yangjian102621@gmail.com', '用户@550657', '246545e5a2feefe1ef205d7f0a5324ab43d209cbdc8541a04b2bb38aa2d4760d', '/images/avatar/user.png', '8njlj81n', 123, 0, 1, 0, 0, '', '', '', '2024-12-05 14:12:51', '2024-12-05 14:12:51');

-- --------------------------------------------------------

--
-- 表的结构 `geekai_user_login_logs`
--

DROP TABLE IF EXISTS `geekai_user_login_logs`;
CREATE TABLE `geekai_user_login_logs` (
  `id` int NOT NULL,
  `user_id` int NOT NULL COMMENT '用户ID',
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
  ADD UNIQUE KEY `username` (`username`) USING BTREE;

--
-- 表的索引 `geekai_apps`
--
ALTER TABLE `geekai_apps`
  ADD PRIMARY KEY (`id`);

--
-- 表的索引 `geekai_chat_items`
--
ALTER TABLE `geekai_chat_items`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `chat_id` (`chat_id`);

--
-- 表的索引 `geekai_chat_messages`
--
ALTER TABLE `geekai_chat_messages`
  ADD PRIMARY KEY (`id`),
  ADD KEY `chat_id` (`chat_id`);

--
-- 表的索引 `geekai_configs`
--
ALTER TABLE `geekai_configs`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `marker` (`marker`);

--
-- 表的索引 `geekai_files`
--
ALTER TABLE `geekai_files`
  ADD PRIMARY KEY (`id`);

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
  ADD UNIQUE KEY `code` (`code`);

--
-- 表的索引 `geekai_score_logs`
--
ALTER TABLE `geekai_score_logs`
  ADD PRIMARY KEY (`id`);

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
  ADD PRIMARY KEY (`id`);

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
  MODIFY `id` int NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=4;

--
-- 使用表AUTO_INCREMENT `geekai_files`
--
ALTER TABLE `geekai_files`
  MODIFY `id` int NOT NULL AUTO_INCREMENT;

--
-- 使用表AUTO_INCREMENT `geekai_products`
--
ALTER TABLE `geekai_products`
  MODIFY `id` int NOT NULL AUTO_INCREMENT;

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
  MODIFY `id` int NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=53;

--
-- 使用表AUTO_INCREMENT `geekai_user_login_logs`
--
ALTER TABLE `geekai_user_login_logs`
  MODIFY `id` int NOT NULL AUTO_INCREMENT;
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
