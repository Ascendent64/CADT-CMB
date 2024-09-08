import * as grpc from '@grpc/grpc-js';
import { CLIENT_CERT_PATH, GATEWAY_ENDPOINT, HOST_ALIAS, MSP_ID, PRIVATE_KEY_PATH, TLS_CERT_PATH } from './config';
import * as fs from 'fs';

export async function newGrpcConnection(): Promise<grpc.Client> {
    if (TLS_CERT_PATH) {
        const tlsRootCert = await fs.promises.readFile(TLS_CERT_PATH);
        const tlsCredentials = grpc.credentials.createSsl(tlsRootCert);
        return new grpc.Client(GATEWAY_ENDPOINT, tlsCredentials, newGrpcClientOptions());
    }

    return new grpc.Client(GATEWAY_ENDPOINT, grpc.credentials.createInsecure());
}

function newGrpcClientOptions(): grpc.ClientOptions {
    const result: grpc.ClientOptions = {};
    if (HOST_ALIAS) {
        result['grpc.ssl_target_name_override'] = HOST_ALIAS;
    }
    return result;
}
