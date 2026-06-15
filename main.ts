#!/usr/bin/env -S deno serve --allow-all
import { getAvailablePort } from "@std/net"

let port: number | undefined

export default {
    async fetch(request, info) {
        // if (port == null) {
        //     port = getAvailablePort({ preferredPort: 8000 })
        //     const child = new Deno.Command(new URL(import.meta.resolve("./app")), {
        //         args: [],
        //         env: {
        //             HOST: "[::]",
        //             PORT: port.toString()
        //         },
        //         stdin: "inherit",
        //         stdout: "inherit",
        //         stderr: "inherit",
        //     }).spawn()
        //     child.unref()
        //     await new Promise<void>((resolve, reject) => {
        //         const l = () => {
        //             resolve()
        //             Deno.removeSignalListener("SIGUSR2", l)
        //         }
        //         child.status.then(reject, reject).finally(() => {
        //             Deno.removeSignalListener("SIGUSR2", l)
        //         })
        //         Deno.addSignalListener("SIGUSR2", l)
        //     })
        // }
        // const requestURL = new URL(request.url)
        // const newRequestURL = new URL(requestURL.pathname + requestURL.search, `http://[::1]:${port}`)
        // return await fetch(newRequestURL, request)
        return new Response("HELLO")
    },
    onListen(localAddr) {
        console.log({ localAddr })
    }
} satisfies Deno.ServeDefaultExport

