import { Contract } from '@hyperledger/fabric-gateway';

export async function claimUnitByID(contract: Contract, args: string[]): Promise<void> {
    const warehouseUnitID = args[0];
    console.log(`Claiming unit with WarehouseUnitID: ${warehouseUnitID}`);
    const result = await contract.submitTransaction('ClaimUnitByID', warehouseUnitID);

    const asciiArray = Array.from(new Uint8Array(result));
    const readableResult = String.fromCharCode(...asciiArray);

    console.log(`Transaction has been submitted, result is: ${readableResult}`);
}
