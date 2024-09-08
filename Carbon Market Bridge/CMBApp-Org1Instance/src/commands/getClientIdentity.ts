import { Contract } from '@hyperledger/fabric-gateway';

export async function getClientIdentity(contract: Contract): Promise<string> {
    console.log('Fetching client identity');
    const result = await contract.evaluateTransaction('GetClientIdentity');

    const asciiArray = Array.from(new Uint8Array(result));
    const readableResult = String.fromCharCode(...asciiArray);

    console.log(`Client identity fetched: ${readableResult}`);
    return readableResult;
}
