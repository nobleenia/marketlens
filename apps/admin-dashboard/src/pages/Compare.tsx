import { useState, useMemo } from 'react';
import { useSearchParams } from 'react-router';
import { TrendingUp, TrendingDown, Minus, Trophy, Loader2 } from 'lucide-react';
import { useCrops, usePrices } from '../api/hooks';

export default function Compare() {
  const [searchParams, setSearchParams] = useSearchParams();
  const [selectedCrop, setSelectedCrop] = useState(searchParams.get('crop') || '');

  // ── API data ──────────────────────────────────────────────────────
  const { data: crops = [], isLoading: cropsLoading } = useCrops();

  // Fetch ALL prices for the selected crop (across all markets)
  // GET /v1/prices?crop=Tomatoes  → returns one row per market
  const { data: rawPrices, isLoading: pricesLoading } = usePrices(
    selectedCrop || undefined,
  );
  const prices = rawPrices ?? [];

  // ── Find the best market (highest average price) ──────────────────
  const bestPrice = useMemo(() => {
    if (!prices || prices.length === 0) return null;
    return prices.reduce((best, current) => {
      const bestAvg = (best.price_min + best.price_max) / 2;
      const currentAvg = (current.price_min + current.price_max) / 2;
      return currentAvg > bestAvg ? current : best;
    });
  }, [prices]);

  // ── Unit from the first price row (same crop = same unit) ─────────
  const unit = prices && prices.length > 0 ? prices[0].unit : '';

  // ── Dropdown handler with URL sync ────────────────────────────────
  const handleCropChange = (name: string) => {
    setSelectedCrop(name);
    if (name) {
      setSearchParams({ crop: name });
    }
  };

  // ── Helper: trend icon ────────────────────────────────────────────
  const getTrendIcon = (trend: string) => {
    switch (trend) {
      case 'up':
        return <TrendingUp className="w-5 h-5 text-green-600" />;
      case 'down':
        return <TrendingDown className="w-5 h-5 text-red-600" />;
      default:
        return <Minus className="w-5 h-5 text-gray-600" />;
    }
  };

  // ── Helper: confidence badge ──────────────────────────────────────
  const getConfidenceBadge = (confidence: string) => {
    const styles: Record<string, string> = {
      high: 'bg-green-100 text-green-700',
      medium: 'bg-amber-100 text-amber-700',
      low: 'bg-red-100 text-red-700',
    };
    return (
      <span className={`px-2 py-1 rounded-full text-xs font-medium ${styles[confidence] || styles.low}`}>
        {confidence.charAt(0).toUpperCase() + confidence.slice(1)}
      </span>
    );
  };

  return (
    <div className="min-h-screen bg-gradient-to-b from-[#e8f5e9] to-white pb-20 md:pb-8">
      <div className="container mx-auto px-4 py-8">
        <div className="mb-8">
          <h1 className="text-3xl md:text-4xl font-bold text-[#1a5f3f] mb-2">
            Compare Markets
          </h1>
          <p className="text-gray-700 text-lg">
            Find the best prices across different markets
          </p>
        </div>

        {/* Crop Selector */}
        <div className="bg-white rounded-2xl shadow-lg p-6 mb-6">
          <label htmlFor="crop-compare" className="block text-sm font-medium text-gray-700 mb-2">
            Select Crop to Compare
          </label>
          <select
            id="crop-compare"
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

        {/* Loading */}
        {pricesLoading && selectedCrop && (
          <div className="bg-white rounded-2xl shadow-lg p-12 text-center">
            <Loader2 className="w-12 h-12 text-[#1a5f3f] mx-auto mb-4 animate-spin" />
            <p className="text-gray-600">Loading market prices...</p>
          </div>
        )}

        {/* Best Market Highlight */}
        {!pricesLoading && selectedCrop && bestPrice && (
          <div className="bg-gradient-to-r from-amber-50 to-yellow-50 border-2 border-amber-300 rounded-2xl p-6 mb-6">
            <div className="flex items-start gap-4">
              <div className="bg-amber-400 rounded-full p-3">
                <Trophy className="w-6 h-6 text-white" />
              </div>
              <div className="flex-1">
                <h3 className="text-xl font-bold text-gray-900 mb-1">
                  Best Market Today
                </h3>
                <p className="text-lg text-gray-700 mb-2">
                  {bestPrice.market_name}
                </p>
                <div className="flex items-baseline gap-2">
                  <span className="text-3xl font-bold text-[#1a5f3f]">
                    ₦{bestPrice.price_min.toLocaleString()} - ₦{bestPrice.price_max.toLocaleString()}
                  </span>
                  <span className="text-gray-600">{unit}</span>
                </div>
              </div>
            </div>
          </div>
        )}

        {/* Comparison Table */}
        {!pricesLoading && selectedCrop && prices.length > 0 && (
          <div className="bg-white rounded-2xl shadow-lg overflow-hidden">
            <div className="p-6 bg-gradient-to-r from-[#1a5f3f] to-[#2d8659] text-white">
              <h2 className="text-2xl font-bold">{selectedCrop}</h2>
              <p className="opacity-90">Price comparison across all markets</p>
            </div>

            {/* Desktop Table */}
            <div className="hidden md:block overflow-x-auto">
              <table className="w-full">
                <thead className="bg-gray-50 border-b border-gray-200">
                  <tr>
                    <th className="px-6 py-4 text-left text-sm font-semibold text-gray-900">Market</th>
                    <th className="px-6 py-4 text-left text-sm font-semibold text-gray-900">Price Range</th>
                    <th className="px-6 py-4 text-center text-sm font-semibold text-gray-900">Trend</th>
                    <th className="px-6 py-4 text-center text-sm font-semibold text-gray-900">Confidence</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-gray-200">
                  {prices.map((price) => {
                    const isBest = bestPrice?.id === price.id;
                    return (
                      <tr
                        key={price.id}
                        className={`hover:bg-gray-50 transition-colors ${isBest ? 'bg-amber-50' : ''}`}
                      >
                        <td className="px-6 py-4">
                          <div className="flex items-center gap-2">
                            {isBest && <Trophy className="w-5 h-5 text-amber-500" />}
                            <span className="font-medium text-gray-900">{price.market_name}</span>
                          </div>
                        </td>
                        <td className="px-6 py-4">
                          <div className="font-semibold text-[#1a5f3f]">
                            ₦{price.price_min.toLocaleString()} - ₦{price.price_max.toLocaleString()}
                          </div>
                          <div className="text-sm text-gray-600">{price.unit}</div>
                        </td>
                        <td className="px-6 py-4">
                          <div className="flex justify-center">{getTrendIcon(price.trend)}</div>
                        </td>
                        <td className="px-6 py-4">
                          <div className="flex justify-center">{getConfidenceBadge(price.confidence)}</div>
                        </td>
                      </tr>
                    );
                  })}
                </tbody>
              </table>
            </div>

            {/* Mobile Cards */}
            <div className="md:hidden divide-y divide-gray-200">
              {prices.map((price) => {
                const isBest = bestPrice?.id === price.id;
                return (
                  <div key={price.id} className={`p-5 ${isBest ? 'bg-amber-50' : ''}`}>
                    <div className="flex items-start justify-between mb-3">
                      <div>
                        <div className="flex items-center gap-2 mb-1">
                          {isBest && <Trophy className="w-5 h-5 text-amber-500" />}
                          <h3 className="font-semibold text-gray-900">{price.market_name}</h3>
                        </div>
                      </div>
                      {getConfidenceBadge(price.confidence)}
                    </div>

                    <div className="flex items-baseline gap-2 mb-2">
                      <span className="text-2xl font-bold text-[#1a5f3f]">
                        ₦{price.price_min.toLocaleString()} - ₦{price.price_max.toLocaleString()}
                      </span>
                    </div>
                    <div className="text-sm text-gray-600 mb-3">{price.unit}</div>

                    <div className="flex items-center gap-2">
                      {getTrendIcon(price.trend)}
                      <span className="text-sm text-gray-700">
                        {price.trend === 'up' && 'Rising'}
                        {price.trend === 'down' && 'Falling'}
                        {price.trend === 'stable' && 'Stable'}
                      </span>
                    </div>
                  </div>
                );
              })}
            </div>
          </div>
        )}

        {/* No Data */}
        {!pricesLoading && selectedCrop && prices.length === 0 && (
          <div className="bg-white rounded-2xl shadow-lg p-12 text-center">
            <Trophy className="w-16 h-16 text-gray-300 mx-auto mb-4" />
            <h3 className="text-xl font-semibold text-gray-900 mb-2">No Price Data</h3>
            <p className="text-gray-600">No aggregated prices found for {selectedCrop}.</p>
          </div>
        )}

        {/* Empty State */}
        {!selectedCrop && (
          <div className="bg-white rounded-2xl shadow-lg p-12 text-center">
            <Trophy className="w-16 h-16 text-gray-300 mx-auto mb-4" />
            <h3 className="text-xl font-semibold text-gray-900 mb-2">
              Select a Crop to Compare
            </h3>
            <p className="text-gray-600">
              Choose a crop above to see prices across all markets
            </p>
          </div>
        )}
      </div>
    </div>
  );
}