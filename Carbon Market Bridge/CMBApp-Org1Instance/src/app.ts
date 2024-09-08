import * as crypto from 'crypto';
import express, { Request, Response } from 'express';
import bodyParser from 'body-parser';
import { WebSocketServer } from 'ws';
import * as path from 'path';
import mysql, { RowDataPacket } from 'mysql2/promise';
import { connect, Contract, Identity, Signer, signers, Network, ChaincodeEvent, CloseableAsyncIterable } from '@hyperledger/fabric-gateway';
import { createClient } from 'redis';
import { newGrpcConnection } from './connect';
import { CLIENT_CERT_PATH, PRIVATE_KEY_PATH, TLS_CERT_PATH, MSP_ID, CHANNEL_NAME, CHAINCODE_NAME } from './config';
import { claimUnitByID, burnToken, listToken, delistToken, transferToken, listAllTokens, splitUnits, buyToken, getState, createOrder, getClientIdentity, getOrder, removeOrder } from './commands';
import dbConfig from './dbConfig';
import * as fs from 'fs';
import { exec, spawn } from 'child_process';

const app = express();
const port = 3000;

app.use(bodyParser.json());
app.use(express.static('public'));

const wss = new WebSocketServer({ noServer: true });

const TOKEN_PREFIX = 'token:';
const ORDER_PREFIX = 'order:';

let connection: mysql.Connection;
const redisClient = createClient();
const inMemoryCache: Record<string, any> = {};
const debounceTimers: Record<string, NodeJS.Timeout> = {};
const orderDebounceTimers: Record<string, NodeJS.Timeout> = {};

const checkpointKey = 'fabric_checkpoint';

async function connectRedis() {
    try {
        await redisClient.connect();
        console.log('Connected to Redis');
    } catch (err) {
        console.error('Error connecting to Redis:', (err as Error).message);
    }
}

async function connectDatabase() {
    try {
        connection = await mysql.createConnection(dbConfig);
        console.log('Connected to the database');
    } catch (err) {
        console.error('Error connecting to the database:', (err as Error).message);
    }
}

(async () => {
    await connectDatabase();
    await connectRedis();
})();

app.post('/api/flushRedis', async (req: Request, res: Response) => {
    try {
        exec('redis-cli FLUSHALL', (error, stdout, stderr) => {
            if (error) {
                console.error(`Error flushing Redis: ${stderr}`);
                res.status(500).send('Failed to flush Redis');
                return;
            }
            console.log(`Redis flush output: ${stdout}`);
            res.status(200).send('Redis cache flushed successfully');
        });
    } catch (error) {
        if (error instanceof Error) {
            console.error(`Failed to flush Redis: ${error.message}`);
            res.status(500).send(`Failed to flush Redis: ${error.message}`);
        } else {
            console.error('Unknown error:', error);
            res.status(500).send('Failed to flush Redis due to an unknown error');
        }
    }
});

app.post('/api/:command', async (req: Request, res: Response) => {
    const commandName = req.params.command;
    let args = req.body;

    if (commandName === 'calculateRegressionCoefficients') {
        console.log("Received tokens:", JSON.stringify(args.tokens, null, 2));
        console.log("Received regions:", JSON.stringify(args.regions, null, 2));

        if (!Array.isArray(args.tokens) || typeof args.regions !== 'object') {
            return res.status(400).send('Invalid arguments for calculateRegressionCoefficients');
        }
    } else {
        if (!Array.isArray(args)) {
            return res.status(400).send('Arguments must be an array');
        }

        args = args.map(arg => {
            if (typeof arg === 'string') {
                return arg.trim();
            } else if (Array.isArray(arg)) {
                return arg.map(item => {
                    if (typeof item === 'string') {
                        return item.trim();
                    }
                    return item;
                });
            }
            return arg;
        });
    }

    console.log(`Received command: ${commandName} with args: ${JSON.stringify(args)}`);

    try {
        const client = await newGrpcConnection();
        const gateway = connect({
            client,
            identity: await newIdentity(),
            signer: await newSigner(),
        });
        const network = await gateway.getNetwork(CHANNEL_NAME);
        const contract = network.getContract(CHAINCODE_NAME);

        let result;
        switch (commandName) {
            case 'claimUnitByID':
                console.log(`Executing chaincode function: ClaimUnitByID with args: ${JSON.stringify(args)}`);
                result = await claimUnitByID(contract, args);
                break;
            case 'burnToken':
                console.log(`Executing chaincode function: BurnToken with args: ${JSON.stringify(args)}`);
                result = await burnToken(contract, args);
                break;
            case 'listToken':
                console.log(`Executing chaincode function: ListToken with args: ${JSON.stringify(args)}`);
                result = await listToken(contract, args);
                break;
            case 'delistToken':
                console.log(`Executing chaincode function: DelistToken with args: ${JSON.stringify(args)}`);
                result = await delistToken(contract, args);
                break;
            case 'transferToken':
                console.log(`Executing chaincode function: TransferToken with args: ${JSON.stringify(args)}`);
                result = await transferToken(contract, args);
                break;
            case 'splitUnits':
                console.log(`Executing chaincode function: SplitUnits with args: ${JSON.stringify(args)}`);
                result = await splitUnits(contract, args);
                break;
            case 'listAllTokens':
                console.log(`Retrieving all tokens from Redis`);
                result = await listAllTokensFromRedis();
                await addDatabaseInfo(result);
                res.status(200).json(result);
                broadcastLog(`Executed ${commandName} with args ${JSON.stringify(args)}, result: ${JSON.stringify(result)}`);
                return;
            case 'getClientIdentity':
                console.log(`Executing chaincode function: GetClientIdentity`);
                result = await getClientIdentity(contract);
                res.status(200).send(result);
                broadcastLog(`Executed ${commandName} with result: ${result}`);
                return;
            case 'buyToken':
                console.log(`Executing chaincode function: BuyToken with args: ${JSON.stringify(args)}`);
                result = await buyToken(contract, args);
                break;
            case 'createOrder':
                console.log(`Executing chaincode function: CreateOrder with args: ${JSON.stringify(args)}`);
                result = await createOrder(contract, args);
                break;
            case 'getOrder':
                console.log(`Executing chaincode function: GetOrder with args: ${JSON.stringify(args)}`);
                result = await getOrder(contract, args);
                res.status(200).send(result);
                return;
            case 'removeOrder':
                console.log(`Executing chaincode function: RemoveOrder with args: ${JSON.stringify(args)}`);
                result = await removeOrder(contract, args);
                break;
            case 'listAllOrders':
                console.log(`Retrieving all orders from Redis`);
                result = await listAllOrdersFromRedis();
                res.status(200).json(result);
                broadcastLog(`Executed ${commandName} with args ${JSON.stringify(args)}, result: ${JSON.stringify(result)}`);
                return;
            case 'calculateRegressionCoefficients':
                console.log(`Calculating regression coefficients with args: ${JSON.stringify(args)}`);
                const tokens = args.tokens;
                const regions = args.regions;

                if (!Array.isArray(tokens)) {
                    console.error('Tokens is not an array');
                    return res.status(400).json({ error: 'Tokens must be an array' });
                }
                if (typeof regions !== 'object' || regions === null) {
                    console.error('Regions is not a valid object');
                    return res.status(400).json({ error: 'Regions must be a valid object' });
                }

                const scriptPath = path.resolve(__dirname, '..', 'assets', 'regression.py');
                console.log(`Resolved Python script path: ${scriptPath}`);                

                const pythonProcess = spawn('python3', [scriptPath]);
                let pythonData = '';
                let pythonError = '';

                pythonProcess.stdin.write(JSON.stringify({ tokens, regions }));
                pythonProcess.stdin.end();

                pythonProcess.stdout.on('data', (data) => {
                    pythonData += data.toString();
                });

                pythonProcess.stderr.on('data', (data) => {
                    pythonError += data.toString();
                });

                pythonProcess.on('close', (code) => {
                    if (code !== 0) {
                        console.error(`Python script failed with code ${code}: ${pythonError}`);
                        return res.status(500).json({ error: `Python script failed: ${pythonError}` });
                    }
                    try {
                        const result = JSON.parse(pythonData);
                        return res.json(result);
                    } catch (error) {
                        console.error(`Error parsing Python script output: ${pythonData}`);
                        return res.status(500).json({ error: 'Error parsing Python script output' });
                    }
                });
                return;
            default:
                res.status(400).send(`Unknown command: ${commandName}`);
                return;
        }

        const logMessage = `Executed ${commandName} with args: ${JSON.stringify(args)}, result: ${JSON.stringify(result)}`;
        broadcastLog(logMessage);

        res.status(200).send(result);
    } catch (error) {
        let errorMessage = 'An unexpected error occurred';

        if (error instanceof Error) {
            errorMessage = `Failed to execute ${commandName}: ${error.message}`;
        } else {
            errorMessage = `Failed to execute ${commandName}: ${String(error)}`;
        }

        broadcastLog(errorMessage);
        res.status(500).send(errorMessage);
    }
});

interface Token {
    warehouseUnitID: string;
    unitData?: any;
    burned?: boolean;
    listed?: boolean;
    listedPrice?: string;
    tokenName?: string;
    tradingPlatformID?: string;
    owner?: string;
    unitHash?: string;
    claimedAt?: number;
    burnedByOrgName?: string;
    burnedByOrgURL?: string;
}

interface UnitData {
    warehouseUnitId: string;
    issuanceId: string;
}

interface IssuanceData {
    id: string;
    warehouseProjectId: string;
}

interface ProjectData {
    warehouseProjectId: string;
}

interface Unit {
    unitBlockStart: string | null;
    unitBlockEnd: string | null;
    unitCount: number | null;
    warehouseUnitID: string;
    issuanceID: string | null;
    projectLocationID: string | null;
    orgUID: string;
    unitOwner: string | null;
    countryJurisdictionOfOwner: string | null;
    inCountryJurisdictionOfOwner: string | null;
    serialNumberBlock: string | null;
    serialNumberPattern: string | null;
    vintageYear: number | null;
    unitType: string | null;
    marketplace: string | null;
    marketplaceLink: string | null;
    marketplaceIdentifier: string | null;
    unitTags: string | null;
    unitStatus: string | null;
    unitStatusReason: string | null;
    unitRegistryLink: string | null;
    correspondingAdjustmentDeclaration: string | null;
    correspondingAdjustmentStatus: string | null;
    timeStaged: number | null;
    createdAt: string | null;
    updatedAt: string | null;
}

async function addDatabaseInfo(tokens: Token[]): Promise<void> {
    if (tokens.length === 0) return;

    const unitIds = tokens.map(token => `'${token.warehouseUnitID}'`).join(',');
    const unitQuery = `SELECT * FROM units WHERE warehouseUnitId IN (${unitIds})`;

    try {
        const [unitResults] = await connection.query<RowDataPacket[]>(unitQuery);

        const unitsMap: Record<string, UnitData> = unitResults.reduce((acc: Record<string, UnitData>, unit: any) => {
            acc[unit.warehouseUnitId] = unit;
            return acc;
        }, {});

        const issuanceIds = [...new Set(unitResults.map(unit => unit.issuanceId))].filter(Boolean).map(id => `'${id}'`).join(',');

        let issuanceResults: RowDataPacket[] = [];
        if (issuanceIds.length > 0) {
            const issuanceQuery = `SELECT * FROM issuances WHERE id IN (${issuanceIds})`;
            [issuanceResults] = await connection.query<RowDataPacket[]>(issuanceQuery);
        }

        const issuancesMap: Record<string, IssuanceData> = issuanceResults.reduce((acc: Record<string, IssuanceData>, issuance: any) => {
            acc[issuance.id] = issuance;
            return acc;
        }, {});

        const projectIds = [...new Set(issuanceResults.map(issuance => issuance.warehouseProjectId))].filter(Boolean).map(id => `'${id}'`).join(',');

        let projectResults: RowDataPacket[] = [];
        if (projectIds.length > 0) {
            const projectQuery = `SELECT * FROM projects WHERE warehouseProjectId IN (${projectIds})`;
            [projectResults] = await connection.query<RowDataPacket[]>(projectQuery);
        }

        const projectsMap: Record<string, ProjectData> = projectResults.reduce((acc: Record<string, ProjectData>, project: any) => {
            acc[project.warehouseProjectId] = project;
            return acc;
        }, {});

        let locationResults: RowDataPacket[] = [];
        if (projectIds.length > 0) {
            const locationQuery = `SELECT * FROM projectlocations WHERE warehouseProjectId IN (${projectIds})`;
            [locationResults] = await connection.query<RowDataPacket[]>(locationQuery);
        }

        const locationsMap: Record<string, any[]> = locationResults.reduce((acc: Record<string, any[]>, location: any) => {
            if (!acc[location.warehouseProjectId]) {
                acc[location.warehouseProjectId] = [];
            }
            acc[location.warehouseProjectId].push(location);
            return acc;
        }, {});

        let ratingResults: RowDataPacket[] = [];
        if (projectIds.length > 0) {
            const ratingQuery = `SELECT * FROM projectratings WHERE warehouseProjectId IN (${projectIds})`;
            [ratingResults] = await connection.query<RowDataPacket[]>(ratingQuery);
        }

        const ratingsMap: Record<string, any[]> = ratingResults.reduce((acc: Record<string, any[]>, rating: any) => {
            if (!acc[rating.warehouseProjectId]) {
                acc[rating.warehouseProjectId] = [];
            }
            acc[rating.warehouseProjectId].push(rating);
            return acc;
        }, {});

        let cobenefitResults: RowDataPacket[] = [];
        if (projectIds.length > 0) {
            const cobenefitQuery = `SELECT * FROM cobenefits WHERE warehouseProjectId IN (${projectIds})`;
            [cobenefitResults] = await connection.query<RowDataPacket[]>(cobenefitQuery);
        }

        const cobenefitsMap: Record<string, any[]> = cobenefitResults.reduce((acc: Record<string, any[]>, cobenefit: any) => {
            if (!acc[cobenefit.warehouseProjectId]) {
                acc[cobenefit.warehouseProjectId] = [];
            }
            acc[cobenefit.warehouseProjectId].push(cobenefit);
            return acc;
        }, {});

        let estimationResults: RowDataPacket[] = [];
        if (projectIds.length > 0) {
            const estimationQuery = `SELECT * FROM estimations WHERE warehouseProjectId IN (${projectIds})`;
            [estimationResults] = await connection.query<RowDataPacket[]>(estimationQuery);
        }

        const estimationsMap: Record<string, any[]> = estimationResults.reduce((acc: Record<string, any[]>, estimation: any) => {
            if (!acc[estimation.warehouseProjectId]) {
                acc[estimation.warehouseProjectId] = [];
            }
            acc[estimation.warehouseProjectId].push(estimation);
            return acc;
        }, {});

        let relatedProjectResults: RowDataPacket[] = [];
        if (projectIds.length > 0) {
            const relatedProjectQuery = `SELECT * FROM relatedprojects WHERE warehouseProjectId IN (${projectIds})`;
            [relatedProjectResults] = await connection.query<RowDataPacket[]>(relatedProjectQuery);
        }

        const relatedProjectsMap: Record<string, any[]> = relatedProjectResults.reduce((acc: Record<string, any[]>, relatedProject: any) => {
            if (!acc[relatedProject.warehouseProjectId]) {
                acc[relatedProject.warehouseProjectId] = [];
            }
            acc[relatedProject.warehouseProjectId].push(relatedProject);
            return acc;
        }, {});

        tokens.forEach(token => {
            if (unitsMap[token.warehouseUnitID]) {
                token.unitData = unitsMap[token.warehouseUnitID];
                const issuance = issuancesMap[token.unitData.issuanceId];
                if (issuance) {
                    token.unitData.issuanceData = issuance;
                    const project = projectsMap[issuance.warehouseProjectId];
                    if (project) {
                        token.unitData.issuanceData.projectData = project;
                        token.unitData.issuanceData.projectData.locations = locationsMap[project.warehouseProjectId] || [];
                        token.unitData.issuanceData.projectData.ratings = ratingsMap[project.warehouseProjectId] || [];
                        token.unitData.issuanceData.projectData.cobenefits = cobenefitsMap[project.warehouseProjectId] || [];
                        token.unitData.issuanceData.projectData.estimations = estimationsMap[project.warehouseProjectId] || [];
                        token.unitData.issuanceData.projectData.relatedProjects = relatedProjectsMap[project.warehouseProjectId] || [];
                    }
                }
            }
        });

    } catch (error) {
        console.error('Failed to fetch data from database:', (error as Error).message);
        throw error;
    }
}

async function updateCacheForEvent(event: any) {
    const action = event.action;
    const tradingPlatformID = event.tradingPlatformID;
    const orderID = event.orderID;

    if (action === 'ClaimUnitByID') {
        console.log(`Debounce: ClaimUnitByID event for ${tradingPlatformID}, triggering Redis update immediately.`);
        await updateRedisForClaimUnit(tradingPlatformID);
        return;
    }

    if (action === 'OrderFill') {
        const tradingPlatformIDs: string[] = event.tradingPlatformIDs || [];
        const orderIDs: string[] = event.orderIDs || [];

        if (tradingPlatformIDs.length === 0 && orderIDs.length === 0) {
            console.error("OrderFill event is missing 'tradingPlatformIDs' or 'orderIDs'");
            return;
        }

        for (const id of tradingPlatformIDs) {
            if (debounceTimers[id]) {
                clearTimeout(debounceTimers[id]);
                console.log(`Debounce: Cleared existing timer for tradingPlatformID ${id} due to new event.`);
            }

            debounceTimers[id] = setTimeout(async () => {
                try {
                    await updateRedisForToken(id);
                    console.log(`Debounce: Updating Redis for tradingPlatformID ${id} after delay.`);
                    delete debounceTimers[id];
                } catch (error) {
                    console.error(`Failed to update cache for event action ${action}:`, (error as Error).message);
                }
            }, 2000);
        }

        for (const oid of orderIDs) {
            if (orderDebounceTimers[oid]) {
                clearTimeout(orderDebounceTimers[oid]);
                console.log(`Debounce: Cleared existing timer for orderID ${oid} due to new event.`);
            }

            orderDebounceTimers[oid] = setTimeout(async () => {
                try {
                    await updateRedisForOrder(oid);
                    console.log(`Debounce: Updating Redis for orderID ${oid} after delay.`);
                    delete orderDebounceTimers[oid];
                } catch (error) {
                    console.error(`Failed to update cache for OrderFill event:`, (error as Error).message);
                }
            }, 2000);
        }

        return;
    }


    if (action === 'RemoveOrder') {
        console.log(`Debounce: RemoveOrder event for orderID ${orderID}, triggering Redis update immediately.`);
        await updateRedisForOrder(orderID);
        return;
    }

    console.log(`Updating cache for event: ${JSON.stringify(event)}`); 

    if (action === 'SplitUnits') {
        const newTradingPlatformIDs = event.newTradingPlatformIDs;
        if (!Array.isArray(newTradingPlatformIDs) || newTradingPlatformIDs.length === 0) {
            console.error('Invalid newTradingPlatformIDs in SplitUnits event');
            return;
        }

        for (const id of newTradingPlatformIDs) {
            await updateRedisForToken(id);
            console.log(`Debounce: Immediate Redis update for new tradingPlatformID ${id} from SplitUnits.`);
        }

        await updateRedisForToken(event.originalPlatformID);
        console.log(`Debounce: Immediate Redis update for original tradingPlatformID ${event.originalPlatformID} from SplitUnits.`);
        return;
    }

    inMemoryCache[tradingPlatformID] = event;

    if (debounceTimers[tradingPlatformID]) {
        clearTimeout(debounceTimers[tradingPlatformID]);
        console.log(`Debounce: Cleared existing timer for tradingPlatformID ${tradingPlatformID} due to new event.`);
    }

    debounceTimers[tradingPlatformID] = setTimeout(async () => {
        try {
            const cachedEvent = inMemoryCache[tradingPlatformID];

            if (cachedEvent.action === 'ClaimUnitByID') {
                await updateRedisForClaimUnit(tradingPlatformID);
            } else {
                await updateRedisForToken(tradingPlatformID); 
            }

            console.log(`Debounce: Updating Redis for tradingPlatformID ${tradingPlatformID} after delay.`);
            delete inMemoryCache[tradingPlatformID];
            delete debounceTimers[tradingPlatformID];
        } catch (error) {
            console.error(`Failed to update cache for event action ${action}:`, (error as Error).message);
        }
    }, 2000);
}

async function updateRedisForClaimUnit(tradingPlatformID: string) {
    try {
        const tokenState = await getStateForRedis(tradingPlatformID);
        const tokenStateJSON = JSON.stringify(tokenState);

        await redisClient.set(`${TOKEN_PREFIX}${tradingPlatformID}`, tokenStateJSON);

        const logMessage = `Updated Redis for ClaimUnit: ${tradingPlatformID} with data: ${tokenStateJSON}`;
        console.log(logMessage);
        broadcastLog(logMessage);
    } catch (error) {
        console.error(`Failed to update Redis for ClaimUnit event:`, (error as Error).message);
    }
}

async function updateRedisForToken(tradingPlatformID: string) {
    try {
        if (!tradingPlatformID) {
            throw new Error('Invalid tradingPlatformID');
        }
        const key = `${TOKEN_PREFIX}${tradingPlatformID}`;
        console.log(`Fetching state for tradingPlatformID: ${tradingPlatformID}`);
        const tokenState = await getStateForRedis(tradingPlatformID);
        console.log(`Fetched token state for ${tradingPlatformID}: ${JSON.stringify(tokenState)}`);
        await redisClient.set(key, JSON.stringify(tokenState));
        console.log(`Updated Redis for token: ${key} with data: ${JSON.stringify(tokenState)}`);
        broadcastLog(`Updated Redis for token: ${key}`);
    } catch (error) {
        console.error(`Failed to update Redis for token: ${tradingPlatformID}`, (error as Error).message);
    }
}


async function getStateForRedis(tradingPlatformID: string): Promise<any> {
    if (!tradingPlatformID) {
        throw new Error('Invalid tradingPlatformID');
    }
    const client = await newGrpcConnection();
    const gateway = connect({
        client,
        identity: await newIdentity(),
        signer: await newSigner(),
    });
    const network = await gateway.getNetwork(CHANNEL_NAME);
    const contract = network.getContract(CHAINCODE_NAME);
    const tokenState = await getState(contract, tradingPlatformID);

    const completeTokenState = ensureTokenFields(tokenState);

    console.log('Token state before returning:', completeTokenState);

    return completeTokenState;
}

function broadcastLog(message: string): void {
    wss.clients.forEach(client => {
        if (client.readyState === client.OPEN) {
            client.send(message);
        }
    });
}

async function newIdentity(): Promise<Identity> {
    console.log('Creating new identity...');
    console.log(`Reading certificate from: ${CLIENT_CERT_PATH}`);
    const certPath = path.resolve(CLIENT_CERT_PATH);
    const credentials = await fs.promises.readFile(certPath);
    return { mspId: MSP_ID, credentials };
}

async function newSigner(): Promise<Signer> {
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

const server = app.listen(port, () => {
    console.log(`Server is running on port ${port}`);
});

server.on('upgrade', (request, socket, head) => {
    wss.handleUpgrade(request, socket, head, (ws) => {
        wss.emit('connection', ws, request);
    });
});

async function setupEventListener(network: Network) {
    try {
        const lastBlockNumber = await getLastBlockNumberFromRedis();

        const events: CloseableAsyncIterable<ChaincodeEvent> = await network.getChaincodeEvents(CHAINCODE_NAME, {
            startBlock: lastBlockNumber
        });

        for await (const event of events) {
            const eventPayload = Buffer.from(event.payload).toString('utf-8');
            const logMessage = `Event received: ${eventPayload}`;

            console.log(logMessage); 
            broadcastLog(logMessage);  

            const eventData = JSON.parse(eventPayload);
            await updateCacheForEvent(eventData);

            await checkpointChaincodeEvent(event.blockNumber);

            broadcastLog(`Event received: ${eventPayload}`);
        }

        events.close();
    } catch (error) {
        console.error('Error setting up event listener:', (error as Error).message);
    }
}

async function getLastBlockNumberFromRedis(): Promise<bigint> {
    try {
        const data = await redisClient.get(checkpointKey);
        console.log(`Retrieved checkpoint from Redis: ${data}`);
        return data ? BigInt(data) : BigInt(0);
    } catch (error) {
        console.error('Error reading checkpoint from Redis:', (error as Error).message);
        return BigInt(0); 
    }
}

async function checkpointChaincodeEvent(blockNumber: bigint): Promise<void> {
    try {
        console.log(`Updating checkpoint to block number: ${blockNumber}`);
        await redisClient.set(checkpointKey, blockNumber.toString());
        console.log(`Checkpoint updated to block number: ${blockNumber}`);
    } catch (error) {
        console.error('Error writing checkpoint to Redis:', (error as Error).message);
        throw error;
    }
}

(async () => {
    try {
        const client = await newGrpcConnection();
        const gateway = connect({
            client,
            identity: await newIdentity(),
            signer: await newSigner(),
        });
        const network = await gateway.getNetwork(CHANNEL_NAME);
        const contract = network.getContract(CHAINCODE_NAME);

        await setupEventListener(network);
    } catch (error) {
        console.error('Error setting up event listener:', (error as Error).message);
    }
})();

async function fetchUnitByID(warehouseUnitID: string): Promise<Unit> {
    const [rows] = await connection.query<RowDataPacket[]>(
        'SELECT * FROM units WHERE warehouseUnitID = ?',
        [warehouseUnitID]
    );
    const unitRow = rows[0];
    return {
        unitBlockStart: unitRow.unitBlockStart,
        unitBlockEnd: unitRow.unitBlockEnd,
        unitCount: unitRow.unitCount,
        warehouseUnitID: unitRow.warehouseUnitID,
        issuanceID: unitRow.issuanceID,
        projectLocationID: unitRow.projectLocationID,
        orgUID: unitRow.orgUID,
        unitOwner: unitRow.unitOwner,
        countryJurisdictionOfOwner: unitRow.countryJurisdictionOfOwner,
        inCountryJurisdictionOfOwner: unitRow.inCountryJurisdictionOfOwner,
        serialNumberBlock: unitRow.serialNumberBlock,
        serialNumberPattern: unitRow.serialNumberPattern,
        vintageYear: unitRow.vintageYear,
        unitType: unitRow.unitType,
        marketplace: unitRow.marketplace,
        marketplaceLink: unitRow.marketplaceLink,
        marketplaceIdentifier: unitRow.marketplaceIdentifier,
        unitTags: unitRow.unitTags,
        unitStatus: unitRow.unitStatus,
        unitStatusReason: unitRow.unitStatusReason,
        unitRegistryLink: unitRow.unitRegistryLink,
        correspondingAdjustmentDeclaration: unitRow.correspondingAdjustmentDeclaration,
        correspondingAdjustmentStatus: unitRow.correspondingAdjustmentStatus,
        timeStaged: unitRow.timeStaged,
        createdAt: unitRow.createdAt,
        updatedAt: unitRow.updatedAt,
    };
}

async function listAllTokensFromRedis(): Promise<Token[]> {
    const keys = await redisClient.keys(`${TOKEN_PREFIX}*`);
    console.log(`Found keys: ${keys}`);
    const tokens: Token[] = [];

    for (const key of keys) {
        const tokenData = await redisClient.get(key);
        if (tokenData) {
            try {
                const token = JSON.parse(tokenData);
                console.log(`Parsed token data for key ${key}: ${JSON.stringify(token)}`);
                tokens.push(token);
            } catch (error) {
                console.error(`Failed to parse token data for key ${key}:`, error);
            }
        }
    }
    return tokens;
}

function ensureTokenFields(token: any): any {
    return {
        ...token,
        burned: token.burned !== undefined ? token.burned : false,
        burnedAt: token.burnedAt !== undefined ? token.burnedAt : 0,
        burnedByOrgName: token.burnedByOrgName !== undefined ? token.burnedByOrgName : 'NA',
        burnedByOrgURL: token.burnedByOrgURL !== undefined ? token.burnedByOrgURL : 'NA',
        listedPrice: token.listedPrice !== undefined ? token.listedPrice : 'NA',
        tokenName: token.tokenName !== undefined ? token.tokenName : 'NA',
    };
}

async function updateRedisForOrder(orderID: string) {
    try {
        const key = `${ORDER_PREFIX}${orderID}`;
        const client = await newGrpcConnection();
        const gateway = connect({
            client,
            identity: await newIdentity(),
            signer: await newSigner(),
        });
        const network = await gateway.getNetwork(CHANNEL_NAME);
        const contract = network.getContract(CHAINCODE_NAME);
        const orderData = await getOrder(contract, [orderID]);

        if (!orderData.coefficients) {
            orderData.coefficients = {};
        }

        console.log('Order data fetched:', orderData);
        const orderDataJSON = JSON.stringify(orderData);
        await redisClient.set(key, orderDataJSON);

        const logMessage = `Updated Redis for CreateOrder: ${key} with data: ${orderDataJSON}`;
        console.log(logMessage);
        broadcastLog(logMessage);
    } catch (error) {
        console.error(`Failed to update Redis for CreateOrder event:`, (error as Error).message);
    }
}

async function listAllOrdersFromRedis(): Promise<any[]> {
    const keys = await redisClient.keys('*');
    const orders: any[] = [];

    for (const key of keys) {
        if (key === checkpointKey) continue; 

        const orderData = await redisClient.get(key);
        if (orderData) {
            const order = JSON.parse(orderData);
            if (typeof order === 'object' && order !== null && order.hasOwnProperty('orderID')) {
                orders.push(order);
            }
        }
    }

    return orders;
}
