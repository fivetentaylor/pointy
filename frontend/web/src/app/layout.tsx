import "./globals.css";
import type { Metadata } from "next";
import { Inter } from "next/font/google";
import { ThemeProvider } from "@/components/providers/Theme";
import localFont from "next/font/local";

// Font files can be colocated inside of `app`
const marat = localFont({
  src: "./fonts/marat.woff2",
  display: "swap",
  variable: "--font-marat",
});

const maratMedium = localFont({
  src: "./fonts/marat_medium.woff2",
  display: "swap",
  variable: "--font-marat-medium",
});

const inter = Inter({
  subsets: ["latin"],
  display: "swap",
  variable: "--font-inter",
});

export const metadata: Metadata = {
  title: "Revi.so",
  description:
    "A modern writing tool designed for putting together well written work, and collaborating with a team to get it shipped",
  metadataBase: new URL("https://revi.so"),
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <head>
        <link rel="icon" href="/favicon.svg" type="image/svg+xml" />
      </head>
      <body
        className={`${inter.variable} ${marat.variable} ${maratMedium.variable}`}
      >
        {children}
      </body>
    </html>
  );
}
