export * from "./build/server/index.js";

const publicUrl = process.env.DOMAIN;
export const allowedActionOrigins = publicUrl ? [new URL(publicUrl).host] : false;
