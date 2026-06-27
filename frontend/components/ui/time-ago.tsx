"use client";

import { useState, useEffect } from "react";

interface TimeAgoProps {
  date: Date | string | number;
  className?: string;
}

export function TimeAgo({ date, className = "" }: TimeAgoProps) {
  const [timeStr, setTimeStr] = useState("");

  useEffect(() => {
    const parseDate = (d: Date | string | number) => {
      if (d instanceof Date) return d;
      return new Date(d);
    };

    const targetDate = parseDate(date);

    const updateTime = () => {
      const seconds = Math.floor((new Date().getTime() - targetDate.getTime()) / 1000);
      
      let str = "";
      if (seconds < 60) str = "just now";
      else if (seconds < 3600) str = `${Math.floor(seconds / 60)}m ago`;
      else if (seconds < 86400) str = `${Math.floor(seconds / 3600)}h ago`;
      else str = `${Math.floor(seconds / 86400)}d ago`;

      setTimeStr(str);
    };

    // Run immediately
    updateTime();

    // Re-run every minute (60000ms)
    // For extreme accuracy on "just now", you might run every 10-30s instead.
    const intervalId = setInterval(updateTime, 60000);

    return () => clearInterval(intervalId);
  }, [date]);

  // Initial SSR pass will render nothing or a fallback to avoid hydration mismatch,
  // or return a standard text.
  const [mounted, setMounted] = useState(false);
  useEffect(() => setMounted(true), []);

  if (!mounted) return <span className={className}>...</span>;

  return <span className={className}>{timeStr}</span>;
}