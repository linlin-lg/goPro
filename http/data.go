package main

import (
	"encoding/json"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"golang.org/x/net/context"
	"net/http"
	"strings"
	"time"
)

type SettingModel struct {
	ID string 	`json:"id"`
	Set1 string `json:"set1"`
	Set2 string	`json:"set2"`

	IsShow string 		`json:"isShow"`
	Image string 	`json:"image"`
	Name string  	`json:"name"`
}

type InfoModel struct {
	List []SettingModel 	`json:"list"`
	errorCode string

}

func GetData(request *http.Request) []byte {
	directory := request.URL.Path

	if request.Method == "GET" {
		parame := request.URL.Query()
		fmt.Println("\n request parame :", parame["Token"])
	}else if request.Method == "POST" {
		parame := request.PostForm["UserID"]
		fmt.Println("\n request parame :", parame)
	}

	if strings.Contains(directory, "/list") {
		list := moniModel()
		data,err := json.Marshal(list)
		if err != nil {
			panic(err)
		}
		return data

	} else if strings.Contains(directory, "/setting") {
		url := "https://applhb.longhuvip.com/w1/api/index.php?Token=7d5a372bc287f93b35716bc7652a1e1d&PhoneOSNew=2&VerSion=5.4.0.0&a=GetArrangeIndex&apiv=w29&UserID=799015&c=SysAppVersion"
		data := HttpGet(url)
		dataStr := string(data)
		fmt.Println("\njsonStr : ", dataStr)

		var model InfoModel
		err := json.Unmarshal([]byte(dataStr), &model)
		if err != nil {
			panic(err)
		}

		fmt.Println("\njsonModel ==", model)

		jsonData,err := json.Marshal(model)
		if err != nil {
			panic(err)
		}
		return jsonData
	}else if strings.Contains(directory, "/etcd/set") {
		EtcdSetKey("testKey", "1235", 0)
	}else if strings.Contains(directory, "/etcd/get") {
		EtcdGetKey("testKey")
	}

	return []byte(directory)
}

func moniModel() []SettingModel {

	modelList := make([]SettingModel,0)
	modelList = append(modelList, SettingModel{ID: "1", Set1: "name", Set2: "jerry"})
	modelList = append(modelList, SettingModel{ID: "2", Set1: "sli", Set2: "8989"})
	return modelList
}

var (
	etcdClient *clientv3.Client
	machines = []string{"http://127.0.0.1:2380"}
)

func connectEtcd() {

	cli, err := clientv3.New(clientv3.Config{
		Endpoints: machines,
		DialTimeout: 5 * time.Second,
	})

	if err != nil {
		fmt.Println("\nconnet to etcd failed")
		panic(err)
	}
	fmt.Println("\nconnet to etcd success")
	etcdClient = cli
	//defer etcdClient.Close()

	go EtcdWatchKey("testKey")

}

// 设置key值
func EtcdSetKey(key,vaule string, ttl uint64) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	_, err := etcdClient.Put(ctx, key, vaule)
	cancel()

	if err != nil {
		fmt.Println("\nput to etcd failed, err:%v", err)
		return
	}
	fmt.Println("\nput to etcd success")
}

// 获取key值
func EtcdGetKey(key string) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	resp, err := etcdClient.Get(ctx, key)
	cancel()
	if err != nil {
		fmt.Println("\nget from etcd failed, err:%v", err)
		return
	}
	for _,ev := range resp.Kvs {
		fmt.Printf("\n%s:%s", ev.Key, ev.Value)
	}
}

// 监听key值的变化
func EtcdWatchKey(key string) {

	etcdChan := etcdClient.Watch(context.Background(), key)
	for wresp := range etcdChan {
		for _, ev := range wresp.Events {
			fmt.Printf("Type: %s Key:%s Value:%s\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
		}
	}
}