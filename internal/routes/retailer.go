// package routes

// import (
// 	"github.com/levion-studio/paybazaar/internal/database"
// 	"github.com/levion-studio/paybazaar/internal/handlers"
// 	"github.com/levion-studio/paybazaar/internal/middlewares"
// 	"github.com/levion-studio/paybazaar/internal/repositories"
// 	"github.com/levion-studio/paybazaar/pkg"
// )

// func (r *routes) RetailerRoutes(db *database.Database, jwtUtils *pkg.JwtUtils) {
// 	retRepo := repositories.NewRetailerRepository(db, jwtUtils)
// 	retHandler := handlers.NewRetailerHandler(retRepo)

// 	r.Router.POST("/retailer/login", retHandler.LoginRetailerRequest)
// 	rrg := r.Router.Group("/retailer", middlewares.AuthorizationMiddleware(jwtUtils))
// 	rrg.POST("/create", retHandler.CreateRetailerRequest, middlewares.RequireRoles("admin", "master_distributor", "distributor"))
// 	rrg.PUT("/update/:retailer_id", retHandler.UpdateRetailerRequest, middlewares.RequireRoles("admin", "retailer"))
// 	rrg.DELETE("/delete/:retailer_id", retHandler.DeleteRetailerRequest, middlewares.RequireRoles("admin"))
// 	rrg.GET("/get/all", retHandler.ListRetailersRequest, middlewares.RequireRoles("admin"))
// 	rrg.GET("/get/retailer/:retailer_id", retHandler.GetRetailerByIDRequest, middlewares.RequireRoles("admin", "master_distributor", "distributor", "retailer"))
// 	rrg.GET("/get/md/:master_distributor_id", retHandler.ListRetailersByMasterDistributorIDRequest, middlewares.RequireRoles("admin", "master_distributor"))
// 	rrg.GET("/get/distributor/:distributor_id", retHandler.ListRetailersByDistributorIDRequest, middlewares.RequireRoles("admin", "master_distributor", "distributor"))
// 	rrg.GET("/get/dropdown/:distributor_id", retHandler.GetRetailersByDistributorIDForDropdownRequest, middlewares.RequireRoles("admin", "master_distributor", "distributor"))
// 	rrg.PUT("/update/kyc/:retailer_id", retHandler.UpdateKYCStatus, middlewares.RequireRoles("admin"))
// 	rrg.PUT("/update/block/:retailer_id", retHandler.UpdateBlockStatus, middlewares.RequireRoles("admin"))
// 	rrg.PUT("/update/mpin/:retailer_id", retHandler.UpdateMPIN, middlewares.RequireRoles("master_distributor"))
// }
package routes