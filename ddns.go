package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/RiversJin/DDNS4Aliyun/aliyun.dns.go"
)

func main() {
	var AccessKey = os.Getenv("AccessKey")
	var AccessKeySecret = os.Getenv("AccessKeySecret")
	var DomainName = os.Getenv("DomainName")
	var ResourceRecord = os.Getenv("ResourceRecord")
	if AccessKey == "" || AccessKeySecret == "" || DomainName == "" || ResourceRecord == "" {
		fmt.Println("please set AccessKey, AccessKeySecret, DomainName, ResourceRecord")
		return
	}
	fmt.Printf("AccessKey: %s \nAccessKeySecret: %s \nDomainName: %s \nResourceRecord: %s \n", AccessKey, AccessKeySecret, DomainName, ResourceRecord)
	ddnsCtx := &aliyun.DDNS{
		Ram: aliyun.RAM{
			AccessKeyId:     AccessKey,
			AccessKeySecret: AccessKeySecret,
		},
		RegionId: "cn-hangzhou",
	}
	var wait sync.WaitGroup
	wait.Add(1)
	ddnsCtx.AutoSetDomain(DomainName, ResourceRecord, "AAAA", "30m")
	wait.Wait()
}
