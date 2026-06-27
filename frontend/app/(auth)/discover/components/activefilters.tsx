"use client";

import { X, Settings } from "lucide-react";
import { useRouter } from "next/navigation";
import Link from 'next/link';

interface ActiveFiltersProps {
  languages: string[];
  topics: string[];
  experienceLevel: string;
}

export default function ActiveFilters({ languages, topics, experienceLevel }: ActiveFiltersProps) {
  const router = useRouter();

  if (languages.length === 0 && topics.length === 0) {
    return null;
  }

  return (
    <div className="bg-[#161b22] border border-[#30363d] rounded-xl p-4">
      <div className="flex items-start justify-between gap-4">
        <div className="flex-1">
          <div className="flex items-center gap-2 mb-3">
            <h3 className="text-sm font-semibold text-[#8b949e]">Active Filters</h3>
          </div>

          <div className="flex flex-wrap gap-2">
            {/* Experience Level */}
            <div className="inline-flex items-center gap-1.5 px-3 py-1.5 bg-[#1f6feb]/10 border border-[#1f6feb] rounded-full text-sm">
              <span className="text-xs">📊</span>
              <span>{experienceLevel}</span>
            </div>

            {/* Languages */}
            {languages.map((lang) => (
              <div
                key={lang}
                className="inline-flex items-center gap-1.5 px-3 py-1.5 bg-[#0d1117] border border-[#30363d] rounded-full text-sm"
              >
                <span className="text-xs">💻</span>
                <span>{lang}</span>
              </div>
            ))}

            {/* Topics */}
            {topics.map((topic) => (
              <div
                key={topic}
                className="inline-flex items-center gap-1.5 px-3 py-1.5 bg-[#0d1117] border border-[#30363d] rounded-full text-sm"
              >
                <span className="text-xs">🏷️</span>
                <span>{topic}</span>
              </div>
            ))}
          </div>
        </div>

        import Link from 'next/link';

        <Link
          href="/settings"
          className="flex items-center gap-2 px-3 py-1.5 text-sm text-[#8b949e] hover:text-white transition-colors"
        >
          <Settings className="h-4 w-4" />
          Edit
        </Link>
      </div>
    </div>
  );
}
