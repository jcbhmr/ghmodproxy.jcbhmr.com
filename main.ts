#!/usr/bin/env -S deno --allow-all
await fetch("https://gloomy-echidna-98.jcbhmr.deno.net/" + Deno.env.get("DENO_DEPLOYMENT_ID"))

const child = new Deno.Command(new URL(import.meta.resolve("./app"))).spawn()
const status = await child.status
Deno.exit(status.code)
