import { BarChart3, FileText, CheckCircle, AlertTriangle, TrendingUp, Loader2 } from 'lucide-react';
import { useCrops, useMarkets } from '../../api/hooks';
import { priceSubmissions } from '../../data/mockData';

export default function AdminDashboard() {
  const { data: crops = [], isLoading: cropsLoading } = useCrops();
  const { data: markets = [], isLoading: marketsLoading } = useMarkets();

  const stats = {
    totalSubmissions: priceSubmissions.length,
    pending: priceSubmissions.filter((s) => s.status === 'pending').length,
    approved: priceSubmissions.filter((s) => s.status === 'approved').length,
    flagged: priceSubmissions.filter((s) => s.status === 'flagged').length,
    totalMarkets: markets.length,
    totalCrops: crops.length,
  };

  const recentSubmissions = priceSubmissions.slice(0, 5);

  const isLoading = cropsLoading || marketsLoading;

  return (
    <div className="p-6 md:p-8">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900 mb-2">Admin Dashboard</h1>
        <p className="text-gray-600">Monitor and manage price submissions</p>
      </div>

      {/* Stats Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        <div className="bg-white rounded-xl shadow-md p-6 border border-gray-100">
          <div className="flex items-center justify-between mb-3">
            <div className="bg-blue-100 rounded-lg p-3">
              <FileText className="w-6 h-6 text-blue-600" />
            </div>
            <span className="text-2xl font-bold text-gray-900">{stats.totalSubmissions}</span>
          </div>
          <p className="text-gray-600 font-medium">Total Submissions</p>
        </div>

        <div className="bg-white rounded-xl shadow-md p-6 border border-gray-100">
          <div className="flex items-center justify-between mb-3">
            <div className="bg-amber-100 rounded-lg p-3">
              <AlertTriangle className="w-6 h-6 text-amber-600" />
            </div>
            <span className="text-2xl font-bold text-gray-900">{stats.pending}</span>
          </div>
          <p className="text-gray-600 font-medium">Pending Review</p>
        </div>

        <div className="bg-white rounded-xl shadow-md p-6 border border-gray-100">
          <div className="flex items-center justify-between mb-3">
            <div className="bg-green-100 rounded-lg p-3">
              <CheckCircle className="w-6 h-6 text-green-600" />
            </div>
            <span className="text-2xl font-bold text-gray-900">{stats.approved}</span>
          </div>
          <p className="text-gray-600 font-medium">Approved</p>
        </div>

        <div className="bg-white rounded-xl shadow-md p-6 border border-gray-100">
          <div className="flex items-center justify-between mb-3">
            <div className="bg-red-100 rounded-lg p-3">
              <AlertTriangle className="w-6 h-6 text-red-600" />
            </div>
            <span className="text-2xl font-bold text-gray-900">{stats.flagged}</span>
          </div>
          <p className="text-gray-600 font-medium">Flagged</p>
        </div>
      </div>

      {/* Quick Stats */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-8">
        <div className="bg-gradient-to-br from-[#1a5f3f] to-[#2d8659] text-white rounded-xl shadow-lg p-6">
          <div className="flex items-center gap-4">
            <div className="bg-white/20 rounded-lg p-3">
              {isLoading ? (
                <Loader2 className="w-8 h-8 animate-spin" />
              ) : (
                <BarChart3 className="w-8 h-8" />
              )}
            </div>
            <div>
              <p className="text-white/80 mb-1">Active Markets</p>
              <p className="text-4xl font-bold">{isLoading ? '…' : stats.totalMarkets}</p>
            </div>
          </div>
        </div>

        <div className="bg-gradient-to-br from-[#2d8659] to-[#3a9f6f] text-white rounded-xl shadow-lg p-6">
          <div className="flex items-center gap-4">
            <div className="bg-white/20 rounded-lg p-3">
              {isLoading ? (
                <Loader2 className="w-8 h-8 animate-spin" />
              ) : (
                <TrendingUp className="w-8 h-8" />
              )}
            </div>
            <div>
              <p className="text-white/80 mb-1">Tracked Crops</p>
              <p className="text-4xl font-bold">{isLoading ? '…' : stats.totalCrops}</p>
            </div>
          </div>
        </div>
      </div>

      {/* Recent Submissions */}
      <div className="bg-white rounded-xl shadow-lg overflow-hidden">
        <div className="p-6 border-b border-gray-200">
          <h2 className="text-xl font-bold text-gray-900">Recent Submissions</h2>
          <p className="text-sm text-gray-500 mt-1">Demo data — live feed coming in Sprint 1.10</p>
        </div>

        <div className="overflow-x-auto">
          <table className="w-full">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700 uppercase tracking-wider">
                  Crop
                </th>
                <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700 uppercase tracking-wider">
                  Market
                </th>
                <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700 uppercase tracking-wider">
                  Price Range
                </th>
                <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700 uppercase tracking-wider">
                  Source
                </th>
                <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700 uppercase tracking-wider">
                  Status
                </th>
                <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700 uppercase tracking-wider">
                  Time
                </th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-200">
              {recentSubmissions.map((submission) => (
                <tr key={submission.id} className="hover:bg-gray-50">
                  <td className="px-6 py-4 whitespace-nowrap font-medium text-gray-900">
                    {submission.crop}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-gray-700">
                    {submission.market}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-gray-900 font-semibold">
                    ₦{submission.minPrice.toLocaleString()} - ₦{submission.maxPrice.toLocaleString()}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-gray-700">
                    {submission.source}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <span
                      className={`px-3 py-1 rounded-full text-xs font-medium ${
                        submission.status === 'approved'
                          ? 'bg-green-100 text-green-700'
                          : submission.status === 'pending'
                          ? 'bg-amber-100 text-amber-700'
                          : 'bg-red-100 text-red-700'
                      }`}
                    >
                      {submission.status.charAt(0).toUpperCase() + submission.status.slice(1)}
                    </span>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-600">
                    {submission.timestamp}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}