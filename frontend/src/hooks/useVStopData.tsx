import { getTripOdds, getVStopData } from "@/services/vstopService";

import type { TripDetails } from "@/types/trips";
import { useQuery } from "@tanstack/react-query";

export function useVStopData(id: string) {
  return useQuery({
    queryKey: ["vstop", id],
    queryFn: () => getVStopData(id),
  });
}

export function useTripOdds(tripID: string) {
  return useQuery<TripDetails>({
    queryKey: ["tripOdds", tripID],
    queryFn: () => getTripOdds(tripID),
  });
}
