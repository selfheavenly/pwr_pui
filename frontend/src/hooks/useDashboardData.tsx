import type { DashboardResponse } from "@/types/dashboard";
import { getDashboardData } from "@/services/dashboardService";
import { useQuery } from "@tanstack/react-query";

export function useDashboardData() {
  return useQuery<DashboardResponse, Error>({
    queryKey: ["dashboard"],
    queryFn: () => getDashboardData(),
  });
}
