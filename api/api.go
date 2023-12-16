package api

import (
	"github.com/gin-gonic/gin"

	"market_system/api/handler"
	"market_system/config"
	"market_system/storage"
	"market_system/storage/redis"

	_ "market_system/api/docs"

	swaggerfiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

func SetUpApi(r *gin.Engine, cfg *config.Config, strg storage.StorageI, cache *redis.Cache) {

	handler := handler.NewHandler(cfg, strg, cache)

	r.Use(customCORSMiddleware())

	r.POST("/login", handler.Login)

	v1 := r.Group("/v1")
	v1.Use(handler.AuthMiddleware())

	// User ...
	v1.POST("/user", handler.CreateUser)
	v1.GET("/user/:id", handler.GetByIDUser)
	v1.GET("/user", handler.GetListUser)
	v1.PUT("/user/:id", handler.UpdateUser)
	v1.DELETE("/user/:id", handler.DeleteUser)

	// Category ...
	v1.POST("/category", handler.CreateCategory)
	v1.GET("/category/:id", handler.GetByIDCategory)
	v1.GET("/category", handler.GetListCategory)
	v1.PUT("/category/:id", handler.UpdateCategory)
	v1.DELETE("/category/:id", handler.DeleteCategory)

	//branch ...
	v1.POST("/branch", handler.Createbranch)
	v1.GET("/branch/:id", handler.GetByIDbranch)
	v1.GET("/branch", handler.GetListbranch)
	v1.PUT("/branch/:id", handler.Updatebranch)
	v1.DELETE("/branch/:id", handler.Deletebranch)

	//sale_point
	v1.POST("/sale_point", handler.CreateSalePoint)
	v1.GET("/sale_point/:id", handler.GetByIDSalePoint)
	v1.GET("/sale_point", handler.GetListSalePoint)
	v1.PUT("/sale_point/:id", handler.UpdateSalePoint)
	v1.DELETE("/sale_point/:id", handler.DeleteSalePoint)

	//employee
	v1.POST("/employee", handler.CreateEmployee)
	v1.GET("/employee/:id", handler.GetByIDEmployee)
	v1.GET("/employee", handler.GetListEmployee)
	v1.PUT("/employee/:id", handler.UpdateEmployee)
	v1.DELETE("/employee/:id", handler.DeleteEmployee)

	//supplier
	v1.POST("/supplier", handler.CreateSupplier)
	v1.GET("/supplier/:id", handler.GetByIDSupplier)
	v1.GET("/supplier", handler.GetListSupplier)
	v1.PUT("/supplier/:id", handler.UpdateSupplier)
	v1.DELETE("/supplier/:id", handler.DeleteSupplier)

	//product
	v1.POST("/product", handler.CreateProduct)
	v1.GET("/product/:id", handler.GetByIDProduct)
	v1.GET("/product", handler.GetListProduct)
	v1.PUT("/product/:id", handler.UpdateProduct)
	v1.DELETE("/product/:id", handler.DeleteProduct)

	//income
	v1.POST("/income", handler.CreateIncome)
	v1.GET("/income/:id", handler.GetByIDIncome)
	v1.GET("/income", handler.GetListProduct)
	v1.PUT("/income/:id", handler.UpdateIncome)
	v1.DELETE("/income/:id", handler.DeleteIncome)

	//income_product
	v1.POST("/income_product", handler.CreateIncomeProduct)
	v1.GET("/income_product/:id", handler.GetByIDIncomeProduct)
	v1.GET("/income_product", handler.GetListIncomeProduct)
	v1.PUT("/income_product/:id", handler.UpdateIncomeProduct)
	v1.DELETE("/income_product/:id", handler.DeleteIncomeProduct)

	//remainder
	v1.POST("/remainder", handler.CreateRemainder)
	v1.GET("/remainder/:id", handler.GetByIDRemainder)
	v1.GET("/remainder", handler.GetListRemainder)
	v1.PUT("/remainder/:id", handler.UpdateRemainder)
	v1.DELETE("/remainder/:id", handler.DeleteRemainder)

	//shift
	v1.POST("/shift", handler.CreateShift)
	v1.GET("/shift/:id", handler.GetByIDShift)
	v1.GET("/shift", handler.GetListShift)
	v1.PUT("/shift/:id", handler.UpdateShift)
	v1.DELETE("/shift/:id", handler.DeleteShift)

	//sale
	v1.POST("/sale", handler.CreateSale)
	v1.GET("/sale/:id", handler.GetByIDSale)
	v1.GET("/sale", handler.GetListSale)
	v1.PUT("/sale/:id", handler.UpdateSale)
	v1.DELETE("/sale/:id", handler.DeleteSale)

	v1.GET("/sale/scan-barcode/:sale_id", handler.SaleScanBarcode)
	v1.GET("/dosale/:sale_id", handler.Dosale)

	//sale_product
	v1.POST("/sale_products", handler.CreateSaleProduct)
	v1.GET("/sale_products/:id", handler.GetByIDSaleProduct)
	v1.GET("/sale_products", handler.GetListSaleProduct)
	v1.PUT("/sale_products/:id", handler.UpdateSaleProduct)
	v1.DELETE("/sale_products/:id", handler.DeleteSaleProduct)

	//payment
	v1.POST("/payment", handler.CreatePayment)
	v1.GET("/payment/:id", handler.GetByIDPayment)
	v1.GET("/payment", handler.GetListPayment)
	v1.PUT("/payment/:id", handler.UpdatePayment)
	v1.DELETE("/payment/:id", handler.DeletePayment)

	//transaction
	v1.POST("/transaction", handler.CreateTransaction)
	v1.GET("/transaction/:id", handler.GetByIDTransaction)
	v1.GET("/transaction", handler.GetListTransaction)
	v1.PUT("/transaction/:id", handler.UpdateTransaction)
	v1.DELETE("/transaction/:id", handler.DeleteTransaction)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}

func customCORSMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE, HEAD")
		c.Header("Access-Control-Allow-Headers", "Password, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Max-Age", "3600")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
