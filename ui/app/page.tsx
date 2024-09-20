import styles from './page.module.css'
import {  Divider, Grid2, Stack, Typography } from '@mui/material'
import PaperlessHealth from './ui/elements/PaperlessHealth'
import SyncthingHealth from './ui/elements/SyncthingHealth'
import QuickinkButton from './ui/elements/QuickLinkButton'
import KubernetesHealth from './ui/elements/KubernetesHealth'
import NasHealth from './ui/elements/NasHealth'
import TvHealth from './ui/elements/TvHealth'

export default function Home() {
  return (
      <Grid2  container spacing={2} justifyContent={"center"}>
        <Grid2 size={{ xs: 12 }}>
          <Divider variant='fullWidth' className={styles.code} style={{ marginBottom: "50px" }} > HomeLab Overview</Divider>
        </Grid2>
        <Grid2>
          <KubernetesHealth />
        </Grid2>
        <Grid2>
          <NasHealth />
        </Grid2>
        <Grid2>
          <PaperlessHealth />
        </Grid2>
        <Grid2>
          <SyncthingHealth />
        </Grid2>
        <Grid2>
          <TvHealth />
        </Grid2>
        <Grid2 size={{ xs: 12 }}>
          <Divider variant='fullWidth' className={styles.code} style={{ marginTop: "50px", marginBottom: "50px" }}> Links</Divider>
        </Grid2>
        <Grid2 container justifyContent={'space-evenly'} size={{ xs: 10 }}>
          <Grid2>
            <QuickinkButton name="H3" icon_url="/icons/cockpit.png" url="https://10.0.0.14:9090" />
          </Grid2>
          <Grid2>
            <QuickinkButton name="NAS" icon_url="/icons/unraid-logo.png" url="http://10.0.0.19" />
          </Grid2>
          <Grid2>
            <QuickinkButton name="BUC-KVM" icon_url="/icons/cockpit.png" url="https://10.0.0.13:9090" />
          </Grid2>
          <Grid2>
            <QuickinkButton name="Bastion" icon_url="/icons/cockpit.png" url="https://10.0.0.22:9090" />
          </Grid2>
          <Grid2>
            <QuickinkButton name="Evergreen" icon_url="/icons/cockpit.png" url="https://10.0.0.16:9090" />
          </Grid2>
          <Grid2>
            <QuickinkButton name="OKD Evergreen" icon_url="/icons/kubernetes-icon-color.svg" url="http://evergreen-console.10.0.0.16.nip.io/k8s/all-namespaces/core~v1~Pod" />
          </Grid2>
          <Grid2>
            <QuickinkButton name="OKD Prod" icon_url="/icons/kubernetes-icon-color.svg" url="http://prod-console.10.0.0.21.nip.io/k8s/all-namespaces/core~v1~Pod" />
          </Grid2>
          <Grid2>
            <QuickinkButton name="ArgoCD" icon_url="/icons/argo-icon-color.svg" url="http://argocd.10.0.0.21.nip.io" />
          </Grid2>
          <Grid2>
            <QuickinkButton name="Flux" icon_url="/icons/flux-icon-color.svg" url="http://flux.10.0.0.16.nip.io" />
          </Grid2>
          <Grid2>
            <QuickinkButton name="Gatus" icon_url="/icons/gatus.svg" url="http://gatus.10.0.0.16.nip.io" />
          </Grid2>
          <Grid2>
            <QuickinkButton name="Vault" icon_url="/icons/vault.svg" url="https://vault.buc.sh" />
          </Grid2>
          <Grid2>
            <QuickinkButton name="Grafana" icon_url="/icons/grafana-logo.svg" url="https://grafana.buc.sh" />
          </Grid2>
          <Grid2>
            <QuickinkButton name="Prometheus" icon_url="/icons/prometheus-icon-color.svg" url="https://grafana.buc.sh" />
          </Grid2>
          <Grid2>
            <QuickinkButton name="Tracing" icon_url="/icons/grafana-logo.svg" url="https://grafana.buc.sh" />
          </Grid2>
          <Grid2>
            <QuickinkButton name="Minio" icon_url="/icons/minio-logo-old.webp" url="http://10.0.0.19:9001" />
          </Grid2>
          <Grid2>
            <QuickinkButton name="PDF" icon_url="/icons/stirling-pdf.png" url="https://pdf.buc.sh" />
          </Grid2>
            <QuickinkButton name="PiHole" icon_url="/icons/Pi-hole_vector_logo.svg" url="http://pihole.10.0.0.16.nip.io" />
          <Grid2>
            <QuickinkButton name="Wastbin" icon_url="/icons/wastebin-logo.png" url="https://paste.buc.sh" />
          </Grid2>
          <Grid2>
            <QuickinkButton name="Envoy" icon_url="/icons/React Logo.svg" url="https://envoy.buc.sh" />
          </Grid2>
          <Grid2>
            <QuickinkButton name="NTFY" icon_url="/icons/React Logo.svg" url="https://push.buc.sh" />
          </Grid2>
          <Grid2>
            <QuickinkButton name="S3 Browser" icon_url="/icons/React Logo.svg" url="https://r2.buc.sh" />
          </Grid2>
        </Grid2>
      </Grid2>
  )
}
