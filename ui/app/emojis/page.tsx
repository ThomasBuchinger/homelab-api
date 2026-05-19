'use client'
import { toast } from 'react-toastify';
import styles from '../page.module.css'
import { Accordion, AccordionDetails, AccordionSummary, Button, Card, CardContent, CardHeader, Collapse, Divider, Grid, IconButton, Stack, Typography } from '@mui/material'
import { ArrowDownward, ArrowDownwardSharp, ArrowDropDownCircleOutlined, ArrowUpwardSharp } from '@mui/icons-material';
import { use, useState } from 'react';

function copyToClipboard(text: string) {
  navigator.clipboard.writeText(text);

  toast("Copied: " + text)
}

export default function Home() {
  return (
    <Grid container spacing={2} justifyContent={"center"}>
      <Grid size={{ xs: 12 }}>
        <Divider variant='fullWidth' className={styles.code} style={{ marginBottom: "50px" }} >Stuff to Copy</Divider>
      </Grid>
      <Grid container direction={'row'} justifyContent={'space-evenly'}>
        <Grid><Emojis /></Grid>
        <Grid><Umlaute /></Grid>
      </Grid>
    </Grid>
  )
}


function Umlaute() {
  return <Card>
    <CardContent>
      <Button onClick={() => copyToClipboard("ö")}>
        <Typography textTransform={'lowercase'}>ö</Typography>
      </Button>
      <Button onClick={() => copyToClipboard("Ö")}>
        <Typography>Ö</Typography>
      </Button>
    </CardContent>
    <CardContent>
      <Button onClick={() => copyToClipboard("ü")}>
        <Typography textTransform={'lowercase'}>ü</Typography>
      </Button>
      <Button onClick={() => copyToClipboard("Ü")}>
        <Typography>Ü</Typography>
      </Button>
    </CardContent>
    <CardContent>
      <Button onClick={() => copyToClipboard("ä")}>
        <Typography textTransform={'lowercase'}>ä</Typography>
      </Button>
      <Button onClick={() => copyToClipboard("Ä")}>
        <Typography>Ä</Typography>
      </Button>
    </CardContent>
    <CardContent>
      <Button onClick={() => copyToClipboard("ß")}>
        <Typography textTransform={'lowercase'}>ß</Typography>
      </Button>
    </CardContent>
  </Card>
}
function Emojis() {
  const [expand, setExpand] = useState(true);

  return <Card>
    <CardHeader
      action={
        <IconButton onClick={() => setExpand(!expand)}>{expand ? <ArrowUpwardSharp /> : <ArrowDownwardSharp />}</IconButton>
      }
    />
    <CardContent>
      <Button onClick={() => copyToClipboard("¯\\_(ツ)_/¯")}>
        <Typography fontSize={24}>¯\_(ツ)_/¯</Typography>
      </Button>
    </CardContent>
    <CardContent>
      <Button onClick={() => copyToClipboard("(╯°□°)╯︵ ┻━┻")}>
        <Typography fontSize={24}>(╯°□°)╯︵ ┻━┻</Typography>
      </Button>
    </CardContent>
    <Collapse in={expand} title='More'>
      <CardContent>
        <Button onClick={() => copyToClipboard("ಠ_ಠ")}>
          <Typography fontSize={24}>ಠ_ಠ</Typography>
        </Button>
      </CardContent>
      <CardContent>
        <Button onClick={() => copyToClipboard("(¬_¬)")}>
          <Typography fontSize={24}>(¬_¬)</Typography>
        </Button>
      </CardContent>
      <CardContent>
        <Button onClick={() => copyToClipboard("ʘ︵ʘ")}>
          <Typography fontSize={24}>ʘ︵ʘ</Typography>
        </Button>
      </CardContent>
      <CardContent>
        <Button onClick={() => copyToClipboard("(>_<)")}>
          <Typography fontSize={24}>{"(>_<)"}</Typography>
        </Button>
      </CardContent>
      <CardContent>
        <Button onClick={() => copyToClipboard(".·´¯`(>▂<)´¯`·.")}>
          <Typography fontSize={24}>{".·´¯`(>▂<)´¯`·."}</Typography>
        </Button>
      </CardContent>
      <CardContent>
        <Button onClick={() => copyToClipboard("(｡◕‿‿◕｡)")}>
          <Typography fontSize={24}>(｡◕‿‿◕｡)</Typography>
        </Button>
      </CardContent>
      <CardContent>
        <Button onClick={() => copyToClipboard("[̲̅$̲̅(̲̅ ͡° ͜ʖ ͡°̲̅)̲̅$̲̅]")}>
          <Typography fontSize={24}>[̲̅$̲̅(̲̅ ͡° ͜ʖ ͡°̲̅)̲̅$̲̅]</Typography>
        </Button>
      </CardContent>
      <CardContent>
        <Button onClick={() => copyToClipboard("ヽ( ಠ益ಠ )ﾉ")}>
          <Typography fontSize={24}>ヽ( ಠ益ಠ )ﾉ</Typography>
        </Button>
      </CardContent>
      <CardContent>
        <Button onClick={() => copyToClipboard("(ノಠ益ಠ)ノ彡┻━┻")}>
          <Typography fontSize={24}>(ノಠ益ಠ)ノ彡┻━┻</Typography>
        </Button>
      </CardContent>
      <CardContent>
        <Button onClick={() => copyToClipboard("┬─┬ノ( º _ ºノ)")}>
          <Typography fontSize={24}>┬─┬ノ( º _ ºノ)</Typography>
        </Button>
      </CardContent>
      <CardContent>
        <Button onClick={() => copyToClipboard('(ノಠ ∩ಠ)ノ彡( \\o°o)\\')}>
          <Typography fontSize={24}>{'(ノಠ ∩ಠ)ノ彡( \\o°o)\\'}</Typography>
        </Button>
      </CardContent>
      <CardContent>
        <Button onClick={() => copyToClipboard("(≧σ≦)")}>
          <Typography fontSize={24}>(≧σ≦)</Typography>
        </Button>
      </CardContent>
      <CardContent>
        <Button onClick={() => copyToClipboard("(\/) (°,,°) (\/)")}>
          <Typography fontSize={24}>(\/) (°,,°) (\/)</Typography>
        </Button>
      </CardContent>
      <CardContent>
        <Button onClick={() => copyToClipboard("(｡◕‿◕｡)")}>
          <Typography fontSize={24}>(｡◕‿◕｡)</Typography>
        </Button>
      </CardContent>
      <CardContent>
        <Button onClick={() => copyToClipboard("(◕︵◕)")}>
          <Typography fontSize={24}>(◕︵◕)</Typography>
        </Button>
      </CardContent>
      <CardContent>
        <Button onClick={() => copyToClipboard("┻━┻︵ \(°□°)/ ︵ ┻━┻")}>
          <Typography fontSize={24}>┻━┻︵ \(°□°)/ ︵ ┻━┻</Typography>
        </Button>
      </CardContent>
    </Collapse>

  </Card>
}