import React, { useState } from 'react';
import { Navigate } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';
import { Users, Shield, Zap } from 'lucide-react';

const LoginPage: React.FC = () => {
  const { isAuthenticated, login } = useAuth();
  const [rememberMe, setRememberMe] = useState(false);

  if (isAuthenticated) {
    return <Navigate to="/dashboard" replace />;
  }

  const handleLogin = () => {
    login(rememberMe);
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-md w-full space-y-8">
        <div>
          <div className="mx-auto h-12 w-12 flex items-center justify-center rounded-full bg-primary-100">
            <Users className="h-6 w-6 text-primary-600" />
          </div>
          <h2 className="mt-6 text-center text-3xl font-extrabold text-gray-900">
            PCO Arrivals Billboard
          </h2>
          <p className="mt-2 text-center text-sm text-gray-600">
            Sign in with your Planning Center Online account
          </p>
        </div>
        
        <div className="card p-8">
          <div className="space-y-6">
            <div>
              <h3 className="text-lg font-medium text-gray-900 mb-4">
                Welcome to the Arrivals Billboard
              </h3>
              <p className="text-sm text-gray-600 mb-6">
                This application provides real-time check-in updates and billboard displays for your church events.
              </p>
            </div>

            <div className="space-y-4">
              <div className="flex items-center space-x-3">
                <div className="flex-shrink-0">
                  <Shield className="h-5 w-5 text-success-600" />
                </div>
                <div>
                  <p className="text-sm font-medium text-gray-900">Secure OAuth Login</p>
                  <p className="text-xs text-gray-500">Authenticated through Planning Center Online</p>
                </div>
              </div>

              <div className="flex items-center space-x-3">
                <div className="flex-shrink-0">
                  <Zap className="h-5 w-5 text-warning-600" />
                </div>
                <div>
                  <p className="text-sm font-medium text-gray-900">Real-time Updates</p>
                  <p className="text-xs text-gray-500">Live check-in notifications and billboard displays</p>
                </div>
              </div>
            </div>

            <div className="flex items-center">
              <input
                id="remember-me"
                name="remember-me"
                type="checkbox"
                checked={rememberMe}
                onChange={(e) => setRememberMe(e.target.checked)}
                className="h-4 w-4 text-primary-600 focus:ring-primary-500 border-gray-300 rounded"
              />
              <label htmlFor="remember-me" className="ml-2 block text-sm text-gray-900">
                Remember me for 30 days
              </label>
            </div>

            <div>
              <button
                onClick={handleLogin}
                className="group relative w-full flex justify-center py-3 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 transition-colors duration-200"
              >
                <span className="absolute left-0 inset-y-0 flex items-center pl-3">
                  <Users className="h-5 w-5 text-primary-500 group-hover:text-primary-400" />
                </span>
                Sign in with PCO
              </button>
            </div>

            <div className="text-center">
              <p className="text-xs text-gray-500">
                By signing in, you agree to our terms of service and privacy policy.
              </p>
            </div>
          </div>
        </div>

        <div className="text-center">
          <p className="text-xs text-gray-500">
            Need help? Contact your system administrator
          </p>
        </div>
      </div>
    </div>
  );
};

export default LoginPage; 