package request

import (
	"net/http"
	"io/ioutil"
	"net/url"
	"os"
	"io"
	"bytes"
)

func PostData(url string, data url.Values) (ret string, err error) {
	resp, err := http.PostForm(url, data)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	ret = string(body)
	return
}

func GetData(url string) (ret string, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	ret = string(body)
	return

}

func DownFile(url, dst string) (err error) {
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer out.Close()

	resp, err := http.Get(url)
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return
	}
	return
}

func PostJson(url string,b []byte)(ret string, err error){
	body := bytes.NewBuffer([]byte(b))
	res,err := http.Post(url, "application/json;charset=utf-8", body)
	defer res.Body.Close()

	if err != nil {
		return
	}
	result, err := ioutil.ReadAll(res.Body)
	ret = string(result)
	return
}