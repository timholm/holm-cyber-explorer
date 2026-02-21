# Quick Deploy Guide

When you're back at your machine, run these 3 commands:

## 1. Rotate your credentials first
- GitHub: https://github.com/settings/tokens → delete old → create new
- Cloudflare: Dashboard → My Profile → API Tokens → roll the key
- Pi: ssh in and run `passwd`

## 2. Create your config file (fill in your NEW credentials)
```bash
cat > ~/docs-framework/.deploy-env << 'EOF'
PI_TARGET=rpi1@192.168.8.197
PI_PUBLIC_IP=YOUR_PUBLIC_IP_HERE
CF_API_TOKEN=YOUR_NEW_CLOUDFLARE_TOKEN
CF_ZONE_ID=YOUR_ZONE_ID
EOF
```

## 3. Run the deploy
```bash
cd ~/docs-framework && ./deploy.sh
```

Total time: ~4 minutes.
