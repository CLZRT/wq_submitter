package svc

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
	"wq_submitter/configs"
	"wq_submitter/internal/auth"
	"wq_submitter/internal/constant"
	"wq_submitter/internal/viewer"
)

var conf *configs.GlobalConfig

func init() {
	conf = configs.GetGlobalConfig()
}

type BrainServiceAlpha struct {
	Id             int64
	IdeaId         int64
	SimulationData string
}
type BrainServiceRespContainer struct {
	resp http.Response
}
type BrainServiceRetryResp struct {
	Id       string `json:"id"`
	Type     string `json:"type"`
	Settings struct {
		InstrumentType string  `json:"instrumentType"`
		Region         string  `json:"region"`
		Universe       string  `json:"universe"`
		Delay          int     `json:"delay"`
		Decay          int     `json:"decay"`
		Neutralization string  `json:"neutralization"`
		Truncation     float64 `json:"truncation"`
		Pasteurization string  `json:"pasteurization"`
		UnitHandling   string  `json:"unitHandling"`
		NanHandling    string  `json:"nanHandling"`
		MaxTrade       string  `json:"maxTrade"`
		Language       string  `json:"language"`
		Visualization  bool    `json:"visualization"`
	} `json:"settings"`
	Regular string `json:"regular"`
	Status  string `json:"status"`
	Alpha   string `json:"alpha"`
	Message string `json:"message"`
}

type BrainService struct {
	brainAuth *auth.BrainAuth
}

func NewBrainService() *BrainService {
	brainAuth := auth.GetBrainAuth()
	return &BrainService{
		brainAuth: brainAuth,
	}
}

func (brainSvc *BrainService) SimulateAndStoreResult(alpha BrainServiceAlpha) error {
	//模拟
	simulateResp, err := brainSvc.simulate(alpha.SimulationData)

	if err != nil {
		log.Errorf("simulate Failed {%s}", err.Error())
		return err
	}

	//结果

	if conf.ResultConfig.NeedStoreRecords {

		// 获取results
		alphaResultViewer, err := brainSvc.getResults(alpha, simulateResp)
		if err != nil {
			log.Errorf("GetResult failed: %s", err.Error())
			return err
		}

		// 存Results
		if err = StoreAlphaResult(alphaResultViewer); err != nil {
			log.Errorf("StoreAlphaResult failed: %s", err.Error())
			return err
		}
		//	确认提交成功
	} else if err = brainSvc.confirmResult(simulateResp); err != nil {
		log.Errorf("ConfirmResult failed: %s", err.Error())
		return err
	}
	return nil

}
func (brainSvc *BrainService) simulate(alphaDataStr string) (*http.Response, error) {
	defer brainSvc.brainAuth.CheckFreshToken()

	// 构建请求
	req, err := http.NewRequest("POST", "https://api.worldquantbrain.com/simulations",
		strings.NewReader(alphaDataStr))
	if err != nil {
		log.Errorf("New Simulations Request failed: %s", err.Error())
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	// 发送请求 (httpClient 会自动附加 cookie)
	resp, err := brainSvc.brainAuth.HttpClient.Do(req)
	defer func(resp *http.Response) {
		if resp != nil {
			resp.Body.Close()
		}
	}(resp)
	if err != nil {
		log.Errorf("simulate Request failed: %s", err.Error())
		return nil, err
	}

	// 打印响应状态码
	log.Infof("simulate 响应状态码: %d\n", resp.StatusCode)

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("simulate 读取响应体失败: %s", err.Error())

	}

	// 检查请求结果
	if resp.StatusCode >= 400 {
		err := fmt.Errorf("simulate code: %d, message: %s", resp.StatusCode, string(body))
		return nil, err
	}

	// 打印响应内容
	log.Infof("simulate 响应内容: %s\n", string(body))
	return resp, nil

}

func (brainSvc *BrainService) confirmResult(resp *http.Response) error {
	//构建请求
	req, err := http.NewRequest("GET", resp.Header.Get("Location"), strings.NewReader(""))
	if err != nil {
		log.Errorf("New ConfirmResult Request failed: %s", err.Error())
		return err
	}
	// 获取重试时间
	retrySecond, err := strconv.ParseFloat(resp.Header.Get("Retry-After"), 64)
	if err != nil {
		log.Errorf("Parse Retry-After Failed: %s", err.Error())
		return err
	}

	// 发送请求
	_, err = brainSvc.retryGetBasic(req, retrySecond)
	if err != nil {
		log.Errorf("ConfirmResult failed with message: %s", err.Error())
		return err
	}
	return nil

}

// 获取Simulate 结果
func (brainSvc *BrainService) getResults(alpha BrainServiceAlpha, resp *http.Response) (alphaResultViewer *viewer.AlphaResult, err error) {

	//构建请求
	req, err := http.NewRequest("GET", resp.Header.Get("Location"), strings.NewReader(""))
	if err != nil {
		log.Errorf("New GetResult Request failed: %s", err.Error())
		return nil, err
	}
	retrySecond, err := strconv.ParseFloat(resp.Header.Get("Retry-After"), 64)
	if err != nil {
		log.Errorf("Parse Retry-After Failed: %s", err.Error())
		return nil, err
	}

	// 发送请求 (httpClient 会自动附加 cookie)
	originalResult, err := brainSvc.retryGetBasic(req, retrySecond)

	if err != nil {
		log.Errorf("GetResult Request failed: %s", err.Error())
		return nil, err
	}
	alphaCode := originalResult.Alpha

	// 检查返回结果
	if originalResult.Status != constant.StatusComplete {
		err := fmt.Errorf("get result failed %s, message: %s", originalResult.Status, originalResult.Message)
		log.Error(err.Error())
		return nil, err
	}

	// 获取模拟结果
	resultMap, err := brainSvc.getTotalResult(alphaCode)
	if err != nil {
		log.Errorf("getTotalResult failed: %s", err.Error())
		return nil, err
	}
	alphaResultViewer = &viewer.AlphaResult{}
	alphaResultViewer.AlphaId = alpha.Id
	alphaResultViewer.AlphaCode = alphaCode
	alphaResultViewer.IdeaId = alpha.IdeaId
	alphaResultViewer.AlphaDetail = []byte(alpha.SimulationData)

	if basicResult, ok := resultMap[constant.BasicMapKey]; ok {
		alphaResultViewer.BasicResult = *basicResult
	}
	if checkResult, ok := resultMap[constant.CheckMapKey]; ok {
		alphaResultViewer.CheckResult = *checkResult
	}
	if selfCorrelationResult, ok := resultMap[constant.SelfCorrelationMapKey]; ok {
		alphaResultViewer.SelfCorrelation = *selfCorrelationResult
	}
	if prodCorrelationResult, ok := resultMap[constant.ProdCorrelationMapKey]; ok {
		alphaResultViewer.ProdCorrelation = *prodCorrelationResult
	}

	if turnoverResult, ok := resultMap[constant.TurnoverMapKey]; ok {
		alphaResultViewer.Turnover = *turnoverResult
	}
	if sharpeResult, ok := resultMap[constant.SharpeMapKey]; ok {
		alphaResultViewer.Sharpe = *sharpeResult
	}
	if pnlResult, ok := resultMap[constant.PnlMapKey]; ok {
		alphaResultViewer.Pnl = *pnlResult
	}
	if dailyPnlResult, ok := resultMap[constant.DailyPnlMapKey]; ok {
		alphaResultViewer.DailyPnl = *dailyPnlResult
	}
	if yearlyStatsResult, ok := resultMap[constant.YearlyStatsMapKey]; ok {
		alphaResultViewer.YearlyStats = *yearlyStatsResult
	}

	return alphaResultViewer, nil

}
func (brainSvc *BrainService) getTotalResult(alphaCode string) (map[string]*[]byte, error) {
	basicResult, err := brainSvc.getBasicResult(alphaCode)
	if err != nil {
		log.Errorf("getBasicResult failed: %s", err.Error())
		return nil, err
	}
	resultMap := brainSvc.getDetailResultList(alphaCode)
	resultMap[constant.BasicMapKey] = &basicResult
	return resultMap, nil

}
func (brainSvc *BrainService) getBasicResult(alphaCode string) ([]byte, error) {
	maxTimes := conf.ResultConfig.MaxRetryNum
	var resp *http.Response
	var err error
	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()
	var result []byte
	for i := int64(0); i < maxTimes; i++ {
		resp, err = brainSvc.brainAuth.HttpClient.Get(constant.ResultBasicUrl + "/" + alphaCode)
		if err != nil {
			log.Errorf("getBasicResult Request failed: %s", err.Error())
			return nil, err
		}
		// 读取响应体
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Errorf("读取响应体失败: %s", err.Error())
			return nil, err
		}

		if string(body) != "" {
			result = body
			break
		}
	}

	return result, nil

}
func (brainSvc *BrainService) getDetailResultList(alphaCode string) map[string]*[]byte {
	//结果 容器
	detailMap := make(map[string]*[]byte)

	// Turnover
	turnoverResult := brainSvc.getDetailResult(alphaCode, constant.TurnoverRecordUri)
	if turnoverResult != nil {
		detailMap[constant.TurnoverMapKey] = turnoverResult
	}

	// Sharpe
	sharpeResult := brainSvc.getDetailResult(alphaCode, constant.SharpeRecordUri)
	if sharpeResult != nil {
		detailMap[constant.SharpeMapKey] = sharpeResult
	}

	// Pnl
	pnlResult := brainSvc.getDetailResult(alphaCode, constant.PnlRecordUri)
	if sharpeResult != nil {
		detailMap[constant.PnlMapKey] = pnlResult
	}

	// DailyPnl
	dailyPnlResult := brainSvc.getDetailResult(alphaCode, constant.DailyPnlRecordUri)
	if dailyPnlResult != nil {
		detailMap[constant.DailyPnlMapKey] = dailyPnlResult
	}

	// YearlyStats
	yearlyStatsResult := brainSvc.getDetailResult(alphaCode, constant.YearlyStatsRecordUri)
	if dailyPnlResult != nil {
		detailMap[constant.YearlyStatsMapKey] = yearlyStatsResult
	}

	// selfCorrelation 自相关性检查结果
	selfCorrelationResult := brainSvc.getDetailResult(alphaCode, constant.SelfCorrelationUri)
	if selfCorrelationResult != nil {
		detailMap[constant.SelfCorrelationMapKey] = selfCorrelationResult
	}

	// 生产相关性检查结果
	prodCorrelationResult := brainSvc.getDetailResult(alphaCode, constant.ProdCorrelationUri)
	if prodCorrelationResult != nil {
		detailMap[constant.ProdCorrelationMapKey] = prodCorrelationResult
	}

	// 测试通过检查结果
	checkResult := brainSvc.getDetailResult(alphaCode, constant.CheckUri)
	if checkResult != nil {
		detailMap[constant.CheckMapKey] = checkResult
	}

	return detailMap
}

func (brainSvc *BrainService) retryGetBasic(req *http.Request, retrySecond float64) (result *BrainServiceRetryResp, err error) {

	retryTicker := time.NewTicker(time.Duration(retrySecond) * time.Second)
	defer retryTicker.Stop()
	//到时间就重试
	var resp *http.Response
	defer func() {
		if resp != nil {
			resp.Body.Close()

		}
		if err != nil {
			time.Sleep(10 * time.Second)
		}
	}()
	for range retryTicker.C {
		resp, err = brainSvc.brainAuth.HttpClient.Do(req)
		//请求建立过程有问题
		if err != nil {
			log.Errorf("retryGetBasic Request failed: %s", err.Error())
			return nil, err
		}

		if resp.StatusCode >= 500 {
			err = fmt.Errorf("retryGetBasic Request failed %d", resp.StatusCode)
			log.Error(err.Error())
			return nil, err
		}

		//请求建立没问题,请求构建或者服务端有问题
		if resp.StatusCode >= 400 {
			err = fmt.Errorf("retryGetBasic Request failed %d", resp.StatusCode)
			log.Warn(err.Error())
			return nil, err
		}

		// 状态码[200,400),读取响应结构体
		bytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Warnf("retryGetBasic 读取响应体失败: %s", err.Error())
		}
		//解析结构体
		var respData BrainServiceRetryResp
		err = json.Unmarshal(bytes, &respData)
		if err != nil {
			log.Errorf("retryGetBasic 解析Resp Body失败: %s", err.Error())
		}

		//成功获取到状态，请求没问题
		if respData.Status != "" {

			if respData.Status == constant.StatusComplete {
				return &respData, nil
			}

			if respData.Status == constant.StatusError {
				err = fmt.Errorf("status: %s,Message:%s", respData.Status, respData.Message)
				log.Error(err.Error())
				return nil, err
			}
			if respData.Status == constant.StatusWarning {
				log.Warnf("status: %s,Message:%s", respData.Status, respData.Message)
			}

		}

	}
	return nil, fmt.Errorf("can't Get Result")
}
func (brainSvc *BrainService) getDetailResult(alphaCode string, detailName string) (result *[]byte) {
	maxTimes := conf.ResultConfig.MaxRetryNum
	var resp *http.Response
	var err error
	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()
	for i := int64(0); i < maxTimes; i++ {
		resp, err = brainSvc.brainAuth.HttpClient.Get(constant.ResultBasicUrl + "/" + alphaCode + detailName)
		if err != nil {
			log.Errorf("getDetailResult Request failed: %s", err.Error())
			return nil
		}

		// 读取响应体
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Errorf("读取响应体失败: %s", err.Error())
			return nil
		}
		bodyStr := string(body)
		result = &body
		if bodyStr != "" {
			break
		}
	}

	return result
}
