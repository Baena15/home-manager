# Products Specification

## Purpose

Maintain a household product catalog with price history so shopping lists can estimate costs.

## Requirements

### Requirement: Create Product

The system SHALL allow authenticated users to create a product with name, default quantity unit, and optional category.

#### Scenario: Create a new product

- GIVEN an authenticated user
- WHEN `POST /api/v1/products` is called with `{ "name": "Carne picada", "unit": "g", "category": "carnicería" }`
- THEN the product SHALL be stored and returned with a generated id
- AND the response status SHALL be 201

#### Scenario: Reject duplicate product name

- GIVEN a product named "Carne picada" already exists
- WHEN the user tries to create another product with the same name
- THEN the response SHALL return 409 with code `DUPLICATE_PRODUCT`

### Requirement: List Products

The system SHALL allow authenticated users to list products with pagination and optional search by name.

#### Scenario: List products

- GIVEN products exist in the catalog
- WHEN `GET /api/v1/products?limit=20&offset=0` is called
- THEN the response SHALL return a list of products and pagination metadata

### Requirement: Price History

The system SHALL store a price record per product, store name, amount, and timestamp. The latest price SHALL be retrievable.

#### Scenario: Add price to product

- GIVEN an existing product
- WHEN `POST /api/v1/products/{id}/prices` is called with `{ "store": "Mercadona", "amount": 4.50 }`
- THEN the price SHALL be stored and the product's latest price SHALL be updated

#### Scenario: Get product with latest price

- GIVEN a product has at least one price
- WHEN `GET /api/v1/products/{id}` is called
- THEN the response SHALL include the product fields and `latest_price`

#### Scenario: Product without price history

- GIVEN a product has no prices
- WHEN `GET /api/v1/products/{id}` is called
- THEN `latest_price` SHALL be null and the UI SHALL display "sin precio"
