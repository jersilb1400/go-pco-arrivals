# PCO Arrivals Billboard - Frontend

A modern, professional React application for displaying real-time check-ins from Planning Center Online.

## 🚀 Features

- **Real-time Updates**: WebSocket-powered live check-in displays
- **Responsive Design**: Works on desktop, tablet, and mobile devices
- **PWA Support**: Installable as a standalone app with offline capabilities
- **Modern UI/UX**: Beautiful interface with smooth animations
- **OAuth Authentication**: Secure login through Planning Center Online
- **Location Management**: Add and manage multiple billboard locations
- **Professional Styling**: Built with Tailwind CSS and Framer Motion

## 🛠 Tech Stack

- **React 18** with TypeScript
- **Vite** for fast development and building
- **Tailwind CSS** for styling
- **Framer Motion** for animations
- **React Query** for server state management
- **React Router** for navigation
- **Lucide React** for icons
- **WebSocket** for real-time updates

## 📦 Installation

1. **Install dependencies**:
   ```bash
   npm install
   ```

2. **Set up environment variables**:
   ```bash
   cp env.example .env.local
   ```
   
   Edit `.env.local` with your configuration:
   ```env
   VITE_API_URL=http://localhost:3000
   VITE_WS_URL=ws://localhost:3000
   ```

3. **Start development server**:
   ```bash
   npm run dev
   ```

## 🏗 Development

### Available Scripts

- `npm run dev` - Start development server
- `npm run build` - Build for production
- `npm run preview` - Preview production build
- `npm run lint` - Run ESLint
- `npm run type-check` - Run TypeScript type checking

### Project Structure

```
src/
├── components/          # Reusable UI components
├── contexts/           # React contexts (Auth, WebSocket)
├── pages/              # Page components
├── services/           # API and WebSocket services
├── types/              # TypeScript type definitions
├── App.tsx             # Main app component
├── main.tsx            # App entry point
└── index.css           # Global styles
```

### Key Components

- **Layout**: Main application layout with sidebar navigation
- **AuthContext**: Authentication state management
- **WebSocketContext**: Real-time connection management
- **BillboardPage**: Real-time check-in display
- **LocationsPage**: Location management interface
- **SettingsPage**: User settings and preferences

## 🎨 Styling

The application uses Tailwind CSS with a custom design system:

- **Primary Colors**: Blue theme for professional appearance
- **Success Colors**: Green for positive actions
- **Warning Colors**: Orange for alerts
- **Danger Colors**: Red for errors
- **Custom Components**: Reusable button and input styles

## 🔌 API Integration

The frontend communicates with the Go backend through:

- **REST API**: For CRUD operations
- **WebSocket**: For real-time updates
- **OAuth**: For authentication

### API Endpoints

- `/auth/*` - Authentication endpoints
- `/billboard/*` - Billboard and check-in data
- `/api/*` - General API endpoints
- `/health` - Health check endpoints

## 📱 PWA Features

The application is a Progressive Web App with:

- **Offline Support**: Cached resources for offline viewing
- **Installable**: Can be installed on desktop and mobile
- **App-like Experience**: Full-screen mode and native feel
- **Background Sync**: Automatic data synchronization

## 🔒 Security

- **OAuth 2.0**: Secure authentication through PCO
- **HTTPS Only**: All communications encrypted
- **Session Management**: Secure session handling
- **Input Validation**: Client-side validation with server verification

## 🚀 Deployment

### Build for Production

```bash
npm run build
```

The build output will be in the `dist/` directory.

### Environment Configuration

Set these environment variables for production:

```env
VITE_API_URL=https://your-api-domain.com
VITE_WS_URL=wss://your-api-domain.com
VITE_APP_NAME=PCO Arrivals Billboard
```

### Deployment Options

- **Static Hosting**: Deploy to Netlify, Vercel, or similar
- **CDN**: Use Cloudflare or AWS CloudFront
- **Container**: Docker deployment with nginx

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License.

## 🆘 Support

For support and questions:

- Check the documentation
- Review the backend API documentation
- Contact the development team

---

Built with ❤️ for Grace Fellowship
