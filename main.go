package main

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"regexp"
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
	ip, err := GetPublicIp()
	if len(ip) == 0 {
		panic("获取公网ip地址错误")
	}
	if response.TotalCount >= 1 {
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
		// addRequest.Value = GetPublicIp()
		_, er := client.AddDomainRecord(addRequest)
		if er != nil {
			fmt.Println("新增记录失败:", er)
		} else {
			fmt.Println("新增记录成功:", ip)
		}
	}
}

func GetPublicIp() (string, error) {
	var ip string
	var err error
	ip, err = GetPublicIp1()
	if err != nil {
		ip, err = GetPublicIp2()
		if err != nil {
			ip, err = GetPublicIp3()
			if err != nil {
				ip, err = GetPublicIp4()
			}
		}
	}
	return ip, err
}

func GetPublicIp1() (string, error) {
	responseClient, errClient := http.Get("http://www.net.cn/static/customercare/yourip.asp") // 获取外网 IP
	if errClient != nil {
		fmt.Printf("获取外网 IP 失败，请检查网络\n")
		return "", errClient
	}
	defer responseClient.Body.Close()
	content, err := ioutil.ReadAll(responseClient.Body)
	reg := regexp.MustCompile(`\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}`)
	// 利用正则表达式匹配规则对象匹配指定字符串
	res := reg.FindString(string(content))
	if len(res) == 0 {
		return "", err
	}
	return res, nil
}

func GetPublicIp2() (string, error) {
	resp, err := http.Get("http://myexternalip.com/raw")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	content, _ := ioutil.ReadAll(resp.Body)
	return string(content), nil
}

func GetPublicIp3() (string, error) {
	resp, err := http.Get("https://ip.tool.lu/")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	content, _ := ioutil.ReadAll(resp.Body)
	reg := regexp.MustCompile(`\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}`)
	// 利用正则表达式匹配规则对象匹配指定字符串
	res := reg.FindString(string(content))
	if len(res) == 0 {
		return "", err
	}
	return res, nil
}


func GetPublicIp4() (string, error) {
	resp, err := http.Get("http://www.cip.cc/")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	content, _ := ioutil.ReadAll(resp.Body)
	reg := regexp.MustCompile(`\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}`)
	// 利用正则表达式匹配规则对象匹配指定字符串
	res := reg.FindString(string(content))
	if len(res) == 0 {
		return "", err
	}
	return res, nil
}

