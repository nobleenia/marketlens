import { useState } from 'react';
import { useNavigate } from 'react-router';
import { Search, TrendingUp, TrendingDown, Minus, ArrowRight } from 'lucide-react';
import { useCrops, useMarkets, usePrice } from '../api/hooks';
import type { AggregatedPrice } from '../api/types';

/** ── Price result card (shown after "Check Price") ─────────────────── */
function PriceResultCard({
  price,
  onViewTrend,
  onCompare,
}: {
  price: AggregatedPrice;
  onViewTrend: () => void;
  onCompare: () => void;
}) {
  const trendIcon =
    price.trend === 'up'   ? <TrendingUp className="w-5 h-5 text-green-600" /> :
    price.trend === 'down' ? <TrendingDown className="w-5 h-5 text-red-600" /> :
    <Minus className="w-5 h-5 text-gray-500" />;

  return (
    <div className="bg-white rounded-2xl shadow-xl overflow-hidden">
      <div className="bg-gradient-to-r from-[#1a5f3f] to-[#2d8659] text-white p-6">
        <h2 className="text-2xl font-bold">{price.crop_name}</h2>
        <p className="opacity-90">{price.market_name}</p>
      </div>
      <div className="p-6 space-y-4">
        <div className="flex items-baseline gap-2">
          <span className="text-3xl font-bold text-[#1a5f3f]">
            ₦{price.price_min.toLocaleString()} – ₦{price.price_max.toLocaleString()}
          </span>
          <span className="text-gray-500 text-sm">{price.unit}</span>
        </div>
        <div className="flex items-center gap-2 text-sm text-gray-600">
          {trendIcon}
          <span>
            {price.trend === 'up' ? 'Rising' : price.trend === 'down' ? 'Falling' : 'Stable'}
          </span>
          <span className="mx-2">•</span>
          <span className={`px-2 py-0.5 rounded-full text-xs font-medium ${
            price.confidence === 'high'   ? 'bg-green-100 text-green-700' :
            price.confidence === 'medium' ? 'bg-amber-100 text-amber-700' :
            'bg-red-100 text-red-700'
          }`}>
            {price.confidence.charAt(0).toUpperCase() + price.confidence.slice(1)} confidence
          </span>
        </div>
        <div className="flex gap-3 pt-2">
          <button onClick={onViewTrend} className="flex-1 flex items-center justify-center gap-2 px-4 py-3 bg-[#1a5f3f] text-white rounded-lg hover:bg-[#154d33] transition-colors font-medium">
            View Trend <ArrowRight className="w-4 h-4" />
          </button>
          <button onClick={onCompare} className="flex-1 flex items-center justify-center gap-2 px-4 py-3 border-2 border-[#1a5f3f] text-[#1a5f3f] rounded-lg hover:bg-[#e8f5e9] transition-colors font-medium">
            Compare Markets
          </button>
        </div>
      </div>
    </div>
  );
}

/** ── Home page ─────────────────────────────────────────────────────── */
export default function Home() {
  const navigate = useNavigate();

  // ── API data ──────────────────────────────────────────────────────
  const { data: crops = [], isLoading: cropsLoading } = useCrops();
  const { data: markets = [], isLoading: marketsLoading } = useMarkets();

  // ── Local state ───────────────────────────────────────────────────
  // Dropdowns store NAMES because the price endpoint uses names
  const [selectedCrop, setSelectedCrop] = useState('');
  const [selectedMarket, setSelectedMarket] = useState('');
  const [checkNow, setCheckNow] = useState(false);

  // ── Price query — only fires when checkNow is true ────────────────
  const {
    data: priceResult,
    isLoading: priceLoading,
    error: priceError,
  } = usePrice(
    checkNow ? selectedCrop : undefined,
    checkNow ? selectedMarket : undefined,
  );

  const handleCheckPrice = () => {
    if (selectedCrop && selectedMarket) {
      setCheckNow(true);
    }
  };

  const handleViewTrend = () => {
    if (selectedCrop && selectedMarket) {
      navigate(`/trends?crop=${encodeURIComponent(selectedCrop)}&market=${encodeURIComponent(selectedMarket)}`);
    }
  };

  const handleCompare = () => {
    if (selectedCrop) {
      navigate(`/compare?crop=${encodeURIComponent(selectedCrop)}`);
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-b from-[#e8f5e9] to-white pb-20 md:pb-8">
      <div className="container mx-auto px-4 py-8">
        {/* Hero Section */}
        <div className="text-center mb-8 md:mb-12">
          <h1 className="text-3xl md:text-4xl lg:text-5xl font-bold text-[#1a5f3f] mb-3">
            Check Today's Market Prices
          </h1>
          <p className="text-gray-700 text-lg md:text-xl">
            Get real-time agricultural commodity prices across Nigeria
          </p>
        </div>

        {/* Price Checker Form */}
        <div className="max-w-2xl mx-auto bg-white rounded-2xl shadow-xl p-6 md:p-8 mb-8">
          <div className="space-y-5">
            {/* Crop Selector */}
            <div>
              <label htmlFor="crop" className="block text-sm font-medium text-gray-700 mb-2">
                Select Crop
              </label>
              <select
                id="crop"
                value={selectedCrop}
                onChange={(e) => {
                  setSelectedCrop(e.target.value);
                  setCheckNow(false);
                }}
                disabled={cropsLoading}
                className="w-full px-4 py-3 border-2 border-gray-300 rounded-lg focus:border-[#1a5f3f] focus:outline-none text-base"
              >
                <option value="">{cropsLoading ? 'Loading crops...' : 'Choose a crop...'}</option>
                {crops.map((crop) => (
                  <option key={crop.id} value={crop.name}>
                    {crop.name}
                  </option>
                ))}
              </select>
            </div>

            {/* Market Selector */}
            <div>
              <label htmlFor="market" className="block text-sm font-medium text-gray-700 mb-2">
                Select Market
              </label>
              <select
                id="market"
                value={selectedMarket}
                onChange={(e) => {
                  setSelectedMarket(e.target.value);
                  setCheckNow(false);
                }}
                disabled={marketsLoading}
                className="w-full px-4 py-3 border-2 border-gray-300 rounded-lg focus:border-[#1a5f3f] focus:outline-none text-base"
              >
                <option value="">{marketsLoading ? 'Loading markets...' : 'Choose a market...'}</option>
                {markets.map((market) => (
                  <option key={market.id} value={market.name}>
                    {market.name} ({market.state})
                  </option>
                ))}
              </select>
            </div>

            {/* Check Price Button */}
            <button
              onClick={handleCheckPrice}
              disabled={!selectedCrop || !selectedMarket || priceLoading}
              className="w-full bg-[#1a5f3f] text-white px-6 py-4 rounded-lg hover:bg-[#154d33] transition-colors font-semibold text-lg disabled:bg-gray-300 disabled:cursor-not-allowed flex items-center justify-center gap-3"
            >
              <Search className="w-6 h-6" />
              {priceLoading ? 'Checking...' : 'Check Price'}
            </button>
          </div>
        </div>

        {/* Price Error */}
        {priceError && (
          <div className="max-w-2xl mx-auto mb-6 bg-red-50 border border-red-200 text-red-700 rounded-lg p-4 text-center">
            Price not available for this combination. Try a different crop/market.
          </div>
        )}

        {/* Price Result */}
        {checkNow && priceResult && (
          <div className="max-w-3xl mx-auto animate-fadeIn">
            <PriceResultCard
              price={priceResult}
              onViewTrend={handleViewTrend}
              onCompare={handleCompare}
            />
          </div>
        )}

        {/* Quick Info Cards */}
        {!checkNow && (
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4 max-w-5xl mx-auto mt-12">
            <div className="bg-white rounded-xl p-6 shadow-md border border-gray-100">
              <div className="bg-[#e8f5e9] rounded-full w-12 h-12 flex items-center justify-center mb-4">
                <Search className="w-6 h-6 text-[#1a5f3f]" />
              </div>
              <h3 className="font-semibold text-gray-900 mb-2">Real-Time Prices</h3>
              <p className="text-sm text-gray-600">
                Updated daily from markets across Nigeria
              </p>
            </div>

            <div className="bg-white rounded-xl p-6 shadow-md border border-gray-100">
              <div className="bg-[#e8f5e9] rounded-full w-12 h-12 flex items-center justify-center mb-4">
                <Search className="w-6 h-6 text-[#1a5f3f]" />
              </div>
              <h3 className="font-semibold text-gray-900 mb-2">Price Trends</h3>
              <p className="text-sm text-gray-600">
                Track price movements over time
              </p>
            </div>

            <div className="bg-white rounded-xl p-6 shadow-md border border-gray-100">
              <div className="bg-[#e8f5e9] rounded-full w-12 h-12 flex items-center justify-center mb-4">
                <Search className="w-6 h-6 text-[#1a5f3f]" />
              </div>
              <h3 className="font-semibold text-gray-900 mb-2">Market Compare</h3>
              <p className="text-sm text-gray-600">
                Find the best prices across locations
              </p>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}