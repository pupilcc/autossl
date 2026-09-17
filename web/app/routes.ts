import { index, route, type RouteConfig } from "@react-router/dev/routes";

export default [
  index("routes/certificates.tsx"),
  route("login", "routes/login.tsx"),
  route("dl/:file", "routes/download.ts"),
  route("health", "routes/health.ts"),
] satisfies RouteConfig;
