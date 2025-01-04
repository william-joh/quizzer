import { Outlet } from "react-router";

export function CenterLayout() {
  return (
    <div className="container mx-auto py-8 px-8 max-w-[1000px]">
      <Outlet />
    </div>
  );
}
