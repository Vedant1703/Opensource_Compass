"use client";

import PageWrapper from "@/components/ui/page-wrapper";
import EmptyIssuesState from "./components/EmptyIssuesState";

export default function IssuesPage() {
  // Assuming for now that we have no issues to display (simulating the empty state)
  const hasIssues = false; 

  return (
    <PageWrapper className="relative min-h-screen">
      <div className="relative z-10 space-y-8 pb-20">
        <div className="flex flex-col md:flex-row md:items-end justify-between gap-4 border-b border-[#30363d]/50 pb-8">
            <div>
                <h1 className="text-3xl md:text-4xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-blue-400 via-purple-400 to-pink-400 inline-block">
                    Your Issues Feed
                </h1>
                <p className="mt-2 text-[#8b949e] text-lg max-w-2xl">
                    Track the issues from the repositories you watch
                </p>
            </div>
        </div>

        {/* Render our newly created Empty State when hasIssues is false */}
        {!hasIssues ? (
           <EmptyIssuesState />
        ) : (
           <div>{/* Future populated state would go here */}</div>
        )}
      </div>
    </PageWrapper>
  );
}
