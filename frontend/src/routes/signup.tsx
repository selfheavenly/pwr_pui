import SignupPage from "@/components/pages/SignupPage";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/signup")({
  component: SignupPage,
});
