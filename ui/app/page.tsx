import styles from './page.module.css'
import Grid2 from '@mui/material/Unstable_Grid2/Grid2'
import {  Divider } from '@mui/material'
import { Externalapps } from './ui/things/applicationLinks'
import { HealthStatusPublic } from './ui/elements/public_health_status'
import { ExternalS3 } from './ui/things/s3'

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
          <HealthStatusPublic target={"External API"} />
        </Grid2>
        <Grid2 xs={12}>
          <Divider variant='fullWidth' className={styles.code} style={{ marginTop: "50px", marginBottom: "50px" }}> Applications</Divider>
        </Grid2>
        <Grid2>
          <ExternalS3 internal={false}/>
        </Grid2>
      </Grid2>
    </main>
  )
}
