package fanya

import (
	"regexp"

	"github.com/pkg/errors"
	log "unknwon.dev/clog/v2"
)

type term struct {
	Begin string
	End   string
}

// NewTerm 返回一个新的学期信息
func NewTerm(begin, end string) *term {
	return &term{
		Begin: begin,
		End:   end,
	}
}

func (f *Fanya) getNowTerm() (*term, error) {
	resp, err := f.casSession.Request().Get("http://hdu.fanya.chaoxing.com/courselist/study")
	if err != nil {
		return nil, err
	}
	termSegment := regexp.MustCompile(`<a onclick="research\(this\);" begin="(.*)" end="(.*)" href="javascript:void\(0\);">`).FindAllStringSubmatch(resp.String(), -1)
	if len(termSegment) == 0 || len(termSegment[0]) < 3 {
		return nil, errors.New("获取泛雅学期信息失败")
	}
	from := termSegment[0][1]
	to := termSegment[0][2]

	log.Trace("Term: %v - %v", from, to)
	return NewTerm(from, to), nil
}
