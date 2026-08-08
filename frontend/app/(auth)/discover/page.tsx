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
  const [currentPage, setCurrentPage] = useState(1);
  const [hasNextPage, setHasNextPage] = useState(false);
  const PAGE_SIZE = 9;

  const loadRepos = useCallback(async (page: number = 1) => {
    try {
      setLoading(true);
      setError(null);

      // Fetch user preferences from backend
      const prefs = await getPreferences();

      setLanguages(prefs.languages || []);
      setTopics(prefs.topics || []);
      setExperienceLevel(prefs.experienceLevel || "Beginner");

      // Fetch repositories using preferences
      const repositories = await searchRepositories(prefs.languages, [], prefs.topics, page, PAGE_SIZE);
      setRepos(repositories);
      setHasNextPage(repositories.length === PAGE_SIZE);
      setCurrentPage(page);
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
      loadRepos(1); // reload default recommendations
      return;
    }

    try {
      setLoading(true);
      setError(null);
      const results = await searchRepositoriesByName(query, 1, PAGE_SIZE);
      setRepos(results);
      setHasNextPage(results.length === PAGE_SIZE);
      setCurrentPage(1);
    } catch (err: any) {
      console.error("Failed to search repositories:", err);
      setError("Failed to search repositories. Please try again.");
      setRepos([]);
    } finally {
      setLoading(false);
    }
  }, [loadRepos]);

  useEffect(() => {
    loadRepos(1);
  }, []);

  const handlePageChange = async (newPage: number) => {
    if (searchQuery) {
      try {
        setLoading(true);
        const results = await searchRepositoriesByName(searchQuery, newPage, PAGE_SIZE);
        setRepos(results);
        setHasNextPage(results.length === PAGE_SIZE);
        setCurrentPage(newPage);
        window.scrollTo({ top: 0, behavior: 'smooth' });
      } catch (err: any) {
        console.error("Failed to load page:", err);
      } finally {
        setLoading(false);
      }
    } else {
      await loadRepos(newPage);
      window.scrollTo({ top: 0, behavior: 'smooth' });
    }
  };

  // Check if user has set preferences
  const hasPreferences = languages.length > 0 || topics.length > 0;
  // Show filters only if not searching and has preferences
  const showFilters = !searchQuery && hasPreferences;

  const hasPrev = currentPage > 1;
  const hasNext = hasNextPage;
  const totalPages = currentPage + (hasNext ? 1 : 0); // Approximate total pages for PaginationControls

  return (
    <PageWrapper className="space-y-6">
      {/* Hero Section */}
      <DiscoverHero
        username={username}
        repoCount={repos.length}
        onRefresh={() => loadRepos(currentPage)}
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
        <EmptyState type="error" onRetry={() => loadRepos(currentPage)} />
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
          <RepoGrid repos={repos} />
          <PaginationControls
            currentPage={currentPage}
            totalPages={totalPages}
            onPrev={() => handlePageChange(currentPage - 1)}
            onNext={() => handlePageChange(currentPage + 1)}
            onGoTo={(page: number) => handlePageChange(page)}
            hasPrev={hasPrev}
            hasNext={hasNext}
          />
        </>
      )}
    </PageWrapper>
  );
}
