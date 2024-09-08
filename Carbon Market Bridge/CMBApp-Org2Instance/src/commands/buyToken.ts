import { Contract } from '@hyperledger/fabric-gateway';

export async function buyToken(contract: Contract, args: string[]): Promise<void> {
    const tradingPlatformID = args[0];
    console.log(`Buying token with TradingPlatformID: ${tradingPlatformID}`);
    const result = await contract.submitTransaction('BuyToken', tradingPlatformID);

    const asciiArray = Array.from(new Uint8Array(result));
    const readableResult = String.fromCharCode(...asciiArray);

    console.log(`Transaction has been submitted, result is: ${readableResult}`);
}
