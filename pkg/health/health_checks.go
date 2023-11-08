package health

import (
	"github.com/thomasbuchinger/homelab-api/pkg/common"
)

type ExternalHealthCheckResult struct{
	Health bool
	PassedChecks, TotalChecks int
	Results []singleResult
}
func (healthcheck *ExternalHealthCheckResult) AddResult(result singleResult) {
	healthcheck.Results = append(healthcheck.Results, result)

	// recalculate overall status
	h := true
	passed_count := 0
	for _, r := range healthcheck.Results {
		h = h && r.Passed
		if r.Passed {
			passed_count +=1
		}
	}
	healthcheck.Health = h
	healthcheck.PassedChecks = passed_count
	healthcheck.TotalChecks = len(healthcheck.Results)
}
type singleResult struct{
	Passed bool
	Message string
}

func Ok() ExternalHealthCheckResult {
	return ExternalHealthCheckResult{
		Health: true,
		Results: []singleResult{
			{Passed: true, Message: "No Issues! OK"},
		},
	}
}

func CheckApiPublic() ExternalHealthCheckResult {
	cfg := common.GetServerConfig()
	health := ExternalHealthCheckResult{}

	if cfg.EnableInternalApis == false {
		health.AddResult(singleResult{Passed: true, Message: ""})
	} else {
		health.AddResult(singleResult{Passed: false, Message: "Internal APIs are enabled!"})
	}

	if cfg.EnableLegacyApi == false {
		health.AddResult(singleResult{Passed: true, Message: ""})
	} else {
		health.AddResult(singleResult{Passed: false, Message: "Legacy APIs are enabled"})
	}
	
	return health
}

