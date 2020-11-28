package fanya

import (
	"fmt"
	"strings"

	"github.com/antchfx/htmlquery"
	"github.com/imroc/req"
	log "unknwon.dev/clog/v2"
)

// Courses 课程
type Courses struct {
	Name    string
	Teacher string
	School  string
	Link    string
}

// GetCourseList 返回当前学期课程
func (f *Fanya) GetCourseList(term *term) ([]Courses, error) {
	// 获取第一页所有课程信息
	url := fmt.Sprintf("http://hdu.fanya.chaoxing.com/courselist/study?begin=%s&end=%s", term.Begin, term.End)
	resp, err := f.casSession.Request().Get(url, req.Header{
		"User-Agent": userAgent,
	})
	if err != nil {
		return nil, err
	}

	nodes, err := htmlquery.Parse(resp.Response().Body)
	if err != nil {
		log.Warn("Failed to parse fanya class list: %v", err)
		return nil, err
	}

	courseNameNodes := htmlquery.Find(nodes, `//*[@id="zkaikeshenqing"]/div[1]/ul/li/dl/dt`)
	log.Trace("Find %d courses", len(courseNameNodes))

	courses := make([]Courses, len(courseNameNodes))
	for i, node := range courseNameNodes {
		if node.FirstChild != nil {
			courses[i].Name = strings.TrimSpace(node.FirstChild.Data)
		}
	}

	courseTeacherNodes := htmlquery.Find(nodes, `//*[@id="zkaikeshenqing"]/div[1]/ul/li/dl/dd[1]`)
	for i, node := range courseTeacherNodes {
		if node.FirstChild != nil {
			courses[i].Teacher = strings.TrimSpace(node.FirstChild.Data)
		}
	}

	courseSchoolNodes := htmlquery.Find(nodes, `//*[@id="zkaikeshenqing"]/div[1]/ul/li/dl/dd[2]`)
	for i, node := range courseSchoolNodes {
		if node.FirstChild != nil {
			courses[i].School = strings.TrimSpace(node.FirstChild.Data)
		}
	}

	courseLinkNodes := htmlquery.Find(nodes, `//*[@id="zkaikeshenqing"]/div[1]/ul/li/a[1]`)
	for i, node := range courseLinkNodes {
		for _, attr := range node.Attr {
			if attr.Key == "href" {
				courses[i].Link = strings.TrimSpace(attr.Val)
				break
			}
		}
	}

	return courses, nil
}
