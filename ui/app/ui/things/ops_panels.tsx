import { Button, Card, CardActions, CardContent, CardHeader, Typography } from "@mui/material";
import Image from "next/image";
import { ElementLinkTailscale } from "../elements/status";

export function GitopsEvergreen() {
  return <Card>
    <CardHeader title="Evergreen - Gitops" titleTypographyProps={{ variant: "h6" }}
      avatar={<Image src="/icons/github-mark.svg" alt="Edit DNS Entried" width={25} height={25} />} />
    <CardContent>
      <Typography>Edit DNS Settings</Typography>
      <Button href={"https://github.com/ThomasBuchinger/voodoo-gitops/blob/main/gitops/pihole/pihole-helm.yaml#L48"} target="_blank" >DNS config</Button>
    </CardContent>
    <CardContent >
      <Typography>Certificates for Prod</Typography>
      <Button href={"https://github.com/ThomasBuchinger/voodoo-gitops/blob/main/gitops/stuff/prod/certs.yaml"} target="_blank" >Certificates</Button>
      <Button href={"https://github.com/ThomasBuchinger/voodoo-gitops/blob/main/gitops/vault/instance/patch-vault-content.yaml"} target="_blank" >Vault Import</Button>
    </CardContent>
  </Card>
}

export function HomepageOps() {
  return <Card>
    <CardHeader title="Homepage" titleTypographyProps={{ variant: "h6" }}
      avatar={<Image src="/icons/github-mark.svg" alt="Edit DNS Entried" width={25} height={25} />} />
    <CardContent>
      <Typography>Renovate PRs</Typography>
      <Button href={"https://github.com/ThomasBuchinger/homepage-react/pulls"} target="_blank" >Open PRs</Button>
      <Button href={"http://evergreen-console.10.0.0.16.nip.io/k8s/ns/homepage/core~v1~Pod"} target="_blank" >Delete Pod</Button>
    </CardContent>
  </Card>
}

export function Tailscale() {
  return (
    <Card>
      <CardHeader title="Tailscale" subheader="infinite.void@outlook.com" subheaderTypographyProps={{ marginRight: "20px" }} action={
        <ElementLinkTailscale enabled={true} />
      } />
      <CardActions>
        <Button variant="outlined" href={""} target="_blank">PiKVM Mini</Button>
        <Button variant="outlined" disabled={true} href={""} target="_blank">PiKVM AIO</Button>
      </CardActions>
    </Card>
  )
}

export function ContainerImages() {
  return <Card>
    <CardHeader title="Images" />
    <CardContent>
      <Typography>docker.io/node:latest: outdated </Typography>
      <Typography>ghcr.io.io/thomasbuchinger/nas-lcd:latest: ok</Typography>
      <Typography>ghcr.io.io/thomasbuchinger/homepage-react:latest: ok</Typography>
      <Typography><Button>Update Repo</Button> argocd:2.5.4: outdated</Typography>
    </CardContent>
  </Card>
}