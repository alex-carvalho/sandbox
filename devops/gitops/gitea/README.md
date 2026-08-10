# Gitea Proof of Concept (PoC)

This project contains a Proof of Concept (PoC) setup for [Gitea](https://gitea.com/) using [PostgreSQL 18](https://www.postgresql.org/) as the database, including support for Gitea Actions (using [act_runner](https://gitea.com/gitea/act_runner)).

## Prerequisites

- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/)

---

## Method 1: Running with Docker Compose (Recommended)

### Step 1: Start Gitea & PostgreSQL
To start the Gitea and PostgreSQL 18 stack, execute the following command:

```bash
docker-compose up -d
```

### Step 2: Complete Gitea First-time Installation
1. Open your browser and navigate to **[http://localhost:3000](http://localhost:3000)**.
2. Review the settings (database connection is pre-configured).
3. Scroll to the bottom and click **"Install Gitea"**.
4. Register the first user account (this account automatically becomes the Gitea administrator).

### Step 3: Get the Actions Runner Registration Token
1. Log in to Gitea with your admin account.
2. Go to **Site Administration** (wrench icon in the top right) **> Actions > Runners**.
3. Click **"Create new Runner"**.
4. Copy the registration token shown.

### Step 4: Start the Actions Runner
To spin up the runner with the registration token:

```bash
GITEA_RUNNER_REGISTRATION_TOKEN="your_token_here" docker-compose --profile runner up -d
```

To stop all services (including the runner):

```bash
docker-compose --profile runner down
```
