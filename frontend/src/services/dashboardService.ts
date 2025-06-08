// const BACKEND_URL = import.meta.env.VITE_BACKEND_URL || "http://127.0.0.1:8000";

import type { BetBrief, RecentBetsBrief } from "@/types/bets";

import type { StopSummary } from "@/types/stops";
import type { User } from "@/types/user";

// /api/user/me - User profile & balance
export async function getUserBalance(): Promise<User> {
  return {
    user_id: 101,
    google_id: "google-oauth-id",
    email: "user@example.com",
    name: "Jan Kowalski",
    balance: 45231.89,
    ytd_change_percent: 20.1,
  };
}
// /api/bets - Recent bet overviews

export async function getRecentBetsBrief(
  page: number,
  size: number
): Promise<RecentBetsBrief> {
  const total = 32; // simulate a backend total count
  const start = page * size;
  const end = Math.min(start + size, total);

  const bets: BetBrief[] = Array.from({ length: end - start }, (_, i) => {
    const index = start + i;

    const statusOptions: BetBrief["status"][] = ["won", "lost", "pending"];
    const status = statusOptions[index % statusOptions.length];

    const amount = 100 + (index % 5) * 50; // varied amount
    const rate = parseFloat((1.4 + ((index * 13) % 6) * 0.2).toFixed(2)); // more variance
    let result: number;

    if (status === "won") result = parseFloat((amount * rate).toFixed(2));
    else if (status === "lost") result = -amount;
    else result = 0;

    const tramId = ["03", "18", "9", "10", "7", "0"][index % 6];
    const destinations = [
      "Księże Małe",
      "Gaj",
      "Leśnica",
      "Krzyki",
      "Biskupin",
      "Nowy Dwór",
      "Borek",
    ];
    const stops = ["Rynek", "Grunwald", "Dominikańska", "Plac JP2"];

    return {
      bet_id: index + 1,
      bet_amount: amount,
      bet_rate: rate,
      bet_result: result,
      placed_at: new Date(Date.now() - index * 86400000).toISOString(),
      status,
      tram_lane_id: tramId,
      tram_lane_destination: destinations[index % destinations.length],
      stop_id: 100 + index,
      stop_name: stops[index % stops.length],
      actual_delay: `${Math.floor(Math.random() * 5) + 1}min`,
    };
  });

  return {
    data: bets,
    page,
    next_page: end < total,
    page_count: Math.ceil(total / size),
  };
}

// /api/stops - Virtual stop list
export async function getStops(): Promise<StopSummary[]> {
  return [
    {
      stop_id: 1,
      stop_name: "Pl. Jana Pawła II",
      lines: ["18", "21", "22"],
    },
    {
      stop_id: 2,
      stop_name: "Rynek",
      lines: ["18", "21", "22"],
    },
    {
      stop_id: 3,
      stop_name: "Pasaż Grunwaldzki",
      lines: ["18", "21", "22"],
    },
    {
      stop_id: 4,
      stop_name: "Galeria Dominikańska",
      lines: ["18", "21", "22"],
    },
  ];
}
