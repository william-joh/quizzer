import { useCurrentUser } from "@/contexts/userContext";
import { Outlet } from "react-router";

export function AuthLayout() {
  const { currentUser } = useCurrentUser();

  if (!currentUser) {
    return (
      <div className="container mx-auto py-8 px-8 max-w-[1000px]">
        <div className="text-center">
          <h1 className="text-3xl font-bold">
            You must be logged in to view this page
          </h1>
        </div>
      </div>
    );
  }

  return <Outlet />;
}
