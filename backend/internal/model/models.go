package model

import "time"

// City 服务城市
type City struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:64" json:"name"`
	Province  string    `gorm:"size:32" json:"province"`
	Code      string    `gorm:"size:16" json:"code"`
	Status    int8      `gorm:"default:1" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (City) TableName() string { return "cities" }

// User 用户账号
type User struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	Phone     string    `gorm:"size:20" json:"phone"`
	Password  string    `gorm:"size:128" json:"-"`
	Nickname  string    `gorm:"size:64" json:"nickname"`
	Avatar    string    `gorm:"size:512" json:"avatar"`
	Role      string    `gorm:"size:16" json:"role"`
	CityID    *uint64   `json:"city_id"`
	Status    int8      `gorm:"default:1" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (User) TableName() string { return "users" }

// Merchant 建材商家
type Merchant struct {
	ID              uint64    `gorm:"primaryKey" json:"id"`
	UserID          uint64    `json:"user_id"`
	CityID          uint64    `json:"city_id"`
	Name            string    `gorm:"size:128" json:"name"`
	Logo            string    `gorm:"size:512" json:"logo"`
	Description     string    `gorm:"type:text" json:"description"`
	ContactPhone    string    `gorm:"size:20" json:"contact_phone"`
	Address         string    `gorm:"size:256" json:"address"`
	BusinessHours   string    `gorm:"size:128" json:"business_hours"`
	DeliveryRadius  float64   `gorm:"type:decimal(8,2)" json:"delivery_radius"`
	MinOrderAmount  float64   `gorm:"type:decimal(10,2)" json:"min_order_amount"`
	Status          int8      `json:"status"`
	Rating          float64   `gorm:"type:decimal(3,2)" json:"rating"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func (Merchant) TableName() string { return "merchants" }

// Category 商品分类
type Category struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	ParentID  uint64    `json:"parent_id"`
	Name      string    `gorm:"size:64" json:"name"`
	Icon      string    `gorm:"size:512" json:"icon"`
	Sort      int       `json:"sort"`
	Status    int8      `gorm:"default:1" json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

func (Category) TableName() string { return "categories" }

// Product 商品
type Product struct {
	ID          uint64    `gorm:"primaryKey" json:"id"`
	MerchantID  uint64    `json:"merchant_id"`
	CategoryID  uint64    `json:"category_id"`
	Name        string    `gorm:"size:256" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	CoverImage  string    `gorm:"size:512" json:"cover_image"`
	Price       float64   `gorm:"type:decimal(12,2)" json:"price"`
	SalePrice   float64   `gorm:"type:decimal(12,2)" json:"sale_price"`
	Unit        string    `gorm:"size:16" json:"unit"`
	Stock       int       `json:"stock"`
	SalesCount  int       `json:"sales_count"`
	Status      int8      `gorm:"default:1" json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	Merchant *Merchant      `gorm:"foreignKey:MerchantID" json:"merchant,omitempty"`
	Media    []ProductMedia `gorm:"foreignKey:ProductID" json:"media,omitempty"`
}

func (Product) TableName() string { return "products" }

// ProductMedia 商品媒体
type ProductMedia struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	ProductID uint64    `json:"product_id"`
	Type      string    `gorm:"size:16" json:"type"`
	URL       string    `gorm:"size:512" json:"url"`
	Sort      int       `json:"sort"`
	CreatedAt time.Time `json:"created_at"`
}

func (ProductMedia) TableName() string { return "product_media" }

// Order 订单
type Order struct {
	ID             uint64     `gorm:"primaryKey" json:"id"`
	OrderNo        string     `gorm:"size:32" json:"order_no"`
	UserID         uint64     `json:"user_id"`
	MerchantID     uint64     `json:"merchant_id"`
	TotalAmount    float64    `gorm:"type:decimal(12,2)" json:"total_amount"`
	DiscountAmount float64    `gorm:"type:decimal(12,2)" json:"discount_amount"`
	DeliveryFee    float64    `gorm:"type:decimal(10,2)" json:"delivery_fee"`
	PayAmount      float64    `gorm:"type:decimal(12,2)" json:"pay_amount"`
	Status         int8       `json:"status"`
	PayMethod      string     `gorm:"size:32" json:"pay_method"`
	PayAt          *time.Time `json:"pay_at"`
	Remark         string     `gorm:"size:512" json:"remark"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`

	Items    []OrderItem `gorm:"foreignKey:OrderID" json:"items,omitempty"`
	Delivery *Delivery   `gorm:"foreignKey:OrderID" json:"delivery,omitempty"`
}

func (Order) TableName() string { return "orders" }

// OrderItem 订单明细
type OrderItem struct {
	ID           uint64  `gorm:"primaryKey" json:"id"`
	OrderID      uint64  `json:"order_id"`
	ProductID    uint64  `json:"product_id"`
	ProductName  string  `gorm:"size:256" json:"product_name"`
	ProductImage string  `gorm:"size:512" json:"product_image"`
	Price        float64 `gorm:"type:decimal(12,2)" json:"price"`
	Quantity     int     `json:"quantity"`
	Subtotal     float64 `gorm:"type:decimal(12,2)" json:"subtotal"`
}

func (OrderItem) TableName() string { return "order_items" }

// Delivery 配送
type Delivery struct {
	ID            uint64     `gorm:"primaryKey" json:"id"`
	OrderID       uint64     `json:"order_id"`
	ReceiverName  string     `gorm:"size:64" json:"receiver_name"`
	ReceiverPhone string     `gorm:"size:20" json:"receiver_phone"`
	Address       string     `gorm:"size:512" json:"address"`
	DeliveryType  string     `gorm:"size:16" json:"delivery_type"`
	Status        int8       `json:"status"`
	DriverName    string     `gorm:"size:64" json:"driver_name"`
	DriverPhone   string     `gorm:"size:20" json:"driver_phone"`
	DeliveredAt   *time.Time `json:"delivered_at"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

func (Delivery) TableName() string { return "deliveries" }

// AfterSale 售后
type AfterSale struct {
	ID         uint64    `gorm:"primaryKey" json:"id"`
	OrderID    uint64    `json:"order_id"`
	UserID     uint64    `json:"user_id"`
	MerchantID uint64    `json:"merchant_id"`
	Type       string    `gorm:"size:16" json:"type"`
	Reason     string    `gorm:"size:512" json:"reason"`
	Images     string    `gorm:"type:json" json:"images"`
	Status     int8      `json:"status"`
	Reply      string    `gorm:"size:512" json:"reply"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (AfterSale) TableName() string { return "after_sales" }

// ChatSession 客服会话
type ChatSession struct {
	ID          uint64    `gorm:"primaryKey" json:"id"`
	SessionType string    `gorm:"size:16" json:"session_type"`
	UserID      uint64    `json:"user_id"`
	MerchantID  *uint64   `json:"merchant_id"`
	Status      int8      `gorm:"default:1" json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (ChatSession) TableName() string { return "chat_sessions" }

// ChatMessage 客服消息
type ChatMessage struct {
	ID         uint64    `gorm:"primaryKey" json:"id"`
	SessionID  uint64    `json:"session_id"`
	SenderID   uint64    `json:"sender_id"`
	SenderRole string    `gorm:"size:16" json:"sender_role"`
	Content    string    `gorm:"type:text" json:"content"`
	MsgType    string    `gorm:"size:16" json:"msg_type"`
	CreatedAt  time.Time `json:"created_at"`
}

func (ChatMessage) TableName() string { return "chat_messages" }

// Promotion 促销
type Promotion struct {
	ID            uint64    `gorm:"primaryKey" json:"id"`
	MerchantID    uint64    `json:"merchant_id"`
	ProductID     *uint64   `json:"product_id"`
	Title         string    `gorm:"size:128" json:"title"`
	DiscountType  string    `gorm:"size:16" json:"discount_type"`
	DiscountValue float64   `gorm:"type:decimal(10,2)" json:"discount_value"`
	StartAt       time.Time `json:"start_at"`
	EndAt         time.Time `json:"end_at"`
	Status        int8      `gorm:"default:1" json:"status"`
	CreatedAt     time.Time `json:"created_at"`
}

func (Promotion) TableName() string { return "promotions" }
