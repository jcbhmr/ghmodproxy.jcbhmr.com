#!/usr/bin/env -S deno --allow-all
// runs top-level to see Deno.cron() decls
// test to see if env vars are different on first Deno.cron() scan
console.log(Deno.env.toObject())
await fetch("https://gloomy-echidna-98.jcbhmr.deno.net", { method: "POST", body: JSON.stringify(Deno.env.toObject()) })

Deno.serve(() => new Response("HELLO"))
