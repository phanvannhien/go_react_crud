import { BrowserRouter, Routes, Route, Navigate, Link, Outlet } from 'react-router-dom';
import { AuthProvider, useAuth } from './lib/auth-context';
import LoginPage from './features/auth/pages/LoginPage';
import RegisterPage from './features/auth/pages/RegisterPage';
import UserListPage from './features/user/pages/UserListPage';
import ProductListPage from './features/product/pages/ProductListPage';

function ProtectedRoute() {
  const { isAuthenticated } = useAuth();
  if (!isAuthenticated) return <Navigate to="/login" replace />;
  return <Outlet />;
}

function Layout() {
  const { user, logout, isAuthenticated } = useAuth();

  if (!isAuthenticated) return <Outlet />;

  return (
    <div className="min-h-screen">
      <nav className="bg-[var(--color-surface)] border-b border-[var(--color-border)] px-6 py-3">
        <div className="max-w-6xl mx-auto flex items-center justify-between">
          <div className="flex items-center gap-6">
            <span className="text-xl font-bold bg-gradient-to-r from-indigo-400 to-purple-400 bg-clip-text text-transparent">
              Dashboard
            </span>
            <div className="flex gap-4">
              <Link to="/users" className="text-sm text-[var(--color-text-muted)] hover:text-[var(--color-text)] transition-colors">Users</Link>
              <Link to="/products" className="text-sm text-[var(--color-text-muted)] hover:text-[var(--color-text)] transition-colors">Products</Link>
            </div>
          </div>
          <div className="flex items-center gap-4">
            <span className="text-sm text-[var(--color-text-muted)]">{user?.email}</span>
            <span className={`text-xs px-2 py-0.5 rounded-full ${user?.role === 'admin' ? 'bg-purple-500/20 text-purple-400' : 'bg-blue-500/20 text-blue-400'
              }`}>
              {user?.role}
            </span>
            <button
              onClick={logout}
              className="text-sm text-red-400 hover:text-red-300 transition-colors"
            >
              Logout
            </button>
          </div>
        </div>
      </nav>
      <Outlet />
    </div>
  );
}

export default function App() {
  return (
    <AuthProvider>
      <BrowserRouter>
        <Routes>
          <Route element={<Layout />}>
            {/* Public routes */}
            <Route path="/login" element={<LoginPage />} />
            <Route path="/register" element={<RegisterPage />} />

            {/* Protected routes */}
            <Route element={<ProtectedRoute />}>
              <Route path="/" element={<Navigate to="/users" replace />} />
              <Route path="/users" element={<UserListPage />} />
              <Route path="/products" element={<ProductListPage />} />
            </Route>
          </Route>
        </Routes>
      </BrowserRouter>
    </AuthProvider>
  );
}
