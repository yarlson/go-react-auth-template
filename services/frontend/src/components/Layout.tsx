import React from "react";
import { ModeToggle } from "@/components/ModeToggle";

interface LayoutProps {
  children: React.ReactNode;
}

const Layout: React.FC<LayoutProps> = ({ children }) => {
  return (
    <div className="min-h-screen flex flex-col">
      <div className="absolute top-4 right-4">
        <ModeToggle />
      </div>
      <div className="flex-grow flex items-center justify-center">
        {children}
      </div>
    </div>
  );
};

export default Layout;
