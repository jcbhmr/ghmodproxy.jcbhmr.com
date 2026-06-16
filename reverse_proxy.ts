export interface ReverseProxyOptions {
    rewrite?(request: Request): Request | PromiseLike<Request>
    fetch?(input: RequestInfo | URL, init?: RequestInit): Promise<Response>
}

export default class ReverseProxy {
    #rewrite?: (request: Request) => Request | PromiseLike<Request>
    #fetch: (input: RequestInfo | URL, init?: RequestInit) => Promise<Response>
    constructor(options: ReverseProxyOptions = {}) {
        const { rewrite, fetch = globalThis.fetch } = options
        this.#rewrite = rewrite
        this.#fetch = fetch
    }
}