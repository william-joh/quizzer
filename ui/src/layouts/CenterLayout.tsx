import { Outlet } from "react-router";

export function CenterLayout() {
  return (
    <div className="container mx-auto max-w-[1000px]">
      <Outlet />
    </div>
  );
}
