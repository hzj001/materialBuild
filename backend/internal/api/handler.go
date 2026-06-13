package api

import (
	"strconv"

	"material-build/internal/middleware"
	"material-build/internal/model"
	"material-build/internal/repository"
	"material-build/internal/service"
	"material-build/pkg/config"
	"material-build/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type Handler struct {
	cfg        *config.Config
	repo       *repository.Repository
	rdb        *redis.Client
	authSvc    *service.AuthService
	productSvc *service.ProductService
	orderSvc   *service.OrderService
}

// ========== 认证 ==========

type loginReq struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" binding:"required"`
}

func (h *Handler) Login(c *gin.Context) {
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误")
		return
	}
	result, err := h.authSvc.Login(req.Phone, req.Password, req.Role)
	if err != nil {
		response.Fail(c, 401, err.Error())
		return
	}
	response.OK(c, result)
}

type registerReq struct {
	Phone    string  `json:"phone" binding:"required"`
	Password string  `json:"password" binding:"required"`
	Nickname string  `json:"nickname"`
	Role     string  `json:"role" binding:"required"`
	CityID   *uint64 `json:"city_id"`
}

func (h *Handler) Register(c *gin.Context) {
	var req registerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误")
		return
	}
	result, err := h.authSvc.Register(service.RegisterInput{
		Phone: req.Phone, Password: req.Password,
		Nickname: req.Nickname, Role: req.Role, CityID: req.CityID,
	})
	if err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	response.OK(c, result)
}

// ========== 公共 ==========

func (h *Handler) ListCities(c *gin.Context) {
	cities, err := h.repo.ListCities()
	if err != nil {
		response.Fail(c, 500, "查询失败")
		return
	}
	response.OK(c, cities)
}

func (h *Handler) ListCategories(c *gin.Context) {
	cats, err := h.productSvc.Categories()
	if err != nil {
		response.Fail(c, 500, "查询失败")
		return
	}
	response.OK(c, cats)
}

// ========== 客户端 ==========

func (h *Handler) ClientSearchProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	cityID, _ := strconv.ParseUint(c.Query("city_id"), 10, 64)
	categoryID, _ := strconv.ParseUint(c.Query("category_id"), 10, 64)

	list, total, err := h.productSvc.Search(repository.ProductQuery{
		Keyword: c.Query("keyword"), CityID: cityID,
		CategoryID: categoryID, Page: page, PageSize: pageSize,
	})
	if err != nil {
		response.Fail(c, 500, "查询失败")
		return
	}
	response.OK(c, response.PageData(list, total, page, pageSize))
}

func (h *Handler) ClientProductDetail(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	p, err := h.productSvc.Detail(id)
	if err != nil {
		response.Fail(c, 404, "商品不存在")
		return
	}
	response.OK(c, p)
}

func (h *Handler) ClientListMerchants(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	cityID, _ := strconv.ParseUint(c.Query("city_id"), 10, 64)
	list, total, err := h.repo.ListMerchants(cityID, page, pageSize)
	if err != nil {
		response.Fail(c, 500, "查询失败")
		return
	}
	response.OK(c, response.PageData(list, total, page, pageSize))
}

type createOrderReq struct {
	MerchantID    uint64 `json:"merchant_id" binding:"required"`
	Items         []struct {
		ProductID uint64 `json:"product_id"`
		Quantity  int    `json:"quantity"`
	} `json:"items" binding:"required"`
	ReceiverName  string `json:"receiver_name" binding:"required"`
	ReceiverPhone string `json:"receiver_phone" binding:"required"`
	Address       string `json:"address" binding:"required"`
	Remark        string `json:"remark"`
}

func (h *Handler) ClientCreateOrder(c *gin.Context) {
	var req createOrderReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误")
		return
	}
	var items []service.OrderItemInput
	for _, it := range req.Items {
		items = append(items, service.OrderItemInput{ProductID: it.ProductID, Quantity: it.Quantity})
	}
	order, err := h.orderSvc.Create(service.CreateOrderInput{
		UserID: middleware.GetUserID(c), MerchantID: req.MerchantID,
		Items: items, ReceiverName: req.ReceiverName,
		ReceiverPhone: req.ReceiverPhone, Address: req.Address, Remark: req.Remark,
	})
	if err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	response.OK(c, order)
}

func (h *Handler) ClientListOrders(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	list, total, err := h.orderSvc.ListByUser(middleware.GetUserID(c), page, pageSize)
	if err != nil {
		response.Fail(c, 500, "查询失败")
		return
	}
	response.OK(c, response.PageData(list, total, page, pageSize))
}

func (h *Handler) ClientPayOrder(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.orderSvc.Pay(id); err != nil {
		response.Fail(c, 400, "支付失败")
		return
	}
	response.OKMsg(c, "支付成功")
}

type afterSaleReq struct {
	OrderID uint64 `json:"order_id" binding:"required"`
	Type    string `json:"type" binding:"required"`
	Reason  string `json:"reason" binding:"required"`
}

func (h *Handler) ClientCreateAfterSale(c *gin.Context) {
	var req afterSaleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误")
		return
	}
	order, err := h.repo.GetOrderByID(req.OrderID)
	if err != nil {
		response.Fail(c, 404, "订单不存在")
		return
	}
	a := &model.AfterSale{
		OrderID: req.OrderID, UserID: middleware.GetUserID(c),
		MerchantID: order.MerchantID, Type: req.Type, Reason: req.Reason,
	}
	if err := h.repo.CreateAfterSale(a); err != nil {
		response.Fail(c, 500, "提交失败")
		return
	}
	response.OK(c, a)
}

// ========== 商家端 ==========

func (h *Handler) MerchantProfile(c *gin.Context) {
	m, err := h.repo.FindMerchantByUserID(middleware.GetUserID(c))
	if err != nil {
		response.Fail(c, 404, "商家信息不存在")
		return
	}
	response.OK(c, m)
}

type createProductReq struct {
	CategoryID  uint64  `json:"category_id"`
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	CoverImage  string  `json:"cover_image"`
	Price       float64 `json:"price" binding:"required"`
	SalePrice   float64 `json:"sale_price"`
	Unit        string  `json:"unit"`
	Stock       int     `json:"stock"`
}

func (h *Handler) MerchantCreateProduct(c *gin.Context) {
	m, err := h.repo.FindMerchantByUserID(middleware.GetUserID(c))
	if err != nil {
		response.Fail(c, 404, "商家信息不存在")
		return
	}
	var req createProductReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误")
		return
	}
	p, err := h.productSvc.Create(service.CreateProductInput{
		MerchantID: m.ID, CategoryID: req.CategoryID, Name: req.Name,
		Description: req.Description, CoverImage: req.CoverImage,
		Price: req.Price, SalePrice: req.SalePrice, Unit: req.Unit, Stock: req.Stock,
	})
	if err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	response.OK(c, p)
}

func (h *Handler) MerchantListProducts(c *gin.Context) {
	m, err := h.repo.FindMerchantByUserID(middleware.GetUserID(c))
	if err != nil {
		response.Fail(c, 404, "商家信息不存在")
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	list, total, err := h.productSvc.Search(repository.ProductQuery{
		MerchantID: m.ID, Page: page, PageSize: pageSize,
	})
	if err != nil {
		response.Fail(c, 500, "查询失败")
		return
	}
	response.OK(c, response.PageData(list, total, page, pageSize))
}

func (h *Handler) MerchantUpdateProduct(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	p, err := h.repo.GetProductByID(id)
	if err != nil {
		response.Fail(c, 404, "商品不存在")
		return
	}
	var req createProductReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误")
		return
	}
	if req.Name != "" {
		p.Name = req.Name
	}
	if req.Price > 0 {
		p.Price = req.Price
	}
	if req.SalePrice > 0 {
		p.SalePrice = req.SalePrice
	}
	p.Stock = req.Stock
	if err := h.repo.UpdateProduct(p); err != nil {
		response.Fail(c, 500, "更新失败")
		return
	}
	response.OK(c, p)
}

func (h *Handler) MerchantListOrders(c *gin.Context) {
	m, err := h.repo.FindMerchantByUserID(middleware.GetUserID(c))
	if err != nil {
		response.Fail(c, 404, "商家信息不存在")
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	list, total, err := h.orderSvc.ListByMerchant(m.ID, page, pageSize)
	if err != nil {
		response.Fail(c, 500, "查询失败")
		return
	}
	response.OK(c, response.PageData(list, total, page, pageSize))
}

func (h *Handler) MerchantUpdateOrderStatus(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req struct {
		Status int8 `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误")
		return
	}
	if err := h.orderSvc.UpdateStatus(id, req.Status); err != nil {
		response.Fail(c, 500, "更新失败")
		return
	}
	response.OKMsg(c, "状态已更新")
}

func (h *Handler) MerchantListAfterSales(c *gin.Context) {
	m, err := h.repo.FindMerchantByUserID(middleware.GetUserID(c))
	if err != nil {
		response.Fail(c, 404, "商家信息不存在")
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	list, total, err := h.repo.ListAfterSalesByMerchant(m.ID, page, pageSize)
	if err != nil {
		response.Fail(c, 500, "查询失败")
		return
	}
	response.OK(c, response.PageData(list, total, page, pageSize))
}

func (h *Handler) MerchantHandleAfterSale(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req struct {
		Status int8   `json:"status"`
		Reply  string `json:"reply"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误")
		return
	}
	err := h.repo.DB().Model(&model.AfterSale{}).Where("id = ?", id).
		Updates(map[string]interface{}{"status": req.Status, "reply": req.Reply}).Error
	if err != nil {
		response.Fail(c, 500, "处理失败")
		return
	}
	response.OKMsg(c, "已处理")
}

// ========== 管理端 ==========

func (h *Handler) AdminListCities(c *gin.Context) {
	cities, err := h.repo.ListCities()
	if err != nil {
		response.Fail(c, 500, "查询失败")
		return
	}
	response.OK(c, cities)
}

func (h *Handler) AdminCreateCity(c *gin.Context) {
	var req struct {
		Name     string `json:"name" binding:"required"`
		Province string `json:"province"`
		Code     string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误")
		return
	}
	city := &model.City{Name: req.Name, Province: req.Province, Code: req.Code, Status: 1}
	if err := h.repo.CreateCity(city); err != nil {
		response.Fail(c, 500, "创建失败")
		return
	}
	response.OK(c, city)
}

func (h *Handler) AdminUpdateCityStatus(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req struct {
		Status int8 `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误")
		return
	}
	if err := h.repo.UpdateCityStatus(id, req.Status); err != nil {
		response.Fail(c, 500, "更新失败")
		return
	}
	response.OKMsg(c, "已更新")
}

func (h *Handler) AdminListMerchants(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	status := c.Query("status")
	var list []model.Merchant
	var total int64
	q := h.repo.DB().Model(&model.Merchant{})
	if status != "" {
		s, _ := strconv.Atoi(status)
		q = q.Where("status = ?", s)
	}
	q.Count(&total)
	err := q.Offset((page-1)*pageSize).Limit(pageSize).Find(&list).Error
	if err != nil {
		response.Fail(c, 500, "查询失败")
		return
	}
	response.OK(c, response.PageData(list, total, page, pageSize))
}

func (h *Handler) AdminUpdateMerchantStatus(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req struct {
		Status int8 `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误")
		return
	}
	if err := h.repo.UpdateMerchantStatus(id, req.Status); err != nil {
		response.Fail(c, 500, "更新失败")
		return
	}
	response.OKMsg(c, "已更新")
}

func (h *Handler) AdminListOrders(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	var list []model.Order
	var total int64
	h.repo.DB().Model(&model.Order{}).Count(&total)
	err := h.repo.DB().Preload("Items").Offset((page-1)*pageSize).Limit(pageSize).
		Order("created_at DESC").Find(&list).Error
	if err != nil {
		response.Fail(c, 500, "查询失败")
		return
	}
	response.OK(c, response.PageData(list, total, page, pageSize))
}

func (h *Handler) AdminStats(c *gin.Context) {
	var merchantCount, orderCount, userCount int64
	h.repo.DB().Model(&model.Merchant{}).Where("status = 1").Count(&merchantCount)
	h.repo.DB().Model(&model.Order{}).Count(&orderCount)
	h.repo.DB().Model(&model.User{}).Where("role = 'user'").Count(&userCount)
	response.OK(c, gin.H{
		"merchant_count": merchantCount,
		"order_count":    orderCount,
		"user_count":     userCount,
	})
}
