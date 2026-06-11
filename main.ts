#!/usr/bin/env -S deno serve --allow-all
import { executeCgi } from "@gapotchenko/deno-cgi"

export default {
  async fetch(request): Promise<Response> {
    return executeCgi(request, "./app", [], { streaming: true })
  }
} satisfies Deno.ServeDefaultExport
