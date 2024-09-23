import React from "react";
import { Header } from "@/components/Header";

interface LayoutProps {
  children: React.ReactNode;
}

const Layout: React.FC<LayoutProps> = ({ children }) => {
  return (
    <div className="flex flex-col">
      <Header />
      <main className="flex-grow flex justify-center mt-8">{children}</main>
    </div>
  );
};

export default Layout;
