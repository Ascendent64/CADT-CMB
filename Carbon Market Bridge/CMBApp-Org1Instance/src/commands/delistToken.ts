import { Contract } from '@hyperledger/fabric-gateway';

export async function delistToken(contract: Contract, args: string[]): Promise<void> {
    const tradingPlatformID = args[0];
    console.log(`Delisting token with TradingPlatformID: ${tradingPlatformID}`);
    const result = await contract.submitTransaction('DelistToken', tradingPlatformID);

    const asciiArray = Array.from(new Uint8Array(result));
    const readableResult = String.fromCharCode(...asciiArray);

    console.log(`Transaction has been submitted, result is: ${readableResult}`);
}
