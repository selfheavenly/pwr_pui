export type RecentBetsBrief = {
  data: BetBrief[];
  page: number;
  next_page: boolean;
  page_count: number;
};

export interface BetBrief {
  bet_id: number; // Unique Bet ID
  bet_amount: number; // Amount placed
  bet_rate: number; // 1.4, 1.6, etc.
  bet_result: number; // Win or Loss amount
  placed_at: string; // UTC timestamp
  status: "won" | "lost" | "pending";
  tram_lane_id: string; // Tram Lane Number
  tram_lane_destination: string; // Tram Lane Destination
  stop_id: number; // Virtual Stop ID
  stop_name: string; // Virtual Stop Name
  actual_delay: string; // Actual delay
}

// export interface BetDetails extends Bet {
//   line: string;
//   destination: string;
//   stop: string;
//   time: string;
//   odds: number;
//   amount: number;
//   when: string;
//   placed_at: string;
// }

export interface PlaceBetRequest {
  tram_lane_id: string;
  stop_id: number;
  odds: number;
  amount: number;
  interval: string;
}

export interface PlaceBetResponse {
  success: boolean;
  new_balance: number;
  bet_id: number;
}
