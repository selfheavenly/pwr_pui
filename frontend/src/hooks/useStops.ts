import type { StopSummary } from "@/types/stops";
import { getStops } from "@/services/dashboardService";
import { useQuery } from "@tanstack/react-query";

export function useStops() {
  return useQuery<StopSummary[]>({
    queryKey: ["stops"],
    queryFn: getStops,
  });
}
