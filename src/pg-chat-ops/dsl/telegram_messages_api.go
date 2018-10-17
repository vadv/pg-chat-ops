package dsl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

func (d *dslTelegram) sendSimpleTextMessage(chatID, message, messageType string) error {

	values := url.Values{
		"chat_id": {chatID},
		"text":    {message},
	}
	if messageType != `plain` {
		values.Set(`parse_mode`, messageType)
	}

	_, err := d.botRequest(`sendMessage`, values)
	return err
}

func (d *dslTelegram) sendCallbackMessage(chatID, message string, inlineKeyboard []map[string]string, messageType string) error {

	values := url.Values{
		"chat_id": {chatID},
		"text":    {message},
	}
	if messageType != `plain` {
		values.Set(`parse_mode`, messageType)
	}

	// {"inline_keyboard":[[{"text":"подождать", "callback_data":"1"}, {"text":"убить", "callback_data":"2"}]]}
	jsStr := []string{}
	for _, mmap := range inlineKeyboard {
		data, err := json.Marshal(mmap)
		if err != nil {
			return err
		}
		jsStr = append(jsStr, string(data))
	}
	values.Set("reply_markup", fmt.Sprintf(`{"inline_keyboard":[[ %s ]]}`, strings.Join(jsStr, `,`)))

	_, err := d.botRequest(`sendMessage`, values)
	return err
}

func (d *dslTelegram) sendReplyTextMessage(chatID, messageID, message, messageType string) error {

	values := url.Values{
		"chat_id":             {chatID},
		"text":                {message},
		"reply_to_message_id": {messageID},
	}
	if messageType != `plain` {
		values.Set(`parse_mode`, messageType)
	}

	_, err := d.botRequest(`sendMessage`, values)
	return err
}

func (d *dslTelegram) sendPhotoWithMessage(chatID, message, messageType string, fdPhoto io.ReadCloser) error {

	defer fdPhoto.Close()

	// write temp file
	unixts := strconv.FormatInt(time.Now().UnixNano(), 10)
	fdTmpFile, err := ioutil.TempFile("", unixts)
	if err != nil {
		return err
	}
	defer fdTmpFile.Close()
	defer os.Remove(fdTmpFile.Name())
	if _, err := io.Copy(fdTmpFile, fdPhoto); err != nil {
		return err
	}

	// send file
	botUrl := fmt.Sprintf("%sbot%s/sendPhoto", telegramBaseUrl, d.tocken)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(`photo`, fdTmpFile.Name())
	if err != nil {
		return err
	}
	if _, err := fdTmpFile.Seek(0, 0); err != nil {
		return err
	}
	if _, err := io.Copy(part, fdTmpFile); err != nil {
		return err
	}
	if err := writer.WriteField(`chat_id`, chatID); err != nil {
		return err
	}
	if err := writer.WriteField(`caption`, message); err != nil {
		return err
	}
	if err := writer.WriteField(`parse_mode`, messageType); err != nil {
		return err
	}
	err = writer.Close()
	if err != nil {
		return err
	}

	response, err := d.client.Post(botUrl, writer.FormDataContentType(), body)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, response.Body)
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("bad response code: %d, message: %s", response.StatusCode, buf.Bytes())
	}

	return nil
}
