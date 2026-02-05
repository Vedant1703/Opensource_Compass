"use client";

import { useState } from "react";
import TechChip from "../../../onboarding/components/techchip";

const LANGUAGES = [
  "JavaScript",
  "TypeScript",
  "Python",
  "Java",
  "Go",
];

export default function TechPreferences() {
  const [selected, setSelected] = useState([
    "JavaScript",
    "TypeScript",
  ]);

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
            toggle={() =>
              setSelected(
                selected.includes(lang)
                  ? selected.filter((v) => v !== lang)
                  : [...selected, lang]
              )
            }
          />
        ))}
      </div>
    </div>
  );
}
