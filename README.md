# HudAutomata

RFID-based automation and balance management system with self-service capabilities.

## Features

- 🔐 JWT Authentication & Role-based Access Control
- 👥 User Management with RFID Card Integration
- 💰 Balance Management (Credit/Debit/Refund)
- 🤖 Self-Service Automation API for IoT Devices
- 📊 Dashboard with Statistics & Charts
- 📝 Comprehensive System Logging
- 🔄 Real-time Transaction Processing

## Tech Stack

### Backend
- **Go 1.23** with Gin Web Framework
- **GORM** ORM (PostgreSQL/SQLite)
- **JWT** Authentication
- **Bcrypt** Password Hashing

### Frontend
- **Preact** with Bun runtime
- **DaisyUI** + Tailwind CSS
- **Chart.js** for Analytics
- **Caddy** Web Server

## Quick Start

### Development (SQLite)

```bash
# Backend
cp .env.example .env
go mod download
go run src/main.go

# Frontend
cd frontend
bun install
bun run dev
```

### Production (Docker)

```bash
# Set environment variables
export DB_PASSWORD=your-secure-password
export JWT_SECRET=your-jwt-secret

# Start all services
docker-compose up -d
```

## Default Admin Credentials

**Username:** `admin`  
**Password:** `admin123`

⚠️ **Change these credentials immediately after first login!**

## API Documentation

### Authentication
- `POST /api/v1/auth/login` - Admin login
- `POST /api/v1/auth/logout` - Admin logout
- `GET /api/v1/auth/me` - Get current admin

### Automation (IoT)
- `POST /api/v1/automation/scan` - RFID scan & service
- `POST /api/v1/automation/check-balance` - Check balance only

### Users
- `GET /api/v1/users` - List users
- `POST /api/v1/users` - Create user
- `GET /api/v1/users/:id` - Get user
- `PUT /api/v1/users/:id` - Update user
- `DELETE /api/v1/users/:id` - Delete user

### Transactions
- `GET /api/v1/transactions` - List transactions
- `POST /api/v1/transactions` - Create transaction
- `GET /api/v1/transactions/:id` - Get transaction

### Dashboard
- `GET /api/v1/dashboard/stats` - Get statistics
- `GET /api/v1/dashboard/charts` - Get chart data
- `GET /api/v1/dashboard/recent` - Recent activities

## IoT Integration Example

```cpp
// Arduino/ESP32 Example
#include <HTTPClient.h>

String rfidId = "ABC123456";
HTTPClient http;

http.begin("http://backend:8080/api/v1/automation/scan");
http.addHeader("Content-Type", "application/json");

String payload = "{\"rfid_card_id\":\"" + rfidId + "\",\"service_cost\":10.5,\"description\":\"Service\"}";
int httpCode = http.POST(payload);

if (httpCode == 200) {
  // Service approved - activate relay
  digitalWrite(RELAY_PIN, HIGH);
} else {
  // Insufficient balance or error
  Serial.println("Service denied");
}
```

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `HOST` | `0.0.0.0` | Server host |
| `PORT` | `8080` | Server port |
| `DB_DRIVER` | `sqlite` | Database driver (sqlite/postgres) |
| `DB_HOST` | `localhost` | PostgreSQL host |
| `DB_PORT` | `5432` | PostgreSQL port |
| `DB_USER` | `huduser` | Database user |
| `DB_PASSWORD` | - | Database password |
| `DB_NAME` | `hudautomata` | Database name |
| `JWT_SECRET` | - | JWT secret key |
| `CORS_ORIGINS` | `*` | Allowed CORS origins |

## License

MIT License - see LICENSE file for details

## Contributing

Pull requests are welcome! For major changes, please open an issue first.

---

**Version:** 1.0.0  
**Status:** ✅ Production Ready  
**Last Updated:** 2025-10-02  
**Documentation:** See [PROD.md](PROD.md) for complete technical documentation


