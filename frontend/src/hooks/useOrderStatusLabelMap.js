import { useMemo } from 'react';
import { useQuery } from '@tanstack/react-query';
import { orderService } from '@/services/orderService';
import { buildOrderStatusLabelMap } from '@/lib/orderStatusDisplay';

const STALE_MS = 5 * 60 * 1000;

/**
 * Public order statuses (active): single source aligned with backend dictionary.
 * React Query dedupes by queryKey across components.
 */
export function useOrderStatusLabelMap(options = {}) {
  const { enabled = true } = options;
  const query = useQuery({
    queryKey: ['order-statuses', 'public'],
    queryFn: () => orderService.getPublicOrderStatuses(),
    staleTime: STALE_MS,
    enabled,
  });

  const labelByCode = useMemo(
    () => buildOrderStatusLabelMap(query.data),
    [query.data]
  );

  return { labelByCode, ...query };
}
