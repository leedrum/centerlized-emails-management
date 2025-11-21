
# Centralized Email Management - Microservices Architecture Diagram

```mermaid
flowchart LR

    subgraph User
        U[End User]
    end

    subgraph Auth[Auth Service]
        A1[OAuth2 / SAML Login]
        A2[Manage Tenant Credentials]
    end

    subgraph Fetcher[Mail Fetcher Service]
        F1[Fetch Emails from Gmail/Outlook/EWS]
        F2[Save EML to S3/MinIO]
        F3[Publish email_fetched to Kafka]
    end

    subgraph Processor[Email Processor Service]
        P1[Download EML from S3]
        P2[Parse EML (subject, body, attachments)]
        P3[Save Email → Postgres (tenant schema)]
        P4[Publish email_processed]
    end

    subgraph Indexer[Email Indexer Service]
        I1[Consume email_processed]
        I2[Load Email from Postgres]
        I3[Index Email → Elasticsearch]
        I4[DLQ for failed messages]
    end

    subgraph Infra[Infrastructure]
        K[(Kafka)]
        PG[(PostgreSQL: tenant schemas)]
        ES[(Elasticsearch)]
        S3[(S3 / MinIO)]
    end

    U -->|Login| Auth
    Auth --> Fetcher

    Fetcher -->|email_fetched| K --> Processor
    Processor -->|email_processed| K --> Indexer

    Fetcher --> S3
    Processor --> PG
    Indexer --> ES

```
