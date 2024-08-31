"use client"
import { Button, Card, CardActionArea, CardActions, CardContent, Stack, Typography } from "@mui/material";
import useSWR from "swr";
import { FallbackCardError, FallbackCardLoading } from "./FallbackCards";
import { ComponentContentValue, ComponentHeader } from "../parts/Component";
import { format } from "util";



export default function NasHealth() {
  const NAME = "nas.buc.sh"
  const API_URL = "/api/component/nasv3"
  const fetcher = async (url: string) => fetch(url).then(res => res.json())
  const { data, error, isLoading } = useSWR(API_URL, fetcher, { refreshInterval: 10000, fallbackData: { status: "INVALID", reason: "not implemented" } })

  if (isLoading) {
    return FallbackCardLoading(NAME)
  }

  if (error) {
    return FallbackCardError(NAME, API_URL, error)
  }
  type nasDisk = { display_name: string, smart_ok: boolean, btrfs_error: boolean, capacity_total: Number, capacity_free: Number }
  const inTB = (x: any) => { return (Number(x) / 1000000000000).toFixed(2) }
  const sortfunc = (a: nasDisk, b: nasDisk) => { return a.display_name > b.display_name ? 1 : -1 }
  const percent = (Number(data.disk_free) / Number(data.disk_total) * 100).toFixed(1)

  return (

    <Card>
      <CardActionArea target="_blank" href={data.url}>
        {ComponentHeader(NAME, "/icons/unraid-logo.png", API_URL, data.status, data.reason)}
        <CardContent>
          <ComponentContentValue label={"Parity Check"} value={data.parity_status} />
          <ComponentContentValue label={"Backup"} value={data.backup_status + " - X days"} />
          <ComponentContentValue label={"Total"} value={inTB(data.disk_free) + " / " + inTB(data.disk_total) + " TB (" + percent + "%)"} />
        </CardContent>
        <CardContent>
          Disks (Free Space)
          {data.disks.sort(sortfunc).map((dir: nasDisk) => {
            return <ComponentContentValue key={dir.display_name} label={format("%s(%s)", dir.display_name, dir.btrfs_error ? "OK" : "ERR")} value={format("%s / %s TB", inTB(dir.capacity_free), inTB(dir.capacity_total))} />
          })}
        </CardContent>
      </CardActionArea>
      <CardActions>
        <Button href={data.backup_url} target="_blank">Backup</Button>
      </CardActions>
    </Card>
  );
}


