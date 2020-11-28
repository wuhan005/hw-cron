package cas

import (
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/imroc/req"
	log "unknwon.dev/clog/v2"
)

const userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 11_0_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.67 Safari/537.36"

// Session 数字杭电统一认证会话
type Session struct {
	No       string  // 学号
	Password string  // 密码
	Service  Service // 登录后跳转服务

	request *req.Req
}

// NewSession 返回一个新的数字杭电统一认证会话
func NewSession(no, password string) (*Session, error) {
	session := &Session{
		No:       no,
		Password: password,
		request:  req.New(),
	}

	err := session.Login()
	if err != nil {
		return nil, err
	}

	return session, nil
}

// Request 返回 req 请求
func (s *Session) Request() *req.Req {
	return s.request
}

// Login 数字杭电登录
func (s *Session) Login() error {
	lt, err := s.getLoginTicket()
	if err != nil {
		return err
	}

	log.Trace("Login Ticket: %v", lt)

	rsa, err := desEncrypt(s.No + s.Password + lt)
	if err != nil {
		return err
	}

	params := url.Values{
		"rsa":       []string{rsa},
		"ul":        []string{strconv.Itoa(len(s.No))},
		"pl":        []string{strconv.Itoa(len(s.Password))},
		"lt":        []string{lt},
		"execution": []string{"e1s1"},
		"_eventId":  []string{"submit"},
	}

	url := "https://cas.hdu.edu.cn/cas/login"
	resp, err := s.request.Post(url, req.Header{
		"User-Agent":   userAgent,
		"Content-Type": "application/x-www-form-urlencoded",
	}, params.Encode())
	if err != nil {
		return err
	}

	body := resp.String()
	if strings.Contains(body, "抱歉！您的请求出现了异常，请稍后再试。") {
		return CAS_BAD_REQUEST
	}
	if strings.Contains(body, "用户名密码错误") {
		return CAS_ACCOUNT_ERROR
	}

	return nil
}

// ServiceLogin 登录 CAS 服务
func (s *Session) ServiceLogin(svc Service) error {
	svc.SetCasSession(s)

	serviceURL := svc.GetServiceURL()
	resp, err := s.request.Get(serviceURL, req.Header{
		"User-Agent": userAgent,
	})
	if err != nil {
		log.Error("Failed to login to service %s: %v", serviceURL, err)
		return err
	}

	return svc.LoginCallback(resp.String())
}

func (s *Session) getLoginTicket() (string, error) {
	url := "https://cas.hdu.edu.cn/cas/login"
	resp, err := s.request.Get(url, req.Header{
		"User-Agent": userAgent,
	})
	if err != nil {
		return "", err
	}
	body := resp.String()

	if strings.Contains(body, "不允许使用CAS来认证您访问的目标应用") {
		return "", CAS_UNAUTHORISE_SERVICE
	}

	loginTicketGroup := regexp.MustCompile(`<input type="hidden" id="lt" name="lt" value="(.*)"`).FindStringSubmatch(body)
	if len(loginTicketGroup) != 2 {
		return "", CAS_LOGIN_TICKET_NOT_FOUND
	}

	return loginTicketGroup[1], nil
}
