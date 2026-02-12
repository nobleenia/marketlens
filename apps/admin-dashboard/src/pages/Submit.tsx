import { useState } from 'react';
import { Send, CheckCircle, AlertCircle, Loader2 } from 'lucide-react';
import { useCrops, useMarkets } from '../api/hooks';
import { apiFetch } from '../api/client';

export default function Submit() {
  const { data: crops = [], isLoading: cropsLoading } = useCrops();
  const { data: markets = [], isLoading: marketsLoading } = useMarkets();

  const [cropName, setCropName] = useState('');
  const [marketName, setMarketName] = useState('');
  const [price, setPrice] = useState('');
  const [unit, setUnit] = useState('kg');
  const [notes, setNotes] = useState('');
  const [submitting, setSubmitting] = useState(false);
  const [result, setResult] = useState<'success' | 'error' | null>(null);
  const [errorMsg, setErrorMsg] = useState('');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!cropName || !marketName || !price) return;

    setSubmitting(true);
    setResult(null);
    setErrorMsg('');

    try {
      await apiFetch('/v1/observations', {
        method: 'POST',
        body: JSON.stringify({
          crop_name: cropName,
          market_name: marketName,
          price: parseFloat(price),
          currency: 'NGN',
          unit,
          source: 'web',
          notes: notes || undefined,
        }),
      });
      setResult('success');
      setPrice('');
      setNotes('');
    } catch (err) {
      setResult('error');
      setErrorMsg(err instanceof Error ? err.message : 'Submission failed');
    } finally {
      setSubmitting(false);
    }
  };

  const selectedCropUnit = crops.find((c) => c.name === cropName)?.unit;
  if (selectedCropUnit && selectedCropUnit !== unit) {
    setUnit(selectedCropUnit);
  }

  return (
    <div className="min-h-screen bg-gradient-to-b from-[#e8f5e9] to-white pb-20 md:pb-8">
      <div className="container mx-auto px-4 py-8 max-w-2xl">
        <div className="mb-8">
          <h1 className="text-3xl md:text-4xl font-bold text-[#1a5f3f] mb-2">
            Submit a Price
          </h1>
          <p className="text-gray-700 text-lg">
            Help fellow farmers by reporting current market prices
          </p>
        </div>

        <form onSubmit={handleSubmit} className="bg-white rounded-2xl shadow-lg p-6 space-y-6">
          {/* Crop */}
          <div>
            <label htmlFor="crop" className="block text-sm font-medium text-gray-700 mb-2">
              Crop *
            </label>
            <select
              id="crop"
              value={cropName}
              onChange={(e) => setCropName(e.target.value)}
              disabled={cropsLoading}
              required
              className="w-full px-4 py-3 border-2 border-gray-300 rounded-lg focus:border-[#1a5f3f] focus:outline-none"
            >
              <option value="">{cropsLoading ? 'Loading...' : 'Select a crop'}</option>
              {crops.map((c) => (
                <option key={c.id} value={c.name}>{c.name}</option>
              ))}
            </select>
          </div>

          {/* Market */}
          <div>
            <label htmlFor="market" className="block text-sm font-medium text-gray-700 mb-2">
              Market *
            </label>
            <select
              id="market"
              value={marketName}
              onChange={(e) => setMarketName(e.target.value)}
              disabled={marketsLoading}
              required
              className="w-full px-4 py-3 border-2 border-gray-300 rounded-lg focus:border-[#1a5f3f] focus:outline-none"
            >
              <option value="">{marketsLoading ? 'Loading...' : 'Select a market'}</option>
              {markets.map((m) => (
                <option key={m.id} value={m.name}>{m.name} ({m.state})</option>
              ))}
            </select>
          </div>

          {/* Price + Unit */}
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label htmlFor="price" className="block text-sm font-medium text-gray-700 mb-2">
                Price (₦) *
              </label>
              <input
                id="price"
                type="number"
                min="1"
                step="any"
                value={price}
                onChange={(e) => setPrice(e.target.value)}
                required
                placeholder="e.g. 5000"
                className="w-full px-4 py-3 border-2 border-gray-300 rounded-lg focus:border-[#1a5f3f] focus:outline-none"
              />
            </div>
            <div>
              <label htmlFor="unit" className="block text-sm font-medium text-gray-700 mb-2">
                Unit
              </label>
              <select
                id="unit"
                value={unit}
                onChange={(e) => setUnit(e.target.value)}
                className="w-full px-4 py-3 border-2 border-gray-300 rounded-lg focus:border-[#1a5f3f] focus:outline-none"
              >
                <option value="kg">kg</option>
                <option value="basket">basket</option>
                <option value="tuber">tuber</option>
                <option value="bag">bag (100kg)</option>
                <option value="mudu">mudu</option>
              </select>
            </div>
          </div>

          {/* Notes */}
          <div>
            <label htmlFor="notes" className="block text-sm font-medium text-gray-700 mb-2">
              Notes (optional)
            </label>
            <textarea
              id="notes"
              value={notes}
              onChange={(e) => setNotes(e.target.value)}
              placeholder="e.g. Morning price, retail..."
              rows={3}
              className="w-full px-4 py-3 border-2 border-gray-300 rounded-lg focus:border-[#1a5f3f] focus:outline-none resize-none"
            />
          </div>

          {/* Result messages */}
          {result === 'success' && (
            <div className="flex items-center gap-3 bg-green-50 border border-green-200 rounded-lg p-4">
              <CheckCircle className="w-5 h-5 text-green-600 shrink-0" />
              <p className="text-green-700 font-medium">Price submitted successfully! It will be reviewed by an admin.</p>
            </div>
          )}
          {result === 'error' && (
            <div className="flex items-center gap-3 bg-red-50 border border-red-200 rounded-lg p-4">
              <AlertCircle className="w-5 h-5 text-red-600 shrink-0" />
              <p className="text-red-700 font-medium">{errorMsg}</p>
            </div>
          )}

          {/* Submit */}
          <button
            type="submit"
            disabled={submitting || !cropName || !marketName || !price}
            className="w-full py-3 bg-[#1a5f3f] text-white rounded-lg font-medium hover:bg-[#155232] transition-colors disabled:opacity-50 flex items-center justify-center gap-2"
          >
            {submitting ? (
              <><Loader2 className="w-5 h-5 animate-spin" /> Submitting...</>
            ) : (
              <><Send className="w-5 h-5" /> Submit Price</>
            )}
          </button>
        </form>

        <p className="text-center text-sm text-gray-500 mt-6">
          Your submission will be reviewed before being included in price aggregations.
        </p>
      </div>
    </div>
  );
}