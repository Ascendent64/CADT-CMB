import { Contract } from '@hyperledger/fabric-gateway';

export async function listAllTokens(contract: Contract, args: string[]): Promise<Uint8Array> {
    console.log('Listing all tokens');
    const result = await contract.evaluateTransaction('ListAllTokens');
    return result;
}
