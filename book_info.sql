-- 创建图书信息表
CREATE TABLE IF NOT EXISTS book_information (
    -- 主键和标识字段
    id DECIMAL PRIMARY KEY,
    book_id VARCHAR(50) NOT NULL UNIQUE,
    book_barcode VARCHAR(50) NOT NULL UNIQUE,
    
    -- 图书基本信息
    title VARCHAR(255) NOT NULL,
    publication_number VARCHAR(50),
    primary_author VARCHAR(100),
    classification_number VARCHAR(50),
    language_code VARCHAR(20),
    edition VARCHAR(50),
    
    -- 出版信息
    publisher VARCHAR(100),
    publication_place VARCHAR(100),
    publication_date DATE,
    distribution_unit VARCHAR(100),
    
    -- 其他信息
    notes TEXT,
    
    -- 创建时间戳
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 添加字段注释
COMMENT ON TABLE book_information IS 'Book information table';
COMMENT ON COLUMN book_information.id IS 'Unique identifier for the record (数据序列)';
COMMENT ON COLUMN book_information.book_id IS 'Book ID (图书编号)';
COMMENT ON COLUMN book_information.book_barcode IS 'Book barcode (图书条形码)';
COMMENT ON COLUMN book_information.title IS 'Book title (正标题)';
COMMENT ON COLUMN book_information.publication_number IS 'Publication number (图书出版号)';
COMMENT ON COLUMN book_information.primary_author IS 'Primary author (第一作者)';
COMMENT ON COLUMN book_information.classification_number IS 'Classification number (分类号)';
COMMENT ON COLUMN book_information.language_code IS 'Language code (语种码)';
COMMENT ON COLUMN book_information.edition IS 'Edition (版次)';
COMMENT ON COLUMN book_information.publisher IS 'Publisher (出版社)';
COMMENT ON COLUMN book_information.publication_place IS 'Publication place (出版地)';
COMMENT ON COLUMN book_information.publication_date IS 'Publication date (出版日期)';
COMMENT ON COLUMN book_information.distribution_unit IS 'Distribution unit (发行单位)';
COMMENT ON COLUMN book_information.notes IS 'Notes (备注)';

-- 创建索引
CREATE INDEX idx_book_info_book_id ON book_information(book_id);
CREATE INDEX idx_book_info_barcode ON book_information(book_barcode);
CREATE INDEX idx_book_info_title ON book_information(title);
CREATE INDEX idx_book_info_author ON book_information(primary_author);
CREATE INDEX idx_book_info_classification ON book_information(classification_number);

-- 创建更新时间戳触发器
CREATE OR REPLACE FUNCTION update_book_info_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_book_information_updated_at
    BEFORE UPDATE ON book_information
    FOR EACH ROW
    EXECUTE FUNCTION update_book_info_updated_at_column(); 