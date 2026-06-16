#!/usr/bin/env -S deno --allow-all
import { execa } from "npm:execa@9.6.1"
import * as streamPromises from "node:stream/promises"
import * as fs from "node:fs"
import * as fsPromises from "node:fs/promises"
import * as path from "node:path"
import * as os from "node:os"

const response = await fetch("https://go.dev/dl/go1.26.4.linux-amd64.tar.gz")
if (response.status !== 200) {
    throw new DOMException(`${response.url} ${response.status} ${response.statusText}`)
}
if (!response.body) {
    throw new DOMException(`${response.url} no body`)
}
await streamPromises.pipeline(
    response.body,
    fs.createWriteStream(path.join(os.homedir(), "go1.26.4.linux-amd64.tar.gz"))
)
await fsPromises.rm(path.join(os.homedir(), "go"), { recursive: true, force: true })
await execa({ verbose: "short", stdio: "inherit", cwd: os.homedir() })`tar -xzf go1.26.4.linux-amd64.tar.gz`
await fsPromises.rm(path.join(os.homedir(), "go1.26.4.linux-amd64.tar.gz"))
