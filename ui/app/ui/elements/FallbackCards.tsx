"use client"
import { Avatar, Card, CardContent, CardHeader, LinearProgress, Skeleton, Typography } from "@mui/material";
import Link from "next/link";

export function FallbackCardLoading(title: string) {
  return <Card>
    <CardHeader
      // avatar={<CircularProgress variant="indeterminate" />}
      avatar={<Avatar><Skeleton variant="circular" /></Avatar>}
      title={title}
      titleTypographyProps={{ variant: "h5", marginX: "10px" }}
      />
    <CardContent>
      <LinearProgress />
    </CardContent>
    <CardContent>
      <Skeleton variant="rounded" />
    </CardContent>
  </Card>
}


export function FallbackCardError(title: string, api_url: string, error: any) {
  return <Card>
    <CardHeader
      avatar={<Avatar sx={{bgcolor: "red"}}>E</Avatar>}
      title={title}
      titleTypographyProps={{ variant: "h5", marginX: "10px" }}
    />
    <CardContent>
      <Typography variant="body1">Failed to load Data:</Typography>
      <Link style={{color: "blue"}} target="_blank" href={api_url}>{api_url}</Link>
    </CardContent>
  </Card>
}