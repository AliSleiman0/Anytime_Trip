# CI/CD Setup Guide for Anytime Trip

This guide will help you set up the CI/CD pipeline for the Anytime Trip project.

## Overview

The CI/CD pipeline handles:
- ✅ Backend (Go) - Testing, linting, and building
- ✅ Frontend Admin - Building Tailwind CSS
- ✅ Frontend Superadmin - Validation
- 🔒 Security scanning
- 🚀 Automated deployment to staging and production

## Prerequisites

### 1. GitHub Repository Setup
- Repository must be hosted on GitHub
- Admin access to repository settings

### 2. Server Requirements
- Staging server (optional but recommended)
- Production server
- SSH access to both servers
- Go 1.21+ installed on servers
- Node.js 20+ installed on servers
- Systemd service for backend

## GitHub Secrets Configuration

Navigate to your GitHub repository:
**Settings → Secrets and variables → Actions → New repository secret**

### Required Secrets

#### Staging Environment
```
STAGING_HOST          - IP address or domain (e.g., 192.168.1.100 or staging.anytimetrip.com)
STAGING_USER          - SSH username (e.g., deploy)
STAGING_SSH_KEY       - Private SSH key for authentication
STAGING_PORT          - SSH port (default: 22)
```

#### Production Environment
```
PRODUCTION_HOST       - IP address or domain
PRODUCTION_USER       - SSH username
PRODUCTION_SSH_KEY    - Private SSH key for authentication
PRODUCTION_PORT       - SSH port (default: 22)
```

## Server Setup

### 1. SSH Key Generation

On your local machine:
```bash
# Generate SSH key pair
ssh-keygen -t ed25519 -C "github-actions-deploy" -f ~/.ssh/github_deploy

# Copy public key to server
ssh-copy-id -i ~/.ssh/github_deploy.pub user@your-server
```

Add the **private key** (`~/.ssh/github_deploy`) content to GitHub Secrets.

### 2. Server Directory Structure

On both staging and production servers:
```bash
# Create project directory
mkdir -p ~/anytime-trip
cd ~/anytime-trip

# Clone repository
git clone https://github.com/your-username/anytime-trip.git .

# Create backup directory
mkdir -p ~/backups

# Create log file
touch ~/deployment.log
```

### 3. Systemd Service Setup

Create `/etc/systemd/system/anytime-backend.service`:

```ini
[Unit]
Description=Anytime Trip Backend Service
After=network.target

[Service]
Type=simple
User=deploy
WorkingDirectory=/home/deploy/anytime-trip/backend
ExecStart=/home/deploy/anytime-trip/backend/bin/server
Restart=always
RestartSec=5

# Environment variables
Environment="PORT=8080"
Environment="ENV=production"

# Security
NoNewPrivileges=true
PrivateTmp=true

[Install]
WantedBy=multi-user.target
```

Enable and start the service:
```bash
sudo systemctl daemon-reload
sudo systemctl enable anytime-backend
sudo systemctl start anytime-backend
```

### 4. Nginx Configuration (Optional)

For serving frontend and proxying backend:

```nginx
server {
    listen 80;
    server_name anytimetrip.com;

    # Frontend Admin
    location /admin {
        alias /home/deploy/anytime-trip/frontend/admin;
        try_files $uri $uri/ /admin/index.html;
    }

    # Frontend Superadmin
    location /superadmin {
        alias /home/deploy/anytime-trip/frontend/superadmin;
        try_files $uri $uri/ /superadmin/index.html;
    }

    # Backend API
    location /api {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
    }
}
```

## Workflow Triggers

### Automatic Triggers

1. **Push to main** → Runs tests + deploys to production
2. **Push to staging** → Runs tests + deploys to staging
3. **Push to develop** → Runs tests only
4. **Pull Request** → Runs tests only

### Manual Triggers

Go to **Actions → CI/CD Pipeline → Run workflow**

Options:
- Deploy to staging
- Deploy to production
- Rollback production (restores from latest backup)

## Branch Strategy

```
main (production)
  ↑
develop
  ↑
staging
  ↑
feature branches
```

### Recommended Workflow
1. Create feature branch from `develop`
2. Make changes and push
3. Create PR to `develop` → CI runs
4. Merge to `develop` → CI runs
5. When ready for staging: merge `develop` → `staging` → Auto-deploys
6. When ready for production: merge `staging` → `main` → Auto-deploys

## Testing the Pipeline

### 1. Test CI (No Deployment)

```bash
git checkout -b test-ci
# Make a small change
echo "# Test" >> README.md
git add README.md
git commit -m "test: CI pipeline"
git push origin test-ci
```

Create a PR to `develop` and watch the Actions tab.

### 2. Test Staging Deployment

```bash
git checkout staging
git merge develop
git push origin staging
```

Monitor in Actions tab. Check staging server after deployment.

### 3. Test Production Deployment

```bash
git checkout main
git merge staging
git push origin main
```

⚠️ **Warning**: This deploys to production!

## Monitoring Deployments

### GitHub Actions Dashboard
- Go to **Actions** tab
- Click on a workflow run
- View logs for each job

### Server Logs
```bash
# Backend service logs
sudo journalctl -u anytime-backend -f

# Deployment log
tail -f ~/deployment.log

# System logs
tail -f /var/log/syslog
```

## Troubleshooting

### Build Fails on Backend
```bash
# Check Go version
go version

# Verify dependencies
cd backend
go mod verify
go mod tidy
```

### Frontend Build Fails
```bash
# Check Node version
node --version

# Clear cache and reinstall
cd frontend/admin
rm -rf node_modules package-lock.json
npm install
```

### Deployment Fails (SSH Issues)
```bash
# Test SSH connection locally
ssh -i ~/.ssh/github_deploy user@server

# Check server firewall
sudo ufw status

# Check SSH service
sudo systemctl status ssh
```

### Service Won't Start
```bash
# Check service status
sudo systemctl status anytime-backend

# View service logs
sudo journalctl -u anytime-backend -n 50

# Check if port is in use
sudo netstat -tulpn | grep 8080
```

## Rollback Procedure

### Automatic Rollback (Using Backup)
1. Go to **Actions** → **CI/CD Pipeline**
2. Click **Run workflow**
3. Select `rollback` for deploy_environment
4. Click **Run workflow**

### Manual Rollback
```bash
# On production server
cd ~/anytime-trip

# View available backups
ls -lh ~/backups/

# Restore specific backup
tar -xzf ~/backups/anytime-trip-20260309-143000.tar.gz -C ~/anytime-trip

# Restart service
sudo systemctl restart anytime-backend
```

## Security Best Practices

1. **Never commit secrets** to the repository
2. **Rotate SSH keys** regularly
3. **Use environment-specific secrets** for staging/production
4. **Enable branch protection** on main and staging branches
5. **Require PR reviews** before merging
6. **Enable security scanning** (Dependabot, CodeQL)
7. **Monitor deployment logs** for suspicious activity
8. **Keep backups** for at least 30 days

## Maintenance

### Weekly Tasks
- Review deployment logs
- Check backup integrity
- Update dependencies

### Monthly Tasks
- Rotate SSH keys
- Review security scan results
- Test rollback procedure

### As Needed
- Update Go version in workflow
- Update Node.js version in workflow
- Adjust resource limits
- Scale infrastructure

## Support

For issues or questions:
1. Check workflow logs in Actions tab
2. Review server logs
3. Consult this documentation
4. Contact DevOps team

## Next Steps

- [ ] Set up GitHub Secrets
- [ ] Configure servers
- [ ] Test CI pipeline
- [ ] Test staging deployment
- [ ] Test production deployment
- [ ] Document environment variables
- [ ] Set up monitoring/alerts
- [ ] Configure backup retention
- [ ] Implement health checks
- [ ] Set up SSL certificates
