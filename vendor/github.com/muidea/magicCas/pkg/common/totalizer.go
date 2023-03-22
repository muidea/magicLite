package common

import (
	"time"

	"github.com/muidea/magicCommon/event"

	"github.com/muidea/magicCas/pkg/model"
)

const TotalizerModule = "/kernel/totalizer"

const (
	QueryTotalizer        = "/kernel/totalizer/query/"
	QuerySummaryTotalizer = "/kernel/totalizer/summary/query/"
	CreateTotalizer       = "/kernel/totalizer/create/"
	DeleteTotalizer       = "/kernel/totalizer/delete/"
)

const (
	TotalizeRealtime = 1
	TotalizeWeek     = 2
	TotalizeMonth    = 3
	TotalizeDaily    = 4
)

type Totalizer interface {
	Owner() string
	Trigger(event event.Event) (err error)
	Period(typeVal int, timeStamp time.Time) (err error)
	Get(typeVal int) (ret *model.Totalizer, err error)
	Same(owner, trigger string) bool
}

type TotalizeParam struct {
	Owner   string `json:"owner"`
	Trigger string `json:"trigger"`
	Period  []int  `json:"period"`
}

type TotalizerView struct {
	ID        int     `json:"id"`
	Owner     string  `json:"owner"`
	Type      int     `json:"type"`
	TimeStamp int64   `json:"timeStamp"`
	Value     float64 `json:"value"`
	Catalog   int     `json:"catalog"`
}

func (s *TotalizerView) FromTotalizer(ptr *model.Totalizer) {
	s.ID = ptr.ID
	s.Owner = ptr.Owner
	s.Type = ptr.Type
	s.TimeStamp = ptr.TimeStamp
	s.Value = ptr.Value
	s.Catalog = ptr.Catalog
}

type PeriodParam struct {
	Period int `json:"period"`
	Count  int `json:"count"`
}

type QuerySummaryParam struct {
	Owner  string         `json:"owner"`
	Period []*PeriodParam `json:"period"`
}

type PeriodSummary struct {
	Period  int                `json:"period"`
	Current *model.Totalizer   `json:"current"`
	Trend   []*model.Totalizer `json:"trend"`
}

type QuerySummaryResult struct {
	Owner   string           `json:"owner"`
	Summary []*PeriodSummary `json:"summary"`
}

type PeriodSummaryView struct {
	Period  int              `json:"period"`
	Current *TotalizerView   `json:"current"`
	Trend   []*TotalizerView `json:"trend"`
}

func (s *PeriodSummaryView) FromPeriodSummary(ptr *PeriodSummary) {
	if ptr == nil {
		return
	}

	s.Period = ptr.Period
	if ptr.Current != nil {
		s.Current = &TotalizerView{}
		s.Current.FromTotalizer(ptr.Current)
	}
	for _, val := range ptr.Trend {
		view := &TotalizerView{}
		view.FromTotalizer(val)
		s.Trend = append(s.Trend, view)
	}
}
