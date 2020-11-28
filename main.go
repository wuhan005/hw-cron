package main

import (
	"fmt"

	"github.com/wuhan005/hw-cron/cas"
	"github.com/wuhan005/hw-cron/fanya"
	log "unknwon.dev/clog/v2"
)

func init() {
	_ = log.NewConsole()
}

func main() {
	fanya := fanya.New()

	cas, err := cas.NewSession("", "")
	if err != nil {
		log.Fatal("Failed to login: %v", err)
	}

	err = cas.ServiceLogin(fanya)
	if err != nil {
		log.Fatal("Failed to login to fanya: %v", err)
	}

	courses, err := fanya.GetCourseList()
	if err != nil {
		log.Fatal("Failed to get courses list: %v", err)
	}

	fmt.Println(fanya.GetHomeworks(courses[2]))

}
