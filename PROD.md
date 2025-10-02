# HudAutomata - RFID Admin Panel Üretim Planı

## 📋 Proje Özeti

**HudAutomata**, RFID kart tabanlı bir self-service otomasyon ve bakiye yönetim sistemidir. Kullanıcılar RFID kartlarını okutarak yeterli bakiyeleri varsa otomatik olarak hizmet alırlar. Yetersiz bakiye durumunda admin'den talep ederler ve admin web paneli üzerinden kullanıcılara bakiye tanımlar.

### 🔄 Sistem Akışı

1. **Kullanıcı** → RFID kartını okutma cihazına okutur
2. **Sistem** → RFID kartı ile kullanıcıyı tanımlar ve bakiyeyi kontrol eder
3. **Yeterli Bakiye Varsa** → Hizmet otomatik verilir, bakiyeden düşüm yapılır
4. **Yetersiz Bakiye** → Kullanıcı admin'e başvurur
5. **Admin** → Web panelinden kullanıcıya bakiye yükler
6. **Kullanıcı** → Tekrar RFID okutup hizmet alır

### 🎯 Temel Özellikler

1. **RFID Kullanıcı Yönetimi**
   - Kullanıcı ekleme/düzenleme/silme
   - RFID kart numarası ile kullanıcı ilişkilendirme
   - Kullanıcı profil bilgileri (isim **zorunlu**, email/telefon **opsiyonel**)

2. **Bakiye Yönetimi**
   - Admin tarafından kullanıcılara bakiye yükleme
   - Otomasyon sistemi tarafından otomatik bakiye düşümü
   - Bakiye geçmişi görüntüleme
   - Toplu bakiye işlemleri
   - Minimum bakiye uyarıları

3. **Self-Service Otomasyon**
   - RFID kart okutma API endpoint'i
   - Bakiye kontrolü ve otomatik hizmet verme
   - Yetersiz bakiye bildirimi
   - İşlem başarı/hata durumu dönüşü
   - Gerçek zamanlı işlem logu

4. **Admin Yönetimi**
   - Admin kullanıcı oluşturma
   - Rol bazlı yetkilendirme (Super Admin, Admin)
   - Admin aktivite logları

5. **Sistem Logları**
   - Tüm işlemlerin detaylı loglanması (RFID okutma, bakiye düşümü, admin işlemleri)
   - Filtreleme ve arama özellikleri
   - Export (CSV/JSON) desteği
   - Başarılı/başarısız işlem ayrımı

6. **Dashboard & Raporlama**
   - Toplam kullanıcı sayısı
   - Toplam bakiye miktarı
   - Günlük/haftalık/aylık işlem grafikleri
   - Son işlemler listesi
   - RFID okutma istatistikleri

---

## 🏗️ Teknik Mimari

### Backend Stack (Go)

**Framework & Libraries:**
```
- Gin Web Framework (HTTP router & middleware)
- GORM (ORM - SQLite/PostgreSQL desteği)
- JWT-Go (Authentication)
- Bcrypt (Password hashing)
- Zap/Logrus (Structured logging)
- Go-Validator (Input validation)
- CORS Middleware
```

**Database Schema:**
```sql
-- users table
CREATE TABLE users (
    id UUID PRIMARY KEY,
    rfid_card_id VARCHAR(255) UNIQUE NOT NULL,  -- ZORUNLU
    name VARCHAR(255) NOT NULL,                  -- ZORUNLU (tek bir isim alanı)
    email VARCHAR(255),                          -- OPSİYONEL
    phone VARCHAR(20),                           -- OPSİYONEL
    balance DECIMAL(10,2) DEFAULT 0.00,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- admins table
CREATE TABLE admins (
    id UUID PRIMARY KEY,
    username VARCHAR(100) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(20) DEFAULT 'admin', -- 'super_admin', 'admin'
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_login TIMESTAMP
);

-- transactions table
CREATE TABLE transactions (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id),
    admin_id UUID REFERENCES admins(id),           -- NULL ise otomasyon sistemi tarafından yapılmış
    type VARCHAR(20) NOT NULL,                     -- 'credit' (admin yükler), 'debit' (sistem düşer), 'refund'
    amount DECIMAL(10,2) NOT NULL,
    balance_before DECIMAL(10,2) NOT NULL,
    balance_after DECIMAL(10,2) NOT NULL,
    description TEXT,
    source VARCHAR(50) DEFAULT 'admin',            -- 'admin', 'automation', 'system'
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- system_logs table
CREATE TABLE system_logs (
    id UUID PRIMARY KEY,
    admin_id UUID REFERENCES admins(id),
    action VARCHAR(100) NOT NULL,
    resource VARCHAR(100),
    resource_id VARCHAR(255),
    details JSONB,
    ip_address VARCHAR(45),
    user_agent TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- sessions table (for JWT blacklist/refresh tokens)
CREATE TABLE sessions (
    id UUID PRIMARY KEY,
    admin_id UUID REFERENCES admins(id),
    token_hash VARCHAR(255) UNIQUE NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- indexes
CREATE INDEX idx_users_rfid ON users(rfid_card_id);
CREATE INDEX idx_users_balance ON users(balance);
CREATE INDEX idx_transactions_user ON transactions(user_id);
CREATE INDEX idx_transactions_created ON transactions(created_at DESC);
CREATE INDEX idx_logs_admin ON system_logs(admin_id);
CREATE INDEX idx_logs_created ON system_logs(created_at DESC);
```

### Frontend Stack (Preact + Bun)

**UI Framework & Libraries:**
```json
{
  "dependencies": {
    "preact": "^10.27.2",
    "preact-router": "^4.1.2",
    "@preact/signals": "^1.3.2",
    
    // Modern UI Kit (seçenekler)
    "daisyui": "^4.12.14",        // Tailwind-based (önerilen)
    // VEYA
    "shadcn-ui-preact": "latest",  // Headless UI components
    
    // State Management
    "zustand": "^5.0.3",
    
    // Form Handling
    "react-hook-form": "^7.54.2",
    "zod": "^3.24.1",
    
    // HTTP Client
    "ky": "^1.7.3",
    
    // Charts & Visualization
    "chart.js": "^4.4.8",
    "preact-chartjs-2": "^1.0.0",
    
    // Icons
    "lucide-preact": "^0.468.0",
    
    // Utilities
    "date-fns": "^4.1.0",
    "clsx": "^2.1.1"
  },
  "devDependencies": {
    "@preact/preset-vite": "^2.10.2",
    "vite": "^7.1.7",
    "tailwindcss": "^3.4.17",
    "autoprefixer": "^10.4.20",
    "postcss": "^8.5.1"
  }
}
```

**UI Kit Seçimi: DaisyUI (Tailwind CSS tabanlı)**
- Modern ve responsive komponentler
- Dark/Light mode desteği
- Preact ile mükemmel uyum
- Küçük bundle size
- Kolay özelleştirme

---

## 🔌 Backend API Endpoints

### Authentication
```
POST   /api/v1/auth/login          # Admin login
POST   /api/v1/auth/logout         # Admin logout
POST   /api/v1/auth/refresh        # Refresh token
GET    /api/v1/auth/me             # Get current admin info
```

### Users (RFID Users)
```
GET    /api/v1/users               # List all users (pagination, search, filter)
POST   /api/v1/users               # Create new user
GET    /api/v1/users/:id           # Get user by ID
PUT    /api/v1/users/:id           # Update user
DELETE /api/v1/users/:id           # Delete user
GET    /api/v1/users/rfid/:cardId  # Get user by RFID card
GET    /api/v1/users/:id/balance   # Get user balance & history
```

### Automation (RFID Self-Service)
```
POST   /api/v1/automation/scan           # RFID kart okutma ve hizmet verme
POST   /api/v1/automation/check-balance  # Sadece bakiye kontrolü (RFID ile)
GET    /api/v1/automation/history        # Otomasyon işlem geçmişi
```

**Scan Endpoint Request/Response:**
```json
// Request
{
  "rfid_card_id": "ABC123456",
  "service_cost": 10.50,  // Hizmet bedeli
  "description": "Çamaşır makinesi kullanımı"
}

// Response - Başarılı
{
  "success": true,
  "user_id": "uuid",
  "user_name": "Ahmet Yılmaz",
  "balance_before": 50.00,
  "balance_after": 39.50,
  "transaction_id": "uuid",
  "message": "Hizmet verildi"
}

// Response - Yetersiz Bakiye
{
  "success": false,
  "user_id": "uuid",
  "user_name": "Ahmet Yılmaz",
  "current_balance": 5.00,
  "required_amount": 10.50,
  "deficit": 5.50,
  "message": "Yetersiz bakiye. Lütfen yöneticiye başvurun."
}

// Response - Kullanıcı Bulunamadı
{
  "success": false,
  "message": "RFID kartı kayıtlı değil"
}
```

### Transactions
```
POST   /api/v1/transactions        # Create new transaction (credit/debit)
GET    /api/v1/transactions        # List all transactions (pagination, filter)
GET    /api/v1/transactions/:id    # Get transaction details
GET    /api/v1/users/:id/transactions  # Get user transactions
```

### Admins
```
GET    /api/v1/admins              # List all admins (super_admin only)
POST   /api/v1/admins              # Create new admin (super_admin only)
GET    /api/v1/admins/:id          # Get admin details
PUT    /api/v1/admins/:id          # Update admin
DELETE /api/v1/admins/:id          # Delete admin (super_admin only)
```

### System Logs
```
GET    /api/v1/logs                # List system logs (filter, search)
GET    /api/v1/logs/export         # Export logs (CSV/JSON)
```

### Dashboard
```
GET    /api/v1/dashboard/stats     # Get dashboard statistics
GET    /api/v1/dashboard/charts    # Get chart data
GET    /api/v1/dashboard/recent    # Get recent activities
```

### Health Check
```
GET    /health                     # Health check endpoint
GET    /api/v1/ping               # API ping
```

---

## 🎨 Frontend Sayfa Yapısı

### 1. **Login Page** (`/login`)
- Admin kullanıcı adı ve şifre girişi
- "Beni hatırla" checkbox
- Modern, minimal design
- Logo ve branding

### 2. **Dashboard** (`/`)
- İstatistik kartları (toplam kullanıcı, toplam bakiye, günlük işlem)
- Grafikler (günlük/haftalık/aylık işlem grafikleri)
- Son işlemler tablosu
- Hızlı aksiyonlar (kullanıcı ekle, bakiye yükle)

### 3. **Users Management** (`/users`)
- Kullanıcı listesi (tablo view)
- Arama ve filtreleme (RFID, isim, bakiye)
- Hızlı bakiye yükleme butonu
- Kullanıcı ekleme modal (RFID ID + İsim zorunlu, diğerleri opsiyonel)
- Kullanıcı düzenleme/silme
- Bakiye durumu göstergesi (yeşil: yeterli, kırmızı: düşük)

### 4. **User Detail** (`/users/:id`)
- Kullanıcı profil bilgileri
- Mevcut bakiye (büyük ve belirgin)
- İşlem geçmişi (admin yüklemeleri ve otomasyon düşümleri ayrı gösterilir)
- Bakiye yükleme formu
- Son RFID okutma zamanı
- QR kod (opsiyonel)

### 5. **Transactions** (`/transactions`)
- Tüm işlemler listesi
- Filtreleme (tarih aralığı, işlem tipi, kullanıcı, kaynak)
- İşlem kaynağı badge'i (Admin, Otomasyon, Sistem)
- Export özelliği
- Detaylı görünüm
- Otomasyon işlemlerini farklı renkte göster

### 6. **Admins** (`/admins`)
- Admin kullanıcı listesi
- Yeni admin oluşturma (super_admin only)
- Rol yönetimi
- Admin aktivite görüntüleme

### 7. **System Logs** (`/logs`)
- Sistem logları tablosu
- Gelişmiş filtreleme (log tipi: admin, automation, error)
- RFID okutma logları
- Başarılı/başarısız işlem filtreleme
- Log detay modal
- Export (CSV/JSON)

### 8. **Settings** (`/settings`)
- Sistem ayarları
- Kullanıcı profil ayarları
- Güvenlik ayarları (şifre değiştirme)

---

## 🎨 UI/UX Tasarım Prensipleri

### Renk Paleti
```css
/* Primary Colors */
--primary: #3B82F6;        /* Blue */
--primary-focus: #2563EB;
--primary-content: #FFFFFF;

/* Secondary Colors */
--secondary: #8B5CF6;      /* Purple */
--secondary-focus: #7C3AED;
--secondary-content: #FFFFFF;

/* Accent Colors */
--accent: #10B981;         /* Green */
--accent-focus: #059669;
--accent-content: #FFFFFF;

/* Neutral Colors */
--base-100: #FFFFFF;       /* Background */
--base-200: #F3F4F6;       /* Secondary background */
--base-300: #E5E7EB;       /* Border */
--base-content: #1F2937;   /* Text */

/* Semantic Colors */
--success: #10B981;
--warning: #F59E0B;
--error: #EF4444;
--info: #3B82F6;
```

### Tipografi
```css
/* Font Family */
--font-sans: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;

/* Font Sizes */
--text-xs: 0.75rem;    /* 12px */
--text-sm: 0.875rem;   /* 14px */
--text-base: 1rem;     /* 16px */
--text-lg: 1.125rem;   /* 18px */
--text-xl: 1.25rem;    /* 20px */
--text-2xl: 1.5rem;    /* 24px */
--text-3xl: 1.875rem;  /* 30px */
--text-4xl: 2.25rem;   /* 36px */
```

### Layout System
- **Sidebar Navigation**: Fixed, collapsible
- **Content Area**: Fluid, max-width container
- **Responsive Breakpoints**: 
  - Mobile: < 640px
  - Tablet: 640px - 1024px
  - Desktop: > 1024px
- **Grid System**: CSS Grid & Flexbox kombinasyonu

---

## 📦 Docker Configuration

### Backend Dockerfile (Güncelleme)
```dockerfile
# Multi-stage build
FROM golang:1.23.0-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git gcc musl-dev

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build with optimizations
ENV CGO_ENABLED=1
RUN go build -ldflags="-s -w" -o hudautomata src/main.go

# Final stage
FROM alpine:3.20

WORKDIR /app

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates tzdata

# Copy binary
COPY --from=builder /app/hudautomata /usr/local/bin/hudautomata

# Create non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

EXPOSE 8080

ENV GIN_MODE=release
ENV TZ=Europe/Istanbul

ENTRYPOINT ["hudautomata"]
CMD ["--host", "0.0.0.0", "--port", "8080"]
```

### Frontend Dockerfile (Güncelleme)
```dockerfile
# Build stage
FROM oven/bun:1-alpine AS builder

WORKDIR /app

# Install dependencies
COPY package.json bun.lock ./
RUN bun install --frozen-lockfile

# Copy source and build
COPY . .
RUN bun run build

# Production stage
FROM caddy:2-alpine

WORKDIR /srv

# Copy built files
COPY --from=builder /app/dist /srv

# Copy Caddyfile
COPY Caddyfile /etc/caddy/Caddyfile

EXPOSE 80 443

CMD ["caddy", "run", "--config", "/etc/caddy/Caddyfile", "--adapter", "caddyfile"]
```

### Docker Compose
```yaml
version: '3.9'

services:
  # PostgreSQL Database
  db:
    image: postgres:16-alpine
    container_name: hudautomata_db
    environment:
      POSTGRES_DB: hudautomata
      POSTGRES_USER: huduser
      POSTGRES_PASSWORD: ${DB_PASSWORD:-changeme}
      PGDATA: /var/lib/postgresql/data/pgdata
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./scripts/init.sql:/docker-entrypoint-initdb.d/init.sql
    ports:
      - "5432:5432"
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U huduser"]
      interval: 10s
      timeout: 5s
      retries: 5
    networks:
      - hudautomata_network

  # Backend API
  backend:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: hudautomata_backend
    environment:
      - GIN_MODE=release
      - DB_HOST=db
      - DB_PORT=5432
      - DB_USER=huduser
      - DB_PASSWORD=${DB_PASSWORD:-changeme}
      - DB_NAME=hudautomata
      - JWT_SECRET=${JWT_SECRET:-your-secret-key}
      - CORS_ORIGINS=http://localhost,http://localhost:80
    ports:
      - "8080:8080"
    depends_on:
      db:
        condition: service_healthy
    networks:
      - hudautomata_network
    restart: unless-stopped

  # Frontend
  frontend:
    build:
      context: ./frontend
      dockerfile: Dockerfile
    container_name: hudautomata_frontend
    environment:
      - VITE_API_URL=http://localhost:8080
    ports:
      - "80:80"
      - "443:443"
    depends_on:
      - backend
    networks:
      - hudautomata_network
    restart: unless-stopped

volumes:
  postgres_data:
    driver: local

networks:
  hudautomata_network:
    driver: bridge
```

---

## 🔐 Güvenlik Özellikleri

1. **Authentication & Authorization**
   - JWT token tabanlı kimlik doğrulama
   - Refresh token mekanizması
   - Rol bazlı erişim kontrolü (RBAC)
   - Session management

2. **Input Validation**
   - Tüm inputlar server-side validate edilecek
   - XSS koruması
   - SQL injection koruması (GORM parametreli sorgular)
   - CSRF token (gerekirse)

3. **Rate Limiting**
   - Login endpoint için rate limiting
   - API endpoint'leri için genel rate limiting

4. **Password Security**
   - Bcrypt hashing (cost factor: 12)
   - Güçlü şifre politikası
   - Şifre değiştirme zorunluluğu (opsiyonel)

5. **HTTPS**
   - Caddy ile otomatik SSL/TLS
   - HTTP to HTTPS redirect

6. **Logging & Auditing**
   - Tüm kritik işlemler loglanacak
   - IP adresi ve user agent kaydı
   - Audit trail

---

## 📈 Geliştirme Aşamaları

### Phase 1: Temel Altyapı (3-4 gün)
- [ ] Database schema oluşturma
- [ ] Backend API iskelet yapısı
- [ ] Authentication sistemi (JWT)
- [ ] Frontend layout ve routing
- [ ] Docker compose kurulumu

### Phase 2: Core Features (5-6 gün)
- [ ] Kullanıcı CRUD işlemleri
- [ ] Bakiye yükleme/çıkarma API
- [ ] Transaction yönetimi
- [ ] Frontend sayfaları (Users, Transactions)
- [ ] Form validasyonları

### Phase 3: Admin & Logging (3-4 gün)
- [ ] Admin yönetimi
- [ ] Rol bazlı yetkilendirme
- [ ] System logs
- [ ] Admin activity tracking

### Phase 4: Dashboard & Analytics (3-4 gün)
- [ ] Dashboard istatistikleri
- [ ] Chart.js entegrasyonu
- [ ] Raporlama özellikleri
- [ ] Export fonksiyonları

### Phase 5: UI/UX Polish (2-3 gün)
- [ ] Responsive design optimizasyonu
- [ ] Dark mode
- [ ] Animasyonlar ve transitions
- [ ] Loading states
- [ ] Error handling UI

### Phase 6: Testing & Deployment (2-3 gün)
- [ ] Unit testler
- [ ] Integration testler
- [ ] E2E testler (opsiyonel)
- [ ] Production deployment
- [ ] Dokümantasyon

**Toplam Süre:** ~18-24 gün (yaklaşık 3-4 hafta)

---

## 🚀 Deployment Strategy

### Development Environment
```bash
# Backend
cd hudautomata
go run src/main.go

# Frontend
cd frontend
bun run dev
```

### Production Deployment
```bash
# Docker Compose ile
docker-compose up -d

# Veya manuel
docker build -t hudautomata-backend .
docker build -t hudautomata-frontend ./frontend
```

### Environment Variables
```env
# Backend
DB_HOST=localhost
DB_PORT=5432
DB_USER=huduser
DB_PASSWORD=secure_password
DB_NAME=hudautomata
JWT_SECRET=your-very-secure-secret-key
JWT_EXPIRATION=24h
CORS_ORIGINS=http://localhost:3000,https://your-domain.com
PORT=8080

# Frontend
VITE_API_URL=http://localhost:8080
VITE_APP_NAME=HudAutomata
```

---

## 📝 İlk Admin Kullanıcısı

Sistem ilk kurulumda seed data ile bir super admin oluşturulacak:

```go
// scripts/seed.go
username: "admin"
password: "admin123" // İlk girişte değiştirilmeli
email: "admin@hudautomata.local"
role: "super_admin"
```

---

## 🎯 Başarı Kriterleri

✅ **Fonksiyonel Gereksinimler**
- Tüm CRUD işlemleri çalışıyor
- Authentication ve authorization güvenli
- Responsive tasarım (mobile, tablet, desktop)
- Tüm API endpoint'leri test edilmiş

✅ **Non-Fonksiyonel Gereksinimler**
- Sayfa yükleme süresi < 2 saniye
- API response time < 200ms (ortalama)
- Modern ve kullanıcı dostu arayüz
- Hatasız çalışan Docker ortamı

✅ **Kod Kalitesi**
- Clean code prensipleri
- Yeterli error handling
- Logging mekanizması
- Güvenlik best practices

---

## 📚 Ek Notlar

### RFID Kart Okuyucu Entegrasyonu

**Otomasyon Cihazı (Arduino/ESP32/Raspberry Pi):**
1. RFID kart okuyucu modülü bağlı
2. Kullanıcı kartını okuttuğunda:
   - RFID ID'yi okur
   - Backend'e POST isteği gönderir (`/api/v1/automation/scan`)
   - Hizmet bedelini parametre olarak gönderir
3. Backend response'a göre:
   - Başarılı → Hizmeti verir (röle aktif, motor çalışır vs.)
   - Başarısız → Hata mesajı gösterir (LCD/LED)

**Admin Panel (Frontend):**
- RFID okuyucu simülasyonu input alanı (test için)
- Gerçek okuyucu entegrasyonu:
  - Web Serial API (modern browsers)
  - USB HID device okuyucusu
  - Websocket üzerinden microcontroller bağlantısı

**Örnek Arduino/ESP32 Kodu:**
```cpp
// RFID okutulduğunda
String rfidId = readRFID();
HTTPClient http;
http.begin("http://backend:8080/api/v1/automation/scan");
http.addHeader("Content-Type", "application/json");

String payload = "{\"rfid_card_id\":\"" + rfidId + "\",\"service_cost\":10.5,\"description\":\"Çamaşır makinesi\"}";
int httpCode = http.POST(payload);

if (httpCode == 200) {
  // Hizmet ver
  digitalWrite(RELAY_PIN, HIGH);
} else {
  // Hata göster
  lcd.print("Yetersiz bakiye!");
}
```

### Genişletme Fikirleri (v2.0)
- [ ] Email/SMS bildirimleri
- [ ] QR kod ile ödeme
- [ ] Mobile app (React Native / Flutter)
- [ ] Detaylı raporlama ve analytics
- [ ] Multi-tenant support
- [ ] Backup/restore özelliği

---

## 🐳 Docker Deployment

### Production Setup

**Docker Compose Mimarisi:**
```
[Browser/IoT Device] 
    ↓
[Caddy:80] (Frontend + Reverse Proxy)
    ↓
    ├── /api/* → [Backend:8080] (Go/Gin API)
    │              ↓
    └── /* → [Static Files] (Preact SPA)
                   ↓
              [PostgreSQL:5432] (Database)
```

### Kurulum ve Çalıştırma

#### 1. Environment Variables Ayarlama

`.env` dosyası oluşturun:
```bash
# Database
DB_PASSWORD=your-secure-password-here

# JWT
JWT_SECRET=your-super-secret-jwt-key-here
```

⚠️ **UYARI:** Production'da bu değerleri mutlaka değiştirin!

#### 2. Docker Compose ile Çalıştırma

```bash
# Tüm servisleri build et ve başlat
docker compose up --build -d

# Logları izle
docker compose logs -f

# Sadece backend logları
docker compose logs -f backend

# Sadece frontend logları
docker compose logs -f frontend

# Container durumunu kontrol et
docker compose ps

# Servisleri durdur ve temizle
docker compose down -v
```

#### 3. Servisler

| Servis | Port | Açıklama |
|--------|------|----------|
| **Frontend** | 80, 443 | Caddy web server (Preact UI + API Proxy) |
| **Backend** | 8080 | Go/Gin REST API |
| **Database** | 5432 | PostgreSQL 16 |

#### 4. API Endpoints

**Reverse Proxy ile API erişimi:**
- Production: `http://localhost/api/v1/*`
- Backend Direct: `http://localhost:8080/api/v1/*`

**Health Check:**
```bash
curl http://localhost/api/health
```

**Login Test:**
```bash
curl -X POST http://localhost/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'
```

### Docker Compose Konfigürasyonu

**Services:**

1. **PostgreSQL Database**
   - Image: `postgres:16-alpine`
   - Volume: `postgres_data` (persistent storage)
   - Health check: `pg_isready -U huduser`
   - Auto-creates: `huduser` user and `hudautomata` database

2. **Backend API**
   - Build: Multi-stage Dockerfile (Go 1.23 → Debian slim)
   - Environment: Production mode (GIN_MODE=release)
   - Depends on: Database (waits for health check)
   - Auto-migration: Database tables created on startup
   - Default admin: username=`admin`, password=`admin123`

3. **Frontend**
   - Build: Bun → Vite build → Caddy serve
   - Reverse proxy: `/api/*` → `backend:8080`
   - SPA fallback: All routes → `index.html`
   - Static asset caching: 1 year for JS/CSS
   - GZIP/Zstd compression enabled

### Sorun Giderme

#### Frontend default sayfa gösteriyorsa:
```bash
docker compose down
docker compose build --no-cache frontend
docker compose up -d
```

#### Database bağlantı hatası:
```bash
# DB loglarını kontrol et
docker compose logs db

# Backend loglarını kontrol et
docker compose logs backend
```

#### API 502 Bad Gateway hatası:
```bash
# Backend ayakta mı kontrol et
docker compose ps

# Caddy config'i kontrol et
docker compose exec frontend cat /etc/caddy/Caddyfile

# Servisleri yeniden başlat
docker compose restart
```

#### Database şifre hatası (password authentication failed):
```bash
# Tüm volume'leri temizle ve yeniden oluştur
docker compose down -v
docker compose up --build -d
```

### Önemli Dosyalar

**Modified Files (Docker Setup):**
- `docker-compose.yml` - Servis orkestrasyon
- `Dockerfile` - Backend multi-stage build
- `frontend/Dockerfile` - Frontend multi-stage build (Bun + Caddy)
- `frontend/Caddyfile` - Reverse proxy + SPA routing
- `frontend/.dockerignore` - Build optimizasyonu
- `frontend/src/lib/api.js` - Production/dev API URL logic
- `.env` - Production secrets

### Build Optimizasyonları

**Backend:**
- Static binary (CGO_ENABLED=0)
- Multi-stage build (builder + runtime)
- Minimal image (Debian slim)
- Go mod cache layer

**Frontend:**
- Bun dependency installation (--frozen-lockfile)
- Vite production build
- Asset minification & tree-shaking
- Caddy static file serving
- GZIP/Zstd compression

---

## 🧪 Test & Doğrulama

### Başarı Kriterleri

- [x] Frontend `http://localhost` adresinde HudAutomata UI açılıyor
- [x] Login sayfası çalışıyor (admin/admin123)
- [x] POST `/api/v1/auth/login` başarılı (JWT token dönüyor)
- [x] Database bağlantısı başarılı
- [x] Default admin user otomatik oluşturuluyor
- [x] RFID kullanıcı ekleme çalışıyor
- [x] Dashboard stats gösteriliyor
- [x] Caddy reverse proxy doğru çalışıyor
- [x] Production build optimize edilmiş

### Manual Test Senaryosu

1. **Login Test**
   ```bash
   curl -X POST http://localhost/api/v1/auth/login \
     -H "Content-Type: application/json" \
     -d '{"username":"admin","password":"admin123"}'
   ```

2. **RFID User Oluşturma**
   ```bash
   curl -X POST http://localhost/api/v1/users \
     -H "Authorization: Bearer YOUR_TOKEN" \
     -H "Content-Type: application/json" \
     -d '{
       "rfid_card_id": "RFID12345",
       "name": "Test User",
       "balance": 100.00
     }'
   ```

3. **RFID Automation Scan (IoT Device)**
   ```bash
   curl -X POST http://localhost/api/v1/automation/scan \
     -H "Content-Type: application/json" \
     -d '{
       "rfid_card_id": "RFID12345",
       "service_cost": 25.50,
       "description": "Çamaşır makinesi"
     }'
   ```

4. **Dashboard Stats**
   ```bash
   curl http://localhost/api/v1/dashboard/stats \
     -H "Authorization: Bearer YOUR_TOKEN"
   ```

---

## 🤝 Katkıda Bulunma

Proje açık kaynak olarak geliştirilecek. Katkılar için:
1. Fork the repository
2. Create feature branch
3. Commit changes
4. Push to branch
5. Create Pull Request

---

**Son Güncelleme:** 2025-10-02  
**Versiyon:** 1.0.0  
**Durum:** ✅ Production Ready  
**Yazar:** HudAutomata Team

