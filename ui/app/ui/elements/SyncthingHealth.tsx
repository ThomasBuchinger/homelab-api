"use client"
import { Button, Card, CardActionArea, CardActions, CardContent, Stack, Typography } from "@mui/material";
import useSWR from "swr";
import { FallbackCardError, FallbackCardLoading } from "./FallbackCards";
import { ComponentContentValue, ComponentHeader } from "../parts/Component";
import axios, { AxiosError } from "axios";
import { Bounce, toast, ToastContainer } from "react-toastify";
import { format } from "util";


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

  type syncthingObject = { display_name: string, status: string }
  const sortfunc = (a: syncthingObject, b: syncthingObject) => { return a.display_name < b.display_name ? 1 : -1}


  const handleRestart = async (e: any) => {
    e.preventDefault();

    try {
      const res = await axios.delete("/api/component/syncthing/restart")
      console.log(res.data);
      toast.success("Restarted")
    } catch (err: any) {
      console.log(err);
      toast.error(format("Failed to Restart Syncthing Pod: %s / %s", err.response.data.status || "CALL_FAILED", err.response.data.reason || "") )
    }
  };

  return (
    <Card>
      <CardActionArea href={data.url} target="_blank">
        {ComponentHeader("Syncthing", "/icons/syncthing-logo-128.png", API_URL, data.status, data.reason)}
        <CardContent>
          <Typography>Devices:</Typography>
          {data.devices.sort(sortfunc).map((dev: syncthingObject) => {
            return <ComponentContentValue key={dev.display_name} label={dev.display_name} value={dev.status} />
          })}
          <Typography>Folders:</Typography>
          {data.folders.sort(sortfunc).map((dir: syncthingObject) => {
            return <ComponentContentValue key={dir.display_name} label={dir.display_name} value={dir.status} />
          })}
        </CardContent>
      </CardActionArea>
      <CardActions>
        <Button onClick={handleRestart} >Restart</Button>
      </CardActions>

    </Card>
  );
}


