"use client"
import { Button, CardActionArea, CardActions, CardMedia, Paper, Typography } from "@mui/material";


export default function QuickinkButton({name, icon_url, url}: {name: string, icon_url: string, url: string}) {

  return (
    <Paper square={false} sx={{ paddingTop: "5px", paddingX: "5px" }}>
      <CardActionArea href={url} title={name} target="_blank">
        <CardMedia image={icon_url} sx={{ width: "50px", height: "50px", margin: "auto"}} />
        <CardActions>
          <Typography variant="button">{name}</Typography>
        </CardActions>
      </CardActionArea>
    </Paper>
  );
}


