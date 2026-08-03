ALTER TABLE `geekai_apps` ADD `icon` VARCHAR(255) NULL DEFAULT NULL COMMENT '应用图标' AFTER `score`;
ALTER TABLE `geekai_apps` ADD `model_name` VARCHAR(100) NOT NULL COMMENT '模型名称' AFTER `icon`;
ALTER TABLE `geekai_chat_items` ADD `conversation_id` VARCHAR(100) DEFAULT NULL COMMENT '会话ID(coze/dify)' AFTER `chat_id`;
ALTER TABLE `geekai_apps` CHANGE `params` `configs` TEXT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '智能体配置参数';
ALTER TABLE `geekai_apps` DROP `token`;
ALTER TABLE `geekai_apps` DROP `api_url`;
ALTER TABLE `geekai_apps` DROP `bot_id`;
ALTER TABLE `geekai_apps` DROP `model_name`;

-- v1.0.2
ALTER TABLE `geekai_apps` ADD `bot_id` VARCHAR(30) NOT NULL COMMENT '机器人ID（coze 专用）' AFTER `type`;

-- 扣费模式功能
ALTER TABLE `geekai_apps` ADD `billing_mode` VARCHAR(20) NOT NULL DEFAULT 'immediate' COMMENT '扣费模式：immediate立即扣费，file_suffix文件后缀触发，string_marker字符串标记触发' AFTER `score`;
ALTER TABLE `geekai_apps` ADD `billing_config` TEXT NULL COMMENT '扣费配置JSON' AFTER `billing_mode`;