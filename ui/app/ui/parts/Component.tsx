"use client"
import { Circle } from "@mui/icons-material";
import { Avatar, CardHeader, Icon, IconButton, Stack, Typography } from "@mui/material";
import Image from "next/image";

export function ComponentHeader(title: string, image_url: string, api_url: string, status: string, reason: string, width = 32, height = 32) {
  if (status == null || status === "") {
    status = "INVALID"
  }
  type Mapping = Record<string, string>
  const ICON_MAPPING: Mapping = {
    "OK": "green",
    "WARN": "yellow",
    "INVALID": "blue",
    "DOWN": "red",
  }
  const status_icon = <Circle htmlColor={ICON_MAPPING[status]} />


  return <CardHeader
    avatar=<Avatar sx={{ bgcolor: "white" }}><Image width={width} height={height} src={image_url} alt={title}></Image></Avatar>
    title={title}
    titleTypographyProps={{ variant: "h5", marginX: "10px" }}
    action={
      <IconButton href={api_url} target="_blank" title={reason || ""}>{status_icon}</IconButton>
    }
  />
}


export function ComponentContentValue({label, value}: {label: string, value: string}) {
  return <Stack justifyContent={"space-between"} direction={"row"}>
    <Typography component={"span"} variant="overline" align="left">{label}:</Typography>
    <Typography component={"span"} variant="overline" align="right" marginLeft={"10px"}>{value}</Typography>
  </Stack>
}

export function ComponentCamera({ label, wifi_rating, fps, ha_status}: { label: string, wifi_rating: string, fps: string, ha_status: number }) {
  return <Stack justifyContent={"space-between"} direction={"row"}>
    <span>
      <img src={"/icons/camera.svg"} height={18} width={24}></img>
      <Typography component={"span"} variant="overline" align="left" marginLeft={"10px"}>{label}</Typography>
    </span>
    <Typography component={"span"} variant="overline" align="right" marginLeft={"10px"}>
      <Typography component={"span"} variant="overline" align="right" sx={{ verticalAlign: "top" }}>{fps} FPS  </Typography>
      <img src={"/icons/wifi_" + wifi_rating + ".svg"} alt={"strength: " + wifi_rating} width={24} ></img>
      <ComponentHomeAssistantStatus ha_status={ha_status} />
    </Typography>
  </Stack>
}
export function ComponentDoorbell({ label, wifi_rating, ha_status}: { label: string, wifi_rating: string, ha_status: number}) {
  return <Stack justifyContent={"space-between"} direction={"row"}>
    <Typography component={"span"} variant="overline" align="left">{label}:</Typography>
    <Typography component={"span"} variant="overline" align="right" marginLeft={"10px"}>
      <img src={"/icons/wifi_" + wifi_rating + ".svg"} alt={"strength: " + wifi_rating} width={24} ></img>
      <ComponentHomeAssistantStatus ha_status={ha_status}/>
    </Typography>
  </Stack>
}

export function ComponentHomeAssistantStatus({ha_status}: {ha_status: number}) {
  const color = ha_status == 1 ? "green": "red"
  return <img src={"/icons/HomeAssistant_status_" + color + ".svg"} alt={"HomeAssistant Status: " + color} width={24} ></img>
}