"use client";

import { Github } from "lucide-react";
import { Button } from "@/components/ui/button";
import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";

export default function GitHubProfile() {
  const router = useRouter();
  const [user, setUser] = useState<{ username: string; avatar: string } | null>(null);

  useEffect(() => {
    const username = localStorage.getItem("githubUsername");
    const avatar = localStorage.getItem("githubAvatar");
    if (username) {
      setUser({ username, avatar: avatar || "" });
    }
  }, []);

  const handleDisconnect = () => {
    // Clear all localStorage items
    localStorage.removeItem("authToken");
    localStorage.removeItem("githubUsername");
    localStorage.removeItem("githubAvatar");
    localStorage.removeItem("techStack");
    localStorage.removeItem("frameworks");
    localStorage.removeItem("domains");
    
    // Redirect to home page
    router.push("/");
  };

  if (!user) {
    return (
        <div className="bg-[#0d1117] border border-[#30363d] rounded-xl p-6">
            <p className="text-[#8b949e]">Not connected to GitHub</p>
        </div>
    )
  }

  return (
    <div className="bg-[#0d1117] border border-[#30363d] rounded-xl p-6">
      <div className="flex items-center gap-4">
        <div className="h-12 w-12 rounded-full bg-[#161b22] flex items-center justify-center overflow-hidden">
             {user.avatar ? (
                 <img src={user.avatar} alt="Avatar" className="h-full w-full object-cover" />
             ) : (
                <Github className="h-6 w-6" />
             )}
        </div>

        <div className="flex-1">
          <p className="font-medium">
            Connected as {user.username}
          </p>
          <p className="text-sm text-[#8b949e]">
            Read-only access to public repositories
          </p>
        </div>

        <Button 
          variant="outline" 
          className="border-[#30363d] hover:border-red-500 hover:text-red-500"
          onClick={handleDisconnect}
        >
          Disconnect
        </Button>
      </div>
    </div>
  );
}
