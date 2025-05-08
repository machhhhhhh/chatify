import "./globals.css";
import { ReactNode } from "react";

export const metadata = {
  title: "Chatify",
  description: "Your awesome app",
};

export default function RootLayout({ children }: { children: ReactNode }) {
  return (
    <html lang="en">
      <body className="antialiased bg-gray-50 text-gray-800">{children}</body>
    </html>
  );
}
