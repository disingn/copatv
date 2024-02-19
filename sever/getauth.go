package sever

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"runtime"
)

type AuthBata struct {
	Name         string `json:"name"`
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}
type GetAuthData struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Token      string `json:"token"`
		Url        string `json:"url"`
		ExpireTime string `json:"expire_time"`
	} `json:"data"`
}

type GetVersionData struct {
	Message    string `json:"message"`
	Token      string `json:"token"`
	Url        string `json:"url"`
	ExpireTime string `json:"expire_time"`
}

func Getauth(name, token string) (error, GetVersionData) {
	url := "https://gc.aiu.im/api/auth/"
	method := "POST"
	var clid string

	switch runtime.GOOS {
	case "windows":
		c, err := GetWinHardwareUUID()
		if err != nil {

			return err, GetVersionData{}
		}
		clid = c
	case "darwin":
		c, err := GetMacHardwareUUID()
		if err != nil {

			return err, GetVersionData{}
		}
		clid = c
	}
	//log.Print("clid:", clid)
	payload := AuthBata{
		Name:         name,
		ClientId:     clid,
		ClientSecret: token,
	}
	payloadStr, _ := json.Marshal(payload)
	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewReader(payloadStr))

	if err != nil {

		return err, GetVersionData{}
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Host", "gc.aiu.im")
	req.Header.Add("Connection", "keep-alive")

	res, err := client.Do(req)
	if err != nil {

		return err, GetVersionData{}
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {

		return err, GetVersionData{}
	}
	if res.StatusCode != 200 {

		return fmt.Errorf("授权码已使用或无权限"), GetVersionData{}
	}
	var data GetAuthData
	err = json.Unmarshal(body, &data)
	if err != nil {
		return err, GetVersionData{}
	}
	return nil, GetVersionData{
		Message:    data.Message,
		Token:      data.Data.Token,
		Url:        data.Data.Url,
		ExpireTime: data.Data.ExpireTime,
	}

}

func OpenBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "rundll32"
		args = []string{"url.dll,FileProtocolHandler"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}
