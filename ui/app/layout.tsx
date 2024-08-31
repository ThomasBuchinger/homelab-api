import type { Metadata } from 'next'
import { Inter } from 'next/font/google'
import './globals.css'
import TitleBar from './ui/elements/Titelbar'
import { Bounce, ToastContainer } from 'react-toastify'
import "react-toastify/dist/ReactToastify.css";


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
        <ToastContainer
          position="top-right"
          autoClose={5000}
          hideProgressBar={false}
          newestOnTop={false}
          closeOnClick
          rtl={false}
          pauseOnFocusLoss
          draggable
          pauseOnHover
          theme="light"
          transition={Bounce}
           />
        {children}
      </body>
    </html>
  )
}
