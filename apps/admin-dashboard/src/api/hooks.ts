import { useQuery } from "@tanstack/react-query";
import { apiFetch } from "./client";
import type { Crop, Market, AggregatedPrice } from "./types";

// ── Fetch all crops (used by every page's dropdown) ──────────────────
export function useCrops() {
  return useQuery({
    queryKey: ["crops"],
    queryFn: () => apiFetch<Crop[]>("/v1/crops"),
    staleTime: 10 * 60 * 1000, // crops rarely change — cache 10 min
  });
}

// ── Fetch all markets (used by every page's dropdown) ────────────────
export function useMarkets() {
  return useQuery({
    queryKey: ["markets"],
    queryFn: () => apiFetch<Market[]>("/v1/markets"),
    staleTime: 10 * 60 * 1000,
  });
}

// ── Fetch aggregated prices with optional filters ────────────────────
// GET /v1/prices?crop=Tomatoes&market=Bodija+Market&from=2026-01-01&to=2026-02-11
//
// Used by:
//   - Trends.tsx  → usePrices(cropName, marketName, fromDate, toDate)
//   - Compare.tsx → usePrices(cropName) to get all markets for one crop
export function usePrices(
  crop?: string,
  market?: string,
  from?: string,
  to?: string
) {
  const params = new URLSearchParams();
  if (crop) params.set("crop", crop);
  if (market) params.set("market", market);
  if (from) params.set("from", from);
  if (to) params.set("to", to);
  const qs = params.toString();

  return useQuery({
    queryKey: ["prices", crop, market, from, to],
    queryFn: () =>
      apiFetch<AggregatedPrice[]>(`/v1/prices${qs ? `?${qs}` : ""}`),
    enabled: !!crop, // don't fire until at least a crop is selected
  });
}

// ── Fetch single latest price for a crop+market pair ─────────────────
// GET /v1/prices/Tomatoes/Bodija%20Market
//
// Used by Home.tsx after user clicks "Check Price"
export function usePrice(cropName?: string, marketName?: string) {
  return useQuery({
    queryKey: ["price", cropName, marketName],
    queryFn: () =>
      apiFetch<AggregatedPrice>(
        `/v1/prices/${encodeURIComponent(cropName!)}/${encodeURIComponent(marketName!)}`
      ),
    enabled: !!cropName && !!marketName, // only fire when both selected
  });
}