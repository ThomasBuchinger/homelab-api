"use client"
import { Button, Card, CardActionArea, CardActions, CardContent, Stack, Typography } from "@mui/material";
import useSWR from "swr";
import { FallbackCardError, FallbackCardLoading } from "./FallbackCards";
import { ComponentContentValue, ComponentHeader } from "../parts/Component";
import { format } from "util";



export default function TvHealth() {
  const NAME = "tv.buc.sh (Dummy) "
  const API_URL = "/api/public/bff/nasv3"
  const fetcher = async (url: string) => fetch(url).then(res => res.json())
  const { data, error, isLoading } = useSWR(API_URL, fetcher, { refreshInterval: 10000, fallbackData: { status: "INVALID", reason: "not implemented" } })

  if (isLoading) {
    return FallbackCardLoading(NAME)
  }

  if (error) {
    return FallbackCardError(NAME, API_URL, error)
  }
  
  return (
    <Card>
      <CardActionArea target="_blank" href={"player.10.0.0.24.nip.io"}>
        {ComponentHeader(NAME, "/icons/unraid-logo.png", API_URL, data.status, data.reason)}
        <CardContent>
          <ComponentContentValue label={"Movies / Shows"} value={"42"} />
          <ComponentContentValue label={"Watching"} value={"1"} />
        </CardContent>
      </CardActionArea>
      <CardActions>
        <Stack>
          <Button href="" target="_blank">Torrent</Button>
          <Button href="" target="_blank">UseNet</Button>
          <Button href="" target="_blank">Movies</Button>
          <Button href="" target="_blank">Filme</Button>
          <Button href="" target="_blank">Shows</Button>
          <Button href="" target="_blank">Serien</Button>
          <Button href="" target="_blank">Music</Button>
        </Stack>
      </CardActions>
    </Card>
  );
}


