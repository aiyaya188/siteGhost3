package utils

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var HttpClient *http.Client

func CreateHttpClient() {
	HttpClient = &http.Client{
		//Timeout: 500 * time.Millisecond,
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			//AllowHTTP: true,
			DisableKeepAlives:     true,
			MaxIdleConnsPerHost:   -1,
			ResponseHeaderTimeout: 10 * time.Second,
		},
	}
}

func DoPostForm(host string, resource string, reqBody map[string]string) (string, error) {
	data := url.Values{}
	if reqBody != nil {
		for k, v := range reqBody {
			data.Set(k, v)
		}
	}
	u, _ := url.ParseRequestURI(host)
	u.Path = resource
	urlStr := u.String() // "https://api.com/user/"
	//client := &http.Client{}

	//client := &http.Client{Timeout: 10 * time.Second}
	r, _ := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode())) // URL-encoded payload
	r.Header.Add("User-agent", `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.130 Safari/537.36`)
	r.Header.Add("Referer", host)
	//resp, err := client.Do(r)
	resp, err := HttpClient.Do(r)
	if err != nil {
		return "", err
	}
	if resp.Body != nil {
		respon, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()
		return string(respon), nil
	}
	return "", nil

	//fmt.Println(resp.Status)
}

func DoSimplePost(addr string, Header map[string]string, body string, jsonFormat bool) (string, error) {

	req, err := http.NewRequest("POST", addr, strings.NewReader(body))
	if err != nil {
		return "", err
	}

	if Header != nil {
		for k, v := range Header {
			req.Header.Set(k, v)
		}
	}
	if jsonFormat {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Content-Length", strconv.Itoa(len(body)))
	//client := http.Client{}
	_, err1 := HttpClient.Do(req)
	if err1 != nil {
		return "", err1
	}
	return "", nil
}

func DoPost(addr string, Header map[string]string, body string, jsonFormat bool) (string, error) {
	req, err := http.NewRequest("POST", addr, strings.NewReader(body))
	if err != nil {
		return "", err
	}
	if Header != nil {
		for k, v := range Header {
			req.Header.Set(k, v)
		}
	}
	if jsonFormat {
		req.Header.Set("Content-Type", "application/json")
	} else {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	req.Header.Set("Content-Length", strconv.Itoa(len(body)))
	resp, err := HttpClient.Do(req)
	if err != nil {
		return "", err
	}
	if resp.Body != nil {
		respon, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		/*
			if resp == nil {
				return "", err
			}
		*/
		resp.Body.Close()
		return string(respon), nil
	}
	return "", nil
}
func PostForm() {
	apiUrl := "https://api.com"
	resource := "/user/"
	data := url.Values{}
	data.Set("name", "foo")
	data.Set("surname", "bar")

	u, _ := url.ParseRequestURI(apiUrl)
	u.Path = resource
	urlStr := u.String() // "https://api.com/user/"

	//client := &http.Client{}
	r, _ := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode())) // URL-encoded payload
	r.Header.Add("Authorization", "auth_token=\"XXXXXXX\"")
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, _ := HttpClient.Do(r)
	fmt.Println(resp.Status)
}

func IPFromRequest(r *http.Request, shift int) string {
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	forwordsIP := strings.Split(xForwardedFor, ",")
	if len(forwordsIP) < 2 {
		fmt.Println("err")
		return ""
	}
	newIPs := make([]string, 0)
	lenIP := len(forwordsIP)
	for lenIP > 0 {
		//fmt.Println("index:", i+len(forwordsIP)-1)
		newIPs = append(newIPs, strings.TrimSpace(forwordsIP[lenIP-1]))
		lenIP--
	}
	return newIPs[shift]
}
func IsMobile(r *http.Request) string {
	headers := make(map[string]string)
	if len(r.Header) > 0 {
		for k, v := range r.Header {
			headers[k] = v[0]
			//fmt.Printf("%s=%s\n", k, v[0])
		}
	}
	var is_mobile = false
	via := strings.ToLower(headers["VIA"])
	accept := strings.ToUpper(headers["Accept"])
	HTTP_X_WAP_PROFILE := headers["X_WAP_PROFILE"]
	HTTP_PROFILE := headers["PROFILE"]
	HTTP_USER_AGENT := headers["User-Agent"]
	if via != "" && strings.Index(via, "wap") != -1 {
		is_mobile = true
	} else if accept != "" && strings.Index(accept, "VND.WAP.WML") != -1 {
		is_mobile = true
	} else if HTTP_X_WAP_PROFILE != "" || HTTP_PROFILE != "" {
		is_mobile = true
	} else if HTTP_USER_AGENT != "" {
		reg := regexp.MustCompile(`(?i:(blackberry|configuration\/cldc|hp |hp-|htc |htc_|htc-|iemobile|kindle|midp|mmp|motorola|mobile|nokia|opera mini|opera |Googlebot-Mobile|YahooSeeker\/M1A1-R2D2|android|iphone|ipod|mobi|palm|palmos|pocket|portalmmm|ppc;|smartphone|sonyericsson|sqh|spv|symbian|treo|up.browser|up.link|vodafone|windows ce|xda |xda_|MicroMessenger))`)
		if len(reg.FindAllString(HTTP_USER_AGENT, -1)) > 0 {
			is_mobile = true
		}

	}

	if is_mobile == true {
		return "wap"
	}
	return "pc"

}

func ClientIP(r *http.Request) string {
	/*
		xForwardedFor := r.Header.Get("X-Forwarded-For")
		log.Infof("xForwardedFor:%v", xForwardedFor)
		ip := strings.TrimSpace(strings.Split(xForwardedFor, ",")[0])
		if ip != "" {
			return ip
		}

		ip = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
		log.Infof("realIP:%v", ip)
		if ip != "" {
			return ip
		}
	*/
	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}

	return ""
}

//DoGet 发送post请求 addr:目标地址;Header:指定头信息，如果没有赋予nil;body:发送内容;jsonFormat:是否为json格式
func DoGet(addr string, Header map[string]string) (string, error) {
	//var err error
	req, err := http.NewRequest("GET", addr, nil)
	if err != nil {
		return "", err
	}

	if Header != nil {
		for k, v := range Header {
			req.Header.Set(k, v)
		}
	}
	/*
		transport := &http.Transport{
			TLSHandshakeTimeout: 10 * time.Second,
		}
	*/
	resp, err := HttpClient.Do(req)
	if err != nil {
		if resp != nil {
			resp.Body.Close()
		}
		return "", err
	}

	respon, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		if resp != nil {
			resp.Body.Close()
		}
		return "", err
	}
	resp.Body.Close()
	//fmt.Printf("response:", string(respon))
	return string(respon), nil
}

//DoGet 发送post请求 addr:目标地址;Header:指定头信息，如果没有赋予nil;body:发送内容;jsonFormat:是否为json格式
func DoGetOld(addr string, Header map[string]string) (string, error) {
	//var err error
	req, err := http.NewRequest("GET", addr, nil)
	if err != nil {
		return "", err
	}

	if Header != nil {
		for k, v := range Header {
			req.Header.Set(k, v)
		}
	}
	/*
		transport := &http.Transport{
			TLSHandshakeTimeout: 10 * time.Second,
		}
	*/

	client := &http.Client{Timeout: 5 * time.Second}
	//client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	respon, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	//fmt.Printf("response:", string(respon))
	return string(respon), nil
}

/*
func ClientPublicIP(r *http.Request) string {
	var ip string
	for _, ip = range strings.Split(r.Header.Get("X-Forwarded-For"), ",") {
		ip = strings.TrimSpace(ip)
		if ip != "" && !HasLocalIPddr(ip) {
			return ip
		}
	}

	ip = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	if ip != "" && !HasLocalIPddr(ip) {
		return ip
	}

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		if !HasLocalIPddr(ip) {
			return ip
		}
	}

	return ""
}
*/
/*
// HasLocalIPddr 检测 IP 地址字符串是否是内网地址
func HasLocalIPddr(ip string) bool {
	return HasLocalIP(net.ParseIP(ip))
}

// HasLocalIP 检测 IP 地址是否是内网地址
func HasLocalIP(ip net.IP) bool {
	for _, network := range localNetworks {
		if network.Contains(ip) {
			return true
		}
	}

	return ip.IsLoopback()
}
*/
