# Investment Portfolio Rebalancer

This project is a simple Go application to help you rebalance your investment portfolio based on your target asset class allocations and a new contribution amount. It reads your current assets and asset classes from CSV files, then suggests how to allocate a new investment to maintain your desired balance.

## Features

- Reads asset classes and assets from CSV files.
- Calculates current portfolio allocation.
- Suggests how to allocate a new contribution to maintain target percentages.
- Allocation suggestions are based on asset scores within each class.

## Getting Started

### Prerequisites

- Go 1.18 or newer installed on your system.

### CSV File Format

You need two CSV files: one for asset classes and one for assets.

#### Example: `classes.csv`

| id | name         | target_pct |
|----|--------------|------------|
| 1  | Fixed Income | 60         |
| 2  | Stocks       | 30         |
| 3  | Crypto       | 10         |

**File content:**
#### Example: `assets.csv`

| id | class_id | name         | score | quantity | unit_price |
|----|----------|--------------|-------|----------|------------|
| 1  | 1        | CDB Banco X  | 8     | 10000    | 1          |
| 2  | 1        | Tesouro Selic| 9     | 5000     | 1          |
| 3  | 2        | PETR4        | 7     | 200      | 30         |
| 4  | 2        | ITUB4        | 8     | 150      | 25         |
| 5  | 3        | Bitcoin      | 10    | 0.1      | 250000     |
| 6  | 3        | Ethereum     | 9     | 0.5      | 15000      |

**File content:**
