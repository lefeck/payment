package main

import (
	"github.com/jinzhu/gorm"
	"github.com/micro/go-micro/v2"
	service2 "payment/domain/service"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/wrapper/monitoring/prometheus/v2"
	"github.com/opentracing/opentracing-go"
	"github.com/wangjinh/common"
	consul2 "github.com/micro/go-plugins/registry/consul/v2"
	opentracing2 "github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	ratelimit "github.com/micro/go-plugins/wrapper/ratelimiter/uber/v2"
	"payment/domain/repository"
	"payment/handler"
	"payment/proto/payment"
	"strconv"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)
var QPS = 1000
func main() {
	//设置配置中心
	consulConf,err := common.GetConsulConfig("192.168.10.168",8500,"/micro/config")
	if err !=nil {
		log.Fatal(err)
	}
	//设置注册中心
	consul := consul2.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string {
			"192.168.10.168:8500",
		}
	})
	//设置链路追踪
	t,io,err := common.NewTracer("go.micro.service.payment","localhost:6434")
	if err !=nil {
		log.Fatal(err)
	}
	defer io.Close()

	opentracing.SetGlobalTracer(t)

	//连接数据库
	mysqlInfo :=common.GetMysqlFromConsul(consulConf,"mysql")

	mysqlInfoPort := strconv.FormatInt(mysqlInfo.Port,10)
	db,err := gorm.Open("mysql",mysqlInfo.User+":"+mysqlInfo.Pwd+"@("+mysqlInfo.Host+":"+mysqlInfoPort+")/"+mysqlInfo.Database+"?charset=utf8&parseTime=True&loc=Local")
	if err !=nil {
		log.Error(err)
	}
	defer db.Close()
	db.SingularTable(true)

	//初始化表
	if err :=repository.NewPaymentRepository(db).InitTable(); err !=nil {
		log.Fatal(err)
	}

	//监控
	common.PrometheusBoot(9082)

	// New Service
	service := micro.NewService(
		micro.Name("go.micro.service.payment"),
		micro.Version("latest"),
		micro.Address("0.0.0.0:8356"),
		//注册consul
		micro.Registry(consul),
		//添加链路追踪
		micro.WrapHandler(opentracing2.NewHandlerWrapper(opentracing.GlobalTracer())),
		//添加限流
		micro.WrapHandler(ratelimit.NewHandlerWrapper(QPS)),
		//添加prometheus监控
		micro.WrapHandler(prometheus.NewHandlerWrapper()),
	)

	// Initialise service
	service.Init()
	paymentDataService := service2.NewPaymentDataService(repository.NewPaymentRepository(db))
	// Register Handler
	payment.RegisterPaymentHandler(service.Server(), &handler.Payment{PaymentDataService: paymentDataService})

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
