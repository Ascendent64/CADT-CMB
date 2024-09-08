import * as path from 'path';
import * as fs from 'fs';

interface Settings {
    gatewayEndpoint: string;
    mspId: string;
    clientCertPath: string;
    privateKeyPath: string;
    tlsCertPath: string;
    channelName: string;
    chaincodeName: string;
    hostAlias: string;
}

const settings: Settings = JSON.parse(fs.readFileSync(path.resolve(__dirname, '..', 'assets', 'settings.json'), 'utf8'));

export const GATEWAY_ENDPOINT = settings.gatewayEndpoint;
export const MSP_ID = settings.mspId;
export const CLIENT_CERT_PATH = path.resolve(settings.clientCertPath);
export const PRIVATE_KEY_PATH = path.resolve(settings.privateKeyPath);
export const TLS_CERT_PATH = path.resolve(settings.tlsCertPath);
export const CHANNEL_NAME = settings.channelName;
export const CHAINCODE_NAME = settings.chaincodeName;
export const HOST_ALIAS = settings.hostAlias;

console.log("Client Cert Path:", CLIENT_CERT_PATH);
console.log("Private Key Path:", PRIVATE_KEY_PATH);
console.log("TLS Cert Path:", TLS_CERT_PATH);
