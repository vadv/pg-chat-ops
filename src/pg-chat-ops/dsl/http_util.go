package dsl

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	lua "github.com/yuin/gopher-lua"
)

/*
variables:
	* HTTP_PROXY
	* SSL_INSECURE
*/

const HTTP_USER_AGENT = `zbgate`
const HTTP_TIMEOUT = time.Second * 120

func setHttpRequestHeaders(req *http.Request, L *lua.LState) error {
	req.Header.Set("User-Agent", HTTP_USER_AGENT)
	req.Header.Set("Content-Type", "application/json")
	return nil
}

/*
	table {
		method: GET/POST,
		url: http://,
		body: <string>,
		ignore_ssl: false,
		proxy: http://user:password@proxy
		basic_auth_user: <string>,
		basic_auth_password: <string>,
		headers: {

		}
}
*/
func getHTTPRequest(t *lua.LTable) (*http.Response, error) {
	// request
	var req *http.Request
	var reqMethod, reqUrl, reqBody string

	// get method
	luaVerb := t.RawGetString(`method`)
	if val, ok := luaVerb.(lua.LString); ok {
		reqMethod = strings.ToUpper(string(val))
	}
	// get url
	luaUrl := t.RawGetString(`url`)
	if val, ok := luaUrl.(lua.LString); ok {
		reqUrl = string(val)
	}
	// get body
	luaBody := t.RawGetString(`body`)
	if val, ok := luaBody.(lua.LString); ok {
		reqBody = string(val)
	}

	var reqErr error
	switch reqMethod {
	case `GET`:
		req, reqErr = http.NewRequest(`GET`, reqUrl, nil)
	case `POST`:
		buf := bytes.NewBuffer([]byte(reqBody))
		req, reqErr = http.NewRequest(`POST`, reqUrl, buf)
	default:
		reqErr = fmt.Errorf(`unknown method: "%s"`, reqMethod)
	}
	if reqErr != nil {
		return nil, reqErr
	}
	req.Header.Set(`User-Agent`, HTTP_USER_AGENT)

	// create client
	client := &http.Client{
		Timeout: HTTP_TIMEOUT,
	}

	var clientTransport *http.Transport
	// set proxy
	luaProxy := t.RawGetString(`proxy`)
	if val, ok := luaProxy.(lua.LString); ok {
		proxyUrl, err := url.Parse(string(val))
		if err != nil {
			return nil, err
		}
		clientTransport.Proxy = http.ProxyURL(proxyUrl)
	}

	// ignore ssl
	luaIgnoreSSL := t.RawGetString(`ignore_ssl`)
	if val, ok := luaIgnoreSSL.(lua.LBool); ok {
		if bool(val) {
			clientTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		}
	}

	// set basic auth
	var basicAuthUser, basicAuthPassword string
	useBasicAuth := false
	luaAuthUser := t.RawGetString(`basic_auth_user`)
	if val, ok := luaAuthUser.(lua.LString); ok {
		useBasicAuth = true
		basicAuthUser = string(val)
	}
	luaAuthPassword := t.RawGetString(`basic_auth_password`)
	if val, ok := luaAuthPassword.(lua.LString); ok {
		useBasicAuth = true
		basicAuthPassword = string(val)
	}
	if useBasicAuth {
		req.SetBasicAuth(basicAuthUser, basicAuthPassword)
	}

	// make request
	client.Transport = clientTransport

	return client.Do(req)
}

func (d *dslConfig) dslHttpRequest(L *lua.LState) int {
	config := L.CheckTable(1)
	response, err := getHTTPRequest(config)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("http error: %s\n", err.Error())))
		return 2
	}
	defer response.Body.Close()
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("http read response error: %s\n", err.Error())))
		return 2
	}
	// write response
	result := L.NewTable()
	L.SetField(result, "code", lua.LNumber(response.StatusCode))
	L.SetField(result, "body", lua.LString(string(data)))
	L.Push(result)
	return 1
}

func (d *dslConfig) dslHttpEscape(L *lua.LState) int {
	query := L.CheckString(1)
	escapedUrl := url.QueryEscape(query)
	L.Push(lua.LString(escapedUrl))
	return 1
}

func (d *dslConfig) dslHttpUnEscape(L *lua.LState) int {
	query := L.CheckString(1)
	url, err := url.QueryUnescape(query)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("unescape error: %s\n", err.Error())))
		return 2
	}
	L.Push(lua.LString(url))
	return 1
}
