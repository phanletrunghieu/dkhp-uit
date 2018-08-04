package util

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
)

func HttpPost(client *http.Client, url string, data string) (string, error) {
	// create request
	req, err := http.NewRequest("POST", url, strings.NewReader(data))
	if err != nil {
		return "", err
	}

	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	req.Header.Add("cache-control", "no-cache")

	// create client
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	receivData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(receivData), nil
}

func HttpGet(client *http.Client, url string) (string, error) {
	// create request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("cache-control", "no-cache")

	// create client
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	receivData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(receivData), nil
}

func GetFormToken(data string) (string, error) {
	arrStr1 := strings.Split(data, "form_token")
	if len(arrStr1) <= 1 {
		return "", errors.New("Can't find form_token")
	}

	arrStr2 := strings.Split(arrStr1[1], "value=\"")
	if len(arrStr2) <= 1 {
		return "", errors.New("Can't find form_token")
	}
	runes := []rune(arrStr2[1])
	index := strings.Index(arrStr2[1], "\"")

	form_token := string(runes[0:index])

	return form_token, nil
}

func GetResultOfDKHP(data string) (string, error) {
	data = strings.Replace(data, "\r", "", len(data))
	data = strings.Replace(data, "\n", "", len(data))

	sep := "<div id=\"console\""
	arrStr := strings.Split(data, sep)
	if len(arrStr) <= 1 {
		return "", errors.New("Can't find result")
	}
	data = sep + arrStr[1]

	sep = "</div></div>"
	arrStr = strings.Split(data, sep)
	if len(arrStr) <= 1 {
		return "", errors.New("Can't find result")
	}
	data = arrStr[0] + sep

	// beauty
	data = strings.Replace(data, "  ", "", len(data))
	data = strings.Replace(data, "</div>", "", len(data))
	data = strings.Replace(data, "> <", "><", len(data))
	data = strings.Replace(data, "<div id=\"console\" class=\"clearfix\">", "", len(data))
	data = strings.Replace(data, "<div class=\"alert alert-block alert-success\">", "", len(data))
	data = strings.Replace(data, "<div class=\"alert alert-block alert-error\">", "", len(data))
	data = strings.Replace(data, "<a class=\"close\" data-dismiss=\"alert\" href=\"#\">x</a>", "", len(data))

	data = strings.Replace(data, "<h2 class=\"element-invisible\">Status message</h2>Danh sách lớp đã đăng ký thành công:", "Thành công:\n", len(data))
	data = strings.Replace(data, "<h2 class=\"element-invisible\">Error message</h2>Danh sách lớp bị lỗi:", "Thất bại:\n", len(data))
	data = strings.Replace(data, "<ul>", "", len(data))
	data = strings.Replace(data, "</ul>", "", len(data))
	data = strings.Replace(data, "<li>", "\t- ", len(data))
	data = strings.Replace(data, "</li>", "\n", len(data))

	return data, nil
}
