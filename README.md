# 🧩 System Overview of Centerlized Email Management --- Email Aggregation Microservices Architecture

## 🏗️ Main Architecture

-   **API Gateway** --- Single entry point, routes requests to the
    correct service.
-   **Kafka** --- Central event bus for asynchronous communication
    between services.
-   **S3 / MinIO** --- Stores `.eml` files, attachments, and images.
-   **MongoDB** --- Stores full email metadata and body content.
-   **ElasticSearch** --- Supports full-text search (title, sender,
    receiver, body text).
-   **PostgreSQL** --- Stores user, tenant, and configuration metadata.

------------------------------------------------------------------------

## 🧱 Core Microservices

### 1️⃣ Auth Service

**Responsibilities:** - Manages users and authentication (JWT /
OAuth2). - Integrates with **Google / Microsoft / OKTA (SAML /
OIDC)**. - Multi-tenant support. - Issues access tokens for API Gateway.

**Database:** PostgreSQL\
**Kafka Topics:** - `user.created`, `tenant.created`, `user.deleted`

------------------------------------------------------------------------

### 2️⃣ Tenant Service

**Responsibilities:** - Manages tenant (company) information, domains,
and mail server configurations. - Allows tenant admins to manage users
and mailbox connections (OAuth tokens).

**Database:** PostgreSQL\
**Kafka Topics:** - `tenant.config.updated` (notifies Mail Fetcher
Service)

------------------------------------------------------------------------

### 3️⃣ Mail Fetcher Service

**Responsibilities:** - Connects to external mail servers (EWS, Gmail
API, Outlook) using tenant config. - Downloads `.eml` files (via polling
or webhook). - Uploads `.eml` files to S3. - Publishes Kafka event:
`email.eml.uploaded`.

**Kafka Topics:** - `email.eml.uploaded`

------------------------------------------------------------------------

### 4️⃣ Email Processor Service

**Responsibilities:** - Subscribes to `email.eml.uploaded` topic. -
Downloads `.eml` file from S3. - Parses and extracts metadata (title,
sender, receiver, datetime, body, attachments, images). - Uploads
decoded images/attachments to S3. - Saves metadata and body to
MongoDB. - Publishes Kafka event: `email.index.requested`.

**Database:** MongoDB\
**Kafka Topics:** - Input: `email.eml.uploaded` - Output:
`email.index.requested`

------------------------------------------------------------------------

### 5️⃣ Email Indexer Service

**Responsibilities:** - Listens to `email.index.requested` events. -
Indexes emails (title, sender, receiver, datetime, body text) into
ElasticSearch.

**Database:** ElasticSearch\
**Kafka Topics:** - Input: `email.index.requested`

------------------------------------------------------------------------

### 6️⃣ Email API Service

**Responsibilities:** - Provides REST APIs for users: - Fetch emails
(per inbox / all inboxes) - View email details - Send new emails -
Search emails (via ElasticSearch) - Reads data from MongoDB and
ElasticSearch.

**Database:** MongoDB, ElasticSearch\
**Kafka Topics (optional):** - `email.sent` (for tracking or audit)

------------------------------------------------------------------------

### 7️⃣ Notification Service (optional)

**Responsibilities:** - Sends notifications when new emails arrive (via
WebSocket or email). - Subscribes to `email.eml.uploaded` or
`email.processed`.

------------------------------------------------------------------------

### 8️⃣ Audit / Activity Service (optional)

**Responsibilities:** - Logs user actions: login, read, send, delete. -
Supports compliance tracking (GDPR, SOC2). - Subscribes to major Kafka
topics.

------------------------------------------------------------------------

## 🔄 Workflow

``` plaintext
1. User logs in → Auth Service (OAuth2/SAML)
2. Tenant configures mail → Tenant Service
3. Mail Fetcher → downloads `.eml` → uploads to S3 → emits Kafka: email.eml.uploaded
4. Email Processor → parses `.eml` → saves Mongo → emits Kafka: email.index.requested
5. Email Indexer → stores searchable data in ElasticSearch
6. Email API → serves list/search/view/send endpoints
```

------------------------------------------------------------------------

## 🧠 Design Notes

-   **Isolation per tenant:** Use `tenant_id` as namespace in Kafka,
    MongoDB, ElasticSearch.
-   **Idempotency:** Mail Fetcher and Processor must handle duplicate
    processing safely.
-   **Retry & DLQ:** Configure Kafka Dead-letter Queue for resilience.
-   **Scalability:**
    -   Mail Fetcher & Processor scale per tenant or mailbox.
    -   Email API scales horizontally (read-heavy).

------------------------------------------------------------------------
