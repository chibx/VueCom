# VueCom

A deployable set of packages that help developers set up any e-commerce application.

> [!WARNING]
> This project is currently in development.

## Badges

[![MIT License](https://img.shields.io/badge/License-MIT-green.svg)](https://choosealicense.com/licenses/mit/)

## Overview

VueCom is an e-commerce platform designed to provide a comprehensive and modular solution for building scalable online stores. It aims to offer parity with established platforms like Magento and Shopify while leveraging modern technologies and a microservices architecture.

## Features

This project aims to implement a rich set of features, categorized as follows:

### Core E-commerce Functionality

*   [ ] **Catalog Management:** Product listings, categorization, attributes, variants (with extensive SKU generation), bundling, digital products, advanced search and filtering, product reviews and ratings.
*   [ ] **Inventory Management:** Real-time stock tracking, multi-warehouse support, low-stock alerts, backorders, supplier integration, automated reordering, and IoT-enabled omnichannel syncing.
*   [ ] **Order Management:** Processing, fulfillment, tracking, returns/refunds, invoicing, multi-channel synchronization, and advanced order status workflows.
*   [ ] **Payment Processing:** Integration with multiple gateways (credit cards, PayPal, PayStack), secure checkout, fraud detection, tax calculations, and Level 1 PCI compliance.
*   [ ] **Shipping & Logistics:** Carrier integrations, real-time rate calculations, label printing, dropshipping, international shipping with dynamic tax/shipping zones, and omnichannel fulfillment.
*   [ ] **Customer Management:** Accounts, profiles, segmentation, loyalty programs, guest checkout, B2B company accounts with credit limits, wish lists, and saved carts.
*   [ ] **Marketing & SEO:** Promo codes, discount management, email marketing, SEO tools (custom URLs, meta tags, sitemaps), abandoned cart recovery, analytics integration, and A/B testing.
*   [ ] **Storefront & Admin:** Customizable themes, drag-and-drop page builder, mobile-responsive design, multilingual/multi-currency support, SSR for SEO, headless commerce (GraphQL/REST APIs), intuitive admin dashboard, role-based access control, and reporting dashboards.

### B2B Specific Features

*   [ ] Custom pricing per customer/group
*   [ ] Bulk ordering interface
*   [ ] Quote management
*   [ ] Advanced approval workflows
*   [ ] Punchout catalogs (OCI/cXML)

### Performance & Security

*   [ ] Caching strategies and CDN integration
*   [ ] Fast loading times (low TTFB)
*   [ ] Auto-scaling infrastructure for high traffic
*   [ ] SSL and data encryption
*   [ ] PCI compliance (Level 1)
*   [ ] GDPR tools
*   [ ] CAPTCHA integration
*   [ ] Audit logging
*   [ ] Web application firewall

### Integrations

*   [ ] POS system integration
*   [ ] CRM/ERP system integration (SAP, Salesforce)
*   [ ] Social media commerce (Instagram, TikTok)
*   [ ] Third-party app marketplace
*   [ ] Unified multi-channel selling

### Analytics & Reporting

*   [ ] Built-in dashboards (sales, customer behavior, inventory trends)
*   [ ] A/B testing reporting
*   [ ] Data import/export (CSV/Excel bulk uploads)

### Global Expansion

*   [ ] Multi-store setup
*   [ ] Localization (language and regional formats)
*   [ ] Cross-border compliance tools
*   [ ] Multi-storefront management with automated currency conversion

<!-- ### Advanced & Differentiating Features

*   **AI & ML Capabilities:** AI-Powered Visual Search, Predictive Inventory AI, AI Content Generation, AI Insights Dashboard, Dynamic pricing automation, and modular recommendation services.
*   **Sustainability & Ethics:** Carbon footprint calculators, eco-friendly product filtering, supplier ethics/sustainability scoring, and eco-badges.
*   **Search & Customization:** Advanced Fitment Search, Bulk Variant Editing, and automated local pickup customization.
*   **Marketplace & Vendor Support:** Native multi-vendor marketplace functionality, commission tracking, and vendor management.
*   **Immersive & Emerging Tech:** AR/VR product try-on, Blockchain payment, and decentralized hosting.
*   **Subscription & Business Models:** Advanced subscription models, rental models, and hybrid ownership/rental options.
*   **Privacy & Security:** Privacy-focused data handling, anonymized analytics, zero-knowledge proofs, and encrypted customer data.
*   **Social & Omnichannel:** Built-in social commerce hub, user-generated content moderation, influencer collaboration tools, and microservices for real-time sync with physical stores/IoT. -->

## Technical Implementation

### Architecture

*   [ ] Microservices architecture
*   [ ] Headless-first design (GraphQL/REST APIs)
*   [ ] Composable architecture (MACH principles)
*   [ ] Multi-tenancy support (isolated tenants via UUIDs, row-level security)
*   [ ] Separate control plane and data plane
*   [ ] Hybrid SaaS & Self-Hosted model (Open Core)

### Development & Operations

*   [ ] Docker containerization
*   [ ] Kubernetes orchestration
*   [ ] CI/CD pipeline
*   [ ] Comprehensive API documentation
*   [ ] Webhooks support with admin UI
*   [ ] Developer experience tools and documentation
*   [ ] "Eject to Self-Hosted" option for SaaS users
*   [ ] Pre-integrated tech stacks (Algolia, Akeneo)

### Extensibility

*   [ ] Plugin/extension architecture
*   [ ] Support for dynamic loading (gRPC/HTTP hooks, WASM)
*   [ ] Plugins as separate microservices
*   [ ] API for extensions

### Deployment & Scaling

*   [ ] Auto-scaling infrastructure
*   [ ] Blue-green deployment support
*   [ ] Multi-region deployment capabilities

### Recommendation Service

*   [ ] Modular recommendation service ("Lite" for self-hosted, "Deep" for SaaS)
*   [ ] Decoupled via message queue (e.g., RabbitMQ)
*   [ ] Upsell AI/ML features to self-hosted users via API key

### Implementation Status

The project is currently in the **Development Started** phase.

## Tech Stack

**Frontend:** Vue 3, Pinia, Shadcn-Vue
**Backend:** Fiber (Golang), PostgreSQL, Redis, RabbitMQ

## Requirements

*   **Docker**
*   **Docker Compose**
*   **Node.js >= v20.10.0**
*   **Pnpm**
*   **Golang >= v1.22.1**
*   **PostgreSQL**
*   **Redis**
*   **RabbitMQ**
*   **Prometheus _(Optional)_**

## Database Structure (Modular Monolith with PostgreSQL)

The project utilizes a modular monolith architecture with PostgreSQL as the primary database. Each core service (Catalog, Inventory, Orders, Users) is designed to manage its own database schema, promoting separation of concerns and independent evolution. The SQL scripts located in `deploy-config/sql/` illustrate this modular approach:

*   `00-add-user.sh`: Script for adding a dedicated `vuecom` user to PostgreSQL.
*   `00-init.sql`: Initial database setup script.
*   `01-catalog-00.sql`: SQL script for the **Catalog** service database schema.
*   `01-catalog-01.sql`: Additional SQL script for the **Catalog** service.
*   `02-inventory-00.sql`: SQL script for the **Inventory** service database schema.
*   `03-orders-00.sql`: SQL script for the **Orders** service database schema.
*   `04-users-00.sql`: SQL script for the **Users** service database schema.
*   `04-users-01.sql`: Additional SQL script for the **Users** service.
*   `grant.sql`: SQL script for granting necessary permissions to the `vuecom` user.

This structure allows for a clear distinction and management of data pertinent to each service, even within a single PostgreSQL instance.

## Prerequisites

*   You create a password for the user **`vuecom`** in the PostgreSQL database by assigning it to the `APP_PG_PASSWORD` environment variable in the `./backend/services/gateway/.env` file.
*   For the Redis instance, you'll need to add a connection password in a newly created file, not going to be committed to the repository, in the `./deploy-config/redis/secrets.conf` file using the format in the [`secrets.conf.example`](./deploy-config/redis/secrets.conf.example) file.
*   For as many actively used/initialized services, you set up the respective environment variables `./backend/services/gateway/.env` file as stated in the [.env.example](./backend/services/gateway/.env.example) file.

**Note:** You can spin up the necessary services (e.g., PostgreSQL, Redis, RabbitMQ) by running the following command in the project root:

```bash
  docker-compose up service1 service2 ... serviceN
```

## Run Locally

*This is if you have the necessary dependencies installed and set up already.*

> [!WARNING]
> The backend development section requires that you have Docker and Docker Compose installed and set up with images of PostgreSQL, Redis, RabbitMQ, and Prometheus (optional).

**Clone the project**

```bash
  git clone https://github.com/chibx/VueCom
```

**Go to the project directory**

```bash
  cd VueCom
```

**Install dependencies (Frontend)**

```bash
  pnpm install
```

**Install dependencies (Backend)**

```bash
  pnpm install:server
```

**Start the development server (client & server)**

```bash
  pnpm dev
```

## Contributing

For Golang development, we recommend using `air` for hot reloading to improve your development experience. You can find it at [https://github.com/air-verse/air](https://github.com/air-verse/air).

## Authors

- [@chibx](https://www.github.com/chibx)
