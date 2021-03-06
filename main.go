package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "gitlab.com/beehplus/sql-compose/docs"
	"gitlab.com/beehplus/sql-compose/restapi"
	"github.com/gin-contrib/cors"
	"os"
	"time"
)

//env
type Specification struct {
	Debug      bool
	Port       string
	BasePath   string
	Dsn        string
	User       string
	Rate       float32
	Timeout    time.Duration
	ColorCodes map[string]int
}

// @title sql-compose-api
// @version 1.0
// @description This is a api for sql-compose.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
func main() {
	var s Specification
	if err := envconfig.Process("sqlcompose", &s); err != nil {
		log.Fatal(err)
	}

	log.Info(s)
	log.SetLevel(log.DebugLevel)

	//init db
	db, err := sqlx.Connect("mysql", s.Dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	//b, _ := base64.StdEncoding.DecodeString("MjAyMDA1MjY3OQ==")
	//fmt.Println(string(b))

	router := gin.Default()

	url := ginSwagger.URL("http://localhost" +
		s.Port + "/swagger/doc.json")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	handler := restapi.NewHandler(db)

	// 跨域
	router.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "X-Requested-With", "X-CSRF-TOKEN"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowAllOrigins:  true,
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/doc", handler.GetDocList)
	router.POST("/doc/:uuid", handler.UpdateDoc)
	router.PATCH("/doc", handler.AddDoc)
	router.POST("/doc", handler.PostDoc)
	router.GET("/doc/:uuid", handler.GetDocDetailByUuid)
	router.DELETE("/doc/:uuid", handler.DeleteDoc)

	router.GET("/dns", handler.GetDbConfigList)
	router.DELETE("/dns/:uuid", handler.DeleteDbConfigByUUID)
	router.POST("/dns/:uuid", handler.UpdateDbConfigByUUID)
	router.POST("/dns", handler.AddDbConfig)

	router.POST(s.BasePath+"*path", handler.GetResult)

	if err := router.Run(s.Port); err != nil {
		log.Fatal(err)
	}
}

func init() {
	//log format json
	log.SetFormatter(&log.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		//PrettyPrint:     false,
	})
	log.SetOutput(os.Stdout)
	//set to true to show where the log is printed in the code
	log.SetReportCaller(true)
}

