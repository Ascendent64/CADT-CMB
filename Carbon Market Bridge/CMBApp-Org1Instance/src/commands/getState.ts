import { Contract } from '@hyperledger/fabric-gateway';

export async function getState(contract: Contract, tradingPlatformID: string): Promise<any> {
    console.log(`Fetching state for TradingPlatformID: ${tradingPlatformID}`);
    const result = await contract.evaluateTransaction('GetState', tradingPlatformID);

    const decoder = new TextDecoder('utf-8');
    const readableResult = decoder.decode(result);

    console.log(`State fetched, result is: ${readableResult}`);
    return JSON.parse(readableResult);
}
