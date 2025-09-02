package scheduler

import (
	"context"
	"github.com/panjf2000/ants"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
	"wq_submitter/configs"
	"wq_submitter/internal/submitter"
	"wq_submitter/internal/svc"
	"wq_submitter/internal/viewer"
)

//扫描ValidIdea并作比较创建或者关闭IdeaSubmitter

type IdeaScheduler struct {
	ticker           time.Ticker
	once             sync.Once
	ideaMap          map[int64]*submitter.Submitter
	concurrencyLimit int64
	ctx              context.Context
	cancelFunc       context.CancelFunc
	workerPool       *ants.Pool
}

var conf *configs.GlobalConfig

func init() {
	conf = configs.GetGlobalConfig()
}
func NewIdeaScheduler(context context.Context, cancelFunc context.CancelFunc) *IdeaScheduler {
	workerPool, err := ants.NewPool(5)
	if err != nil {
		log.Errorf("IdeaScheduler NewIdeaScheduler Error: %s", err.Error())
		return nil
	}
	return &IdeaScheduler{
		ticker:           *time.NewTicker(time.Second),
		once:             sync.Once{},
		ideaMap:          make(map[int64]*submitter.Submitter),
		concurrencyLimit: conf.AppConfig.Concurrency,
		ctx:              context,
		cancelFunc:       cancelFunc,
		workerPool:       workerPool,
	}
}

func (s *IdeaScheduler) Run() {
	err := s.workerPool.Submit(s.work)
	if err != nil {
		log.Errorf("IdeaScheduler Run Error: %s", err.Error())
		return
	}
}

func (s *IdeaScheduler) work() {
	defer func() {
		if r := recover(); r != nil {
			s.Stop()
		}
	}()
	for range s.ticker.C {

		ideas := s.scanValidIdea()
		//  判断是否超过并发限制
		if s.isOverConcurrencyLimit(ideas) {
			log.Error("IdeaScheduler isOverConcurrencyLimit")
			panic("IdeaScheduler isOverConcurrencyLimit")
		}
		//  找到当前在跑,但是不在新的Ideas里面的,关闭的时候要保证正在跑的Alpha跑完
		tempMap := s.storeIdea2Map(ideas)
		for ideaId := range s.ideaMap {
			if _, ok := tempMap[ideaId]; ok {
				continue
			}
			s.deleteSubmitter(s.ideaMap[ideaId].Idea)
		}

		//  找到当前不在跑，但是未来要跑的Idea
		//  已经存在的IdeaSubmitter的并发度更新由它自己完成
		for _, ideaViewer := range ideas {
			if _, ok := s.ideaMap[ideaViewer.ID]; ok {
				continue
			}
			s.newSubmitter(ideaViewer)
		}
	}

}

func (s *IdeaScheduler) Stop() {
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("IdeaScheduler Stop Error: %s", r)
		}
	}()

	s.once.Do(func() {
		s.ticker.Stop()
		s.cancelFunc()
		err := s.workerPool.Release()
		if err != nil {
			log.Errorf("IdeaScheduler Stop Error: %s", err.Error())
			return
		}
	})

}

func (s *IdeaScheduler) isOverConcurrencyLimit(ideaViewers []viewer.Idea) bool {
	var curConcurrency int64
	for _, ideaViewer := range ideaViewers {
		curConcurrency += ideaViewer.ConcurrencyNum
	}
	return curConcurrency > s.concurrencyLimit
}
func (s *IdeaScheduler) storeIdea2Map(ideaViewers []viewer.Idea) map[int64]struct{} {
	ideaMap := make(map[int64]struct{})
	for _, ideaViewer := range ideaViewers {
		ideaMap[ideaViewer.ID] = struct{}{}
	}
	return ideaMap

}
func (s *IdeaScheduler) scanValidIdea() []viewer.Idea {
	ideas, err := svc.GetNeedRunIdea()
	if err != nil {
		log.Errorf("IdeaScheduler scanValidIdea Error: %s", err.Error())
		return nil
	}
	return ideas
}

func (s *IdeaScheduler) newSubmitter(ideaViewer viewer.Idea) {
	ctx, cancelFunc := context.WithCancel(s.ctx)
	s.ideaMap[ideaViewer.ID] = submitter.NewSubmitter(ctx, cancelFunc, ideaViewer)
	s.ideaMap[ideaViewer.ID].Run()
}
func (s *IdeaScheduler) deleteSubmitter(ideaViewer viewer.Idea) {
	err := s.ideaMap[ideaViewer.ID].Stop()
	if err != nil {
		log.Errorf("IdeaScheduler deleteSubmitter Error: %s", err.Error())
		return
	}
	delete(s.ideaMap, ideaViewer.ID)

}
