import { Contract } from '@hyperledger/fabric-gateway';

export async function splitUnits(contract: Contract, args: string[]): Promise<void> {
    const [tradingPlatformID, method, value] = args;
    console.log(`Splitting units for TradingPlatformID: ${tradingPlatformID} using method: ${method} with value: ${value}`);

    try {
        const result = await contract.submitTransaction('SplitUnits', tradingPlatformID, method, value.toString());
        const asciiArray = Array.from(new Uint8Array(result));
        const readableResult = String.fromCharCode(...asciiArray);
        console.log(`Transaction has been submitted, result is: ${readableResult}`);
    } catch (error) {
        console.error(`Failed to submit transaction SplitUnits:`, error);
        throw error;
    }
}
