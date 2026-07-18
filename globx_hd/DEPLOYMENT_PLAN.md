# GlobX HD — Deployment Plan

> **Two-Server Architecture on Hostinger VPS**  
> KVM2 → App Server (Frontend + Backend + n8n)  
> KVM4 → Database Server (PostgreSQL)

---

## Table of Contents

1. [Architecture Overview](#1-architecture-overview)
2. [Prerequisites](#2-prerequisites)
3. [Phase 1 — Initial VPS Setup (Both Servers)](#3-phase-1--initial-vps-setup-both-servers)
4. [Phase 2 — Database Server (KVM4)](#4-phase-2--database-server-kvm4)
5. [Phase 3 — App Server (KVM2) — Backend](#5-phase-3--app-server-kvm2--backend)
6. [Phase 4 — App Server (KVM2) — Frontend](#6-phase-4--app-server-kvm2--frontend)
7. [Phase 5 — App Server (KVM2) — n8n](#7-phase-5--app-server-kvm2--n8n)
8. [Phase 6 — Nginx Reverse Proxy & SSL](#8-phase-6--nginx-reverse-proxy--ssl)
9. [Phase 7 — Systemd Services (Auto-Start on Reboot)](#9-phase-7--systemd-services-auto-start-on-reboot)
10. [Phase 8 — Firewall Configuration](#10-phase-8--firewall-configuration)
11. [Phase 9 — DNS Setup](#11-phase-9--dns-setup)
12. [Phase 10 — Final Verification](#12-phase-10--final-verification)
13. [Server Communication Explained](#13-server-communication-explained)
14. [Maintenance & Useful Commands](#14-maintenance--useful-commands)
15. [Troubleshooting](#15-troubleshooting)
16. [Architecture Diagram](#16-architecture-diagram)

---

## 1. Architecture Overview

```
                    ┌─── INTERNET ───┐
                    │                │
                    ▼                │
        ┌──────────────────────┐    │
        │   KVM2 (App Server)  │    │
        │                      │    │
        │  ┌────────────────┐  │    │
        │  │   Nginx        │  │    │  (Reverse Proxy + SSL)
        │  │   Port 80/443  │  │    │
        │  └──┬─────┬───┬───┘  │    │
        │     │     │   │      │    │
        │     ▼     ▼   ▼      │    │
        │  ┌────┐ ┌──┐ ┌───┐  │    │
        │  │Vue │ │Go│ │n8n│  │    │
        │  │:80 │ │  │ │   │  │    │
        │  │    │ │  │ │   │  │    │
        │  │dist│ │  │ │   │  │    │
        │  │    │ │:8│ │:5 │  │    │
        │  │file│ │0 │ │6 │  │    │
        │  │s   │ │8 │ │7 │  │    │
        │  │    │ │0 │ │8 │  │    │
        │  └────┘ └──┘ └───┘  │    │
        │            │         │    │
        └────────────┼─────────┘    │
                     │ PostgreSQL    │
                     │ connection    │
                     │ (port 5432)   │
                     ▼              │
        ┌──────────────────────┐    │
        │  KVM4 (DB Server)    │    │
        │                      │    │
        │  ┌────────────────┐  │    │
        │  │  PostgreSQL    │  │    │
        │  │  Port 5432     │  │    │
        │  │                │  │    │
        │  │  DB: globx_hd  │  │    │
        │  └────────────────┘  │    │
        │                      │    │
        └──────────────────────┘    │
                                    │
```

### What Runs Where

| Server | Component | Port | Purpose |
|--------|-----------|------|---------|
| **KVM2** | Nginx | 80, 443 | Reverse proxy, serves frontend, routes API & n8n |
| **KVM2** | Go Backend | 8080 (internal) | API server |
| **KVM2** | n8n | 5678 (internal) | Workflow automation |
| **KVM4** | PostgreSQL | 5432 | Database (accepts connections from KVM2 only) |

### How the Two Servers Talk

KVM2 connects to KVM4's **public IP** on port **5432** (PostgreSQL). We lock this down with firewall rules so **only KVM2's IP** can reach the database. No one else on the internet can connect to it.

> **Hostinger Note:** Hostinger KVM VPS instances don't have a private network between them by default. We use public IPs + strict firewall rules to simulate private networking securely.

---

## 2. Prerequisites

### What You Need Before Starting

- [ ] Both VPS servers are provisioned and you can access them
- [ ] You know the **public IP** of both servers (find in Hostinger panel)
- [ ] You have **root SSH access** to both servers
- [ ] A **domain name** (e.g., `helpdesk.yourdomain.com`) — optional but recommended for SSL
- [ ] Your code is pushed to a **Git repository** (GitHub/GitLab)

### Notation Used in This Guide

```
KVM2_IP = Your KVM2 public IP (e.g., 154.12.xxx.xxx)
KVM4_IP = Your KVM4 public IP (e.g., 154.12.yyy.yyy)
```

> **Replace `KVM2_IP` and `KVM4_IP` with your actual IPs everywhere in this guide.**

### Recommended OS

Both VPS should run **Ubuntu 22.04 LTS** or **Ubuntu 24.04 LTS** (select during VPS setup in Hostinger).

---

## 3. Phase 1 — Initial VPS Setup (Both Servers)

> **Do these steps on BOTH KVM2 and KVM4.**

### Step 1: Connect via SSH

From your local machine (PowerShell or terminal):

```bash
# Connect to KVM2
ssh root@KVM2_IP

# (In a separate terminal) Connect to KVM4
ssh root@KVM4_IP
```

If using Windows, you can use **PuTTY** or **Windows Terminal** with built-in SSH.

### Step 2: Update the System

```bash
apt update && apt upgrade -y
```

### Step 3: Create a Non-Root User (Security Best Practice)

```bash
# Create user 'deploy'
adduser deploy

# Give sudo privileges
usermod -aG sudo deploy
```

### Step 4: Set Timezone

```bash
timedatectl set-timezone Asia/Kolkata
```

### Step 5: Install Essential Tools

```bash
apt install -y curl wget git ufw nano htop
```

---

## 4. Phase 2 — Database Server (KVM4)

> **All commands in this section run on KVM4.**

### Step 1: Install PostgreSQL

```bash
# Install PostgreSQL
apt install -y postgresql postgresql-contrib

# Verify it's running
systemctl status postgresql
```

You should see `active (running)`.

### Step 2: Configure PostgreSQL User & Database

```bash
# Switch to postgres system user
sudo -u postgres psql
```

Now you're inside the PostgreSQL shell. Run these SQL commands:

```sql
-- Create a strong password (CHANGE THIS!)
CREATE USER globx_user WITH PASSWORD 'YourSuperStrongPassword123!';

-- Create the database
CREATE DATABASE globx_hd OWNER globx_user;

-- Grant all privileges
GRANT ALL PRIVILEGES ON DATABASE globx_hd TO globx_user;

-- Exit
\q
```

> **IMPORTANT:** Replace `YourSuperStrongPassword123!` with a real, strong password. Save it — you'll need it for the backend config.

### Step 3: Allow Remote Connections from KVM2

PostgreSQL by default only accepts connections from `localhost`. We need to allow KVM2 to connect.

#### Edit `postgresql.conf`:

```bash
# Find the config file
sudo nano /etc/postgresql/*/main/postgresql.conf
```

Find the line:
```
#listen_addresses = 'localhost'
```

Change it to:
```
listen_addresses = '*'
```

> This tells PostgreSQL to listen on all network interfaces. Don't worry — we'll restrict who can actually connect in the next step.

#### Edit `pg_hba.conf` (Access Control):

```bash
sudo nano /etc/postgresql/*/main/pg_hba.conf
```

Add this line **at the end** of the file:

```
# Allow KVM2 to connect to globx_hd database
host    globx_hd    globx_user    KVM2_IP/32    scram-sha-256
```

> **Replace `KVM2_IP` with your actual KVM2 IP address.**  
> The `/32` means "only this exact IP". No one else can connect.

**Example:** If KVM2's IP is `154.12.100.50`:
```
host    globx_hd    globx_user    154.12.100.50/32    scram-sha-256
```

### Step 4: Restart PostgreSQL

```bash
sudo systemctl restart postgresql
```

### Step 5: Configure Firewall on KVM4

```bash
# Enable UFW firewall
ufw allow OpenSSH
ufw allow from KVM2_IP to any port 5432

# Enable the firewall
ufw enable

# Verify rules
ufw status
```

> This means: Allow SSH from anywhere + Allow PostgreSQL ONLY from KVM2.

### Step 6: Test the Connection (from KVM4 itself)

```bash
psql -U globx_user -d globx_hd -h localhost
```

Enter the password when prompted. If you see the `globx_hd=>` prompt, the database is working.

Type `\q` to exit.

---

## 5. Phase 3 — App Server (KVM2) — Backend

> **All commands in this section run on KVM2.**

### Step 1: Install Go

```bash
# Download Go (check https://go.dev/dl/ for latest version)
wget https://go.dev/dl/go1.23.4.linux-amd64.tar.gz

# Extract to /usr/local
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.23.4.linux-amd64.tar.gz

# Add to PATH (add to ~/.bashrc for persistence)
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Verify
go version
```

### Step 2: Install PostgreSQL Client (for testing connection)

```bash
apt install -y postgresql-client
```

### Step 3: Test Connection to KVM4 Database

```bash
psql -U globx_user -d globx_hd -h KVM4_IP
```

If you can connect, the two servers are communicating. Type `\q` to exit.

> **If this fails:** Check KVM4's firewall (`ufw status`), pg_hba.conf, and postgresql.conf. Make sure you used the correct KVM2 IP.

### Step 4: Clone Your Repository

```bash
cd /home/deploy
git clone https://github.com/YOUR_USERNAME/YOUR_REPO.git globx_hd
cd globx_hd/ticketing_backend
```

### Step 5: Create the Environment File

```bash
nano .env.email
```

Paste and edit:

```env
# Database - POINT TO KVM4
DATABASE_URL=host=KVM4_IP user=globx_user password=YourSuperStrongPassword123! dbname=globx_hd port=5432 sslmode=disable TimeZone=Asia/Kolkata

# Server
SERVER_ADDRESS=:8080

# JWT Secret (generate a random string - CHANGE THIS!)
JWT_SECRET=your-very-long-random-secret-string-change-me-to-something-unique

# n8n Integration
N8N_API_KEY=your-n8n-api-key-change-me

# Gemini AI
GEMINI_API_KEY=your-gemini-api-key

# Email Notifications (fill with your SMTP details)
EMAIL_SMTP_SERVER=smtp.gmail.com
EMAIL_SMTP_PORT=587
EMAIL_USERNAME=your-email@gmail.com
EMAIL_PASSWORD=your-app-password
EMAIL_NOTIFICATION_ADDRESS=support@yourdomain.com

# File Uploads
UPLOAD_DIR=./uploads
```

> **CRITICAL:** Replace `KVM4_IP` with the actual IP of your KVM4 server.

### Step 6: Build the Backend

```bash
cd /home/deploy/globx_hd/ticketing_backend

# Download dependencies
go mod download

# Build the binary
go build -o globx-backend ./cmd/main.go

```

### Step 7: Create Upload Directory

```bash
mkdir -p /home/deploy/globx_hd/ticketing_backend/uploads
```

### Step 8: Test Run

```bash
./globx-backend
```

You should see GORM migration logs and `Listening on :8080`. Press `Ctrl+C` to stop.

---

## 6. Phase 4 — App Server (KVM2) — Frontend

> **All commands in this section run on KVM2.**

### Step 1: Install Node.js

```bash
# Install Node.js 22.x (LTS)
curl -fsSL https://deb.nodesource.com/setup_22.x | sudo -E bash -
apt install -y nodejs

# Verify
node --version
npm --version
```

### Step 2: Build the Frontend

```bash
cd /home/deploy/globx_hd/ticketing_frontend

# Install dependencies
npm install
```

Before building, configure the API URL:

```bash
# Create production env file
nano .env.production
```

```env
VITE_API_BASE=https://helpdesk.yourdomain.com/api
VITE_API_URL=helpdesk.yourdomain.com
```

> **If you don't have a domain yet**, use your KVM2 IP temporarily:
> ```env
> VITE_API_BASE=http://KVM2_IP/api
> VITE_API_URL=KVM2_IP
> ```

Now build:

```bash
npm run build-only
```

This creates a `dist/` folder with static files. Nginx will serve these.

### Step 3: Copy Build to Nginx Directory

```bash
sudo mkdir -p /var/www/globx-hd
sudo cp -r dist/* /var/www/globx-hd/
```

---

## 7. Phase 5 — App Server (KVM2) — n8n

> **All commands in this section run on KVM2.**

### Step 1: Install n8n

```bash
# Install n8n globally
npm install -g n8n
```

### Step 2: Create n8n Environment File

```bash
mkdir -p /home/deploy/.n8n
nano /home/deploy/.n8n/.env
```

```env
# n8n Configuration
N8N_HOST=0.0.0.0
N8N_PORT=5678
N8N_PROTOCOL=https
WEBHOOK_URL=https://helpdesk.yourdomain.com/n8n/
N8N_BASIC_AUTH_ACTIVE=true
N8N_BASIC_AUTH_USER=admin
N8N_BASIC_AUTH_PASSWORD=YourN8nPassword123!

# Use PostgreSQL for n8n data (optional but recommended)
DB_TYPE=postgresdb
DB_POSTGRESDB_HOST=KVM4_IP
DB_POSTGRESDB_PORT=5432
DB_POSTGRESDB_DATABASE=n8n_db
DB_POSTGRESDB_USER=globx_user
DB_POSTGRESDB_PASSWORD=YourSuperStrongPassword123!
```

> **If using PostgreSQL for n8n**, create the `n8n_db` database on KVM4:
> ```bash
> # On KVM4:
> sudo -u postgres psql
> CREATE DATABASE n8n_db OWNER globx_user;
> \q
> ```
> Also update `pg_hba.conf` on KVM4 to allow access to `n8n_db`:
> ```
> host    n8n_db      globx_user    KVM2_IP/32    scram-sha-256
> ```
> Then restart PostgreSQL: `sudo systemctl restart postgresql`

> **Simpler alternative:** Skip the DB config lines above and n8n will use SQLite locally (fine for small setups).

### Step 3: Test n8n

```bash
n8n start
```

You should see n8n starting on port 5678. Press `Ctrl+C` to stop.

---

## 8. Phase 6 — Nginx Reverse Proxy & SSL

> **All commands in this section run on KVM2.**

### What Nginx Does Here

Nginx sits in front of everything and routes traffic:
- `helpdesk.yourdomain.com/` → Serves Vue.js frontend files
- `helpdesk.yourdomain.com/api/*` → Proxies to Go backend (port 8080)
- `helpdesk.yourdomain.com/ws/*` → Proxies WebSocket to Go backend
- `helpdesk.yourdomain.com/n8n/*` → Proxies to n8n (port 5678)

### Step 1: Install Nginx

```bash
apt install -y nginx
```

### Step 2: Create Nginx Configuration

```bash
sudo nano /etc/nginx/sites-available/globx-hd
```

Paste:

```nginx
server {
    listen 80;
    server_name helpdesk.yourdomain.com;  # CHANGE to your domain or KVM2_IP

    # --- Frontend (Vue.js static files) ---
    root /var/www/globx-hd;
    index index.html;

    # Vue Router: All frontend routes fall back to index.html
    location / {
        try_files $uri $uri/ /index.html;
    }

    # --- Backend API ---
    location /api/ {
        # Strip /api prefix before forwarding
        # e.g., /api/auth/login → /auth/login on backend
        rewrite ^/api/(.*) /$1 break;

        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        # File upload size (adjust if needed)
        client_max_body_size 10M;
    }

    # --- WebSocket ---
    location /ws/ {
        proxy_pass http://127.0.0.1:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $websocket_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        # WebSocket timeout (keep alive for 1 hour)
        proxy_read_timeout 3600s;
        proxy_send_timeout 3600s;
    }

    # --- n8n ---
    location /n8n/ {
        rewrite ^/n8n/(.*) /$1 break;

        proxy_pass http://127.0.0.1:5678;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        # n8n WebSocket support
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
    }

    # --- n8n Webhooks (n8n receives webhooks at /webhook/*) ---
    location /webhook/ {
        proxy_pass http://127.0.0.1:5678;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}

# WebSocket upgrade map
map $http_upgrade $websocket_upgrade {
    default upgrade;
    ''      close;
}
```

> **IMPORTANT:** Replace `helpdesk.yourdomain.com` with your actual domain. If you don't have one yet, use your KVM2 IP address.

### Step 3: Enable the Site

```bash
# Enable the config
sudo ln -s /etc/nginx/sites-available/globx-hd /etc/nginx/sites-enabled/

# Remove the default site
sudo rm /etc/nginx/sites-enabled/default

# Test config for errors
sudo nginx -t

# Reload Nginx
sudo systemctl reload nginx
```

### Step 4: Install SSL Certificate (Free with Let's Encrypt)

> **Requires a domain name pointing to KVM2's IP.** Skip this if using IP only.

```bash
# Install Certbot
apt install -y certbot python3-certbot-nginx

# Get SSL certificate (follow the prompts)
sudo certbot --nginx -d helpdesk.yourdomain.com

# Certbot automatically modifies your Nginx config to add SSL
# It also sets up auto-renewal
```

Verify auto-renewal:
```bash
sudo certbot renew --dry-run
```

---

## 9. Phase 7 — Systemd Services (Auto-Start on Reboot)

These ensure your backend, n8n, and PostgreSQL all start automatically when the server reboots.

### Backend Service (KVM2)

```bash
sudo nano /etc/systemd/system/globx-backend.service
```

```ini
[Unit]
Description=GlobX HD Backend
After=network.target

[Service]
Type=simple
User=deploy
WorkingDirectory=/home/deploy/globx_hd/ticketing_backend
ExecStart=/home/deploy/globx_hd/ticketing_backend/globx-backend
Restart=always
RestartSec=5
Environment=GIN_MODE=release

[Install]
WantedBy=multi-user.target
```

```bash
sudo systemctl daemon-reload
sudo systemctl enable globx-backend
sudo systemctl start globx-backend

# Check status
sudo systemctl status globx-backend
```

### n8n Service (KVM2)

```bash
sudo nano /etc/systemd/system/n8n.service
```

```ini
[Unit]
Description=n8n Workflow Automation
After=network.target

[Service]
Type=simple
User=deploy
ExecStart=/usr/bin/n8n start
Restart=always
RestartSec=5
Environment=N8N_PORT=5678

[Install]
WantedBy=multi-user.target
```

```bash
sudo systemctl daemon-reload
sudo systemctl enable n8n
sudo systemctl start n8n

# Check status
sudo systemctl status n8n
```

### PostgreSQL (KVM4)

PostgreSQL auto-starts by default after installation. Verify:

```bash
# On KVM4
sudo systemctl enable postgresql
sudo systemctl status postgresql
```

---

## 10. Phase 8 — Firewall Configuration

### KVM2 (App Server) Firewall

```bash
# Allow SSH
ufw allow OpenSSH

# Allow HTTP and HTTPS (for Nginx)
ufw allow 80/tcp
ufw allow 443/tcp

# DO NOT open 8080 or 5678 to the public - Nginx handles routing
# These ports are only accessed internally via 127.0.0.1

# Enable firewall
ufw enable

# Verify
ufw status
```

Expected output:
```
To                         Action      From
--                         ------      ----
OpenSSH                    ALLOW       Anywhere
80/tcp                     ALLOW       Anywhere
443/tcp                    ALLOW       Anywhere
```

### KVM4 (DB Server) Firewall

Already configured in Phase 2, but verify:

```bash
ufw status
```

Expected:
```
To                         Action      From
--                         ------      ----
OpenSSH                    ALLOW       Anywhere
5432                       ALLOW       KVM2_IP
```

> **Key Point:** Port 5432 is ONLY accessible from KVM2's IP. The database is invisible to the rest of the internet.

---

## 11. Phase 9 — DNS Setup

If you have a domain, set up DNS records in your domain registrar (or Hostinger DNS):

| Type | Name | Value | TTL |
|------|------|-------|-----|
| A | helpdesk | KVM2_IP | 3600 |

Example: If your domain is `yourdomain.com` and you add an A record for `helpdesk`, your app will be at `helpdesk.yourdomain.com`.

> **No domain?** You can access everything via `http://KVM2_IP` directly. SSL won't work without a domain though.

---

## 12. Phase 10 — Final Verification

### Checklist

Run these checks in order:

#### 1. Database (from KVM2)
```bash
psql -U globx_user -d globx_hd -h KVM4_IP -c "SELECT 1;"
```
✅ Should return `1`.

#### 2. Backend
```bash
curl http://localhost:8080/n8n/health
```
✅ Should return a JSON health check response.

#### 3. Backend via Nginx
```bash
curl http://localhost/api/n8n/health
```
✅ Same response as above, but routed through Nginx.

#### 4. Frontend
Open browser: `http://helpdesk.yourdomain.com` (or `http://KVM2_IP`)  
✅ Should see the login page.

#### 5. n8n
Open browser: `http://helpdesk.yourdomain.com/n8n/`  
✅ Should see n8n interface.

#### 6. WebSocket
Login to the app, check browser console for:
```
✅ WebSocket connected
```

#### 7. SSL (if configured)
Open browser: `https://helpdesk.yourdomain.com`  
✅ Should see padlock icon in address bar.

---

## 13. Server Communication Explained

Here's exactly how data flows between the two servers:

```
User's Browser
     │
     │ HTTPS request (port 443)
     ▼
┌─────────────────────────────────────────────────────────┐
│ KVM2 - Nginx (port 80/443)                              │
│                                                         │
│   /             → Serves Vue.js static files            │
│   /api/*        → Proxies to Go backend (127.0.0.1:8080)│
│   /ws/*         → Proxies WebSocket to Go (127.0.0.1:8080)│
│   /n8n/*        → Proxies to n8n (127.0.0.1:5678)       │
│   /webhook/*    → Proxies to n8n (127.0.0.1:5678)       │
└─────────────────────────────────────────────────────────┘
                          │
                          │ Go backend makes DB queries
                          │ using DATABASE_URL in .env
                          │
                          │ TCP connection to KVM4_IP:5432
                          │ (goes over public internet but
                          │  encrypted by PostgreSQL + firewalled)
                          ▼
┌─────────────────────────────────────────────────────────┐
│ KVM4 - PostgreSQL (port 5432)                           │
│                                                         │
│   Only accepts connections from KVM2_IP                  │
│   Uses scram-sha-256 password authentication             │
│   Firewall blocks ALL other IPs                          │
└─────────────────────────────────────────────────────────┘
```

### Why is this secure?

1. **Firewall (UFW):** KVM4 only allows port 5432 from KVM2's exact IP
2. **PostgreSQL pg_hba.conf:** Only `globx_user` from KVM2's IP can connect to `globx_hd` database
3. **Password auth:** Even if someone spoofs KVM2's IP, they need the database password
4. **No public DB port:** Port 5432 is invisible to port scanners (blocked by firewall for all other IPs)

### Optional: Extra Security with SSL for PostgreSQL

For production with sensitive data, you can enable SSL on the PostgreSQL connection:

```env
# In .env.email on KVM2, change sslmode:
DATABASE_URL=host=KVM4_IP user=globx_user password=... dbname=globx_hd port=5432 sslmode=require TimeZone=Asia/Kolkata
```

This encrypts all data between KVM2 and KVM4. PostgreSQL supports SSL out of the box (you'll need to configure certificates on KVM4).

---

## 14. Maintenance & Useful Commands

### Checking Service Status (KVM2)
```bash
# Backend
sudo systemctl status globx-backend

# n8n
sudo systemctl status n8n

# Nginx
sudo systemctl status nginx
```

### Viewing Logs (KVM2)
```bash
# Backend logs (live)
sudo journalctl -u globx-backend -f

# n8n logs (live)
sudo journalctl -u n8n -f

# Nginx access logs
sudo tail -f /var/log/nginx/access.log

# Nginx error logs
sudo tail -f /var/log/nginx/error.log
```

### Restarting Services (KVM2)
```bash
sudo systemctl restart globx-backend
sudo systemctl restart n8n
sudo systemctl restart nginx
```

### Deploying Updates

#### Method A: Full Pull (Updates all files)
Use this if you want to keep the entire local repository in sync with the remote:

```bash
cd /home/deploy/globx_hd
git pull origin master
```

#### Method B: Selective Pull (Updates only specific files)
Use this if you have local modifications on the VPS that you don't want to lose, or if you only want to pull specific files:

```bash
cd /home/deploy/globx_hd

# 1. Fetch remote changes without merging
git fetch origin master

# 2. Checkout/restore only the specific updated files from the remote branch
git checkout origin/master -- ticketing_backend/internal/models/models.go
git checkout origin/master -- ticketing_backend/internal/handlers/contacts.go
git checkout origin/master -- ticketing_backend/internal/handlers/auth.go
git checkout origin/master -- ticketing_backend/internal/handlers/n8n_webhook.go
git checkout origin/master -- ticketing_backend/internal/services/audit_service.go
git checkout origin/master -- ticketing_backend/migrations/009_make_contact_auth_optional.sql
git checkout origin/master -- ticketing_backend/internal/handlers/audit_testing_endpoint.go
git checkout origin/master -- ticketing_frontend/src/components/Contacts.vue
git checkout origin/master -- ticketing_frontend/src/components/shared/ContactCreateModal.vue
```

#### After Pulling: Rebuild & Restart Services

##### 1. Backend:
```bash
cd /home/deploy/globx_hd/ticketing_backend
go build -o globx-backend ./cmd/main.go
sudo systemctl restart globx-backend
```

##### 2. Frontend:
```bash
cd /home/deploy/globx_hd/ticketing_frontend
npm install
npm run build-only
sudo cp -r dist/* /var/www/globx-hd/
# No restart needed - Nginx automatically serves the updated static files
```

### Database Backup (KVM4)
```bash
# Create a backup
pg_dump -U globx_user -h localhost globx_hd > /home/deploy/backups/globx_hd_$(date +%Y%m%d).sql

# Restore from backup
psql -U globx_user -h localhost globx_hd < /home/deploy/backups/globx_hd_20250712.sql
```

#### Automated Daily Backup (KVM4):
```bash
# Create backup directory
mkdir -p /home/deploy/backups

# Create backup script
nano /home/deploy/backup.sh
```

```bash
#!/bin/bash
BACKUP_DIR="/home/deploy/backups"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
export PGPASSWORD="YourSuperStrongPassword123!"

pg_dump -U globx_user -h localhost globx_hd > "$BACKUP_DIR/globx_hd_$TIMESTAMP.sql"

# Keep only last 7 days of backups
find $BACKUP_DIR -name "globx_hd_*.sql" -mtime +7 -delete
```

```bash
chmod +x /home/deploy/backup.sh

# Add to crontab (runs daily at 2 AM)
crontab -e
# Add this line:
0 2 * * * /home/deploy/backup.sh
```

---

## 15. Troubleshooting

### "Connection refused" when backend tries to reach database

| Check | Command (on KVM4) |
|-------|-------------------|
| PostgreSQL running? | `sudo systemctl status postgresql` |
| Listening on all interfaces? | `grep listen_addresses /etc/postgresql/*/main/postgresql.conf` |
| pg_hba.conf has KVM2 IP? | `cat /etc/postgresql/*/main/pg_hba.conf \| grep KVM2_IP` |
| Firewall allows KVM2? | `sudo ufw status` |
| Correct port? | `ss -tlnp \| grep 5432` |

### "502 Bad Gateway" from Nginx

| Check | Command (on KVM2) |
|-------|-------------------|
| Backend running? | `sudo systemctl status globx-backend` |
| Backend listening on 8080? | `ss -tlnp \| grep 8080` |
| Backend logs? | `sudo journalctl -u globx-backend -n 50` |

### Frontend loads but API calls fail

- Check browser console (F12 → Console) for errors
- Verify `VITE_API_BASE` in `.env.production` matches your Nginx config
- Test API directly: `curl http://localhost/api/auth/login`

### WebSocket won't connect

- Check Nginx WebSocket config (upgrade headers)
- Check browser console for WS connection errors
- Test: `curl -i -N -H "Connection: Upgrade" -H "Upgrade: websocket" http://localhost/ws/notifications`

### n8n can't reach backend API

- From KVM2, test: `curl http://127.0.0.1:8080/n8n/health`
- In n8n workflows, use `http://127.0.0.1:8080` as the backend URL (internal communication)

---

## 16. Architecture Diagram

### Request Flow Summary

```
┌──────────────────────────────────────────────────────────────────┐
│                        USER'S BROWSER                            │
│                                                                  │
│   https://helpdesk.yourdomain.com                                │
│   ├── /                    → Login page, dashboard, etc.         │
│   ├── /api/auth/login      → Login API                           │
│   ├── /api/manager/tickets → Ticket list API                     │
│   ├── /ws/notifications    → Real-time notifications             │
│   └── /n8n/               → n8n workflow editor                  │
└─────────────────────────────┬────────────────────────────────────┘
                              │
                              ▼ (HTTPS / port 443)
┌─────────────────────────────────────────────────────────────────┐
│                     KVM2 — APP SERVER                            │
│                                                                  │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │ NGINX (Reverse Proxy + SSL Termination)                  │    │
│  │                                                         │    │
│  │  /           → /var/www/globx-hd/index.html (Vue SPA)   │    │
│  │  /api/*      → proxy_pass http://127.0.0.1:8080         │    │
│  │  /ws/*       → proxy_pass ws://127.0.0.1:8080           │    │
│  │  /n8n/*      → proxy_pass http://127.0.0.1:5678         │    │
│  │  /webhook/*  → proxy_pass http://127.0.0.1:5678         │    │
│  └──────────┬──────────────────────┬───────────────────────┘    │
│             │                      │                             │
│             ▼                      ▼                             │
│  ┌──────────────────┐   ┌──────────────────┐                   │
│  │ GO BACKEND        │   │ n8n              │                   │
│  │ (port 8080)       │   │ (port 5678)      │                   │
│  │                   │   │                  │                   │
│  │ • REST API        │   │ • Email → Ticket │                   │
│  │ • WebSocket Hub   │   │   workflows      │                   │
│  │ • Auth (JWT)      │   │ • Calls backend  │                   │
│  │ • File uploads    │   │   API internally  │                   │
│  │ • Audit logging   │   │ • Calls Gemini   │                   │
│  └────────┬──────────┘   │   AI API         │                   │
│           │               └────────┬─────────┘                   │
│           │                        │                             │
│           │    SQL queries         │    SQL queries (if PG mode) │
│           └────────┬───────────────┘                             │
│                    │                                             │
└────────────────────┼─────────────────────────────────────────────┘
                     │
                     │ TCP connection to KVM4_IP:5432
                     │ (firewall: ONLY KVM2 IP allowed)
                     ▼
┌─────────────────────────────────────────────────────────────────┐
│                     KVM4 — DATABASE SERVER                       │
│                                                                  │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │ POSTGRESQL (port 5432)                                   │    │
│  │                                                         │    │
│  │  Database: globx_hd    (application data)                │    │
│  │  Database: n8n_db      (n8n workflow data) [optional]    │    │
│  │                                                         │    │
│  │  User: globx_user                                        │    │
│  │  Auth: scram-sha-256                                     │    │
│  │  Access: KVM2_IP/32 ONLY                                 │    │
│  └─────────────────────────────────────────────────────────┘    │
│                                                                  │
│  Firewall: SSH (anywhere) + 5432 (KVM2_IP only)                 │
│  Daily backups via cron at 2 AM                                  │
└─────────────────────────────────────────────────────────────────┘
```

### Port Summary

| Server | Port | Accessible From | Service |
|--------|------|----------------|---------|
| KVM2 | 22 | Anywhere (SSH) | SSH |
| KVM2 | 80 | Anywhere | Nginx HTTP (redirects to 443) |
| KVM2 | 443 | Anywhere | Nginx HTTPS |
| KVM2 | 8080 | localhost only | Go Backend (internal) |
| KVM2 | 5678 | localhost only | n8n (internal) |
| KVM4 | 22 | Anywhere (SSH) | SSH |
| KVM4 | 5432 | KVM2_IP only | PostgreSQL |

---

### Quick Reference: Full Deployment Commands

<details>
<summary><b>Click to expand: All commands in order</b></summary>

#### KVM4 (Database) — Do this first
```bash
# 1. SSH in
ssh root@KVM4_IP

# 2. Update
apt update && apt upgrade -y

# 3. Set timezone
timedatectl set-timezone Asia/Kolkata

# 4. Install PostgreSQL
apt install -y postgresql postgresql-contrib

# 5. Create DB + user
sudo -u postgres psql -c "CREATE USER globx_user WITH PASSWORD 'YourPassword';"
sudo -u postgres psql -c "CREATE DATABASE globx_hd OWNER globx_user;"
sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE globx_hd TO globx_user;"

# 6. Edit postgresql.conf: listen_addresses = '*'
sudo sed -i "s/#listen_addresses = 'localhost'/listen_addresses = '*'/" /etc/postgresql/*/main/postgresql.conf

# 7. Add KVM2 to pg_hba.conf
echo "host    globx_hd    globx_user    KVM2_IP/32    scram-sha-256" | sudo tee -a /etc/postgresql/*/main/pg_hba.conf

# 8. Restart PostgreSQL
sudo systemctl restart postgresql

# 9. Firewall
ufw allow OpenSSH
ufw allow from KVM2_IP to any port 5432
ufw enable
```

#### KVM2 (App) — Do this second
```bash
# 1. SSH in
ssh root@KVM2_IP

# 2. Update + install tools
apt update && apt upgrade -y
apt install -y curl wget git ufw nano htop nginx postgresql-client certbot python3-certbot-nginx

# 3. Set timezone
timedatectl set-timezone Asia/Kolkata

# 4. Install Go
wget https://go.dev/dl/go1.23.4.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.23.4.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc && source ~/.bashrc

# 5. Install Node.js
curl -fsSL https://deb.nodesource.com/setup_22.x | sudo -E bash -
apt install -y nodejs

# 6. Install n8n
npm install -g n8n

# 7. Clone repo
cd /home && mkdir -p deploy && cd deploy
git clone YOUR_REPO_URL globx_hd

# 8. Build backend
cd globx_hd/ticketing_backend
# Create .env.email (see Phase 3, Step 5)
go mod download
go build -o globx-backend ./cmd/main.go
mkdir -p uploads

# 9. Build frontend
cd ../ticketing_frontend
# Create .env.production (see Phase 4, Step 2)
npm install
npm run build-only
sudo mkdir -p /var/www/globx-hd
sudo cp -r dist/* /var/www/globx-hd/

# 10. Configure Nginx (see Phase 6, Step 2)
# 11. Create systemd services (see Phase 7)
# 12. Enable firewall
ufw allow OpenSSH
ufw allow 80/tcp
ufw allow 443/tcp
ufw enable

# 13. Start services
sudo systemctl start globx-backend
sudo systemctl start n8n

# 14. SSL (if domain is ready)
sudo certbot --nginx -d helpdesk.yourdomain.com
```

</details>

---

### Important Reminders

> ⚠️ **Change all passwords** in this guide to real, strong passwords  
> ⚠️ **Never commit `.env` files** to Git  
> ⚠️ **Keep both servers updated**: `apt update && apt upgrade -y` regularly  
> ⚠️ **Monitor disk space**: `df -h` (especially uploads folder and backups)  
> ⚠️ **Test backups**: Periodically restore a backup to verify it works  

---

*End of Deployment Plan*
