import { Link, Outlet, useLocation } from 'react-router';
import {
  LayoutDashboard,
  FileText,
  CheckSquare,
  Store,
  Sprout,
  Menu,
  X,
  ScrollText,
} from 'lucide-react';
import { useState } from 'react';

const navItems = [
  { to: '/admin', label: 'Dashboard', icon: <LayoutDashboard className="w-5 h-5" /> },
  { to: '/admin/submissions', label: 'Price Submissions', icon: <FileText className="w-5 h-5" /> },
  { to: '/admin/audit', label: 'Audit Log', icon: <ScrollText className="w-5 h-5" /> },
  { to: '/admin/prices', label: 'Aggregated Prices', icon: <CheckSquare className="w-5 h-5" /> },
  { to: '/admin/markets', label: 'Markets', icon: <Store className="w-5 h-5" /> },
  { to: '/admin/crops', label: 'Crops', icon: <Sprout className="w-5 h-5" /> },
];

export default function AdminLayout() {
  const location = useLocation();
  const [sidebarOpen, setSidebarOpen] = useState(false);
  
  return (
    <div className="min-h-screen bg-gray-100">
      {/* Mobile Header */}
      <div className="lg:hidden bg-white border-b border-gray-200 px-4 py-3 flex items-center justify-between sticky top-0 z-50">
        <Link to="/" className="flex items-center gap-2">
          <div className="bg-[#1a5f3f] text-white rounded-lg p-2">
            <LayoutDashboard className="w-5 h-5" />
          </div>
          <span className="font-bold text-[#1a5f3f]">MarketLens Admin</span>
        </Link>
        <button
          onClick={() => setSidebarOpen(!sidebarOpen)}
          className="p-2 hover:bg-gray-100 rounded-lg"
        >
          {sidebarOpen ? <X className="w-6 h-6" /> : <Menu className="w-6 h-6" />}
        </button>
      </div>
      
      <div className="flex">
        {/* Sidebar */}
        <aside
          className={`
            fixed lg:sticky top-0 left-0 h-screen bg-white border-r border-gray-200 z-40
            transition-transform duration-300 ease-in-out
            ${sidebarOpen ? 'translate-x-0' : '-translate-x-full lg:translate-x-0'}
            w-64
          `}
        >
          {/* Desktop Logo */}
          <div className="hidden lg:block p-6 border-b border-gray-200">
            <Link to="/" className="flex items-center gap-2">
              <div className="bg-[#1a5f3f] text-white rounded-lg p-2">
                <LayoutDashboard className="w-6 h-6" />
              </div>
              <div>
                <h1 className="font-bold text-[#1a5f3f]">MarketLens</h1>
                <p className="text-xs text-gray-600">Admin Panel</p>
              </div>
            </Link>
          </div>
          
          {/* Navigation */}
          <nav className="p-4 space-y-2">
            {navItems.map((item) => {
              const isActive = location.pathname === item.to;
              return (
                <Link
                  key={item.to}
                  to={item.to}
                  onClick={() => setSidebarOpen(false)}
                  className={`
                    flex items-center gap-3 px-4 py-3 rounded-lg transition-colors
                    ${
                      isActive
                        ? 'bg-[#1a5f3f] text-white'
                        : 'text-gray-700 hover:bg-gray-100'
                    }
                  `}
                >
                  {item.icon}
                  <span className="font-medium">{item.label}</span>
                </Link>
              );
            })}
          </nav>
          
          {/* Back to Public Site */}
          <div className="absolute bottom-0 left-0 right-0 p-4 border-t border-gray-200">
            <Link
              to="/"
              className="flex items-center justify-center gap-2 px-4 py-3 text-[#1a5f3f] border-2 border-[#1a5f3f] rounded-lg hover:bg-[#e8f5e9] transition-colors font-medium"
            >
              ← Back to Public Site
            </Link>
          </div>
        </aside>
        
        {/* Mobile Overlay */}
        {sidebarOpen && (
          <div
            className="fixed inset-0 bg-black bg-opacity-50 z-30 lg:hidden"
            onClick={() => setSidebarOpen(false)}
          />
        )}
        
        {/* Main Content */}
        <main className="flex-1 min-h-screen lg:h-screen overflow-auto">
          <Outlet />
        </main>
      </div>
    </div>
  );
}
