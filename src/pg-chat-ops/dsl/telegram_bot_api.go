package dsl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/url"
	"strconv"
	"time"
)

func (d *dslTelegram) botRequest(method string, params url.Values) (*dslTelegramBotAPIResponse, error) {
	now := time.Now()
	url := fmt.Sprintf("%sbot%s/%s", telegramBaseUrl, d.tocken, method)
	response, err := d.client.PostForm(url, params)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	result := &dslTelegramBotAPIResponse{}
	var buf bytes.Buffer
	_, err = io.Copy(&buf, response.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %s", err.Error())
	}
	if err := json.Unmarshal(buf.Bytes(), result); err != nil {
		return nil, fmt.Errorf("parse response: %s", err.Error())
	}
	if !result.Ok {
		return nil, fmt.Errorf("response: %s", result.Description)
	}
	log.Printf("[INFO] processed %s bot request[%s]: %.2f\n", d.toString(), method, time.Now().Sub(now).Seconds())
	return result, nil
}

func (d *dslTelegram) getUpdates() ([]*dslTelegramBotUpdate, error) {
	params := url.Values{}
	params.Add(`limit`, strconv.Itoa(100))
	params.Add(`offset`, strconv.Itoa(d.offset))
	response, err := d.botRequest("getUpdates", params)
	if err != nil {
		return nil, err
	}
	result := make([]*dslTelegramBotUpdate, 0)
	if err := json.Unmarshal(response.Result, &result); err != nil {
		return nil, err
	}
	// update offset
	for _, upd := range result {
		if upd.UpdateID >= d.offset {
			d.offset = upd.UpdateID + 1
		}
	}
	return result, nil
}
