import "./styles.css";

import { RouterProvider, createRouter } from "@tanstack/react-router";

import ReactDOM from "react-dom/client";
// src/main.tsx
import { StrictMode } from "react";
import { routeTree } from "./routeTree.gen"; // auto-generated by `tsr`

// Create the router
const router = createRouter({
  routeTree,
  defaultPreload: "intent",
  context: {},
  defaultPreloadStaleTime: 0,
  scrollRestoration: true,
  defaultStructuralSharing: true,
});

// Declare module for type-safe hooks like `useNavigate()`
declare module "@tanstack/react-router" {
  interface Register {
    router: typeof router;
  }
}

// Mount the app
const rootElement = document.getElementById("app");
if (rootElement) {
  const root = ReactDOM.createRoot(rootElement);
  root.render(
    <StrictMode>
      <RouterProvider router={router} />
    </StrictMode>
  );
}
