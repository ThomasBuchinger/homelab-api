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

          <Box sx={{ flexGrow: 1, display: { xs: 'none', md: 'flex' } }} alignItems={"center"}>
            {/* <Button href='/public' sx={{ my: 2, color: 'white', display: 'block' }} >Public</Button> */}
          </Box>

          {/* <Box>
            {UserInfo(viewInternal, setViewInternal)}
          </Box> */}
        </Toolbar>
      </Box>
    </AppBar>

  );
}
