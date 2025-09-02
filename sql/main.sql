
create database world_quant;
use world_quant;
create table if not exists `idea` (
    `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'table:alpha_template ID',
    `idea_alpha_template` longtext NOT NULL COMMENT 'alpha 模板表达式',
    `idea_title` varchar(255) COMMENT 'alpha Idea title',
    `idea_desc` varchar(255) COMMENT 'alpha Idea desc',
    `start_idx` int(10) default -1 Comment 'alphaList-start-index',
    `end_idx` int(10) default -1 Comment 'alphaList-end-index',
    `next_idx` int(10) default -1 Comment 'next-index',
    `success_num` int(10) default 0 COMMENT '已经提交的,成功的数量',
    `fail_num` int(10) default 0 COMMENT '已经提交的,失败的数量',
    `is_finished` TINYINT unsigned DEFAULT 0 NOT NULL COMMENT '0表未完成，1表已完成',
    `concurrency_num`   int unsigned default 0 Comment 'max thread num for submit this idea',
    `created_at` DATETIME NULL DEFAULT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NULL DEFAULT NULL,
    `deleted_at` DATETIME NULL DEFAULT NULL,
    PRIMARY KEY (`id`)
    );

create table  if not exists `alpha`(
    `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'alpha id',
    `simulation_env` json NOT NULL COMMENT '模拟环境',
    `alpha` longtext NOT NULL COMMENT 'alpha表达式',
    `idea_id` int(10) unsigned COMMENT 'table:idea ID',
    `simulation_data` json NOT NULL COMMENT '模拟数据=模拟环境+alpha表达式',
    `test_period` varchar(50) COMMENT '模拟测试周期',
    `is_submitted` TINYINT unsigned DEFAULT 0 NOT NULL COMMENT '0表未提交，1表已提交，2表示提交失败',
    `is_deleted` TINYINT unsigned DEFAULT 0 NOT NULL COMMENT '0表未删除，1表已删除',
    `deleted_at` DATETIME NULL DEFAULT NULL,
    `updated_at` DATETIME NULL DEFAULT NULL,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
)ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COMMENT='alpha表';

CREATE INDEX  alpha_idea_id_idx ON alpha (idea_id);


create table  if not exists `alpha_result` (
    `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
    `idea_id` int(10) unsigned NOT NULL COMMENT 'table:idea ID',
    `alpha_id`  int(10) unsigned NOT NULL COMMENT '关联的alpha id',
    `alpha_detail` json NOT NULL COMMENT 'alpha表达式及其环境',
    `alpha_code` varchar(50) NOT NULL COMMENT 'alpha代码,brain平台唯一标识一个回测过的alpha',
    `basic_result` json  NOT NULL COMMENT '基本测试结果',
    `check_result` json  COMMENT '检查结果',
    `self_correlation` json COMMENT  '自相关性结果',
    `prod_correlation` json COMMENT  '生产相关性结果',
    `turnover` json  COMMENT  'turnover 详细数据',
    `sharpe`     json COMMENT  'sharpe 详细数据',
    `pnl`     json COMMENT  'pnl 详细数据',
    `daily_pnl`     json   COMMENT  'daily_pnl 详细数据',
    `yearly_stats`     json   COMMENT  'yearly_stats 详细数据',
    PRIMARY KEY (`id`),
    UNIQUE KEY `alpha_id` (`alpha_id`)
);
