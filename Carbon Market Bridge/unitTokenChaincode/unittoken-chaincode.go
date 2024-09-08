package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type Token struct {
	TradingPlatformID string   `json:"tradingPlatformID"`
	WarehouseUnitID   string   `json:"warehouseUnitID"`
	UnitHash          string   `json:"unitHash"`
	Burned            bool     `json:"burned"`
	BurnedByOrgName   string   `json:"burnedByOrgName,omitempty"`
	BurnedByOrgURL    string   `json:"burnedByOrgURL,omitempty"`
	BurnedAt          int64    `json:"burnedAt"`
	Owner             string   `json:"owner"`
	Listed            bool     `json:"listed"`
	UnitCount         int64    `json:"unitCount"`
	ClaimedAt         int64    `json:"claimedAt"`
	ListedUnitPrice   string   `json:"listedUnitPrice,omitempty"`
	TokenName         string   `json:"tokenName,omitempty"`
	SplitInto         []string `json:"splitInto"`
	SplitFrom         string   `json:"splitFrom"`
}

type Order struct {
	OrderID          string                        `json:"orderID"`
	PaymentToken     string                        `json:"paymentToken"`
	MaxPrice         float64                       `json:"maxPrice"`
	Quantity         int                           `json:"quantity"`
	VintageYearRange string                        `json:"vintageYearRange"`
	OrderType        string                        `json:"orderType"`
	Coefficients     map[string]map[string]float64 `json:"coefficients"`
	Owner            string                        `json:"owner"`
	EscrowedAmount   int                           `json:"escrowedAmount"`
	Filled           bool                          `json:"filled"`
	FilledWith       []string                      `json:"filledWith"`
}

type ProjectData struct {
	CurrentRegistry string
	Sector          string
	Country         string
	VintageYear     int64
}

type Unit struct {
	UnitBlockStart                     string `json:"UnitBlockStart"`
	UnitBlockEnd                       string `json:"UnitBlockEnd"`
	UnitCount                          int64  `json:"UnitCount"`
	WarehouseUnitID                    string `json:"WarehouseUnitID"`
	IssuanceID                         string `json:"IssuanceID"`
	ProjectLocationID                  string `json:"ProjectLocationID"`
	OrgUID                             string `json:"OrgUID"`
	UnitOwner                          string `json:"UnitOwner"`
	CountryJurisdictionOfOwner         string `json:"CountryJurisdictionOfOwner"`
	InCountryJurisdictionOfOwner       string `json:"InCountryJurisdictionOfOwner"`
	SerialNumberBlock                  string `json:"SerialNumberBlock"`
	SerialNumberPattern                string `json:"SerialNumberPattern"`
	VintageYear                        int64  `json:"VintageYear"`
	UnitType                           string `json:"UnitType"`
	Marketplace                        string `json:"Marketplace"`
	MarketplaceLink                    string `json:"MarketplaceLink"`
	MarketplaceIdentifier              string `json:"MarketplaceIdentifier"`
	UnitTags                           string `json:"UnitTags"`
	UnitStatus                         string `json:"UnitStatus"`
	UnitStatusReason                   string `json:"UnitStatusReason"`
	UnitRegistryLink                   string `json:"UnitRegistryLink"`
	CorrespondingAdjustmentDeclaration string `json:"CorrespondingAdjustmentDeclaration"`
	CorrespondingAdjustmentStatus      string `json:"CorrespondingAdjustmentStatus"`
	TimeStaged                         int64  `json:"TimeStaged"`
	CreatedAt                          string `json:"CreatedAt"`
	UpdatedAt                          string `json:"UpdatedAt"`
}

type PaymentToken struct {
	TokenID string `json:"tokenID"`
	Creator string `json:"creator"`
	Balance int    `json:"balance"`
}

type DatabaseDetails struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbName"`
}

type MarketplaceChaincode struct {
	contractapi.Contract
}

var db *sql.DB

func initializeDB() error {
	cfg := mysql.Config{
		User:                 "User",
		Passwd:               "password123",
		Net:                  "tcp",
		Addr:                 "172.23.209.124:3306",
		DBName:               "CADTDatabase",
		AllowNativePasswords: true,
	}

	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return fmt.Errorf("Failed to open the database: %v", err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		return fmt.Errorf("Failed to ping the database: %v", pingErr)
	}

	fmt.Println("Connected to the database successfully!")
	return nil
}

func (m *MarketplaceChaincode) ClaimUnitByID(ctx contractapi.TransactionContextInterface, warehouseUnitID string) (string, error) {
	log.Println("Starting ClaimUnitByID function")
	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		log.Printf("Failed to get client identity: %v", err)
		return "", fmt.Errorf("failed to get client identity: %v", err)
	}
	log.Printf("Client ID: %s", clientID)

	// Check if the unit is already claimed
	claimedStatus, err := ctx.GetStub().GetState("CLAIMED_" + warehouseUnitID)
	if err != nil {
		log.Printf("Failed to read from world state: %v", err)
		return "", fmt.Errorf("failed to read from world state: %v", err)
	}
	if claimedStatus != nil {
		log.Printf("Unit %s already claimed", warehouseUnitID)
		return "unit already claimed", nil
	}

	// Fetch the unit from the database
	log.Println("Fetching unit from database")
	unit, err := m.fetchUnitByID(ctx, warehouseUnitID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("Unit does not exist in database")
			return "unit does not exist", nil
		}
		log.Printf("Failed to fetch unit: %v", err)
		return "", fmt.Errorf("failed to fetch unit: %v", err)
	}

	log.Printf("Fetched unit: %+v", unit)

	log.Printf("Checking if client ID matches MarketplaceIdentifier")
	log.Printf("Client ID: %s, MarketplaceIdentifier: %s", clientID, unit.MarketplaceIdentifier)
	if unit.MarketplaceIdentifier != clientID {
		log.Printf("Caller is not authorized to claim the unit: %s", warehouseUnitID)
		return "caller is not authorized to claim the unit", nil
	}

	log.Println("Generating TradingPlatformID using transaction timestamp")
	i := 0
	newID, err := generateUniqueID(ctx, i)
	tradingPlatformID := newID
	log.Printf("Generated TradingPlatformID: %s", tradingPlatformID)

	log.Printf("Checking if TradingPlatformID %s already exists in the ledger", tradingPlatformID)
	existingTokenBytes, err := ctx.GetStub().GetState(tradingPlatformID)
	if err != nil {
		log.Printf("Failed to check existing token in ledger: %v", err)
		return "", fmt.Errorf("failed to check existing token in ledger: %v", err)
	}
	if existingTokenBytes != nil {
		log.Printf("TradingPlatformID %s already exists in ledger", tradingPlatformID)
		log.Printf("Existing token data for TradingPlatformID: %s", string(existingTokenBytes))
		return "TradingPlatformID already exists in ledger", nil
	}

	timestamp, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		log.Printf("Failed to get transaction timestamp: %v", err)
		return "", fmt.Errorf("failed to get transaction timestamp: %v", err)
	}
	claimedAt := timestamp.Seconds

	hashUnitData := func(data string) string {
		hash := sha256.New()
		hash.Write([]byte(data))
		return hex.EncodeToString(hash.Sum(nil))
	}

	unitData := concatenateUnitData(unit)
	unitHash := hashUnitData(unitData)
	token := &Token{
		TradingPlatformID: tradingPlatformID,
		WarehouseUnitID:   unit.WarehouseUnitID,
		UnitHash:          unitHash,
		Burned:            false,
		BurnedByOrgName:   "NA",
		BurnedByOrgURL:    "NA",
		BurnedAt:          0,
		Owner:             clientID,
		Listed:            false,
		UnitCount:         unit.UnitCount,
		ClaimedAt:         claimedAt,
		ListedUnitPrice:   "NA",
		TokenName:         "NA",
		SplitInto:         []string{},
		SplitFrom:         "",
	}

	tokenBytes, err := json.Marshal(token)
	if err != nil {
		log.Printf("Failed to marshal token: %v", err)
		return "", fmt.Errorf("failed to marshal token: %v", err)
	}
	log.Printf("Token marshaled: %s", string(tokenBytes))

	log.Println("Saving token to ledger")
	err = ctx.GetStub().PutState(tradingPlatformID, tokenBytes)
	if err != nil {
		log.Printf("Failed to put token to world state: %v", err)
		return "", fmt.Errorf("failed to put token to world state: %v", err)
	}

	log.Println("Marking the unit as claimed")
	err = ctx.GetStub().PutState("CLAIMED_"+warehouseUnitID, []byte("claimed"))
	if err != nil {
		log.Printf("Failed to mark unit as claimed in world state: %v", err)
		return "", fmt.Errorf("failed to mark unit as claimed in world state: %v", err)
	}

	eventPayload := map[string]interface{}{
		"action":            "ClaimUnitByID",
		"tradingPlatformID": tradingPlatformID,
		"warehouseUnitID":   warehouseUnitID,
	}
	eventPayloadBytes, _ := json.Marshal(eventPayload)
	err = ctx.GetStub().SetEvent("TokenEvent", eventPayloadBytes)
	if err != nil {
		log.Printf("Failed to emit event: %v", err)
		return "", fmt.Errorf("failed to emit event: %v", err)
	}

	log.Println("ClaimUnitByID function completed successfully")
	return "unit successfully claimed", nil
}

func (m *MarketplaceChaincode) BurnToken(ctx contractapi.TransactionContextInterface, tradingPlatformID, orgName, orgURL string) (string, error) {
	log.Printf("Starting BurnToken function for TradingPlatformID: %s", tradingPlatformID)

	tokenBytes, err := ctx.GetStub().GetState(tradingPlatformID)
	if err != nil {
		log.Printf("Failed to read from world state: %v", err)
		return "", fmt.Errorf("failed to read from world state: %v", err)
	}
	if tokenBytes == nil {
		log.Printf("No state found for TradingPlatformID: %s", tradingPlatformID)
		return "no state found for TradingPlatformID", nil
	}

	var token Token
	err = json.Unmarshal(tokenBytes, &token)
	if err != nil {
		log.Printf("Failed to unmarshal token: %v", err)
		return "", fmt.Errorf("failed to unmarshal token: %v", err)
	}

	log.Printf("Current owner of token %s: %s", tradingPlatformID, token.Owner)

	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		log.Printf("Failed to get client identity: %v", err)
		return "", fmt.Errorf("failed to get client identity: %v", err)
	}
	log.Printf("Caller client ID: %s", clientID)
	if token.Owner != clientID {
		log.Printf("Caller is not the owner of token %s", tradingPlatformID)
		return "caller is not the owner of the token", nil
	}

	if token.Listed {
		log.Printf("Token %s is listed and cannot be burned", tradingPlatformID)
		return "token is listed and cannot be burned", nil
	}

	timestamp, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		log.Printf("Failed to get transaction timestamp: %v", err)
		return "", fmt.Errorf("failed to get transaction timestamp: %v", err)
	}
	burnedAt := timestamp.Seconds

	token.Burned = true
	token.BurnedByOrgName = orgName
	token.BurnedByOrgURL = orgURL
	token.BurnedAt = burnedAt

	tokenBytes, err = json.Marshal(token)
	if err != nil {
		log.Printf("Failed to marshal token: %v", err)
		return "", fmt.Errorf("failed to marshal token: %v", err)
	}

	err = ctx.GetStub().PutState(tradingPlatformID, tokenBytes)
	if err != nil {
		log.Printf("Failed to put token to world state: %v", err)
		return "", fmt.Errorf("failed to put token to world state: %v", err)
	}

	eventPayload := map[string]interface{}{
		"action":            "BurnToken",
		"tradingPlatformID": tradingPlatformID,
	}
	eventPayloadBytes, _ := json.Marshal(eventPayload)
	err = ctx.GetStub().SetEvent("TokenEvent", eventPayloadBytes)
	if err != nil {
		log.Printf("Failed to emit event: %v", err)
		return "", fmt.Errorf("failed to emit event: %v", err)
	}

	log.Printf("BurnToken function completed successfully for TradingPlatformID: %s", tradingPlatformID)
	return "token successfully burned", nil
}

func (m *MarketplaceChaincode) GetClientIdentity(ctx contractapi.TransactionContextInterface) (string, error) {
	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return "", fmt.Errorf("failed to get client identity: %v", err)
	}
	log.Printf("Client ID: %s", clientID)
	return clientID, nil
}

func (m *MarketplaceChaincode) TransferToken(ctx contractapi.TransactionContextInterface, tradingPlatformID, newOwner string) (string, error) {
	log.Println("Starting TransferToken function")

	tokenBytes, err := ctx.GetStub().GetState(tradingPlatformID)
	if err != nil {
		log.Printf("Failed to read from world state: %v", err)
		return "", fmt.Errorf("failed to read from world state: %v", err)
	}
	if tokenBytes == nil {
		log.Printf("No token found for TradingPlatformID: %s", tradingPlatformID)
		return "token does not exist", nil
	}

	var token Token
	err = json.Unmarshal(tokenBytes, &token)
	if err != nil {
		log.Printf("Failed to unmarshal token: %v", err)
		return "", fmt.Errorf("failed to unmarshal token: %v", err)
	}

	log.Printf("Current owner of token %s: %s", tradingPlatformID, token.Owner)

	if token.Burned {
		log.Printf("Token %s is burned and cannot be transferred", tradingPlatformID)
		return "token is burned and cannot be transferred", nil
	}

	if token.Listed {
		log.Printf("Token %s is listed and cannot be transferred", tradingPlatformID)
		return "token is listed and cannot be transferred", nil
	}

	callerID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		log.Printf("Failed to get client identity: %v", err)
		return "", fmt.Errorf("failed to get client identity: %v", err)
	}
	log.Printf("Caller ID for transfer: %s", callerID)

	if token.Owner != callerID && callerID != "marketplace-client-id" {
		log.Printf("Caller %s is not the owner or authorized entity for token %s", callerID, tradingPlatformID)
		return "caller is not authorized to transfer the token", nil
	}

	token.Owner = newOwner

	tokenBytes, err = json.Marshal(token)
	if err != nil {
		log.Printf("Failed to marshal token: %v", err)
		return "", fmt.Errorf("failed to marshal token: %v", err)
	}

	err = ctx.GetStub().PutState(tradingPlatformID, tokenBytes)
	if err != nil {
		log.Printf("Failed to update token in world state: %v", err)
		return "", fmt.Errorf("failed to update token in world state: %v", err)
	}

	eventPayload := map[string]interface{}{
		"action":            "TransferToken",
		"tradingPlatformID": tradingPlatformID,
		"newOwner":          newOwner,
	}
	eventPayloadBytes, _ := json.Marshal(eventPayload)
	err = ctx.GetStub().SetEvent("TokenEvent", eventPayloadBytes)
	if err != nil {
		log.Printf("Failed to emit event: %v", err)
		return "", fmt.Errorf("failed to emit event: %v", err)
	}

	log.Printf("TransferToken function completed successfully. New owner of token %s: %s", tradingPlatformID, newOwner)
	return "token successfully transferred", nil
}

func (m *MarketplaceChaincode) ListToken(ctx contractapi.TransactionContextInterface, tradingPlatformID string, listedUnitPrice string, tokenName string, matches []map[string]interface{}) (string, error) {
	log.Printf("Starting ListToken function with args: tradingPlatformID=%s, listedUnitPrice=%s, tokenName=%s, matches=%v", tradingPlatformID, listedUnitPrice, tokenName, matches)

	tokenBytes, err := ctx.GetStub().GetState(tradingPlatformID)
	if err != nil {
		return "", fmt.Errorf("failed to get token state: %v", err)
	}
	if tokenBytes == nil {
		return "token does not exist", nil
	}

	var token Token
	err = json.Unmarshal(tokenBytes, &token)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal token: %v", err)
	}

	log.Printf("Fetched token: %+v", token)

	if token.Burned {
		return "token is burned and cannot be listed", nil
	}

	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return "", fmt.Errorf("failed to get client identity: %v", err)
	}
	if token.Owner != clientID {
		return "caller is not the owner of the token", nil
	}

	token.Listed = true
	token.ListedUnitPrice = listedUnitPrice
	token.TokenName = tokenName

	tokenBytes, err = json.Marshal(token)
	if err != nil {
		return "", fmt.Errorf("failed to marshal updated token: %v", err)
	}
	err = ctx.GetStub().PutState(tradingPlatformID, tokenBytes)
	if err != nil {
		return "", fmt.Errorf("failed to update token in world state: %v", err)
	}

	log.Printf("Updated token state: %+v", token)

	matchResults, err := m.VerifyListingMatches(ctx, tradingPlatformID, tokenName, listedUnitPrice, matches)
	if err != nil {
		return "", fmt.Errorf("failed to verify listing matches: %v", err)
	}
	log.Printf("Match results: %+v", matchResults)

	var tradingPlatformIDs []string
	var orderIDs []string
	tradingPlatformIDs = append(tradingPlatformIDs, tradingPlatformID)

	filledQuantity := int64(0)

	for _, match := range matches {
		orderID, ok := match["OrderID"].(string)
		if !ok {
			log.Printf("Skipping match due to missing OrderID: %+v", match)
			continue
		}

		orderBytes, err := ctx.GetStub().GetState(orderID)
		if err != nil {
			log.Printf("Failed to fetch order state for OrderID: %s, error: %v", orderID, err)
			continue
		}
		if orderBytes == nil {
			log.Printf("No order found for OrderID: %s", orderID)
			continue
		}

		var order Order
		err = json.Unmarshal(orderBytes, &order)
		if err != nil {
			log.Printf("Failed to unmarshal order for OrderID: %s, error: %v", orderID, err)
			continue
		}

		log.Printf("Fetched order: %+v", order)

		isMatch, ok := matchResults[orderID]
		if !ok || !isMatch {
			log.Printf("OrderID: %s is not a valid match", orderID)
			continue
		}

		if order.Filled {
			log.Printf("OrderID: %s is already filled", orderID)
			continue
		}

		quantityToFill := int64(order.Quantity)
		if quantityToFill > (token.UnitCount - filledQuantity) {
			quantityToFill = token.UnitCount - filledQuantity
		}

		log.Printf("Quantity to fill for OrderID: %s is %d", orderID, quantityToFill)

		var filledWithID string

		if quantityToFill > 0 {
			firstNewUnitCount := quantityToFill
			secondNewUnitCount := token.UnitCount - firstNewUnitCount

			firstNewID, err := generateUniqueID(ctx, 1)
			if err != nil {
				return "", fmt.Errorf("failed to generate unique ID: %v", err)
			}

			firstNewToken := &Token{
				TradingPlatformID: firstNewID,
				WarehouseUnitID:   token.WarehouseUnitID,
				UnitHash:          token.UnitHash,
				Burned:            false,
				Owner:             order.Owner,
				Listed:            false,
				UnitCount:         firstNewUnitCount,
				ClaimedAt:         token.ClaimedAt,
				ListedUnitPrice:   "NA",
				TokenName:         "NA",
				SplitInto:         []string{},
				SplitFrom:         tradingPlatformID,
			}

			token.SplitInto = append(token.SplitInto, firstNewID)
			firstNewTokenBytes, err := json.Marshal(firstNewToken)
			if err != nil {
				return "", fmt.Errorf("failed to marshal new token: %v", err)
			}

			err = ctx.GetStub().PutState(firstNewID, firstNewTokenBytes)
			if err != nil {
				return "", fmt.Errorf("failed to put new token to world state: %v", err)
			}

			token.UnitCount = secondNewUnitCount
			if token.UnitCount == 0 {
				token.ListedUnitPrice = "NA"
				token.TokenName = "NA"
			}

			tokenBytes, err = json.Marshal(token)
			if err != nil {
				log.Printf("Failed to marshal updated token after splitting: %v", err)
				continue
			}

			err = ctx.GetStub().PutState(tradingPlatformID, tokenBytes)
			if err != nil {
				log.Printf("Failed to update token in world state after splitting: %v", err)
				continue
			}
			filledWithID = firstNewID
			tradingPlatformIDs = append(tradingPlatformIDs, firstNewID)
		}

		orderIDs = append(orderIDs, orderID)
		log.Printf("OrderID: %s, Quantity before fill: %d", orderID, order.Quantity)
		order.Quantity -= int(quantityToFill)
		log.Printf("OrderID: %s, Quantity after fill: %d", orderID, order.Quantity)

		order.FilledWith = append(order.FilledWith, filledWithID)

		totalPrice, err := strconv.ParseFloat(listedUnitPrice, 64)
		if err != nil {
			log.Printf("Failed to parse listedUnitPrice: %v", err)
			continue
		}
		totalPrice *= float64(quantityToFill)

		log.Printf("Total price for OrderID: %s is %f", orderID, totalPrice)

		order.EscrowedAmount -= int(totalPrice)
		log.Printf("OrderID: %s, EscrowedAmount after adjustment: %d", orderID, order.EscrowedAmount)

		sellerBalanceKey := token.Owner + "_" + order.PaymentToken
		sellerBalanceBytes, err := ctx.GetStub().GetState(sellerBalanceKey)
		if err != nil {
			log.Printf("Failed to read seller's balance from world state: %v", err)
			continue
		}

		var sellerBalance int
		if sellerBalanceBytes != nil {
			err = json.Unmarshal(sellerBalanceBytes, &sellerBalance)
			if err != nil {
				log.Printf("Failed to unmarshal seller's balance: %v", err)
				continue
			}
		}

		log.Printf("Seller's balance before update: %d", sellerBalance)

		sellerBalance += int(totalPrice)
		sellerBalanceBytes, err = json.Marshal(sellerBalance)
		if err != nil {
			log.Printf("Failed to marshal seller's balance: %v", err)
			continue
		}

		err = ctx.GetStub().PutState(sellerBalanceKey, sellerBalanceBytes)
		if err != nil {
			log.Printf("Failed to update seller's balance in world state: %v", err)
			continue
		}

		log.Printf("Seller's balance after update: %d", sellerBalance)

		if order.Quantity == 0 {
			order.Filled = true

			if order.EscrowedAmount > 0 {
				buyerBalanceKey := order.Owner + "_" + order.PaymentToken
				buyerBalanceBytes, err := ctx.GetStub().GetState(buyerBalanceKey)
				if err != nil {
					log.Printf("Failed to get buyer's balance: %v", err)
					continue
				}

				var buyerBalance int
				if buyerBalanceBytes != nil {
					err = json.Unmarshal(buyerBalanceBytes, &buyerBalance)
					if err != nil {
						log.Printf("Failed to unmarshal buyer's balance: %v", err)
						continue
					}
				}

				buyerBalance += order.EscrowedAmount
				buyerBalanceBytes, err = json.Marshal(buyerBalance)
				if err != nil {
					log.Printf("Failed to marshal buyer's balance: %v", err)
					continue
				}

				err = ctx.GetStub().PutState(buyerBalanceKey, buyerBalanceBytes)
				if err != nil {
					log.Printf("Failed to update buyer's balance in world state: %v", err)
					continue
				}

				log.Printf("Escrowed amount refunded to buyer: %d", order.EscrowedAmount)
				order.EscrowedAmount = 0
			}
		}

		orderBytes, err = json.Marshal(order)
		if err != nil {
			log.Printf("Failed to marshal updated order: %v", err)
			continue
		}

		err = ctx.GetStub().PutState(orderID, orderBytes)
		if err != nil {
			return "", fmt.Errorf("failed to update order in world state: %v", err)
		}

		filledQuantity += quantityToFill
		if filledQuantity >= token.UnitCount {
			break
		}
	}

	orderFillEventPayload := map[string]interface{}{
		"action":             "OrderFill",
		"tradingPlatformIDs": tradingPlatformIDs,
		"orderIDs":           orderIDs,
	}
	orderFillEventBytes, _ := json.Marshal(orderFillEventPayload)
	err = ctx.GetStub().SetEvent("OrderEvent", orderFillEventBytes)
	if err != nil {
		return "", fmt.Errorf("failed to emit OrderFill event: %v", err)
	}

	log.Println("ListToken function completed successfully")
	return "token successfully listed and orders fulfilled", nil
}

func (m *MarketplaceChaincode) DelistToken(ctx contractapi.TransactionContextInterface, tradingPlatformID string) (string, error) {
	log.Println("Starting DelistToken function")

	tokenBytes, err := ctx.GetStub().GetState(tradingPlatformID)
	if err != nil {
		log.Printf("Failed to read from world state: %v", err)
		return "", fmt.Errorf("failed to read from world state: %v", err)
	}
	if tokenBytes == nil {
		log.Printf("No token found for TradingPlatformID: %s", tradingPlatformID)
		return "token does not exist", nil
	}

	var token Token
	err = json.Unmarshal(tokenBytes, &token)
	if err != nil {
		log.Printf("Failed to unmarshal token: %v", err)
		return "", fmt.Errorf("failed to unmarshal token: %v", err)
	}

	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		log.Printf("Failed to get client identity: %v", err)
		return "", fmt.Errorf("failed to get client identity: %v", err)
	}
	log.Printf("Caller client ID: %s", clientID)
	if token.Owner != clientID {
		log.Printf("Caller is not the owner of token %s", tradingPlatformID)
		return "caller is not the owner of the token", nil
	}

	token.Listed = false
	token.ListedUnitPrice = ""
	token.TokenName = ""

	tokenBytes, err = json.Marshal(token)
	if err != nil {
		log.Printf("Failed to marshal token: %v", err)
		return "", fmt.Errorf("failed to marshal token: %v", err)
	}

	err = ctx.GetStub().PutState(tradingPlatformID, tokenBytes)
	if err != nil {
		log.Printf("Failed to update token in world state: %v", err)
		return "", fmt.Errorf("failed to update token in world state: %v", err)
	}

	eventPayload := map[string]interface{}{
		"action":            "DelistToken",
		"tradingPlatformID": tradingPlatformID,
	}
	eventPayloadBytes, _ := json.Marshal(eventPayload)
	err = ctx.GetStub().SetEvent("TokenEvent", eventPayloadBytes)
	if err != nil {
		log.Printf("Failed to emit event: %v", err)
		return "", fmt.Errorf("failed to emit event: %v", err)
	}

	log.Println("DelistToken function completed successfully")
	return "token successfully delisted", nil
}

func (m *MarketplaceChaincode) BuyToken(ctx contractapi.TransactionContextInterface, tradingPlatformID string) (string, error) {
	log.Println("Starting BuyToken function")

	tokenBytes, err := ctx.GetStub().GetState(tradingPlatformID)
	if err != nil {
		log.Printf("Failed to read from world state: %v", err)
		return "", fmt.Errorf("failed to read from world state: %v", err)
	}
	if tokenBytes == nil {
		log.Printf("No token found for TradingPlatformID: %s", tradingPlatformID)
		return "token does not exist", nil
	}

	var token Token
	err = json.Unmarshal(tokenBytes, &token)
	if err != nil {
		log.Printf("Failed to unmarshal token: %v", err)
		return "", fmt.Errorf("failed to unmarshal token: %v", err)
	}

	if !token.Listed {
		log.Printf("Token %s is not listed for sale", tradingPlatformID)
		return "token is not listed for sale", nil
	}

	buyerID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		log.Printf("Failed to get client identity: %v", err)
		return "", fmt.Errorf("failed to get client identity: %v", err)
	}
	log.Printf("Buyer client ID: %s", buyerID)

	buyerBalance, err := m.GetBalance(ctx, buyerID, token.TokenName)
	if err != nil {
		log.Printf("Failed to get buyer's balance: %v", err)
		return "", fmt.Errorf("failed to get buyer's balance: %v", err)
	}

	price, err := strconv.Atoi(token.ListedUnitPrice)
	if err != nil {
		log.Printf("Failed to convert listed price to int: %v", err)
		return "", fmt.Errorf("failed to convert listed price to int: %v", err)
	}

	if buyerBalance < price {
		log.Printf("Buyer %s has insufficient funds of %s", buyerID, token.TokenName)
		return "buyer has insufficient funds", nil
	}

	transferResult, err := m.TransferTokens(ctx, token.TokenName, token.Owner, price)
	if err != nil || transferResult != "tokens successfully transferred" {
		log.Printf("Failed to transfer tokens from buyer to seller: %v", err)
		return "", fmt.Errorf("failed to transfer tokens from buyer to seller: %v", err)
	}

	token.Listed = false
	token.Owner = buyerID
	token.ListedUnitPrice = ""
	token.TokenName = ""

	tokenBytes, err = json.Marshal(token)
	if err != nil {
		log.Printf("Failed to marshal token: %v", err)
		return "", fmt.Errorf("failed to marshal token: %v", err)
	}

	err = ctx.GetStub().PutState(tradingPlatformID, tokenBytes)
	if err != nil {
		log.Printf("Failed to update token in world state: %v", err)
		return "", fmt.Errorf("failed to update token in world state: %v", err)
	}

	eventPayload := map[string]interface{}{
		"action":            "BuyToken",
		"tradingPlatformID": tradingPlatformID,
		"buyerID":           buyerID,
	}
	eventPayloadBytes, _ := json.Marshal(eventPayload)
	err = ctx.GetStub().SetEvent("TokenEvent", eventPayloadBytes)
	if err != nil {
		log.Printf("Failed to emit event: %v", err)
		return "", fmt.Errorf("failed to emit event: %v", err)
	}

	log.Println("BuyToken function completed successfully")
	return "token successfully purchased", nil
}

func (m *MarketplaceChaincode) GetState(ctx contractapi.TransactionContextInterface, tradingPlatformID string) (*Token, error) {
	log.Printf("Starting GetState function for TradingPlatformID: %s", tradingPlatformID)

	tokenBytes, err := ctx.GetStub().GetState(tradingPlatformID)
	if err != nil {
		log.Printf("Failed to read from world state: %v", err)
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if tokenBytes == nil {
		log.Printf("No state found for TradingPlatformID: %s", tradingPlatformID)
		return nil, fmt.Errorf("no state found for TradingPlatformID: %s", tradingPlatformID)
	}

	var token Token
	err = json.Unmarshal(tokenBytes, &token)
	if err != nil {
		log.Printf("Failed to unmarshal token: %v", err)
		return nil, fmt.Errorf("failed to unmarshal token: %v", err)
	}

	if token.BurnedByOrgName == "" {
		token.BurnedByOrgName = "NA"
	}
	if token.BurnedByOrgURL == "" {
		token.BurnedByOrgURL = "NA"
	}
	if token.BurnedAt == 0 {
		token.BurnedAt = 0
	}
	if token.ListedUnitPrice == "" {
		token.ListedUnitPrice = "NA"
	}
	if token.TokenName == "" {
		token.TokenName = "NA"
	}

	log.Printf("GetState function completed successfully for TradingPlatformID: %s. Token data: %+v", tradingPlatformID, token)
	return &token, nil
}

func (m *MarketplaceChaincode) ListAllTokens(ctx contractapi.TransactionContextInterface) (string, error) {
	log.Println("Starting ListAllTokens function")
	startKey := ""
	endKey := ""

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)
	if err != nil {
		return "", err
	}
	defer resultsIterator.Close()

	var tokens []*Token
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return "", err
		}

		log.Printf("Raw query result: %s", string(queryResponse.Value))

		var token Token
		err = json.Unmarshal(queryResponse.Value, &token)
		if err != nil {
			log.Printf("Skipping non-token entry: %s", string(queryResponse.Value))
			continue
		}
		tokens = append(tokens, &token)
	}

	tokenBytes, err := json.Marshal(tokens)
	if err != nil {
		return "", err
	}

	log.Println("ListAllTokens function completed successfully")
	return string(tokenBytes), nil
}

func concatenateUnitData(unit *Unit) string {
	values := []string{
		unit.UnitBlockStart,
		unit.UnitBlockEnd,
		fmt.Sprintf("%d", unit.UnitCount),
		unit.WarehouseUnitID,
		unit.IssuanceID,
		unit.ProjectLocationID,
		unit.OrgUID,
		unit.UnitOwner,
		unit.CountryJurisdictionOfOwner,
		unit.InCountryJurisdictionOfOwner,
		unit.SerialNumberBlock,
		unit.SerialNumberPattern,
		fmt.Sprintf("%d", unit.VintageYear),
		unit.UnitType,
		unit.Marketplace,
		unit.MarketplaceLink,
		unit.MarketplaceIdentifier,
		unit.UnitTags,
		unit.UnitStatus,
		unit.UnitStatusReason,
		unit.UnitRegistryLink,
		unit.CorrespondingAdjustmentDeclaration,
		unit.CorrespondingAdjustmentStatus,
		fmt.Sprintf("%d", unit.TimeStaged),
		unit.CreatedAt,
		unit.UpdatedAt,
	}
	return strings.Join(values, "|")
}

func (m *MarketplaceChaincode) fetchUnitByID(ctx contractapi.TransactionContextInterface, warehouseUnitID string) (*Unit, error) {
	jsonResult, err := m.ConsensusQueryRow(ctx, `
        SELECT JSON_OBJECT(
            'UnitBlockStart', unitBlockStart,
            'UnitBlockEnd', unitBlockEnd,
            'UnitCount', unitCount,
            'WarehouseUnitID', warehouseUnitId,
            'IssuanceID', issuanceId,
            'ProjectLocationID', projectLocationId,
            'OrgUID', orgUid,
            'UnitOwner', unitOwner,
            'CountryJurisdictionOfOwner', countryJurisdictionOfOwner,
            'InCountryJurisdictionOfOwner', inCountryJurisdictionOfOwner,
            'SerialNumberBlock', serialNumberBlock,
            'SerialNumberPattern', serialNumberPattern,
            'VintageYear', vintageYear,
            'UnitType', unitType,
            'Marketplace', marketplace,
            'MarketplaceLink', marketplaceLink,
            'MarketplaceIdentifier', marketplaceIdentifier,
            'UnitTags', unitTags,
            'UnitStatus', unitStatus,
            'UnitStatusReason', unitStatusReason,
            'UnitRegistryLink', unitRegistryLink,
            'CorrespondingAdjustmentDeclaration', correspondingAdjustmentDeclaration,
            'CorrespondingAdjustmentStatus', correspondingAdjustmentStatus,
            'TimeStaged', timeStaged,
            'CreatedAt', createdAt,
            'UpdatedAt', updatedAt
        ) AS result
        FROM units WHERE warehouseUnitID = ?`, warehouseUnitID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch unit by ID: %v", err)
	}

	var unit Unit
	err = json.Unmarshal([]byte(jsonResult), &unit)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal unit JSON: %v", err)
	}

	return &unit, nil
}

func generateUniqueID(ctx contractapi.TransactionContextInterface, iterator int) (string, error) {
	log.Println("Generating unique ID using transaction timestamp")
	timestamp, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		log.Printf("Failed to get transaction timestamp: %v", err)
		return "", fmt.Errorf("failed to get transaction timestamp: %v", err)
	}
	iteratorBytes := []byte(strconv.Itoa(iterator))
	timeBytes := append([]byte(strconv.FormatInt(timestamp.Seconds, 10)), iteratorBytes...)
	hash := sha256.New()
	hash.Write(timeBytes)
	uniqueID := hex.EncodeToString(hash.Sum(nil))[:10]
	log.Printf("Generated unique ID: %s", uniqueID)
	return uniqueID, nil
}

func (m *MarketplaceChaincode) SplitUnits(ctx contractapi.TransactionContextInterface, tradingPlatformID string, method string, value int) (string, error) {
	log.Printf("Starting SplitUnits function for TradingPlatformID: %s", tradingPlatformID)

	tokenBytes, err := ctx.GetStub().GetState(tradingPlatformID)
	if err != nil {
		log.Printf("Failed to read from world state: %v", err)
		return "", fmt.Errorf("failed to read from world state: %v", err)
	}
	if tokenBytes == nil {
		log.Printf("No token found for TradingPlatformID: %s", tradingPlatformID)
		return "token does not exist", nil
	}

	var token Token
	err = json.Unmarshal(tokenBytes, &token)
	if err != nil {
		log.Printf("Failed to unmarshal token: %v", err)
		return "", fmt.Errorf("failed to unmarshal token: %v", err)
	}

	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		log.Printf("Failed to get client identity: %v", err)
		return "", fmt.Errorf("failed to get client identity: %v", err)
	}
	log.Printf("Caller client ID: %s", clientID)
	if token.Owner != clientID {
		log.Printf("Caller is not the owner of token %s", tradingPlatformID)
		return "caller is not the owner of the token", nil
	}

	if token.Listed {
		log.Printf("Token %s is listed and cannot be split", tradingPlatformID)
		return "token is listed and cannot be split", nil
	}

	if token.Burned {
		log.Printf("Token %s is burned and cannot be split", tradingPlatformID)
		return "token is burned and cannot be split", nil
	}

	log.Printf("Splitting token using method: %s, value: %d", method, value)
	var newTokens []*Token
	var newTokenIDs []string
	switch method {
	case "ratio":
		firstPart := (token.UnitCount * int64(value)) / 100
		secondPart := token.UnitCount - firstPart

		log.Printf("First part: %d, Second part: %d", firstPart, secondPart)
		for i := 0; i < 2; i++ {
			newUnitCount := firstPart
			if i == 1 {
				newUnitCount = secondPart
			}
			newID, err := generateUniqueID(ctx, i)
			if err != nil {
				return "", fmt.Errorf("failed to generate unique ID: %v", err)
			}
			newTokens = append(newTokens, &Token{
				TradingPlatformID: newID,
				WarehouseUnitID:   token.WarehouseUnitID,
				UnitHash:          token.UnitHash,
				Burned:            false,
				Owner:             clientID,
				Listed:            false,
				UnitCount:         newUnitCount,
				ClaimedAt:         token.ClaimedAt,
				SplitInto:         []string{},
				SplitFrom:         tradingPlatformID,
			})
			newTokenIDs = append(newTokenIDs, newID)
		}

	case "divide":
		if value <= 0 || token.UnitCount%int64(value) != 0 {
			return "", fmt.Errorf("invalid value for division")
		}
		partCount := token.UnitCount / int64(value)

		log.Printf("Splitting into %d parts, each with %d units", value, partCount)
		for i := 0; i < value; i++ {
			newID, err := generateUniqueID(ctx, i)
			if err != nil {
				return "", fmt.Errorf("failed to generate unique ID: %v", err)
			}
			newTokens = append(newTokens, &Token{
				TradingPlatformID: newID,
				WarehouseUnitID:   token.WarehouseUnitID,
				UnitHash:          token.UnitHash,
				Burned:            false,
				Owner:             clientID,
				Listed:            false,
				UnitCount:         partCount,
				ClaimedAt:         token.ClaimedAt,
				SplitInto:         []string{},
				SplitFrom:         tradingPlatformID,
			})
			newTokenIDs = append(newTokenIDs, newID)
		}

	default:
		return "", fmt.Errorf("invalid split method")
	}

	token.UnitCount = 0
	token.SplitInto = newTokenIDs

	tokenBytes, err = json.Marshal(token)
	if err != nil {
		log.Printf("Failed to marshal original token: %v", err)
		return "", fmt.Errorf("failed to marshal original token: %v", err)
	}

	err = ctx.GetStub().PutState(tradingPlatformID, tokenBytes)
	if err != nil {
		log.Printf("Failed to update original token in world state: %v", err)
		return "", fmt.Errorf("failed to update original token in world state: %v", err)
	}

	for _, newToken := range newTokens {
		newTokenBytes, err := json.Marshal(newToken)
		if err != nil {
			log.Printf("Failed to marshal new token: %v", err)
			return "", fmt.Errorf("failed to marshal new token: %v", err)
		}

		err = ctx.GetStub().PutState(newToken.TradingPlatformID, newTokenBytes)
		if err != nil {
			log.Printf("Failed to put new token to world state: %v", err)
			return "", fmt.Errorf("failed to put new token to world state: %v", err)
		}
	}

	eventPayload := map[string]interface{}{
		"action":                "SplitUnits",
		"originalPlatformID":    tradingPlatformID,
		"newTradingPlatformIDs": newTokenIDs,
	}
	eventPayloadBytes, _ := json.Marshal(eventPayload)
	err = ctx.GetStub().SetEvent("TokenEvent", eventPayloadBytes)
	if err != nil {
		log.Printf("Failed to emit event: %v", err)
		return "", fmt.Errorf("failed to emit event: %v", err)
	}

	log.Println("SplitUnits function completed successfully")
	return "units successfully split", nil
}

func (m *MarketplaceChaincode) CreateOrder(ctx contractapi.TransactionContextInterface, paymentToken string, maxPrice float64, quantity int, vintageYearRange string, orderType string, coefficients map[string]map[string]float64, matches []map[string]interface{}) (string, error) {
	log.Printf("Starting CreateOrder function with args: paymentToken=%s, maxPrice=%f, quantity=%d, vintageYearRange=%s, orderType=%s, coefficients=%+v, matches=%+v", paymentToken, maxPrice, quantity, vintageYearRange, orderType, coefficients, matches)

	orderID, err := generateUniqueID(ctx, 0)
	if err != nil {
		log.Printf("Failed to generate unique OrderID: %v", err)
		return "", fmt.Errorf("failed to generate unique OrderID: %v", err)
	}
	log.Printf("Generated OrderID: %s", orderID)

	owner, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		log.Printf("Failed to get client identity: %v", err)
		return "", fmt.Errorf("failed to get client identity: %v", err)
	}
	log.Printf("Order owner: %s", owner)

	if coefficients == nil {
		coefficients = make(map[string]map[string]float64)
		log.Println("Coefficients were nil, initialized an empty map")
	}

	totalEscrowAmount := int(maxPrice * float64(quantity))
	log.Printf("Total escrow amount calculated: %d", totalEscrowAmount)

	balanceKey := owner + "_" + paymentToken
	balanceBytes, err := ctx.GetStub().GetState(balanceKey)
	if err != nil {
		log.Printf("Failed to read balance from world state: %v", err)
		return "", fmt.Errorf("failed to read balance from world state: %v", err)
	}

	var balance int
	if balanceBytes != nil {
		err = json.Unmarshal(balanceBytes, &balance)
		if err != nil {
			log.Printf("Failed to unmarshal balance: %v", err)
			return "", fmt.Errorf("failed to unmarshal balance: %v", err)
		}
		log.Printf("Balance for key %s: %d", balanceKey, balance)
	} else {
		log.Printf("No existing balance found for key %s", balanceKey)
	}

	if balance < totalEscrowAmount {
		log.Printf("Insufficient balance. Required: %d, Available: %d", totalEscrowAmount, balance)
		return "", fmt.Errorf("insufficient balance to create order")
	}

	balance -= totalEscrowAmount
	log.Printf("Balance after deducting escrow amount: %d", balance)

	balanceBytes, err = json.Marshal(balance)
	if err != nil {
		log.Printf("Failed to marshal balance: %v", err)
		return "", fmt.Errorf("failed to marshal balance: %v", err)
	}

	err = ctx.GetStub().PutState(balanceKey, balanceBytes)
	if err != nil {
		log.Printf("Failed to update balance in world state: %v", err)
		return "", fmt.Errorf("failed to update balance in world state: %v", err)
	}

	order := &Order{
		OrderID:          orderID,
		PaymentToken:     paymentToken,
		MaxPrice:         maxPrice,
		Quantity:         quantity,
		VintageYearRange: vintageYearRange,
		OrderType:        orderType,
		Coefficients:     coefficients,
		Owner:            owner,
		EscrowedAmount:   totalEscrowAmount,
		Filled:           false,
		FilledWith:       []string{},
	}

	orderBytes, err := json.Marshal(order)
	if err != nil {
		log.Printf("Failed to marshal order: %v", err)
		return "", fmt.Errorf("failed to marshal order: %v", err)
	}

	err = ctx.GetStub().PutState(orderID, orderBytes)
	if err != nil {
		log.Printf("Failed to put order to world state: %v", err)
		return "", fmt.Errorf("failed to put order to world state: %v", err)
	}
	log.Printf("Order saved in world state: %+v", order)

	var tradingPlatformIDs []string
	filledQuantity := 0

	for _, match := range matches {
		tradingPlatformID := match["TradingPlatformID"].(string)
		originalTradingPlatformID := tradingPlatformID
		desiredPrice := match["DesiredPrice"].(float64)

		log.Printf("Processing match: TradingPlatformID=%s, DesiredPrice=%f", tradingPlatformID, desiredPrice)

		tokenBytes, err := ctx.GetStub().GetState(tradingPlatformID)
		if err != nil {
			log.Printf("Failed to read from world state: %v", err)
			continue
		}
		if tokenBytes == nil {
			log.Printf("No token found for TradingPlatformID: %s", tradingPlatformID)
			continue
		}

		var token Token
		err = json.Unmarshal(tokenBytes, &token)
		if err != nil {
			log.Printf("Failed to unmarshal token: %v", err)
			continue
		}
		log.Printf("Fetched token: %+v", token)

		listedPrice, err := strconv.ParseFloat(token.ListedUnitPrice, 64)
		if err != nil {
			log.Printf("Failed to convert listed price to float: %v", err)
			continue
		}
		log.Printf("Token listed price: %f", listedPrice)

		if !token.Listed || token.TokenName != paymentToken || listedPrice > desiredPrice {
			log.Printf("Token does not match criteria or listed price is too high")
			continue
		}

		quantityToFill := int(token.UnitCount)
		if quantityToFill > (order.Quantity - filledQuantity) {
			quantityToFill = order.Quantity - filledQuantity
			log.Printf("Adjusting quantity to fill: %d", quantityToFill)

			firstNewUnitCount := int64(quantityToFill)
			secondNewUnitCount := token.UnitCount - firstNewUnitCount

			firstNewID, err := generateUniqueID(ctx, 1)
			if err != nil {
				log.Printf("Failed to generate unique ID for split token: %v", err)
				return "", fmt.Errorf("failed to generate unique ID: %v", err)
			}

			firstNewToken := &Token{
				TradingPlatformID: firstNewID,
				WarehouseUnitID:   token.WarehouseUnitID,
				UnitHash:          token.UnitHash,
				Burned:            false,
				Owner:             owner,
				Listed:            false,
				UnitCount:         firstNewUnitCount,
				ClaimedAt:         token.ClaimedAt,
				ListedUnitPrice:   "NA",
				TokenName:         "NA",
				SplitInto:         []string{},
				SplitFrom:         tradingPlatformID,
			}

			token.SplitInto = append(token.SplitInto, firstNewID)
			log.Printf("Token split into new token with ID: %s", firstNewID)

			firstNewTokenBytes, err := json.Marshal(firstNewToken)
			if err != nil {
				log.Printf("Failed to marshal new token: %v", err)
				return "", fmt.Errorf("failed to marshal new token: %v", err)
			}

			err = ctx.GetStub().PutState(firstNewID, firstNewTokenBytes)
			if err != nil {
				log.Printf("Failed to put new token to world state: %v", err)
				return "", fmt.Errorf("failed to put new token to world state: %v", err)
			}

			token.UnitCount = secondNewUnitCount

			tokenBytes, err = json.Marshal(token)
			if err != nil {
				log.Printf("Failed to marshal updated original token: %v", err)
				continue
			}

			err = ctx.GetStub().PutState(tradingPlatformID, tokenBytes)
			if err != nil {
				log.Printf("Failed to update original token in world state: %v", err)
				continue
			}

			tradingPlatformID = firstNewID
			token = *firstNewToken
		}

		if originalTradingPlatformID != tradingPlatformID {
			tradingPlatformIDs = append(tradingPlatformIDs, originalTradingPlatformID)
		}

		tradingPlatformIDs = append(tradingPlatformIDs, tradingPlatformID)
		log.Printf("Trading platform IDs updated: %+v", tradingPlatformIDs)

		amountToTransfer := int(listedPrice * float64(quantityToFill))
		log.Printf("Amount to transfer for token: %d", amountToTransfer)

		totalEscrowAmount -= amountToTransfer
		log.Printf("Remaining escrow amount after transfer: %d", totalEscrowAmount)

		sellerBalanceKey := token.Owner + "_" + paymentToken
		sellerBalanceBytes, err := ctx.GetStub().GetState(sellerBalanceKey)
		if err != nil {
			log.Printf("Failed to read seller's balance from world state: %v", err)
			continue
		}

		var sellerBalance int
		if sellerBalanceBytes != nil {
			err = json.Unmarshal(sellerBalanceBytes, &sellerBalance)
			if err != nil {
				log.Printf("Failed to unmarshal seller's balance: %v", err)
				continue
			}
		}
		log.Printf("Seller's balance before update: %d", sellerBalance)

		sellerBalance += amountToTransfer
		sellerBalanceBytes, err = json.Marshal(sellerBalance)
		if err != nil {
			log.Printf("Failed to marshal seller's balance: %v", err)
			continue
		}

		err = ctx.GetStub().PutState(sellerBalanceKey, sellerBalanceBytes)
		if err != nil {
			log.Printf("Failed to update seller's balance in world state: %v", err)
			continue
		}
		log.Printf("Seller's balance updated to: %d", sellerBalance)

		token.Listed = false
		token.Owner = owner
		token.ListedUnitPrice = ""
		token.TokenName = ""

		tokenBytes, err = json.Marshal(token)
		if err != nil {
			log.Printf("Failed to marshal token: %v", err)
			continue
		}

		err = ctx.GetStub().PutState(tradingPlatformID, tokenBytes)
		if err != nil {
			log.Printf("Failed to update token in world state: %v", err)
			continue
		}
		log.Printf("Token state updated after filling order: %+v", token)

		filledQuantity += quantityToFill
		log.Printf("Filled quantity updated: %d", filledQuantity)

		order.Quantity -= quantityToFill
		order.FilledWith = append(order.FilledWith, tradingPlatformID)
		log.Printf("Order filled with token ID: %s", tradingPlatformID)

		if filledQuantity >= quantity {
			break
		}
	}

	if filledQuantity < quantity {
		order.EscrowedAmount = totalEscrowAmount
		log.Printf("Escrowed amount after partial order fill: %d", order.EscrowedAmount)
	} else {
		remainingEscrowAmount := order.EscrowedAmount

		if remainingEscrowAmount > 0 {

			balance += remainingEscrowAmount
			balanceBytes, err := json.Marshal(balance)
			if err != nil {
				log.Printf("Failed to marshal updated balance for refund: %v", err)
				return "", fmt.Errorf("failed to marshal balance: %v", err)
			}

			err = ctx.GetStub().PutState(balanceKey, balanceBytes)
			if err != nil {
				log.Printf("Failed to refund remaining escrow to buyer: %v", err)
				return "", fmt.Errorf("failed to refund remaining escrow to buyer: %v", err)
			}
			log.Printf("Refunded remaining escrow of %d to buyer", remainingEscrowAmount)
		}

		order.EscrowedAmount = 0
		log.Printf("Escrowed amount after full order fill: %d", order.EscrowedAmount)

		order.Filled = true
		log.Printf("Order is fully filled")
	}

	orderBytes, err = json.Marshal(order)
	if err != nil {
		log.Printf("Failed to marshal updated order: %v", err)
		return "", fmt.Errorf("failed to marshal updated order: %v", err)
	}

	err = ctx.GetStub().PutState(orderID, orderBytes)
	if err != nil {
		log.Printf("Failed to update order in world state: %v", err)
		return "", fmt.Errorf("failed to update order in world state: %v", err)
	}

	orderBytes, err = json.Marshal(order)
	if err != nil {
		log.Printf("Failed to marshal updated order: %v", err)
		return "", fmt.Errorf("failed to marshal updated order: %v", err)
	}

	err = ctx.GetStub().PutState(orderID, orderBytes)
	if err != nil {
		log.Printf("Failed to update order in world state: %v", err)
		return "", fmt.Errorf("failed to update order in world state: %v", err)
	}
	log.Printf("Final order state updated: %+v", order)

	orderFillEventPayload := map[string]interface{}{
		"action":             "OrderFill",
		"orderIDs":           []string{orderID},
		"tradingPlatformIDs": tradingPlatformIDs,
	}
	orderFillEventBytes, _ := json.Marshal(orderFillEventPayload)
	err = ctx.GetStub().SetEvent("OrderEvent", orderFillEventBytes)
	if err != nil {
		log.Printf("Failed to emit OrderFill event: %v", err)
		return "", fmt.Errorf("failed to emit OrderFill event: %v", err)
	}
	log.Println("OrderFill event emitted successfully")

	log.Println("CreateOrder function completed successfully")
	return "order successfully created", nil
}

func (m *MarketplaceChaincode) RemoveOrder(ctx contractapi.TransactionContextInterface, orderID string) (string, error) {
	log.Println("Starting RemoveOrder function")

	orderBytes, err := ctx.GetStub().GetState(orderID)
	if err != nil {
		log.Printf("Failed to read order from world state: %v", err)
		return "", fmt.Errorf("failed to read order from world state: %v", err)
	}
	if orderBytes == nil {
		log.Printf("No order found for OrderID: %s", orderID)
		return "order does not exist", nil
	}

	var order Order
	err = json.Unmarshal(orderBytes, &order)
	if err != nil {
		log.Printf("Failed to unmarshal order: %v", err)
		return "", fmt.Errorf("failed to unmarshal order: %v", err)
	}

	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		log.Printf("Failed to get client identity: %v", err)
		return "", fmt.Errorf("failed to get client identity: %v", err)
	}
	log.Printf("Caller client ID: %s", clientID)
	if order.Owner != clientID {
		log.Printf("Caller is not the owner of order %s", orderID)
		return "caller is not the owner of the order", nil
	}

	order.Filled = true

	balanceKey := order.Owner + "_" + order.PaymentToken
	balanceBytes, err := ctx.GetStub().GetState(balanceKey)
	if err != nil {
		return "", fmt.Errorf("failed to read balance from world state: %v", err)
	}

	var balance int
	if balanceBytes != nil {
		err = json.Unmarshal(balanceBytes, &balance)
		if err != nil {
			return "", fmt.Errorf("failed to unmarshal balance: %v", err)
		}
	}

	balance += order.EscrowedAmount
	order.EscrowedAmount = 0

	balanceBytes, err = json.Marshal(balance)
	if err != nil {
		return "", fmt.Errorf("failed to marshal balance: %v", err)
	}

	err = ctx.GetStub().PutState(balanceKey, balanceBytes)
	if err != nil {
		return "", fmt.Errorf("failed to update balance in world state: %v", err)
	}

	orderBytes, err = json.Marshal(order)
	if err != nil {
		log.Printf("Failed to marshal updated order: %v", err)
		return "", fmt.Errorf("failed to marshal updated order: %v", err)
	}

	err = ctx.GetStub().PutState(orderID, orderBytes)
	if err != nil {
		log.Printf("Failed to update order in world state: %v", err)
		return "", fmt.Errorf("failed to update order in world state: %v", err)
	}

	eventPayload := map[string]interface{}{
		"action":  "RemoveOrder",
		"orderID": orderID,
	}
	eventPayloadBytes, _ := json.Marshal(eventPayload)
	err = ctx.GetStub().SetEvent("OrderEvent", eventPayloadBytes)
	if err != nil {
		log.Printf("Failed to emit event: %v", err)
		return "", fmt.Errorf("failed to emit event: %v", err)
	}

	log.Println("RemoveOrder function completed successfully")
	return "order successfully removed and escrow amount returned", nil
}

func (m *MarketplaceChaincode) GetOrder(ctx contractapi.TransactionContextInterface, orderID string) (*Order, error) {
	log.Printf("Starting GetOrder function for OrderID: %s", orderID)

	orderBytes, err := ctx.GetStub().GetState(orderID)
	if err != nil {
		log.Printf("Failed to read from world state: %v", err)
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if orderBytes == nil {
		log.Printf("No state found for OrderID: %s", orderID)
		return nil, fmt.Errorf("no state found for OrderID: %s", orderID)
	}

	var order Order
	err = json.Unmarshal(orderBytes, &order)
	if err != nil {
		log.Printf("Failed to unmarshal order: %v", err)
		return nil, fmt.Errorf("failed to unmarshal order: %v", err)
	}

	log.Printf("GetOrder function completed successfully for OrderID: %s. Order data: %+v", orderID, order)
	return &order, nil
}

func (m *MarketplaceChaincode) VerifyListingMatches(ctx contractapi.TransactionContextInterface, tradingPlatformID string, tokenName string, listedUnitPrice string, matches []map[string]interface{}) (map[string]bool, error) {
	log.Printf("Starting VerifyListingMatches function for TradingPlatformID: %s with tokenName: %s, listedUnitPrice: %s, matches: %+v", tradingPlatformID, tokenName, listedUnitPrice, matches)

	tokenBytes, err := ctx.GetStub().GetState(tradingPlatformID)
	if err != nil {
		log.Printf("Failed to read token from world state: %v", err)
		return nil, fmt.Errorf("failed to read token from world state: %v", err)
	}
	if tokenBytes == nil {
		log.Printf("No token found for TradingPlatformID: %s", tradingPlatformID)
		return nil, fmt.Errorf("token does not exist")
	}

	var token Token
	log.Printf("Unmarshalling token data")
	err = json.Unmarshal(tokenBytes, &token)
	if err != nil {
		log.Printf("Failed to unmarshal token: %v", err)
		return nil, fmt.Errorf("failed to unmarshal token: %v", err)
	}

	log.Printf("Fetched token: %+v", token)

	registry, sector, country, vintageYear, err := m.fetchRegistrySectorCountryAndVintage(ctx, token.WarehouseUnitID)
	if err != nil {
		log.Printf("Failed to fetch registry, sector, country, and vintage year for WarehouseUnitID %s: %v", token.WarehouseUnitID, err)
		return nil, fmt.Errorf("failed to fetch registry, sector, country, and vintage year: %v", err)
	}

	log.Printf("Registry: %s, Sector: %s, Country: %s, VintageYear: %d", registry, sector, country, vintageYear)

	results := make(map[string]bool)

	for _, match := range matches {
		orderID, ok := match["OrderID"].(string)
		if !ok {
			log.Printf("Failed to parse OrderID in match: %+v", match)
			continue
		}

		log.Printf("Fetching order state for OrderID: %s", orderID)
		orderBytes, err := ctx.GetStub().GetState(orderID)
		if err != nil {
			log.Printf("Failed to read order from world state: %v", err)
			continue
		}
		if orderBytes == nil {
			log.Printf("No order found for OrderID: %s", orderID)
			continue
		}

		var order Order
		log.Printf("Unmarshalling order data")
		err = json.Unmarshal(orderBytes, &order)
		if err != nil {
			log.Printf("Failed to unmarshal order: %v", err)
			continue
		}

		log.Printf("Fetched order: %+v", order)

		if tokenName != order.PaymentToken {
			log.Printf("Payment token mismatch: TokenName = %s, Order PaymentToken = %s", tokenName, order.PaymentToken)
			results[orderID] = false
			continue
		}

		if order.OrderType == "Information Agnostic" {
			vintageYearStr := strconv.FormatInt(vintageYear, 10)
			log.Printf("Comparing vintage year: Token's VintageYear = %s", vintageYearStr)

			vintageYearRange := order.VintageYearRange
			if vintageYearRange != "" {
				years := strings.Split(vintageYearRange, " - ")
				if len(years) == 2 {
					startYear, errStart := strconv.Atoi(years[0])
					endYear, errEnd := strconv.Atoi(years[1])

					if errStart != nil || errEnd != nil || vintageYear < int64(startYear) || vintageYear > int64(endYear) {
						log.Printf("Vintage year %d is outside the range %d - %d", vintageYear, startYear, endYear)
						results[orderID] = false
						continue
					}
				} else {
					log.Printf("Invalid vintage year range format: %s", vintageYearRange)
					results[orderID] = false
					continue
				}
			}

			log.Printf("Order is Information Agnostic. Skipping registry, sector, and country checks, using order.MaxPrice.")

			listedPrice, err := strconv.ParseFloat(listedUnitPrice, 64)
			if err != nil {
				log.Printf("Failed to parse listed price: %v", err)
				results[orderID] = false
				continue
			}

			if order.MaxPrice >= listedPrice {
				results[orderID] = true
				log.Printf("OrderID: %s is a match. MaxPrice %f >= ListedPrice %f", orderID, order.MaxPrice, listedPrice)
			} else {
				results[orderID] = false
				log.Printf("OrderID: %s is not a match. MaxPrice %f < ListedPrice %f", orderID, order.MaxPrice, listedPrice)
			}
			continue
		}

		matchFound := true

		if registryCoefficient, exists := order.Coefficients["Registry"]; exists {
			if _, registryMatches := registryCoefficient[registry]; !registryMatches {
				log.Printf("Registry mismatch: Expected %s, Got %s", registry, registry)
				matchFound = false
			}
		} else {
			log.Printf("No registry found in order coefficients")
			matchFound = false
		}

		if matchFound {
			if sectorCoefficient, exists := order.Coefficients["Sector"]; exists {
				if _, sectorMatches := sectorCoefficient[sector]; !sectorMatches {
					log.Printf("Sector mismatch: Expected %s, Got %s", sector, sector)
					matchFound = false
				}
			} else {
				log.Printf("No sector found in order coefficients")
				matchFound = false
			}
		}

		if matchFound {
			if countryCoefficient, exists := order.Coefficients["Location"]; exists {
				if _, countryMatches := countryCoefficient[country]; !countryMatches {
					log.Printf("Country mismatch: Expected %s, Got %s", country, country)
					matchFound = false
				}
			} else {
				log.Printf("No country found in order coefficients")
				matchFound = false
			}
		}

		if matchFound {
			vintageYearStr := strconv.FormatInt(vintageYear, 10)
			log.Printf("Comparing vintage year: Token's VintageYear = %s", vintageYearStr)

			vintageYearRange := order.VintageYearRange
			if vintageYearRange != "" {
				years := strings.Split(vintageYearRange, " - ")
				if len(years) == 2 {
					startYear, errStart := strconv.Atoi(years[0])
					endYear, errEnd := strconv.Atoi(years[1])

					if errStart != nil || errEnd != nil || vintageYear < int64(startYear) || vintageYear > int64(endYear) {
						log.Printf("Vintage year %d is outside the range %d - %d", vintageYear, startYear, endYear)
						matchFound = false
					}
				} else {
					log.Printf("Invalid vintage year range format: %s", vintageYearRange)
					matchFound = false
				}
			} else {
				log.Printf("No vintage year range found in order")
				matchFound = false
			}
		}

		if matchFound {
			log.Printf("Calculating price")
			calculatedPrice := calculatePrice(order.Coefficients, registry, sector, country, order.MaxPrice)

			listedPrice, err := strconv.ParseFloat(listedUnitPrice, 64)
			if err != nil {
				log.Printf("Failed to parse listed price: %v", err)
				results[orderID] = false
				continue
			}

			if calculatedPrice >= listedPrice {
				results[orderID] = true
				log.Printf("OrderID: %s is a match with calculated price %f >= listed price %f", orderID, calculatedPrice, listedPrice)
			} else {
				results[orderID] = false
				log.Printf("OrderID: %s is not a match, calculated price %f < listed price %f", orderID, calculatedPrice, listedPrice)
			}
		} else {
			results[orderID] = false
		}
	}

	log.Printf("VerifyListingMatches function completed for TradingPlatformID: %s with results: %+v", tradingPlatformID, results)
	return results, nil
}

func (m *MarketplaceChaincode) fetchRegistrySectorCountryAndVintage(ctx contractapi.TransactionContextInterface, warehouseUnitID string) (string, string, string, int64, error) {
	log.Printf("Starting fetchRegistrySectorCountryAndVintage for WarehouseUnitID: %s", warehouseUnitID)

	unitQuery := `
        SELECT JSON_OBJECT(
            'IssuanceID', IssuanceID,
            'VintageYear', vintageYear
        ) AS result
        FROM units
        WHERE WarehouseUnitID = ?`

	jsonResult, err := m.ConsensusQueryRow(ctx, unitQuery, warehouseUnitID)
	if err != nil {
		log.Printf("Failed to fetch issuance ID and vintage year: %v", err)
		return "", "", "", 0, fmt.Errorf("failed to fetch issuance ID and vintage year: %v", err)
	}

	var unitData struct {
		IssuanceID  string `json:"IssuanceID"`
		VintageYear int64  `json:"VintageYear"`
	}
	err = json.Unmarshal([]byte(jsonResult), &unitData)
	if err != nil {
		log.Printf("Failed to unmarshal unit JSON: %v", err)
		return "", "", "", 0, fmt.Errorf("failed to unmarshal unit JSON: %v", err)
	}

	log.Printf("Fetched issuanceID: %s, vintageYear: %d", unitData.IssuanceID, unitData.VintageYear)

	issuanceQuery := `
        SELECT JSON_OBJECT(
            'WarehouseProjectID', WarehouseProjectId
        ) AS result
        FROM issuances
        WHERE id = ?`

	jsonResult, err = m.ConsensusQueryRow(ctx, issuanceQuery, unitData.IssuanceID)
	if err != nil {
		log.Printf("Failed to fetch WarehouseProjectId: %v", err)
		return "", "", "", 0, fmt.Errorf("failed to fetch WarehouseProjectId: %v", err)
	}

	var issuanceData struct {
		WarehouseProjectID string `json:"WarehouseProjectID"`
	}
	err = json.Unmarshal([]byte(jsonResult), &issuanceData)
	if err != nil {
		log.Printf("Failed to unmarshal issuance JSON: %v", err)
		return "", "", "", 0, fmt.Errorf("failed to unmarshal issuance JSON: %v", err)
	}

	log.Printf("Fetched warehouseProjectID: %s", issuanceData.WarehouseProjectID)

	projectQuery := `
        SELECT JSON_OBJECT(
            'CurrentRegistry', p.CurrentRegistry,
            'Sector', p.Sector,
            'Country', l.Country
        ) AS result
        FROM projects p
        JOIN projectlocations l ON p.WarehouseProjectId = l.WarehouseProjectId
        WHERE p.WarehouseProjectId = ?`

	jsonResult, err = m.ConsensusQueryRow(ctx, projectQuery, issuanceData.WarehouseProjectID)
	if err != nil {
		log.Printf("Failed to fetch registry, sector, and country: %v", err)
		return "", "", "", 0, fmt.Errorf("failed to fetch registry, sector, and country: %v", err)
	}

	var projectData struct {
		CurrentRegistry string `json:"CurrentRegistry"`
		Sector          string `json:"Sector"`
		Country         string `json:"Country"`
	}
	err = json.Unmarshal([]byte(jsonResult), &projectData)
	if err != nil {
		log.Printf("Failed to unmarshal project JSON: %v", err)
		return "", "", "", 0, fmt.Errorf("failed to unmarshal project JSON: %v", err)
	}

	log.Printf("Fetched registry: %s, sector: %s, country: %s", projectData.CurrentRegistry, projectData.Sector, projectData.Country)

	region, err := getRegionFromCountry(ctx, projectData.Country)
	if err != nil {
		log.Printf("Failed to get region from country: %v", err)
		return "", "", "", 0, fmt.Errorf("failed to get region from country: %v", err)
	}

	log.Printf("Determined region: %s", region)

	return projectData.CurrentRegistry, projectData.Sector, region, unitData.VintageYear, nil
}

func calculatePrice(coefficients map[string]map[string]float64, registry, sector, country string, maxPrice float64) float64 {
	log.Printf("Calculating price with registry=%s, sector=%s, country=%s, maxPrice=%f", registry, sector, country, maxPrice)

	basePrice := coefficients["registry"][registry] *
		coefficients["sector"][sector] *
		coefficients["country"][country]

	finalPrice := basePrice + maxPrice

	log.Printf("Calculated final price: %f", finalPrice)

	return finalPrice
}

func getRegionFromCountry(ctx contractapi.TransactionContextInterface, country string) (string, error) {
	log.Printf("Determining region for country: %s", country)

	regionCountries := map[string][]string{
		"Africa":         {"Algeria", "Angola", "Benin", "Botswana", "Burkina Faso", "Burundi", "Cabo Verde", "Cameroon", "Central African Republic", "Chad", "Comoros", "Congo", "Côte d'Ivoire", "Democratic Republic of Congo", "Djibouti", "Egypt", "Equatorial Guinea", "Eritrea", "Eswatini", "Ethiopia", "Gabon", "Gambia", "Ghana", "Guinea", "Guinea-Bissau", "Kenya", "Lesotho", "Liberia", "Libya", "Madagascar", "Malawi", "Mali", "Mauritania", "Mauritius", "Morocco", "Mozambique", "Namibia", "Niger", "Nigeria", "Rwanda", "Sao Tome and Principe", "Senegal", "Seychelles", "Sierra Leone", "Somalia", "South Africa", "South Sudan", "Sudan", "Togo", "Tunisia", "Uganda", "United Republic of Tanzania", "Zambia", "Zimbabwe"},
		"Asia":           {"Afghanistan", "Armenia", "Azerbaijan", "Bahrain", "Bangladesh", "Bhutan", "Brunei Darussalam", "Cambodia", "China", "Georgia", "India", "Indonesia", "Iran", "Iraq", "Israel", "Japan", "Jordan", "Kazakhstan", "Kuwait", "Kyrgyzstan", "Lao People's Democratic Republic", "Lebanon", "Malaysia", "Maldives", "Mongolia", "Myanmar", "Nepal", "Oman", "Pakistan", "Philippines", "Qatar", "Republic of Korea", "Saudi Arabia", "Singapore", "Sri Lanka", "Syrian Arab Republic", "Tajikistan", "Thailand", "Timor-Leste", "Turkmenistan", "United Arab Emirates", "Uzbekistan", "Viet Nam", "Yemen"},
		"Europe":         {"Albania", "Andorra", "Austria", "Belarus", "Belgium", "Bosnia and Herzegovina", "Bulgaria", "Croatia", "Cyprus", "Czech Republic", "Denmark", "Estonia", "Finland", "France", "Germany", "Greece", "Hungary", "Iceland", "Ireland", "Italy", "Latvia", "Liechtenstein", "Lithuania", "Luxembourg", "Malta", "Monaco", "Montenegro", "Netherlands", "North Macedonia", "Norway", "Poland", "Portugal", "Republic of Moldova", "Romania", "Russian Federation", "San Marino", "Serbia", "Slovakia", "Slovenia", "Spain", "Sweden", "Switzerland", "Turkey", "Ukraine", "United Kingdom"},
		"Americas":       {"Antigua and Barbuda", "Argentina", "Bahamas", "Barbados", "Belize", "Bolivia", "Brazil", "Canada", "Chile", "Colombia", "Costa Rica", "Cuba", "Dominica", "Dominican Republic", "Ecuador", "El Salvador", "Grenada", "Guatemala", "Guyana", "Haiti", "Honduras", "Jamaica", "Mexico", "Nicaragua", "Panama", "Paraguay", "Peru", "Puerto Rico", "Saint Kitts and Nevis", "Saint Lucia", "Saint Vincent and the Grenadines", "Suriname", "Trinidad and Tobago", "United States of America", "Uruguay", "Venezuela"},
		"Oceania":        {"Australia", "Fiji", "Kiribati", "Marshall Islands", "Micronesia", "Nauru", "New Zealand", "Palau", "Papua New Guinea", "Samoa", "Solomon Islands", "Tonga", "Tuvalu", "Vanuatu"},
		"European Union": {"European Union"},
		"Antarctica":     {"Antarctica"},
		"Not Specified":  {"Not Specified"},
	}

	for region, countries := range regionCountries {
		for _, c := range countries {
			if c == country {
				log.Printf("Country %s is in region %s", country, region)
				return region, nil
			}
		}
	}

	log.Printf("Country %s not found in any region", country)
	return "", fmt.Errorf("country not found in any region")
}

func (m *MarketplaceChaincode) ConsensusSubmit(ctx contractapi.TransactionContextInterface, host, port, user, password, dbName string) (string, error) {
	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return "", fmt.Errorf("failed to get client identity: %v", err)
	}

	existingDB, err := ctx.GetStub().GetState("DB_" + clientID)
	if err != nil {
		return "", fmt.Errorf("failed to check existing database for clientID: %v", err)
	}
	if existingDB != nil {
		return "clientID already has a submitted database", nil
	}

	dbDetails := DatabaseDetails{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		DBName:   dbName,
	}

	dbDetailsBytes, err := json.Marshal(dbDetails)
	if err != nil {
		return "", fmt.Errorf("failed to marshal database details: %v", err)
	}

	err = ctx.GetStub().PutState("DB_"+clientID, dbDetailsBytes)
	if err != nil {
		return "", fmt.Errorf("failed to save database details: %v", err)
	}

	return "database successfully submitted", nil
}

func (m *MarketplaceChaincode) ConsensusQueryRow(ctx contractapi.TransactionContextInterface, sqlQuery string, args ...interface{}) (string, error) {
	log.Println("Starting ConsensusQueryRow function")

	resultFrequency := make(map[string]int)
	resultValues := make(map[string]string)

	log.Println("Querying the default database")
	var defaultResult string
	defaultRow := db.QueryRow(sqlQuery, args...)
	err := defaultRow.Scan(&defaultResult)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No rows found in default database.")
		} else {
			log.Printf("Default database query failed: %v", err)
		}
	} else {
		resultFrequency[defaultResult]++
		resultValues[defaultResult] = defaultResult
		log.Printf("Default database returned result: %s", defaultResult)
	}

	log.Println("Querying user-submitted databases")
	resultsIterator, err := ctx.GetStub().GetStateByRange("DB_", "DB_~")
	if err != nil {
		return "", fmt.Errorf("failed to get database submissions: %v", err)
	}
	defer resultsIterator.Close()

	if !resultsIterator.HasNext() {
		log.Println("No user-submitted databases found.")
	}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			log.Printf("Failed to get next result from user database: %v", err)
			continue
		}

		var dbDetails DatabaseDetails
		err = json.Unmarshal(queryResponse.Value, &dbDetails)
		if err != nil {
			log.Printf("Failed to unmarshal database details: %v", err)
			continue
		}

		log.Printf("Querying user database: %s", dbDetails.DBName)
		userDB, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbDetails.User, dbDetails.Password, dbDetails.Host, dbDetails.Port, dbDetails.DBName))
		if err != nil {
			log.Printf("Failed to connect to user database %s: %v", dbDetails.DBName, err)
			continue
		}
		defer userDB.Close()

		var userResult string
		userRow := userDB.QueryRow(sqlQuery, args...)
		err = userRow.Scan(&userResult)
		if err == sql.ErrNoRows {
			log.Printf("No rows found in user database: %s", dbDetails.DBName)
			continue
		} else if err != nil {
			log.Printf("User database query failed for %s: %v", dbDetails.DBName, err)
			continue
		}

		resultFrequency[userResult]++
		resultValues[userResult] = userResult
		log.Printf("User database %s returned result: %s", dbDetails.DBName, userResult)
	}

	var consensusResult string
	maxFrequency := 0
	for result, frequency := range resultFrequency {
		if frequency > maxFrequency {
			consensusResult = result
			maxFrequency = frequency
		}
	}

	if consensusResult == "" {
		log.Println("No consensus result found.")
		return "", fmt.Errorf("no results found in any database")
	}

	log.Printf("Consensus result: %s with frequency: %d", consensusResult, maxFrequency)

	return resultValues[consensusResult], nil
}

func (m *MarketplaceChaincode) CreateToken(ctx contractapi.TransactionContextInterface, tokenID string) (string, error) {
	log.Println("Starting CreateToken function")
	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		log.Printf("Failed to get client identity: %v", err)
		return "", fmt.Errorf("failed to get client identity: %v", err)
	}
	log.Printf("Client ID: %s", clientID)

	tokenBytes, err := ctx.GetStub().GetState(tokenID)
	if err != nil {
		log.Printf("Failed to read from world state: %v", err)
		return "", fmt.Errorf("failed to read from world state: %v", err)
	}
	if tokenBytes != nil {
		log.Printf("Token %s already exists", tokenID)
		return "token already exists", nil
	}

	token := &PaymentToken{
		TokenID: tokenID,
		Creator: clientID,
		Balance: 0,
	}
	tokenBytes, err = json.Marshal(token)
	if err != nil {
		log.Printf("Failed to marshal token: %v", err)
		return "", fmt.Errorf("failed to marshal token: %v", err)
	}

	err = ctx.GetStub().PutState(tokenID, tokenBytes)
	if err != nil {
		log.Printf("Failed to put token to world state: %v", err)
		return "", fmt.Errorf("failed to put token to world state: %v", err)
	}

	log.Println("CreateToken function completed successfully")
	return "token successfully created", nil
}

func (m *MarketplaceChaincode) MintTokens(ctx contractapi.TransactionContextInterface, tokenID string, amount int) (string, error) {
	log.Println("Starting MintTokens function")
	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		log.Printf("Failed to get client identity: %v", err)
		return "", fmt.Errorf("failed to get client identity: %v", err)
	}
	log.Printf("Client ID: %s", clientID)

	tokenBytes, err := ctx.GetStub().GetState(tokenID)
	if err != nil {
		log.Printf("Failed to read from world state: %v", err)
		return "", fmt.Errorf("failed to read from world state: %v", err)
	}
	if tokenBytes == nil {
		log.Printf("No token found for TokenID: %s", tokenID)
		return "token does not exist", nil
	}

	var token PaymentToken
	err = json.Unmarshal(tokenBytes, &token)
	if err != nil {
		log.Printf("Failed to unmarshal token: %v", err)
		return "", fmt.Errorf("failed to unmarshal token: %v", err)
	}

	if token.Creator != clientID {
		log.Printf("Caller is not the creator of token %s", tokenID)
		return "caller is not the creator of the token", nil
	}

	token.Balance += amount

	tokenBytes, err = json.Marshal(token)
	if err != nil {
		log.Printf("Failed to marshal token: %v", err)
		return "", fmt.Errorf("failed to marshal token: %v", err)
	}

	err = ctx.GetStub().PutState(tokenID, tokenBytes)
	if err != nil {
		log.Printf("Failed to put token to world state: %v", err)
		return "", fmt.Errorf("failed to put token to world state: %v", err)
	}

	log.Printf("MintTokens function completed successfully for TokenID: %s", tokenID)
	return "tokens successfully minted", nil
}

func (m *MarketplaceChaincode) GetBalance(ctx contractapi.TransactionContextInterface, clientID string, tokenID string) (int, error) {
	log.Println("Starting GetBalance function")

	balanceKey := clientID + "_" + tokenID

	balanceBytes, err := ctx.GetStub().GetState(balanceKey)
	if err != nil {
		log.Printf("Failed to read from world state: %v", err)
		return 0, fmt.Errorf("failed to read from world state: %v", err)
	}
	if balanceBytes == nil {
		log.Printf("No balance found for ClientID: %s and TokenID: %s", clientID, tokenID)
		return 0, nil
	}

	var balance int
	err = json.Unmarshal(balanceBytes, &balance)
	if err != nil {
		log.Printf("Failed to unmarshal balance: %v", err)
		return 0, fmt.Errorf("failed to unmarshal balance: %v", err)
	}

	log.Printf("GetBalance function completed successfully for TokenID: %s and ClientID: %s", tokenID, clientID)
	return balance, nil
}

func (m *MarketplaceChaincode) TransferTokens(ctx contractapi.TransactionContextInterface, tokenID, newOwner string, amount int) (string, error) {
	log.Println("Starting TransferTokens function")
	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		log.Printf("Failed to get client identity: %v", err)
		return "", fmt.Errorf("failed to get client identity: %v", err)
	}
	log.Printf("Client ID: %s", clientID)

	senderBalanceKey := clientID + "_" + tokenID

	senderBalanceBytes, err := ctx.GetStub().GetState(senderBalanceKey)
	if err != nil {
		log.Printf("Failed to read sender's balance from world state: %v", err)
		return "", fmt.Errorf("failed to read sender's balance from world state: %v", err)
	}
	if senderBalanceBytes == nil {
		log.Printf("No balance found for ClientID: %s and TokenID: %s", clientID, tokenID)
		return "insufficient balance", nil
	}

	var senderBalance int
	err = json.Unmarshal(senderBalanceBytes, &senderBalance)
	if err != nil {
		log.Printf("Failed to unmarshal sender's balance: %v", err)
		return "", fmt.Errorf("failed to unmarshal sender's balance: %v", err)
	}

	if senderBalance < amount {
		log.Printf("Insufficient balance for token %s", tokenID)
		return "insufficient balance", nil
	}

	senderBalance -= amount
	senderBalanceBytes, err = json.Marshal(senderBalance)
	if err != nil {
		log.Printf("Failed to marshal sender's balance: %v", err)
		return "", fmt.Errorf("failed to marshal sender's balance: %v", err)
	}

	err = ctx.GetStub().PutState(senderBalanceKey, senderBalanceBytes)
	if err != nil {
		log.Printf("Failed to update sender's balance in world state: %v", err)
		return "", fmt.Errorf("failed to update sender's balance in world state: %v", err)
	}

	receiverBalanceKey := newOwner + "_" + tokenID

	receiverBalanceBytes, err := ctx.GetStub().GetState(receiverBalanceKey)
	if err != nil {
		log.Printf("Failed to read new owner's balance from world state: %v", err)
		return "", fmt.Errorf("failed to read new owner's balance from world state: %v", err)
	}

	var receiverBalance int
	if receiverBalanceBytes != nil {
		err = json.Unmarshal(receiverBalanceBytes, &receiverBalance)
		if err != nil {
			log.Printf("Failed to unmarshal new owner's balance: %v", err)
			return "", fmt.Errorf("failed to unmarshal new owner's balance: %v", err)
		}
	}

	receiverBalance += amount
	receiverBalanceBytes, err = json.Marshal(receiverBalance)
	if err != nil {
		log.Printf("Failed to marshal new owner's balance: %v", err)
		return "", fmt.Errorf("failed to marshal new owner's balance: %v", err)
	}

	err = ctx.GetStub().PutState(receiverBalanceKey, receiverBalanceBytes)
	if err != nil {
		log.Printf("Failed to update new owner's balance in world state: %v", err)
		return "", fmt.Errorf("failed to update new owner's balance in world state: %v", err)
	}

	log.Printf("TransferTokens function completed successfully. Token %s transferred to %s", tokenID, newOwner)
	return "tokens successfully transferred", nil
}

func (m *MarketplaceChaincode) TransferMintedTokens(ctx contractapi.TransactionContextInterface, tokenID string, amount int) (string, error) {
	log.Println("Starting TransferMintedTokens function")
	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		log.Printf("Failed to get client identity: %v", err)
		return "", fmt.Errorf("failed to get client identity: %v", err)
	}
	log.Printf("Client ID: %s", clientID)

	tokenBytes, err := ctx.GetStub().GetState(tokenID)
	if err != nil {
		log.Printf("Failed to read from world state: %v", err)
		return "", fmt.Errorf("failed to read from world state: %v", err)
	}
	if tokenBytes == nil {
		log.Printf("No token found for TokenID: %s", tokenID)
		return "token does not exist", nil
	}

	var token PaymentToken
	err = json.Unmarshal(tokenBytes, &token)
	if err != nil {
		log.Printf("Failed to unmarshal token: %v", err)
		return "", fmt.Errorf("failed to unmarshal token: %v", err)
	}

	if token.Creator != clientID {
		log.Printf("Caller is not the creator of token %s", tokenID)
		return "caller is not the creator of the token", nil
	}

	if token.Balance < amount {
		log.Printf("Insufficient balance for token %s", tokenID)
		return "insufficient balance", nil
	}

	token.Balance -= amount

	tokenBytes, err = json.Marshal(token)
	if err != nil {
		log.Printf("Failed to marshal token: %v", err)
		return "", fmt.Errorf("failed to marshal token: %v", err)
	}

	err = ctx.GetStub().PutState(tokenID, tokenBytes)
	if err != nil {
		log.Printf("Failed to put token to world state: %v", err)
		return "", fmt.Errorf("failed to put token to world state: %v", err)
	}

	balanceKey := clientID + "_" + tokenID
	balanceBytes, err := ctx.GetStub().GetState(balanceKey)
	if err != nil {
		log.Printf("Failed to read creator's balance from world state: %v", err)
		return "", fmt.Errorf("failed to read creator's balance from world state: %v", err)
	}

	var balance int
	if balanceBytes != nil {
		err = json.Unmarshal(balanceBytes, &balance)
		if err != nil {
			log.Printf("Failed to unmarshal creator's balance: %v", err)
			return "", fmt.Errorf("failed to unmarshal creator's balance: %v", err)
		}
	}

	balance += amount
	balanceBytes, err = json.Marshal(balance)
	if err != nil {
		log.Printf("Failed to marshal creator's balance: %v", err)
		return "", fmt.Errorf("failed to marshal creator's balance: %v", err)
	}

	err = ctx.GetStub().PutState(balanceKey, balanceBytes)
	if err != nil {
		log.Printf("Failed to update creator's balance in world state: %v", err)
		return "", fmt.Errorf("failed to update creator's balance in world state: %v", err)
	}

	log.Printf("TransferMintedTokens function completed successfully for TokenID: %s", tokenID)
	return "minted tokens successfully transferred to creator's account", nil
}

func main() {
	marketplaceChaincode := new(MarketplaceChaincode)

	err := initializeDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	chaincode, err := contractapi.NewChaincode(marketplaceChaincode)
	if err != nil {
		log.Panicf("Error creating chaincode: %v", err)
	}

	if err := chaincode.Start(); err != nil {
		log.Panicf("Error starting chaincode: %v", err)
	}
}
