import { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { useAuthStore } from '@/stores/useAuthStore';
import { useUIStore } from '@/stores/useUIStore';
import { ApiService } from '@/services/api';

export default function LoginPage() {
  const navigate = useNavigate();
  const setUser = useAuthStore((state) => state.setUser);
  const addToast = useUIStore((state) => state.addToast);
  const [isLoading, setIsLoading] = useState(false);
  const [isLogin, setIsLogin] = useState(true);
  const [form, setForm] = useState({ email: '', password: '', username: '' });

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);

    try {
      if (isLogin) {
        const response = await ApiService.login(form.email, form.password);
        setUser(response.user, response.token);
        addToast('Login successful!', 'success');
        navigate('/');
      } else {
        const response = await ApiService.register(form.username, form.email, form.password);
        setUser(response.user, response.token);
        addToast('Registration successful!', 'success');
        navigate('/');
      }
    } catch (error) {
      addToast(error instanceof Error ? error.message : 'An error occurred', 'error');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-900">
      <div className="w-full max-w-md">
        <div className="bg-gray-800 rounded-lg shadow-xl p-8">
          <h1 className="text-3xl font-bold text-center mb-8">Voice Chat</h1>

          <form onSubmit={handleSubmit} className="space-y-4">
            {!isLogin && (
              <input
                type="text"
                placeholder="Username"
                value={form.username}
                onChange={(e) => setForm({ ...form, username: e.target.value })}
                className="w-full px-4 py-2 bg-gray-700 rounded border border-gray-600 focus:border-blue-500 outline-none"
                required
              />
            )}

            <input
              type="email"
              placeholder="Email"
              value={form.email}
              onChange={(e) => setForm({ ...form, email: e.target.value })}
              className="w-full px-4 py-2 bg-gray-700 rounded border border-gray-600 focus:border-blue-500 outline-none"
              required
            />

            <input
              type="password"
              placeholder="Password"
              value={form.password}
              onChange={(e) => setForm({ ...form, password: e.target.value })}
              className="w-full px-4 py-2 bg-gray-700 rounded border border-gray-600 focus:border-blue-500 outline-none"
              required
            />

            <button
              type="submit"
              disabled={isLoading}
              className="w-full px-4 py-2 bg-blue-600 hover:bg-blue-700 rounded font-semibold disabled:opacity-50"
            >
              {isLoading ? 'Loading...' : isLogin ? 'Login' : 'Register'}
            </button>
          </form>

          <p className="text-center mt-4 text-gray-400">
            {isLogin ? "Don't have an account? " : 'Already have an account? '}
            <button
              onClick={() => {
                setIsLogin(!isLogin);
                setForm({ email: '', password: '', username: '' });
              }}
              className="text-blue-400 hover:text-blue-300"
            >
              {isLogin ? 'Register' : 'Login'}
            </button>
          </p>
        </div>
      </div>
    </div>
  );
}
