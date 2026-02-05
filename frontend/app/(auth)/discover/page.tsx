import FilterBar from "../issues/components/filterbar";
import RepoGrid from "../issues/components/repogrid";

export default function DiscoverPage() {
  return (
    <div className="space-y-8">
      {/* Page Header */}
      <div>
        <h1 className="text-2xl md:text-3xl font-semibold tracking-tight">
          Discover repositories
        </h1>
        <p className="mt-2 text-[#8b949e] max-w-2xl">
          Personalized open-source projects based on your tech stack
          and interests. Beginner-friendly and actively maintained.
        </p>
      </div>

      {/* Filter Bar */}
      <FilterBar />

      {/* Repo Grid */}
      <RepoGrid />
    </div>
  );
}
