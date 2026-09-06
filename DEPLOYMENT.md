# 🚀 MI-Tech (StarTech Clone) — Production Deployment Guide
**বাংলা এবং ইংরেজিতে সম্পূর্ণ প্রডাকশন ডেপ্লয়মেন্ট গাইড**

---

## 📋 Table of Contents / সূচিপত্র
1. [Server Requirements / সার্ভার রিকোয়ারমেন্ট](#1-server-requirements)
2. [VPS Initial Setup / সার্ভার প্রস্তুতি](#2-vps-initial-setup)
3. [Domain & DNS Configuration / ডোমেইন কনফিগারেশন](#3-domain--dns-configuration)
4. [Clone Project & Setup Environment / প্রজেক্ট ক্লোন ও .env তৈরি](#4-clone-project--setup-environment)
5. [Deploy with Docker Compose / ডকার দিয়ে ডেপ্লয়](#5-deploy-with-docker-compose)
6. [Free SSL Setup (HTTPS) / ফ্রি এসএসএল সার্টিফিকেট](#6-free-ssl-setup-https)
7. [Automated Daily Database Backups / অটোমেটিক ডেটাবেস ব্যাকআপ](#7-automated-daily-database-backups)
8. [Auto-Start on Reboot (Systemd) / সার্ভার রিবুটে অটো-স্টার্ট](#8-auto-start-on-reboot-systemd)
9. [Useful Maintenance Commands / প্রয়োজনীয় কমান্ড](#9-useful-maintenance-commands)
10. [Troubleshooting / সাধারণ সমস্যা ও সমাধান](#10-troubleshooting)

---

## 1. Server Requirements
যেকোনো ক্লাউড প্রোভাইডার (Hostinger VPS, DigitalOcean, Hetzner, AWS EC2, Linode) ব্যবহার করতে পারেন:

| Resource | Minimum | Recommended |
|---|---|---|
| **OS** | Ubuntu 22.04 LTS / 24.04 LTS | Ubuntu 22.04 LTS |
| **RAM** | 1 GB (with 2GB swap) | 2 GB or higher |
| **CPU** | 1 vCPU | 2 vCPU |
| **Storage** | 20 GB SSD | 40 GB+ NVMe SSD |

---

## 2. VPS Initial Setup
আপনার সার্ভারে SSH দিয়ে প্রবেশ করুন:
```bash
ssh root@YOUR_SERVER_IP
```

### ১. সার্ভার আপডেট করুন:
```bash
sudo apt update && sudo apt upgrade -y
sudo apt install -y curl wget git ufw
```

### ২. Docker & Docker Compose ইনস্টল করুন:
```bash
# Docker official install script
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# Verify Docker installation
docker --version
docker compose version
```

### ৩. Firewall (UFW) কনফিগার করুন:
```bash
sudo ufw allow 22/tcp    # SSH
sudo ufw allow 80/tcp    # HTTP
sudo ufw allow 443/tcp   # HTTPS
sudo ufw enable
```

### ৪. (Optional কিন্তু বাঞ্ছনীয়) 2GB Swap তৈরি করুন (যাতে 1GB RAM ক্র্যাশ না করে):
```bash
sudo fallocate -l 2G /swapfile
sudo chmod 600 /swapfile
sudo mkswap /swapfile
sudo swapon /swapfile
echo '/swapfile none swap sw 0 0' | sudo tee -a /etc/fstab
```

---

## 3. Domain & DNS Configuration
আপনার ডোমেইন প্রোভাইডারে (Hostinger, Namecheap, Cloudflare ইত্যাদি) গিয়ে **DNS Management** এ নিচের রেকর্ডগুলো যোগ করুন:

| Type | Name / Host | Value / Points to | TTL |
|---|---|---|---|
| **A** | `@` | `YOUR_SERVER_IP` | Auto / 300 |
| **A** | `www` | `YOUR_SERVER_IP` | Auto / 300 |

*নোট: DNS পরিবর্তন বিশ্বব্যাপী কার্যকর হতে ৫ মিনিট থেকে ২৪ ঘণ্টা লাগতে পারে।*

---

## 4. Clone Project & Setup Environment
প্রজেক্ট সার্ভারে ডাউনলোড করুন:
```bash
cd /var/www
git clone https://github.com/jaberpatwary/startech-clone.git
cd startech-clone
```

### `.env.production` ফাইল তৈরি করুন:
```bash
cp .env.production.example .env.production
nano .env.production
```

**নিচের মানগুলো আপনার নিজস্ব তথ্য দিয়ে পরিবর্তন করুন:**
```ini
APP_ENV=production
DOMAIN_NAME=yourdomain.com
CLIENT_URL=https://yourdomain.com
ALLOWED_ORIGINS=https://yourdomain.com,https://www.yourdomain.com

POSTGRES_DB=startech_prod
POSTGRES_USER=startech_user
POSTGRES_PASSWORD=YOUR_STRONG_DATABASE_PASSWORD_HERE

# Database URL with same credentials
DATABASE_URL=postgres://startech_user:YOUR_STRONG_DATABASE_PASSWORD_HERE@postgres:5432/startech_prod?sslmode=disable

# একটি শক্তিশালী JWT Secret তৈরি করতে চালান: openssl rand -hex 32
JWT_SECRET=paste_generated_random_64_char_secret_here

UPLOAD_DIR=/app/uploads
VITE_API_URL=
```
`Ctrl + O` চেপে Save করুন এবং `Ctrl + X` চেপে বের হন।

---

## 5. Deploy with Docker Compose
আমরা একটি অটোমেটেড ডেপ্লয় স্ক্রিপ্ট যুক্ত করেছি। এটি এক্সিকিউট পারমিশন দিয়ে রান করুন:

```bash
chmod +x deploy.sh scripts/backup-db.sh
./deploy.sh
```

অথবা ম্যানুয়ালি Docker Compose দিয়ে চালাতে পারেন:
```bash
docker compose --env-file .env.production up -d --build
```

কন্টেইনারের স্ট্যাটাস দেখুন:
```bash
docker compose ps
```
সবগুলো কন্টেইনার (`mitech-postgres`, `mitech-backend`, `mitech-frontend`) **healthy** বা **running** দেখাবে! 🎉

এখন ব্রাউজারে `http://YOUR_SERVER_IP` ওপেন করলেই ওয়েবসাইট দেখতে পাবেন।

---

## 6. Free SSL Setup (HTTPS)
Let's Encrypt ও Certbot দিয়ে সম্পূর্ণ ফ্রি SSL সার্টিফিকেট লাগান:

```bash
# Certbot ইনস্টল করুন
sudo apt install -y certbot

# সার্টিফিকেট সংগ্রহ করুন
sudo certbot certonly --standalone -d yourdomain.com -d www.yourdomain.com
```
*(নোট: এই কমান্ড চালানোর সময় পোর্ট ৮০ ফ্রি থাকতে হবে। প্রয়োজনে সাময়িক সময়ের জন্য `docker compose stop frontend` দিয়ে বন্ধ করে সার্টিফিকেট নিয়ে আবার চালু করতে পারেন।)*

সার্টিফিকেট পাওয়ার পর:
```bash
# ssl.conf সক্রিয় করুন
cp nginx/conf.d/ssl.conf.template nginx/conf.d/ssl.conf
sed -i 's/YOUR_DOMAIN_NAME/yourdomain.com/g' nginx/conf.d/ssl.conf

# ফ্রন্টএন্ড কন্টেইনার রিস্টার্ট করুন
docker compose restart frontend
```

### অটো-রিনিউয়াল টেস্ট:
```bash
sudo certbot renew --dry-run
```

---

## 7. Automated Daily Database Backups
প্রতিদিন রাত ৩:০০ টায় অটোমেটিক ডেটাবেস ব্যাকআপ নেওয়ার জন্য Cronjob সেট করুন:

```bash
crontab -e
```
ফাইলের শেষে এই লাইনটি যোগ করুন:
```cron
0 3 * * * /var/www/startech-clone/scripts/backup-db.sh >> /var/log/db-backup.log 2>&1
```

*ব্যাকআপগুলো `/var/backups/mitech-db/` ফোল্ডারে সেভ হবে এবং ৭ দিনের বেশি পুরনো ফাইল অটোমেটিক ডিলিট হবে।*

---

## 8. Auto-Start on Reboot (Systemd)
সার্ভার রিস্টার্ট বা রিবুট হলে যাতে ওয়েবসাইট নিজে থেকেই ব্যাকগ্রাউন্ডে চালু হয়ে যায়:

```bash
sudo cp systemd/mitech.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable mitech
sudo systemctl start mitech
```

---

## 9. Useful Maintenance Commands

```bash
# সব কন্টেইনারের লাইভ লগ দেখতে:
docker compose logs -f

# শুধু ব্যাকএন্ডের লগ দেখতে:
docker compose logs -f backend

# কন্টেইনার বন্ধ করতে:
docker compose down

# নতুন কোড পুশ করার পর সার্ভারে আপডেট দিতে:
./deploy.sh

# ম্যানুয়ালি ইনস্ট্যান্ট ডেটাবেস ব্যাকআপ নিতে:
./scripts/backup-db.sh
```

---

## 10. Troubleshooting

### সমস্যা ১: `port 80 already in use`
**সমাধান:** আগে থেকে অ্যাপাচি বা এনজিনিক্স চালু আছে কিনা দেখুন:
```bash
sudo systemctl stop apache2 2>/dev/null || true
sudo systemctl stop nginx 2>/dev/null || true
sudo systemctl disable apache2 2>/dev/null || true
```

### সমস্যা ২: ডেটাবেস কানেকশন ফেইল (`connection refused`)
**সমাধান:** `postgres` কন্টেইনারটি সুস্থভাবে চালু হতে ৫-১০ সেকেন্ড সময় নেয়। `docker compose ps` দিয়ে দেখুন `(healthy)` আছে কিনা। প্রয়োজনে `docker compose logs postgres` দেখুন।

### সমস্যা ৩: ডিস্ক স্পেস শেষ হয়ে যাচ্ছে
**সমাধান:** পুরনো আনইউজড ডকার ইমেজ পরিষ্কার করুন:
```bash
docker system prune -af --volumes
```

---

### 🛡️ Built-in Security Features:
- ✅ **Nginx Reverse Proxy**: SPA client caching, gzip compression, and rate buffering.
- ✅ **PostgreSQL Connection Pool**: MaxOpenConns, MaxIdleConns, connection lifetime limits.
- ✅ **HttpOnly Cookies**: XSS token extraction prevention.
- ✅ **Rate Limiting**: Echo memory store DDoS & brute-force protection.
- ✅ **Security Headers**: HSTS, X-Frame-Options (Clickjacking defense), nosniff.
- ✅ **Graceful Shutdown**: Zero lost requests or severed DB transactions during rolling updates.
- ✅ **Health Checks**: Integrated Docker & Nginx health endpoints with real DB ping.
