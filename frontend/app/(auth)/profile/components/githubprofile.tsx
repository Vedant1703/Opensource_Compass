
import { Github } from "lucide-react";
import { Button } from "@/components/ui/button";

export default function GitHubProfile() {
  return (
    <div className="bg-[#0d1117] border border-[#30363d] rounded-xl p-6">
      <div className="flex items-center gap-4">
        <div className="h-12 w-12 rounded-full bg-[#161b22] flex items-center justify-center">
          <Github className="h-6 w-6" />
        </div>

        <div className="flex-1">
          <p className="font-medium">
            Connected as vedantkulkarni
          </p>
          <p className="text-sm text-[#8b949e]">
            Read-only access to public repositories
          </p>
        </div>

        <Button variant="outline" className="border-[#30363d]">
          Reconnect
        </Button>
      </div>
    </div>
  );
}
