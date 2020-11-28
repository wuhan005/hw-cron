package fanya

import (
	"net/url"
	"regexp"
	"strings"

	"github.com/imroc/req"
	"github.com/pkg/errors"
	log "unknwon.dev/clog/v2"

	"github.com/wuhan005/hw-cron/cas"
)

const userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 11_0_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.67 Safari/537.36"

var _ cas.Service = new(Fanya)

type Fanya struct {
	casSession *cas.Session
}

func (f *Fanya) GetServiceURL() string {
	return "http://hdu.fanya.chaoxing.com/sso/hdu"
}

func (f *Fanya) SetCasSession(casSession *cas.Session) {
	f.casSession = casSession
}

// LoginCallback 泛雅登录后跳转为前端 JS 提交登录表单，这里需要模拟。
func (f *Fanya) LoginCallback(body string) error {
	// http://hdu.fanya.chaoxing.com/sso/logindsso
	body, err := f.submitForm(body)
	if err != nil {
		return err
	}

	// http://passport2.chaoxing.com/loginfanya
	body, err = f.submitForm(body)
	if err != nil {
		return err
	}

	// 检查登录状态
	resp, err := f.casSession.Request().Get("http://hdu.fanya.chaoxing.com/topjs?index=quote", req.Header{
		"User-Agent": userAgent,
		"Referer":    "http://hdu.fanya.chaoxing.com/portal",
	})
	if err != nil {
		return err
	}

	if strings.Contains(resp.String(), "数字杭电登录") {
		return errors.Errorf("泛雅 SSO 登录失败")
	}

	return nil
}

func (f *Fanya) submitForm(body string) (string, error) {
	urls := regexp.MustCompile(`<form action="(.*)" method="post"`).FindStringSubmatch(body)
	if len(urls) != 2 {
		return "", errors.Errorf("Failed to parse request url: %v", urls)
	}
	requestURL := urls[1]
	log.Trace("Submit form to login %v", requestURL)

	params := url.Values{}
	fields := regexp.MustCompile(`<input type="hidden" name="(.*)" value="(.*)"/>`).FindAllStringSubmatch(body, -1)
	for _, field := range fields {
		if len(field) < 3 {
			return "", errors.Errorf("Failed to parse request param: %v", field)
		}
		params.Set(field[1], field[2])
	}

	resp, err := f.casSession.Request().Post(requestURL, req.Header{
		"User-Agent":   userAgent,
		"Content-Type": "application/x-www-form-urlencoded",
	}, params.Encode())
	if err != nil {
		return "", err
	}
	return resp.String(), nil
}

// New 返回一个新的泛雅会话
func New() *Fanya {
	fanya := &Fanya{}
	return fanya
}
