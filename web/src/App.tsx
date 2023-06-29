import { createBrowserRouter, RouterProvider } from "react-router-dom";
import { DashboardPage } from "./pages/dashboard";
import { ServerPage } from "./pages/server";

import React from "react";

export const App: React.FC<{}> = () => {
  const router = createBrowserRouter([
    {
      path: "/",
      element: <DashboardPage />,
    },
    {
      path: "/server/:id",
      element: <ServerPage />,
    },
    {
      path: "/node/:id",
      element: <ServerPage />,
    },
    {
      path: "about",
      element: <div>About</div>,
    },
  ]);

  return <RouterProvider router={router} />;
};
