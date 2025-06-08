import { type RateOdds } from "./rates";

export interface DepartureDetails {
  tram_id: string;
  stop_id: number;
  stop_name: string;
  line: string;
  destination: string;
  arrival_time: string;
  odds: RateOdds[];
  balance: number; // user's current balance
}
