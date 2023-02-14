package crawler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func apiPost(apiUrl string, form url.Values) map[string]interface{} {
	resp, err := http.PostForm(apiUrl, form)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	res := make(map[string]interface{})
	json.Unmarshal(body, &res)
	return res
}
