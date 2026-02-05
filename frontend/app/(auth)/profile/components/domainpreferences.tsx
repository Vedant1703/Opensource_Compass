"use client";

import { useState } from "react";
import TechChip from "../../../onboarding/components/techchip";

const DOMAINS = [
  "Backend Systems",
  "Frontend",
  "DevOps",
  "AI / ML",
  "Open Source Tooling",
];

export default function DomainPreferences() {
  const [selected, setSelected] = useState([
    "Backend Systems",
    "Open Source Tooling",
  ]);

  return (
    <div className="space-y-4">
      <h2 className="text-xl font-semibold">
        Domains & interests
      </h2>

      <div className="flex flex-wrap gap-2">
        {DOMAINS.map((domain) => (
          <TechChip
            key={domain}
            label={domain}
            selected={selected.includes(domain)}
            toggle={() =>
              setSelected(
                selected.includes(domain)
                  ? selected.filter((v) => v !== domain)
                  : [...selected, domain]
              )
            }
          />
        ))}
      </div>
    </div>
  );
}
