import type { User } from "@/types/user";
import { getUserBalance } from "@/services/dashboardService";
import { useQuery } from "@tanstack/react-query";

export function useUserBalance() {
  return useQuery<User>({
    queryKey: ["userBalance"],
    queryFn: getUserBalance,
  });
}
