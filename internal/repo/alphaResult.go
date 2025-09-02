package repo

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"wq_submitter/internal/model"
)

type AlphaResultRepo struct {
}

func NewAlphaResultRepo() *AlphaResultRepo {
	return &AlphaResultRepo{}
}
func (resultRepo *AlphaResultRepo) Add(_ context.Context, alphaResult *model.AlphaResult) (int64, error) {

	result := db.Create(alphaResult)
	if result.Error != nil {
		log.Errorf("failed to add alpha result for alphaId %d: %v", alphaResult.AlphaId, result.Error)
		return result.RowsAffected, result.Error
	}
	if result.RowsAffected == 0 {
		err := fmt.Errorf("add AlphaResult for alphaId %d: no rows affected", alphaResult.AlphaId)
		log.Error(err)
		return result.RowsAffected, err
	}
	return result.RowsAffected, nil

}
