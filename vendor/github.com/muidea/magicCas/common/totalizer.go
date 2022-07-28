package common

import "github.com/muidea/magicCas/model"

const (
	CreateTotalizer = "/base/totalizer/create/"
	DeleteTotalizer = "/base/totalizer/delete/:id"
	UpdateTotalizer = "/base/totalizer/update/:id"
	QueryTotalizer  = "/base/totalizer/query/"
)

const TotalizerModule = "/kernel/totalizer"

type TotalizerView struct {
	ID        int     `json:"id"`
	Owner     string  `json:"owner"`
	Type      int     `json:"type"`
	TimeStamp int64   `json:"timeStamp"`
	Value     float64 `json:"value"`
	Catalog   int     `json:"catalog"`
}

func (s *TotalizerView) FromTotalizer(ptr *model.Totalizer) {
	if ptr == nil {
		return
	}
	s.ID = ptr.ID
	s.Owner = ptr.Owner
	s.Type = ptr.Type
	s.TimeStamp = ptr.TimeStamp
	s.Value = ptr.Value
	s.Catalog = ptr.Catalog
}
