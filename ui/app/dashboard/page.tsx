import Image from 'next/image'
import styles from '../page.module.css'
import Grid2 from '@mui/material/Unstable_Grid2/Grid2'
import { Button, Card, CardActionArea, CardContent, CardHeader, Divider, Link } from '@mui/material'
import { Externalapps, InfrastructureApps, ProductionApps, UtilityApps } from '../ui/things/applicationLinks'
import { Evergreen, NAS, ProdK8s } from '../ui/things/status'
import PlaceholderComponent from '../ui/elements/placeholder'

export default function Dashboard() {
  return (
    <main className={styles.main}>
      <Grid2 container spacing={2}>
        <Grid2>
          <Externalapps />
        </Grid2>
        <Grid2>
          <ProductionApps />
        </Grid2>
        <Grid2>
          <UtilityApps />
        </Grid2>
        <Grid2>
          <InfrastructureApps />
        </Grid2>
      </Grid2>

      <Grid2 container spacing={2} justifyContent={"center"} xs={12}>
        <Grid2 xs={12}>
          <Divider variant='fullWidth' className={styles.code} style={{ marginTop: "50px", marginBottom: "50px" }}> Status Overview</Divider>
        </Grid2>
        <Grid2>
          <Evergreen />
        </Grid2>
        <Grid2>
          <NAS />
        </Grid2>
        <Grid2>
          <ProdK8s />
        </Grid2>
        <Grid2>
          <PlaceholderComponent name={"Status: API"} todos={["SecurityCheck passed"]} />
        </Grid2>
      </Grid2>
    </main>
  )
}
