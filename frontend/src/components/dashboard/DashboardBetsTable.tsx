import { FaCheckCircle, FaCircle, FaClock } from "react-icons/fa";
import {
  Pagination,
  PaginationContent,
  PaginationItem,
  PaginationNext,
  PaginationPrevious,
} from "@/components/ui/pagination";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";

import { Badge } from "@/components/ui/badge";
import { type BetBrief } from "@/types/bets";
import { cn } from "@/lib/utils";
import { useRecentBets } from "@/hooks/useRecentBets";
import { useState, type JSX } from "react";

// Relative date formatter
const formatRelativeDate = (dateString: string) => {
  const now = new Date();
  const placed = new Date(dateString);
  const diff = Math.floor((now.getTime() - placed.getTime()) / 1000);

  const ranges: [number, Intl.RelativeTimeFormatUnit][] = [
    [60, "seconds"],
    [3600, "minutes"],
    [86400, "hours"],
    [604800, "days"],
    [2592000, "weeks"],
    [31536000, "months"],
    [Infinity, "years"],
  ];

  for (const [range, unit] of ranges) {
    if (diff < range) {
      const value = Math.floor(diff / (range / (unit === "seconds" ? 1 : 60)));
      return new Intl.RelativeTimeFormat("pl", { numeric: "auto" }).format(
        -value,
        unit
      );
    }
  }
};

const statusMap: Record<
  BetBrief["status"],
  {
    label: string;
    icon: JSX.Element;
    className: string;
  }
> = {
  won: {
    label: "Wygrany",
    icon: <FaCheckCircle className="text-green-500 mr-1" />,
    className: "text-green-500 border border-green-700/30 bg-green-900/10",
  },
  lost: {
    label: "Przegrany",
    icon: <FaCircle className="text-red-500 mr-1" />,
    className: "text-red-500 border border-red-700/30 bg-red-900/10",
  },
  pending: {
    label: "W toku",
    icon: <FaClock className="text-muted-foreground mr-1" />,
    className: "text-muted-foreground border border-zinc-600/40 bg-zinc-800/30",
  },
};

export function BetsTable() {
  const [page, setPage] = useState(0);
  const pageSize = 5;

  const { data, isLoading, isError } = useRecentBets(page, pageSize);
  const bets = data?.data ?? [];
  const hasNext = data?.next_page ?? false;
  const totalPages = data?.page_count ?? 0;
  const currentPage = page + 1;

  return (
    <div className="w-full">
      <h3 className="font-semibold text-lg mb-4">Wszystkie zakłady</h3>

      {isLoading && <p className="text-sm text-muted">Ładowanie...</p>}
      {isError && (
        <p className="text-sm text-red-500">Błąd podczas ładowania zakładów.</p>
      )}

      {!isLoading && !isError && (
        <>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Data</TableHead>
                <TableHead>Linia</TableHead>
                <TableHead>Przystanek</TableHead>
                <TableHead>Kierunek</TableHead>
                <TableHead>Opóźnienie</TableHead>
                <TableHead>Kurs</TableHead>
                <TableHead>Stawka</TableHead>
                <TableHead>Wynik</TableHead>
                <TableHead>Status</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {bets.map((bet) => {
                const status = statusMap[bet.status];
                return (
                  <TableRow key={bet.bet_id}>
                    <TableCell className="text-sm text-muted-foreground">
                      {formatRelativeDate(bet.placed_at)}
                    </TableCell>
                    <TableCell className="text-center">
                      <div className="inline-flex items-center bg-zinc-800 px-2 py-1 rounded-md text-xs">
                        {bet.tram_lane_id}
                      </div>
                    </TableCell>
                    <TableCell>{bet.stop_name}</TableCell>
                    <TableCell>{bet.tram_lane_destination}</TableCell>
                    <TableCell>{bet.actual_delay}</TableCell>
                    <TableCell>{bet.bet_rate.toFixed(2)}</TableCell>
                    <TableCell>{bet.bet_amount} zł</TableCell>
                    <TableCell
                      className={cn(
                        "font-medium",
                        bet.bet_result > 0 ? "text-green-400" : "text-red-400"
                      )}
                    >
                      {bet.bet_result.toFixed(2)} zł
                    </TableCell>
                    <TableCell>
                      <span
                        className={cn(
                          "inline-flex items-center rounded-full px-2 py-0.5 text-xs font-medium",
                          status.className
                        )}
                      >
                        {status.icon}
                        {status.label}
                      </span>
                    </TableCell>
                  </TableRow>
                );
              })}
            </TableBody>
          </Table>

          <Pagination className="mt-6 justify-center">
            <PaginationContent>
              <PaginationItem>
                <PaginationPrevious
                  href="#"
                  onClick={(e) => {
                    e.preventDefault();
                    if (page > 0) setPage((p) => p - 1);
                  }}
                  className={cn({
                    "pointer-events-none opacity-50": page === 0,
                  })}
                />
              </PaginationItem>

              <PaginationItem>
                <span className="text-sm text-zinc-300">
                  Strona {currentPage} z {totalPages}
                </span>
              </PaginationItem>

              <PaginationItem>
                <PaginationNext
                  href="#"
                  onClick={(e) => {
                    e.preventDefault();
                    if (hasNext) setPage((p) => p + 1);
                  }}
                  className={cn({ "pointer-events-none opacity-50": !hasNext })}
                />
              </PaginationItem>
            </PaginationContent>
          </Pagination>
        </>
      )}
    </div>
  );
}
