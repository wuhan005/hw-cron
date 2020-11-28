package main

import (
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

	cas.ServiceLogin(fanya)

}
