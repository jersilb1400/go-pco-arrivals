# Render Deployment Checklist

## ✅ Pre-Deployment Checklist

### MongoDB Atlas Setup
- [ ] MongoDB Atlas cluster created
- [ ] Database user created with read/write permissions
- [ ] Network access configured (0.0.0.0/0)
- [ ] Connection string obtained and tested

### PCO Configuration
- [ ] PCO application created in Planning Center
- [ ] Client ID and Secret obtained
- [ ] Redirect URI configured: `go-pco-arrivals/auth/callback`
- [ ] Required scopes configured

### Render Configuration
- [ ] Render account created
- [ ] GitHub repository connected to Render
- [ ] Environment variables configured in Render dashboard
- [ ] Auto-deploy enabled

### Code Preparation
- [ ] All changes committed to main branch
- [ ] render.yaml file present in repository
- [ ] Dockerfile.render present in repository
- [ ] .env.render file created (for reference)

## 🚀 Deployment Steps

1. **Push to GitHub**
   ```bash
   git add .
   git commit -m "Configure for Render deployment"
   git push origin main
   ```

2. **Deploy on Render**
   - Render will automatically detect the render.yaml file
   - Build will start automatically
   - Monitor the build logs for any issues

3. **Configure Environment Variables**
   In Render dashboard, set these variables:
   - `MONGODB_URI`: Your MongoDB Atlas connection string
   - `MONGODB_DATABASE`: Your database name
   - `PCO_CLIENT_ID`: Your PCO Client ID
   - `PCO_CLIENT_SECRET`: Your PCO Client Secret
   - `SESSION_SECRET`: YIVWPwOmXOr8UY3cwXD7eJeVGFiL/vuV6/gLdrOBnzA=
   - `JWT_SECRET`: 88Os98iSRtVVtHnAYHvPt6c9eatogKIK2BVGum4l74M=

4. **Test the Deployment**
   - Visit your Render URL
   - Test login functionality
   - Verify WebSocket connections
   - Check MongoDB Atlas for data storage

## 🔧 Post-Deployment Verification

- [ ] Application loads without errors
- [ ] PCO login works correctly
- [ ] WebSocket connections establish
- [ ] Real-time updates work
- [ ] Data is being stored in MongoDB Atlas
- [ ] Health check endpoint responds
- [ ] SSL/HTTPS is working

## 📊 Monitoring

- Monitor Render logs for any errors
- Check MongoDB Atlas metrics
- Verify application performance
- Test all features thoroughly

## 🆘 Troubleshooting

### Common Issues:
1. **Build Failures**: Check render.yaml syntax
2. **Connection Errors**: Verify MongoDB Atlas settings
3. **Authentication Issues**: Check PCO credentials
4. **Environment Variables**: Ensure all required vars are set

### Support Resources:
- Render Documentation: https://render.com/docs
- MongoDB Atlas Documentation: https://docs.atlas.mongodb.com
- PCO API Documentation: https://developer.planningcenteronline.com
