package api

import (
	"fmt"
	_ "myproject/admin-api-gateway/api/docs"
	"myproject/admin-api-gateway/api/handlers/tokens"
	v1 "myproject/admin-api-gateway/api/handlers/v1"
	"myproject/admin-api-gateway/api/middleware"
	"myproject/admin-api-gateway/config"
	"myproject/admin-api-gateway/pkg/logger"
	"myproject/admin-api-gateway/services"
	"myproject/admin-api-gateway/storage/postgresrepo"
	"myproject/admin-api-gateway/storage/repo"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/util"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/gin-gonic/gin"
)

type Option struct {
	InMemory       repo.InMemoryStorageI
	Cfg            config.Config
	Logger         logger.Logger
	ServiceManager services.IServiceManager
	Postgres       postgresrepo.AdminStorageI
}

// Constructor
// @Title Clinic system
// @version 1.0
// @description api-gateway
// @host localhost:7070
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func New(option Option) *gin.Engine {
	psqlString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		option.Cfg.PostgresHost,
		option.Cfg.PostgresPort,
		option.Cfg.PostgresUser,
		option.Cfg.PostgresPassword,
		option.Cfg.PostgresDatabase)

	adapter, err := gormadapter.NewAdapter("postgres", psqlString, true)
	if err != nil {
		option.Logger.Fatal("error while creating a new adapter for casbin")
	}

	casbinEnforcer, err := casbin.NewEnforcer(option.Cfg.AuthConfigPath, adapter)
	if err != nil {
		option.Logger.Error("erro while creating a new casbin enforcer")
	}

	casbinEnforcer.GetRoleManager().AddMatchingFunc("keyMatch", util.KeyMatch)
	casbinEnforcer.GetRoleManager().AddMatchingFunc("keyMatch3", util.KeyMatch3)

	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	jwtHandler := tokens.JWTHandler{
		SignInKey: option.Cfg.SignInKey,
		Log:       option.Logger,
	}

	handlerV1 := v1.New(&v1.HandlerV1Config{
		InMemoryStorage: option.InMemory,
		Log:             option.Logger,
		ServiceManager:  option.ServiceManager,
		Cfg:             option.Cfg,
		JwtHandler:      jwtHandler,
		Casbin:          casbinEnforcer,
		Postgres:        option.Postgres,
	})

	api := router.Group("/v1")

	router.Static("/media", "./media")                    //unauthorized

	//Rbac
	api.GET("/rbac/roles", handlerV1.ListRoles)                 //superadmin
	api.GET("/rbac/policies/:role", handlerV1.ListRolePolicies) //superadmin
	api.POST("/rbac/add/policy", handlerV1.AddPolicyToRole)     //superadmin
	api.DELETE("/rbac/delete/policy", handlerV1.DeletePolicy)   //superadmin

	//Auth
	api.POST("/auth/create", handlerV1.CreateAdmin)            //superadmin
	api.DELETE("/auth/delete", handlerV1.DeleteAdmin)          //superadmin
	api.POST("/auth/login", handlerV1.LoginAdmin)              //unauthorized
	api.GET("/auth/admins/:page/:limit", handlerV1.ListAdmins) //admin
	api.GET("auth/get/:id", handlerV1.GetAdmin)                //admin
	api.PUT("auth/update", handlerV1.Update)                   //admin

	//User
	api.Use(middleware.Auth(casbinEnforcer, option.Cfg))
	api.POST("/register", handlerV1.Register)                   //unauthorized
	api.GET("/verify/:email/:code", handlerV1.Verify)           //unauthorized
	api.POST("/login", handlerV1.Login)                         //unauthorized
	api.POST("/user/create", handlerV1.CreateUser)              //admin
	api.GET("/user/:id", handlerV1.GetUserById)                 //user
	api.PUT("/user/update/:id", handlerV1.UpdateUser)           //user
	api.DELETE("/user/delete/:id", handlerV1.DeleteUser)        //user
	api.GET("/users/:page/:limit/:filter", handlerV1.ListUsers) //admin
	api.POST("/user/password", handlerV1.ChangePassword)        //user
	api.POST("/user/refresh", handlerV1.UpdateRefreshToken)     //user

	//Doctor
	api.POST("/doctor/register", handlerV1.RegisterDoctor)                               //unauthorized
	api.GET("/doctor/verify/{email}/{code}", handlerV1.VerifyDoctor)                     //unauthorized
	api.POST("/doctor/login", handlerV1.LoginDoctor)                                     //unauthorized
	api.POST("/doctor/create", handlerV1.CreateDoctor)                                   //admin, superadmin
	api.GET("/doctor/:id", handlerV1.GetDoctorById)                                      //doctor, user, operator, admin, superadmin
	api.PUT("/doctor/update/:id", handlerV1.UpdateDoctor)                                //doctor, admin, superadmin
	api.DELETE("/doctor/delete/:id", handlerV1.DeleteDoctor)                             //doctor, admin, superadmin
	api.GET("/doctors/:page/:limit", handlerV1.ListDoctors)                              //user, doctor, operator, admin, superadmin
	api.GET("/doctors/:page/:limit/:department_id", handlerV1.ListDoctorsByDepartmentId) //user, doctor, operator, admin, superadmin
	router.POST("/doctor/upload", handlerV1.UploadFile)                                     //unauthorized

	//Department
	api.POST("/department/create", handlerV1.CreateDepartment)       //admin, superadmin
	api.GET("/department/:id", handlerV1.GetDepartmentById)          //doctor, user, admin, operator, superadmin
	api.PUT("/department/update/:id", handlerV1.UpdateDepartment)    //admin, superadmin
	api.DELETE("/department/delete/:id", handlerV1.DeleteDepartment) //admin, superadmin
	api.GET("/departments/:page/:limit", handlerV1.ListDepartments)  //user, doctor, operator, admin, superadmin
	api.POST("/department/upload", handlerV1.UploadDepartmentFile)   //unauthorized

	//Specialization
	api.POST("/specialization/create", handlerV1.CreateSpecializaion)                                    //admin, superadmin
	api.GET("/specialization/:id", handlerV1.GetSpecializationById)                                      //user, doctor, operator, admin, superadmin
	api.PUT("/specialization/update/:id", handlerV1.UpdateSpecialization)                                //admin, superadmin
	api.DELETE("/specialization/delete/:id", handlerV1.DeleteSpecialization)                             //admin, superadmin
	api.GET("/specializations/:page/:limit", handlerV1.ListSpecializations)                              //user, doctor, operator, admin, superadmin
	api.GET("/specializations/:page/:limit/:department_id", handlerV1.ListSpecializationsByDepartmentId) //user, doctor, operator, admin, superadmin

	//Specialization price
	api.POST("/specprice/create", handlerV1.CreateSpecPrice)       //admin, superadmin
	api.GET("/specprice/:id", handlerV1.GetSpecPriceById)          //user, doctor, operator, admin, superadmin
	api.PUT("/specprice/update/:id", handlerV1.UpdateSpecPrice)    //admin, superadmin
	api.DELETE("/specprice/delete/:id", handlerV1.DeleteSpecPrice) //admin, superadmin
	api.GET("/specprices/:page/:limit", handlerV1.ListSpecPrices)  //user, doctor, operator, admin, superadmin

	url := ginSwagger.URL("swagger/doc.json")
	api.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	return router
}
