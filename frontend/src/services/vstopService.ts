import type { StopDetails } from "@/types/stops";
import type { TripDetails } from "@/types/trips";

export async function getVStopData(id: string): Promise<StopDetails> {
  const data = await import(`@/mock/vstop.json`);
  return data.default as StopDetails;
}

export async function getTripOdds(tripID: string): Promise<TripDetails> {
  const data = await import(`@/mock/vstop-trips/${tripID}.json`);
  return data.default as TripDetails;
}
