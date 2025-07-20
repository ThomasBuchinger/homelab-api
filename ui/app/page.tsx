import styles from './page.module.css'
import { Divider, Grid, Stack, Typography } from '@mui/material'
import PaperlessHealth from './ui/elements/PaperlessHealth'
import SyncthingHealth from './ui/elements/SyncthingHealth'
import QuickinkButton from './ui/elements/QuickLinkButton'
import KubernetesHealth from './ui/elements/KubernetesHealth'
import NasHealth from './ui/elements/NasHealth'
import TvHealth from './ui/elements/TvHealth'

export default function Home() {
  return (
    <Grid container spacing={2} justifyContent={"center"}>
      <Grid size={{ xs: 12 }}>
        <Divider variant='fullWidth' className={styles.code} style={{ marginBottom: "50px" }} > HomeLab Overview</Divider>
      </Grid>
      <Grid>
        <KubernetesHealth />
      </Grid>
      <Grid>
        <NasHealth />
      </Grid>
      <Grid>
        <PaperlessHealth />
      </Grid>
      <Grid>
        <SyncthingHealth />
      </Grid>
      <Grid>
        <TvHealth />
      </Grid>
      <Grid size={{ xs: 12 }}>
        <Divider variant='fullWidth' className={styles.code} style={{ marginTop: "50px", marginBottom: "50px" }}> Links</Divider>
      </Grid>
      <Grid container justifyContent={'space-evenly'} size={{ xs: 10 }}>
        <Grid>
          <QuickinkButton name="H3" icon_url="/icons/cockpit.png" url="https://10.0.0.14:9090" />
        </Grid>
        <Grid>
          <QuickinkButton name="NAS" icon_url="/icons/unraid-logo.png" url="http://10.0.0.19" />
        </Grid>
        <Grid>
          <QuickinkButton name="BUC-KVM" icon_url="/icons/cockpit.png" url="https://10.0.0.13:9090" />
        </Grid>
        <Grid>
          <QuickinkButton name="BUC-KVM (IPMI)" icon_url="/icons/Baseboard-Management-Controller.jpg" url="https://10.0.0.18" />
        </Grid>
        <Grid>
          <QuickinkButton name="Bastion" icon_url="/icons/cockpit.png" url="https://10.0.0.22:9090" />
        </Grid>
        <Grid>
          <QuickinkButton name="Evergreen" icon_url="/icons/cockpit.png" url="https://10.0.0.16:9090" />
        </Grid>
        {/* <Grid2>
          <QuickinkButton name="OKD Evergreen" icon_url="/icons/kubernetes-icon-color.svg" url="http://evergreen-console.10.0.0.16.nip.io/k8s/all-namespaces/core~v1~Pod" />
        </Grid2>
        <Grid2>
          <QuickinkButton name="OKD Prod" icon_url="/icons/kubernetes-icon-color.svg" url="http://prod-console.10.0.0.21.nip.io/k8s/all-namespaces/core~v1~Pod" />
        </Grid2> */}
        <Grid>
          <QuickinkButton name="ArgoCD" icon_url="/icons/argo-icon-color.svg" url="http://argocd.10.0.0.21.nip.io" />
        </Grid>
        <Grid>
          <QuickinkButton name="Flux" icon_url="/icons/flux-icon-color.svg" url="http://flux.10.0.0.16.nip.io" />
        </Grid>
        <Grid>
          <QuickinkButton name="Grafana" icon_url="/icons/grafana-logo.svg" url="https://grafana.buc.sh" />
        </Grid>
        <Grid>
          <QuickinkButton name="Prometheus" icon_url="/icons/prometheus-icon-color.svg" url="https://grafana.buc.sh" />
        </Grid>
        <Grid>
          <QuickinkButton name="Minio" icon_url="/icons/minio-logo-old.webp" url="http://10.0.0.19:9001" />
        </Grid>
        <Grid>
          <QuickinkButton name="PDF" icon_url="/icons/stirling-pdf.png" url="https://pdf.buc.sh" />
        </Grid>
        <Grid>
          <QuickinkButton name="PiHole" icon_url="/icons/Pi-hole_vector_logo.svg" url="http://pihole.10.0.0.16.nip.io" />
        </Grid>
        <Grid>
          <QuickinkButton name="Public Files" icon_url="/icons/hfs.ico" url="http://files.buc.sh" />
        </Grid>
        <Grid>
          <QuickinkButton name="Wastebin" icon_url="/icons/wastebin-logo.png" url="https://paste.buc.sh" />
        </Grid>
        <Grid>
          <QuickinkButton name="Envoy" icon_url="/icons/envoy-icon-color.png" url="https://envoy.buc.sh" />
        </Grid>
        <Grid>
          <QuickinkButton name="NTFY" icon_url="/icons/ntfy-logo.svg" url="https://push.buc.sh" />
        </Grid>
      </Grid>
    </Grid>
  )
}
