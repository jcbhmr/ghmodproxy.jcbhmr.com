#!/usr/bin/env -S deno --allow-all
import { execa } from "npm:execa@9.6.1"
import * as path from "node:path"
import * as os from "node:os"

Deno.exit(0)

await execa({
    verbose: "short",
    stdio: "inherit",
    env: {
        GOLANG_VERSION: undefined,
        GOROOT: undefined,
        GOTOOLCHAIN: undefined,
        GOPATH: path.join(os.homedir(), "gopath"),
        CGO_ENABLED: "0",
    },
})`${path.join(os.homedir(), "go")}/bin/go build .`
