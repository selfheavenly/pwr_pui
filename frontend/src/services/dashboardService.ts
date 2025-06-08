// const BACKEND_URL = import.meta.env.VITE_BACKEND_URL || "http://127.0.0.1:8000";

import type { DashboardResponse } from "@/types/dashboard";

// Mock data for development
const dashboard_data: DashboardResponse = {
  balance: "45,231.89",
  changePercent: "20.1%",
  lastBets: [
    {
      line: "03",
      destination: "Księże Małe",
      stop: "Rynek",
      time: "0:42",
      odds: 1.78,
      amount: 890,
      when: "wczoraj",
    },
    {
      line: "18",
      destination: "Gaj",
      stop: "Rynek",
      time: "0:51",
      odds: 1.78,
      amount: 1780,
      when: "wczoraj",
    },
  ],
  virtualStops: [
    { name: "Pl. Jana Pawła II", lines: [18, 21, 22] },
    { name: "Rynek", lines: [18, 21, 22] },
    { name: "Pasaż Grunwaldzki", lines: [18, 21, 22] },
    { name: "Galeria Dominikańska", lines: [18, 21, 22] },
  ],
};

/**
 * Fetches Dashboard data from the backend or returns mock data.
 */
export async function getDashboardData(): Promise<DashboardResponse> {
  // In real use: switch to live API by removing return line below
  return dashboard_data;

  // Example of real request (commented out until needed)
  /*
  try {
    const response = await fetch(`${BACKEND_URL}/dashboard`);
    
    if (!response.ok) {
      const errorData = await response.json().catch(() => null);
      throw new Error(
        `Failed to fetch dashboard data: ${errorData?.detail || response.statusText}`
      );
    }

    return await response.json();
  } catch (error) {
    console.error("Error in getDashboardData:", error);
    throw error;
  }
  */
}
