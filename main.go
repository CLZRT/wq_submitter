package main

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"

	"wq_submitter/configs"
	"wq_submitter/internal/scheduler"
	"wq_submitter/router"
)

func init() {
	configs.InitGlobalConfig()
}
func main() {

	config := configs.GetGlobalConfig()

	log.Infof("The service %s starting", config.AppConfig.AppName)

	ctx, cancelFunc := context.WithCancel(context.Background())
	ideaScheduler := scheduler.NewIdeaScheduler(ctx, cancelFunc)
	ideaScheduler.Run()
	defer func(ideaScheduler *scheduler.IdeaScheduler) {
		ideaScheduler.Stop()
	}(ideaScheduler)

	r := router.SetRouter()
	if err := r.Run(fmt.Sprintf(":%d", config.AppConfig.Port)); err != nil {
		log.Errorf("server run error: %v", err)
	}

}
