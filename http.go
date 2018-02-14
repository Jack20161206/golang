package common

import (
	"crypto/tls"
	"crypto/x509"
	"image"
	"io/ioutil"
	"net"
	"net/http"
	c_url "net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)


/**
* 模拟get请求
* @url		string   需要抓取的url
* @charset	string	字符编码
 */
func HttpGetBody(url string, charset string) (string, int) {
	src := ""
	httpStart := true

	statusCode := 101

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {

		Logs(url + err.Error())
		req, err = http.NewRequest("GET", url, nil)
		if err != nil {
			httpStart = false
			//两次抓取都失败了，需要返回一个空
			return "", statusCode
		}
	}

	if req != nil && req.Body != nil {
		defer req.Body.Close()
	}

	//只有连接成功后，才会写入头的读取字节流
	if httpStart == true {
		//模拟一个host
		sHost := ""
		u, err := c_url.Parse(url)
		if err != nil {
			Logs(url + err.Error())
		} else {
			sHost = u.Host
		}

		if len(charset) > 0 {
			req.Header.Set("Content-Type", "charset="+charset)
		} else {
			req.Header.Set("Content-Type", "charset=UTF-8")
		}
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
		req.Header.Set("Connection", "keep-alive")
		req.Header.Set("Host", sHost)
		req.Header.Set("Referer", sHost)
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.135 Safari/537.36")
		tr := &http.Transport{DisableKeepAlives: true,
			Dial: func(netw, addr string) (net.Conn, error) {
				c, err := net.DialTimeout(netw, addr, time.Second*5) //设置建立连接超时
				if err != nil {
					return nil, err
				}
				c.SetDeadline(time.Now().Add(30 * time.Second)) //设置发送接收数据超时
				return c, nil
			}}
		client := &http.Client{Transport: tr}
		resp, err := client.Do(req)

		if err != nil {
			Logs(url + err.Error())
		} else {

			statusCode = resp.StatusCode
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				Logs(url + err.Error())
			} else {
				src = string(body)
			}
		}
		if resp != nil && resp.Body != nil {
			defer resp.Body.Close()
		}
	}
	return src, statusCode
}


/**
* 模拟POST请求
* @url		string	需要发送请求的url
* @param 	string	需要传递的参数,格式：key=val&key=val
* @gbk		bool 	如果需要抓取的网页是gbk的编码，则需要进行转码
 */
func HttpPostBodyForOutApi(url string, param string) (string, int) {
	src := ""
	statusCode := 101

	req, err := http.NewRequest("POST", url, nil)

	if err != nil {
		Logs(url + err.Error())
		req, err = http.NewRequest("POST", url, nil)
		if err != nil {
			return "", statusCode
		}
	}

	if req != nil && req.Body != nil {
		defer req.Body.Close()
	}

	//模拟一个host
	sHost := ""
	u, err := c_url.Parse(url)
	if err != nil {
		Logs(url + err.Error())
	} else {
		sHost = u.Host
	}

	req.Header.Set("Host", sHost)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.135 Safari/537.36")
	req.Header.Set("Content-Type", "text/plain; charset=UTF-8")
	buf := make([]byte, len(param))
	req.Body = ioutil.NopCloser(strings.NewReader(param))

	req.Body.Read(buf)

	tr := &http.Transport{DisableKeepAlives: true,
		Dial: func(netw, addr string) (net.Conn, error) {
			c, err := net.DialTimeout(netw, addr, time.Second*5) //设置建立连接超时
			if err != nil {
				return nil, err
			}
			c.SetDeadline(time.Now().Add(30 * time.Second)) //设置发送接收数据超时
			return c, nil
		}}
	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		Logs(url + err.Error())
	} else {

		statusCode = resp.StatusCode
		contents, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			Logs(url + err.Error())
		} else {
			src = string(contents)
		}
	}
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	return src, statusCode
}


/**
* 模拟https的POST请求
* @url		string	需要发送请求的url
* @param 	string	需要传递的参数,格式：key=val&key=val
* @gbk		bool 	如果需要抓取的网页是gbk的编码，则需要进行转码
 */
func HttpsPostBody(url, param string, certPath []byte, keyPath []byte, headSet map[string]string, gbk bool) (string, int) {

	src := ""
	statusCode := 101
	//加载安全证书
	cert, err := tls.X509KeyPair(certPath, keyPath)

	if err != nil {
		Logs("证书解密错误->" + err.Error())
		return "", statusCode
	}

	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(certPath)
	_tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
	}

	tr := &http.Transport{TLSClientConfig: _tlsConfig,
		DisableKeepAlives: true,
		Dial: func(netw, addr string) (net.Conn, error) {
			c, err := net.DialTimeout(netw, addr, time.Second*5) //设置建立连接超时
			if err != nil {
				return nil, err
			}
			c.SetDeadline(time.Now().Add(30 * time.Second)) //设置发送接收数据超时
			return c, nil
		}}
	client := &http.Client{Transport: tr}
	httpStart := true

	//模拟post请求
	req, err := http.NewRequest("POST", url, strings.NewReader(param))
	if err != nil {
		httpStart = false
		Logs(url + err.Error())
		//两次连接都失败了，需要返回一个空
		return "", statusCode
	}

	if req != nil && req.Body != nil {
		defer req.Body.Close()
	}

	//只有连接成功后，才会写入头的读取字节流
	if httpStart == true {
		//模拟一个host
		sHost := ""
		u, err := c_url.Parse(url)
		if err != nil {
			Logs(url + err.Error())
		} else {
			sHost = u.Host
		}
		req.Header.Set("Host", sHost)
		req.Header.Set("Referer", sHost)
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.135 Safari/537.36")
		for k, v := range headSet {
			req.Header.Add(k, v)
		}

		if gbk == true {
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=GBK")
		} else {
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
		}

		resp, err := client.Do(req)
		if err != nil {
			Logs(url + err.Error())
		} else {

			statusCode = resp.StatusCode
			contents, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				Logs(url + err.Error())
			} else {
				src = string(contents)
			}
		}
		if resp != nil && resp.Body != nil {
			defer resp.Body.Close()
		}
	}
	return src, statusCode
}
