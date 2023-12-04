package gohttp

import (
	"fmt"
	"net/http"
	"strings"
	

)

func httpDo(url string) {
	// http.Clicent:是一个HTTP客户端，客户机比往返器(比如传输)更高级,还处理HTTP细节，比如cookie，重定向，长连接。
		client := &http.Client{}
	 // 使用 NewRequest 设置头参数、cookie之类的数据，
		 req, err := http.NewRequest("POST", url, strings.NewReader("name=cjb"))
		if err != nil {
			// handle error
		}
	
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Cookie", "name=anny")
	
		resp, err := client.Do(req)
	
		defer resp.Body.Close()
	
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			// handle error
		}
	
		fmt.Println(string(body))
	}
	