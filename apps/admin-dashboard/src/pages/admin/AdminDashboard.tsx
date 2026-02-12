import { BarChart3, FileText, CheckCircle, AlertTriangle, TrendingUp, Loader2 } from 'lucide-react';
import { useCrops, useMarkets, useObservations } from '../../api/hooks';

export default function AdminDashboard() {
  const { data: crops = [], isLoading: cropsLoading } = useCrops();
  const { data: markets = [], isLoading: marketsLoading } = useMarkets();
  const { data: obsData, isLoading: obsLoading } = useObservations();

  const observations = obsData?.data ?? [];

  const stats = {
    totalSubmissions: obsData?.total ?? 0,
    pending: observations.filter((o) => o.status === 'pending').length,
    approved: observations.filter((o) => o.status === 'approved').length,
    flagged: observations.filter((o) => o.status === 'flagged').length,
    totalMarkets: markets.length,
    totalCrops: crops.length,
  };

  const recentSubmissions = observations.slice(0, 5);
  const isLoading = cropsLoading || marketsLoading || obsLoading;

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
            <span className="text-2xl font-bold text-gray-900">
              {isLoading ? '…' : stats.totalSubmissions}
            </span>
          </div>
          <p className="text-gray-600 font-medium">Total Submissions</p>
        </div>

        <div className="bg-white rounded-xl shadow-md p-6 border border-gray-100">
          <div className="flex items-center justify-between mb-3">
            <div className="bg-amber-100 rounded-lg p-3">
              <AlertTriangle className="w-6 h-6 text-amber-600" />
            </div>
            <span className="text-2xl font-bold text-gray-900">
              {isLoading ? '…' : stats.pending}
            </span>
          </div>
          <p className="text-gray-600 font-medium">Pending Review</p>
        </div>

        <div className="bg-white rounded-xl shadow-md p-6 border border-gray-100">
          <div className="flex items-center justify-between mb-3">
            <div className="bg-green-100 rounded-lg p-3">
              <CheckCircle className="w-6 h-6 text-green-600" />
            </div>
            <span className="text-2xl font-bold text-gray-900">
              {isLoading ? '…' : stats.approved}
            </span>
          </div>
          <p className="text-gray-600 font-medium">Approved</p>
        </div>

        <div className="bg-white rounded-xl shadow-md p-6 border border-gray-100">
          <div className="flex items-center justify-between mb-3">
            <div className="bg-red-100 rounded-lg p-3">
              <AlertTriangle className="w-6 h-6 text-red-600" />
            </div>
            <span className="text-2xl font-bold text-gray-900">
              {isLoading ? '…' : stats.flagged}
            </span>
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
        </div>

        {isLoading ? (
          <div className="p-12 text-center">
            <Loader2 className="w-8 h-8 text-[#1a5f3f] mx-auto mb-3 animate-spin" />
            <p className="text-gray-500">Loading...</p>
          </div>
        ) : recentSubmissions.length === 0 ? (
          <div className="p-12 text-center">
            <p className="text-gray-500">No submissions yet.</p>
          </div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead className="bg-gray-50">
                <tr>
                  <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700 uppercase tracking-wider">Crop</th>
                  <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700 uppercase tracking-wider">Market</th>
                  <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700 uppercase tracking-wider">Price</th>
                  <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700 uppercase tracking-wider">Source</th>
                  <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700 uppercase tracking-wider">Status</th>
                  <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700 uppercase tracking-wider">Date</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-200">
                {recentSubmissions.map((obs) => (
                  <tr key={obs.id} className="hover:bg-gray-50">
                    <td className="px-6 py-4 whitespace-nowrap font-medium text-gray-900">
                      {obs.crop_name}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-gray-700">
                      {obs.market_name}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-gray-900 font-semibold">
                      ₦{obs.price.toLocaleString()} / {obs.unit}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-gray-700">
                      {obs.source}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <span
                        className={`px-3 py-1 rounded-full text-xs font-medium ${
                          obs.status === 'approved'
                            ? 'bg-green-100 text-green-700'
                            : obs.status === 'pending'
                            ? 'bg-amber-100 text-amber-700'
                            : obs.status === 'flagged'
                            ? 'bg-red-100 text-red-700'
                            : 'bg-gray-100 text-gray-700'
                        }`}
                      >
                        {obs.status.charAt(0).toUpperCase() + obs.status.slice(1)}
                      </span>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-600">
                      {new Date(obs.observed_at).toLocaleDateString('en-NG', {
                        month: 'short', day: 'numeric', year: 'numeric',
                      })}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </div>
  );
}