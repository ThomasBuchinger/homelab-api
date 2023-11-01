import Image from 'next/image'
import styles from './page.module.css'
import Grid2 from '@mui/material/Unstable_Grid2/Grid2'
import { Button, Card, CardActionArea, CardContent, CardHeader, Link } from '@mui/material'

export default function Home() {
  return (
    <main className={styles.main}>
      <Grid2 container spacing={2}>
        <Grid2>
          <Card>
            <CardHeader title="S3" />
            <CardContent>There is an S3 Bucket somewhere</CardContent>
            <CardActionArea>
              <Button href='/'>Open</Button>
            </CardActionArea>
          </Card>
        </Grid2>
      </Grid2>
    </main>
  )
}
