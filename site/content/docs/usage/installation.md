---
title: "Installation & Setup"
weight: 1
---

# Installation and Setup

Getting Lynx FIM up and running involves three main steps: building the binary, setting up your environment, and initializing the configuration.

## 1. Build from Source

Since I built Lynx in Go, you can compile it into a single, static binary. 

### Prerequisites
- Go 1.22 or higher installed on your system.

### Steps
1. Clone the repository:
   ```bash
   git clone https://github.com/ItsAdam01/Lynx.git
   cd Lynx
   ```
2. Build the binary using the provided Makefile:
   ```bash
   make build
   ```
The executable will be located in the `bin/` directory.

## 2. Setting up the Alert Webhook (Discord)

Lynx supports sending real-time alerts to Discord. Here is how I set mine up for testing:

1. Open Discord and go to your **Server Settings**.
2. Navigate to **Integrations** > **Webhooks**.
3. Click **New Webhook**.
4. Give it a name (like "Lynx Security Bot") and select the channel where alerts should appear.
5. Click **Copy Webhook URL**. You will need this for your configuration.

## 3. Environment and Secrets

Lynx uses an environment variable to store the HMAC secret key. This key is used to sign your baseline file so it can't be tampered with.

I recommend using a `.env` file for local development (though you should never commit it!).

### Sample .env File
```bash
# The secret key used for baseline integrity
LYNX_HMAC_SECRET="your-super-long-random-secret-key"
```

To load it before running Lynx:
```bash
export LYNX_HMAC_SECRET="your-super-long-random-secret-key"
```

## 4. Initialize Lynx

Run the init command to generate a default `config.yaml`:
```bash
./bin/lynx init
```

Open `config.yaml` and paste your Discord Webhook URL into the `webhook_url` field. You can also customize which directories Lynx should watch.
