-- ブログ用のテーブル

CREATE TABLE IF NOT EXISTS articles (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    content TEXT NOT NULL,
    image_url VARCHAR(2083),
    is_publish BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_articles_slug ON articles(slug);


CREATE TABLE IF NOT EXISTS tags (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS articles_tags (
    article_id INT NOT NULL,
    tag_id INT NOT NULL,
    PRIMARY KEY (article_id, tag_id),
    FOREIGN KEY (article_id) REFERENCES articles (id) ON DELETE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES tags (id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 更新時に `updated_at` を自動更新する関数
DROP TRIGGER IF EXISTS trigger_update_articles ON articles;
DROP TRIGGER IF EXISTS trigger_update_about ON about;

DROP FUNCTION IF EXISTS update_modified_column;

CREATE FUNCTION update_modified_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- `articles` テーブルの更新時に `update_modified_column` を実行するトリガー

CREATE TRIGGER trigger_update_articles
BEFORE UPDATE ON articles
FOR EACH ROW
EXECUTE FUNCTION update_modified_column();

-- 自己紹介用のテーブル

CREATE TABLE IF NOT EXISTS about (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    roma VARCHAR(255),
    description TEXT,
    color VARCHAR(255) NOT NULL,
    image_url VARCHAR(2083),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

DROP TYPE IF EXISTS skill_type;

CREATE TYPE skill_type AS ENUM ('Language', 'FrameWorks', 'Tools');

CREATE TABLE IF NOT EXISTS skills (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(255) NOT NULL,
    has_image BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS contact (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    link VARCHAR(2083) NOT NULL,
    has_image BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER trigger_update_about
BEFORE UPDATE ON about
FOR EACH ROW
EXECUTE FUNCTION update_modified_column();

-- プロダクト紹介用のテーブル

CREATE TABLE IF NOT EXISTS product (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    image_url VARCHAR(2083),
    github VARCHAR(2083),
    blog VARCHAR(2083),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
