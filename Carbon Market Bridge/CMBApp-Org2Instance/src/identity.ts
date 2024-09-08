import * as crypto from 'crypto';
import { Identity, signers, Signer } from '@hyperledger/fabric-gateway';
import { CLIENT_CERT_PATH, MSP_ID, PRIVATE_KEY_PATH } from './config';
import * as fs from 'fs';
import * as path from 'path';

export async function newIdentity(): Promise<Identity> {
    console.log('Creating new identity...');
    console.log(`Reading certificate from: ${CLIENT_CERT_PATH}`);
    const certPath = path.resolve(CLIENT_CERT_PATH);
    const credentials = await fs.promises.readFile(certPath);
    return { mspId: MSP_ID, credentials };
}

export async function newSigner(): Promise<Signer> {
    console.log('Creating new signer...');
    console.log(`Reading private key from: ${PRIVATE_KEY_PATH}`);
    const keyPath = path.resolve(PRIVATE_KEY_PATH);
    const privateKeyPem = await fs.promises.readFile(keyPath);

    const privateKey = crypto.createPrivateKey({
        key: privateKeyPem,
        format: 'pem'
    });

    return signers.newPrivateKeySigner(privateKey);
}
