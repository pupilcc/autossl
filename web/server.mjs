export * from "./build/server/index.js";

const consoleUrl = process.env.PUBLIC_CONSOLE_URL;
export const allowedActionOrigins = consoleUrl ? [new URL(consoleUrl).host] : false;
