package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/smallnest/rpcx/client"
)

var (
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

func main() {
	flag.Parse()
	d := client.NewConsulDiscovery(*basePath, "Arith", []string{*consulAddr}, nil)
	xclient := client.NewXClient("Arith", client.Failtry, client.RandomSelect, d, client.DefaultOption)
	defer xclient.Close()
	nums := Num{
		Numa: 222,
		Numb: 333,
	}
	reply := new(Reply)
	err := xclient.Call(context.Background(), "Add", nums, reply)
	if err != nil {
		log.Printf("error:%v", err)
	}
	fmt.Printf("%d + %d = %d\n%s\n", nums.Numa, nums.Numb, reply.Sum, reply.Str)
}
