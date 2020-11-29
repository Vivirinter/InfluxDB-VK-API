package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	influxdb2 "github.com/influxdata/influxdb-client-go"
	"io/ioutil"
	"net/http"
	"time"
)

var groupId string
var vkToken string
var jsonResponse string
var groupResponse VkResponse
var grafanaUrl string
var grafanaToken string
var grafanaOrg string
var grafanaBucket string

type Options = map[string][]string

type GroupResponse struct {
	Id      int
	Name    string
	Members int `json:"members_count"`
}

type VkResponse struct {
	Response []GroupResponse `json:"response"`
}

func methodCall(method string, options map[string][]string) string {

	var opts = Options{"access_token": {vkToken}, "v": {"5.103"}}

	for k, v := range options {
		opts[k] = v
	}

	resp, err := http.PostForm("https://api.vk.com/method/"+method, opts)
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return string(body)
}

func writeDataInfluxDb(orgName string, bucketName string, groupId string, subscribers int, clientUrl string, token string) {
	client := influxdb2.NewClient(clientUrl, token)

	writeAPI := client.WriteAPIBlocking(orgName, bucketName)

	p := influxdb2.NewPoint("VkGroups",
		map[string]string{"group_id": groupId},
		map[string]interface{}{"subscribers": subscribers},
		time.Now())

	err := writeAPI.WritePoint(context.Background(), p)
	if err != nil {
		fmt.Printf("Write error: %s\n", err.Error())
	}

	client.Close()
}

func main() {

	//Line
	flag.StringVar(&vkToken, "token", "", "Enter token")

	flag.StringVar(&groupId, "group_id", "", "Enter group id")

	flag.StringVar(&grafanaUrl, "grafana_url", "http://grafana/", "Grafana")

	flag.StringVar(&grafanaOrg, "grafana_org", "", "GrafanaOrg")

	flag.StringVar(&grafanaToken, "grafana_token", "", "GrafanaToken")

	flag.StringVar(&grafanaBucket, "grafana_bucket", "", "GrafanaBucket")

	flag.Parse()

	jsonResponse = methodCall("groups.getById", Options{"group_id": {groupId}, "fields": {"members_count"}})

	json.Unmarshal([]byte(jsonResponse), &groupResponse)

	writeDataInfluxDb(grafanaOrg, grafanaBucket, groupId, groupResponse.Response[0].Members, grafanaUrl, grafanaToken)

}
