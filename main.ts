#!/usr/bin/env -S deno --allow-all
import * as process from "node:process"
import * as url from "node:url"
import * as child_process from "node:child_process"
import * as os from "node:os"

const command = url.fileURLToPath(import.meta.resolve("./app"));
const status = child_process.spawnSync(command, process.argv, { stdio: "inherit" })
if (status.error) {
  throw status.error
} else if (status.signal != null) {
  process.exit(128 + os.constants.signals[status.signal])
} else {
  process.exit(status.status ?? 100);
}
