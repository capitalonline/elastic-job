package models

import (
	"bytes"
	"encoding/json"
	"log"
	"testing"
)

func TestFmt(t *testing.T) {
	f := &FailureReport{
		Hostname:     "我",
		SubObject:    "你他",
		Ip:           "dasdas",
		Level:        "1",
		LogTimestamp: "倒是",
		Customername: "你",
		Tag1:         "他",
		Message:      "我",
	}
	obj := []FlumeObj{{Body: f.ToString()}}
	bf := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(bf)
	jsonEncoder.SetEscapeHTML(true)

	if err := jsonEncoder.Encode(obj); err != nil {
	}
	log.Printf("%s", bf.String())

}
