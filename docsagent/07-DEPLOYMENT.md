# Deployment Guide

## Production Build

```bash
# Build binary
go build -ldflags="-s -w" -o domesv2 ./cmd/main.go

# Atau dengan UPX compression (jika terinstall)
go build -ldflags="-s -w" -o domesv2 ./cmd/main.go
upx --best domesv2
```

## Server Setup

### 1. Upload ke Server

```bash
# Buat direktori
sudo mkdir -p /opt/domesv2/logs
sudo chmod 755 /opt/domesv2/logs

# Upload file
scp domesv2 user@server:/opt/domesv2/
scp .env user@server:/opt/domesv2/
```

### 2. Systemd Service

File: `/etc/systemd/system/domesv2.service`

```ini
[Unit]
Description=DOMESv2 Backend API
After=network.target mysql.service redis-server.service

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/domesv2
Environment=ENV=production
ExecStart=/opt/domesv2/domesv2
Restart=always
RestartSec=5
StandardOutput=journal+console
StandardError=journal+console

[Install]
WantedBy=multi-user.target
```

```bash
sudo systemctl daemon-reload
sudo systemctl enable domesv2
sudo systemctl start domesv2
sudo systemctl status domesv2
```

### 3. Nginx Reverse Proxy

File: `/etc/nginx/sites-available/domesv2.yourdomain.com`

```nginx
server {
    listen 80;
    server_name domesv2.yourdomain.com;

    location / {
        proxy_pass http://localhost:3000;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

```bash
sudo ln -s /etc/nginx/sites-available/domesv2.yourdomain.com /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl restart nginx
```

### 4. SSL dengan Let's Encrypt

```bash
sudo certbot --nginx -d domesv2.yourdomain.com
```

### 5. Firewall (UFW)

```bash
sudo ufw allow 22/tcp
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw enable
```

## Monitoring

```bash
# Systemd logs
journalctl -u domesv2 -f

# Application logs
tail -f /opt/domesv2/logs/app.log

# Health check
curl https://domesv2.yourdomain.com/api/health-check
```

## Update Aplikasi

```bash
git pull
go build -ldflags="-s -w" -o domesv2 ./cmd/main.go
sudo systemctl restart domesv2
```
