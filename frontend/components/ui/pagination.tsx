interface PaginationControlsProps {
  currentPage: number;
  totalPages: number;
  onPrev: () => void;
  onNext: () => void;
  onGoTo: (page: number) => void;
  hasPrev: boolean;
  hasNext: boolean;
}

export function PaginationControls({ currentPage, totalPages, onPrev, onNext, onGoTo, hasPrev, hasNext }: PaginationControlsProps) {
  if (totalPages <= 1) return null;

  return (
    <div className="flex items-center justify-center gap-2 mt-8">
      <button
        onClick={onPrev}
        disabled={!hasPrev}
        className="px-4 py-2 rounded-lg bg-[#161b22] border border-[#30363d] text-[#c9d1d9] disabled:opacity-40 hover:border-[#2f81f7] transition-colors text-sm"
      >
        ← Prev
      </button>

      {/* Page number buttons */}
      {Array.from({ length: totalPages }, (_, i) => i + 1).map((page) => (
        <button
          key={page}
          onClick={() => onGoTo(page)}
          className={`w-9 h-9 rounded-lg text-sm font-medium transition-colors ${
            page === currentPage
              ? 'bg-[#2f81f7] text-white'
              : 'bg-[#161b22] border border-[#30363d] text-[#c9d1d9] hover:border-[#2f81f7]'
          }`}
        >
          {page}
        </button>
      ))}

      <button
        onClick={onNext}
        disabled={!hasNext}
        className="px-4 py-2 rounded-lg bg-[#161b22] border border-[#30363d] text-[#c9d1d9] disabled:opacity-40 hover:border-[#2f81f7] transition-colors text-sm"
      >
        Next →
      </button>
    </div>
  );
}