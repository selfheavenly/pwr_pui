import { Card, CardContent } from "@/components/ui/card";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";

export default function LoginPage() {
  return (
    <div className="min-h-screen bg-black flex flex-col items-center justify-center text-white px-4">
      <div className="w-full max-w-md">
        <Card className="bg-black border-none shadow-none">
          <CardContent className="p-0 space-y-6">
            <div className="text-center">
              <h1 className="text-2xl font-bold">Log in</h1>
              <p className="text-gray-400">Welcome back. Please log in.</p>
            </div>

            <div className="space-y-4">
              <Input placeholder="Email" className="bg-zinc-900 text-white" />
              <Input
                placeholder="Password"
                type="password"
                className="bg-zinc-900 text-white"
              />
              <Button className="w-full bg-white text-black hover:bg-zinc-100">
                Log In
              </Button>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
