import { Contract } from '@hyperledger/fabric-gateway';

export async function removeOrder(contract: Contract, args: string[]): Promise<void> {
    const orderID = args[0];
    console.log(`Removing order with OrderID: ${orderID}`);
    const result = await contract.submitTransaction('RemoveOrder', orderID);

    const asciiArray = Array.from(new Uint8Array(result));
    const readableResult = String.fromCharCode(...asciiArray);

    console.log(`Transaction has been submitted, result is: ${readableResult}`);
}
