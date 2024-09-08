import { Contract } from '@hyperledger/fabric-gateway';

export async function burnToken(contract: Contract, args: string[]): Promise<void> {
    const [tradingPlatformID, orgName, orgURL] = args;
    console.log(`Burning token with TradingPlatformID: ${tradingPlatformID}`);
    const result = await contract.submitTransaction('BurnToken', tradingPlatformID, orgName, orgURL);

    const asciiArray = Array.from(new Uint8Array(result));
    const readableResult = String.fromCharCode(...asciiArray);

    console.log(`Transaction has been submitted, result is: ${readableResult}`);
}
