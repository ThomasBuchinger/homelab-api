"use client"
import { useEffect, useState } from 'react';
import { AppBar, Box, Button, Chip, Toolbar, Typography } from '@mui/material';
import { Cottage, Public } from '@mui/icons-material';

export default function TitleBar() {
  const [viewInternal, setViewInternal] = useState(false)
  return (
    <AppBar position="static">
      <Box sx={{
        paddingX: "50px",
        backgroundImage: "/icons/kubernetes-hero-texture.png",
        backgroundColor: "#303030"
      }}>
        <Toolbar disableGutters >
          <Typography
            variant="h6"
            noWrap
            component="a"
            href="/"
            sx={{
              mr: 2,
              display: { xs: 'none', md: 'flex' },
              fontFamily: 'monospace',
              fontWeight: 700,
              letterSpacing: '.3rem',
              color: 'inherit',
              textDecoration: 'none',
            }}
            >
            BUC LAB Homepage
          </Typography>

          <Box sx={{ flexGrow: 1, display: { xs: 'none', md: 'flex' }}} alignItems={"center"}>
            {viewInternal ? <Button href='/dashboard' sx={{ my: 2, color: 'white', display: 'block' }} >Dashboard</Button> : <></>}
            {viewInternal ? <Button href='/status' sx={{ my: 2, color: 'white', display: 'block' }} >Status</Button> : <></>}
            {viewInternal ? <Button href='/apps' sx={{ my: 2, color: 'white', display: 'block' }} >Application</Button> : <></>}
            {viewInternal ? <Button href='/notes' sx={{ my: 2, color: 'white', display: 'block' }} >Notes</Button> : <></ >}
          </Box>

          <Box>
            {UserInfo(viewInternal, setViewInternal)}
          </Box>
        </Toolbar>
      </Box>
    </AppBar>

);
}

function UserInfo(view_internal: boolean, setViewMode: CallableFunction) {
  const [clientCfg, setClientCfg] = useState({ country:"", internal: false, ip: "" })
  useEffect(() => {
    fetch("/api/public/client-config")
      .then(response => response.json())
      .then(data => { setViewMode(data.internal); return setClientCfg(data)})
  }, [])
  const toggleFunc = () => {
    if (clientCfg.internal) {
      setViewMode(!view_internal)
    } else {
      setViewMode(false)
    }
  }
  const label = clientCfg.ip + " | " + clientCfg.country

  const chip = clientCfg.internal && view_internal ?
    <Chip variant={'filled'} color={"info"}    icon={<Cottage />} label={label} /> :
    <Chip variant={'filled'} color={"warning"} icon={<Public />} label={label} />


  return <Button onClick={toggleFunc} disabled={!clientCfg.internal}
    sx={{
      mr: 2,
      display: { xs: 'none', md: 'flex' },
      fontFamily: 'monospace',
      fontWeight: 700,
      // letterSpacing: '.3rem',
      color: 'inherit',
      textDecoration: 'none',
    }}>
    {chip}
  </Button>
}