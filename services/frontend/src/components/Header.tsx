import { Link } from "react-router-dom";
import { ModeToggle } from "@/components/ModeToggle";

export function Header() {
  return (
    <header className="sticky top-0 z-40 shadow-sm">
      <nav className="container mx-auto px-4 py-4 flex justify-between items-center">
        <Link
          to="/profile"
          className="text-2xl font-bold text-transparent bg-clip-text bg-gradient-to-r from-purple-400 to-pink-600"
        >
          Sample App
        </Link>
        <div className="flex items-center space-x-4">
          <Link to="/profile">Profile</Link>
          <Link to="/logout">Logout</Link>
          <ModeToggle />
        </div>
      </nav>
    </header>
  );
}
