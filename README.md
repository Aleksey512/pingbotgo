# Telegram Ping Monitor Bot

<img width="75%" alt="image" src="https://github.com/user-attachments/assets/c6376c39-52e5-4465-a3ef-fbf602b40da7" />

## Overview

This is a Telegram bot that monitors server availability by performing ping checks and sends notifications to subscribed users. The bot allows users to subscribe/unsubscribe to server status updates and provides real-time monitoring of specified servers.

## Features

- **Server Monitoring**: Regularly pings servers listed in the configuration
- **Status Notifications**: Sends Telegram alerts when server status changes
- **Subscription Management**: Users can subscribe/unsubscribe to notifications
- **Configuration**: Easy setup via `.env` file

## Prerequisites

- Go 1.24.2
- Telegram Bot Token (obtain from [@BotFather](https://t.me/BotFather))

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/Aleksey512/pingbotgo
   cd pingbotgo
   ```

2. Install Go dependencies:
   ```bash
   go mod download
   ```

3. Create a `.env` file in root of project:
   ```ini
    # Sqlite storage
    SQLLITE_PATH=sqlite3.db
    
    # Redis configuration if replace Storage class in internal/app
    # REDIS_PASSWORD=your_redis_password # Optional
    # REDIS_USERNAME=your_redis_username # Optional
    # HOST_REDIS=redis
    # PORT_REDIS=6379
    
    
    # Telegram bot configuration
    TELEGRAM_BOT_TOKEN=your_telegram_bot_token_here
    
    # Server list (override in code or extend here)
    # Format: {"SERVER_NAME": "IP"}
    # Example:
    # SERVERS={"ЯНДЕКС":"ya.ru","GOOGLE":"google.com"}
    SERVERS={}
   ```

## Usage

1. Start the bot:
   ```bash
   go run cmd/app/main.go
   ```

2. Interact with the bot in Telegram:
   - `/start` - Show welcome message
   - `/subscribe` - Subscribe to server status notifications
   - `/unsubscribe` - Unsubscribe from notifications
   - `/ping_now` - Check current server statuses
   - `/config` - Check current servers to ping

## Bot Commands

| Command       | Description                          |
|---------------|--------------------------------------|
| `/start`      | Show welcome message and help        |
| `/subscribe`  | Subscribe to server status updates   |
| `/unsubscribe`| Stop receiving notifications         |
| `/ping_now`   | Check current server statuses        |
| `/config`     | Check current servers to ping        |


## Architecture

1. **Ping Service**: Regularly checks server availability
2. **Notification Service**: Sends alerts to subscribed users
3. **Telegram Bot Handler**: Processes user commands and interactions
4. **Redis Storage**: Manages user subscriptions and status history

## License

[MIT License](./LICENSE)

## Support

For issues or feature requests, please open an issue on GitHub.
