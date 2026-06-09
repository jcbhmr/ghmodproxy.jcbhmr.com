#!/usr/bin/env -S deno --allow-all
import { execa } from "npm:execa@9.6.1"

await execa({
    verbose: "short",
    stdio: "inherit",
    env: {
        GOLANG_VERSION: undefined,
        GOROOT: undefined,
        GOTOOLCHAIN: undefined,
        GOPATH: undefined,
    },
})`./go/bin/go build .`
