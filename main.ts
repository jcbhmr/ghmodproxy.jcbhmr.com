#!/usr/bin/env -S deno serve --allow-all
import { getAvailablePort } from "@std/net"

const port = getAvailablePort({ preferredPort: 8000 })

const child = new Deno.Command(new URL(import.meta.resolve("./app")), {
    args: [],
    env: {
        HOST: "[::]",
        PORT: port.toString()
    },
    stdin: "inherit",
    stdout: "inherit",
    stderr: "inherit",
}).spawn()
child.unref()

let ready: Promise<void> | undefined
function getReady(): Promise<void> {
    ready ??= (async () => {
        const signal = AbortSignal.timeout(500)
        while (true) {
            signal.throwIfAborted()
            try {
                await fetch(`http://[::1]:${port}`, { method: "HEAD", signal }).then(() => false, () => true)
            } catch (error) {
                if (error instanceof TypeError && error.message.includes("os error 111")) {
                    continue;
                } else {
                    throw error
                }
            }
        }
    })()
    return ready
}

export default {
    async fetch(request) {
        await getReady()
        const requestURL = new URL(request.url)
        const newRequestURL = new URL(requestURL.pathname + requestURL.search, `http://[::1]:${port}`)
        return await fetch(newRequestURL, request)
    }
} satisfies Deno.ServeDefaultExport

