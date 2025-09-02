package viewer

import (
	"encoding/json"
)

type AlphaResult struct {
	IdeaId      int64           //  table:idea id
	AlphaId     int64           //  关联的alpha id
	AlphaDetail json.RawMessage //  alpha表达式及其环境
	AlphaCode   string

	//  alpha代码,brain平台唯一标识一个回测过的alpha
	BasicResult     json.RawMessage //  基本测试结果
	CheckResult     json.RawMessage //  检查结果
	SelfCorrelation json.RawMessage //  自相关性结果
	ProdCorrelation json.RawMessage
	Turnover        json.RawMessage //  turnover 详细数据
	Sharpe          json.RawMessage //  sharpe 详细数据
	Pnl             json.RawMessage //  pnl 详细数据
	DailyPnl        json.RawMessage //  daily_pnl 详细数据
	YearlyStats     json.RawMessage //  yearly_stats 详细数据
}
