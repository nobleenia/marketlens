export interface Crop {
    id: string;
    name: string;
    unit: string;
}

export interface Market {
    id: string;
    name: string;
    state: string;
    country: string;
    latitude: number;
    longitude: number;
}

export interface AggregatedPrice {
    id: string;
    crop_id: string;
    crop_name: string;
    market_id: string;
    market_name: string;
    period: string;
    period_start: string;
    period_end: string;
    price_min: number;
    price_max: number;
    price_mean: number;
    price_median: number;
    currency: string;
    unit: string;
    confidence: string;
    sample_size: number;
    trend: string;
    created_at: string;
    updated_at: string;
}

export interface Trend {
    date: string;
    price: number;
}

export interface PriceObservation {
  id: string;
  crop_id: string;
  crop_name: string;
  market_id: string;
  market_name: string;
  observed_at: string;
  price: number;
  currency: string;
  unit: string;
  source: string;
  reporter_id: string;
  notes: string;
  confidence_score: number;
  status: string;       // pending | approved | flagged | rejected
  created_at: string;
}

export interface AuditLog {
  id: string;
  admin_id: string;
  action: string;
  entity_type: string;
  entity_id: string;
  old_value: string;
  new_value: string;
  reason: string;
  created_at: string;
}

export interface PaginatedResponse<T> {
  data: T[];
  total: number;
  limit: number;
  offset: number;
}