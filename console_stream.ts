export default class ConsoleStream extends WritableStream<string> {
    constructor(methodName: "log" | "error" | "warn" | "info" | "debug" = "log") {
        super({
            write(chunk, _controller) {
                console[methodName](chunk)
            },
        })
    }
}
