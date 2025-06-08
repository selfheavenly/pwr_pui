import { getVStopData } from "@/services/vstopService";
import { useQuery } from "@tanstack/react-query";

export function useVStopData() {
  return useQuery({
    queryKey: ["vstop"],
    queryFn: getVStopData,
  });
}
