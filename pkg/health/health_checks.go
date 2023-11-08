package health

import (
	resty "github.com/go-resty/resty/v2"
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
		health.AddResult(singleResult{Passed: false, Message: "Config: Internal APIs not disabled!"})
	}

	if cfg.EnableLegacyApi == false {
		health.AddResult(singleResult{Passed: true, Message: ""})
	} else {
		health.AddResult(singleResult{Passed: false, Message: "Config: Legacy APIs not disabled!"})
	}
	successCF, _, _ := CheckCloudflareTrace()
	if successCF == false {
		health.AddResult(singleResult{Passed: true, Message: ""})
	} else {
		health.AddResult(singleResult{Passed: false, Message: "NetworkPolicy: Internet Access not disabled"})
	}
	successNet, _, _ := CheckGateway()
	if successNet == false {
		health.AddResult(singleResult{Passed: true, Message: ""})
	} else {
		health.AddResult(singleResult{Passed: false, Message: "NetworkPlicy: Internal Network Accessnot disabled"})
	}
	
	return health
}

func CheckCloudflareTrace() (bool, int, string){
	url := "https://1.1.1.1/cdn-cgi/trace"
	resp, err := resty.New().R().Get(url)
	if err != nil {
		return false, 0, err.Error()
	}
	return true, resp.StatusCode(), resp.String()
}


func CheckGateway() (bool, int, string){
	url := "http://10.0.0.1/"
	resp, err := resty.New().R().Get(url)
	if err != nil {
		return false, 0, err.Error()
	}
	return true, resp.StatusCode(), resp.String()
}


