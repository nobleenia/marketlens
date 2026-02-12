import { useState } from 'react';
import { Loader2, TrendingUp, TrendingDown, Minus } from 'lucide-react';
import { useCrops, useMarkets, usePrices } from '../../api/hooks';

export default function PricesPage() {
  const { data: crops = [], isLoading: cropsLoading } = useCrops();
  const { data: markets = [], isLoading: marketsLoading } = useMarkets();

  const [selectedCrop, setSelectedCrop] = useState('');
  const [selectedMarket, setSelectedMarket] = useState('');

  const { data: prices = [], isLoading: pricesLoading } = usePrices(
    selectedCrop || undefined,
    selectedMarket || undefined,
  );

  const isLoading = cropsLoading || marketsLoading;

  const TrendIcon = ({ trend }: { trend: string }) => {
    if (trend === 'up') return <TrendingUp className="w-4 h-4 text-green-600" />;
    if (trend === 'down') return <TrendingDown className="w-4 h-4 text-red-600" />;
    return <Minus className="w-4 h-4 text-gray-400" />;
  };

  return (
    <div className="p-6 md:p-8">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900 mb-2">Aggregated Prices</h1>
        <p className="text-gray-600">Browse computed price aggregations</p>
      </div>

      {/* Filters */}
      <div className="bg-white rounded-xl shadow-md p-6 mb-6">
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label htmlFor="crop-filter" className="block text-sm font-medium text-gray-700 mb-2">Crop</label>
            <select
              id="crop-filter"
              value={selectedCrop}
              onChange={(e) => setSelectedCrop(e.target.value)}
              disabled={isLoading}
              className="w-full px-4 py-3 border-2 border-gray-300 rounded-lg focus:border-[#1a5f3f] focus:outline-none"
            >
              <option value="">{cropsLoading ? 'Loading...' : 'All crops'}</option>
              {crops.map((c) => (
                <option key={c.id} value={c.name}>{c.name}</option>
              ))}
            </select>
          </div>
          <div>
            <label htmlFor="market-filter" className="block text-sm font-medium text-gray-700 mb-2">Market</label>
            <select
              id="market-filter"
              value={selectedMarket}
              onChange={(e) => setSelectedMarket(e.target.value)}
              disabled={isLoading}
              className="w-full px-4 py-3 border-2 border-gray-300 rounded-lg focus:border-[#1a5f3f] focus:outline-none"
            >
              <option value="">{marketsLoading ? 'Loading...' : 'All markets'}</option>
              {markets.map((m) => (
                <option key={m.id} value={m.name}>{m.name} ({m.state})</option>
              ))}
            </select>
          </div>
        </div>
      </div>

      {/* Prompt */}
      {!selectedCrop && !selectedMarket && (
        <div className="bg-white rounded-xl shadow-lg p-12 text-center">
          <p className="text-gray-500 text-lg">Select a crop or market to view aggregated prices.</p>
        </div>
      )}

      {/* Loading */}
      {pricesLoading && (selectedCrop || selectedMarket) && (
        <div className="bg-white rounded-xl shadow-lg p-12 text-center">
          <Loader2 className="w-10 h-10 text-[#1a5f3f] mx-auto mb-3 animate-spin" />
          <p className="text-gray-600">Loading prices...</p>
        </div>
      )}

      {/* Empty */}
      {!pricesLoading && (selectedCrop || selectedMarket) && prices.length === 0 && (
        <div className="bg-white rounded-xl shadow-lg p-12 text-center">
          <p className="text-gray-500 text-lg">No aggregated prices found for this selection.</p>
        </div>
      )}

      {/* Results */}
      {!pricesLoading && prices.length > 0 && (
        <div className="bg-white rounded-xl shadow-lg overflow-hidden">
          <div className="p-6 border-b border-gray-200">
            <h2 className="text-lg font-semibold text-gray-900">{prices.length} price record{prices.length !== 1 ? 's' : ''}</h2>
          </div>

          {/* Desktop Table */}
          <div className="hidden md:block overflow-x-auto">
            <table className="w-full">
              <thead className="bg-gray-50 border-b border-gray-200">
                <tr>
                  <th className="px-6 py-4 text-left text-xs font-semibold text-gray-700 uppercase tracking-wider">Crop</th>
                  <th className="px-6 py-4 text-left text-xs font-semibold text-gray-700 uppercase tracking-wider">Market</th>
                  <th className="px-6 py-4 text-left text-xs font-semibold text-gray-700 uppercase tracking-wider">Period</th>
                  <th className="px-6 py-4 text-left text-xs font-semibold text-gray-700 uppercase tracking-wider">Min</th>
                  <th className="px-6 py-4 text-left text-xs font-semibold text-gray-700 uppercase tracking-wider">Max</th>
                  <th className="px-6 py-4 text-left text-xs font-semibold text-gray-700 uppercase tracking-wider">Mean</th>
                  <th className="px-6 py-4 text-left text-xs font-semibold text-gray-700 uppercase tracking-wider">Confidence</th>
                  <th className="px-6 py-4 text-left text-xs font-semibold text-gray-700 uppercase tracking-wider">Samples</th>
                  <th className="px-6 py-4 text-left text-xs font-semibold text-gray-700 uppercase tracking-wider">Trend</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-200">
                {prices.map((p) => (
                  <tr key={p.id} className="hover:bg-gray-50">
                    <td className="px-6 py-4 font-medium text-gray-900">{p.crop_name}</td>
                    <td className="px-6 py-4 text-gray-700">{p.market_name}</td>
                    <td className="px-6 py-4 text-sm text-gray-600">
                      {p.period_start.slice(0, 10)} → {p.period_end.slice(0, 10)}
                    </td>
                    <td className="px-6 py-4 font-semibold text-gray-900">₦{p.price_min.toLocaleString()}</td>
                    <td className="px-6 py-4 font-semibold text-gray-900">₦{p.price_max.toLocaleString()}</td>
                    <td className="px-6 py-4 font-semibold text-[#1a5f3f]">₦{p.price_mean.toLocaleString()}</td>
                    <td className="px-6 py-4">
                      <span className={`px-3 py-1 rounded-full text-xs font-medium ${
                        p.confidence === 'high' ? 'bg-green-100 text-green-700'
                        : p.confidence === 'medium' ? 'bg-amber-100 text-amber-700'
                        : 'bg-red-100 text-red-700'
                      }`}>
                        {p.confidence}
                      </span>
                    </td>
                    <td className="px-6 py-4 text-sm text-gray-600">{p.sample_size}</td>
                    <td className="px-6 py-4"><TrendIcon trend={p.trend} /></td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>

          {/* Mobile Cards */}
          <div className="md:hidden divide-y divide-gray-200">
            {prices.map((p) => (
              <div key={p.id} className="p-5">
                <div className="flex items-start justify-between mb-2">
                  <div>
                    <h3 className="font-semibold text-gray-900">{p.crop_name}</h3>
                    <p className="text-sm text-gray-600">{p.market_name}</p>
                  </div>
                  <TrendIcon trend={p.trend} />
                </div>
                <div className="grid grid-cols-3 gap-3 mb-2">
                  <div>
                    <p className="text-xs text-gray-500">Min</p>
                    <p className="font-semibold text-gray-900">₦{p.price_min.toLocaleString()}</p>
                  </div>
                  <div>
                    <p className="text-xs text-gray-500">Mean</p>
                    <p className="font-semibold text-[#1a5f3f]">₦{p.price_mean.toLocaleString()}</p>
                  </div>
                  <div>
                    <p className="text-xs text-gray-500">Max</p>
                    <p className="font-semibold text-gray-900">₦{p.price_max.toLocaleString()}</p>
                  </div>
                </div>
                <div className="flex items-center justify-between text-sm text-gray-500">
                  <span>{p.period_start.slice(0, 10)} → {p.period_end.slice(0, 10)}</span>
                  <span className={`px-2 py-0.5 rounded-full text-xs font-medium ${
                    p.confidence === 'high' ? 'bg-green-100 text-green-700'
                    : p.confidence === 'medium' ? 'bg-amber-100 text-amber-700'
                    : 'bg-red-100 text-red-700'
                  }`}>{p.confidence}</span>
                </div>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}