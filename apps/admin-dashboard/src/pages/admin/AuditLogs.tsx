import { useState } from 'react';
import { ScrollText, Loader2 } from 'lucide-react';
import { useAuditLogs } from '../../api/hooks';

export default function AuditLogs() {
  const [entityType, setEntityType] = useState('');
  const { data, isLoading, isError } = useAuditLogs(entityType || undefined);

  const logs = data?.data ?? [];

  return (
    <div className="p-6 md:p-8">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900 mb-2">Audit Log</h1>
        <p className="text-gray-600">All admin actions are recorded here</p>
      </div>

      {/* Filter */}
      <div className="bg-white rounded-xl shadow-md mb-6 p-4 flex items-center gap-4 flex-wrap">
        <label htmlFor="entity-filter" className="font-medium text-gray-700">Entity:</label>
        <select
          id="entity-filter"
          value={entityType}
          onChange={(e) => setEntityType(e.target.value)}
          className="px-4 py-2 border-2 border-gray-300 rounded-lg focus:border-[#1a5f3f] focus:outline-none"
        >
          <option value="">All</option>
          <option value="price_observation">Price Observations</option>
          <option value="aggregated_price">Aggregated Prices</option>
        </select>
      </div>

      {/* Loading */}
      {isLoading && (
        <div className="bg-white rounded-xl shadow-lg p-12 text-center">
          <Loader2 className="w-10 h-10 text-[#1a5f3f] mx-auto mb-3 animate-spin" />
          <p className="text-gray-600">Loading audit logs...</p>
        </div>
      )}

      {/* Error */}
      {isError && (
        <div className="bg-red-50 border border-red-200 rounded-xl p-6 text-center">
          <p className="text-red-700 font-medium">Failed to load audit logs.</p>
        </div>
      )}

      {/* Empty */}
      {!isLoading && !isError && logs.length === 0 && (
        <div className="bg-white rounded-xl shadow-lg p-12 text-center">
          <ScrollText className="w-12 h-12 text-gray-300 mx-auto mb-4" />
          <p className="text-gray-500 text-lg">No audit records yet.</p>
          <p className="text-gray-400 text-sm mt-1">Actions will appear here once admins approve, flag, or reject observations.</p>
        </div>
      )}

      {/* Log Table */}
      {!isLoading && !isError && logs.length > 0 && (
        <div className="bg-white rounded-xl shadow-lg overflow-hidden">
          {/* Desktop */}
          <div className="hidden md:block overflow-x-auto">
            <table className="w-full">
              <thead className="bg-gray-50 border-b border-gray-200">
                <tr>
                  <th className="px-6 py-4 text-left text-xs font-semibold text-gray-700 uppercase tracking-wider">Time</th>
                  <th className="px-6 py-4 text-left text-xs font-semibold text-gray-700 uppercase tracking-wider">Admin</th>
                  <th className="px-6 py-4 text-left text-xs font-semibold text-gray-700 uppercase tracking-wider">Action</th>
                  <th className="px-6 py-4 text-left text-xs font-semibold text-gray-700 uppercase tracking-wider">Entity</th>
                  <th className="px-6 py-4 text-left text-xs font-semibold text-gray-700 uppercase tracking-wider">Change</th>
                  <th className="px-6 py-4 text-left text-xs font-semibold text-gray-700 uppercase tracking-wider">Reason</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-200">
                {logs.map((log) => (
                  <tr key={log.id} className="hover:bg-gray-50">
                    <td className="px-6 py-4 text-sm text-gray-600 whitespace-nowrap">
                      {new Date(log.created_at).toLocaleString('en-NG', {
                        month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit',
                      })}
                    </td>
                    <td className="px-6 py-4 font-medium text-gray-900">{log.admin_id}</td>
                    <td className="px-6 py-4">
                      <span className={`px-3 py-1 rounded-full text-xs font-medium ${
                        log.action === 'approved' ? 'bg-green-100 text-green-700'
                        : log.action === 'flagged' ? 'bg-amber-100 text-amber-700'
                        : log.action === 'rejected' ? 'bg-red-100 text-red-700'
                        : 'bg-blue-100 text-blue-700'
                      }`}>
                        {log.action}
                      </span>
                    </td>
                    <td className="px-6 py-4 text-sm text-gray-600">
                      {log.entity_type.replace('_', ' ')}
                      <span className="block text-xs text-gray-400 font-mono">{log.entity_id.slice(0, 8)}…</span>
                    </td>
                    <td className="px-6 py-4 text-sm text-gray-600">
                      {log.old_value && log.new_value ? (
                        <span>
                          {JSON.parse(log.old_value).status} → {JSON.parse(log.new_value).status}
                        </span>
                      ) : '—'}
                    </td>
                    <td className="px-6 py-4 text-sm text-gray-600 max-w-xs truncate">
                      {log.reason || '—'}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>

          {/* Mobile Cards */}
          <div className="md:hidden divide-y divide-gray-200">
            {logs.map((log) => (
              <div key={log.id} className="p-5">
                <div className="flex items-start justify-between mb-2">
                  <span className={`px-3 py-1 rounded-full text-xs font-medium ${
                    log.action === 'approved' ? 'bg-green-100 text-green-700'
                    : log.action === 'flagged' ? 'bg-amber-100 text-amber-700'
                    : log.action === 'rejected' ? 'bg-red-100 text-red-700'
                    : 'bg-blue-100 text-blue-700'
                  }`}>
                    {log.action}
                  </span>
                  <span className="text-xs text-gray-500">
                    {new Date(log.created_at).toLocaleString('en-NG', {
                      month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit',
                    })}
                  </span>
                </div>
                <p className="text-sm text-gray-700 mb-1">
                  <strong>{log.admin_id}</strong> on {log.entity_type.replace('_', ' ')}
                </p>
                {log.reason && <p className="text-sm text-gray-500">{log.reason}</p>}
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}