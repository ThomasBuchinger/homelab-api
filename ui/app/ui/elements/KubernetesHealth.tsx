"use client"
import { Button, Card, CardActionArea, CardActions, CardContent, Stack, Typography } from "@mui/material";
import useSWR from "swr";
import { FallbackCardError, FallbackCardLoading } from "./FallbackCards";
import { ComponentContentValue, ComponentHeader } from "../parts/Component";


export default function KubernetesHealth() {
  const NAME="Kubernetes"
  const API_URL="/api/public/bff/kubernetes"
  const fetcher = async (url: string) => fetch(url).then(res => res.json())
  const { data, error, isLoading } = useSWR(API_URL, fetcher, { refreshInterval: 10000, fallbackData: { status: "INVALID", messages: [] } })

  if (isLoading) {
    return FallbackCardLoading(NAME)
  }

  if (error) {
    return FallbackCardError(NAME, API_URL, error)
  }

  return (
    <Card>
      {ComponentHeader("Kubernetes", "/icons/kubernetes-icon-color.svg", API_URL, data.status, data.reason)}
      <CardContent>
        Evergreen
        <ComponentContentValue label={"Healthy Pods"} value={data.evergreen.pod_healthy + " / " + data.evergreen.pod_total} />
        <ComponentContentValue label={"Healthy PVCs"} value={data.evergreen.pvc_healthy + " / " + data.evergreen.pvc_total} />
      </CardContent>
      <CardContent>
        Prod
        <ComponentContentValue label={"Healthy Pods"} value={data.prod.pod_healthy + " / " + data.prod.pod_total} />
        <ComponentContentValue label={"Healthy PVCs"} value={data.prod.pvc_healthy + " / " + data.prod.pvc_total} />
      </CardContent>
      <CardActions>
        <Button href={data.evergreen_url} target="_blank">Evergreen</Button>
        <Button href={data.prod_url} target="_blank">Prod</Button>
      </CardActions>
    </Card>
  );
}


