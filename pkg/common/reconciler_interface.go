package common

import "strings"

type Reconciler interface {
	GetStatus() ReconcilerStatus
	GetReason() string
	Reconcile() (bool, int)
	// GetData() interface{}
}

type ReconcilerStatus string

const (
	ReconcilerStatusNew           ReconcilerStatus = "NEW"
	ReconcilerStatusConfigInvalid ReconcilerStatus = "INVALID"
	ReconcilerStatusOK            ReconcilerStatus = "OK"
	ReconcilerStatusDown          ReconcilerStatus = "DOWN"
	ReconcilerStatusError         ReconcilerStatus = "ERROR"
	ReconcilerStatusWarn          ReconcilerStatus = "WARN"
	ReconcilerStatusInfo          ReconcilerStatus = "INFO"
)

func MergeScrapersStatus(scrapers ...Reconciler) (ReconcilerStatus, string) {
	severity := map[ReconcilerStatus]int{
		ReconcilerStatusOK:            0,
		ReconcilerStatusNew:           1,
		ReconcilerStatusInfo:          2,
		ReconcilerStatusWarn:          3,
		ReconcilerStatusConfigInvalid: 4,
		ReconcilerStatusDown:          5,
		ReconcilerStatusError:         6,
	}

	if len(scrapers) == 0 {
		return ReconcilerStatusNew, ""
	}

	merged := scrapers[0].GetStatus()
	var problems []string
	for _, s := range scrapers {
		st := s.GetStatus()
		if severity[st] > severity[merged] {
			merged = st
		}
		if st != ReconcilerStatusOK && st != ReconcilerStatusNew {
			problems = append(problems, s.GetReason())
		}
	}

	reason := "OK"
	if len(problems) > 0 {
		reason = strings.Join(problems, "; ")
	}
	return merged, reason
}
