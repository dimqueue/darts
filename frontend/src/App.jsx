import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { AuthProvider, useAuth } from './contexts/AuthContext';
import { ThemeProvider } from './contexts/ThemeContext';
import { GameProvider } from './contexts/GameContext';
import { ToastProvider } from './contexts/ToastContext';
import ErrorBoundary from './components/ErrorBoundary';
import { ToastContainer } from './components/ui/Toast';

import AuthPage from './pages/AuthPage';
import GamePage from './pages/GamePage';
import GameModesPage from './pages/GameModesPage';
import ProfilePage from './pages/ProfilePage';
import LeaderboardPage from './pages/LeaderboardPage';
import SettingsPage from './pages/SettingsPage';

function ProtectedRoute({ children }) {
    const { isAuthenticated, loading } = useAuth();

    if (loading) {
        return (
            <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-purple-50 to-pink-50">
                <div className="animate-spin w-8 h-8 border-4 border-purple-500 border-t-transparent rounded-full" />
            </div>
        );
    }

    if (!isAuthenticated) {
        return <Navigate to="/auth" replace />;
    }

    return children;
}

function PublicRoute({ children }) {
    const { isAuthenticated, loading } = useAuth();

    if (loading) {
        return (
            <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-purple-50 to-pink-50">
                <div className="animate-spin w-8 h-8 border-4 border-purple-500 border-t-transparent rounded-full" />
            </div>
        );
    }

    if (isAuthenticated) {
        return <Navigate to="/game" replace />;
    }

    return children;
}

function AppRoutes() {
    return (
        <Routes>
            <Route
                path="/auth"
                element={
                    <PublicRoute>
                        <AuthPage />
                    </PublicRoute>
                }
            />

            <Route
                path="/game"
                element={
                    <ProtectedRoute>
                        <GamePage />
                    </ProtectedRoute>
                }
            />
            <Route
                path="/modes"
                element={
                    <ProtectedRoute>
                        <GameModesPage />
                    </ProtectedRoute>
                }
            />
            <Route
                path="/profile"
                element={
                    <ProtectedRoute>
                        <ProfilePage />
                    </ProtectedRoute>
                }
            />
            <Route
                path="/leaderboard"
                element={
                    <ProtectedRoute>
                        <LeaderboardPage />
                    </ProtectedRoute>
                }
            />
            <Route
                path="/settings"
                element={
                    <ProtectedRoute>
                        <SettingsPage />
                    </ProtectedRoute>
                }
            />

            <Route path="/" element={<Navigate to="/game" replace />} />
            <Route path="*" element={<Navigate to="/game" replace />} />
        </Routes>
    );
}

export default function App() {
    return (
        <ErrorBoundary>
            <BrowserRouter basename={import.meta.env.BASE_URL}>
                <AuthProvider>
                    <ThemeProvider>
                        <ToastProvider>
                            <GameProvider>
                                <AppRoutes />
                                <ToastContainer />
                            </GameProvider>
                        </ToastProvider>
                    </ThemeProvider>
                </AuthProvider>
            </BrowserRouter>
        </ErrorBoundary>
    );
}
