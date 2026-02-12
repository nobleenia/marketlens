import { createBrowserRouter } from 'react-router';

import PublicLayout from './layouts/PublicLayout';
import Home from './pages/Home';
import Trends from './pages/Trends';
import Compare from './pages/Compare';
import About from './pages/About';

import AdminLayout from './pages/admin/AdminLayout';
import AdminDashboard from './pages/admin/AdminDashboard';
import Submissions from './pages/admin/Submissions';
import AuditLogs from './pages/admin/AuditLogs';

export const router = createBrowserRouter([
  {
    element: <PublicLayout />,
    children: [
      { path: '/', element: <Home /> },
      { path: '/trends', element: <Trends /> },
      { path: '/compare', element: <Compare /> },
      { path: '/about', element: <About /> },
    ],
  },
  {
    path: '/admin',
    element: <AdminLayout />,
    children: [
      { index: true, element: <AdminDashboard /> },
      { path: 'submissions', element: <Submissions /> },
      { path: 'audit', element: <AuditLogs /> },
      // placeholders for Sprint 1.10
      // { path: 'prices', element: <AggregatedPrices /> },
      // { path: 'markets', element: <MarketsPage /> },
      // { path: 'crops', element: <CropsPage /> },
    ],
  },
]);
