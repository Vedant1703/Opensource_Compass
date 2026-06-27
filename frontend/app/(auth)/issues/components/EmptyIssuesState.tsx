import { Bug } from 'lucide-react';
import Link from 'next/link';

export default function EmptyIssuesState() {
  return (
    <div className="flex flex-col items-center justify-center py-24 text-center">
      <div className="w-16 h-16 rounded-full bg-[#161b22] border border-[#30363d] flex items-center justify-center mb-6">
        <Bug className="w-8 h-8 text-[#8b949e]" />
      </div>
      <h2 className="text-xl font-semibold text-[#c9d1d9] mb-2">No issues found</h2>
      <p className="text-sm text-[#8b949e] max-w-sm mb-6">
        Watch some repositories from the Discover page to start seeing their open issues here.
      </p>
      <Link
        href="/discover"
        className="inline-flex items-center gap-2 px-4 py-2 rounded-lg bg-[#238636] text-white text-sm hover:bg-[#2ea043] transition-colors"
      >
        Go to Discover
      </Link>
    </div>
  );
}