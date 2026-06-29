#!/usr/bin/env -S deno serve --allow-all
import { routeRadix, type Route } from "@std/http/unstable-route"
import { basename } from "@std/path/posix"
import { Octokit } from "octokit"
import { HttpRangeReader, ZipReader, type FileEntry } from "@zip.js/zip.js"

const octokit = new Octokit({
    auth: Deno.env.get("GITHUB_TOKEN"),
})

function unescapePath(input: string): string {
    return input.replaceAll(/!([a-z])/g, (match0, match1: string) => match1.toUpperCase())
}

function unescapeVersion(input: string): string {
    return input.replaceAll(/!([a-z])/g, (match0, match1: string) => match1.toUpperCase())
}

const routes = [
    {
        method: "GET",
        pattern: new URLPattern({ pathname: "/" }),
        async handler(request, params, info) {
            return new Response("Hello, Alan Turing!")
        },
    },
    {
        method: "GET",
        pattern: new URLPattern({ pathname: "/healthz" }),
        async handler(request, params, info) {
            return Response.json({ status: "pass" }, {
                headers: {
                    "Content-Type": "application/health+json"
                }
            })
        },
    },
    {
        method: "GET",
        pattern: new URLPattern({ pathname: "/:owner/:repo/:moduleEscaped(.*?)/@v/list" }),
        async handler(request, params, info) {
            const module = unescapePath(params.pathname.groups.moduleEscaped!)
            const searchParams = new URL(request.url).searchParams
            const subdirectory = searchParams.get("subdirectory") ?? ""
            const prefix = subdirectory !== "" ? `${subdirectory}/` : ""
            const assetName = basename(module) + ".zip"

            const releases = await octokit.paginate(octokit.rest.repos.listReleases, {
                owner: params.pathname.groups.owner!,
                repo: params.pathname.groups.repo!,
            })
            const versions = releases
                .filter(x => x.tag_name.startsWith(prefix))
                .filter(x => x.assets.find(y => y.name === assetName))
                .map(x => x.tag_name.slice(prefix.length))
            if (versions.length === 0) {
                return new Response(null, { status: 404 })
            }
            return new Response(versions.join("\n") + "\n", {
                headers: {
                    "Content-Type": "text/plain; charset=utf-8"
                }
            })
        }
    },
    {
        method: "GET",
        pattern: new URLPattern({ pathname: "/:owner/:repo/:moduleEscaped(.*?)/@latest" }),
        async handler(request, params, info) {
            const module = unescapePath(params.pathname.groups.moduleEscaped!)
            const searchParams = new URL(request.url).searchParams
            const subdirectory = searchParams.get("subdirectory") ?? ""
            const prefix = subdirectory !== "" ? `${subdirectory}/` : ""
            const assetName = basename(module) + ".zip"

            const releases = await octokit.paginate(octokit.rest.repos.listReleases, {
                owner: params.pathname.groups.owner!,
                repo: params.pathname.groups.repo!,
            })
            const release = releases
                .filter(x => x.tag_name.startsWith(prefix))
                .find(x => x.assets.find(y => y.name === assetName))
            if (!release) {
                return new Response(null, { status: 404 })
            }
            return Response.json({
                Version: release.tag_name.slice(prefix.length),
                Time: release.published_at ?? undefined,
            })
        }
    },
    {
        method: "GET",
        pattern: new URLPattern({ pathname: "/:owner/:repo/:moduleEscaped(.*?)/@v/:versionEscaped(.*?).info" }),
        async handler(request, params, info) {
            const module = unescapePath(params.pathname.groups.moduleEscaped!)
            const version = unescapeVersion(params.pathname.groups.versionEscaped!)
            const searchParams = new URL(request.url).searchParams
            const subdirectory = searchParams.get("subdirectory") ?? ""
            const prefix = subdirectory !== "" ? `${subdirectory}/` : ""
            const assetName = basename(module) + ".zip"

            const { data: release } = await octokit.rest.repos.getReleaseByTag({
                owner: params.pathname.groups.owner!,
                repo: params.pathname.groups.repo!,
                tag: prefix + version,
            })
            if (!release.assets.find(x => x.name === assetName)) {
                return new Response(null, { status: 404 })
            }
            return Response.json({
                Version: release.tag_name.slice(prefix.length),
                Time: release.published_at ?? undefined,
            })
        },
    },
    {
        method: "GET",
        pattern: new URLPattern({ pathname: "/:owner/:repo/:moduleEscaped(.*?)/@v/:versionEscaped(.*?).mod" }),
        async handler(request, params, info) {
            const module = unescapePath(params.pathname.groups.moduleEscaped!)
            const version = unescapeVersion(params.pathname.groups.versionEscaped!)
            const searchParams = new URL(request.url).searchParams
            const subdirectory = searchParams.get("subdirectory") ?? ""
            const prefix = subdirectory !== "" ? `${subdirectory}/` : ""
            const assetName = basename(module) + ".zip"

            const { data: release } = await octokit.rest.repos.getReleaseByTag({
                owner: params.pathname.groups.owner!,
                repo: params.pathname.groups.repo!,
                tag: prefix + params.pathname.groups.version!,
            })
            const asset = release.assets.find(x => x.name === assetName)
            if (!asset) {
                return new Response(null, { status: 404 })
            }

            const zip = new ZipReader(new HttpRangeReader(asset.browser_download_url))
            let entry: FileEntry | undefined;
            for await (const e of zip.getEntriesGenerator()) {
                if (e.directory) {
                    continue
                }
                if (e.filename === `${module}@${version}/go.mod`) {
                    entry = e
                    break;
                }
            }
            if (!entry) {
                return new Response(null, { status: 404 })
            }

            const stream = new TransformStream()
            entry.getData(stream.writable).finally(async () => {
                await zip.close()
            })
            return new Response(stream.readable, {
                headers: {
                    "Content-Type": "text/plain; charset=utf-8",
                    "Content-Length": BigInt(entry.uncompressedSize).toString()
                }
            })
        },
    },
    {
        method: "GET",
        pattern: new URLPattern({ pathname: "/:owner/:repo/:moduleEscaped(.*?)/@v/:versionEscaped(.*?).zip" }),
        async handler(request, params, info) {
            const module = unescapePath(params.pathname.groups.moduleEscaped!)
            const version = unescapeVersion(params.pathname.groups.versionEscaped!)
            const searchParams = new URL(request.url).searchParams
            const subdirectory = searchParams.get("subdirectory") ?? ""
            const prefix = subdirectory !== "" ? `${subdirectory}/` : ""
            const assetName = basename(module) + ".zip"

            const { data: release } = await octokit.rest.repos.getReleaseByTag({
                owner: params.pathname.groups.owner!,
                repo: params.pathname.groups.repo!,
                tag: prefix + params.pathname.groups.version!,
            })
            const asset = release.assets.find(x => x.name === assetName)
            if (!asset) {
                return new Response(null, { status: 404 })
            }

            return Response.redirect(asset.browser_download_url, 307)
        },
    },
] satisfies Route[]

const router = routeRadix(routes, async (request) => {
    return new Response(null, { status: 404 })
})

export default {
    async fetch(request, info) {
        return await router(request, info)
    }
} satisfies Deno.ServeDefaultExport

