# Tenant Service

## Installation

## APIs:

1. Create Tenant
- generate human readable tenant name following
  ```
  tenant_20251031_142530_7b1c
  tenant_20251031_142530_f9a2
  tenant_20251031_142531_102c
  ```
- Create a job which will clone the schema for each service via `kafka`

2. Delete Tenant
- Mark tenant as deleted
- Create a job which will delete all the data for for each service via `kafka` (optional)
