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