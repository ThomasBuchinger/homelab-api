"use client"
import { Button, Card, CardActionArea, CardActions, CardContent, Stack, Typography } from "@mui/material";
import useSWR from "swr";
import { FallbackCardError, FallbackCardLoading } from "./FallbackCards";
import { ComponentContentValue, ComponentHeader } from "../parts/Component";
import { duration } from "@mui/material";



export default function TvHealth() {
  const NAME = "tv.buc.sh"
  const API_URL = "/api/public/bff/jellyfin"
  const fetcher = async (url: string) => fetch(url).then(res => res.json())
  const { data, error, isLoading } = useSWR(API_URL, fetcher, { refreshInterval: 10000, fallbackData: { status: "INVALID", reason: "not implemented" } })

  if (isLoading) {
    return FallbackCardLoading(NAME)
  }

  if (error) {
    return FallbackCardError(NAME, API_URL, error)
  }

  var jellyfin_media = parseInt(data.jellyfin_movies, 10) + parseInt(data.jellyfin_shows, 10);
  
  return (
    <Card>
      <CardActionArea target="_blank" href={data.jellyfin_url}>
        {ComponentHeader(NAME, "/icons/jellyfin-logo.svg", API_URL, data.status, data.reason)}
      </CardActionArea>
      <CardContent>
          <ComponentContentValue label={"Watchers / Total"} value={data.jellyfin_watching + " / " + jellyfin_media} />
        <CardActionArea target="_blank" href={data.movies_url}>
          <ComponentContentValue label={"Movies"} value={"100/100"} />
        </CardActionArea>
        <CardActionArea target="_blank" href={data.shows_url}>
          <ComponentContentValue label={"Shows"} value={"100/100"} />
        </CardActionArea>
        <CardActionArea target="_blank" href={data.downloads_url}>
          <ComponentContentValue label={"Downloads"} value={data.downloads_running + " / " + secondsToStr(100990)} />
          {/* data.downloads_estimated_finish */}
        </CardActionArea>
      </CardContent>
      <CardActions >
        <Button href={data.jellyfin_url}>Player</Button>
        <Button href={data.jellyseer_url}>Requests</Button>
      </CardActions>
    </Card>
  );
}

function secondsToStr(input: number) {
  const years = Math.floor(input / 31536000),
    days = Math.floor((input %= 31536000) / 86400),
    hours = Math.floor((input %= 86400) / 3600),
    minutes = Math.floor((input %= 3600) / 60),
    seconds = input % 60;

  if (days || hours || seconds || minutes) {
    return (years ? years + "y " : "") +
      (days ? days + "d " : "") +
      (hours ? hours + "h " : "") +
      (minutes ? minutes + "m " : "") +
      (seconds ? seconds + "s " : "");
  }

  return "< 1s";
}