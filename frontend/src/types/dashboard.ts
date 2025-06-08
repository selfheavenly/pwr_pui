export interface LastBet {
  line: string;
  destination: string;
  stop: string;
  time: string; // e.g. "0:42"
  odds: number; // e.g. 1.78
  amount: number; // e.g. 890
  when: string; // e.g. "wczoraj"
}

export interface VirtualStop {
  name: string;
  lines: number[]; // e.g. [18, 21, 22]
}

export interface DashboardResponse {
  balance: string; // e.g. "45,231.89"
  changePercent: string; // e.g. "20.1%"
  lastBets: LastBet[];
  virtualStops: VirtualStop[];
}
