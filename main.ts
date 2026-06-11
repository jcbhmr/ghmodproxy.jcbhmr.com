#!/usr/bin/env -S deno serve --allow-all
import * as childProcess from "node:child_process"
import { fcgi } from "fcgi"

const child = childProcess.spawn("./app")

export default {
  async fetch(request): Promise<Response> {
    const response = await fcgi.fetch({
      addr: "localhost:7045",
    }, request)
    return new Response(response.body, {
      status: response.status,
      statusText: response.statusText,
      headers: response.headers,
    })
  }
} satisfies Deno.ServeDefaultExport
