import { Grid } from '@mui/material'
import { Evergreen, H3, KVM, NAS, ProdK8s } from '../ui/things/status'
import PlaceholderComponent from '../ui/elements/placeholder'
import Grid2 from '@mui/material/Unstable_Grid2/Grid2'

export default function Home() {
  return (
    <Grid2 container spacing={2} margin={"100px"}>
      <Grid2 container xs={12} >
      <Grid2>
        <Evergreen />
      </Grid2>
      <Grid2>
        <NAS />
      </Grid2>
      <Grid2>
        <ProdK8s />
      </Grid2>
    </Grid2>

    <Grid2 container xs={12} >
      <Grid2>
        <H3 />
      </Grid2>
      <Grid2>
        <KVM />
      </Grid2>
      <Grid2>
        <PlaceholderComponent name="Passive KVM" todos={[
          "Cockpit",
          "Nodeexporter"
        ]} />
      </Grid2>
      <Grid2>
        <PlaceholderComponent name="Network" todos={[
          "A1 router",
          "NG7 Router",
          "NG7oben AccessPoint",
          "Drucker",
          "Internet"
        ]} />
      </Grid2>
      <Grid2>
        <PlaceholderComponent name="Bastion" todos={[
          "External HTTP",
          "External S3",
          "NG7oben AccessPoint",
          "Drucker",
          "Internet"
        ]} />
      </Grid2>
      </Grid2>
    </Grid2>
  )
}
