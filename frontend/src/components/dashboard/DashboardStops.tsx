import { Card } from "@/components/ui/card";
import { Link } from "@tanstack/react-router";
import type { StopSummary } from "@/types/stops";

interface Props {
  virtualStops: StopSummary[];
}

export function DashboardStops({ virtualStops }: Props) {
  return (
    <div>
      <h3 className="font-semibold mb-1">Wirtualne przystanki</h3>
      <p className="text-zinc-400 text-sm mb-3">
        Wybierz przystanek, aby obstawić kolejne zdarzenia wrocławskiego MPK.
      </p>
      <div className="grid grid-cols-2 gap-3">
        {virtualStops.map((stop) => (
          <Link
            key={stop.stop_id}
            to="/vstop/$id"
            params={{ id: String(stop.stop_id) }}
            className="block"
          >
            <Card className="bg-black p-4 space-y-2 border border-zinc-800 hover:border-white/20 transition">
              <div className="font-medium">{stop.stop_name}</div>
              <div className="flex flex-wrap gap-2">
                {stop.lines.map((line, i) => (
                  <div
                    key={i}
                    className="bg-zinc-800 px-2 py-1 rounded-md text-xs"
                  >
                    {line}
                  </div>
                ))}
              </div>
            </Card>
          </Link>
        ))}
      </div>
    </div>
  );
}
