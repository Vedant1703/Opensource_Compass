"use client";

import { useEffect, useState, useCallback } from "react";
import { searchRepositories, searchRepositoriesByName, Repository } from "@/lib/api/github-service";
import { getPreferences } from "@/lib/api/preferences";
import DiscoverHero from "./components/discoverhero";
import SearchInput from "./components/search-input";
import ActiveFilters from "./components/activefilters";
import RepoGrid from "./components/repogrid";
import SkeletonCard from "./components/skeletoncard";
import EmptyState from "./components/emptystate";
import PageWrapper from "@/components/ui/page-wrapper";
import { usePagination } from '@/lib/hooks/usePagination';
import { PaginationControls } from '@/components/ui/pagination';

export default function DiscoverPage() {
  const [repos, setRepos] = useState<Repository[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Preferences
  const [languages, setLanguages] = useState<string[]>([]);
  const [topics, setTopics] = useState<string[]>([]);
  const [experienceLevel, setExperienceLevel] = useState<string>("Beginner");
  const [username, setUsername] = useState<string | null>(null);
  const [searchQuery, setSearchQuery] = useState("");

  const loadRepos = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);

      // Fetch user preferences from backend
      const prefs = await getPreferences();

      setLanguages(prefs.languages || []);
      setTopics(prefs.topics || []);
      setExperienceLevel(prefs.experienceLevel || "Beginner");

      // Fetch repositories using preferences
      const repositories = await searchRepositories(prefs.languages, [], prefs.topics);
      setRepos(repositories);
    } catch (err: any) {
      console.error("Failed to load repositories:", err);
      // If error is 401, it will be handled by the API client or global error boundary
      // For now, we just show a generic error
      setError("Failed to load repositories. Please check your connection.");
    } finally {
      setLoading(false);
    }
  }, []);

  const handleSearch = useCallback(async (query: string) => {
    setSearchQuery(query);
    if (!query) {
      loadRepos(); // reload default recommendations
      return;
    }

    try {
      setLoading(true);
      setError(null);
      const results = await searchRepositoriesByName(query);
      setRepos(results);
    } catch (err: any) {
      console.error("Failed to search repositories:", err);
      setError("Failed to search repositories. Please try again.");
      setRepos([]);
    } finally {
      setLoading(false);
    }
  }, [loadRepos]);

  useEffect(() => {
    loadRepos();
  }, []);

  // Check if user has set preferences
  const hasPreferences = languages.length > 0 || topics.length > 0;
  // Show filters only if not searching and has preferences
  const showFilters = !searchQuery && hasPreferences;

  const { paginatedItems, currentPage, totalPages, goToPage, nextPage, prevPage, hasPrev, hasNext } = usePagination({
    items: repos,
    pageSize: 9,
  });

  return (
    <PageWrapper className="space-y-6">
      {/* Hero Section */}
      <DiscoverHero
        username={username}
        repoCount={repos.length}
        onRefresh={loadRepos}
        isLoading={loading}
      />

      {/* Search Input */}
      <div className="max-w-5xl mx-auto px-4 my-8">
        <SearchInput onSearch={handleSearch} />
      </div>

      {/* Active Filters */}
      {showFilters && !loading && (
        <ActiveFilters
          languages={languages}
          topics={topics}
          experienceLevel={experienceLevel}
        />
      )}

      {/* Loading State */}
      {loading && (
        <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-6">
          {[...Array(6)].map((_, i) => (
            <SkeletonCard key={i} />
          ))}
        </div>
      )}

      {/* Error State */}
      {error && !loading && (
        <EmptyState type="error" onRetry={loadRepos} />
      )}

      {/* No Preferences State */}
      {!hasPreferences && !loading && !error && (
        <EmptyState type="no-preferences" />
      )}

      {/* No Repos Found State */}
      {hasPreferences && repos.length === 0 && !loading && !error && (
        <EmptyState type="no-repos" />
      )}

      {/* Repo Grid */}
      {repos.length > 0 && !loading && !error && (
        <>
          <RepoGrid repos={paginatedItems} />
          <PaginationControls
            currentPage={currentPage}
            totalPages={totalPages}
            onPrev={prevPage}
            onNext={nextPage}
            onGoTo={goToPage}
            hasPrev={hasPrev}
            hasNext={hasNext}
          />
        </>
      )}
    </PageWrapper>
  );
}
