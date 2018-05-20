//注意：	
// daemon程序必须在go build后才可正常运行
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	metrics "github.com/rcrowley/go-metrics"
	rpcx "github.com/smallnest/rpcx/server"
	"github.com/smallnest/rpcx/serverplugin"
	"github.com/takama/daemon"
)

const (
	// 服务名称 随便定义
	name        = "myservice"
	description = "My rpcx Service"
)

var stdlog, errlog *log.Logger

var (
	addr       = flag.String("addr", "localhost:8972", "server addr")
	consulAddr = flag.String("consulAddr", "localhost:8500", "consul addr")
	basePath   = flag.String("basePath", "/repc_consul", "prefix path")
)

type Num struct {
	Numa int
	Numb int
}

type Reply struct {
	Sum int
	Str string
}

type Arith struct{}

func (a *Arith) Add(ctx context.Context, nums Num, reply *Reply) error {
	reply.Sum = nums.Numa + nums.Numb
	reply.Str = "hello rpcxConsul"
	return nil
}

// Service结构体已经包含了daemon
type Service struct {
	daemon.Daemon
}

// Manage用于控制和管理daemon
func (service *Service) Manage() (string, error) {

	usage := "Usage: myservice install | remove | start | stop | status"

	// 在命令行输入相应的命令时，执行对应的命令
	if len(os.Args) > 1 {
		command := os.Args[1]
		switch command {
		case "install":
			return service.Install()
		case "remove":
			return service.Remove()
		case "start":
			return service.Start()
		case "stop":
			return service.Stop()
		case "status":
			return service.Status()
		default:
			return usage, nil
		}
	}
	startsever()
	return usage, nil
}

// 注册rpcx服务端
func startsever() {
	flag.Parse()
	s := rpcx.NewServer()
	addRegistryPlugin(s)
	s.RegisterName("Arith", new(Arith), "")
	err := s.Serve("tcp", *addr)
	if err != nil {
		fmt.Println(err)
	}
}

//设置consul插件
func addRegistryPlugin(s *rpcx.Server) {
	r := &serverplugin.ConsulRegisterPlugin{
		ServiceAddress: "tcp@" + *addr,
		ConsulServers:  []string{*consulAddr},
		BasePath:       *basePath,
		Metrics:        metrics.NewRegistry(),
		UpdateInterval: time.Minute,
	}
	err := r.Start()
	if err != nil {
		log.Fatal(err)
	}
	s.Plugins.Add(r)
}

func init() {
	stdlog = log.New(os.Stdout, "", 0)
	errlog = log.New(os.Stderr, "", 0)
}

//main函数运行 基本不用改
func main() {
	srv, err := daemon.New(name, description)
	if err != nil {
		errlog.Println("Error: ", err)
		os.Exit(1)
	}
	service := &Service{srv}
	status, err := service.Manage()
	if err != nil {
		errlog.Println(status, "\nError: ", err)
		os.Exit(1)
	}
	fmt.Println(status)
}
