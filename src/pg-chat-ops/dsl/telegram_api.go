package dsl

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

const telegramBaseUrl = `https://api.telegram.org/`

type dslTelegram struct {
	tocken    string
	proxy     string
	ignoreSsl bool
	//internal
	client *http.Client
	offset int
}

func (d *dslTelegram) toString() string {
	name := strings.Split(d.tocken, ":")[0][3:]
	return fmt.Sprintf("telegram[`%sXXXX:XXXXX`]", name)
}

type dslTelegrams struct {
	sync.Mutex
	list map[string]*dslTelegram
}

var listOfTelegrams = &dslTelegrams{list: make(map[string]*dslTelegram, 0)}

// cache должен быть общий для всех плагинов
func newDSLTelegram(tocken string) *dslTelegram {
	listOfTelegrams.Lock()
	defer listOfTelegrams.Unlock()
	if result, ok := listOfTelegrams.list[tocken]; ok {
		return result
	}
	result := &dslTelegram{tocken: tocken}
	listOfTelegrams.list[tocken] = result
	return result
}

func (d *dslTelegram) buildHttpClient() error {

	client := &http.Client{
		Timeout: HTTP_TIMEOUT,
	}
	transport := &http.Transport{}
	if d.proxy != `` {
		proxyUrl, err := url.Parse(d.proxy)
		if err != nil {
			return err
		}
		transport.Proxy = http.ProxyURL(proxyUrl)
	}
	if d.ignoreSsl {
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
	client.Transport = transport

	d.client = client
	return nil
}
