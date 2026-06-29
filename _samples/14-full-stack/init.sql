# MySQL 初始化脚本 / MySQL init script
# docker-compose 启动时自动执行 / Auto-executed on docker-compose up

CREATE DATABASE IF NOT EXISTS stargo_demo CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE stargo_demo;

-- 用户表 / Users table (前缀 user_ 由 stargo 自动管理)
-- Table prefix "user_" is auto-managed by stargo store config
CREATE TABLE IF NOT EXISTS user_users (
    id       BIGINT AUTO_INCREMENT PRIMARY KEY,
    name     VARCHAR(128) NOT NULL DEFAULT '',
    email    VARCHAR(256) NOT NULL DEFAULT '',
    avatar   VARCHAR(512) NOT NULL DEFAULT '',
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    INDEX idx_name (name),
    INDEX idx_email (email)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 订单表 / Orders table (前缀 order_ 由 stargo 自动管理)
-- Table prefix "order_" is auto-managed by stargo store config
CREATE TABLE IF NOT EXISTS order_orders (
    id        BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_id   BIGINT NOT NULL DEFAULT 0,
    user_name VARCHAR(128) NOT NULL DEFAULT '',
    product   VARCHAR(256) NOT NULL DEFAULT '',
    amount    DECIMAL(10,2) NOT NULL DEFAULT 0.00,
    status    VARCHAR(32) NOT NULL DEFAULT 'pending',
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    INDEX idx_user_id (user_id),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 商品表 / Products table (前缀 product_ 由 stargo 管理)
CREATE TABLE IF NOT EXISTS product_products (
    id         BIGINT AUTO_INCREMENT PRIMARY KEY,
    name       VARCHAR(256) NOT NULL DEFAULT '',
    `desc`     VARCHAR(1024) NOT NULL DEFAULT '',
    price      DECIMAL(10,2) NOT NULL DEFAULT 0.00,
    stock      INT NOT NULL DEFAULT 0,
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    INDEX idx_name (name),
    INDEX idx_price (price)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- UID 生成器表 (用于分布式自增 ID)
-- UID generator table (for distributed auto-increment ID)
CREATE TABLE IF NOT EXISTS uid (
    business_id VARCHAR(64) NOT NULL,
    max_id      BIGINT NOT NULL DEFAULT 0,
    step        INT NOT NULL DEFAULT 100,
    PRIMARY KEY (business_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 插入测试数据 / Insert test data
INSERT IGNORE INTO user_users (id, name, email, avatar) VALUES
(1, 'Alice', 'alice@example.com', 'https://i.pravatar.cc/150?u=alice'),
(2, 'Bob',   'bob@example.com',   'https://i.pravatar.cc/150?u=bob'),
(3, 'Carol', 'carol@example.com', 'https://i.pravatar.cc/150?u=carol');

INSERT IGNORE INTO order_orders (id, user_id, user_name, product, amount, status) VALUES
(1, 1, 'Alice', 'MacBook Pro',    12999.00, 'done'),
(2, 1, 'Alice', 'AirPods Pro',     1999.00, 'shipped'),
(3, 2, 'Bob',   'iPhone 15',       6999.00, 'paid');

INSERT IGNORE INTO product_products (id, name, `desc`, price, stock) VALUES
(1, 'MacBook Pro',  'Apple M3 Pro chip, 18GB RAM, 512GB SSD',         12999.00, 50),
(2, 'AirPods Pro',  'Active Noise Cancellation, USB-C',                1999.00, 200),
(3, 'iPhone 15',   'A16 Bionic, 48MP camera, 128GB',                  6999.00, 100),
(4, 'iPad Air',    'M2 chip, 10.9-inch Liquid Retina, 128GB',         4799.00, 80),
(5, 'Apple Watch', 'Series 9, GPS, 45mm, Midnight Aluminium',         3199.00, 60);
