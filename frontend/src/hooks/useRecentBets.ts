import type { RecentBetsBrief } from "@/types/bets";
import { getRecentBetsBrief } from "@/services/dashboardService";
import { useQuery } from "@tanstack/react-query";

export function useRecentBets(page = 0, size = 2) {
  return useQuery<RecentBetsBrief>({
    queryKey: ["recentBets", page, size],
    queryFn: () => getRecentBetsBrief(page, size),
  });
}
