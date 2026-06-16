#!/usr/bin/env -S deno --allow-all
import { TextLineStream } from "@std/streams";

const server = Deno.serve({
    onListen(localAddr) {
        console.log("Listening for HTTP on %s", `${localAddr.transport}:${localAddr.hostname}:${localAddr.port}`)
    }
}, () => Response.error())
await server.shutdown()

const child = new Deno.Command(new URL(import.meta.resolve("./app")), {
    args: [],
    env: {},
    stdin: "inherit",
    stdout: "piped",
    stderr: "piped",
}).spawn()
void child.stdout.pipeThrough(new TextDecoderStream()).pipeThrough(new TextLineStream()).pipeTo(new WritableStream({
    write(chunk, _controller) {
        console.log(chunk)
    }
}))
void child.stderr.pipeThrough(new TextDecoderStream()).pipeThrough(new TextLineStream()).pipeTo(new WritableStream({
    write(chunk, _controller) {
        console.error(chunk)
    }
}))
Deno.addSignalListener("SIGINT", () => {
    child.kill("SIGINT")
})
Deno.addSignalListener("SIGTERM", () => {
    child.kill("SIGTERM")
})

