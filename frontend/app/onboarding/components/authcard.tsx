import { Github } from "lucide-react";
import { Button } from "@/components/ui/button";

export default function AuthCard() {
  return (
    <div className="bg-[#161b22] border border-[#30363d] rounded-xl p-6 space-y-4">
      <div className="flex items-center gap-3">
        <Github className="h-6 w-6 text-white" />
        <h1 className="text-xl font-semibold">
          Connect your GitHub account
        </h1>
      </div>

      <p className="text-sm text-[#8b949e]">
        We use GitHub data to recommend repositories and beginner-friendly
        issues that match your interests.
      </p>

      <Button
        className="bg-[#238636] hover:bg-[#2ea043] w-full"
        onClick={() => {
          window.location.href = `${process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'}/auth/github`;
        }}
      >
        Continue with GitHub
      </Button>

      <p className="text-xs text-[#6e7681]">
        Read-only access • No private repositories
      </p>
    </div>
  );
}
