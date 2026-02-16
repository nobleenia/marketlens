import { useState } from 'react';
import { useCreateCrop } from '../../api/hooks';

export default function CreateCropForm() {
  const [name, setName] = useState('');
  const [unit, setUnit] = useState('kg');
  const createCrop = useCreateCrop();

  const submit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!name.trim()) return;
    try {
      await createCrop.mutateAsync({ name: name.trim(), unit });
      setName('');
      // optionally show toast or message
      alert('Crop created');
    } catch (err: any) {
      alert(err.message || 'Failed to create crop');
    }
  };

  return (
    <form onSubmit={submit} className="space-y-4">
      <div>
        <label className="block text-sm">Crop name</label>
        <input value={name} onChange={(e) => setName(e.target.value)} className="input" />
      </div>
      <div>
        <label className="block text-sm">Unit</label>
        <input value={unit} onChange={(e) => setUnit(e.target.value)} className="input" />
      </div>
      <button type="submit" className="btn" disabled={createCrop.isPending}>
        {createCrop.isPending ? 'Creating...' : 'Create Crop'}
      </button>
    </form>
  );
}