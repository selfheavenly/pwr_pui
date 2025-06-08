import { Card, CardContent } from "@/components/ui/card";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { useState } from "react";
import { useVStopData } from "@/hooks/useVStopData";

export default function VStopPage() {
  const { data, isLoading } = useVStopData();
  const [selectedOdds, setSelectedOdds] = useState<number | null>(null);
  const [betAmount, setBetAmount] = useState("");

  if (isLoading || !data) return <div className="text-white">Loading...</div>;

  const handleBet = () => {
    const amount = parseFloat(betAmount);
    if (isNaN(amount) || amount > data.balance) return;
    // logic to place bet here
    alert(`Betting ${amount} on odds ${selectedOdds}`);
  };

  return (
    <div className="min-h-screen bg-black text-white p-4 space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-xl font-semibold">
          Home / <span className="text-white font-bold">{data.stopName}</span>
        </h1>
        <div className="flex space-x-1">
          {data.lines.map((line, idx) => (
            <span
              key={idx}
              className="bg-zinc-800 text-white px-2 py-1 rounded text-xs"
            >
              {line}
            </span>
          ))}
        </div>
      </div>

      <img
        src={data.imageUrl}
        alt="tram"
        className="rounded-xl w-full object-cover h-48"
      />

      <div>
        <h2 className="text-sm text-gray-400 mb-2">Tablica odjazdów</h2>
        {data.departures.map((dep, idx) => (
          <Card key={idx} className="bg-zinc-900 mb-4">
            <CardContent className="p-4">
              <div className="flex justify-between items-center">
                <div className="flex items-center space-x-2">
                  <span className="bg-zinc-800 px-2 py-1 rounded text-sm">
                    {dep.line}
                  </span>
                  <span>{dep.destination}</span>
                </div>
                <span className="text-sm text-zinc-400">
                  ⏱ {dep.arrivalTime}
                </span>
              </div>

              <div className="flex justify-between mt-4 gap-2">
                {dep.odds.map((odd, i) => (
                  <Button
                    key={i}
                    variant={selectedOdds === odd.value ? "default" : "outline"}
                    onClick={() => setSelectedOdds(odd.value)}
                    className="flex-1 flex flex-col items-center p-2"
                  >
                    <span className="text-xs">{odd.label}</span>
                    <span className="text-lg font-semibold">
                      {odd.value.toFixed(2)}
                    </span>
                  </Button>
                ))}
              </div>

              {idx === 1 && (
                <div className="mt-4 space-y-2">
                  <p className="text-center text-sm text-gray-400">
                    Dostępne:{" "}
                    {data.balance.toLocaleString("pl-PL", {
                      style: "currency",
                      currency: "PLN",
                    })}
                  </p>
                  <div className="flex gap-2">
                    <Input
                      placeholder="Kwota zakładu"
                      className="bg-zinc-800 text-white"
                      value={betAmount}
                      onChange={(e) => setBetAmount(e.target.value)}
                    />
                    <Button onClick={handleBet} className="bg-white text-black">
                      Obstaw
                    </Button>
                  </div>
                  {parseFloat(betAmount) > data.balance && (
                    <p className="text-center text-sm text-red-500">
                      Podana kwota przekracza dostępne saldo konta.
                    </p>
                  )}
                </div>
              )}
            </CardContent>
          </Card>
        ))}
      </div>
    </div>
  );
}
