"use client";

import { ExternalLink, Tag } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Issue } from "@/lib/api/core-service";

interface FeaturedIssuesProps {
  issues: Issue[];
  repoId: string;
}

export default function FeaturedIssues({ issues, repoId }: FeaturedIssuesProps) {
  if (!issues || issues.length === 0) {
    return null;
  }


  const displayIssues = issues.slice(0, 3);

  return (
    <div id="issues" className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold">Good First Issues</h2>
          <p className="text-[#8b949e] mt-1">Start contributing with these beginner-friendly issues</p>
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {displayIssues.map((issue, index) => {
          const targetUrl = issue.issue_url 
            ?? (issue.issue_number ? `https://github.com/${repoId}/issues/${issue.issue_number}` : `https://github.com/${repoId}/issues`);

          return (
          <div
            key={index}
            className="bg-[#0d1117]/60 backdrop-blur-sm border border-[#30363d] rounded-xl p-5 hover:border-[#2f81f7] hover:shadow-lg hover:shadow-[#2f81f7]/10 transition-all group"
          >
            <div className="space-y-3">
              <h3 className="font-medium leading-snug line-clamp-2 group-hover:text-[#2f81f7] transition-colors">
                {issue.title}
              </h3>

              {/* AI Summary if available */}
              {issue.ai?.summary && (
                <p className="text-sm text-[#8b949e] line-clamp-2">
                  {issue.ai.summary}
                </p>
              )}

              {/* Labels */}
              {issue.labels && issue.labels.length > 0 && (
                <div className="flex flex-wrap gap-2">
                  {issue.labels.slice(0, 2).map((label) => (
                    <span
                      key={label}
                      className="inline-flex items-center gap-1 px-2 py-1 text-xs rounded-full bg-[#161b22] border border-[#30363d] text-[#8b949e]"
                    >
                      <Tag className="h-3 w-3" />
                      {label}
                    </span>
                  ))}
                  {issue.labels.length > 2 && (
                    <span className="px-2 py-1 text-xs rounded-full bg-[#161b22] border border-[#30363d] text-[#8b949e]">
                      +{issue.labels.length - 2}
                    </span>
                  )}
                </div>
              )}

              {/* AI Metadata */}
              {issue.ai && (issue.ai.difficulty || issue.ai.estimated_time) && (
                <div className="flex items-center gap-3 text-xs text-[#8b949e]">
                  {issue.ai.difficulty && (
                    <span className="px-2 py-0.5 rounded bg-[#161b22]">
                      {issue.ai.difficulty}
                    </span>
                  )}
                  {issue.ai.estimated_time && (
                    <span className="px-2 py-0.5 rounded bg-[#161b22]">
                      ⏱️ {issue.ai.estimated_time}
                    </span>
                  )}
                </div>
              )}
            </div>

            <div className="mt-4">
              <Button
                variant="secondary"
                className="w-full bg-[#21262d] hover:bg-[#30363d] text-[#c9d1d9] border border-[#30363d]"
                onClick={() => window.open(targetUrl, '_blank')}
              >
                <ExternalLink className="h-4 w-4 mr-2" />
                  {issue.issue_number ? `View Issue #${issue.issue_number}` : 'View Issues'}
              </Button>
            </div>
          </div>
        )})}
      </div>
    </div>
  );
}