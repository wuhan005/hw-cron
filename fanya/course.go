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
func (f *Fanya) GetCourseList() ([]Courses, error) {
	//term, err := f.getNowTerm()
	//if err != nil {
	//	return nil, err
	//}

	// 获取第一页所有课程信息
	url := fmt.Sprintf("http://hdu.fanya.chaoxing.com/courselist/study")
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
	courses := make([]Courses, len(courseNameNodes))
	for i, node := range courseNameNodes {
		courses[i].Name = strings.TrimSpace(node.FirstChild.Data)
	}

	courseTeacherNodes := htmlquery.Find(nodes, `//*[@id="zkaikeshenqing"]/div[1]/ul/li/dl/dd[1]`)
	for i, node := range courseTeacherNodes {
		courses[i].Teacher = strings.TrimSpace(node.FirstChild.Data)
	}

	courseSchoolNodes := htmlquery.Find(nodes, `//*[@id="zkaikeshenqing"]/div[1]/ul/li/dl/dd[2]`)
	for i, node := range courseSchoolNodes {
		courses[i].School = strings.TrimSpace(node.FirstChild.Data)
	}

	courseLinkNodes := htmlquery.Find(nodes, `//*[@id="zkaikeshenqing"]/div[1]/ul/li/a`)
	for i, node := range courseLinkNodes {
		for _, attr := range node.Attr {
			if attr.Key == "href" {
				courses[i].Link = strings.TrimSpace(attr.Val)
				break
			}
		}
	}

	log.Trace("Find %d courses", len(courses))

	return courses, nil
}
