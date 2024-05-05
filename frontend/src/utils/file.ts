export class FileReaderEx extends FileReader {
    constructor() {
        super();
    }

    private readAs(blob: Blob, ctx: 'readAsArrayBuffer' | 'readAsDataURL' | 'readAsText') {
        return new Promise((res, rej) => {
            super.addEventListener("load", ({target}) => res(target?.result));
            super.addEventListener("error", ({target}) => rej(target?.error));

            super[ctx](blob);
        });
    }

    readAsArrayBuffer(blob: Blob) {
        return this.readAs(blob, "readAsArrayBuffer");
    }

    readAsDataURL(blob: Blob) {
        return this.readAs(blob, "readAsDataURL");
    }

    readAsText(blob: Blob) {
        return this.readAs(blob, "readAsText");
    }
}
