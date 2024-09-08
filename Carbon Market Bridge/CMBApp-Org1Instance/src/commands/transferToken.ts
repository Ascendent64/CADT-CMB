import { Contract } from '@hyperledger/fabric-gateway';

export async function transferToken(contract: Contract, args: string[]): Promise<void> {
    const [tradingPlatformID, newOwner] = args;
    console.log(`Transferring token with TradingPlatformID: ${tradingPlatformID} to new owner: ${newOwner}`);
    const result = await contract.submitTransaction('TransferToken', tradingPlatformID, newOwner);

    const asciiArray = Array.from(new Uint8Array(result));
    const readableResult = String.fromCharCode(...asciiArray);

    console.log(`Transaction has been submitted, result is: ${readableResult}`);
}
