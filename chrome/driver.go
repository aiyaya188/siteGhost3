package chrome

import (
	"errors"
	"fmt"
	"strings"
	"time"

	log "github.com/aiyaya188/go-libs/logger"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

func TodayCookie(str string) (string, error) {
	//str := "zt3uk5s4m22708123l214999"
	i := strings.Index(str, "m")
	sortType := "T"
	filterType := "CB"
	if i == -1 {
		return "", errors.New("uid invalid")
	}
	str = str[i+1:]
	//fmt.Println("str:", str)
	j := strings.Index(str, "l")
	if j == -1 {
		return "", errors.New("uid invalid")
	}
	str = str[:j]
	//fmt.Println("str:", str)
	cookie := "gamePoint_" + str + "=" + time.Now().Format("2006-01-01") + `%2A2%2A0;SortType@` + str + "=" + sortType + ";filterType@" + str + "=" + filterType
	return cookie, nil
	//fmt.Println("cookie:", cookie)
}

//GetWebDriver ...
func GetWebDriver(wdPath string, loadImag bool, ifHeadless bool, setLog bool) (selenium.WebDriver, error) {
	opts := []selenium.ServiceOption{}
	caps := selenium.Capabilities{
		"browserName": "chrome",
	}
	//暂时不要
	//caps.SetLogLevel(log.Performance, log.All)
	var err error
	// 禁止加载图片，加快渲染速度
	//imagCaps := map[string]interface{}{
	//"profile.managed_default_content_settings.images": 2,
	//}
	var chromParam []string
	//ifHeadless := common.GetConfig("system", "headless").String()
	if ifHeadless {
		chromParam = []string{
			"--headless", // 设置Chrome无头模式
			"--no-sandbox",
			//		"--proxy-server=socks://127.0.0.1:9050",
			"--user-agent=Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_2) AppleWebKit/604.4.7 (KHTML, like Gecko) Version/11.0.2 Safari/604.4.7", // 模拟user-agent，防反爬
		}
	} else {
		chromParam = []string{
			"--no-sandbox",
			//"--user-agent=Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_2) AppleWebKit/604.4.7 (KHTML, like Gecko) Version/11.0.2 Safari/604.4.7", // 模拟user-agent，防反爬
		}
	}
	//randomUA := browser.Random()
	//randomUA := browser.Chrome()
	//chromParam = append(chromParam, randomUA)
	chromeCaps := chrome.Capabilities{
		//Prefs: imagCaps,
		Path: "",
		Args: chromParam,
		//Prefs: download,//{'profile.default_content_settings.popups': 0, 'download.default_directory': 'd:\\'}
	}
	caps.AddChrome(chromeCaps)
	// 启动chromedriver，端口号可自定义
	//_, err = selenium.NewChromeDriverService("chromedriver", 9515, opts...)
	_, err = selenium.NewChromeDriverService(wdPath, 9515, opts...)
	log.Info("new hromedriver end")
	if err != nil {
		log.Info(err)
		return nil, errors.New("启动webdriver出错")
	}
	// 调起chrome浏览器
	webDriver, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", 9515))
	if err != nil {
		log.Info(err)
		return nil, errors.New("启动webdriver端口出错")
	}
	return webDriver, nil
}

func Run(wd selenium.WebDriver, addr string) error {
	//wd.SetPageLoadTimeout(120 * time.Second)
	log.Info("begin")
	err := wd.Get(addr)
	if err != nil {
		return errors.New("open timeout")
	}

	log.Info("end")
	wd.MaximizeWindow("")
	title, _ := wd.Title()
	log.Info("title:", title)
	source, err := wd.PageSource()
	if err != nil {
		return err
	}
	log.Infof("打开地址:%s,ok! ", addr)
	log.Infof("抓取内容:%s,ok! ", source)
	return nil
	/*
		reader := bytes.NewReader([]byte(source))
		scanner := bufio.NewScanner(reader)
		m3u8 := ""
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, "html5player.setVideoHLS") {
				//m3u8 = line
				log.Info("line:", line)
				pos1 := strings.Index(line, `'`)
				if pos1 != -1 {
					line = line[pos1+1:]
					pos2 := strings.LastIndex(line, `'`)
					if pos2 != -1 {
						m3u8 = line[:pos2]
					}
				}
			break
			} else {
				continue
			}
		}
		if m3u8 != "" {
			log.Info("m3u8:", m3u8)
			downPath, err1 := b.ProcessM3u8(m3u8, title, addr)
			if err1 != nil {
				return err1
			}
			//回调函数
			handle(vid, downPath, title)
			return nil
		}
		return errors.New("找不到m3u8")
	*/
}

/*
func (b *XSpider) run(wd selenium.WebDriver, addr string, vid uint, handle func(vid uint, downPath string, title string)) error {
	log.Info("pigspideropen addr: ", addr)
	wd.SetPageLoadTimeout(120 * time.Second)
	err := wd.Get(addr)
	if err != nil {
		return errors.New("open timeout")
	}
	wd.MaximizeWindow("")
	title, _ := wd.Title()
	log.Info("title:", title)
	source, err := wd.PageSource()
	if err != nil {
		return err
	}
	log.Infof("pigspideropen addr:%s,ok! ", addr)
	reader := bytes.NewReader([]byte(source))
	scanner := bufio.NewScanner(reader)
	m3u8 := ""
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "html5player.setVideoHLS") {
			//m3u8 = line
			log.Info("line:", line)
			pos1 := strings.Index(line, `'`)
			if pos1 != -1 {
				line = line[pos1+1:]
				pos2 := strings.LastIndex(line, `'`)
				if pos2 != -1 {
					m3u8 = line[:pos2]
				}
			}
		break
		} else {
			continue
		}
	}
	if m3u8 != "" {
		log.Info("m3u8:", m3u8)
		downPath, err1 := b.ProcessM3u8(m3u8, title, addr)
		if err1 != nil {
			return err1
		}
		//回调函数
		handle(vid, downPath, title)
		return nil
	}
	return errors.New("找不到m3u8")
}
*/
