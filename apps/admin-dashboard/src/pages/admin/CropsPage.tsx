import { Loader2, Sprout } from 'lucide-react';
import { useCrops } from '../../api/hooks';

export default function CropsPage() {
  const { data: crops = [], isLoading, isError } = useCrops();

  return (
    <div className="p-6 md:p-8">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900 mb-2">Crops</h1>
        <p className="text-gray-600">All tracked crops in the system</p>
      </div>

      {isLoading && (
        <div className="bg-white rounded-xl shadow-lg p-12 text-center">
          <Loader2 className="w-10 h-10 text-[#1a5f3f] mx-auto mb-3 animate-spin" />
          <p className="text-gray-600">Loading crops...</p>
        </div>
      )}

      {isError && (
        <div className="bg-red-50 border border-red-200 rounded-xl p-6 text-center">
          <p className="text-red-700 font-medium">Failed to load crops.</p>
        </div>
      )}

      {!isLoading && !isError && (
        <div className="bg-white rounded-xl shadow-lg overflow-hidden">
          <div className="p-6 border-b border-gray-200 flex items-center justify-between">
            <h2 className="text-lg font-semibold text-gray-900">
              {crops.length} crop{crops.length !== 1 ? 's' : ''} registered
            </h2>
          </div>

          {/* Desktop Table */}
          <div className="hidden md:block overflow-x-auto">
            <table className="w-full">
              <thead className="bg-gray-50 border-b border-gray-200">
                <tr>
                  <th className="px-6 py-4 text-left text-xs font-semibold text-gray-700 uppercase tracking-wider">Name</th>
                  <th className="px-6 py-4 text-left text-xs font-semibold text-gray-700 uppercase tracking-wider">Unit</th>
                  <th className="px-6 py-4 text-left text-xs font-semibold text-gray-700 uppercase tracking-wider">ID</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-200">
                {crops.map((crop) => (
                  <tr key={crop.id} className="hover:bg-gray-50">
                    <td className="px-6 py-4">
                      <div className="flex items-center gap-3">
                        <div className="bg-green-100 rounded-lg p-2">
                          <Sprout className="w-4 h-4 text-green-600" />
                        </div>
                        <span className="font-medium text-gray-900">{crop.name}</span>
                      </div>
                    </td>
                    <td className="px-6 py-4 text-gray-700">{crop.unit}</td>
                    <td className="px-6 py-4 text-xs text-gray-400 font-mono">{crop.id.slice(0, 8)}…</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>

          {/* Mobile Cards */}
          <div className="md:hidden divide-y divide-gray-200">
            {crops.map((crop) => (
              <div key={crop.id} className="p-5 flex items-center gap-4">
                <div className="bg-green-100 rounded-lg p-3">
                  <Sprout className="w-5 h-5 text-green-600" />
                </div>
                <div>
                  <p className="font-semibold text-gray-900">{crop.name}</p>
                  <p className="text-sm text-gray-500">Unit: {crop.unit}</p>
                </div>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}