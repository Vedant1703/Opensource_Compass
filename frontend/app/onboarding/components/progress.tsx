import { Check } from "lucide-react";

export default function ProgressStepper({
  currentStep,
}: {
  currentStep: number;
}) {
  const steps = ["Connect GitHub", "Choose Tech Stack", "Discover Repos"];

  return (
    <div className="flex items-center justify-between w-full">
      {steps.map((step, idx) => {
        const stepNum = idx + 1;
        const isCompleted = stepNum < currentStep;
        const isCurrent = stepNum === currentStep;

        return (
          <div key={step} className="flex-1 flex items-center last:flex-none">
            <div className="flex items-center gap-3">
              <div
                className={`h-8 w-8 rounded-full flex items-center justify-center text-sm font-medium transition-colors duration-300 ${
                  isCompleted
                    ? "bg-[#238636] text-white"
                    : isCurrent
                    ? "bg-[#2f81f7] text-white ring-2 ring-[#2f81f7] ring-offset-2 ring-offset-[#0d1117]"
                    : "bg-[#161b22] text-[#6e7681] border border-[#30363d]"
                }`}
              >
                {isCompleted ? <Check className="h-4 w-4" /> : stepNum}
              </div>
              <span
                className={`text-sm font-medium hidden sm:block ${
                  isCompleted || isCurrent ? "text-white" : "text-[#6e7681]"
                }`}
              >
                {step}
              </span>
            </div>
            {idx !== steps.length - 1 && (
  <div
    className={`flex-1 mx-4 h-[1px] hidden sm:block transition-all duration-1000 ease-in-out ${
      idx < currentStep 
        ? 'bg-[#238636]' 
        : idx === currentStep 
        ? 'bg-gradient-to-r from-[#238636] via-[#238636] to-[#30363d]' 
        : 'bg-[#30363d]' 
    }`}
  />
)}
          </div>
        );
      })}
    </div>
  );
}
