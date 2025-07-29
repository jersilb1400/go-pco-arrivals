# MongoDB Atlas Setup Guide

## 1. Create MongoDB Atlas Cluster

1. Go to [MongoDB Atlas](https://cloud.mongodb.com)
2. Sign in to your account
3. Click "Build a Database"
4. Choose "FREE" tier (M0)
5. Select your preferred cloud provider and region
6. Click "Create"

## 2. Configure Database Access

1. Go to "Database Access" in the left sidebar
2. Click "Add New Database User"
3. Choose "Password" authentication
4. Create a username and password
5. Set privileges to "Read and write to any database"
6. Click "Add User"

## 3. Configure Network Access

1. Go to "Network Access" in the left sidebar
2. Click "Add IP Address"
3. For Render deployment, add: `0.0.0.0/0` (allows all IPs)
4. Click "Confirm"

## 4. Get Connection String

1. Go to "Database" in the left sidebar
2. Click "Connect"
3. Choose "Connect your application"
4. Copy the connection string
5. Replace `<password>` with your database user password
6. Replace `<dbname>` with your database name

## 5. Update Render Environment Variables

In your Render dashboard:
1. Go to your service
2. Click "Environment"
3. Add the following variables:
   - `MONGODB_URI`: Your connection string
   - `MONGODB_DATABASE`: Your database name
   - `PCO_CLIENT_ID`: Your PCO Client ID
   - `PCO_CLIENT_SECRET`: Your PCO Client Secret
   - `SESSION_SECRET`: Generated secret
   - `JWT_SECRET`: Generated secret

## 6. Deploy to Render

1. Push your code to GitHub
2. Connect your GitHub repository to Render
3. Render will automatically deploy using the `render.yaml` configuration
4. Your app will be available at your Render URL

## 7. Test the Deployment

1. Visit your Render URL
2. Test the login functionality
3. Verify WebSocket connections work
4. Check that data is being stored in MongoDB Atlas

## Troubleshooting

- **Connection Issues**: Ensure your MongoDB Atlas cluster is in the same region as your Render service
- **Authentication Errors**: Double-check your database username and password
- **Network Access**: Make sure you've added `0.0.0.0/0` to allowed IP addresses
- **Environment Variables**: Verify all required environment variables are set in Render
