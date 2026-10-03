import { createRouter } from "@tanstack/react-router";
import { HomeRoute } from "./routes/home";
import { ApiRoute } from "./routes/api";
import { RootRoute } from "./routes/root";

const routeTree = RootRoute.addChildren([HomeRoute, ApiRoute]);
export const router = createRouter({
  routeTree,
  basepath: import.meta.env.BASE_URL,
  defaultPreload: "intent",
});

declare module "@tanstack/react-router" {
  interface Register {
    router: typeof router;
  }
}
