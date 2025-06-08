import { Card, CardContent } from "@/components/ui/card";

import { Button } from "@/components/ui/button";
import { Github } from "lucide-react";
import { Input } from "@/components/ui/input";
import { useNavigate } from "@tanstack/react-router";

export default function SignupPage() {
  const navigate = useNavigate();

  const handleLogin = (e: React.MouseEvent<HTMLButtonElement>): void => {
    e.preventDefault();
    navigate({
      to: "/login",
    });
  };

  return (
    <div className="min-h-screen bg-black flex flex-col items-center justify-center text-white px-4">
      <div className="w-full max-w-md">
        <div className="flex justify-between mb-6">
          <span className="text-gray-400 font-semibold">Auth</span>
          <button onClick={handleLogin} className="text-white">
            Login
          </button>
        </div>

        <Card className="bg-black border-none shadow-none">
          <CardContent className="p-0 space-y-6">
            <div className="text-center">
              <h1 className="text-2xl font-bold">Create an account</h1>
              <p className="text-gray-400">
                Enter your email below to create your account
              </p>
            </div>

            <div className="space-y-4">
              <Input
                placeholder="name@example.com"
                className="bg-zinc-900 text-white"
              />
              <Button className="w-full bg-white text-black hover:bg-zinc-100">
                Sign In with Email
              </Button>
            </div>

            <div className="flex items-center justify-center space-x-2">
              <hr className="flex-grow border-zinc-700" />
              <span className="text-xs text-zinc-500">OR CONTINUE WITH</span>
              <hr className="flex-grow border-zinc-700" />
            </div>

            <Button
              variant="outline"
              className="w-full border-zinc-700 bg-zinc-900 text-white hover:bg-zinc-800"
            >
              <Github className="mr-2 h-4 w-4" /> Github
            </Button>

            <p className="text-center text-xs text-gray-400">
              By clicking continue, you agree to our{" "}
              <a href="#" className="underline">
                Terms of Service
              </a>{" "}
              and{" "}
              <a href="#" className="underline">
                Privacy Policy
              </a>
              .
            </p>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
