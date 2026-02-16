import { useState } from 'react';
import { useCreateMarket } from '../../api/hooks';

export default function CreateMarketForm() {
  const [name, setName] = useState('');
  const [state, setState] = useState('');
  const [country, setCountry] = useState('NG');
  const [lat, setLat] = useState('');
  const [lng, setLng] = useState('');
  const createMarket = useCreateMarket();

  const submit = async (e: React.FormEvent) => {
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
      setName('');
      setState('');
      setLat('');
      setLng('');
      alert('Market created');
    } catch (err: any) {
      alert(err.message || 'Failed to create market');
    }
  };

  return (
    <form onSubmit={submit} className="space-y-4">
      <div>
        <label className="block text-sm">Market name</label>
        <input value={name} onChange={(e) => setName(e.target.value)} className="input" />
      </div>
      <div>
        <label className="block text-sm">State</label>
        <input value={state} onChange={(e) => setState(e.target.value)} className="input" />
      </div>
      <div>
        <label className="block text-sm">Country</label>
        <input value={country} onChange={(e) => setCountry(e.target.value)} className="input" />
      </div>
      <div className="grid grid-cols-2 gap-2">
        <div>
          <label className="block text-sm">Latitude</label>
          <input value={lat} onChange={(e) => setLat(e.target.value)} className="input" />
        </div>
        <div>
          <label className="block text-sm">Longitude</label>
          <input value={lng} onChange={(e) => setLng(e.target.value)} className="input" />
        </div>
      </div>
      <button type="submit" className="btn" disabled={createMarket.isPending}>
        {createMarket.isPending ? 'Creating...' : 'Create Market'}
      </button>
    </form>
  );
}