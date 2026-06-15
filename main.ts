#!/usr/bin/env -S deno --allow-all
import * as http from "node:http"
import * as path from "node:path"
import cgi from "cgi"

const server = http.createServer(cgi(path.resolve("./app")))
server.listen();
