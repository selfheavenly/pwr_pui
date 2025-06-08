import { type BetDetails } from "./bets";

export interface VirtualStop {
  name: string;
  lines: string[];
}

export interface DashboardData {
  balance: number;
  changePercent: string;
  lastBets: BetDetails[];
  virtualStops: VirtualStop[];
}
