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
    stdout: "piped",
    stderr: "piped",
}).spawn()
const firstByteSeen = Promise.withResolvers<void>()
void child.stdout.pipeThrough(new TransformStream({
    transform(chunk, controller) {
        firstByteSeen.resolve()
        controller.enqueue(chunk)
    },
    cancel(reason) {
        firstByteSeen.reject(reason)
    }
})).pipeTo(Deno.stdout.writable, { preventAbort: true, preventCancel: true, preventClose: true })
void child.stderr.pipeThrough(new TransformStream({
    transform(chunk, controller) {
        firstByteSeen.resolve()
        controller.enqueue(chunk)
    },
    cancel(reason) {
        firstByteSeen.reject(reason)
    }
})).pipeTo(Deno.stderr.writable, { preventAbort: true, preventCancel: true, preventClose: true })
child.unref()

export default {
    async fetch(request) {
        const requestURL = new URL(request.url)
        await firstByteSeen.promise;
        const newRequestURL = new URL(requestURL.pathname + requestURL.search, `http://[::1]:${port}`)
        return await fetch(newRequestURL, request)
    }
} satisfies Deno.ServeDefaultExport

