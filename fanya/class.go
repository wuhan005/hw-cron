package fanya

import "fmt"

// Class 课程
type Class struct {
}

// GetCourseList 返回当前学期课程
func (f *Fanya) GetCourseList() ([]Class, error) {
	term, err := f.getNowTerm()
	fmt.Println(term)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
