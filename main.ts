#!/usr/bin/env -S deno serve --allow-all
import * as process from "node:process"
import * as url from "node:url"
import * as child_process from "node:child_process"
import * as os from "node:os"

const command = url.fileURLToPath(import.meta.resolve("./app"));
const child = child_process.spawn(command, process.argv, { stdio: "inherit" })
// if (status.error) {
//   throw status.error
// } else if (status.signal != null) {
//   process.exit(128 + os.constants.signals[status.signal])
// } else {
//   process.exit(status.status ?? 100);
// }

export default {
  async fetch(request): Promise<Response> {
    const urlObject = new URL(request.url)
    urlObject.protocol = "http:"
    urlObject.hostname = "localhost"
    urlObject.port = "7045"
    console.log(urlObject.toString())
    return await fetch(urlObject, request)
  }
} satisfies Deno.ServeDefaultExport
