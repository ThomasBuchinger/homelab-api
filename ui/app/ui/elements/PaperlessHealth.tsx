"use client"
import { CheckBox, Circle } from "@mui/icons-material";
import { Avatar, Box, Button, Card, CardActionArea, CardActions, CardContent, CardHeader, CardMedia, IconButton, Stack, Table, TableCell, TableRow, Typography } from "@mui/material";
import Image from "next/image";
import useSWR from "swr";
import { FallbackCardError, FallbackCardLoading } from "./FallbackCards";
import { ComponentContentValue, ComponentHeader } from "../parts/Component";


export default function PaperlessHealth() {
  const NAME="Paperless"
  const API_URL="/api/public/bff/paperless"
  const fetcher = async (url: string) => fetch(url).then(res => res.json())
  const { data, error, isLoading } = useSWR(API_URL, fetcher, { refreshInterval: 10000, fallbackData: { status: "invalid", messages: [] } })

  if (isLoading) {
    return FallbackCardLoading(NAME)
  }

  if (error) {
    return FallbackCardError(NAME, API_URL, error)
  }

  return (
    <Card>
      <CardActionArea href={data.url} target="_blank">
        {ComponentHeader("Paperless", "/icons/paperless-square.svg", API_URL, data.status, data.reason)}
        <CardContent>
          <ComponentContentValue label="Total Documents:" value={data.total_documents || "-"} />
        </CardContent>
      </CardActionArea>
      <CardActions>
        <Button href={data.alt_url}>NIP</Button>
      </CardActions>
    </Card>
  );
}


