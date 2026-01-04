-- 添加 description 字段到 favicon_cache 表
ALTER TABLE favicon_cache ADD COLUMN description TEXT COMMENT '网站描述';
