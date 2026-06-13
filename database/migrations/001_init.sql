-- 建材通 数据库初始化脚本
-- MySQL 8.0+

CREATE DATABASE IF NOT EXISTS material_build DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE material_build;

-- ========== 基础 ==========

CREATE TABLE cities (
    id          BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name        VARCHAR(64)  NOT NULL COMMENT '城市名称',
    province    VARCHAR(32)  NOT NULL DEFAULT '' COMMENT '省份',
    code        VARCHAR(16)  NOT NULL COMMENT '城市编码',
    status      TINYINT      NOT NULL DEFAULT 1 COMMENT '1启用 0禁用',
    created_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_code (code)
) ENGINE=InnoDB COMMENT='服务城市';

CREATE TABLE users (
    id          BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    phone       VARCHAR(20)  NOT NULL COMMENT '手机号',
    password    VARCHAR(128) NOT NULL DEFAULT '' COMMENT '密码哈希',
    nickname    VARCHAR(64)  NOT NULL DEFAULT '',
    avatar      VARCHAR(512) NOT NULL DEFAULT '',
    role        ENUM('user','merchant','admin') NOT NULL DEFAULT 'user',
    city_id     BIGINT UNSIGNED DEFAULT NULL COMMENT '所属城市',
    status      TINYINT      NOT NULL DEFAULT 1 COMMENT '1正常 0禁用',
    created_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_phone_role (phone, role),
    KEY idx_city (city_id)
) ENGINE=InnoDB COMMENT='用户账号';

-- ========== 商家 ==========

CREATE TABLE merchants (
    id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id         BIGINT UNSIGNED NOT NULL COMMENT '关联商家账号',
    city_id         BIGINT UNSIGNED NOT NULL,
    name            VARCHAR(128) NOT NULL COMMENT '店铺名称',
    logo            VARCHAR(512) NOT NULL DEFAULT '',
    description     TEXT,
    contact_phone   VARCHAR(20)  NOT NULL DEFAULT '' COMMENT '联系电话',
    address         VARCHAR(256) NOT NULL DEFAULT '',
    business_hours  VARCHAR(128) NOT NULL DEFAULT '08:00-20:00',
    delivery_radius DECIMAL(8,2) NOT NULL DEFAULT 10.00 COMMENT '配送半径(km)',
    min_order_amount DECIMAL(10,2) NOT NULL DEFAULT 0 COMMENT '起送金额',
    status          TINYINT NOT NULL DEFAULT 0 COMMENT '0待审核 1营业 2休息 3禁用',
    rating          DECIMAL(3,2) NOT NULL DEFAULT 5.00,
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    KEY idx_city (city_id),
    KEY idx_user (user_id),
    KEY idx_status (status)
) ENGINE=InnoDB COMMENT='建材商家';

-- ========== 商品 ==========

CREATE TABLE categories (
    id          BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    parent_id   BIGINT UNSIGNED NOT NULL DEFAULT 0,
    name        VARCHAR(64) NOT NULL,
    icon        VARCHAR(512) NOT NULL DEFAULT '',
    sort        INT NOT NULL DEFAULT 0,
    status      TINYINT NOT NULL DEFAULT 1,
    created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB COMMENT='商品分类';

CREATE TABLE products (
    id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    merchant_id     BIGINT UNSIGNED NOT NULL,
    category_id     BIGINT UNSIGNED NOT NULL DEFAULT 0,
    name            VARCHAR(256) NOT NULL,
    description     TEXT,
    cover_image     VARCHAR(512) NOT NULL DEFAULT '',
    price           DECIMAL(12,2) NOT NULL COMMENT '原价',
    sale_price      DECIMAL(12,2) NOT NULL COMMENT '售价',
    unit            VARCHAR(16) NOT NULL DEFAULT '件' COMMENT '单位',
    stock           INT NOT NULL DEFAULT 0,
    sales_count     INT NOT NULL DEFAULT 0,
    status          TINYINT NOT NULL DEFAULT 1 COMMENT '1上架 0下架',
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    KEY idx_merchant (merchant_id),
    KEY idx_category (category_id),
    FULLTEXT KEY ft_name (name)
) ENGINE=InnoDB COMMENT='商品';

CREATE TABLE product_media (
    id          BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    product_id  BIGINT UNSIGNED NOT NULL,
    type        ENUM('image','video') NOT NULL DEFAULT 'image',
    url         VARCHAR(512) NOT NULL,
    sort        INT NOT NULL DEFAULT 0,
    created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    KEY idx_product (product_id)
) ENGINE=InnoDB COMMENT='商品媒体';

CREATE TABLE promotions (
    id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    merchant_id     BIGINT UNSIGNED NOT NULL,
    product_id      BIGINT UNSIGNED DEFAULT NULL COMMENT 'NULL表示全店',
    title           VARCHAR(128) NOT NULL,
    discount_type   ENUM('percent','fixed') NOT NULL DEFAULT 'percent',
    discount_value  DECIMAL(10,2) NOT NULL,
    start_at        DATETIME NOT NULL,
    end_at          DATETIME NOT NULL,
    status          TINYINT NOT NULL DEFAULT 1,
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    KEY idx_merchant (merchant_id)
) ENGINE=InnoDB COMMENT='促销活动';

-- ========== 订单 ==========

CREATE TABLE orders (
    id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    order_no        VARCHAR(32) NOT NULL COMMENT '订单号',
    user_id         BIGINT UNSIGNED NOT NULL,
    merchant_id     BIGINT UNSIGNED NOT NULL,
    total_amount    DECIMAL(12,2) NOT NULL COMMENT '商品总额',
    discount_amount DECIMAL(12,2) NOT NULL DEFAULT 0,
    delivery_fee    DECIMAL(10,2) NOT NULL DEFAULT 0,
    pay_amount      DECIMAL(12,2) NOT NULL COMMENT '实付',
    status          TINYINT NOT NULL DEFAULT 0 COMMENT '0待支付 1已支付 2配送中 3已完成 4已取消 5售后中',
    pay_method      VARCHAR(32) NOT NULL DEFAULT '',
    pay_at          DATETIME DEFAULT NULL,
    remark          VARCHAR(512) NOT NULL DEFAULT '',
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_order_no (order_no),
    KEY idx_user (user_id),
    KEY idx_merchant (merchant_id),
    KEY idx_status (status)
) ENGINE=InnoDB COMMENT='订单';

CREATE TABLE order_items (
    id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    order_id        BIGINT UNSIGNED NOT NULL,
    product_id      BIGINT UNSIGNED NOT NULL,
    product_name    VARCHAR(256) NOT NULL,
    product_image   VARCHAR(512) NOT NULL DEFAULT '',
    price           DECIMAL(12,2) NOT NULL,
    quantity        INT NOT NULL DEFAULT 1,
    subtotal        DECIMAL(12,2) NOT NULL,
    KEY idx_order (order_id)
) ENGINE=InnoDB COMMENT='订单明细';

CREATE TABLE deliveries (
    id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    order_id        BIGINT UNSIGNED NOT NULL,
    receiver_name   VARCHAR(64) NOT NULL,
    receiver_phone  VARCHAR(20) NOT NULL,
    address         VARCHAR(512) NOT NULL,
    delivery_type   ENUM('merchant','platform') NOT NULL DEFAULT 'merchant',
    status          TINYINT NOT NULL DEFAULT 0 COMMENT '0待发货 1配送中 2已送达',
    driver_name     VARCHAR(64) NOT NULL DEFAULT '',
    driver_phone    VARCHAR(20) NOT NULL DEFAULT '',
    delivered_at    DATETIME DEFAULT NULL,
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_order (order_id)
) ENGINE=InnoDB COMMENT='配送信息';

-- ========== 售后 ==========

CREATE TABLE after_sales (
    id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    order_id        BIGINT UNSIGNED NOT NULL,
    user_id         BIGINT UNSIGNED NOT NULL,
    merchant_id     BIGINT UNSIGNED NOT NULL,
    type            ENUM('refund','return','exchange','complaint') NOT NULL,
    reason          VARCHAR(512) NOT NULL,
    images          JSON DEFAULT NULL COMMENT '凭证图片',
    status          TINYINT NOT NULL DEFAULT 0 COMMENT '0待处理 1处理中 2已完成 3已拒绝',
    reply           VARCHAR(512) NOT NULL DEFAULT '',
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    KEY idx_order (order_id),
    KEY idx_merchant (merchant_id)
) ENGINE=InnoDB COMMENT='售后工单';

-- ========== 客服 ==========

CREATE TABLE chat_sessions (
    id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    session_type    ENUM('official','merchant') NOT NULL DEFAULT 'official',
    user_id         BIGINT UNSIGNED NOT NULL,
    merchant_id     BIGINT UNSIGNED DEFAULT NULL,
    status          TINYINT NOT NULL DEFAULT 1 COMMENT '1进行中 0已关闭',
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    KEY idx_user (user_id)
) ENGINE=InnoDB COMMENT='客服会话';

CREATE TABLE chat_messages (
    id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    session_id      BIGINT UNSIGNED NOT NULL,
    sender_id       BIGINT UNSIGNED NOT NULL,
    sender_role     ENUM('user','merchant','admin') NOT NULL,
    content         TEXT NOT NULL,
    msg_type        ENUM('text','image') NOT NULL DEFAULT 'text',
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    KEY idx_session (session_id)
) ENGINE=InnoDB COMMENT='客服消息';

-- ========== 收款 ==========

CREATE TABLE merchant_settlements (
    id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    merchant_id     BIGINT UNSIGNED NOT NULL,
    order_id        BIGINT UNSIGNED NOT NULL,
    amount          DECIMAL(12,2) NOT NULL,
    platform_fee    DECIMAL(10,2) NOT NULL DEFAULT 0,
    settle_amount   DECIMAL(12,2) NOT NULL COMMENT '商家实收',
    status          TINYINT NOT NULL DEFAULT 0 COMMENT '0待结算 1已结算',
    settled_at      DATETIME DEFAULT NULL,
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    KEY idx_merchant (merchant_id)
) ENGINE=InnoDB COMMENT='商家结算';

-- ========== 初始数据 ==========

INSERT INTO cities (name, province, code) VALUES
('遵义', '贵州', '520300'),
('绵阳', '四川', '510700'),
('南阳', '河南', '411300');

INSERT INTO categories (parent_id, name, sort) VALUES
(0, '水泥砂石', 1),
(0, '瓷砖地板', 2),
(0, '油漆涂料', 3),
(0, '管材管件', 4),
(0, '五金工具', 5),
(0, '门窗定制', 6);

-- 默认管理员 账号: 13800000000  密码: admin123
INSERT INTO users (phone, password, nickname, role, status) VALUES
('13800000000', '$2a$10$OYSNoaEWDDcTCBpZ3B6nqeYuqhZJjeZRjx4eUI16KIqlQOqocJ5XO', '系统管理员', 'admin', 1);
