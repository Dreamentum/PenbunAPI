# SQL Add-on v3 — `id_status` Business Status Column

> **Status:** ✅ EXECUTED — Applied to PENBUN MSSQL Server on 2026-06-10  
> All ALTER statements below have been run. Go models and controllers updated accordingly.

## Purpose

`id_status` replaces the original intent of `is_active` on tables that require **business workflow state tracking**. This is distinct from `is_active` (a technical boolean flag for record enable/disable).

| Column | Type | Purpose |
|--------|------|---------|
| `is_active` | `BIT` default `1` | Technical: is this record enabled/disabled |
| `id_status` | `NVARCHAR(20)` default `'ACTIVE'` | Business: current workflow state |

## Tables Requiring `id_status`

### Transaction Tables (need state tracking)

| Table | Statuses | Notes |
|-------|----------|-------|
| `tb_receive` | `PENDING`, `RECEIVED`, `CANCELLED` | Receive Order header |
| `tb_receive_item` | `PENDING`, `RECEIVED`, `CANCELLED` | Receive Order line items |
| `tb_order` | `DRAFT`, `CONFIRMED`, `PROCESSING`, `SHIPPED`, `DELIVERED`, `CANCELLED` | Sales Order header |
| `tb_order_item` | `PENDING`, `ALLOCATED`, `SHIPPED`, `CANCELLED` | Sales Order lines |
| `tb_deliver` | `PENDING`, `DELIVERED`, `CANCELLED` | Delivery note |
| `tb_return` | `PENDING`, `APPROVED`, `REJECTED`, `COMPLETED` | Return/Refund |
| `tb_invoice` | `DRAFT`, `SENT`, `PAID`, `OVERDUE`, `CANCELLED` | Invoice/AR |

### Master Data Tables (richer states beyond active/inactive)

| Table | Statuses | Notes |
|-------|----------|-------|
| `tb_product` | `ACTIVE`, `DISCONTINUED`, `OUT_OF_STOCK` | Product lifecycle |
| `tb_customer` | `ACTIVE`, `BLOCKED`, `VIP` | Customer relationship state |
| `tb_vendor` | `ACTIVE`, `BLACKLISTED` | Vendor relationship state |
| `tb_discount` | `ACTIVE`, `EXPIRED`, `SCHEDULED` | Promotion lifecycle |

### Tables That Keep `is_active` Only

Simple reference/lookup tables with no workflow states:

- `tb_product_category`
- `tb_product_format_type`
- `tb_product_group`
- `tb_product_sku`
- `tb_unit_type`
- `tb_customer_type`
- `tb_vendor_type`
- `tb_warehouse`
- `tb_company`
- `tb_reference`
- `tb_discount_type`
- `tb_users`

## Executed SQL (2026-06-10)

```sql
-- Added to 4 master data tables for business workflow state tracking
ALTER TABLE tb_product  ADD id_status NVARCHAR(20) NOT NULL DEFAULT ('ACTIVE');
ALTER TABLE tb_customer ADD id_status NVARCHAR(20) NOT NULL DEFAULT ('ACTIVE');
ALTER TABLE tb_vendor   ADD id_status NVARCHAR(20) NOT NULL DEFAULT ('ACTIVE');
ALTER TABLE tb_discount ADD id_status NVARCHAR(20) NOT NULL DEFAULT ('ACTIVE');
```

## SQL Template (for future tables)

```sql
-- Add id_status to a table
ALTER TABLE tb_<table> ADD id_status NVARCHAR(20) NOT NULL DEFAULT ('ACTIVE');

-- Optional CHECK constraint
ALTER TABLE tb_<table> ADD CONSTRAINT CK_tb_<table>_id_status
CHECK (id_status IN ('ACTIVE', 'CANCELLED', 'PENDING' /* ... */));

-- Index for status-based queries
CREATE NONCLUSTERED INDEX IX_tb_<table>_id_status ON tb_<table> (id_status)
INCLUDE (autoID);
```

## is_active Migration (executed 2026-06-10)

The following tables had `id_status` renamed to `is_active` to fix a column naming inconsistency:

```sql
EXEC sp_rename 'tb_customer_type.id_status', 'is_active', 'COLUMN';
EXEC sp_rename 'tb_reference.id_status',     'is_active', 'COLUMN';
EXEC sp_rename 'tb_users.id_status',         'is_active', 'COLUMN';
```

These 3 tables are lookup/reference tables that only need a technical `BIT` flag — they do **not** need `id_status` as a business column.

## Migration for Already-Migrated Tables

The following tables had `id_status` renamed to `is_active` in the v3.0.0 migration. If they now need `id_status` as a business column, add it back as a **new column** (don't rename `is_active` back):

```sql
-- tb_customer_type, tb_reference, tb_users: keep is_active, add id_status only if needed
-- Currently these tables only need is_active (lookup tables)
```

## Go Implementation (executed 2026-06-10)

### Model Update

```go
// Add to struct after IsActive
IsActive bool   `json:"is_active"`
IDStatus string `json:"id_status"`    // business status: ACTIVE, CANCELLED, EXPIRED, etc.
```

### Controller SELECT Pattern

```go
// Add id_status between is_active and update_by
rows, err := config.DB.Query(`SELECT ..., is_active, id_status, update_by, ... FROM tb_<entity> WHERE is_delete = 0`)
// Scan:
rows.Scan(&item.IsActive, &item.IDStatus, &item.UpdateBy, ...)
```

### Controller INSERT Pattern

```go
// Insert with default 'ACTIVE' if empty string supplied
INSERT INTO tb_<entity> (..., id_status, update_by)
VALUES (..., COALESCE(NULLIF(?, ''), 'ACTIVE'), ?)
// Args: ..., item.IDStatus, item.UpdateBy
```

### Controller UPDATE Pattern

```go
// Partial update — empty string keeps current value
UPDATE tb_<entity> SET ..., id_status = COALESCE(NULLIF(?, ''), id_status), update_by = ?
// Args: ..., item.IDStatus, item.UpdateBy, id
```

### Filter by status

```sql
WHERE is_delete = 0 AND id_status = 'ACTIVE'         -- active records
WHERE is_delete = 0 AND id_status IN ('ACTIVE', 'PENDING')  -- multi-status filter
```
