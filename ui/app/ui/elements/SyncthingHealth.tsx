"use client"
import { Button, Card, CardActionArea, CardActions, CardContent, Stack, Typography } from "@mui/material";
import useSWR from "swr";
import { FallbackCardError, FallbackCardLoading } from "./FallbackCards";
import { ComponentContentValue, ComponentHeader } from "../parts/Component";


export default function SyncthingHealth() {
  const NAME="Syncthing"
  const API_URL="/api/component/syncthing"
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
        {ComponentHeader("Syncthing", "/icons/syncthing-logo-128.png", API_URL, data.status, data.reason)}
        <CardContent>
          <Typography>Devices:</Typography>
          {data.devices.map((dev: any) => {
            return <ComponentContentValue key={dev.DisplayName} label={dev.DisplayName} value={dev.Status} />
          })}
          <Typography>Folders:</Typography>
          {data.folders.map((dir: any) => {
            return <ComponentContentValue key={dir.DisplayName} label={dir.DisplayName} value={dir.Status} />
          })}
        </CardContent>
      </CardActionArea>
      <CardActions>
        <Button href={data.alt_url}>NIP</Button>
      </CardActions>
    </Card>
  );
}


