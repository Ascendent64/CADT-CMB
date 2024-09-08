import { Contract } from '@hyperledger/fabric-gateway';

export async function createOrder(contract: Contract, args: string[]): Promise<void> {
    const [paymentToken, maxPrice, quantity, vintageYearRange, orderType, coefficientsString, matchesString] = args;
    let coefficients = JSON.parse(coefficientsString);
    let matches = JSON.parse(matchesString);

    if (!coefficients) {
        coefficients = {};
    }
    if (!matches) {
        matches = [];
    }

    console.log(`Creating order with PaymentToken: ${paymentToken}, MaxPrice: ${maxPrice}, Quantity: ${quantity}, VintageYearRange: ${vintageYearRange}, OrderType: ${orderType}, Coefficients: ${JSON.stringify(coefficients)}, Matches: ${JSON.stringify(matches)}`);
    
    const result = await contract.submitTransaction('CreateOrder', paymentToken, maxPrice, quantity, vintageYearRange, orderType, JSON.stringify(coefficients), JSON.stringify(matches));
    
    const asciiArray = Array.from(new Uint8Array(result));
    const readableResult = String.fromCharCode(...asciiArray);
    
    console.log(`Transaction has been submitted, result is: ${readableResult}`);
}
