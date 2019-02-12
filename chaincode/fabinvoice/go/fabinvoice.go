package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	// "github.com/hyperledger/fabric/core/chaincode/lib/cid"
	sc "github.com/hyperledger/fabric/protos/peer"
)

type SmartContract struct {
}

// Define the invoice structure, with 11 properties.  Structure tags are used by encoding/json library

type Invoice struct {
	InvoiceNumber   string `json:"invoiceNumber"`
	BilledTo        string `json:"billedTo"`
	InvoiceDate     string `json:"invoiceDate"`
	InvoiceAmount   string `json:"invoiceAmount"`
	ItemDescription string `json:"itemDescription"`
	GR              string `json:"gr"`
	IsPaid          string `json:"isPaid"`
	PaidAmount      string `json:"paidAmount"`
	Repaid          string `json:"repaid"`
	RepaymentAmount string `json:"repaymentAmount"`
}

/*
 * The Init method is called when the Smart Contract "fabinvoice" is instantiated by the blockchain network
 * Best practice is to have any Ledger initialization in separate function -- see initLedger()
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {

	return shim.Success(nil)

}

func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()

	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "queryAllInvoices" {

		return s.queryAllInvoices(APIstub)

	} else if function == "newInvoice" {

		return s.newInvoice(APIstub, args)

	} else if function == "getUser" {

		return s.getUser(APIstub, args)

	} else if function == "createInvoiceWithJsonInput" {

		return s.getUser(APIstub, args)

	} else if function == "isGoodsReceived" {

		return s.isGoodsReceived(APIstub, args)

	} else if function == "isPaidToSupplier" {

		return s.isPaidToSupplier(APIstub, args)

	} else if function == "isPaidToBank" {

		return s.isPaidToBank(APIstub, args)

	} else if function == "getHistoryForInvoice" {

		return s.getHistoryForInvoice(APIstub, args)

	}

	return shim.Error("Invalid Smart Contract function name.")

}

func getQueryResultForQueryString(APIstub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

	resultsIterator, err := APIstub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryRecords
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}

		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	return buffer.Bytes(), nil
}

//display all invoices from Id INV0 to INV999
func (s *SmartContract) queryAllInvoices(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "INV0"
	endKey := "INV999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Invoice\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryAllInvoices:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

//Create Invoice or Raise Invoices
func (s *SmartContract) newInvoice(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 11 {
		return shim.Error("Incorrect number of arguments. Expecting 11")
	}

	var invoice = Invoice{InvoiceNumber: args[1], BilledTo: args[2], InvoiceDate: args[3], InvoiceAmount: args[4], ItemDescription: args[5], GR: args[6], IsPaid: args[7], PaidAmount: args[8], Repaid: args[9], RepaymentAmount: args[10]}

	invoiceAsBytes, _ := json.Marshal(invoice)
	APIstub.PutState(args[0], invoiceAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) getUser(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	return shim.Success(nil)

}

//Create Invoice with Json
func (s *SmartContract) createInvoiceWithJsonInput(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}
	fmt.Println("args[1] > ", args[1])
	invoiceAsBytes := []byte(args[1])
	invoice := Invoice{}
	err := json.Unmarshal(invoiceAsBytes, &invoice)

	if err != nil {
		return shim.Error("Error During Invoice Unmarshall")
	}
	APIstub.PutState(args[0], invoiceAsBytes)
	return shim.Success(nil)
}

//Goods Received
func (s *SmartContract) isGoodsReceived(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	invoiceAsBytes, _ := APIstub.GetState(args[0])
	invoice := Invoice{}

	json.Unmarshal(invoiceAsBytes, &invoice)
	invoice.GR = args[1]

	invoiceAsBytes, _ = json.Marshal(invoice)
	APIstub.PutState(args[0], invoiceAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) isPaidToSupplier(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {

		return shim.Error("Invalid number of arguments, Expecting 2")

	}

	invoiceAsBytes, _ := APIstub.GetState(args[0])
	invoice := Invoice{}

	json.Unmarshal(invoiceAsBytes, &invoice)
	invoice.PaidAmount = args[1]
	paid, _ := strconv.ParseFloat(args[1], 32)
	invoiceAmount, _ := strconv.ParseFloat(invoice.InvoiceAmount, 32)

	if paid >= invoiceAmount {

		return shim.Error("Invoice Amount is Less than Paid Amount")

	}

	invoice.IsPaid = "Y"
	invoiceAsBytes, _ = json.Marshal(invoice)
	APIstub.PutState(args[0], invoiceAsBytes)

	return shim.Success(nil)

}

func (s *SmartContract) isPaidToBank(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {

		return shim.Error("Invalid number of arguments, Expecting 2 ")

	}

	invoiceAsBytes, _ := APIstub.GetState(args[0])
	invoice := Invoice{}

	json.Unmarshal(invoiceAsBytes, &invoice)
	invoice.RepaymentAmount = args[1]

	ripaid, _ := strconv.ParseFloat(args[1], 32)
	invoiceAmount, _ := strconv.ParseFloat(invoice.InvoiceAmount, 32)

	if ripaid < invoiceAmount {

		return shim.Error("Invoice Amount is greater than Paid Amount")

	}

	invoice.Repaid = "Y"
	invoiceAsBytes, _ = json.Marshal(invoice)
	APIstub.PutState(args[0], invoiceAsBytes)

	return shim.Success(nil)

}

func (s *SmartContract) getHistoryForInvoice(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	carKey := args[0]

	resultsIterator, err := APIstub.GetHistoryForKey(carKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing historic values for the car
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")
		buffer.WriteString(string(response.Value))

		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}

	buffer.WriteString("]")

	return shim.Success(buffer.Bytes())

}

func main() {

	err := shim.Start(new(SmartContract))

	if err != nil {

		fmt.Printf("Error creating new Smart Contract: %s", err)

	}
}
