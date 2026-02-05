import { Button } from "@/components/ui/button";

export default function DangerZone() {
  return (
    <div className="border border-[#f85149] rounded-xl p-6 bg-[#0d1117]/50">
      <h2 className="text-lg font-semibold text-[#f85149]">
        Danger zone
      </h2>

      <p className="mt-2 text-sm text-[#8b949e]">
        Disconnecting GitHub will remove personalized
        recommendations and saved issues.
      </p>

      <Button
        variant="outline"
        className="mt-4 border-[#f85149] text-[#f85149] hover:bg-[#f85149]/10"
      >
        Disconnect GitHub
      </Button>
    </div>
  );
}
