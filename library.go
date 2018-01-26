package library

import (
	"errors"
	"net/url"

	"github.com/gocolly/colly"
)

type Library struct {
	URL string
	c   *colly.Collector
}

// NewLibrary 新建一个图书馆对象
func NewLibrary(studentID, password string) (*Library, error) {
	c := colly.NewCollector()
	lib := &Library{c: c}

	lib.URL = getURL()

	if lib.URL == "" {
		return nil, errors.New("获取url失败！")
	}

	// login
	err := lib.c.Post(lib.URL, map[string]string{
		"func":             "login-session",
		"login_source":     "bor-info",
		"bor_id":           studentID,
		"bor_verification": password,
		"bor_library":      "SCU50",
	})
	if err != nil {
		return nil, err
	}
	return lib, nil
}

func (lib *Library) GetLoan() {
	lib.c.OnHTML("body > center > center > table:nth-child(6) > tbody > tr", func (e *colly.HTMLElement)  {
		
	})
}

func (lib *Library) GetLoanAll() {

}

func (lib *Library) Loan() {

}

func (lib *Library) LoanAll() {

}

func ()  {
	
}

func getURL() string {
	c := colly.NewCollector()
	urlstr := ""
	c.OnHTML("#header > a:nth-child(1)", func(e *colly.HTMLElement) {
		urlstr = e.Attr("href")

		if urlstr != "" {
			uri, err := url.Parse(urlstr)
			if err == nil {
				urlstr = "http://opac.scu.edu.cn:8080" + uri.EscapedPath()
			}
		}
	})
	c.Visit("http://opac.scu.edu.cn:8080/F")
	return urlstr
}
