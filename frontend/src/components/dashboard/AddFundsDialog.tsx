import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { cn } from "@/lib/utils";
import { useState } from "react";

interface Props {
  open: boolean;
  setOpen: (open: boolean) => void;
}

export function AddFundsDialog({ open, setOpen }: Props) {
  const [amount, setAmount] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState("");

  const handleSubmit = () => {
    const value = parseFloat(amount.replace(",", "."));
    if (isNaN(value) || value < 1) {
      setError("Minimalna kwota to 1 zł");
      return;
    }

    setIsLoading(true);
    setTimeout(() => {
      console.log(`Dodano ${value} PLN`);
      setIsLoading(false);
      setOpen(false);
      setAmount("");
    }, 1000);
  };

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogContent className="max-w-md">
        <DialogHeader>
          <DialogTitle>Dodaj środki</DialogTitle>
        </DialogHeader>

        <div className="space-y-4 mt-2">
          <div>
            <Label htmlFor="amount">Kwota (PLN)</Label>
            <Input
              id="amount"
              type="number"
              min="1"
              step="0.01"
              value={amount}
              onChange={(e) => setAmount(e.target.value)}
              className={cn(error && "border-red-500")}
            />
            {error && <p className="text-sm text-red-500 mt-1">{error}</p>}
          </div>

          <Button
            onClick={handleSubmit}
            disabled={isLoading}
            className="w-full"
          >
            {isLoading ? "Przetwarzanie..." : "Dodaj środki"}
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  );
}
