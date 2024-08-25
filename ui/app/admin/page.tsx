"use client"
import { ArrowForward, ContentPaste, PlayArrow, Send, SettingsEthernet, Share } from '@mui/icons-material';
import { Avatar, Box, Button, Card, CardActionArea, CardActions, CardContent, CardHeader, Container, Dialog, DialogActions, DialogContent, DialogTitle, FormControl, InputLabel, MenuItem, OutlinedInput, Paper, Select, SelectChangeEvent, TextField, Typography } from '@mui/material'
import Grid2 from '@mui/material/Unstable_Grid2/Grid2'
import Link from 'next/link'
import { useState } from 'react';


export default function Home() {
  const [open, setOpen] = useState(false);
  const [exposeTimeout, setExposeTimeout] = useState<number>(60);

  const handleChange = (event: SelectChangeEvent<typeof exposeTimeout>) => {
    setExposeTimeout(Number(event.target.value));
  };

  const handleClickOpen = () => {
    setOpen(true);
  };

  const handleClose = (event: React.SyntheticEvent<unknown>, reason?: string) => {
    if (reason !== 'backdropClick') {
      setOpen(false);
    }
  };

  const api_data = {
    "name": "r2",
    "web_url": "https://dash.cloud.buc.sh",
    "api_url": "https://dash.cloud.buc.sh:1234",
    "exposed": true,
    "expireation": (new Date((new Date()).valueOf() + exposeTimeout * 60 * 1000)).toString()
  }
  const ui_data = {
    "s3": {
      "title": "S3 WebUI",
      "icon": "/icons/minio-logo-old.webp"
    }
  }


  return (
    <main >
      <Grid2 container spacing={2} justifyContent={"center"} alignContent={"center"} xs={12} height={"1000px"}>
        <Grid2 >
          <Card>
            <CardHeader 
                title={"S3 WebUI"} 
                subheader="https://dash.cloud.buc.sh" 
                titleTypographyProps={{ variant: "h6" }}
                subheaderTypographyProps={{ variant: "h6", paddingRight: "30px" }}
                avatar ={<Avatar src='/icons/minio-logo-old.webp' />}
                />
            <CardContent>
              <Grid2 container justifyContent={"space-between"}>
                <Grid2><Typography variant='body1'>Status:</Typography></Grid2>
                <Grid2><Typography variant='body1'>Exposed | {exposeTimeout} minutes</Typography></Grid2>
              </Grid2>
              <Grid2 container justifyContent={"space-between"}>
                <Grid2><Typography variant='body1'>Policy:</Typography></Grid2>
                <Grid2><Typography variant='body1'>Basic Auth</Typography></Grid2>
              </Grid2>
            </CardContent>
            <CardActions>
              <Box display={"flex"} flexDirection={"row"} flexGrow={1} justifyContent={"space-evenly"}>

                <Button startIcon={<Share />} variant='outlined' onClick={handleClickOpen}>Expose</Button>
                {ExposeServiceForDurationDialog(exposeTimeout, open, handleClose, handleChange)}
                <Button endIcon={<PlayArrow />} variant='contained' href={"https://dash.cloud.buc.sh"}>GOTO </Button>
              </Box>
            </CardActions>
          </Card>
        </Grid2>
      </Grid2>
    </main>
  )
}

function ExposeServiceForDurationDialog(selectedExposeTimeput, open, handleClose, handleChange){
  const now = new Date()
  const expireTime = new Date(now.valueOf() + (selectedExposeTimeput * 60 * 1000))

  return <Dialog disableEscapeKeyDown open={open} onClose={handleClose}>
    <DialogTitle>Expose Service: S3 WebUI</DialogTitle>
    <DialogContent>
      <Box component="form" sx={{ display: 'flex', flexWrap: 'nowrap' }}>
        <FormControl sx={{ m: 1, minWidth: 120 }}>
          <InputLabel htmlFor="expire_duration">Expires in: </InputLabel>
          <Select native value={selectedExposeTimeput}
            onChange={handleChange}
            input={<OutlinedInput label="Expires in:" id="expire_duration" />}
          >
            <option value={10}>10m</option>
            <option value={30}>30m</option>
            <option value={60}>1h</option>
            <option value={600}>10h</option>
            <option value={1440}>24h</option>
          </Select>
        </FormControl>
        <Box display={"flex"} flexGrow={1} justifyContent={"center"} alignContent={"center"}>
          <Typography alignSelf={"center"}>{expireTime.getHours()}:{expireTime.getMinutes()}</Typography>
        </Box>
      </Box>
    </DialogContent>
    <DialogActions>
      <Button
        onClick={handleClose}
      >Cancel</Button>
      <Button
        onClick={handleClose}
      >Ok</Button>
    </DialogActions>
  </Dialog>


}


function AppRoutingStatus(){

}