import { Contract } from '@hyperledger/fabric-gateway';

export async function getOrder(contract: Contract, args: string[]): Promise<any> {
    const orderID = args[0];
    console.log(`Fetching order with OrderID: ${orderID}`);
    const result = await contract.evaluateTransaction('GetOrder', orderID);

    const asciiArray = Array.from(new Uint8Array(result));
    const readableResult = String.fromCharCode(...asciiArray);

    console.log(`Transaction has been evaluated, result is: ${readableResult}`);
    return JSON.parse(readableResult);
}
