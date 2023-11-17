import { ArrowLeftSharp, Help, HelpOutline, HelpOutlineSharp, HelpOutlined } from "@mui/icons-material";
import { Button, Card, CardActionArea, CardContent, CardHeader, IconButton, Tooltip } from "@mui/material";
import Image from "next/image";
import Link from "next/link";

export function ExternalS3({internal}: {internal: boolean}) {
  return ExternalAppIcon({title: "S3 API", hostname: "r2", icon: "minio-logo-old.webp", internal: internal})
}

export function ExternalAppIcon({title, hostname, icon, internal=false}: {title: string, hostname: string, icon: string, internal:boolean}) {
  const avatar = <Image src={"/icons/"+icon} alt="logo" width={50} height={50} ></Image>

  // 49791
  const element =<Card>
    <CardHeader avatar={avatar} title={title} subheader={"Hostname: "+hostname} titleTypographyProps={{ variant: "h6" }} subheaderTypographyProps={{ variant: "button" }}
      action={
        <Tooltip title="Port: 4-9-7-9-1" leaveDelay={1000}>
          <IconButton><HelpOutlineSharp /> </IconButton>
        </Tooltip>
      }
      />
  </Card>

  return internal ? <Link href={"https://"+hostname+"cloud.buc.sh"}>{element}</Link> : element

}
