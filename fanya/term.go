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

func (f *Fanya) GetAllTerm() ([]*term, error) {
	resp, err := f.casSession.Request().Get("http://hdu.fanya.chaoxing.com/courselist/study")
	if err != nil {
		return nil, err
	}
	termSegment := regexp.MustCompile(`<a onclick="research\(this\);" begin="(.*)" end="(.*)" href="javascript:void\(0\);">`).FindAllStringSubmatch(resp.String(), -1)
	if len(termSegment) == 0 || len(termSegment[0]) < 3 {
		return nil, errors.New("获取泛雅学期信息失败")
	}

	terms := make([]*term, len(termSegment))
	for i := range terms {
		terms[i] = NewTerm(termSegment[i][1], termSegment[i][2])
	}
	return terms, nil
}

func (f *Fanya) GetNowTerm() (*term, error) {
	terms, err := f.GetAllTerm()
	if err != nil {
		return nil, err
	}
	term := terms[0]

	log.Trace("Term: %v - %v", term.Begin, term.End)
	return term, nil
}
