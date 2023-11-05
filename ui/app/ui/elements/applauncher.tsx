import { Box, Card, CardContent, CardHeader, IconButton, Stack, Typography } from "@mui/material";
import Image from "next/image";

export function AppLauncherComponent(props: any) {
  return <Card>
    <CardHeader title={props.name} titleTypographyProps={{ variant: "h6" }} />
    <CardContent >
      <Stack direction="row" spacing={5} justifyContent={"space-around"}>
        {props.children}
      </Stack>
    </CardContent>
  </Card>
}
export function ApplicationButton({ href, icon, text }: { href: string, icon: string, text: string }) {
  const size = 50
  return <IconButton href={href} target='_blank'>
    <Box>
      <Image src={"/icons/" + icon} alt={text} width={size} height={size} />
      <Typography color={"primary"}>{text}</Typography>
    </Box>
  </IconButton>
}
