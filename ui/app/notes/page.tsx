import { Card, CardContent, CardHeader, Grid, Typography } from '@mui/material'


export default function Home() {
  return (
    <Grid container spacing={2} justifyContent={"space-around"}>
      <Grid item>
      <Card>
        <CardHeader title="Notes and TODOs"/>
        <CardContent>
          <Typography>Hello World</Typography>
          </CardContent>
      </Card>

      </Grid>
    </Grid>
  )
}
