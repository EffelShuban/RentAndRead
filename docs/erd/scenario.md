# Scenario End to End

## Register

```bash
curl --location 'http://localhost:8080/api/users/register' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name": "Jane Doe",
    "email": "jane@example.com",
    "password": "securepassword"
}'
```

## Login

```bash
curl --location 'http://localhost:8080/api/users/login' \
--header 'Content-Type: application/json' \
--data-raw '{
    "email": "jane@example.com",
    "password": "securepassword"
}'
```

## Top Up

```bash
curl --location 'http://localhost:8080/api/wallet/topup' \
--header 'Authorization: Bearer your_jwt_token_here' \
--header 'Content-Type: application/json' \
--data '{
    "amount": 50000
}'

curl --location 'http://localhost:8080/webhooks/xendit' \
--header 'Content-Type: application/json' \
--data '{
    "external_id": "TOPUP-USERID-TIMESTAMP",
    "status": "PAID",
    "amount": 50000
}'
```

## Rent a Book

```bash
curl --location 'http://localhost:8080/api/rentals' \
--header 'Authorization: Bearer your_jwt_token_here' \
--header 'Content-Type: application/json' \
--data '{
    "book_id": 1,
    "days": 5
}'
```

## View Rented Books

```bash
curl --location 'http://localhost:8080/api/users/books' \
--header 'Authorization: Bearer your_jwt_token_here'
```
