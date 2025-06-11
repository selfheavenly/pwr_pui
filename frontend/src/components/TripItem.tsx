import {
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "@/components/ui/accordion";
import type { Dispatch, SetStateAction } from "react";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import type { StopDeparture } from "@/types/stops";
import { useTripOdds } from "@/hooks/useVStopData";

interface Props {
  dep: StopDeparture;
  idx: number;
  selectedOdds: number | null;
  setSelectedOdds: Dispatch<SetStateAction<number | null>>;
  betAmount: string;
  setBetAmount: Dispatch<SetStateAction<string>>;
  handleBet: (balance: number) => void;
}

export default function TripItem({
  dep,
  idx,
  selectedOdds,
  setSelectedOdds,
  betAmount,
  setBetAmount,
  handleBet,
}: Props) {
  const { data: oddsData, isLoading: oddsLoading } = useTripOdds(dep.tram_id);

  const numericValue = parseFloat(betAmount);
  const isEmpty = !betAmount || isNaN(numericValue) || numericValue <= 0;
  const isTooMuch = oddsData && numericValue > oddsData.balance;

  return (
    <AccordionItem value={`item-${idx}`} className="border-none rounded-xl">
      <AccordionTrigger className="rounded-lg bg-zinc-900 px-4 py-3 text-left flex justify-between items-center text-sm hover:no-underline">
        <div className="flex items-center gap-3">
          <span className="bg-zinc-800 px-2 py-1 rounded text-xs font-medium">
            {dep.line.padStart(2, "0")}
          </span>
          <span className="font-medium text-base">{dep.destination}</span>
        </div>
        <span className="text-xs text-zinc-400">⏱ {dep.arrival_time}</span>
      </AccordionTrigger>

      <AccordionContent className="bg-zinc-950 px-4 py-4 space-y-4 rounded-b-xl">
        {oddsLoading ? (
          <div className="text-sm text-gray-400 text-center">
            Ładowanie kursów...
          </div>
        ) : oddsData && oddsData.odds.length > 0 ? (
          <>
            {/* Kursy */}
            <div className="grid grid-cols-2 sm:grid-cols-4 gap-2">
              {oddsData.odds.map((odd, i) => (
                <Button
                  key={i}
                  variant="outline"
                  onClick={() => setSelectedOdds(odd.value)}
                  className={`h-16 flex flex-col items-center p-3 text-white font-semibold ${
                    selectedOdds === odd.value
                      ? "border-2 border-white"
                      : "border border-zinc-700"
                  }`}
                >
                  <span className="text-xs font-normal">{odd.label}</span>
                  <span className="text-lg">{odd.value.toFixed(2)}</span>
                </Button>
              ))}
            </div>

            {/* Saldo */}
            <div className="text-center text-sm text-gray-400">
              Dostępne:{" "}
              {oddsData.balance.toLocaleString("pl-PL", {
                style: "currency",
                currency: "PLN",
              })}
            </div>

            {/* Zakład */}
            <div className="flex flex-col sm:flex-row gap-2">
              <Input
                type="number"
                step="any"
                min="0"
                placeholder="Kwota zakładu"
                className="bg-zinc-800 text-white w-full"
                value={betAmount}
                onChange={(e) => setBetAmount(e.target.value)}
              />
              <Button
                onClick={() => handleBet(oddsData.balance)}
                disabled={isEmpty || isTooMuch}
                className="bg-white text-black px-6 disabled:opacity-50"
              >
                Obstaw
              </Button>
            </div>

            {/* Błąd */}
            {!isEmpty && isTooMuch && (
              <p className="text-center text-sm text-red-500">
                Podana kwota przekracza dostępne saldo konta.
              </p>
            )}
          </>
        ) : (
          <p className="text-center text-gray-500 text-sm">
            Brak dostępnych kursów.
          </p>
        )}
      </AccordionContent>
    </AccordionItem>
  );
}
