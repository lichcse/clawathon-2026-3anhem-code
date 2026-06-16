import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { useAuthStore } from '@/stores/useAuthStore';
import { useUIStore } from '@/stores/useUIStore';
import LoginPage from '@/pages/LoginPage';
import LandingPage from '@/pages/LandingPage';
import RoomPage from '@/pages/RoomPage';
import ToastContainer from '@/components/common/ToastContainer';

function App() {
  const isAuthenticated = useAuthStore((state) => state.isAuthenticated);
  const { toasts } = useUIStore();

  return (
    <BrowserRouter>
      <div className="min-h-screen bg-gray-900 text-gray-100">
        <Routes>
          <Route
            path="/login"
            element={!isAuthenticated ? <LoginPage /> : <Navigate to="/" />}
          />
          <Route
            path="/"
            element={isAuthenticated ? <LandingPage /> : <Navigate to="/login" />}
          />
          <Route
            path="/rooms/:roomId"
            element={isAuthenticated ? <RoomPage /> : <Navigate to="/login" />}
          />
          <Route path="*" element={<Navigate to="/" />} />
        </Routes>
      </div>
      <ToastContainer toasts={toasts} />
    </BrowserRouter>
  );
}

export default App;
