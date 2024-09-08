export class ExpectedError extends Error {
    constructor(message?: string) {
        super(message);
        this.name = ExpectedError.name;
    }
}