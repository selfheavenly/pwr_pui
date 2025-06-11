import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbSeparator,
} from "@/components/ui/breadcrumb";

import { Accordion } from "@/components/ui/accordion";
import { DashboardHeader } from "../dashboard/DashboardHeader";
import TripItem from "@/components/TripItem";
import { useParams } from "@tanstack/react-router";
import { useState } from "react";
import { useVStopData } from "@/hooks/useVStopData";

export default function VStopPage() {
  const { id } = useParams({ from: "/vstop/$id" });
  const { data, isLoading } = useVStopData(id);
  const [selectedOdds, setSelectedOdds] = useState<number | null>(null);
  const [betAmount, setBetAmount] = useState("");

  const handleBet = (balance: number) => {
    const amount = parseFloat(betAmount);
    if (isNaN(amount) || amount > balance || !selectedOdds) return;
    alert(`Obstawiasz ${amount} zł na kurs ${selectedOdds}`);
    setBetAmount("");
    setSelectedOdds(null);
  };

  if (isLoading || !data) {
    return <div className="text-white p-4">Ładowanie danych przystanku...</div>;
  }

  return (
    <div className="min-h-screen bg-black text-white p-4 space-y-6">
      <DashboardHeader />

      {/* 🧭 Breadcrumb */}
      <Breadcrumb>
        <BreadcrumbList>
          <BreadcrumbItem>
            <BreadcrumbLink href="/dashboard">MPKBet</BreadcrumbLink>
          </BreadcrumbItem>
          <BreadcrumbSeparator />
          <BreadcrumbItem>
            <BreadcrumbLink className="text-white font-semibold" href="#">
              {data.stop_name}
            </BreadcrumbLink>
          </BreadcrumbItem>
        </BreadcrumbList>
      </Breadcrumb>

      {/* 🖼️ Obraz przystanku z numerami linii */}
      <div className="relative">
        <img
          src={data.image_url}
          alt={data.stop_name}
          className="rounded-xl h-48 w-full object-cover"
        />
        <div className="absolute top-4 right-4 flex gap-2">
          {data.lines.map((line, idx) => (
            <div
              key={idx}
              className="bg-zinc-800 text-white px-2 py-1 rounded text-xs"
            >
              {line.padStart(2, "0")}
            </div>
          ))}
        </div>
      </div>

      {/* 🚌 Lista odjazdów */}
      <div>
        <h2 className="text-sm text-gray-400 mb-2">Tablica odjazdów</h2>
        <Accordion type="multiple" className="space-y-2">
          {data.departures.map((dep, idx) => (
            <TripItem
              key={dep.tram_id}
              dep={dep}
              idx={idx}
              selectedOdds={selectedOdds}
              setSelectedOdds={setSelectedOdds}
              betAmount={betAmount}
              setBetAmount={setBetAmount}
              handleBet={handleBet}
            />
          ))}
        </Accordion>
      </div>
    </div>
  );
}
