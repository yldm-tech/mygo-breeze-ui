import { createRouter } from "@tanstack/react-router";
import { HomeRoute } from "./routes/home";
import { CliMcpRoute } from "./routes/cli-mcp";
import { DocsRoute } from "./routes/docs";
import { ComponentsRoute, ComponentDocRoute } from "./routes/components-docs";
import { RootRoute } from "./routes/root";

const routeTree = RootRoute.addChildren([HomeRoute, DocsRoute, ComponentsRoute, ComponentDocRoute, CliMcpRoute]);
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
