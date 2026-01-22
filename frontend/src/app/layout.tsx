import type { Metadata } from "next";
import { Geist, Geist_Mono } from "next/font/google";
import "./globals.css";

const geistSans = Geist({
    variable: "--font-geist-sans",
    subsets: ["latin"],
});

const geistMono = Geist_Mono({
    variable: "--font-geist-mono",
    subsets: ["latin"],
});

export const metadata: Metadata = {
    title: "Open-Gerege | Монголын хөгжүүлэгчдэд зориулсан платформ",
    description: "Enterprise-grade backend API template болон вэб аппликейшн framework. Go болон Next.js дээр суурилсан, Clean Architecture зарчмаар бүтээгдсэн нээлттэй эхийн төсөл.",
    keywords: ["Open Source", "Go", "Next.js", "API", "Backend", "Frontend", "Mongolia", "Gerege"],
    authors: [{ name: "Gerege Core Team" }],
    openGraph: {
        title: "Open-Gerege",
        description: "Монголын хөгжүүлэгчдэд зориулсан платформ",
        type: "website",
        locale: "mn_MN",
    },
};

export default function RootLayout({
    children,
}: Readonly<{
    children: React.ReactNode;
}>) {
    return (
        <html lang="mn">
            <body
                className={`${geistSans.variable} ${geistMono.variable} antialiased`}
            >
                {children}
            </body>
        </html>
    );
}
