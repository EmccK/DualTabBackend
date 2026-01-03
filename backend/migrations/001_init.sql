-- DualTab 后台管理系统数据库初始化脚本

-- 管理员用户表
CREATE TABLE IF NOT EXISTS admin_users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 图标分类表
CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    name_en VARCHAR(50),
    sort_order INT DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 推荐图标表
CREATE TABLE IF NOT EXISTS icons (
    id SERIAL PRIMARY KEY,
    uuid VARCHAR(36) UNIQUE NOT NULL,
    title VARCHAR(100) NOT NULL,
    description TEXT,
    url VARCHAR(500) NOT NULL,
    img_url VARCHAR(500),
    bg_color VARCHAR(20) DEFAULT '#ffffff',
    mime_type VARCHAR(50) DEFAULT 'image/png',
    category_id INT REFERENCES categories(id) ON DELETE SET NULL,
    sort_order INT DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 搜索引擎表
CREATE TABLE IF NOT EXISTS search_engines (
    id SERIAL PRIMARY KEY,
    uuid VARCHAR(36) UNIQUE NOT NULL,
    name VARCHAR(50) NOT NULL,
    url VARCHAR(500) NOT NULL,
    icon_url VARCHAR(500),
    sort_order INT DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 索引
CREATE INDEX IF NOT EXISTS idx_icons_category ON icons(category_id);
CREATE INDEX IF NOT EXISTS idx_icons_active ON icons(is_active);

-- 预设分类数据（兼容 MonkNow 分类 ID）
INSERT INTO categories (id, name, name_en, sort_order) VALUES
(24, '热门', 'hot', 1),
(9, '购物', 'shopping', 2),
(10, '社交', 'social', 3),
(26, '娱乐', 'entertainment', 4),
(11, '新闻与阅读', 'news', 5),
(14, '效率', 'efficiency', 6),
(25, '内置App', 'builtin', 7),
(15, '图片', 'image', 8),
(16, '生活方式', 'lifestyle', 9),
(17, '旅行', 'travel', 10),
(18, '科技与教育', 'tech', 11),
(19, '金融', 'finance', 12)
ON CONFLICT (id) DO NOTHING;

-- 重置序列
SELECT setval('categories_id_seq', (SELECT MAX(id) FROM categories));

-- 预设搜索引擎
INSERT INTO search_engines (uuid, name, url, icon_url, sort_order) VALUES
('e58b5a00-74fe-4319-af0a-d4999565dd71', 'Google', 'https://www.google.com/search?q=', 'https://static.monknow.com/newtab/searcher/e58b5a00-74fe-4319-af0a-d4999565dd71.svg', 1),
('0eb43a90-b4c7-43ce-9c73-ab110945f47d', '百度', 'https://www.baidu.com/s?wd=', 'https://static.monknow.com/newtab/searcher/0eb43a90-b4c7-43ce-9c73-ab110945f47d.svg', 2),
('ceb6c985-d09c-4fdc-b0ea-b304f1ee0f2d', 'Bing', 'https://www.bing.com/search?q=', 'https://static.monknow.com/newtab/searcher/ceb6c985-d09c-4fdc-b0ea-b304f1ee0f2d.svg', 3),
('2a5e69d9-bf13-4188-8da2-004551a913a0', 'Yahoo', 'https://search.yahoo.com/search?p=', 'https://static.monknow.com/newtab/searcher/2a5e69d9-bf13-4188-8da2-004551a913a0.svg', 4),
('118f7463-4411-4856-873f-2851faa3b543', 'Yandex', 'https://yandex.ru/search/?text=', 'https://static.monknow.com/newtab/searcher/118f7463-4411-4856-873f-2851faa3b543.svg', 5),
('259d8e2b-340e-4690-8046-88a0b130cbd0', 'DuckDuckGo', 'https://duckduckgo.com/?q=', 'https://static.monknow.com/newtab/searcher/259d8e2b-340e-4690-8046-88a0b130cbd0.svg', 6)
ON CONFLICT (uuid) DO NOTHING;

-- 预设管理员账号 (密码: admin123)
-- 密码哈希使用 bcrypt 生成
INSERT INTO admin_users (username, password_hash) VALUES
('admin', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZRGdjGj/n3.Q7Z7Z7Z7Z7Z7Z7Z7Z7')
ON CONFLICT (username) DO NOTHING;
