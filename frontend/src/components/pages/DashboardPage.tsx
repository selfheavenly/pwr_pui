import { Card, CardContent } from "@/components/ui/card";

import { Button } from "@/components/ui/button";
import { useDashboardData } from "@/hooks/useDashboardData";

export default function DashboardPage() {
  const { data, isLoading } = useDashboardData();

  if (isLoading) return <div className="text-white">Loading...</div>;

  return (
    <div className="min-h-screen bg-black text-white px-4 py-6 space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-xl font-semibold">MPKBet</h1>
        <div className="h-10 w-10 rounded-full bg-white/10 flex items-center justify-center">
          <span className="text-white">👤</span>
        </div>
      </div>

      <div className="text-center">
        <h2 className="text-4xl font-bold">{data?.balance} PLN</h2>
        <p className="text-green-500">+{data?.changePercent} (YTD)</p>
      </div>

      <div>
        <h3 className="font-semibold mb-2">Ostatnie zakłady</h3>
        <div className="space-y-3">
          {data?.lastBets.map((bet, idx) => (
            <Card key={idx} className="bg-zinc-900">
              <CardContent className="p-4 space-y-1">
                <div className="flex justify-between items-center">
                  <div className="flex items-center space-x-2">
                    <div className="bg-zinc-800 px-2 py-1 rounded-md text-sm">
                      {bet.line}
                    </div>
                    <span className="font-medium">{bet.destination}</span>
                  </div>
                  <div className="text-green-500 text-sm text-right">
                    <p>{bet.when}</p>
                    <p className="font-semibold">+ {bet.amount} PLN</p>
                  </div>
                </div>
                <div className="text-sm text-zinc-400 flex items-center space-x-4">
                  <span>{bet.stop}</span>
                  <span>⏱ {bet.time}</span>
                  <span>x {bet.odds}</span>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
        <Button variant="outline" className="w-full mt-4">
          Więcej statystyk
        </Button>
      </div>

      <div>
        <h3 className="font-semibold mb-1">Wirtualne przystanki</h3>
        <p className="text-zinc-400 text-sm mb-3">
          Wybierz przystanek, aby obstawić kolejne zdarzenia wrocławskiego MPK.
        </p>
        <div className="grid grid-cols-2 gap-3">
          {data?.virtualStops.map((stop, idx) => (
            <Card key={idx} className="bg-zinc-900 p-4 space-y-2">
              <div className="font-medium">{stop.name}</div>
              <div className="flex space-x-2">
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
          ))}
        </div>
      </div>
    </div>
  );
}
