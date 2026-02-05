"use client";

import { useEffect } from "react";
import { useRouter, useSearchParams } from "next/navigation";

export default function AuthSuccessPage() {
  const router = useRouter();
  const searchParams = useSearchParams();

  useEffect(() => {
    const token = searchParams.get("token");

    if (token) {
      localStorage.setItem("authToken", token);
      // Redirect to onboarding or wherever appropriate
      router.push("/onboarding"); 
    } else {
      // Handle error case - maybe redirect to login with error
       router.push("/?error=auth_failed");
    }
  }, [router, searchParams]);

  return (
    <div className="flex h-screen items-center justify-center bg-[#0d1117] text-white">
      <div className="text-center space-y-4">
        <h1 className="text-2xl font-semibold">Authenticating...</h1>
        <p className="text-[#8b949e]">Please wait while we log you in.</p>
        {/* Optional: Add a spinner here */}
      </div>
    </div>
  );
}
