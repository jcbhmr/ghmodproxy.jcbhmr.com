#!/usr/bin/env -S deno --allow-all
import { TextLineStream } from "@std/streams"
// await fetch("https://gloomy-echidna-98.jcbhmr.deno.net/" + Deno.env.get("DENO_DEPLOYMENT_ID"))

let port: string | undefined
async function startChild(): Promise<string> {
    if (port == null) {
        const controller = new AbortController()
        const child = new Deno.Command(new URL(import.meta.resolve("./app")), {
            args: [],
            env: {
                PORT: "0"
            },
            stdin: "inherit",
            stdout: "piped",
            stderr: "inherit",
            signal: controller.signal
        }).spawn()
        child.unref()
        globalThis.addEventListener("unload", (_event) => {
            controller.abort()
        }, { passive: true })
        for await (const line of child.stdout.pipeThrough(new TextDecoderStream()).pipeThrough(new TextLineStream())) {
            port = line;
            break;
        }
        if (port == null) {
            throw new TypeError("port is null")
        }
    }
    return port
}


export default {
    async fetch(request) {
        const requestURL = new URL(request.url)
        const port = await startChild()
        const newRequestURL = new URL(requestURL.pathname + requestURL.search, `http://localhost:${port}`)
        console.log("%s => %s", requestURL, newRequestURL)
        return await fetch(newRequestURL, request)
    }
} satisfies Deno.ServeDefaultExport

