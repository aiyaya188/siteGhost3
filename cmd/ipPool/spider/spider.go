package spider

import (
	"fmt"

	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

func GetTargetContent(target string) (err error) {
	var (
		driverPath = "chromedriver"
		port       = 9222
	)
	//var content string
	//var err error
	service, err := selenium.NewChromeDriverService(driverPath, port, []selenium.ServiceOption{}...)
	if nil != err {
		fmt.Println("start a chromedriver service failed", err.Error())
		return
	}
	defer func() {
		_ = service.Stop()
	}()

	// Connect to the WebDriver instance running locally.
	caps := selenium.Capabilities{"browserName": "chrome"}
	//禁止图片加载，加快渲染速度
	imagCaps := map[string]interface{}{
		"profile.managed_default_content_settings.images": 2,
	}
	chromeCaps := chrome.Capabilities{
		Prefs: imagCaps,
		Path:  "",
		Args: []string{
			"--headless",
			"--no-sandbox",
			"--disable-gpu-sandbox",
			"--user-agent=Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.79 Safari/537.36",
		},
	}
	//以上是设置浏览器参数
	caps.AddChrome(chromeCaps)
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	if err != nil {
		fmt.Println("connect to the webDriver failed", err.Error())
		return
	}
	defer func() {
		_ = wd.Quit()
	}()
	fmt.Println("抓取地址:", target)
	err = wd.Get(target)
	if err != nil {
		fmt.Println("get page failed", err.Error())
		return
	}
	fmt.Println("打开完毕")
	err = wd.Wait(func(wdTemp selenium.WebDriver) (b bool, e error) {
		tit, err := wdTemp.Title()
		if err != nil {
			return false, nil
		}
		if tit == "" {
			return false, nil
		}
		return true, nil
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	source, _ := wd.PageSource()
	fmt.Println("source:", source)
	/*
		fmt.Println("获取body")
		tbody, err := wd.FindElement(selenium.ByTagName, "tbody")
		if err != nil {
			return
		}
		fmt.Println("获取tr")
		tr, err := tbody.FindElements(selenium.ByTagName, "tr")
		if err != nil {
			return
		}

		fmt.Println("获取td")
		for _, v := range tr {
			tds, err := v.FindElements(selenium.ByTagName, "td")
			if err != nil {
				continue
			}
			if len(tds) < 6 {
				continue
			}
			text, _ := tds[0].Text()
			fmt.Println("0:", text)

			text, _ = tds[1].Text()
			fmt.Println("1:", text)

			text, _ = tds[2].Text()
			fmt.Println("2:", text)

			text, _ = tds[3].Text()
			fmt.Println("3:", text)

			text, _ = tds[4].Text()
			fmt.Println("4:", text)

			text, _ = tds[5].Text()
			fmt.Println("5:", text)
		}
	*/

	fmt.Println("获取完毕")
	return
}

func Displayed(by, elementName string) func(selenium.WebDriver) (bool, error) {
	return func(wd selenium.WebDriver) (ok bool, err error) {
		var el selenium.WebElement
		el, err = wd.FindElement(by, elementName)
		if err != nil {
			return
		}
		ok, err = el.IsDisplayed()
		return
	}
}
