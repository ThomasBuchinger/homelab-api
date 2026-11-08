"use client"
import { Button, Card, CardActionArea, CardActions, CardContent, Stack, Typography } from "@mui/material";
import useSWR from "swr";
import { FallbackCardError, FallbackCardLoading } from "./FallbackCards";
import { ComponentCamera, ComponentContentValue, ComponentDoorbell, ComponentHeader } from "../parts/Component";
import { format } from "util";



export default function SmarthomeHealth() {
  const NAME = "Smart Home"
  const API_URL = "/api/public/bff/smarthome"
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
      <CardActionArea target="_blank" href={data.url}>
        {ComponentHeader(NAME, "/icons/New_Home_Assistant_logo.svg", API_URL, data.status, data.reason)}
        <CardContent>
          <ComponentContentValue label={"NG7a_IOT"} value={data.accesspoint_up == 1 ? "up" : "down"} />
          <ComponentContentValue label={"Expected Devices"} value={"2"} />
          <ComponentContentValue label={"Devices (Wifi / DHCP)"} value={data.devices_wifi_active + " / " + data.devices_dhcp_active} />
          <ComponentCamera label="Carport" wifi_rating={data.wifi_signal_carportcam} fps={data.frigate_camera_fps_carport} ha_status={data.homeassistant_status_carport_camera} />
          <ComponentDoorbell label="Doorbell NG7" wifi_rating={data.wifi_signal_ng7_doorbell_shelly} ha_status={data.homeassistant_status_doorbell_ng7_shelly}/>
        </CardContent>
      </CardActionArea>
      <CardActions>
        <Button href={data.backup_url} target="_blank">OpenWrt</Button>
      </CardActions>
    </Card>
  );
}


