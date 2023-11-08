import { Card, CardContent, CardHeader } from "@mui/material";
import Image from "next/image";
import Link from "next/link";

export function ExternalS3({internal}: {internal: boolean}) {
  return ExternalAppIcon({title: "S3 API", hostname: "r2", icon: "minio-logo-old.webp", internal: internal})
}

export function ExternalAppIcon({title, hostname, icon, internal=false}: {title: string, hostname: string, icon: string, internal:boolean}) {
  const avatar = <Image src={"/icons/"+icon} alt="logo" width={50} height={50} ></Image>

  const element =<Card>
    <CardHeader avatar={avatar} title={title} subheader={"Hostname: "+hostname} titleTypographyProps={{ variant: "h6" }} subheaderTypographyProps={{ variant: "button" }} />
  </Card>

  return internal ? <Link href={"https://"+hostname+"cloud.buc.sh"}>{element}</Link> : element

}
