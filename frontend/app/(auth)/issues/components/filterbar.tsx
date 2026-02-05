import {
  Select,
  SelectTrigger,
  SelectContent,
  SelectItem,
  SelectValue,
} from "@/components/ui/select";
import { Switch } from "@/components/ui/switch";

export default function FilterBar() {
  return (
    <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between bg-[#161b22] border border-[#30363d] rounded-xl p-4">
      
      {/* Left Filters */}
      <div className="flex flex-wrap gap-3">
        <Select>
          <SelectTrigger className="w-[160px] bg-[#0d1117] border-[#30363d]">
            <SelectValue placeholder="Language" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="javascript">JavaScript</SelectItem>
            <SelectItem value="python">Python</SelectItem>
            <SelectItem value="java">Java</SelectItem>
            <SelectItem value="go">Go</SelectItem>
          </SelectContent>
        </Select>

        <Select>
          <SelectTrigger className="w-[160px] bg-[#0d1117] border-[#30363d]">
            <SelectValue placeholder="Difficulty" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="easy">Easy</SelectItem>
            <SelectItem value="medium">Medium</SelectItem>
          </SelectContent>
        </Select>
      </div>

      {/* Right Toggles */}
      <div className="flex items-center gap-6">
        <div className="flex items-center gap-2">
          <Switch />
          <span className="text-sm text-[#8b949e]">
            Beginner-friendly
          </span>
        </div>

        <div className="flex items-center gap-2">
          <Switch />
          <span className="text-sm text-[#8b949e]">
            Recently active
          </span>
        </div>
      </div>
    </div>
  );
}
