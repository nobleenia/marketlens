import { Link, NavLink, Outlet } from 'react-router';
import { Home, TrendingUp, GitCompareArrows, Info, Settings, Menu, X, Send } from 'lucide-react';
import { useState } from 'react';

const navLinks = [
  { to: '/', label: 'Prices', icon: <Home className="w-5 h-5" /> },
  { to: '/trends', label: 'Trends', icon: <TrendingUp className="w-5 h-5" /> },
  { to: '/compare', label: 'Compare', icon: <GitCompareArrows className="w-5 h-5" /> },
  { to: '/submit', label: 'Submit', icon: <Send className="w-5 h-5" /> },
  { to: '/about', label: 'About', icon: <Info className="w-5 h-5" /> },
];

export default function PublicLayout() {
  const [menuOpen, setMenuOpen] = useState(false);

  return (
    <div className="min-h-screen flex flex-col">
      {/* ── Top Navigation Bar ─────────────────────────────────── */}
      <header className="bg-white border-b border-gray-200 sticky top-0 z-50">
        <div className="container mx-auto px-4 flex items-center justify-between h-16">
          {/* Logo */}
          <Link to="/" className="flex items-center gap-2">
            <div className="bg-[#1a5f3f] text-white rounded-lg p-1.5">
              <TrendingUp className="w-5 h-5" />
            </div>
            <div className="leading-tight">
              <span className="font-bold text-[#1a5f3f] text-lg">MarketLens</span>
              <span className="block text-[10px] text-gray-500 -mt-1">Price Intelligence</span>
            </div>
          </Link>

          {/* Desktop nav */}
          <nav className="hidden md:flex items-center gap-1">
            {navLinks.map((link) => (
              <NavLink
                key={link.to}
                to={link.to}
                end={link.to === '/'}
                className={({ isActive }) =>
                  `flex items-center gap-2 px-4 py-2 rounded-lg text-sm font-medium transition-colors ${
                    isActive
                      ? 'bg-[#e8f5e9] text-[#1a5f3f]'
                      : 'text-gray-600 hover:bg-gray-100'
                  }`
                }
              >
                {link.icon}
                {link.label}
              </NavLink>
            ))}
          </nav>

          {/* Admin button (desktop) */}
          <div className="hidden md:block">
            <Link
              to="/admin"
              className="flex items-center gap-2 px-4 py-2 border-2 border-[#1a5f3f] text-[#1a5f3f] rounded-lg hover:bg-[#e8f5e9] transition-colors text-sm font-medium"
            >
              <Settings className="w-4 h-4" />
              Admin
            </Link>
          </div>

          {/* Mobile hamburger */}
          <button
            className="md:hidden p-2 hover:bg-gray-100 rounded-lg"
            onClick={() => setMenuOpen(!menuOpen)}
          >
            {menuOpen ? <X className="w-6 h-6" /> : <Menu className="w-6 h-6" />}
          </button>
        </div>

        {/* Mobile menu */}
        {menuOpen && (
          <nav className="md:hidden border-t border-gray-200 bg-white px-4 pb-4 space-y-1">
            {navLinks.map((link) => (
              <NavLink
                key={link.to}
                to={link.to}
                end={link.to === '/'}
                onClick={() => setMenuOpen(false)}
                className={({ isActive }) =>
                  `flex items-center gap-3 px-4 py-3 rounded-lg text-sm font-medium ${
                    isActive
                      ? 'bg-[#e8f5e9] text-[#1a5f3f]'
                      : 'text-gray-600 hover:bg-gray-100'
                  }`
                }
              >
                {link.icon}
                {link.label}
              </NavLink>
            ))}
            <Link
              to="/admin"
              onClick={() => setMenuOpen(false)}
              className="flex items-center gap-3 px-4 py-3 rounded-lg text-sm font-medium border-2 border-[#1a5f3f] text-[#1a5f3f] hover:bg-[#e8f5e9] mt-2"
            >
              <Settings className="w-5 h-5" />
              Admin
            </Link>
          </nav>
        )}
      </header>

      {/* ── Page content ───────────────────────────────────────── */}
      <main className="flex-1">
        <Outlet />
      </main>
    </div>
  );
}
