# CI/CD Pipeline Setup

This repository includes a GitHub Actions workflow that automatically deploys both your Go backend service and React frontend to a remote server when code is pushed to the `main` branch.

## Required GitHub Secrets

You need to configure the following secrets in your GitHub repository settings (`Settings > Secrets and variables > Actions`):

### Server Configuration
- `SERVER_HOST`: The IP address or hostname of your remote server
- `SERVER_USER`: The username to SSH into your server (e.g., `ubuntu`, `root`)
- `SERVER_SSH_KEY`: The private SSH key for accessing your server
- `SERVER_PORT`: The SSH port (usually `22`)

## How to Set Up Secrets

1. Go to your GitHub repository
2. Click on `Settings` tab
3. In the left sidebar, click `Secrets and variables` > `Actions`
4. Click `New repository secret` for each secret above

### SSH Key Setup

1. Generate an SSH key pair on your local machine:
   ```bash
   ssh-keygen -t rsa -b 4096 -f ~/.ssh/deploy_key
   ```

2. Copy the public key to your server:
   ```bash
   ssh-copy-id -i ~/.ssh/deploy_key.pub user@your-server-ip
   ```

3. Copy the private key content and add it as `SERVER_SSH_KEY` secret:
   ```bash
   cat ~/.ssh/deploy_key
   ```

## Pipeline Process

The CI/CD pipeline will:

1. **Test Backend**: Run Go tests to ensure backend code quality
2. **Test Frontend**: Install dependencies, lint, and build React frontend
3. **Build & Deploy**: 
   - Build frontend production assets
   - Build Docker image for Go backend
   - Copy both frontend and backend to your server
   - Deploy backend container on port 8080
   - Deploy frontend to `/var/www/html`
   - Verify both deployments

## Server Requirements

Your remote server needs:
- Docker installed
- Nginx installed (for serving frontend)
- SSH access enabled
- Port 80 (frontend) and 8080 (backend API) available

### Server Setup Commands

Run these commands on your server to prepare for deployment:

```bash
# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# Install Nginx
sudo apt update
sudo apt install nginx -y

# Create web directory
sudo mkdir -p /var/www/html

# Copy the nginx configuration
sudo cp nginx.conf /etc/nginx/sites-available/default
sudo nginx -t
sudo systemctl restart nginx
sudo systemctl enable nginx

# Open firewall ports
sudo ufw allow 80
sudo ufw allow 8080
```

## Architecture

- **Frontend**: React/Vite app **built to static files** and served by Nginx on port 80
- **Backend**: Go service running in Docker on port 8080  
- **API Proxy**: Nginx proxies `/api/*` requests to the Go backend

### Production Frontend Deployment

The frontend is **NOT** running as a Vite dev server in production. Instead:

1. **Build Phase**: `npm run build` creates optimized static files in `frontend/dist/`
2. **Deploy Phase**: Static files are extracted to `/var/www/html/` 
3. **Serve Phase**: Nginx serves these static files directly (no Node.js process running)

This is the correct production approach for React/Vite applications.

## Manual Testing

### Backend
```bash
# Build and test backend
docker build -t go-image-service:latest .
docker run -d --name go-image-service -p 8080:8080 go-image-service:latest
curl http://localhost:8080/api/health
```

### Frontend
```bash
# Build and test frontend
cd frontend
npm install
npm run build
npm run preview
```

### Full Stack
- Frontend: http://your-server-ip (port 80)
- Backend API: http://your-server-ip:8080
- API through proxy: http://your-server-ip/api/