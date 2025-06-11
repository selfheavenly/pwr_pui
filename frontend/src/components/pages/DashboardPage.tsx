import { BetsTable } from "../dashboard/DashboardBetsTable";
import { DashboardBalance } from "@/components/dashboard/DashboardBalance";
import { DashboardBets } from "@/components/dashboard/DashboardBets";
import { DashboardHeader } from "@/components/dashboard/DashboardHeader";
import { DashboardStops } from "@/components/dashboard/DashboardStops";
import { useRecentBets } from "@/hooks/useRecentBets";
import { useStops } from "@/hooks/useStops";
import { useUserBalance } from "@/hooks/useUserBalance";

export default function DashboardPage() {
  const { data: user, isLoading: userLoading } = useUserBalance();
  const { data: betsData, isLoading: betsLoading } = useRecentBets(0, 2);
  const { data: stops, isLoading: stopsLoading } = useStops();

  const isLoading = userLoading || betsLoading || stopsLoading;

  if (isLoading || !user || !betsData || !stops) {
    return <div className="text-white">Loading...</div>;
  }

  return (
    <div className="relative min-h-screen text-white px-4 py-6 space-y-6 overflow-hidden bg-black">
      <DashboardHeader />
      <DashboardBalance
        balance={user.balance}
        changePercent={user.ytd_change_percent}
      />
      <DashboardBets lastBets={betsData.data} />
      <BetsTable />
      <DashboardStops virtualStops={stops} />
    </div>
  );
}
