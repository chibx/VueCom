
# VueCom

A deployable set of packages that help developer set up any e-commerce application.

> [!WARNING]
> This project is currently in development.

<!-- ## Badges -->

[![MIT License](https://img.shields.io/badge/License-MIT-green.svg)](https://choosealicense.com/licenses/mit/)

### Features

- [ ] **Product Catalog**
- [X] **User Authentication**
- [ ] **Order Management**
- [ ] **Payment Processing**
- [ ] **Admin Panel**

### Requirements

- **Docker**
- **Docker Compose**
- **Node.js >= v20.10.0**
- **Pnpm**
- **Golang >= v1.22.1**
- **PostgreSQL**
- **Redis**
- **RabbitMQ**
- **Prometheus**

### Authors

- [@chibx](https://www.github.com/chibx)




### Tech Stack

**Frontend:** Vue 3, Pinia, Shadcn-Vue
**Backend:** Fiber (Golang), PostgreSQL, Redis, RabbitMQ


### Prerequisites

- You create a password for the user **`vuecom`** in the PostgreSQL database.
- For the Redis instance, you'll neeed to add a connection password in a newly created file, not going to be committed to the repository, in the `./deploy-config/redis/secrets.conf` file using the format in the [`secrets.conf.example`](./deploy-config/redis/secrets.conf.example) file.
- For as many actively used/initialized service, you set up the respective environment variables `./backend/services/gateway_service/.env` file as stated in the [.env.example](./backend/services/gateway_service/.env.example) file.


### Run Locally

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
  pnpm run install:server
```

**Start the development server (client & server)**

```bash
  pnpm run dev
```
