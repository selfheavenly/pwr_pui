// const BACKEND_URL = import.meta.env.VITE_BACKEND_URL || "http://127.0.0.1:8000";

import type { BetBrief, RecentBetsBrief } from "@/types/bets";

import type { StopSummary } from "@/types/stops";
import type { User } from "@/types/user";
// Import the mock data directly
import mockBetsData from "@/mock/bets.json";
import mockStopsData from "@/mock/stops.json"; // Assuming this path is correct and it's an array of StopSummary

// Assuming this path is correct and it's an array of BetBrief

// /api/user/me - User profile & balance
export async function getUserBalance(): Promise<User> {
  // This remains a mock for now, as requested.
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
  page: number, // The current page number (0-indexed)
  size: number // The number of items per page
): Promise<RecentBetsBrief> {
  // Use the imported mock data as the source
  const allBets: BetBrief[] = mockBetsData as BetBrief[]; // Cast to BetBrief[] for type safety

  const total = allBets.length; // Total number of available bets from the mock data
  const start = page * size; // Calculate the starting index for the current page
  const end = Math.min(start + size, total); // Calculate the ending index, ensuring it doesn't exceed total

  // Slice the array to get only the bets for the current page
  const betsForPage: BetBrief[] = allBets.slice(start, end);

  return {
    data: betsForPage, // The bets for the requested page
    page: page, // The current page number
    next_page: end < total, // True if there are more pages after this one
    page_count: Math.ceil(total / size), // Total number of pages
  };
}

// /api/stops - Virtual stop list
export async function getStops(): Promise<StopSummary[]> {
  // Return the imported mock stops data directly
  return mockStopsData as StopSummary[]; // Cast to StopSummary[] for type safety
}
