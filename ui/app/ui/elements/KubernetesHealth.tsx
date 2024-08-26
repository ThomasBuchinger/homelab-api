"use client"
import { Button, Card, CardActionArea, CardActions, CardContent, Stack, Typography } from "@mui/material";
import useSWR from "swr";
import { FallbackCardError, FallbackCardLoading } from "./FallbackCards";
import { ComponentContentValue, ComponentHeader } from "../parts/Component";


export default function KubernetesHealth() {
  const NAME="Kubernetes"
  const API_URL="/api/component/kubernetes"
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
      <CardActionArea href={data.url} target="_blank">
        {ComponentHeader("Kubernetes", "/icons/kubernetes-icon-color.svg", API_URL, data.status, data.reason)}
        <CardContent>
          <ComponentContentValue label={"Healthy Pods"} value={data.pod_healthy + " / " + data.pod_total} />
          <ComponentContentValue label={"Healthy PVCs"} value={data.pvc_healthy + " / " + data.pvc_total} />
        </CardContent>
      </CardActionArea>
      <CardActions>
        <Button href={data.evergreen_url}>Evergreen</Button>
        <Button href={data.prod_url}>Prod</Button>
      </CardActions>
    </Card>
  );
}


