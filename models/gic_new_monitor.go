package models

import (
	"bytes"
	"encoding/json"
	"github.com/mongodb-job/config"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

type (
	FailureReport struct {
		//body = {
		//            "hostname": self.instance.get("by_user_name"),
		//            "subObject": "{0} ({1})".format(self.instance["svc_type"],self.instance["sub_product_name"]),
		//            "ip": "{}:{}".format(self.instance.get("ip"), self.instance.get("port")),
		//            "level": "Alert",
		//            "logTimestamp": self.start_time,
		//            "customername": self.customer_name,
		//            "tag1": self.event_type,
		//            "message": self._format_gic_monitor_message(),
		//        }

		Hostname     string `json:"hostname"`
		SubObject    string `json:"subObject"`
		Ip           string `json:"ip"`
		Level        string `json:"level"`
		LogTimestamp string `json:"logTimestamp"`
		Customername string `json:"customername"`
		Tag1         string `json:"tag1"`
		Message      string `json:"message"`
	}
	FlumeObj struct {
		Body string `json:"body"`
	}
)

func (report *FailureReport) SendFailureToGICNewMonitor() {
	client := &http.Client{}
	obj := []FlumeObj{{Body: report.ToString()}}
	j, err := json.Marshal(obj)
	if err != nil {
		logrus.Errorf("调用一线告警接口失败：%+v", err)
		return
	}
	req, err := http.NewRequest(http.MethodPost, config.Conf.Api.GicNewMonitorSystemUrl, bytes.NewReader(j))
	if err != nil {
		logrus.Errorf("调用一线告警接口失败：%+v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	resp, err := client.Do(req)
	if err != nil {
		logrus.Errorf("调用一线告警接口失败：%+v", err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		all, err1 := ioutil.ReadAll(resp.Body)
		if err1 != nil {
			logrus.Errorf("调用一线告警接口失败：%+v， %d", err1, resp.StatusCode)
			return
		}
		logrus.Errorf("调用一线告警接口失败：%+v", string(all))
		return
	}
	logrus.Infof("调用一线告警接口成功，参数：%s", string(j))
	return
}

func (report *FailureReport) ToString() string{

	bf := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(bf)
	//jsonEncoder.SetEscapeHTML(false)

	if err := jsonEncoder.Encode(report); err != nil {
		return ""
	}
	//marshal, err := json.Marshal(report)
	//if err != nil {
	//	return ""
	//}
	//ascii := strconv.QuoteToASCII(bf.String())
	return bf.String()
}
