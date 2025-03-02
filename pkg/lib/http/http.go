package httpUtils

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	authUtils "github.com/easysoft/zentaoatf/internal/pkg/helper/auth"

	"github.com/bitly/go-simplejson"
	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	logUtils "github.com/easysoft/zentaoatf/pkg/lib/log"
	"github.com/fatih/color"
)

func Get(url string) (ret []byte, err error) {
	ret, _, err = GetCheckForward(url)
	return
}

func GetCheckForward(url string) (ret []byte, isForward bool, err error) {
	logUtils.InfofIfVerbose("===DEBUG=== request: %s", url)

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logUtils.InfofIfVerbose(color.RedString("get request failed, error: %s.", err.Error()))
		return
	}

	if commConsts.ExecFrom == commConsts.FromZentao {
		authUtils.AddBearTokenIfNeeded(req)
	} else {
		if !strings.Contains(url, "/tokens") {
			req.Header.Add(commConsts.Token, commConsts.SessionId)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		logUtils.InfofIfVerbose(color.RedString("get request failed, error: %s.", err.Error()))
		return
	}
	defer resp.Body.Close()

	isForward = req.URL.Path != resp.Request.URL.Path

	if !IsSuccessCode(resp.StatusCode) {
		logUtils.InfofIfVerbose(color.RedString("read response failed, StatusCode: %d.", resp.StatusCode))
		err = errors.New(resp.Status)
		return
	}

	ret, err = ioutil.ReadAll(resp.Body)
	logUtils.InfofIfVerbose("===DEBUG=== response: %s", logUtils.ConvertUnicode(ret))

	if err != nil {
		logUtils.InfofIfVerbose(color.RedString("read response failed, error ", err.Error()))
		return
	}

	if len(ret) == 0 {
		return
	}

	jsn, err := simplejson.NewJson(ret)
	if err != nil {
		return
	}
	errMsg, _ := jsn.Get("error").String()
	if strings.ToLower(errMsg) == "unauthorized" {
		err = errors.New(commConsts.UnAuthorizedErr.Message)
		return
	}

	return
}

func Post(url string, data interface{}) (ret []byte, err error) {
	return PostOrPut(url, "POST", data)
}
func Put(url string, data interface{}) (ret []byte, err error) {
	return PostOrPut(url, "PUT", data)
}

func PostOrPut(url string, method string, data interface{}) (ret []byte, err error) {
	logUtils.InfofIfVerbose("===DEBUG=== request: %s", url)

	client := &http.Client{}

	dataBytes, err := json.Marshal(data)
	logUtils.InfofIfVerbose("===DEBUG=== data: %s", string(dataBytes))

	if err != nil {
		logUtils.InfofIfVerbose(color.RedString("marshal request failed, error: %s.", err.Error()))
		return
	}

	dataStr := string(dataBytes)

	req, err := http.NewRequest(method, url, strings.NewReader(dataStr))
	if err != nil {
		logUtils.InfofIfVerbose(color.RedString("post request failed, error: %s.", err.Error()))
		return
	}

	if commConsts.ExecFrom == commConsts.FromZentao {
		authUtils.AddBearTokenIfNeeded(req)
	} else {
		if !strings.Contains(url, "/tokens") {
			req.Header.Add(commConsts.Token, commConsts.SessionId)
		}
	}

	//req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		logUtils.InfofIfVerbose(color.RedString("post request failed, error: %s.", err.Error()))
		return
	}

	if !IsSuccessCode(resp.StatusCode) {
		logUtils.InfofIfVerbose(color.RedString("post request return '%s'.", resp.Status))
		err = errors.New(resp.Status)
		return
	}

	ret, err = ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	logUtils.InfofIfVerbose("===DEBUG=== response: %s", logUtils.ConvertUnicode(ret))

	if err != nil {
		logUtils.InfofIfVerbose(color.RedString("read response failed, error: %s.", err.Error()))
		return
	}

	if len(ret) == 0 {
		return
	}

	jsn, err := simplejson.NewJson(ret)
	if err != nil {
		return
	}
	errMsg, _ := jsn.Get("error").String()
	if strings.ToLower(errMsg) == "unauthorized" {
		err = errors.New(commConsts.UnAuthorizedErr.Message)
		return
	}

	return
}

func IsSuccessCode(code int) (success bool) {
	return code >= 200 && code <= 299
}
