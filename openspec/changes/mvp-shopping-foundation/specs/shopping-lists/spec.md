# Shopping Lists Specification

## Purpose

Allow users to create shopping lists, add products from the catalog, and calculate a closed total based on the latest known prices.

## Requirements

### Requirement: Create List

The system SHALL allow authenticated users to create a named shopping list.

#### Scenario: Create shopping list

- GIVEN an authenticated user
- WHEN `POST /api/v1/lists` is called with `{ "name": "Compra semanal" }`
- THEN a list SHALL be created with status `active` and returned with id

### Requirement: Add Item

The system SHALL allow adding a product to a list with a quantity. If no custom price is provided, the latest price SHALL be used.

#### Scenario: Add product with latest price

- GIVEN an active list and a product with latest price 4.50 €
- WHEN `POST /api/v1/lists/{id}/items` is called with `{ "product_id": "...", "quantity": 2 }`
- THEN the item SHALL be stored with unit price 4.50 € and total 9.00 €

#### Scenario: Add product with custom price

- GIVEN an active list and a product
- WHEN the user adds the product with `{ "product_id": "...", "quantity": 1, "custom_price": 5.00 }`
- THEN the item SHALL use 5.00 € as unit price and total 5.00 €

### Requirement: Calculate Closed Total

The system SHALL return the sum of all item totals as the list's `estimated_total`.

#### Scenario: Calculate total for populated list

- GIVEN a list has two items with totals 9.00 € and 5.00 €
- WHEN `GET /api/v1/lists/{id}` is called
- THEN `estimated_total` SHALL be 14.00 €

### Requirement: Manage Items

The system SHALL allow marking items as purchased and removing items.

#### Scenario: Mark item purchased

- GIVEN a list has an unpurchased item
- WHEN `PATCH /api/v1/lists/{id}/items/{item_id}` is called with `{ "purchased": true }`
- THEN the item SHALL be marked purchased

#### Scenario: Remove item

- GIVEN a list has an item
- WHEN `DELETE /api/v1/lists/{id}/items/{item_id}` is called
- THEN the item SHALL be removed and the estimated total recalculated
