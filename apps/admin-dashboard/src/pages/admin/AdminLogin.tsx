import { useState } from 'react';
import { Link } from 'react-router-dom';
import { Lock, Eye, EyeOff, Home } from 'lucide-react';

interface Props {
  onLogin: (key: string) => void;
}

export default function AdminLogin({ onLogin }: Props) {
  const [key, setKey] = useState('');
  const [showKey, setShowKey] = useState(false);
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!key.trim()) {
      setError('API key is required');
      return;
    }

    setLoading(true);
    setError('');

    try {
      // Try calling an admin endpoint to validate the key
      const API_BASE = import.meta.env.VITE_API_URL || 'http://localhost:8081';
      const res = await fetch(`${API_BASE}/v1/admin/observations?limit=1`, {
        headers: { 'X-API-Key': key },
      });

      if (res.ok) {
        onLogin(key);
      } else if (res.status === 401) {
        setError('Invalid API key');
      } else {
        // No auth configured on backend (dev mode) — allow through
        onLogin(key);
      }
    } catch {
      // API might not be running — allow through for dev
      onLogin(key);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-b from-[#e8f5e9] to-white flex items-center justify-center p-4 relative">
      <Link
        to="/"
        className="absolute top-6 left-6 flex items-center gap-2 text-[#1a5f3f] hover:text-[#155232] font-medium transition-colors"
      >
        <Home className="w-5 h-5" />
        <span>Back to site</span>
      </Link>

      <div className="bg-white rounded-2xl shadow-xl p-8 w-full max-w-md">
        <div className="text-center mb-8">
          <div className="bg-[#1a5f3f] text-white rounded-xl p-4 w-16 h-16 mx-auto mb-4 flex items-center justify-center">
            <Lock className="w-8 h-8" />
          </div>
          <h1 className="text-2xl font-bold text-gray-900">Admin Access</h1>
          <p className="text-gray-600 mt-2">Enter your API key to continue</p>
        </div>

        <form onSubmit={handleSubmit} className="space-y-6">
          <div>
            <label htmlFor="api-key" className="block text-sm font-medium text-gray-700 mb-2">
              API Key
            </label>
            <div className="relative">
              <input
                id="api-key"
                type={showKey ? 'text' : 'password'}
                value={key}
                onChange={(e) => setKey(e.target.value)}
                placeholder="Enter your admin API key"
                className="w-full px-4 py-3 border-2 border-gray-300 rounded-lg focus:border-[#1a5f3f] focus:outline-none pr-12"
              />
              <button
                type="button"
                onClick={() => setShowKey(!showKey)}
                className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600"
              >
                {showKey ? <EyeOff className="w-5 h-5" /> : <Eye className="w-5 h-5" />}
              </button>
            </div>
          </div>

          {error && (
            <p className="text-red-600 text-sm font-medium">{error}</p>
          )}

          <button
            type="submit"
            disabled={loading}
            className="w-full py-3 bg-[#1a5f3f] text-white rounded-lg font-medium hover:bg-[#155232] transition-colors disabled:opacity-50"
          >
            {loading ? 'Verifying...' : 'Access Admin Panel'}
          </button>
        </form>

        <p className="text-center text-xs text-gray-400 mt-6">
          Contact your system administrator for access credentials.
        </p>
      </div>
    </div>
  );
}