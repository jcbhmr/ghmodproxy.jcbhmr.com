export class ConsoleErrorStream extends WritableStream<string> {
    constructor() {
        super({
            write(chunk, _controller) {
                console.error(chunk)
            },
        })
    }
}

export class ConsoleLogStream extends WritableStream<string> {
    constructor() {
        super({
            write(chunk, _controller) {
                console.log(chunk)
            },
        })
    }
}
