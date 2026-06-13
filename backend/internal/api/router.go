package api

import (
	"material-build/internal/middleware"
	"material-build/internal/repository"
	"material-build/internal/service"
	"material-build/pkg/config"
	redispkg "material-build/pkg/redis"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func SetupRouter(cfg *config.Config, repo *repository.Repository, rdb *redis.Client) *gin.Engine {
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	r.Use(middleware.CORS(cfg))

	authSvc := service.NewAuthService(repo, cfg)
	productSvc := service.NewProductService(repo)
	orderSvc := service.NewOrderService(repo)

	h := &Handler{
		cfg:        cfg,
		repo:       repo,
		rdb:        rdb,
		authSvc:    authSvc,
		productSvc: productSvc,
		orderSvc:   orderSvc,
	}

	r.GET("/health", func(c *gin.Context) {
		status := gin.H{"status": "ok", "service": "material-build", "mysql": "ok"}
		if rdb != nil {
			if err := redispkg.Ping(rdb); err != nil {
				status["redis"] = "error"
				status["status"] = "degraded"
			} else {
				status["redis"] = "ok"
			}
		}
		c.JSON(200, status)
	})

	v1 := r.Group("/api/v1")
	{
		// 公共
		common := v1.Group("/common")
		{
			common.GET("/cities", h.ListCities)
			common.GET("/categories", h.ListCategories)
		}

		// 认证
		auth := v1.Group("/auth")
		{
			auth.POST("/login", h.Login)
			auth.POST("/register", h.Register)
		}

		// 客户端
		client := v1.Group("/client")
		{
			client.GET("/products", h.ClientSearchProducts)
			client.GET("/products/:id", h.ClientProductDetail)
			client.GET("/merchants", h.ClientListMerchants)

			clientAuth := client.Group("")
			clientAuth.Use(middleware.Auth(cfg), middleware.RequireRole("user"))
			{
				clientAuth.POST("/orders", h.ClientCreateOrder)
				clientAuth.GET("/orders", h.ClientListOrders)
				clientAuth.POST("/orders/:id/pay", h.ClientPayOrder)
				clientAuth.POST("/after-sales", h.ClientCreateAfterSale)
			}
		}

		// 商家端
		merchant := v1.Group("/merchant")
		merchant.Use(middleware.Auth(cfg), middleware.RequireRole("merchant"))
		{
			merchant.GET("/profile", h.MerchantProfile)
			merchant.POST("/products", h.MerchantCreateProduct)
			merchant.GET("/products", h.MerchantListProducts)
			merchant.PUT("/products/:id", h.MerchantUpdateProduct)
			merchant.GET("/orders", h.MerchantListOrders)
			merchant.PUT("/orders/:id/status", h.MerchantUpdateOrderStatus)
			merchant.GET("/after-sales", h.MerchantListAfterSales)
			merchant.PUT("/after-sales/:id", h.MerchantHandleAfterSale)
		}

		// 管理端
		admin := v1.Group("/admin")
		admin.Use(middleware.Auth(cfg), middleware.RequireRole("admin"))
		{
			admin.GET("/cities", h.AdminListCities)
			admin.POST("/cities", h.AdminCreateCity)
			admin.PUT("/cities/:id/status", h.AdminUpdateCityStatus)
			admin.GET("/merchants", h.AdminListMerchants)
			admin.PUT("/merchants/:id/status", h.AdminUpdateMerchantStatus)
			admin.GET("/orders", h.AdminListOrders)
			admin.GET("/stats", h.AdminStats)
		}
	}

	return r
}
