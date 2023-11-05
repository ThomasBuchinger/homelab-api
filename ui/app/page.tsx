import styles from './page.module.css'
import Grid2 from '@mui/material/Unstable_Grid2/Grid2'
import {  Divider } from '@mui/material'
import { Externalapps } from './ui/things/applicationLinks'
import { HealthStatusPublic } from './ui/elements/public_health_status'

export default function Home() {
  return (
    <main className={styles.main}>
      <Grid2  container spacing={2} justifyContent={"center"} xs={12}>
        <Grid2 xs={12}>
          <Divider variant='fullWidth' className={styles.code} style={{ marginTop: "50px", marginBottom: "50px" }}> Status Overview</Divider>
        </Grid2>
        <Grid2 xs={2}>
          <HealthStatusPublic target={"Servers"} />
        </Grid2>
        <Grid2 xs={2}>
          <HealthStatusPublic target={"Network"} />
        </Grid2>
        <Grid2 xs={2}>
          <HealthStatusPublic target={"API"} />
        </Grid2>
\        <Grid2 xs={12}>
          <Divider variant='fullWidth' className={styles.code} style={{ marginTop: "50px", marginBottom: "50px" }}> Applications</Divider>
        </Grid2>
        <Grid2>
          <Externalapps />
        </Grid2>
      </Grid2>
    </main>
  )
}
