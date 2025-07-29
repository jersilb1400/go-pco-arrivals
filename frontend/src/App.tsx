import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { AuthProvider } from './contexts/AuthContext';
import { AppProvider } from './contexts/AppContext';
import { WebSocketProvider } from './contexts/WebSocketContext';
import Layout from './components/Layout';
import LoginPage from './pages/LoginPage';
import DashboardPage from './pages/DashboardPage';
import BillboardPage from './pages/BillboardPage';
import LocationsPage from './pages/LocationsPage';
import LocationStatusPage from './pages/LocationStatusPage';
import SettingsPage from './pages/SettingsPage';
import AdminPage from './pages/AdminPage';
import ProtectedRoute from './components/ProtectedRoute';
import './index.css';

// Create a client
const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: 1,
      refetchOnWindowFocus: false,
      staleTime: 5 * 60 * 1000, // 5 minutes
    },
  },
});

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <AuthProvider>
        <AppProvider>
          <WebSocketProvider>
            <Router>
              <div className="min-h-screen bg-gray-50">
                <Routes>
                  <Route path="/login" element={<LoginPage />} />
                  <Route
                    path="/"
                    element={
                      <ProtectedRoute>
                        <Layout />
                      </ProtectedRoute>
                    }
                  >
                    <Route index element={<Navigate to="/dashboard" replace />} />
                    <Route path="dashboard" element={<DashboardPage />} />
                                      <Route path="billboard" element={<Navigate to="/billboard/all" replace />} />
                  <Route path="billboard/:locationId" element={<BillboardPage />} />
                  <Route path="admin" element={<AdminPage />} />
                  <Route path="locations" element={<LocationsPage />} />
                  <Route path="location-status" element={<LocationStatusPage />} />
                  <Route path="settings" element={<SettingsPage />} />
                  </Route>
                </Routes>
              </div>
            </Router>
          </WebSocketProvider>
        </AppProvider>
      </AuthProvider>
    </QueryClientProvider>
  );
}

export default App;
