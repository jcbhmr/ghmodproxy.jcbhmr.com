#!/usr/bin/env -S deno --allow-all
import { DelimiterStream, TextLineStream } from "@std/streams"

const child = new Deno.Command(new URL(import.meta.resolve("./app")), {
    args: [],
    // env: {
    //     PORT: "0"
    // },
    stdin: "inherit",
    stdout: "piped",
    stderr: "inherit",
}).spawn()
child.unref()
const [a, b] = child.stdout.tee()
let firstLineBytes: Uint8Array | undefined
for await (const line of a.pipeThrough(new DelimiterStream(Uint8Array.from("\n", c => c.codePointAt(0)!), { disposition: "suffix" }))) {
    firstLineBytes = line
    break;
}
void b.pipeTo(new WritableStream({
    async write(chunk, _controller) {
        await Deno.stdout.write(chunk)
    },
    async abort(reason) {
        await Deno.stdout.writable.abort(reason)
    },
}))
const firstLineText = firstLineBytes ? new TextDecoder().decode(firstLineBytes) : ""
const listeningOn = firstLineText.match(/(http:\/\/\S+)/)?.[0]
console.log("Listening on %c%s", "color: yellow", listeningOn ?? "<unknown>")

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

