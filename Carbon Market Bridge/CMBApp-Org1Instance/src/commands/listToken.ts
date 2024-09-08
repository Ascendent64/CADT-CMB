import { Contract } from '@hyperledger/fabric-gateway';

export async function listToken(contract: Contract, args: string[]): Promise<void> {
    const [tradingPlatformID, price, tokenName, matches = '[]'] = args;
    console.log(`Listing token with TradingPlatformID: ${tradingPlatformID}, Price: ${price}, Token Name: ${tokenName}, Matches: ${matches}`);

    const matchesJson = JSON.stringify(JSON.parse(matches));

    const result = await contract.submitTransaction('ListToken', tradingPlatformID, price, tokenName, matchesJson);

    const asciiArray = Array.from(new Uint8Array(result));
    const readableResult = String.fromCharCode(...asciiArray);

    console.log(`Transaction has been submitted, result is: ${readableResult}`);
}
