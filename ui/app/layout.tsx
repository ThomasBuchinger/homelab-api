import type { Metadata } from 'next'
import { Inter } from 'next/font/google'
import './globals.css'
import TitleBar from './ui/things/titlebar'


const inter = Inter({ subsets: ['latin'] })

export const metadata: Metadata = {
  title: 'BUC HomeLAB',
  description: 'Homelab Dashboard and API',
}

export default function RootLayout({children,}: {children: React.ReactNode}) {
  return (
    <html lang="en">
      <body className={inter.className}>
        <TitleBar />
        {children}
        </body>
    </html>
  )
}
