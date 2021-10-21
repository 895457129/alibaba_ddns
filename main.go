package main

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"time"
)

type Config struct {
	DomainName string `json:"DomainName"`
	AccessKeyId string
	SecretKey string
	REGION   string
	RR string
	Type string
	UpdateDuration int
}

func main() {
	configFile, _ := ioutil.ReadFile("config.yml")
	c := Config{}
	err := yaml.Unmarshal(configFile, &c)
	if err != nil {
		panic(err)
	}
	updateDNS(c)
	for range time.Tick(time.Duration(c.UpdateDuration) * time.Minute) {
		updateDNS(c)
	}
}

func updateDNS(c Config)  {
	client, err := alidns.NewClientWithAccessKey("REGION_ID", c.AccessKeyId, c.SecretKey)
	if err != nil {
		panic(err)
	}
	q := alidns.CreateDescribeDomainRecordsRequest()
	q.DomainName = c.DomainName
	q.Type = "A"
	response, err := client.DescribeDomainRecords(q)
	if err != nil {
		panic(err)
	}
	if response.TotalCount >= 1 {
		ip := GetPublicIp()
		RecordId := response.DomainRecords.Record[0].RecordId
		changeRequest := alidns.CreateUpdateDomainRecordRequest()
		changeRequest.RecordId = RecordId
		changeRequest.RR = c.RR
		changeRequest.Value = ip
		changeRequest.Type = c.Type
		_, er := client.UpdateDomainRecord(changeRequest)
		if er != nil {
			fmt.Println("修改记录失败:", er)
		} else {
			fmt.Println("修改记录成功: ", ip)
		}
	} else {
		addRequest := alidns.CreateAddDomainRecordRequest()
		addRequest.DomainName = c.DomainName
		addRequest.RR = c.RR
		addRequest.Type = c.Type
		addRequest.Value = GetPublicIp()
		ip := GetPublicIp()
		_, er := client.AddDomainRecord(addRequest)
		if er != nil {
			fmt.Println("新增记录失败:", er)
		} else {
			fmt.Println("新增记录成功:", ip)
		}
	}
}

func GetPublicIp() string {
	responseClient, errClient := http.Get("http://ip.dhcp.cn/?ip") // 获取外网 IP
	if errClient != nil {
		fmt.Printf("获取外网 IP 失败，请检查网络\n")
		panic(errClient)
	}
	defer responseClient.Body.Close()
	body, _ := ioutil.ReadAll(responseClient.Body)
	return fmt.Sprintf("%s", string(body))
}

