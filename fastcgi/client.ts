export default class Client {
    #connectOptions
    constructor(options: Parameters<typeof Deno.connect>[0]) {
        this.#connectOptions = options
    }

    async fetch(input: RequestInfo | URL, init?: RequestInit): Promise<Response> {
        const request = input instanceof Request && init === undefined ? input : new Request(input, init)
        const urlObject = new URL(request.url)

        const params: Record<string, string> = {
            REQUEST_METHOD: request.method,
            REQUEST_SCHEME: urlObject.protocol.replace(/:$/, ""),
            HTTP_HOST: urlObject.hostname,
            REQUEST_URI: urlObject.pathname + urlObject.search,
            QUERY_STRING: urlObject.search.replace(/^\?/, ""),
            SERVER_SOFTWARE: `fastcgi-client/1.0`,
        }
        for (const [name, value] of request.headers) {
            if (name === "content-type") {
                params.CONTENT_TYPE = value;
            } else if (name === "content-length") {
                params.CONTENT_LENGTH = value;
            } else {
                params[`HTTP_${name.replaceAll("-", "_").toUpperCase()}`] = value;
            }
        }

        const connection = await Deno.connect(this.#connectOptions)
        await connection.write()
    }
}