- [1. REST API](#1-rest-api)
  - [1.1. Product API](#11-product-api)
    - [1.1.1. GET /api/product](#111-get-apiproduct)
    - [1.1.2. POST /api/product](#112-post-apiproduct)
    - [1.1.3. GET /api/product/:id](#113-get-apiproductid)
    - [1.1.4. PATCH /api/product/:id](#114-patch-apiproductid)
    - [1.1.5. DELETE /api/product/:id](#115-delete-apiproductid)
  - [1.2. Order API](#12-order-api)
    - [1.2.1. GET /api/order](#121-get-apiorder)
    - [1.2.2. POST /api/order](#122-post-apiorder)
    - [1.2.3. GET /api/order/:id](#123-get-apiorderid)
    - [1.2.4. PATCH /api/order/:id](#124-patch-apiorderid)
    - [1.2.5. POST /api/order/:id/cancel](#125-post-apiorderidcancel)
    - [1.2.6. DELETE /api/order/:id](#126-delete-apiorderid)
  - [1.3. Account API](#13-account-api)
    - [1.3.1. POST /api/account/login](#131-post-apiaccountlogin)
    - [1.3.2. POST /api/account/register](#132-post-apiaccountregister)
    - [1.3.3. GET /api/account/:id](#133-get-apiaccountid)
    - [1.3.4. PATCH /api/account/:id](#134-patch-apiaccountid)
    - [1.3.5. GET /api/account/:id/orders](#135-get-apiaccountidorders)
    - [1.3.6. POST /api/account/:id/password](#136-post-apiaccountidpassword)
    - [1.3.7. DELETE /api/account/:id](#137-delete-apiaccountid)
- [2. Models](#2-models)
  - [2.1. Product](#21-product)
    - [2.1.1. ProductStatus](#211-productstatus)
  - [2.2. Order](#22-order)
    - [2.2.1. OrderStatus](#221-orderstatus)
    - [2.2.2. ProductOrder](#222-productorder)
  - [2.3. Account](#23-account)
    - [2.3.1 AccountType](#231-accounttype)
  - [2.4. Session](#24-session)
  - [2.5. Audit Entry](#25-audit-entry)
    - [2.5.1. EntryType](#251-entrytype)
    - [2.5.2. ResourceType](#252-resourcetype)

Specification for the showcase backend server.

# 1. REST API

## 1.1. Product API

### 1.1.1. GET /api/product

List all products

Query params:
- limit - Limit the number of items per page
- page - Page number
  
Returns:
- `200 OK` - JSON-encoded list of products
  
### 1.1.2. POST /api/product

Create a new product

> **Warning**
> Requires Authorization as Staff

> **Note**
> This endpoint, on success, will create an audit entry

Accepts: JSON-encoded product without the ID field

Returns:
- `203 No content` - Product was added successfully
- `400 Bad Request` - Malformed body or it contains invlid data
- `401 Unauthorized` - No `Authorization` header or it has invalid data

### 1.1.3. GET /api/product/:id

Get product information for a specific product

Returns:
- `200 OK` - JSON-encoded product
- `404 Not Found` - Product with given ID does not exists

### 1.1.4. PATCH /api/product/:id

Modify product information of a specific product

> **Warning**
> Requires Authorization as Staff

> **Note**
> This endpoint, on success, will create an audit entry

Accepts: Full or partial JSON-encoded product without the ID field

Returns:
- `200 OK` - JSON-encoded product with modifications applied
- `400 Bad Request` - Malformed body or it contains invalid data
- `401 Unauthorized` - No `Authorization` header or it has invalid data
- `404 Not Found` - Product with given ID does not exists

### 1.1.5. DELETE /api/product/:id

Remove a sepcific product

> **Warning**
> Requires Authorization as Staff

> **Note**
> This endpoint, on success, will create an audit entry

Returns:
- `203 No Content` - Product was successfully deleted
- `401 Unauthorized` - No `Authorization` header or it has invalid data
- `404 Not Found` - Product with given ID does not exists

## 1.2. Order API

### 1.2.1. GET /api/order

Get order list sorted from newest to oldest

> **Warning**
> Requires Authorization as Staff

Query params:
- limit - Limit the number of items per page
- page - Page number

Returns:
- `200 OK` - JSON-encoded list of orders
- `401 Unauthorized` - No `Authorization` header or it has invalid data

### 1.2.2. POST /api/order

Create a new order

> **Warning**
> Requires Authorization

> **Note**
> This endpoint, on success, will create an audit entry

Accepts: Partial JSON-encoded order model

Returns:
- `200 OK` - JSON-encoded order confirmation
- `400 Bad Request` - Malformed body or it contains invalid data
- `401 Unauthorized` - No `Authorization` header or it has invalid data

### 1.2.3. GET /api/order/:id

Get order information for a specific order

> **Warning**
> Requires Authorization.
> Only Staff and order owner can access this information

Returns:
- `200 OK` - JSON-encoded order information
- `401 Unauthorized` - No `Authorization` header or it has invalid data
- `404 Not Found` - Order with given ID does not exists

### 1.2.4. PATCH /api/order/:id

Modify order information of a specific order

> **Warning**
> Requires Authorization as Staff

> **Note**
> This endpoint, on success, will create an audit entry

Accepts: Full or partial JSON-encoded order model withour the ID field

Returns:
- `200 OK` - JSON-encoded order with modifications applied
- `400 Bad Request` - Malformed body or it contains invalid data
- `401 Unauthorized` - No `Authorization` header or it has invalid data
- `404 Not Found` - Order with given ID does not exists

### 1.2.5. POST /api/order/:id/cancel

Cancel order

> **Warning**
> Requires Authorization.
> Only Staff and order owner can access this information

> **Note**
> This endpoint, on success, will create an audit entry

Accepts: JSON-encoded cancelation reason

Returns:
- `200 No content` - Succesfully canceled order
- `400 Bad Request` - Malformed body or it contains invalid data
- `401 Unauthorized` - No `Authorization` header or it has invalid data
- `404 Not Found` - Order with given ID does not exists

### 1.2.6. DELETE /api/order/:id

Remove a sepcific order

> **Warning**
> Requires Authorization as Staff

> **Note**
> This endpoint, on success, will create an audit entry

Returns:
- `203 No Content` - Order was successfully deleted
- `401 Unauthorized` - No `Authorization` header or it has invalid data
- `404 Not Found` - Order with given ID does not exists

## 1.3. Account API

### 1.3.1. POST /api/account/login

Get account access token

> **Note**
> This endpoint, on success, will create an audit entry

Accepts: JSON-encoded account e-mail and hashed password

Returns:
- `200 OK` - Account login confirmation with the account's access token
- `400 Bad Request` - Malformed body or it contains invalid data
- `401 Unauthorized` - Password hash does not match
- `404 Not Found` - Account with the given e-mail does not exists

### 1.3.2. POST /api/account/register

Create a new account

> **Note**
> This endpoint, on success, will create an audit entry

Accepts: JSON-encoded account model without the ID field

Returns:
- `200 OK` - Account creation confirmation with the account's access token
- `400 Bad Request` - Malformed body or it contains invalid data

### 1.3.3. GET /api/account/:id

Get information about a specific account.
Use `@me` as the id to get informations about the current account

> **Warning**
> Requires Authorization.
> Only Staff can see other account information

Returns:
- `200 OK` - JSON-encoded account information
- `401 Unauthorized` - No `Authorization` header or it has invalid data
- `404 Not Found` - Account with given ID does not exists

### 1.3.4. PATCH /api/account/:id

Modify accout information

> **Warning**
> Requires Authorization.
> Only Staff can modify other account information

> **Note**
> This endpoint, on success, will create an audit entry

Accepts: Full or partial JSON-encoded account model withour the ID field

Returns:
- `200 OK` - JSON-encoded account with modifications applied
- `400 Bad Request` - Malformed body or it contains invalid data
- `401 Unauthorized` - No `Authorization` header or it has invalid data
- `404 Not Found` - Account with given ID does not exists

### 1.3.5. GET /api/account/:id/orders

Get a list of orders made by the given account

> **Warning**
> Requires Authorization.
> Only Staff can see other account orders

Query params:
- limit - Limit the number of items per page
- page - Page number

Returns:
- `200 OK` - JSON-encoded list of order information
- `401 Unauthorized` - No `Authorization` header or it has invalid data
- `404 Not Found` - Account with given ID does not exists

### 1.3.6. POST /api/account/:id/password

Change password for the given account

> **Warning**
> Requires Authorization.
> Only Staff can change other account password without the knowlage of the original password

> **Note**
> This endpoint, on success, will create an audit entry

Accepts: JSON-encoded old and new password hashes

Returns:
- `200 OK` - Password change confirmation with a new access token
- `400 Bad Request` - Malformed body or it contains invalid data
- `401 Unauthorized` - No `Authorization` header or it has invalid data
- `404 Not Found` - Account with given ID does not exists

### 1.3.7. DELETE /api/account/:id

Remove account with the given id

> **Warning**
> Requires Authorization.
> Only Staff can remove other accounts

> **Note**
> This endpoint, on success, will create an audit entry

Returns:
- `203 No Content` - Successfully removed account
- `401 Unauthorized` - No `Authorization` header or it has invalid data
- `404 Not Found` - Account with given ID does not exists

# 2. Models

In this document, field names are in PascalCase which isn't a requirement for field names in code

Models field names should be consistent with the language nameing convention

SQL table column names and JSON keys should use **snake_case** instead

All models for SQL tables also include the following fields:

| Field     | Type      | Description
|-----------|-----------|-------------
| ID        | uint      | Primary key
| CreatedAt | datetime  | When the row was created
| UpdatedAt | datetime  | WHen the row as modified
| DeletedAt | datetime? | when the row was soft deleted

## 2.1. Product

| Field       | Type          | Description
|-------------|---------------|-------------
| Status      | [ProductStatus](#211-productstatus) | Product status. In stock, out of stock, discontinued
| Name        | string        | The name of the product
| Description | string        | The entire description of the product in markdown
| Price       | float         | Product's price in EUR
| Images      | array(string) | Product images. Array of valid URLs. Can be represented and stored as a semicolon separated string

A SQL table named `products` is based around this model

### 2.1.1. ProductStatus

ProductStatus is an enum with the following values:

| Name         | Value
|--------------|-------
| InStock      | 1
| OutOfStock   | 2
| Discontinued | 3

## 2.2. Order

| Field           | Type                | Description
|-----------------|---------------------|-------------
| Status          | [OrderStatus](#221-orderstatus) | Order status. Created, paid, in realisation, shipped, completed, canceled
| OrderBy         | uint                | ID of the account who created this order
| ShippingAddress | string              | Shipping address
| InvoceAddress   | string              | Invoice addresss. If empty, use shipping address
| Products        | array([ProductOrder](#222-productorder)) | List of ordered products with the amount and price per piece. Can be represented and stored as a semicolon separated string
| ShippingPrice   | float               | Postage and handling costs
| Total           | float               | Grand total for the order

A SQL table named `orders` is based around this model

### 2.2.1. OrderStatus

OrderStatus is an enum with the following values:

| Name          | Value
|---------------|-------
| Created       | 1
| Paid          | 2
| InRealisation | 3
| Shipped       | 4
| Completed     | 5
| Canceled      | 6

### 2.2.2. ProductOrder

| Field     | Type  | Description
|-----------|-------|-------------
| ProductID | uint  | The ID of the ordered product
| Amount    | uint  | The amount ordered
| Price     | float | Price at the time of the order

Structure can be represented and stored as a comma separated string

## 2.3. Account

| Field   | Type        | Description
|---------|-------------|-------------
| Type    | [AccountType](#231-accounttype) | The type of the account. Staff, user
| Email   | string      | E-Mail address identifing the account
| Hash    | string      | Argon2id hashed password
| Salt    | string      | Salt for the password hash 
| Name    | string      | User's first name
| Surname | string      | User's last name

A SQL table named `accounts` is based around this model

### 2.3.1 AccountType

AccountType is an enum with the following values:

| Name  | Value
|-------|-------
| Staff | 1
| User  | 2

## 2.4. Session

| Field     | Type   | Description
|-----------|--------|-------------
| AccountID | uint   | The account that the session token is assigned to
| Token     | string | The session token

A SQL table named `sessions` is based around this model

Session will be valid only for 30 minutes

## 2.5. Audit Entry

| Field        | Type         | Description
|--------------|--------------|-------------
| Type         | [EntryType](#251-entrytype) | The type of the audit entry. Creation, modification, deletion
| ResourceType | [ResourceType](#252-resourcetype) | The resource type which expirienced change
| ResourceID   | uint         | The ID of the resource
| ExecutorID   | uint         | The ID of the account that made the change

A SQL table named `audit_log` is based around this model

### 2.5.1. EntryType

EntryType is an enum with the following values:

| Name         | Value
|--------------|-------
| Creation     | 1
| Modification | 2
| Deletion     | 3

### 2.5.2. ResourceType

ResourceType is an enum with the following values:

| Name            | Value
|-----------------|-------
| Product         | 1
| Order           | 2
| Account         | 3
| AccountPassword | 4
| Session         | 5

> **Note**
> `AccountPassword` type is used to differenciate account information modification from password chnage
