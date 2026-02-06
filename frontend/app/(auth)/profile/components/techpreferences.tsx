"use client";

import { useState, useEffect } from "react";
import TechChip from "../../../onboarding/components/techchip";

const LANGUAGES = [
  "JavaScript",
  "TypeScript",
  "Python",
  "Java",
  "Go",
];

export default function TechPreferences() {
  const [selected, setSelected] = useState<string[]>([]);

  useEffect(() => {
    const saved = localStorage.getItem("techStack");
    if (saved) {
      try {
        setSelected(JSON.parse(saved));
      } catch (e) {
        console.error("Failed to parse tech stack", e);
      }
    }
  }, []);

  const handleToggle = (lang: string) => {
    const newSelected = selected.includes(lang)
      ? selected.filter((v) => v !== lang)
      : [...selected, lang];
    
    setSelected(newSelected);
    localStorage.setItem("techStack", JSON.stringify(newSelected));
  };

  return (
    <div className="space-y-4">
      <h2 className="text-xl font-semibold">
        Programming languages
      </h2>

      <div className="flex flex-wrap gap-2">
        {LANGUAGES.map((lang) => (
          <TechChip
            key={lang}
            label={lang}
            selected={selected.includes(lang)}
            toggle={() => handleToggle(lang)}
          />
        ))}
      </div>
    </div>
  );
}
