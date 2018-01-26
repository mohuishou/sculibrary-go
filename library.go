package library

import (
	"errors"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/gocolly/colly"
)

type Library struct {
	URL string
	c   *colly.Collector
}

// LoanBook 借阅的书籍
type LoanBook struct {
	Author      string
	Title       string
	PublishYear int     // 出版年
	DueDate     string  // 到期日期
	ReturnDate  string  // 归还日期(借阅历史)
	ReturnTime  string  // 归还时间(借阅历史)
	Arrearage   float64 // 欠费
	Address     string  // 分馆
	Number      string  // 索书号(当前借阅)
}

// NewLibrary 新建一个图书馆对象
func NewLibrary(studentID, password string) (*Library, error) {
	c := colly.NewCollector()
	lib := &Library{c: c}

	lib.URL = getURL()

	if lib.URL == "" {
		return nil, errors.New("获取url失败！")
	}
	loginErr := 0
	lib.c.OnHTML("#feedbackbar", func(e *colly.HTMLElement) {
		if strings.Contains(e.Text, "错") {
			loginErr = 1 //账号或密码错误
		}
	})

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

	if loginErr != 0 {
		return nil, errors.New("账号或密码错误！")
	}
	return lib, nil
}

// GetLoan 获取当前借阅
func (lib *Library) GetLoan() []LoanBook {
	books := make([]LoanBook, 0)
	lib.c.OnHTML("body > center > center table", func(e *colly.HTMLElement) {
		e.DOM.Find("tr").Each(func(i int, s *goquery.Selection) {
			if s.Find("td:nth-child(6)").Text() == "" {
				return
			}
			book := LoanBook{}
			v := reflect.ValueOf(&book)
			elem := v.Elem()
			typeOfBook := elem.Type()
			eq := 2
			for k := 0; k < elem.NumField(); k++ {
				val := strings.TrimSpace(s.Find("td").Eq(eq).Text())
				switch typeOfBook.Field(k).Name {
				case "ReturnDate", "ReturnTime":
				case "PublishYear":
					v, _ := strconv.Atoi(val)
					elem.Field(k).SetInt(int64(v))
					eq++
				case "Arrearage":
					v, _ := strconv.ParseFloat(val, 10)
					elem.Field(k).SetFloat(v)
					eq++
				default:
					elem.Field(k).SetString(val)
					eq++
				}
			}
			books = append(books, book)
		})
	})
	lib.c.Visit(lib.URL + "?func=bor-loan&adm_library=SCU50")
	return books
}

// GetLoanAll 获取历史借阅
func (lib *Library) GetLoanAll() []LoanBook {
	books := make([]LoanBook, 0)
	lib.c.OnHTML("body > center table", func(e *colly.HTMLElement) {
		e.DOM.Find("tr").Each(func(i int, s *goquery.Selection) {
			if s.Find("td:nth-child(6)").Text() == "" {
				return
			}
			book := LoanBook{}
			v := reflect.ValueOf(&book)
			elem := v.Elem()
			typeOfBook := elem.Type()
			eq := 2
			for k := 0; k < elem.NumField(); k++ {
				val := strings.TrimSpace(s.Find("td").Eq(eq).Text())
				switch typeOfBook.Field(k).Name {
				case "Number":
				case "PublishYear":
					v, _ := strconv.Atoi(val)
					elem.Field(k).SetInt(int64(v))
					eq++
				case "Arrearage":
					v, _ := strconv.ParseFloat(val, 10)
					elem.Field(k).SetFloat(v)
					eq++
				default:
					elem.Field(k).SetString(val)
					eq++
				}
			}
			books = append(books, book)
		})
	})
	lib.c.Visit(lib.URL + "?func=bor-history-loan&adm_library=SCU50")
	return books
}

func (lib *Library) Loan() {

}

func (lib *Library) LoanAll() {

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
