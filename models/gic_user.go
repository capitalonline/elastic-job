package models

import (
	"encoding/json"
	"fmt"
	"github.com/mongodb-job/config"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

const (
	_customerInfo = "/getCustomerByNo/"
)

type (
	CustomerInfo struct {
		CustomerNo   string `json:"customerNo"`
		CustomerName string `json:"customerName"`
	}
	CustomerResponse struct {
		Status   string       `json:"status"`
		ErrCode  string       `json:"errCode"`
		ErrMsg   string       `json:"errMsg"`
		Customer CustomerInfo `json:"customer"`
	}
)

func GetCustomer(customerId string) CustomerInfo {
	var(
		c CustomerResponse
	)
	client := &http.Client{}
	uri := fmt.Sprintf("%s%s%s", config.Conf.Api.GicUserUrl, _customerInfo, customerId)
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		logrus.Errorf("查询用户信息失败(%v)", err.Error())
		return c.Customer
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)

	if err != nil {
		logrus.Errorf("查询用户信息失败(%v)", err.Error())
		return c.Customer
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("查询用户信息失败(%v)", err.Error())
		return c.Customer
	}
	if err = json.Unmarshal(body, &c); err != nil {
		logrus.Errorf("查询用户信息失败(%v)", err.Error())
	}
	return c.Customer
}
