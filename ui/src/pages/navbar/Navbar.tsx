import { Button } from "@/components/ui/button";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { useCurrentUser } from "@/contexts/userContext";
import { Link, Outlet, useNavigate } from "react-router";
import Cookies from "js-cookie";

export function Navbar() {
  return (
    <>
      <nav className="fixed top-0 left-0 right-0 bg-white shadow-md z-50 h-14">
        <div className="container mx-auto px-4 py-2 flex justify-between items-center">
          <Link to="/" className="text-xl font-bold cursor-pointer">
            Quizzer
          </Link>
          <div className="hidden md:flex space-x-4">
            <UserButton />
          </div>
        </div>
      </nav>
      <main className="pt-14">
        <Outlet />
      </main>
    </>
  );
}

function UserButton() {
  const { currentUser, setCurrentUser } = useCurrentUser();
  const navigate = useNavigate();

  const handleLogout = () => {
    // TODO: also call the logout endpoint
    // Clear cookies
    Cookies.remove("quizzer_session_id");

    // Clear user context
    setCurrentUser(null);

    // Navigate to login page
    navigate("/login");
  };

  if (!currentUser) {
    return (
      <Button variant="link">
        <Link to="/login">Login</Link>
      </Button>
    );
  }

  return (
    <Popover>
      <PopoverTrigger asChild>
        <Button variant="outline">{currentUser.username}</Button>
      </PopoverTrigger>
      <PopoverContent className="w-80">
        <div className="grid gap-4">
          <div className="space-y-2">
            <h4 className="font-medium leading-none">User</h4>
            <p className="text-sm text-muted-foreground">
              {currentUser.username}
            </p>
          </div>

          <div className="space-y-2">
            <h4 className="font-medium leading-none">Signup Date</h4>
            <p className="text-sm text-muted-foreground">
              {currentUser.signupDate}
            </p>
          </div>

          <Button onClick={handleLogout}>Log out</Button>
        </div>
      </PopoverContent>
    </Popover>
  );
}
