import VStopPage from "@/components/pages/VStopPage";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/vstop")({
  component: VStopPage,
});
