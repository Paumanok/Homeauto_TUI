package client

import (
	"fmt"
	"net/http"
	"strings"
	"strconv"
	"encoding/json"
	"io/ioutil"
	//"reflect"
	"datapaddock.lan/ht_client/models"
)


//in this file we will have the various functions to call api endpoints

 
type Client struct {
	Url string
	Port string
}

func  getJsonArray[T interface{}](items []T, url string ) []T{

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer resp.Body.Close()
	
	err = json.NewDecoder(resp.Body).Decode(&items)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return items
}

func (c *Client) formatUrl(endpoint string) string {
	url := fmt.Sprintf("http://%s:%s%s", c.Url, c.Port, endpoint)
	return url
}

func (c *Client) GetDevices() []models.Device {
	url := c.formatUrl("/devices")
	var devs []models.Device
	devs = getJsonArray[models.Device](devs, url)
	return devs
}

func (c *Client) GetLast() []models.Measurement {
	url := c.formatUrl("/measurements/last")
	var meas []models.Measurement
	meas = getJsonArray[models.Measurement](meas, url)
	return meas
}

func (c *Client) GetSyncInt() int {
	ret := 0
	url := c.formatUrl("/next")
	resp, err:= http.Get(url)
	if err != nil {
		fmt.Println(err)
		return 60
	}
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return 60
	}
	body := string(respData)
	if strings.Contains(body, "sync_time"){
		parts := strings.Split(body, " ")
		if len(parts) > 1 {
			ret, err = strconv.Atoi(parts[1])
			if err != nil {
				fmt.Println(err)
				return 60
			}
		}

	}
	return ret
}
