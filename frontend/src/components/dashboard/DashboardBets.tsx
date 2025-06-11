import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import {
  Pagination,
  PaginationContent,
  PaginationItem,
  PaginationNext,
  PaginationPrevious,
} from "@/components/ui/pagination";

import type { BetBrief } from "@/types/bets";
import { BetCard } from "@/components/bet/BetCard";
import { Button } from "@/components/ui/button";
import { useRecentBets } from "@/hooks/useRecentBets";
import { useState } from "react";

interface Props {
  lastBets: BetBrief[];
}

export function DashboardBets({ lastBets }: Props) {
  const [open, setOpen] = useState(false);
  const [page, setPage] = useState(0);
  const pageSize = 5;

  const { data, isLoading, isError } = useRecentBets(page, pageSize);
  const bets = data?.data ?? [];
  const hasNext = data?.next_page ?? false;

  const totalPages = data?.page_count ?? 0;
  const currentPage = page + 1;

  return (
    <div>
      <h3 className="font-semibold mb-2">Ostatnie zakłady</h3>
      <div className="space-y-3">
        {lastBets.slice(0, 2).map((bet) => (
          <BetCard key={bet.bet_id} bet={bet} />
        ))}
      </div>

      <Dialog
        open={open}
        onOpenChange={(v) => {
          setOpen(v);
          if (!v) setPage(0);
        }}
      >
        <DialogTrigger asChild>
          <Button variant="outline" className="w-full mt-4">
            Więcej statystyk
          </Button>
        </DialogTrigger>

        <DialogContent className="max-w-2xl backdrop-blur-md bg-black/50 text-white border border-white/10">
          <DialogHeader>
            <DialogTitle>Wszystkie zakłady</DialogTitle>
          </DialogHeader>

          {isLoading && <p className="text-sm text-muted">Ładowanie...</p>}
          {isError && (
            <p className="text-sm text-red-500">
              Błąd podczas ładowania zakładów.
            </p>
          )}

          <div className="space-y-3 mt-4">
            {bets.map((bet) => (
              <BetCard key={bet.bet_id} bet={bet} />
            ))}
          </div>

          <Pagination className="mt-6 justify-center">
            <PaginationContent>
              <PaginationItem>
                <PaginationPrevious
                  href="#"
                  onClick={(e) => {
                    e.preventDefault();
                    if (page > 0) setPage((p) => p - 1);
                  }}
                  className={page === 0 ? "pointer-events-none opacity-50" : ""}
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
                  className={!hasNext ? "pointer-events-none opacity-50" : ""}
                />
              </PaginationItem>
            </PaginationContent>
          </Pagination>
        </DialogContent>
      </Dialog>
    </div>
  );
}
