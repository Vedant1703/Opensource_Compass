"use client";

import { useState, useEffect } from "react";

interface SearchInputProps {
  onSearch: (query: string) => void;
  className?: string;
}

export default function SearchInput({ onSearch, className = "" }: SearchInputProps) {
  const [query, setQuery] = useState("");

  useEffect(() => {
    const timer = setTimeout(() => {
      onSearch(query);
    }, 1000); // 1000ms debounce (increased from 500ms to reduce API calls)

    return () => clearTimeout(timer);
  }, [query, onSearch]);

  return (
    <div className={`relative ${className}`}>
      <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
        <svg
          className="h-5 w-5 text-[#6e7681]"
          xmlns="http://www.w3.org/2000/svg"
          viewBox="0 0 20 20"
          fill="currentColor"
          aria-hidden="true"
        >
          <path
            fillRule="evenodd"
            d="M8 4a4 4 0 100 8 4 4 0 000-8zM2 8a6 6 0 1110.89 3.476l4.817 4.817a1 1 0 01-1.414 1.414l-4.816-4.816A6 6 0 012 8z"
            clipRule="evenodd"
          />
        </svg>
      </div>
      <input
        type="text"
        name="search"
        id="search"
        className="block w-full pl-10 pr-3 py-3 rounded-xl text-sm bg-[#161b22] border border-[#30363d] text-[#c9d1d9] placeholder-[#6e7681] focus:outline-none focus:border-[#2f81f7] focus:ring-1 focus:ring-[#2f81f7] transition-colors"
        placeholder="Search repositories (e.g. react, tensorflow)..."
        value={query}
        onChange={(e) => setQuery(e.target.value)}
      />
    </div>
  );
}
