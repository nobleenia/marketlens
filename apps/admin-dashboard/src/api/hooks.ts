import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiFetch } from "./client";
import type { Crop, Market, AggregatedPrice, PriceObservation, AuditLog, PaginatedResponse } from "./types";

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

// ── Admin hooks ─────────────────────────────────────────────────────

export function useObservations(crop?: string, market?: string, status?: string, limit = 50, offset = 0) {
  const params = new URLSearchParams();
  if (crop) params.set('crop', crop);
  if (market) params.set('market', market);
  if (status) params.set('status', status);
  params.set('limit', String(limit));
  params.set('offset', String(offset));

  return useQuery<PaginatedResponse<PriceObservation>>({
    queryKey: ['observations', crop, market, status, limit, offset],
    queryFn: () => apiFetch(`/v1/admin/observations?${params}`),
  });
}

export function useAuditLogs(entityType?: string, entityId?: string, limit = 50, offset = 0) {
  const params = new URLSearchParams();
  if (entityType) params.set('entity_type', entityType);
  if (entityId) params.set('entity_id', entityId);
  params.set('limit', String(limit));
  params.set('offset', String(offset));

  return useQuery<PaginatedResponse<AuditLog>>({
    queryKey: ['audit-logs', entityType, entityId, limit, offset],
    queryFn: () => apiFetch(`/v1/admin/audit?${params}`),
  });
}

export function useUpdateObservationStatus() {
  const queryClient = useQueryClient();

  return useMutation<
    { status: string; new_status: string },
    Error,
    { id: string; status: string; reason: string; adminId?: string }
  >({
    mutationFn: ({ id, status, reason, adminId }) =>
      apiFetch(`/v1/admin/observations/${id}`, {
        method: 'PATCH',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          status,
          reason,
          admin_id: adminId || 'admin',
        }),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['observations'] });
      queryClient.invalidateQueries({ queryKey: ['audit-logs'] });
    },
  });
}