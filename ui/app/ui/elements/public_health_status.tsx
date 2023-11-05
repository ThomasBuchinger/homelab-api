"use client"
import { Circle, Close } from "@mui/icons-material";
import { Card, CardContent, CardHeader, CircularProgress, Skeleton, Typography } from "@mui/material";
import useSWR from "swr";

export function HealthStatusPublic({target}: {target: string}) {
  const fetcher = async (url: string) => fetch(url).then(res => res.json())
  const { data, error, isLoading } = useSWR("/api/public/health?target="+target, fetcher, {refreshInterval: 10000, fallbackData: {healthy: false, messages: []}})

  var status = <Circle color={ data.healthy ? "success" : "error"} />
  if (error) { status = <Close htmlColor='red' />  }
  if (isLoading) { status = <CircularProgress variant='indeterminate' size={15} /> }

  return <Card >
    <CardHeader title={target} titleTypographyProps={{ variant: "h6" }}  avatar={status} />
    {isLoading ?
      <CardContent><Skeleton variant="rectangular" animation={"wave"} /></CardContent> :
      data.messages.map((m: any) => <CardContent key={m}>{m}</CardContent>)
     }
  </Card>
}