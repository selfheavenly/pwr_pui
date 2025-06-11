import { Button } from "@/components/ui/button";
import { FcGoogle } from "react-icons/fc";
import { useNavigate } from "@tanstack/react-router";

export default function LoginPage() {
  const navigate = useNavigate();

  const handleGoogleLogin = () => {
    // Tu normalnie byłoby logowanie przez Google
    console.log("Zaloguj przez Google");

    // Tymczasowe przekierowanie
    navigate({ to: "/dashboard" });
  };

  return (
    <div className="fixed inset-0 z-0">
      {/* 🔁 GIF jako tło */}
      <img
        src="/bgtram.gif"
        alt="Background"
        className="absolute inset-0 w-full h-full object-cover"
      />

      {/* 🔲 Overlay */}
      <div className="absolute inset-0 bg-black/60 backdrop-blur-md" />

      {/* 📦 Karta logowania */}
      <div className="relative z-10 flex items-center justify-center h-full">
        <div className="bg-black/40 backdrop-blur-lg border border-zinc-700 rounded-2xl shadow-xl px-8 py-10 max-w-sm w-full text-center space-y-6">
          <h1 className="text-3xl font-bold text-white">MPKBet</h1>
          <p className="text-gray-300 text-sm">
            Zaloguj się przez Google, aby rozpocząć
          </p>

          <Button
            onClick={handleGoogleLogin}
            className="w-full bg-white text-black hover:bg-zinc-100 flex items-center justify-center gap-2 px-6 py-5 rounded-xl font-medium text-base"
          >
            <FcGoogle className="text-xl" />
            Zaloguj się przez Google
          </Button>
        </div>
      </div>
    </div>
  );
}
