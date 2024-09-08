export function assertDefined<T>(value: T | undefined, message: string | (() => string)): T {
    if (value == undefined) {
        throw new Error(typeof message === 'string' ? message : message());
    }
    return value;
}

export { newIdentity, newSigner } from '../identity';  
