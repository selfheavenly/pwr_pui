import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";

import { AddFundsDialog } from "./AddFundsDialog";
import { Button } from "@/components/ui/button";
import { useState } from "react";

export function DashboardHeader() {
  const [fundsOpen, setFundsOpen] = useState(false);

  const handleLogout = () => {
    // implement your logout logic here
    console.log("Logging out...");
  };

  return (
    <div className="flex justify-between items-center">
      <h1 className="text-xl font-semibold">MPKBet</h1>

      <AddFundsDialog open={fundsOpen} setOpen={setFundsOpen} />

      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Avatar className="h-10 w-10 cursor-pointer border border-white/10 hover:opacity-90 transition">
            <AvatarImage src="/user.png" alt="User" />
            <AvatarFallback>👤</AvatarFallback>
          </Avatar>
        </DropdownMenuTrigger>

        <DropdownMenuContent align="end" className="w-40">
          <DropdownMenuItem onClick={() => setFundsOpen(true)}>
            ➕ Dodaj środki
          </DropdownMenuItem>
          <DropdownMenuItem onClick={handleLogout}>🔓 Wyloguj</DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
    </div>
  );
}
