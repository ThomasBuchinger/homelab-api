"use client"
import { Button, CardActions, CardMedia, Paper } from "@mui/material";


export default function QuickinkButton({name, icon_url, url}: {name: string, icon_url: string, url: string}) {

  return (
    <Paper square={false}>
      <a href={url} title={name} target="_blank">
        <CardMedia image={icon_url} sx={{width: "50px", height: "50px", marginTop: "10px", margin: "auto"}} />
        <CardActions>
          <Button>{name}</Button>
        </CardActions>
      </a>
    </Paper>
  );
}


