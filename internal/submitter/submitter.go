package submitter

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/panjf2000/ants"
	log "github.com/sirupsen/logrus"
	"sync"
	"sync/atomic"
	"time"
	"wq_submitter/configs"
	"wq_submitter/internal/constant"
	"wq_submitter/internal/svc"

	"wq_submitter/internal/viewer"
)

type AlphaTask struct {
	ID             int64
	IdeaID         int64
	SimulationData json.RawMessage
	RetryNum       int64
	Status         int64
}
type SafeChan struct {
	alphaTaskChan chan AlphaTask
	once          sync.Once
	isClosed      bool
	mutex         sync.Mutex
}

func NewSafeChan(len int64) SafeChan {
	return SafeChan{
		alphaTaskChan: make(chan AlphaTask, len),
		once:          sync.Once{},
		isClosed:      false,
		mutex:         sync.Mutex{},
	}

}

func (safeChan *SafeChan) Write(alphaTask AlphaTask) {
	safeChan.mutex.Lock()
	defer safeChan.mutex.Unlock()
	if safeChan.isClosed {
		return
	}
	safeChan.alphaTaskChan <- alphaTask

}

func (safeChan *SafeChan) Close() {
	safeChan.mutex.Lock()
	defer safeChan.mutex.Unlock()
	safeChan.once.Do(func() {
		safeChan.isClosed = true
		close(safeChan.alphaTaskChan)
	})
}

func (safeChan *SafeChan) GetReadChan() <-chan AlphaTask {
	return safeChan.alphaTaskChan

}

type Submitter struct {
	context          context.Context
	cancelFunc       context.CancelFunc
	ScanIdeaSecond   int64
	AlphaChan        SafeChan
	DeadChan         SafeChan
	FinishChan       SafeChan
	FinishAlphaIdMap map[int64]struct{}
	ConcurrencyNum   int64
	ChannelLen       int64
	Idea             viewer.Idea
	NextScanId       atomic.Int64
	scanIdMutex      sync.Mutex
	IdeaNextId       int64
	WorkerPool       *ants.Pool
	CancelFuncList   []context.CancelFunc
	mutex            sync.Mutex
	once             sync.Once
}

func NewSubmitter(ctx context.Context, cancelFunc context.CancelFunc, idea viewer.Idea) *Submitter {
	workerPool, err := ants.NewPool(5)
	conf := configs.GetGlobalConfig()
	if err != nil {
		log.Errorf(" IdeaTitle: %s Init Worker Pool Failed", idea.IdeaTitle)
		return nil
	}
	channelLen := conf.AlphaConfig.ChannelLen
	submitter := &Submitter{
		context:          ctx,
		cancelFunc:       cancelFunc,
		ScanIdeaSecond:   conf.AlphaConfig.ScanIdeaSecond,
		AlphaChan:        NewSafeChan(channelLen),
		DeadChan:         NewSafeChan(channelLen),
		FinishChan:       NewSafeChan(channelLen),
		FinishAlphaIdMap: make(map[int64]struct{}),
		ConcurrencyNum:   idea.ConcurrencyNum,
		Idea:             idea,
		ChannelLen:       channelLen,
		NextScanId:       atomic.Int64{},
		IdeaNextId:       idea.NextIdx,
		WorkerPool:       workerPool,
		CancelFuncList:   make([]context.CancelFunc, 0),
		mutex:            sync.Mutex{},
	}
	//赋值NextScanId
	submitter.NextScanId.Store(idea.NextIdx)
	return submitter

}
func (s *Submitter) Run() bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	//开启消费者
	for i := int64(0); i < s.ConcurrencyNum; i++ {
		subCtx, cancelFunc := context.WithCancel(s.context)
		s.CancelFuncList = append(s.CancelFuncList, cancelFunc)
		go s.executeAlpha(subCtx)
	}
	//开启初始Alpha加载线程,一次性
	// 装载alphas进 alphaChan
	s.initLoadAlpha()

	//开启扫描线程,关闭FinishChan就退出了
	err := s.WorkerPool.Submit(s.afterAlphaFinish)
	if err != nil {
		log.Errorf("submit afterAlphaFinish err: %s", err.Error())
		return false
	}
	//开启重试线程,关闭DeadChan就退出了
	err = s.WorkerPool.Submit(s.retryAlpha)
	if err != nil {
		log.Errorf("submit retryAlpha err: %s", err.Error())
		return false
	}
	//开启扫描线程
	err = s.WorkerPool.Submit(s.scanIdea)
	if err != nil {
		log.Errorf("submit scanIdea err: %s", err.Error())
		return false
	}
	return true

}
func (s *Submitter) Stop() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.cancelFunc()
	s.AlphaChan.Close()
	s.DeadChan.Close()
	s.FinishChan.Close()
	err := s.WorkerPool.Release()
	if err != nil {
		err = fmt.Errorf("fail to close worker pool: %s", err.Error())
		log.Error(err.Error())
		return err
	}
	return nil

}

func (s *Submitter) initLoadAlpha() {

	alphas, err := s.batchGetAlphaById(s.NextScanId.Load(), s.ChannelLen)
	s.NextScanId.Add(int64(len(alphas)))
	if err != nil {
		log.Error(err.Error())
	}

	for _, alpha := range alphas {
		alphaTask := AlphaTask{
			ID:             alpha.ID,
			IdeaID:         alpha.IdeaID,
			SimulationData: alpha.SimulationData,
			RetryNum:       0,
			Status:         0,
		}
		s.AlphaChan.Write(alphaTask)
	}
}

// 后台线程
func (s *Submitter) scanIdea() {
	ticker := time.NewTicker(time.Duration(s.ScanIdeaSecond) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		ideaViewer := svc.GetIdeaById(s.Idea.ID)
		if ideaViewer.ConcurrencyNum != s.ConcurrencyNum {
			s.changeConcurrency(ideaViewer.ConcurrencyNum)
		}
		s.Idea = ideaViewer
	}
}
func (s *Submitter) retryAlpha() {
	config := configs.GetGlobalConfig()

	for alphaTask := range s.DeadChan.GetReadChan() {

		alphaTask.RetryNum++
		// 超过最大重试次数
		if alphaTask.RetryNum > config.AlphaConfig.RetryNum {
			s.FinishChan.Write(alphaTask)
			continue
		}
		// 	未超过最大重试次数
		s.AlphaChan.Write(alphaTask)
	}

}
func (s *Submitter) afterAlphaFinish() {
	for alphaTask := range s.FinishChan.GetReadChan() {
		//更新状态
		s.FinishAlphaIdMap[alphaTask.ID] = struct{}{}
		s.updateIdeaStatus(alphaTask)
		s.updateAlphaStatus(alphaTask.ID, alphaTask.Status)

		//获取新alpha 并写入 alphaChan
		//判断 要获取的alphaId是否超过当前Idea的范围
		//加载下一个要扫描的新Alpha,并更新nextScanId+1
		scanId := s.loadNumAndAdd(&s.NextScanId)
		if scanId <= s.Idea.EndIdx {

			alpha := s.getAlphaById(scanId)
			var newAlphaTask AlphaTask
			if alpha != nil {
				newAlphaTask = AlphaTask{
					ID:             alpha.ID,
					IdeaID:         alpha.IdeaID,
					SimulationData: alpha.SimulationData,
					RetryNum:       0,
					Status:         alpha.IsSubmitted,
				}
				s.AlphaChan.Write(newAlphaTask)
			}

		}

	}
}

// 改变并发
func (s *Submitter) changeConcurrency(concurrencyNum int64) {

	if concurrencyNum > s.ConcurrencyNum {
		s.addWorker(concurrencyNum - s.ConcurrencyNum)
	} else {
		s.shrinkWorker(s.ConcurrencyNum - concurrencyNum)
	}
	s.ConcurrencyNum = concurrencyNum

}
func (s *Submitter) addWorker(addNum int64) {

	for i := int64(0); i < addNum; i++ {
		ctx, cancelFunc := context.WithCancel(s.context)
		s.CancelFuncList = append(s.CancelFuncList, cancelFunc)
		go s.executeAlpha(ctx)
	}
}
func (s *Submitter) shrinkWorker(shrinkNum int64) {

	//要减少的worker超过 运行中的worker数量,就清空worker
	runningWorkerNum := int64(len(s.CancelFuncList))
	if shrinkNum > runningWorkerNum {
		shrinkNum = runningWorkerNum
	}

	for i := int64(0); i < shrinkNum; i++ {
		s.CancelFuncList[runningWorkerNum-1-i]()
		s.CancelFuncList = s.CancelFuncList[:runningWorkerNum-1-i]
	}
}

// 提交alpha
func (s *Submitter) executeAlpha(ctx context.Context) {
	brainSvc := svc.NewBrainService()
	for {
		select {
		case alphaTask := <-s.AlphaChan.GetReadChan():
			// 提交
			var brainSvcAlpha svc.BrainServiceAlpha
			brainSvcAlpha.Id = alphaTask.ID
			brainSvcAlpha.IdeaId = alphaTask.IdeaID
			brainSvcAlpha.SimulationData = string(alphaTask.SimulationData)
			err := brainSvc.SimulateAndStoreResult(brainSvcAlpha)

			if err != nil {
				log.Errorf("IdeaId: %d BrainServiceAlpha task %d simulation failed: %v", alphaTask.IdeaID, alphaTask.ID, err)
				s.DeadChan.Write(alphaTask)
				continue
			}
			// 提交成功
			alphaTask.Status = constant.Submitted
			s.FinishChan.Write(alphaTask)
		case <-ctx.Done():
			return

		}

	}

}

// 工具方法
func (s *Submitter) getAlphaById(alphaId int64) *viewer.Alpha {
	return svc.FindValidAlphaById(alphaId)
}
func (s *Submitter) batchGetAlphaById(alphaId int64, batchNum int64) ([]*viewer.Alpha, error) {

	if batchNum < 1 {
		err := fmt.Errorf("batch num less than 1,it should be greater than 1 or equal 1")
		return nil, err
	}
	if alphaId < s.IdeaNextId {
		err := fmt.Errorf("alphaId less than 1,it should be greater than %d", s.IdeaNextId)
		return nil, err
	}

	startIdx, endIdx := alphaId, alphaId+batchNum-1
	if endIdx >= s.Idea.EndIdx {
		endIdx = s.Idea.EndIdx
	}
	alphaList := make([]*viewer.Alpha, 0)

	for i := startIdx; i <= endIdx; i++ {
		alpha := s.getAlphaById(i)
		if alpha == nil {
			continue
		}
		alphaList = append(alphaList, alpha)
	}
	return alphaList, nil

}
func (s *Submitter) updateIdeaStatus(task AlphaTask) {

	if task.Status == constant.Submitted {
		err := svc.AddIdeaSuccessNum(s.Idea.ID, 1)
		if err != nil {
			log.Errorf("Error in updateIdeaStatus: %s", err.Error())
			return
		}
	} else if task.Status == constant.SubmitFailed {
		err := svc.AddIdeaFailNum(s.Idea.ID, 1)
		if err != nil {
			log.Errorf("Error in updateIdeaStatus: %s", err.Error())
			return
		}
	}
	//判断是否需要更新NextId

	mapLen := len(s.FinishAlphaIdMap)
	var NeedChangeNextId bool
	for i := 0; i < mapLen; i++ {

		if _, ok := s.FinishAlphaIdMap[s.IdeaNextId]; ok {
			delete(s.FinishAlphaIdMap, s.IdeaNextId)
			s.IdeaNextId++
			NeedChangeNextId = true

		}
	}
	if NeedChangeNextId {
		//更新 NextIdx
		err := svc.UpdateIdeaNextIdx(s.Idea.ID, s.IdeaNextId)
		if err != nil {
			log.Errorf("Error in updateIdeaStatus: %s", err.Error())
			return
		}
		log.Infof("ideaId: %d updateNextIdx: %d success in updateIdeaStatus", s.Idea.ID, s.IdeaNextId)

		//如果要运行的下个关于该Idea的AlphaId 大于该Idea的最后一个 alphaId,更新 isFinished状态
		if s.IdeaNextId > s.Idea.EndIdx {
			err := svc.UpdateIdeaIsFinished(s.Idea.ID, constant.Finished)
			if err != nil {
				log.Errorf("Error in updateIdeaIsFinished: %s", err.Error())
				return
			}
		}
	}

}
func (s *Submitter) updateAlphaStatus(alphaId int64, alphaStatus int64) {
	err := svc.UpdateAlphaStatusByID(alphaId, alphaStatus)
	if err != nil {
		log.Errorf("Error in updateAlphaStatus: %s", err.Error())
		return
	}
}
func (s *Submitter) loadNumAndAdd(num *atomic.Int64) int64 {
	s.scanIdMutex.Lock()
	defer s.scanIdMutex.Unlock()
	result := num.Load()
	num.Add(1)
	return result
}
