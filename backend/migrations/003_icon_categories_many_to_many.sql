-- 创建图标-分类多对多关联表
CREATE TABLE IF NOT EXISTS icon_categories (
    icon_id INT NOT NULL REFERENCES icons(id) ON DELETE CASCADE,
    category_id INT NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (icon_id, category_id)
);

-- 索引
CREATE INDEX IF NOT EXISTS idx_icon_categories_icon ON icon_categories(icon_id);
CREATE INDEX IF NOT EXISTS idx_icon_categories_category ON icon_categories(category_id);

-- 迁移现有数据：将 icons 表中的 category_id 迁移到中间表
INSERT INTO icon_categories (icon_id, category_id)
SELECT id, category_id FROM icons WHERE category_id IS NOT NULL AND category_id > 0
ON CONFLICT DO NOTHING;
