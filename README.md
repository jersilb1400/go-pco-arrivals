# PCO Arrivals Billboard

A professional, real-time check-in display system for Planning Center Online, built with Go and React.

## 🎯 Overview

The PCO Arrivals Billboard is a modern web application that displays real-time check-ins from Planning Center Online in an attractive, professional format. Perfect for churches and organizations that want to show arrival information on screens or billboards.

## ✨ Features

### 🎨 Frontend (React)
- **Real-time Updates**: WebSocket-powered live check-in displays
- **Responsive Design**: Works on desktop, tablet, and mobile devices
- **PWA Support**: Installable as a standalone app with offline capabilities
- **Modern UI/UX**: Beautiful interface with smooth animations
- **OAuth Authentication**: Secure login through Planning Center Online
- **Location Management**: Add and manage multiple billboard locations
- **Professional Styling**: Built with Tailwind CSS and Framer Motion

### ⚡ Backend (Go)
- **High Performance**: Built with Go and Fiber for optimal performance
- **Real-time Communication**: WebSocket support for live updates
- **OAuth Integration**: Secure authentication with Planning Center Online
- **Database Support**: SQLite (development) and PostgreSQL (production)
- **Location-based Filtering**: Separate billboards for different locations
- **Session Management**: 30-day "Remember Me" functionality
- **API Caching**: Intelligent caching for better performance
- **Health Monitoring**: Built-in health checks and monitoring

## 🏗 Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Frontend      │    │   Backend       │    │   PCO API       │
│   (React)       │◄──►│   (Go/Fiber)    │◄──►│   (OAuth)       │
│                 │    │                 │    │                 │
│ • Real-time UI  │    │ • REST API      │    │ • Check-ins     │
│ • WebSocket     │    │ • WebSocket     │    │ • User Data     │
│ • PWA Support   │    │ • OAuth         │    │ • Locations     │
│ • Responsive    │    │ • Database      │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## 🚀 Quick Start

### Prerequisites

- **Go 1.21+** for the backend
- **Node.js 18+** for the frontend
- **Planning Center Online** account with API access
- **PCO OAuth App** credentials

### 1. Clone the Repository

```bash
git clone <repository-url>
cd go-pco-arrivals-dashboard
```

### 2. Backend Setup

```bash
# Navigate to backend directory
cd backend

# Install dependencies
go mod download

# Set up environment variables
cp .env.example .env
# Edit .env with your PCO credentials

# Run database migrations
go run cmd/migrate/main.go

# Start the server
go run cmd/server/main.go
```

### 3. Frontend Setup

```bash
# Navigate to frontend directory
cd frontend

# Install dependencies
npm install

# Set up environment variables
cp env.example .env.local
# Edit .env.local with your API URL

# Start development server
npm run dev
```

### 4. Access the Application

- **Frontend**: http://localhost:5173
- **Backend API**: http://localhost:3000
- **Health Check**: http://localhost:3000/health

## 🔧 Configuration

### Backend Environment Variables

```env
# Server Configuration
SERVER_PORT=3000
SERVER_HOST=localhost
SERVER_TRUST_PROXY=false

# Database Configuration
DB_TYPE=sqlite
DB_DSN=./data/billboard.db
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=5
DB_CONN_MAX_LIFETIME=300s

# PCO Configuration
PCO_CLIENT_ID=your_pco_client_id
PCO_CLIENT_SECRET=your_pco_client_secret
PCO_REDIRECT_URI=http://localhost:3000/auth/callback
PCO_AUTH_URL=https://api.planningcenteronline.com/oauth/authorize
PCO_TOKEN_URL=https://api.planningcenteronline.com/oauth/token
PCO_API_BASE_URL=https://api.planningcenteronline.com

# Authentication Configuration
AUTH_SESSION_SECRET=your_session_secret
AUTH_REMEMBER_ME_DAYS=30
AUTH_TOKEN_REFRESH_THRESHOLD=300s

# Real-time Configuration
REALTIME_ENABLED=true
REALTIME_HEARTBEAT_INTERVAL=30s
REALTIME_CONNECTION_TIMEOUT=60s
```

### Frontend Environment Variables

```env
# API Configuration
VITE_API_URL=http://localhost:3000
VITE_WS_URL=ws://localhost:3000

# Application Configuration
VITE_APP_NAME=PCO Arrivals Billboard
VITE_APP_VERSION=1.0.0

# Feature Flags
VITE_ENABLE_PWA=true
VITE_ENABLE_ANALYTICS=false
```

## 📁 Project Structure

```
go-pco-arrivals-dashboard/
├── backend/                 # Go backend application
│   ├── cmd/                # Application entry points
│   ├── internal/           # Private application code
│   │   ├── config/         # Configuration management
│   │   ├── database/       # Database connection and migrations
│   │   ├── handlers/       # HTTP request handlers
│   │   ├── middleware/     # HTTP middleware
│   │   ├── models/         # Database models
│   │   ├── services/       # Business logic services
│   │   └── utils/          # Utility functions
│   ├── web/                # Static web assets
│   ├── go.mod              # Go module file
│   └── go.sum              # Go dependencies checksum
├── frontend/               # React frontend application
│   ├── src/                # Source code
│   │   ├── components/     # Reusable UI components
│   │   ├── contexts/       # React contexts
│   │   ├── pages/          # Page components
│   │   ├── services/       # API and WebSocket services
│   │   └── types/          # TypeScript type definitions
│   ├── public/             # Static assets
│   ├── package.json        # Node.js dependencies
│   └── README.md           # Frontend documentation
├── scripts/                # Deployment and utility scripts
├── docs/                   # Documentation
└── README.md               # This file
```

## 🔌 API Endpoints

### Authentication
- `GET /auth/login` - Initiate OAuth login
- `GET /auth/callback` - OAuth callback handler
- `GET /auth/status` - Get authentication status
- `POST /auth/logout` - Logout user
- `POST /auth/refresh` - Refresh access token

### Billboard
- `GET /billboard/state/:locationID` - Get billboard state
- `GET /billboard/check-ins/:locationID` - Get recent check-ins
- `GET /billboard/stats/:locationID` - Get check-in statistics
- `POST /billboard/sync/:locationID` - Sync PCO check-ins
- `GET /billboard/locations` - Get all locations
- `POST /billboard/locations` - Add new location

### Health
- `GET /health` - Basic health check
- `GET /health/detailed` - Detailed system status

### WebSocket
- `GET /ws` - WebSocket connection for real-time updates
- `GET /ws/billboard/:locationID` - Location-specific WebSocket

## 🚀 Deployment

### Backend Deployment (Render)

1. **Connect to Render**:
   - Link your GitHub repository
   - Set environment variables
   - Deploy automatically

2. **Environment Variables**:
   ```env
   DB_TYPE=postgres
   DB_DSN=postgres://user:pass@host:port/db
   PCO_CLIENT_ID=your_client_id
   PCO_CLIENT_SECRET=your_client_secret
   PCO_REDIRECT_URI=https://your-domain.com/auth/callback
   ```

### Frontend Deployment

1. **Build the application**:
   ```bash
   cd frontend
   npm run build
   ```

2. **Deploy to static hosting**:
   - Netlify, Vercel, or similar
   - Set environment variables for production
   - Configure custom domain

### Docker Deployment

```bash
# Build and run with Docker Compose
docker-compose up -d
```

## 🔒 Security

- **OAuth 2.0**: Secure authentication through PCO
- **HTTPS Only**: All communications encrypted
- **Session Management**: Secure session handling
- **Input Validation**: Server-side validation
- **Rate Limiting**: Protection against abuse
- **CORS Configuration**: Proper cross-origin settings

## 🧪 Testing

### Backend Tests

```bash
cd backend
go test ./...
go test -bench=. ./...
```

### Frontend Tests

```bash
cd frontend
npm test
npm run test:coverage
```

## 📊 Monitoring

- **Health Checks**: Built-in endpoint monitoring
- **Logging**: Structured logging with levels
- **Metrics**: Performance monitoring
- **Error Tracking**: Comprehensive error handling

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License.

## 🆘 Support

For support and questions:

- Check the documentation
- Review the API documentation
- Contact the development team

## 🙏 Acknowledgments

- **Planning Center Online** for their excellent API
- **Grace Fellowship** for the opportunity to build this system
- **Open Source Community** for the amazing tools and libraries

---

Built with ❤️ for Grace Fellowship 