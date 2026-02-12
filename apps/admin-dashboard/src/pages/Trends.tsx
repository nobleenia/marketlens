import { useState, useMemo } from 'react';
import { useSearchParams } from 'react-router';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts';
import { TrendingUp, TrendingDown, Activity, Loader2 } from 'lucide-react';
import { useCrops, useMarkets, usePrices } from '../api/hooks';

export default function Trends() {
  const [searchParams, setSearchParams] = useSearchParams();

  // ── Dropdowns store NAMES (from URL params or user selection) ─────
  const [selectedCrop, setSelectedCrop] = useState(searchParams.get('crop') || '');
  const [selectedMarket, setSelectedMarket] = useState(searchParams.get('market') || '');
  const [period, setPeriod] = useState<7 | 30>(7);

  // ── API data for dropdowns ────────────────────────────────────────
  const { data: crops = [], isLoading: cropsLoading } = useCrops();
  const { data: markets = [], isLoading: marketsLoading } = useMarkets();

  // ── Compute date range from period ────────────────────────────────
  const today = new Date().toISOString().slice(0, 10);
  const fromDate = new Date(Date.now() - period * 86_400_000).toISOString().slice(0, 10);

  // ── Fetch real prices for chart ───────────────────────────────────
  const {
    data: prices,
    isLoading: pricesLoading,
  } = usePrices(
    selectedCrop || undefined,
    selectedMarket || undefined,
    fromDate,
    today,
  );

  // ── Map API data → chart format ───────────────────────────────────
  const chartData = useMemo(() => {
    if (!prices || prices.length === 0) return [];
    return prices
      .map((p) => ({
        date: p.period_start.slice(0, 10),
        price: p.price_mean,
      }))
      .sort((a, b) => a.date.localeCompare(b.date));
  }, [prices]);

  // ── Compute stats from real data ──────────────────────────────────
  const stats = useMemo(() => {
    if (chartData.length === 0) return null;

    const values = chartData.map((d) => d.price);
    const highest = Math.max(...values);
    const lowest = Math.min(...values);
    const avg = values.reduce((a, b) => a + b, 0) / values.length;
    const variance = values.reduce((sum, v) => sum + Math.pow(v - avg, 2), 0) / values.length;
    const stdDev = Math.sqrt(variance);
    const volatility = (stdDev / avg) * 100;

    let volatilityLevel = 'Low';
    if (volatility > 10) volatilityLevel = 'High';
    else if (volatility > 5) volatilityLevel = 'Moderate';

    return { highest, lowest, volatility: volatilityLevel };
  }, [chartData]);

  // ── Dropdown change handlers (also sync URL) ──────────────────────
  const handleCropChange = (name: string) => {
    setSelectedCrop(name);
    setSearchParams({ crop: name, market: selectedMarket });
  };
  const handleMarketChange = (name: string) => {
    setSelectedMarket(name);
    setSearchParams({ crop: selectedCrop, market: name });
  };

  return (
    <div className="min-h-screen bg-gradient-to-b from-[#e8f5e9] to-white pb-20 md:pb-8">
      <div className="container mx-auto px-4 py-8">
        <div className="mb-8">
          <h1 className="text-3xl md:text-4xl font-bold text-[#1a5f3f] mb-2">
            Price Trends
          </h1>
          <p className="text-gray-700 text-lg">
            Track price movements over time
          </p>
        </div>

        {/* Selection Panel */}
        <div className="bg-white rounded-2xl shadow-lg p-6 mb-6">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
            <div>
              <label htmlFor="crop-trend" className="block text-sm font-medium text-gray-700 mb-2">
                Select Crop
              </label>
              <select
                id="crop-trend"
                value={selectedCrop}
                onChange={(e) => handleCropChange(e.target.value)}
                disabled={cropsLoading}
                className="w-full px-4 py-3 border-2 border-gray-300 rounded-lg focus:border-[#1a5f3f] focus:outline-none"
              >
                <option value="">{cropsLoading ? 'Loading...' : 'Choose a crop...'}</option>
                {crops.map((crop) => (
                  <option key={crop.id} value={crop.name}>
                    {crop.name}
                  </option>
                ))}
              </select>
            </div>

            <div>
              <label htmlFor="market-trend" className="block text-sm font-medium text-gray-700 mb-2">
                Select Market
              </label>
              <select
                id="market-trend"
                value={selectedMarket}
                onChange={(e) => handleMarketChange(e.target.value)}
                disabled={marketsLoading}
                className="w-full px-4 py-3 border-2 border-gray-300 rounded-lg focus:border-[#1a5f3f] focus:outline-none"
              >
                <option value="">{marketsLoading ? 'Loading...' : 'Choose a market...'}</option>
                {markets.map((market) => (
                  <option key={market.id} value={market.name}>
                    {market.name} ({market.state})
                  </option>
                ))}
              </select>
            </div>
          </div>

          {/* Period Toggle */}
          <div className="flex items-center gap-3">
            <span className="text-sm font-medium text-gray-700">Period:</span>
            <div className="flex gap-2">
              <button
                onClick={() => setPeriod(7)}
                className={`px-4 py-2 rounded-lg font-medium transition-colors ${
                  period === 7
                    ? 'bg-[#1a5f3f] text-white'
                    : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                }`}
              >
                7 Days
              </button>
              <button
                onClick={() => setPeriod(30)}
                className={`px-4 py-2 rounded-lg font-medium transition-colors ${
                  period === 30
                    ? 'bg-[#1a5f3f] text-white'
                    : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                }`}
              >
                30 Days
              </button>
            </div>
          </div>
        </div>

        {/* Loading State */}
        {pricesLoading && selectedCrop && selectedMarket && (
          <div className="bg-white rounded-2xl shadow-lg p-12 text-center">
            <Loader2 className="w-12 h-12 text-[#1a5f3f] mx-auto mb-4 animate-spin" />
            <p className="text-gray-600">Loading price data...</p>
          </div>
        )}

        {/* Chart and Stats */}
        {!pricesLoading && chartData.length > 0 && (
          <div className="space-y-6">
            {/* Header */}
            <div className="bg-gradient-to-r from-[#1a5f3f] to-[#2d8659] text-white rounded-2xl shadow-lg p-6">
              <h2 className="text-2xl font-bold mb-1">{selectedCrop}</h2>
              <p className="opacity-90">{selectedMarket}</p>
            </div>

            {/* Stats Cards */}
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              <div className="bg-white rounded-xl shadow-md p-6 border border-gray-100">
                <div className="flex items-center gap-3 mb-2">
                  <TrendingUp className="w-6 h-6 text-green-600" />
                  <span className="text-sm text-gray-600">Highest Price</span>
                </div>
                <p className="text-3xl font-bold text-[#1a5f3f]">
                  ₦{stats?.highest.toLocaleString()}
                </p>
              </div>

              <div className="bg-white rounded-xl shadow-md p-6 border border-gray-100">
                <div className="flex items-center gap-3 mb-2">
                  <TrendingDown className="w-6 h-6 text-red-600" />
                  <span className="text-sm text-gray-600">Lowest Price</span>
                </div>
                <p className="text-3xl font-bold text-[#1a5f3f]">
                  ₦{stats?.lowest.toLocaleString()}
                </p>
              </div>

              <div className="bg-white rounded-xl shadow-md p-6 border border-gray-100">
                <div className="flex items-center gap-3 mb-2">
                  <Activity className="w-6 h-6 text-amber-600" />
                  <span className="text-sm text-gray-600">Volatility</span>
                </div>
                <p className="text-3xl font-bold text-[#1a5f3f]">
                  {stats?.volatility}
                </p>
              </div>
            </div>

            {/* Chart */}
            <div className="bg-white rounded-2xl shadow-lg p-6">
              <h3 className="text-xl font-semibold text-gray-900 mb-6">
                Price History ({period} Days)
              </h3>
              <div className="h-80 md:h-96">
                <ResponsiveContainer width="100%" height="100%">
                  <LineChart data={chartData}>
                    <CartesianGrid strokeDasharray="3 3" stroke="#e5e7eb" />
                    <XAxis
                      dataKey="date"
                      stroke="#6b7280"
                      tick={{ fontSize: 12 }}
                      tickFormatter={(value) => {
                        const date = new Date(value);
                        return `${date.getMonth() + 1}/${date.getDate()}`;
                      }}
                    />
                    <YAxis
                      stroke="#6b7280"
                      tick={{ fontSize: 12 }}
                      tickFormatter={(value) => `₦${(value / 1000).toFixed(0)}k`}
                    />
                    <Tooltip
                      contentStyle={{
                        backgroundColor: 'white',
                        border: '1px solid #e5e7eb',
                        borderRadius: '8px',
                      }}
                      formatter={(value) => [`₦${Number(value).toLocaleString()}`, 'Price']}
                      labelFormatter={(label) => {
                        const date = new Date(label);
                        return date.toLocaleDateString('en-NG', {
                          month: 'short',
                          day: 'numeric',
                          year: 'numeric',
                        });
                      }}
                    />
                    <Line
                      type="monotone"
                      dataKey="price"
                      stroke="#1a5f3f"
                      strokeWidth={3}
                      dot={{ fill: '#1a5f3f', r: 4 }}
                      activeDot={{ r: 6 }}
                    />
                  </LineChart>
                </ResponsiveContainer>
              </div>
            </div>
          </div>
        )}

        {/* No Data State */}
        {!pricesLoading && selectedCrop && selectedMarket && chartData.length === 0 && (
          <div className="bg-white rounded-2xl shadow-lg p-12 text-center">
            <Activity className="w-16 h-16 text-gray-300 mx-auto mb-4" />
            <h3 className="text-xl font-semibold text-gray-900 mb-2">
              No Price Data Available
            </h3>
            <p className="text-gray-600">
              No aggregated prices found for {selectedCrop} at {selectedMarket} in the last {period} days.
            </p>
          </div>
        )}

        {/* Empty State — nothing selected */}
        {(!selectedCrop || !selectedMarket) && (
          <div className="bg-white rounded-2xl shadow-lg p-12 text-center">
            <TrendingUp className="w-16 h-16 text-gray-300 mx-auto mb-4" />
            <h3 className="text-xl font-semibold text-gray-900 mb-2">
              Select a Crop and Market
            </h3>
            <p className="text-gray-600">
              Choose options above to view price trends
            </p>
          </div>
        )}
      </div>
    </div>
  );
}