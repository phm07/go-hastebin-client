package pasteclient

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
)

func PasteAndReturnKey(baseUrl, content string) (res string, err error) {
	var docsUrl string
	docsUrl, err = url.JoinPath(baseUrl, "/documents")
	if err != nil {
		return
	}

	var req *http.Request
	req, err = http.NewRequest("POST", docsUrl, strings.NewReader(content))
	if err != nil {
		return
	}

	client := http.Client{}
	var resp *http.Response
	resp, err = client.Do(req)
	if err != nil {
		return
	}
	defer func() {
		err = errors.Join(err, resp.Body.Close())
	}()

	var schema = struct {
		Key string `json:"key"`
	}{}
	err = json.NewDecoder(resp.Body).Decode(&schema)
	if err != nil {
		return
	}

	res = schema.Key
	return
}

func PasteAndReturnUrl(baseUrl, content, ext string) (string, error) {
	key, err := PasteAndReturnKey(baseUrl, content)
	if err != nil {
		return "", err
	}
	if len(ext) > 0 && ext[0] != '.' {
		ext = "." + ext
	}
	return url.JoinPath(baseUrl, "/"+key+ext)
}
