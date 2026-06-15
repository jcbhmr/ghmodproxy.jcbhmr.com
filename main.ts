#!/usr/bin/env -S deno --allow-all
import { TextLineStream } from "@std/streams"

const child = new Deno.Command(new URL(import.meta.resolve("./app")), {
    args: [],
    // env: {
    //     PORT: "0"
    // },
    stdin: "inherit",
    stdout: "inherit",
    stderr: "inherit",
}).spawn()
child.unref()
// let portCached: string | undefined
// async function getPort(): Promise<string> {
//     if (portCached == null) {
//         let firstLine: string | undefined
//         for await (const line of child.stdout.pipeThrough(new TextDecoderStream()).pipeThrough(new TextLineStream())) {
//             firstLine = line
//             break;
//         }
//         if (firstLine == null) {
//             throw new DOMException("child.stdout had no first line", "SyntaxError")
//         }
//         portCached = firstLine
//     }
//     return portCached
// }


// export default {
//     async fetch(request) {
//         const requestURL = new URL(request.url)
//         const port = await getPort()
//         const newRequestURL = new URL(requestURL.pathname + requestURL.search, `http://localhost:${port}`)
//         console.log("%s => %s", requestURL, newRequestURL)
//         return await fetch(newRequestURL, request)
//     }
// } satisfies Deno.ServeDefaultExport

