package common

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
