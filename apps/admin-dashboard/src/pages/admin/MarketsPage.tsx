import { useState } from 'react';
import { Loader2, Store, MapPin, Plus, X } from 'lucide-react';
import { useMarkets, useCreateMarket } from '../../api/hooks';

export default function MarketsPage() {
  const { data: markets = [], isLoading, isError } = useMarkets();
  const [showForm, setShowForm] = useState(false);
  const [name, setName] = useState('');
  const [state, setState] = useState('');
  const [country, setCountry] = useState('NG');
  const [lat, setLat] = useState('');
  const [lng, setLng] = useState('');
  const createMarket = useCreateMarket();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!name.trim() || !state.trim()) return;
    try {
      await createMarket.mutateAsync({
        name: name.trim(),
        state: state.trim(),
        country: country.trim() || 'NG',
        latitude: lat ? Number(lat) : 0,
        longitude: lng ? Number(lng) : 0,
      });
      setName(''); setState(''); setLat(''); setLng('');
      setShowForm(false);
    } catch {
      // error is available via createMarket.error
    }
  };

  // Group by state for a cleaner view
  const byState = markets.reduce<Record<string, typeof markets>>((acc, m) => {
    (acc[m.state] ??= []).push(m);
    return acc;
  }, {});
  const states = Object.keys(byState).sort();

  return (
    <div className="p-6 md:p-8">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900 mb-2">Markets</h1>
        <p className="text-gray-600">All registered markets across Nigeria</p>
      </div>

      {isLoading && (
        <div className="bg-white rounded-xl shadow-lg p-12 text-center">
          <Loader2 className="w-10 h-10 text-[#1a5f3f] mx-auto mb-3 animate-spin" />
          <p className="text-gray-600">Loading markets...</p>
        </div>
      )}

      {isError && (
        <div className="bg-red-50 border border-red-200 rounded-xl p-6 text-center">
          <p className="text-red-700 font-medium">Failed to load markets.</p>
        </div>
      )}

      {!isLoading && !isError && (
        <>
          <div className="bg-white rounded-xl shadow-md p-6 mb-6">
            <div className="flex items-center justify-between">
              <div className="grid grid-cols-2 md:grid-cols-3 gap-4 flex-1">
                <div>
                  <p className="text-sm text-gray-600">Total Markets</p>
                  <p className="text-3xl font-bold text-[#1a5f3f]">{markets.length}</p>
                </div>
                <div>
                  <p className="text-sm text-gray-600">States Covered</p>
                  <p className="text-3xl font-bold text-[#1a5f3f]">{states.length}</p>
                </div>
              </div>
              <button
                onClick={() => setShowForm(!showForm)}
                className="flex items-center gap-2 px-4 py-2 bg-[#1a5f3f] text-white rounded-lg hover:bg-[#155232] transition-colors text-sm font-medium"
              >
                {showForm ? <X className="w-4 h-4" /> : <Plus className="w-4 h-4" />}
                {showForm ? 'Cancel' : 'Add Market'}
              </button>
            </div>
          </div>

          {/* Add Market Form */}
          {showForm && (
            <div className="bg-white rounded-xl shadow-md p-6 mb-6 border-2 border-blue-200 bg-blue-50">
              <h3 className="text-lg font-semibold text-gray-900 mb-4">New Market</h3>
              <form onSubmit={handleSubmit} className="space-y-4">
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">Market Name *</label>
                    <input
                      value={name}
                      onChange={(e) => setName(e.target.value)}
                      placeholder="e.g. Bodija Market"
                      className="w-full px-4 py-2 border-2 border-gray-300 rounded-lg focus:border-[#1a5f3f] focus:outline-none"
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">State *</label>
                    <input
                      value={state}
                      onChange={(e) => setState(e.target.value)}
                      placeholder="e.g. Oyo"
                      className="w-full px-4 py-2 border-2 border-gray-300 rounded-lg focus:border-[#1a5f3f] focus:outline-none"
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">Country</label>
                    <input
                      value={country}
                      onChange={(e) => setCountry(e.target.value)}
                      className="w-full px-4 py-2 border-2 border-gray-300 rounded-lg focus:border-[#1a5f3f] focus:outline-none"
                    />
                  </div>
                  <div className="grid grid-cols-2 gap-2">
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-1">Latitude</label>
                      <input
                        value={lat}
                        onChange={(e) => setLat(e.target.value)}
                        placeholder="0.0"
                        className="w-full px-4 py-2 border-2 border-gray-300 rounded-lg focus:border-[#1a5f3f] focus:outline-none"
                      />
                    </div>
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-1">Longitude</label>
                      <input
                        value={lng}
                        onChange={(e) => setLng(e.target.value)}
                        placeholder="0.0"
                        className="w-full px-4 py-2 border-2 border-gray-300 rounded-lg focus:border-[#1a5f3f] focus:outline-none"
                      />
                    </div>
                  </div>
                </div>
                <div className="flex items-center gap-3">
                  <button
                    type="submit"
                    disabled={createMarket.isPending || !name.trim() || !state.trim()}
                    className="px-6 py-2 bg-[#1a5f3f] text-white rounded-lg hover:bg-[#155232] transition-colors text-sm font-medium disabled:opacity-50"
                  >
                    {createMarket.isPending ? 'Creating...' : 'Create Market'}
                  </button>
                  <button
                    type="button"
                    onClick={() => setShowForm(false)}
                    className="px-6 py-2 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition-colors text-sm"
                  >
                    Cancel
                  </button>
                </div>
                {createMarket.isError && (
                  <p className="text-red-600 text-sm">{createMarket.error?.message || 'Failed to create market'}</p>
                )}
              </form>
            </div>
          )}

          {states.map((state) => (
            <div key={state} className="bg-white rounded-xl shadow-lg overflow-hidden mb-6">
              <div className="p-5 border-b border-gray-200 bg-gray-50">
                <div className="flex items-center gap-2">
                  <MapPin className="w-5 h-5 text-[#1a5f3f]" />
                  <h2 className="text-lg font-semibold text-gray-900">{state}</h2>
                  <span className="text-sm text-gray-500">({byState[state].length})</span>
                </div>
              </div>

              {/* Desktop Table */}
              <div className="hidden md:block overflow-x-auto">
                <table className="w-full">
                  <thead className="bg-gray-50 border-b border-gray-200">
                    <tr>
                      <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700 uppercase tracking-wider">Name</th>
                      <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700 uppercase tracking-wider">Country</th>
                      <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700 uppercase tracking-wider">Coordinates</th>
                      <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700 uppercase tracking-wider">ID</th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-gray-200">
                    {byState[state].map((m) => (
                      <tr key={m.id} className="hover:bg-gray-50">
                        <td className="px-6 py-4">
                          <div className="flex items-center gap-3">
                            <div className="bg-blue-100 rounded-lg p-2">
                              <Store className="w-4 h-4 text-blue-600" />
                            </div>
                            <span className="font-medium text-gray-900">{m.name}</span>
                          </div>
                        </td>
                        <td className="px-6 py-4 text-gray-700">{m.country}</td>
                        <td className="px-6 py-4 text-sm text-gray-600">
                          {m.latitude !== 0 ? `${m.latitude.toFixed(4)}, ${m.longitude.toFixed(4)}` : '—'}
                        </td>
                        <td className="px-6 py-4 text-xs text-gray-400 font-mono">{m.id.slice(0, 8)}…</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>

              {/* Mobile Cards */}
              <div className="md:hidden divide-y divide-gray-200">
                {byState[state].map((m) => (
                  <div key={m.id} className="p-5 flex items-center gap-4">
                    <div className="bg-blue-100 rounded-lg p-3">
                      <Store className="w-5 h-5 text-blue-600" />
                    </div>
                    <div>
                      <p className="font-semibold text-gray-900">{m.name}</p>
                      <p className="text-sm text-gray-500">{m.country}</p>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          ))}
        </>
      )}
    </div>
  );
}