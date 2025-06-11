import { Card, CardContent } from "@/components/ui/card";

import type { BetBrief } from "@/types/bets";
import { formatDistanceToNow } from "date-fns";
import { pl } from "date-fns/locale";

interface Props {
  bet: BetBrief;
}

export function BetCard({ bet }: Props) {
  const when = formatDistanceToNow(new Date(bet.placed_at), {
    addSuffix: true,
    locale: pl,
  });

  const statusLabel =
    bet.status === "won"
      ? "Wygrana"
      : bet.status === "lost"
        ? "Przegrana"
        : "Oczekuje";

  const statusColor =
    bet.status === "won"
      ? "text-green-500"
      : bet.status === "lost"
        ? "text-red-500"
        : "text-zinc-400"; // grey for pending

  return (
    <Card className="bg-black border border-zinc-800">
      <CardContent className="p-3 space-y-1">
        <div className="flex justify-between items-center">
          <div className="flex items-center space-x-2">
            <div className="bg-zinc-800 px-2 py-0.5 rounded text-xs">
              {bet.tram_lane_id}
            </div>
            <span className="text-sm font-medium">
              {bet.tram_lane_destination}
            </span>
          </div>
          <div className="text-right">
            <p className="text-xs text-zinc-400">{when}</p>
            <p className={`text-sm font-semibold ${statusColor}`}>
              {bet.status === "pending"
                ? statusLabel
                : `${bet.status === "won" ? "+" : "-"} ${Math.abs(bet.bet_result)} PLN`}
            </p>
          </div>
        </div>
        <div className="text-xs text-zinc-500 flex items-center space-x-4">
          <span>{bet.stop_name}</span>
          <span>⏱ {bet.actual_delay}</span>
          <span>x {bet.bet_rate}</span>
        </div>
      </CardContent>
    </Card>
  );
}
