import RepoCard from "../../discover/components/repocard"

export default function RepoGrid() {
  // mock array for now
  const repos = Array.from({ length: 6 });

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-6">
      {repos.map((_, idx) => (
        <RepoCard key={idx} />
      ))}
    </div>
  );
}
