# Bezier

Welcome!

Bezier is a *simple* pre-paid card api that allows creating cards and merchants to simulate a transaction workflow. The details of the API can be seen in the sections. To keep things simple, following design choices have been made:

* There is no authentication on the api
* Multiple cards and merchants can be creating by hitting the endpoint
* Transactions can be viewed on the merchant details endpoint
* The data is stored in memory, however the API depends on a simple storage API. So integrating with any DB should be a matter of writing new adapters

Have fun buying coffees!

## Create card
Create endpoint creates a new card with a different card number
```
$ curl -X GET https://bezier.herokuapp.com/v1/card/create
```
Produces:
```json
{
    "card_number": 4921000000000001,
    "error": null
}
```

## Get Card details
Details can be viewed using card number in parameters:
```
$ curl -X GET 'https://bezier.herokuapp.com/v1/card/details?card_number=4921000000000001'
```
Fetches:
```json
{
  "card_details": {
    "card_number": 4921000000000001,
    "available_balance": 3000,
    "blocked_balance": 0,
    "total_loaded": 3000
  },
  "error": ""
}
```

## Topup card
Topup requires a POST:
```
$ curl -X POST \
  https://bezier.herokuapp.com/v1/card/top-up \
  -H 'content-type: application/json' \
  -d '{"card_number": 4921000000000001,"amount": 3000}'
```

## Create Merchant
Merchants can be created similar to cards. Merchant ID returned will be their identifier to perform any transaction related operations

```
$ curl -X GET https://bezier.herokuapp.com/v1/merchant/create
```
```json
{
  "merchant_id": "M1001",
  "error": null
}
```

## Authorize Transaction
It is assumed the user hands merchant the card details (card_number) to make any payment
An authorization request can be made as follows:
```
$ curl -X POST \
  https://bezier.herokuapp.com/v1/merchant/authorize-transaction \
  -H 'content-type: application/json' \
  -d '{"card_number": 4921000000000001,"merchant_id":"M1001","amount": 300}'
```
This returns a transaction id that can be used for further actions:
```json
{"transaction_id":"TX10001","error":""}
```

Merchant's transactions should reflect the authorizations, captures, refunds and reversals:
```
$ curl -X GET https://bezier.herokuapp.com/v1/merchant/transactions?merchant_id=M1001
```
```json
{
  "merchant_activity": {
    "id": "M1001",
    "transactions": [
      {
        "id": "TX10001",
        "card_number": 4921000000000001,
        "authorized": 300,
        "captured": 0,
        "reversed": 0,
        "refunded": 0
      }
    ]
  },
  "error": ""
}
```

And card should reflect the blocked balance:
```
curl -X GET 'https://bezier.herokuapp.com/v1/card/details?card_number=4921000000000001'
```

```json
{
  "card_details": {
    "card_number": 4921000000000001,
    "available_balance": 2700,
    "blocked_balance": 300,
    "total_loaded": 3000
  },
  "error": ""
}
```


## Capture transaction
Capturing moves funds from authorized to captured field in the transaction removing from card's blocked funds
```
$ curl -X POST \
  https://bezier.herokuapp.com/v1/merchant/capture-transaction \
  -H 'content-type: application/json' \
  -d '{"merchant_id":"M1001","transaction_id":"TX10001","amount": 100}'
```
Updates Merchant details to:
```json
{
  "merchant_activity": {
    "id": "M1001",
    "transactions": [
      {
        "id": "TX10001",
        "card_number": 4921000000000001,
        "authorized": 200,
        "captured": 100,
        "reversed": 0,
        "refunded": 0
      }
    ]
  },
  "error": ""
}
```
And unblocks funds on card details:
```json
{
  "card_details": {
    "card_number": 4921000000000001,
    "available_balance": 2700,
    "blocked_balance": 200, // Blocked funds now withdrawn!
    "total_loaded": 3000
  },
  "error": ""
}
```

## Reverse transaction

Reverse moves funds from authorized to reversed field in the transaction and also unblocks the amount on the card
```
$ curl -X POST \
  https://bezier.herokuapp.com/v1/merchant/reverse-transaction \
  -H 'content-type: application/json' \
  -d '{"merchant_id":"M1001","transaction_id":"TX10001","amount": 100}'
```
Updates Merchant details to:
```json
{
  "merchant_activity": {
    "id": "M1001",
    "transactions": [
      {
        "id": "TX10001",
        "card_number": 4921000000000001,
        "authorized": 100,
        "captured": 100,
        "reversed": 100,
        "refunded": 0
      }
    ]
  },
  "error": ""
}
```

And moves funds back to available_balance to:
```json
{
  "card_details": {
    "card_number": 4921000000000001,
    "available_balance": 2800,
    "blocked_balance": 100, // Less money blocked!
    "total_loaded": 3000
  },
  "error": ""
}
```

## Refund captured 

```
$ curl -X POST \
  https://bezier.herokuapp.com/v1/merchant/refund-transaction \
  -H 'content-type: application/json' \
  -d '{"merchant_id":"M1001","transaction_id":"TX10001","amount": 50}'
```

Moves funds to on transaction to refunded:
```json
{
  "merchant_activity": {
    "id": "M1001",
    "transactions": [
      {
        "id": "TX10001",
        "card_number": 4921000000000001,
        "authorized": 100,
        "captured": 50,
        "reversed": 100,
        "refunded": 50
      }
    ]
  },
  "error": ""
}
```

And the card can be used for more coffees!:
```json
{
  "card_details": {
    "card_number": 4921000000000001,
    "available_balance": 2850, // more coffees!
    "blocked_balance": 100,
    "total_loaded": 3000
  },
  "error": ""
}
```

## Development and testing

Clone the repo and run `dep ensure` to get the dependencies

Then to run tests, go to the root of the project and run `go test -v ./...`

To start the server locally, run: `go run main.go`