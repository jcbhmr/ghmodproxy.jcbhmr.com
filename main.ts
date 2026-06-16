#!/usr/bin/env -S deno --allow-all
import { TextLineStream } from "@std/streams";
import ConsoleStream from "./console_stream.ts";

const child = new Deno.Command(new URL(import.meta.resolve("./app")), {
    args: [],
    env: {},
    stdin: "inherit",
    stdout: "piped",
    stderr: "piped",
}).spawn()
void child.stdout.pipeThrough(new TextDecoderStream()).pipeThrough(new TextLineStream()).pipeTo(new ConsoleStream())
void child.stderr.pipeThrough(new TextDecoderStream()).pipeThrough(new TextLineStream()).pipeTo(new ConsoleStream("error"))
Deno.addSignalListener("SIGINT", () => {
    child.kill("SIGINT")
})
Deno.addSignalListener("SIGTERM", () => {
    child.kill("SIGTERM")
})

