# GitHub Secrets Reference

Quick reference for all secrets needed in the CI/CD pipeline.

## How to Add Secrets

1. Go to your GitHub repository
2. Navigate to **Settings** → **Secrets and variables** → **Actions**
3. Click **New repository secret**
4. Add name and value
5. Click **Add secret**

## DigitalOcean Droplet Secrets

> Both staging and production share the same DigitalOcean droplet (`104.248.138.150`).
> Authentication uses **password-based SSH** via `appleboy/ssh-action`.

| Secret Name | Description | Value to Set |
|------------|-------------|------|
| `DEPLOY_HOST` | Droplet IP address | `104.248.138.150` |
| `DEPLOY_USER` | SSH login username | `root` |
| `DEPLOY_PASSWORD` | SSH login password | *(droplet root password)* |

## Optional Application Secrets

If your application needs environment variables during deployment:

| Secret Name | Description | Example |
|------------|-------------|---------|
| `DATABASE_URL` | Database connection string | `postgres://user:pass@host:5432/db` |
| `JWT_SECRET` | JWT signing secret | `your-super-secret-jwt-key-here` |
| `EMAIL_USER` | SMTP email username | `noreply@anytimetrip.com` |
| `EMAIL_PASSWORD` | SMTP email password | `app-specific-password` |
| `API_KEY` | External API key | `sk_live_...` |

## Testing SSH Connection

```bash
# Test password login to the DigitalOcean droplet
ssh root@104.248.138.150

# If successful, the server is reachable and your DEPLOY_PASSWORD secret is correct
```

## Secret Security Checklist

- [ ] All secrets are stored in GitHub Secrets (not in code)
- [ ] `DEPLOY_HOST`, `DEPLOY_USER`, `DEPLOY_PASSWORD` are set in GitHub repository secrets
- [ ] Droplet password is a strong password
- [ ] Server firewall only allows necessary ports (22, 80, 443, 8080)
- [ ] Secrets are documented (but values kept secret)
- [ ] Team members with access are documented

## Rotating Secrets

### When to Rotate:
- Every 90 days (recommended)
- When a team member leaves
- After a security incident
- When secrets may have been exposed

### How to Rotate:

1. **Change the password on the DigitalOcean droplet**
   ```bash
   ssh root@104.248.138.150
   passwd  # enter new password
   ```

2. **Update GitHub Secret**
   - Go to Settings → Secrets → Actions
   - Edit `DEPLOY_PASSWORD`
   - Paste new password
   - Save

3. **Test deployment**
   - Trigger a manual workflow run
   - Verify successful deployment

## Troubleshooting

### "Permission denied (password)"
- Verify `DEPLOY_PASSWORD` secret value matches the droplet root password
- Confirm password login is enabled: `PasswordAuthentication yes` in `/etc/ssh/sshd_config`
- Try logging in manually: `ssh root@104.248.138.150`

### "Host key verification failed"
- `appleboy/ssh-action` handles this automatically

### "Connection refused"
- Verify droplet IP is `104.248.138.150`
- Check firewall allows port 22: `sudo ufw allow ssh`
- Verify SSH service is running: `sudo systemctl status ssh`

### Secret not updating
- Check secret name matches exactly (case-sensitive): `DEPLOY_HOST`, `DEPLOY_USER`, `DEPLOY_PASSWORD`
- Verify secrets are set at repository level (Settings → Secrets → Actions)

## Environment-Specific Secrets

If using GitHub Environments (recommended):

1. Go to **Settings** → **Environments**
2. Create `staging` and `production` environments
3. Add environment-specific secrets to each
4. Add protection rules (require approval for production)

This provides:
- Environment-specific secret values
- Deployment approvals
- Better audit trail
- Deployment URLs

## Backup Your Secrets

**Important**: Store secrets securely offline

### Recommended:
- Password manager (1Password, Bitwarden, LastPass)
- Encrypted file on secure storage
- Company secrets management system (Vault, AWS Secrets Manager)

### Never:
- Plain text files
- Unencrypted email
- Slack/Discord messages
- Shared documents
- Commit to git

## Contact

For secret-related issues:
- DevOps team lead
- Security team
- System administrator

**Remember**: Never share secrets through insecure channels!
