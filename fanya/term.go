package fanya

import (
	"fmt"
	"regexp"
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
	fmt.Println(resp.String())
	termSegment := regexp.MustCompile(`<a onclick="research\(this\);" begin="(.*)" end="(.*)" href="javascript:void\(0\);">`).FindAllString(resp.String(), -1)
	fmt.Println(termSegment)
	return NewTerm("", ""), nil
}
