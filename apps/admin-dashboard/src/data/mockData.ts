// Temporary mock data for admin features until Sprint 1.10 backend endpoints are built.
// Only priceSubmissions remains — crops, markets, and priceData are now served by the API.

export interface PriceSubmission {
  id: string;
  crop: string;
  market: string;
  minPrice: number;
  maxPrice: number;
  source: string;
  status: 'pending' | 'approved' | 'flagged';
  timestamp: string;
}

export const priceSubmissions: PriceSubmission[] = [
  {
    id: 'sub-1',
    crop: 'Tomatoes',
    market: 'Bodija Market',
    minPrice: 3500,
    maxPrice: 4200,
    source: 'USSD',
    status: 'approved',
    timestamp: '2 hours ago',
  },
  {
    id: 'sub-2',
    crop: 'Rice',
    market: 'Mile 12 Market',
    minPrice: 72000,
    maxPrice: 78000,
    source: 'USSD',
    status: 'pending',
    timestamp: '3 hours ago',
  },
  {
    id: 'sub-3',
    crop: 'Maize',
    market: 'Dawanau Market',
    minPrice: 28000,
    maxPrice: 32000,
    source: 'Agent',
    status: 'approved',
    timestamp: '5 hours ago',
  },
  {
    id: 'sub-4',
    crop: 'Yam',
    market: 'Oyingbo Market',
    minPrice: 15000,
    maxPrice: 22000,
    source: 'USSD',
    status: 'flagged',
    timestamp: '6 hours ago',
  },
  {
    id: 'sub-5',
    crop: 'Beans',
    market: 'Yankaba Market',
    minPrice: 55000,
    maxPrice: 62000,
    source: 'Agent',
    status: 'pending',
    timestamp: '8 hours ago',
  },
  {
    id: 'sub-6',
    crop: 'Onions',
    market: 'Bodija Market',
    minPrice: 8000,
    maxPrice: 12000,
    source: 'USSD',
    status: 'approved',
    timestamp: '10 hours ago',
  },
];