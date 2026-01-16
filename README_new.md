
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

### Authors

- [@chibx](https://www.github.com/chibx)




### Tech Stack

**Frontend:** Vue 3, Pinia, Shadcn-Vue
**Backend:** Fiber (Golang), PostgreSQL, Redis


### Run Locally

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
