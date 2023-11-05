import React, { useEffect } from "react"
import { format } from "util"

export class CheckResult {
  state: string = ""
  name: string = ""
  requeue: boolean = true
  last_run: Date = new Date()
  message: string = ""
  resetFunc = () => { };

  constructor(name: string) {
    this.name = name
    this.state = "na"
  }

  public setResult(state: string | undefined, message: string | undefined, requeue = false): CheckResult {
    if (state) { this.state = state }
    if (message) { this.message = message }
    this.requeue = this.requeue || requeue
    return this
  }
  public setInProgress(): CheckResult {
    this.state = "running"
    return this
  }
  public toString(): string {
    return format("%s / %s / %s", this.state, this.name, this.message)
  }
}

export function useHealthCheck(name: string, params: any, checkFunc: (arg0: string, arg1: any) => Promise<CheckResult>) {
  const [result, setResult] = React.useState(new CheckResult(name).setInProgress())

  useEffect(() => {
    checkFunc(name, params)
      .then(r => { console.log("Prformed HealthCheck: " + r); return r })
      .then(r => setResult(r))
  }, [])

  return result
}



export async function CheckDebuggerFunction(name: string, params: any): Promise<CheckResult> {
  var ret = new CheckResult(name)
  ret.state = params.state
  ret.last_run = new Date()
  ret.requeue = false
  ret.message = params.message


  return await ret
}

interface CheckUrlParams {
  url: string
  expected_status?: number
  host: string | undefined
}
export async function CheckUrlViaBackend(name: string, params: CheckUrlParams): Promise<CheckResult> {
  var ret = new CheckResult(name)
  ret.state = "na"
  ret.last_run = new Date()

  const expected_statuscode = params.expected_status || 200
  const host_header = params.host || "na"

  return fetch("/api/legacy/proxy?" + new URLSearchParams({ url: params.url, host: host_header }), {})
    .then(res => res.status)
    .then(status => {
      if (status == expected_statuscode) {
        return ret.setResult("ok", "Success")
      } else {
        return ret.setResult("fail", format("Status Code: %d. Expected: %d", status, expected_statuscode))
      }
    })
    .catch(err => { return ret.setResult("fail", format("Exception calling \"%s\"", params.url, err)) });
}
interface CheckMetricsParams {
  url: string
  metric: string
  labels: {}
  value: string | undefined
}
export async function CheckMetric(name: string, params: CheckMetricsParams): Promise<CheckResult> {
  var ret = new CheckResult(name)
  ret.state = "na"
  ret.last_run = new Date()
  const metric = params.metric
  const expected_labels: {} = params.labels || {}
  const expected_value = params.value

  const filterByMetricName = (all_metrics: string[]) => { return all_metrics.filter(m => m.startsWith(metric)) }
  const filterByLabels = (metrics: string[]) => {
    Object.entries(expected_labels).forEach((kvpair) => {
      metrics = metrics.filter(m => m.includes(format("%s=\"%s\"", kvpair[0], kvpair[1])))
    })
    return metrics
  }
  const match_value = (metrics: string[]) => {
    if (metrics.length === 0) { return ret.setResult("fail", format("metric not found: %s", metric)) }
    if (params.value === undefined) { return ret.setResult("ok", format("Metric exists: %s", metric)) }

    const value = metrics[0].split(" ")
    if (value.length != 2) { return ret.setResult("fail", "Cannot parse metrics: " + metrics) }

    if (value[1] === expected_value) {
      return ret.setResult("ok", format("%s is %s", metric, value[1]))
    } else { return ret.setResult("rail", format("%s: %s != %s", metric, value[1], expected_value)) }
  }


  return fetch("/api/legacy/metrics?" + new URLSearchParams({ url: params.url }), {})
    .then(res => res.text())
    .then(text => text.split("\n"))
    .then(filterByMetricName)
    .then(filterByLabels)
    // .then(o => { console.log(o); return o })
    .then(match_value)
    .catch(err => { return ret.setResult("fail", format("Error getting metrics endpoint \"%s\"", params.url, err)) });
}
