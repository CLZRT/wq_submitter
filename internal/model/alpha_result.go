package model

import "gorm.io/datatypes"

// AlphaResult corresponds to the `alpha_result` table in the database.
type AlphaResult struct {
	ID          int64          `gorm:"column:id" db:"id" json:"id"`
	IdeaId      int64          `gorm:"column:idea_id" db:"idea_id" json:"idea_id"`                                   //  table:idea id
	AlphaId     int64          `gorm:"column:alpha_id" db:"alpha_id" json:"alpha_id"`                                //  关联的alpha id
	AlphaDetail datatypes.JSON `gorm:"column:alpha_detail;type:json not null" db:"alpha_detail" json:"alpha_detail"` //  alpha表达式及其环境
	AlphaCode   string         `gorm:"column:alpha_code" db:"alpha_code" json:"alpha_code"`

	//  alpha代码,brain平台唯一标识一个回测过的alpha
	BasicResult     datatypes.JSON `gorm:"column:basic_result;type:json not null" db:"basic_result" json:"basic_result"`    //  基本测试结果
	CheckResult     datatypes.JSON `gorm:"column:check_result;type:json" db:"check_result" json:"check_result"`             //  检查结果
	SelfCorrelation datatypes.JSON `gorm:"column:self_correlation;type:json" db:"self_correlation" json:"self_correlation"` //  自相关性结果
	ProdCorrelation datatypes.JSON `gorm:"column:prod_correlation;type:json" db:"prod_correlation" json:"prod_correlation"` //  自相关性结果
	Turnover        datatypes.JSON `gorm:"column:turnover;type:json" db:"turnover" json:"turnover"`                         //  turnover 详细数据
	Sharpe          datatypes.JSON `gorm:"column:sharpe;type:json" db:"sharpe" json:"sharpe"`                               //  sharpe 详细数据
	Pnl             datatypes.JSON `gorm:"column:pnl;type:json" db:"pnl" json:"pnl"`                                        //  pnl 详细数据
	DailyPnl        datatypes.JSON `gorm:"column:daily_pnl;type:json" db:"daily_pnl" json:"daily_pnl"`                      //  daily_pnl 详细数据
	YearlyStats     datatypes.JSON `gorm:"column:yearly_stats;type:json" db:"yearly_stats" json:"yearly_stats"`             //  yearly_stats 详细数据
}

func (AlphaResult) TableName() string {
	return "alpha_result"
}
