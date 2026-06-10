export class Header<T extends ArrayBufferLike = ArrayBufferLike> implements ArrayBufferView<T> {
    static BEGIN_REQUEST = 0
    static ABORT_REQUEST = 1
    static END_REQUEST = 2
    static PARAMS = 4
    static STDIN = 5
    static STDOUT = 6
    static STDERR = 7
    static DATA = 8
    static GET_VALUES = 9
    static GET_VALUES_RESULT = 10
    static UNKNOWN_TYPE = 11

    static BYTE_LENGTH = 8
    static from(options: {
        version: 1
        type: number
        requestID: number
        contentLength: number
        paddingLength: number
    }): Header<ArrayBuffer> {
        const self = new Header(new ArrayBuffer(Header.BYTE_LENGTH))
        self.version = options.version
        self.type = options.type
        self.requestID = options.requestID
        self.contentLength = options.contentLength
        self.paddingLength = options.paddingLength
        return self
    }

    #dataView: DataView<T>
    constructor(buffer: T, byteOffset?: number) {
        this.#dataView = new DataView(buffer, byteOffset, Header.BYTE_LENGTH)
    }
    get buffer(): T {
        return this.#dataView.buffer
    }
    get byteOffset(): number {
        return this.#dataView.byteOffset
    }
    get byteLength(): number {
        return this.#dataView.byteLength
    }

    get version(): number {
        return this.#dataView.getUint8(0)
    }
    set version(v: number) {
        this.#dataView.setUint8(0, v)
    }

    get type(): number {
        return this.#dataView.getUint8(1)
    }
    set type(v: number) {
        this.#dataView.setUint8(1, v)
    }

    get requestID(): number {
        return this.#dataView.getUint16(2, false)
    }
    set requestID(v: number) {
        this.#dataView.setUint16(2, v, false)
    }

    get contentLength(): number {
        return this.#dataView.getUint16(4, false)
    }
    set contentLength(v: number) {
        this.#dataView.setUint16(4, v, false)
    }

    get paddingLength(): number {
        return this.#dataView.getUint8(6)
    }
    set paddingLength(v: number) {
        this.#dataView.setUint8(6, v)
    }
}

export default class Record<T extends ArrayBufferLike = ArrayBufferLike> implements ArrayBufferView<T> {
    #dataView: DataView<T>
    #header: Header
    #content: Uint8Array<T> | undefined
    #padding: Uint8Array<T> | undefined
    constructor(buffer: T, byteOffset?: number) {
        this.#header = new Header(buffer, byteOffset)
        this.#dataView = new DataView(buffer, byteOffset, Header.BYTE_LENGTH + this.#header.contentLength + this.#header.paddingLength)
    }
    get buffer(): T {
        return this.#dataView.buffer
    }
    get byteOffset(): number {
        return this.#dataView.byteOffset
    }
    get byteLength(): number {
        return this.#dataView.byteLength
    }

    get header(): Header {
        return this.#header
    }

    get content(): Uint8Array<T> {
        this.#content ??= new Uint8Array(this.#dataView.buffer, this.#dataView.byteOffset + Header.BYTE_LENGTH, this.#header.contentLength)
        return this.#content
    }

    get padding(): Uint8Array<T> {
        this.#padding ??= new Uint8Array(this.#dataView.buffer, this.#dataView.byteOffset + Header.BYTE_LENGTH + this.#header.contentLength, this.#header.paddingLength)
        return this.#padding
    }
}

export class BeginRequestBody<T extends ArrayBufferLike = ArrayBufferLike> implements ArrayBufferView<T> {
    static KEEP_CONN = 1
    static RESPONDER = 1
    static AUTHORIZER = 2
    static FILTER = 3

    static BYTE_LENGTH = 8
    static from(options: {
        role: number
        flags: number
    }): BeginRequestBody<ArrayBuffer> {
        const self = new BeginRequestBody(new ArrayBuffer(BeginRequestBody.BYTE_LENGTH))
        self.role = options.role
        self.flags = options.flags
        return self
    }

    #dataView: DataView<T>
    constructor(buffer: T, byteOffset?: number) {
        this.#dataView = new DataView(buffer, byteOffset, BeginRequestBody.BYTE_LENGTH)
    }
    get buffer(): T {
        return this.#dataView.buffer
    }
    get byteOffset(): number {
        return this.#dataView.byteOffset
    }
    get byteLength(): number {
        return this.#dataView.byteLength
    }

    get role(): number {
        return this.#dataView.getUint16(0, false)
    }
    set role(v: number) {
        this.#dataView.setUint16(0, v, false)
    }

    get flags(): number {
        return this.#dataView.getUint8(2)
    }
    set flags(v: number) {
        this.#dataView.setUint8(2, v)
    }
}

export class BeginRequestRecord<T extends ArrayBufferLike = ArrayBufferLike> implements ArrayBufferView<T> {
    static BYTE_LENGTH = Header.BYTE_LENGTH + BeginRequestBody.BYTE_LENGTH
    static from(options: {
        header: Parameters<typeof Header.from>[0],
        body: Parameters<typeof BeginRequestBody.from>[0],
    }): BeginRequestRecord<ArrayBuffer> {
        const self = new BeginRequestRecord(new ArrayBuffer(BeginRequestRecord.BYTE_LENGTH))
        const header = new Header(self.buffer, 0)
        header.version = options.header.version
        header.type = options.header.type
        header.requestID = options.header.requestID
        header.contentLength = options.header.contentLength
        header.paddingLength = options.header.paddingLength
        const body = new BeginRequestBody(self.buffer, Header.BYTE_LENGTH)
        body.role = options.body.role
        body.flags = options.body.flags
        return self
    }

    #dataView: DataView<T>
    #header: Header | undefined
    #body: BeginRequestBody | undefined
    constructor(buffer: T, byteOffset?: number) {
        Uint32Array.BYTES_PER_ELEMENT
        this.#dataView = new DataView(buffer, byteOffset, BeginRequestRecord.BYTE_LENGTH)
    }
    get buffer(): T {
        return this.#dataView.buffer
    }
    get byteOffset(): number {
        return this.#dataView.byteOffset
    }
    get byteLength(): number {
        return this.#dataView.byteLength
    }

    get header(): Header {
        this.#header ??= new Header(this.#dataView.buffer, this.#dataView.byteOffset)
        return this.#header
    }

    get body(): BeginRequestBody {
        this.#body ??= new BeginRequestBody(this.#dataView.buffer, this.#dataView.byteOffset + Header.BYTE_LENGTH)
        return this.#body
    }
}