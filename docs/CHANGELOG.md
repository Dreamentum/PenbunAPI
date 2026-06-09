# Changelog

All notable changes to PenbunAPI will be documented in this file.

## [3.0.0] - 2026-06-10

### Added
- Complete rewrite in Go using Fiber v2 web framework
- 17 master data modules with 8 standard CRUD functions each (136 endpoints):
  - Product API (Layer 4) — hybrid product system with `count_stock`, `is_active`, `id_status`, FK references to product_group, product_format_type, unit_type, vendor
  - Business ID generation (`generateBusinessID` in Go) matching `USP_GENERATE_BUSINESS_ID` algorithm (prefix + series char + 6-digit running number)
  - `OUTPUT INSERTED.autoID` pattern for identity retrieval on INSERT
- 16 master data modules with 8 standard CRUD functions each (128 endpoints):
  - Publisher, Publisher Type, Customer, Customer Type
  - Vendor, Vendor Type, Book, Book Type
  - Discount, Discount Type, Unit Type, Product Format Type
  - Product Category, Product Group, Warehouse
- `id_status` business status column — `NVARCHAR(20)` with `DEFAULT ('ACTIVE')` on `tb_product`, `tb_customer`, `tb_vendor`, `tb_discount` for workflow state tracking (`ACTIVE`, `DISCONTINUED`, `BLOCKED`, `EXPIRED`, etc.)
- `docs/sql-addon-v3.md` — design document for `id_status` business status column strategy
- JWT authentication with Bearer token validation
- Token blacklist system for logout (in-memory, thread-safe)
- Transaction management with panic recovery and step logging
- Partial update pattern using `COALESCE(NULLIF(...))`
- Soft delete (`is_delete = 1`) with `update_by` tracking
- Graceful shutdown on SIGINT/SIGTERM
- Connection pool configuration (MaxOpenConns=25, MaxIdleConns=10)
- CORS and Recover middleware
- Environment config loading from `.env` with defaults
- Transaction audit logging to `logs/transaction.log`
- Structured project layout: config, models, controllers, routes, middleware, utils
- Parameterized SQL queries (SQL injection prevention)
- 109 unit tests with race detection across all packages
- Response helper utilities for standardized JSON format (`{status, message, data}`)

### Changed
- Migrated from v1.x codebase to v3.0.0 Go backend
- All API routes under `/api/v1/protected/[module]`
- Response format standardized across all endpoints
- `id_status` → `is_active` rename on `tb_customer_type`, `tb_reference`, `tb_users` (all 3 unchanged since they only need `BIT` flag)
- All `IDStatus` Go model fields renamed to `IsActive` (bool), all controller SQL updated to `is_active`
- New `IDStatus` (string) field added to Product, Customer, Vendor, Discount models and controllers for business status tracking

### Removed
- Legacy v1.x Node.js or previous backend code

### Pending
- Remove Book API / Book Type API (superseded by Product API)
- Global error handler middleware
- Receive Note / Receive Item APIs (Layer 5)
- Order / Order Item APIs (Layer 6)
- Global error handler middleware
- Role-based access control (RBAC)
- Performance optimizations (Prefork, Redis cache, Connection Pool tuning)
- Integration tests against live database

## [2.1.0] - 2025-12-16

### Added
- Development standards documentation (standard.md)
- Layered architecture planning (todov2.md)

## [1.7.4] - Previous

### Added
- Initial API modules and planning (todov1.md)
