import { useState } from 'react';
import { CheckCircle, Flag, XCircle, Filter, Loader2 } from 'lucide-react';
import { useObservations, useUpdateObservationStatus } from '../../api/hooks';
import type { PriceObservation } from '../../api/types';

type StatusFilter = 'all' | 'pending' | 'approved' | 'flagged' | 'rejected';

function StatusBadge({ status }: { status: string }) {
  const styles: Record<string, string> = {
    approved: 'bg-green-100 text-green-700',
    pending: 'bg-amber-100 text-amber-700',
    flagged: 'bg-red-100 text-red-700',
    rejected: 'bg-gray-100 text-gray-700',
  };
  return (
    <span className={`px-3 py-1 rounded-full text-xs font-medium ${styles[status] || styles.pending}`}>
      {status.charAt(0).toUpperCase() + status.slice(1)}
    </span>
  );
}

function ActionButtons({ observation }: { observation: PriceObservation }) {
  const mutation = useUpdateObservationStatus();
  const [reason, setReason] = useState('');
  const [showReason, setShowReason] = useState<string | null>(null); // which action needs reason

  const handleAction = (newStatus: string) => {
    if (newStatus === 'flagged' || newStatus === 'rejected') {
      if (!showReason) {
        setShowReason(newStatus);
        return;
      }
    }
    mutation.mutate(
      { id: observation.id, status: newStatus, reason },
      { onSuccess: () => { setShowReason(null); setReason(''); } },
    );
  };

  if (showReason) {
    return (
      <div className="flex items-center gap-2">
        <input
          type="text"
          placeholder="Reason..."
          value={reason}
          onChange={(e) => setReason(e.target.value)}
          className="px-2 py-1 border rounded text-sm w-32"
        />
        <button
          onClick={() => handleAction(showReason)}
          disabled={mutation.isPending}
          className="px-2 py-1 bg-red-600 text-white text-xs rounded hover:bg-red-700 disabled:opacity-50"
        >
          Confirm
        </button>
        <button
          onClick={() => { setShowReason(null); setReason(''); }}
          className="px-2 py-1 bg-gray-200 text-gray-700 text-xs rounded hover:bg-gray-300"
        >
          Cancel
        </button>
      </div>
    );
  }

  return (
    <div className="flex items-center gap-1">
      {observation.status !== 'approved' && (
        <button
          onClick={() => handleAction('approved')}
          disabled={mutation.isPending}
          className="p-2 text-green-600 hover:bg-green-50 rounded-lg transition-colors disabled:opacity-50"
          title="Approve"
        >
          <CheckCircle className="w-4 h-4" />
        </button>
      )}
      {observation.status !== 'flagged' && (
        <button
          onClick={() => handleAction('flagged')}
          disabled={mutation.isPending}
          className="p-2 text-amber-600 hover:bg-amber-50 rounded-lg transition-colors disabled:opacity-50"
          title="Flag"
        >
          <Flag className="w-4 h-4" />
        </button>
      )}
      {observation.status !== 'rejected' && (
        <button
          onClick={() => handleAction('rejected')}
          disabled={mutation.isPending}
          className="p-2 text-red-600 hover:bg-red-50 rounded-lg transition-colors disabled:opacity-50"
          title="Reject"
        >
          <XCircle className="w-4 h-4" />
        </button>
      )}
    </div>
  );
}

export default function Submissions() {
  const [filter, setFilter] = useState<StatusFilter>('all');
  const statusParam = filter === 'all' ? undefined : filter;
  const { data, isLoading, isError } = useObservations(undefined, undefined, statusParam);

  const observations = data?.data ?? [];
  const total = data?.total ?? 0;

  return (
    <div className="p-6 md:p-8">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900 mb-2">Price Submissions</h1>
        <p className="text-gray-600">Review and manage submitted price observations</p>
      </div>

      {/* Filter Tabs */}
      <div className="bg-white rounded-xl shadow-md mb-6 p-4 flex items-center gap-4 flex-wrap">
        <div className="flex items-center gap-2">
          <Filter className="w-5 h-5 text-gray-600" />
          <span className="font-medium text-gray-700">Filter:</span>
        </div>
        <div className="flex gap-2 flex-wrap">
          {(['all', 'pending', 'approved', 'flagged', 'rejected'] as StatusFilter[]).map((f) => {
            const colors: Record<StatusFilter, string> = {
              all: 'bg-[#1a5f3f] text-white',
              pending: 'bg-amber-500 text-white',
              approved: 'bg-green-500 text-white',
              flagged: 'bg-red-500 text-white',
              rejected: 'bg-gray-600 text-white',
            };
            return (
              <button
                key={f}
                onClick={() => setFilter(f)}
                className={`px-4 py-2 rounded-lg font-medium transition-colors ${
                  filter === f ? colors[f] : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                }`}
              >
                {f.charAt(0).toUpperCase() + f.slice(1)}
                {f === 'all' ? ` (${total})` : ''}
              </button>
            );
          })}
        </div>
      </div>

      {/* Loading */}
      {isLoading && (
        <div className="bg-white rounded-xl shadow-lg p-12 text-center">
          <Loader2 className="w-10 h-10 text-[#1a5f3f] mx-auto mb-3 animate-spin" />
          <p className="text-gray-600">Loading observations...</p>
        </div>
      )}

      {/* Error */}
      {isError && (
        <div className="bg-red-50 border border-red-200 rounded-xl p-6 text-center">
          <p className="text-red-700 font-medium">Failed to load observations. Is the API running?</p>
        </div>
      )}

      {/* Empty */}
      {!isLoading && !isError && observations.length === 0 && (
        <div className="bg-white rounded-xl shadow-lg p-12 text-center">
          <p className="text-gray-500 text-lg">No observations found.</p>
        </div>
      )}

      {/* Table */}
      {!isLoading && !isError && observations.length > 0 && (
        <div className="bg-white rounded-xl shadow-lg overflow-hidden">
          {/* Desktop Table */}
          <div className="hidden md:block overflow-x-auto">
            <table className="w-full">
              <thead className="bg-gray-50 border-b border-gray-200">
                <tr>
                  <th className="px-6 py-4 text-left text-xs font-semibold text-gray-700 uppercase tracking-wider">Crop</th>
                  <th className="px-6 py-4 text-left text-xs font-semibold text-gray-700 uppercase tracking-wider">Market</th>
                  <th className="px-6 py-4 text-left text-xs font-semibold text-gray-700 uppercase tracking-wider">Price</th>
                  <th className="px-6 py-4 text-left text-xs font-semibold text-gray-700 uppercase tracking-wider">Source</th>
                  <th className="px-6 py-4 text-left text-xs font-semibold text-gray-700 uppercase tracking-wider">Observed</th>
                  <th className="px-6 py-4 text-left text-xs font-semibold text-gray-700 uppercase tracking-wider">Confidence</th>
                  <th className="px-6 py-4 text-left text-xs font-semibold text-gray-700 uppercase tracking-wider">Status</th>
                  <th className="px-6 py-4 text-right text-xs font-semibold text-gray-700 uppercase tracking-wider">Actions</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-200">
                {observations.map((obs) => (
                  <tr key={obs.id} className="hover:bg-gray-50">
                    <td className="px-6 py-4 font-medium text-gray-900">{obs.crop_name}</td>
                    <td className="px-6 py-4 text-gray-700">{obs.market_name}</td>
                    <td className="px-6 py-4 font-semibold text-gray-900">
                      ₦{obs.price.toLocaleString()} / {obs.unit}
                    </td>
                    <td className="px-6 py-4 text-gray-700">{obs.source}</td>
                    <td className="px-6 py-4 text-sm text-gray-600">
                      {new Date(obs.observed_at).toLocaleDateString('en-NG', {
                        month: 'short', day: 'numeric', year: 'numeric',
                      })}
                    </td>
                    <td className="px-6 py-4 text-sm text-gray-600">
                      {(obs.confidence_score * 100).toFixed(0)}%
                    </td>
                    <td className="px-6 py-4"><StatusBadge status={obs.status} /></td>
                    <td className="px-6 py-4">
                      <div className="flex justify-end">
                        <ActionButtons observation={obs} />
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>

          {/* Mobile Cards */}
          <div className="md:hidden divide-y divide-gray-200">
            {observations.map((obs) => (
              <div key={obs.id} className="p-5">
                <div className="flex items-start justify-between mb-3">
                  <div>
                    <h3 className="font-semibold text-gray-900 mb-1">{obs.crop_name}</h3>
                    <p className="text-sm text-gray-600">{obs.market_name}</p>
                  </div>
                  <StatusBadge status={obs.status} />
                </div>
                <div className="grid grid-cols-2 gap-3 mb-3">
                  <div>
                    <p className="text-xs text-gray-600 mb-1">Price</p>
                    <p className="font-semibold text-gray-900">₦{obs.price.toLocaleString()} / {obs.unit}</p>
                  </div>
                  <div>
                    <p className="text-xs text-gray-600 mb-1">Confidence</p>
                    <p className="font-semibold text-gray-900">{(obs.confidence_score * 100).toFixed(0)}%</p>
                  </div>
                </div>
                <div className="flex items-center justify-between text-sm mb-3">
                  <span className="text-gray-600">{obs.source}</span>
                  <span className="text-gray-600">
                    {new Date(obs.observed_at).toLocaleDateString('en-NG', {
                      month: 'short', day: 'numeric', year: 'numeric',
                    })}
                  </span>
                </div>
                <ActionButtons observation={obs} />
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}