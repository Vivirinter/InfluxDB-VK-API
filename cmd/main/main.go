package main

import (
	"context"
	"encoding/json"
	"flag"
	influxdb2 "github.com/influxdata/influxdb-client-go"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Options = map[string][]string

type GroupResponse struct {
	Id      int
	Name    string
	Members int `json:"members_count"`
}

type VkResponse struct {
	Response []GroupResponse `json:"response"`
}

type VkConfig struct {
	Token   string
	GroupId string
	Version string
}

type GrafanaConfig struct {
	Url    string
	Org    string
	Token  string
	Bucket string
}

type Config struct {
	Vk      VkConfig
	Grafana GrafanaConfig
}

func readConfig() Config {
	config := Config{
		Vk:      VkConfig{},
		Grafana: GrafanaConfig{},
	}

	flag.StringVar(&config.Vk.Token, "token", "", "Enter token")
	flag.StringVar(&config.Vk.GroupId, "group_id", "", "Enter group id")
	flag.StringVar(&config.Grafana.Url, "grafana_url", "http://grafana/", "Grafana")
	flag.StringVar(&config.Grafana.Org, "grafana_org", "", "GrafanaOrg")
	flag.StringVar(&config.Grafana.Token, "grafana_token", "", "GrafanaToken")
	flag.StringVar(&config.Grafana.Bucket, "grafana_bucket", "", "GrafanaBucket")
	flag.StringVar(&config.Vk.Version, "vk_version", "5.199", "VK API version")

	flag.Parse()

	return config
}

func validateConfig(config Config) {
	if config.Vk.Token == "" || config.Vk.GroupId == "" || config.Grafana.Org == "" || config.Grafana.Token == "" || config.Grafana.Bucket == "" {
		log.Fatalln("Missing required parameters")
	}
}

func methodCall(method string, options Options, config VkConfig) (string, error) {
	var opts = Options{"access_token": {config.Token}, "v": {config.Version}}

	for k, v := range options {
		opts[k] = v
	}

	resp, err := http.PostForm("https://api.vk.com/method/"+method, opts)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func writeDataInfluxDb(groupId string, subscribers int, config GrafanaConfig) error {
	client := influxdb2.NewClient(config.Url, config.Token)
	defer client.Close()

	writeAPI := client.WriteAPIBlocking(config.Org, config.Bucket)

	p := influxdb2.NewPoint("VkGroups",
		map[string]string{"group_id": groupId},
		map[string]interface{}{"subscribers": subscribers},
		time.Now())

	err := writeAPI.WritePoint(context.Background(), p)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	config := readConfig()
	validateConfig(config)

	jsonResponse, err := methodCall("groups.getById", Options{"group_id": {config.Vk.GroupId}, "fields": {"members_count"}}, config.Vk)
	if err != nil {
		log.Fatalf("Error calling method: %s\n", err.Error())
	}

	var groupResponse VkResponse
	err = json.Unmarshal([]byte(jsonResponse), &groupResponse)
	if err != nil {
		log.Fatalf("Error unmarshalling response: %s\n", err.Error())
	}

	err = writeDataInfluxDb(config.Vk.GroupId, groupResponse.Response[0].Members, config.Grafana)
	if err != nil {
		log.Fatalf("Error writing data to InfluxDB: %s\n", err.Error())
	}
}
