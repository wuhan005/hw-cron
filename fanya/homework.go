package fanya

import (
	"regexp"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
	"github.com/imroc/req"
	"github.com/pkg/errors"
	log "unknwon.dev/clog/v2"
)

const (
	NONE_EXPIRED = -62135596800

	TODO     = iota // 待做
	FINISHED        // 已完成、待批阅
	EXPIRED         // 已过期
)

// Homework 作业
type Homework struct {
	Title  string
	Begin  time.Time
	End    time.Time
	Status int
}

func (f *Fanya) GetHomeworks(course Courses) ([]Homework, error) {
	resp, err := f.casSession.Request().Get(course.Link, req.Header{
		"User-Agent": userAgent,
		"Referer":    "https://mooc1-2.chaoxing.com/",
	})
	if err != nil {
		return nil, err
	}

	urls := regexp.MustCompile(` href="javascript:;" data="(.*)" title="作业">作业`).FindStringSubmatch(resp.String())
	if len(urls) < 2 {
		return nil, errors.New("获取作业页面 URL 失败")
	}

	url := "https://mooc1-2.chaoxing.com" + urls[1]
	resp, err = f.casSession.Request().Get(url, req.Header{
		"User-Agent": userAgent,
		"Referer":    "https://mooc1-2.chaoxing.com/",
	})
	if err != nil {
		return nil, err
	}

	nodes, err := htmlquery.Parse(resp.Response().Body)
	if err != nil {
		log.Warn("Failed to parse fanya homework list: %v", err)
		return nil, err
	}

	homeworkTitleNodes := htmlquery.Find(nodes, `//*[@id="RightCon"]/div/div/div[2]/ul/li/div[1]/p/a`)
	homeworks := make([]Homework, len(homeworkTitleNodes))
	for i, node := range homeworkTitleNodes {
		for _, attr := range node.Attr {
			if attr.Key == "title" {
				homeworks[i].Title = strings.TrimSpace(attr.Val)
				break
			}
		}
	}

	homeworkBeginTimeNodes := htmlquery.Find(nodes, `//*[@id="RightCon"]/div/div/div[2]/ul/li/div[1]/span[1]/text()`)
	for i, node := range homeworkBeginTimeNodes {
		homeworks[i].Begin = time.Time{}

		t, err := time.Parse("2006-01-02 15:04", node.Data)
		if err != nil {
			continue
		}
		homeworks[i].Begin = t
	}

	homeworkEndTimeNodes := htmlquery.Find(nodes, `//*[@id="RightCon"]/div/div/div[2]/ul/li/div[1]/span[2]/text()`)
	for i, node := range homeworkEndTimeNodes {
		homeworks[i].End = time.Time{}

		tz, _ := time.LoadLocation("Asia/Shanghai")
		t, err := time.ParseInLocation("2006-01-02 15:04", node.Data, tz)
		if err != nil {
			continue
		}
		homeworks[i].End = t
	}

	homeworkStatusNodes := htmlquery.Find(nodes, `//*[@id="RightCon"]/div/div/div[2]/ul/li/div[1]/span[3]/strong/text()`)
	for i, node := range homeworkStatusNodes {
		status := strings.TrimSpace(node.Data)
		switch status {
		case "待做":
			homeworks[i].Status = TODO
		case "待批阅", "已完成":
			homeworks[i].Status = FINISHED
		case "已过期":
			homeworks[i].Status = EXPIRED
		default:
			homeworks[i].Status = TODO
		}
	}

	return homeworks, nil
}
