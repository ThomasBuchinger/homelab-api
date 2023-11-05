import { Button } from "@mui/material";
import { AppLauncherComponent, ApplicationButton,  } from "../elements/applauncher";

export function Externalapps() {
  return <AppLauncherComponent name="External Applications">
    <ApplicationButton href="http://api.cloud.buc.sh" icon="minio-logo-old.webp" text="S3 Bucket" />
  </AppLauncherComponent>
}

export function ProductionApps() {
  return (
    <AppLauncherComponent name="Production">
      <ApplicationButton href="https://paperless.buc.sh" icon="paperless-square.svg" text="Paperless" />
      <ApplicationButton href="http://10.0.0.20:30024/" icon="qbittorrent-new-light.svg" text="qBittorrent" />
      <ApplicationButton href="https://syncthing.buc.sh" icon="syncthing-logo-128.png" text="Syncthing" />
      <ApplicationButton href="http://files.bus.sh" icon="kubernetes-pv-labeled.svg" text="iSCSI Files" />
    </AppLauncherComponent>
  )
}
export function UtilityApps() {
  return (
    <AppLauncherComponent name="Utility">
      <ApplicationButton href="http://pdf.bus.sh" icon="stirling-transparent.svg" text="PDF Tools" />
      <ApplicationButton href="http://paste.buc.sh" icon="wastebin-logo.png" text="Wastebin" />
    </AppLauncherComponent>
  )
}
export function TrialApps() {
  return (
    <AppLauncherComponent name="Trial Apps">
      <Button variant="outlined" href={""} target="_blank">KASM</Button>
      <Button variant="outlined" href={""} target="_blank">LocalAI</Button>
    </AppLauncherComponent>
  )
}
export function InfrastructureApps() {
  return (
    <AppLauncherComponent name="Infrastructure">
      <ApplicationButton href="https://grafana.buc.sh" icon="grafana-logo.svg" text="Grafana" />
      <ApplicationButton href="https://vault.buc.sh/ui/" icon="vault.svg" text="Vault" />
      <ApplicationButton href="http://s3.buc.sh:9002/" icon="minio-logo-old.webp" text="Minio" />
    </AppLauncherComponent>
  )
}
export function OtherApps() {
  return (
    <AppLauncherComponent name="Other">
      <ApplicationButton href="http://status.buc.sh" icon="React Logo.svg" text="Homepage-React" />
      <ApplicationButton href="http://gatus.10.0.0.16.nip.io/" icon="gatus.svg" text="Gatus" />
      <ApplicationButton href="http://pihole.10.0.0.16.nip.io/admin/" icon="Pi-hole_vector_logo.svg" text="Pihole" />
    </AppLauncherComponent>
  )
}

