"use client";

import { useEffect, useState } from "react";
import { useAuth } from "@/contexts/auth-context";
import { getPreferences, UserPreferences } from "@/lib/api/preferences";
import { useRouter } from "next/navigation";
import { Github, LogOut, Settings, Loader2, ExternalLink, Bell, RefreshCw } from "lucide-react";
import Link from "next/link";
import { getWatchlist, WatchedRepo } from "@/lib/api/watchlist";
import { WatchButton } from "@/components/watch-button";
import { syncEmail } from "@/lib/api/preferences";
import { toast } from "react-hot-toast";
import { TimeAgo } from "@/components/ui/time-ago";

import PageWrapper from "@/components/ui/page-wrapper";

export default function ProfilePage() {
  const { user, username, avatar, logout, isLoading: authLoading } = useAuth();
  const router = useRouter();
  const [preferences, setPreferences] = useState<UserPreferences | null>(null);
  const [loading, setLoading] = useState(true);
  const [syncing, setSyncing] = useState(false);

  useEffect(() => {
    const loadPreferences = async () => {
      try {
        const prefs = await getPreferences();
        setPreferences(prefs);
      } catch (error) {
        console.error("Failed to load preferences:", error);
      } finally {
        setLoading(false);
      }
    };

    if (!authLoading) {
      loadPreferences();
    }
  }, [authLoading]);

  const handleLogout = async () => {
    await logout();
  };

  if (authLoading || loading) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <Loader2 className="h-8 w-8 animate-spin text-[#2f81f7]" />
      </div>
    );
  }

  return (
    <PageWrapper>
      <div className="space-y-6 max-w-4xl mx-auto px-6 pt-6">
        {/* Header */}
        <div>
          <h1 className="text-2xl md:text-3xl font-semibold text-[#c9d1d9]">
            Profile
          </h1>
          <p className="mt-2 text-[#8b949e]">
            View your account information and preferences
          </p>
        </div>

        {/* Profile Header Card */}
        <div className="bg-[#0d1117]/60 backdrop-blur-md border border-[#30363d] rounded-xl p-6">
          <div className="flex items-start justify-between">
            <div className="flex items-center gap-4">
              {/* Avatar */}
              <div className="h-20 w-20 rounded-full bg-[#161b22] border-2 border-[#30363d] flex items-center justify-center overflow-hidden">
                {avatar ? (
                  <img
                    src={avatar}
                    alt={username || "User"}
                    className="h-full w-full object-cover"
                  />
                ) : (
                  <div className="text-2xl font-semibold text-[#8b949e]">
                    {username?.[0]?.toUpperCase() || "U"}
                  </div>
                )}
              </div>

              {/* User Info */}
              <div>
                <h2 className="text-xl font-semibold text-[#c9d1d9]">
                  {username || "User"}
                </h2>
                <div className="flex items-center gap-2 mt-1">
                  <Github className="h-4 w-4 text-[#8b949e]" />
                  <span className="text-sm text-[#8b949e]">
                    Connected via GitHub
                  </span>
                </div>
                {user?.email && (
                  <div className="text-sm text-[#8b949e] mt-0.5">
                    {user.email}
                  </div>
                )}
                {preferences?.experienceLevel && (
                  <div className="mt-2">
                    <span className="inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-medium bg-[#238636]/10 border border-[#238636]/20 text-[#3fb950]">
                      {preferences.experienceLevel}
                    </span>
                  </div>
                )}
              </div>
            </div>

            {/* Actions */}
            <div className="flex gap-2">
              <Link
                href="/settings"
                className="inline-flex items-center gap-2 px-4 py-2 rounded-lg bg-[#21262d] border border-[#30363d] text-sm text-[#c9d1d9] hover:bg-[#30363d] transition-colors"
              >
                <Settings className="h-4 w-4" />
                Edit
              </Link>
              <button
                onClick={handleLogout}
                className="inline-flex items-center gap-2 px-4 py-2 rounded-lg bg-[#21262d] border border-[#30363d] text-sm text-[#c9d1d9] hover:bg-[#da3633]/10 hover:border-[#da3633]/30 hover:text-[#f85149] transition-colors"
              >
                <LogOut className="h-4 w-4" />
                Logout
              </button>
            </div>
          </div>
        </div>

        {/* Preferences Section */}
        <div className="bg-[#0d1117]/60 backdrop-blur-md border border-[#30363d] rounded-xl p-6">
          <h3 className="text-lg font-semibold text-[#c9d1d9] mb-4">
            Your Preferences
          </h3>

          {/* Languages */}
          {preferences?.languages && preferences.languages.length > 0 && (
            <div className="mb-6">
              <h4 className="text-sm font-medium text-[#8b949e] mb-3">
                Programming Languages
              </h4>
              <div className="flex flex-wrap gap-2">
                {preferences.languages.map((lang) => (
                  <span
                    key={lang}
                    className="px-3 py-1.5 rounded-lg bg-[#161b22] border border-[#30363d] text-sm text-[#c9d1d9]"
                  >
                    {lang}
                  </span>
                ))}
              </div>
            </div>
          )}

          {/* Topics */}
          {preferences?.topics && preferences.topics.length > 0 && (
            <div>
              <h4 className="text-sm font-medium text-[#8b949e] mb-3">
                Topics & Domains
              </h4>
              <div className="flex flex-wrap gap-2">
                {preferences.topics.map((topic) => (
                  <span
                    key={topic}
                    className="px-3 py-1.5 rounded-lg bg-[#161b22] border border-[#30363d] text-sm text-[#c9d1d9]"
                  >
                    {topic}
                  </span>
                ))}
              </div>
            </div>
          )}

          {/* Empty State */}
          {(!preferences?.languages || preferences.languages.length === 0) &&
            (!preferences?.topics || preferences.topics.length === 0) && (
              <div className="text-center py-8">
                <p className="text-[#8b949e] mb-4">
                  No preferences set yet
                </p>
                <Link
                  href="/settings"
                  className="inline-flex items-center gap-2 px-4 py-2 rounded-lg bg-[#238636] text-white hover:bg-[#2ea043] transition-colors"
                >
                  <Settings className="h-4 w-4" />
                  Set Your Preferences
                </Link>
              </div>
            )}
        </div>

        {/* GitHub Integration Status */}
        <div className="bg-[#0d1117]/60 backdrop-blur-md border border-[#30363d] rounded-xl p-6">
          <h3 className="text-lg font-semibold text-[#c9d1d9] mb-4">
            GitHub Integration
          </h3>
          <div className="flex items-center gap-3">
            <div className="h-10 w-10 rounded-full bg-[#238636]/10 border border-[#238636]/20 flex items-center justify-center">
              <Github className="h-5 w-5 text-[#3fb950]" />
            </div>
            <div>
              <p className="text-sm font-medium text-[#c9d1d9]">
                GitHub Account Connected
              </p>
              <p className="text-xs text-[#8b949e]">
                OAuth access token active
              </p>
            </div>
            <button
              onClick={async () => {
                if (!user?.id) return;
                setSyncing(true);
                try {
                  const result = await syncEmail(user.id);
                  toast.success(`Email synced: ${result.email}`);
                } catch (error) {
                  toast.error("Failed to sync email. Please try logging in again.");
                } finally {
                  setSyncing(false);
                }
              }}
              disabled={syncing}
              className="ml-auto inline-flex items-center gap-2 px-3 py-1.5 rounded-lg bg-[#21262d] border border-[#30363d] text-xs text-[#c9d1d9] hover:bg-[#30363d] transition-colors disabled:opacity-50"
            >
              <RefreshCw className={`h-3 w-3 ${syncing ? "animate-spin" : ""}`} />
              {syncing ? "Syncing..." : "Sync Email"}
            </button>
          </div>
        </div>

        {/* Watchlist Section */}
        <div className="bg-[#0d1117]/60 backdrop-blur-md border border-[#30363d] rounded-xl p-6">
          <h3 className="text-lg font-semibold text-[#c9d1d9] mb-4">
            Watched Repositories
          </h3>
          <Watchlist />
        </div>
      </div>
    </PageWrapper>
  );
}

function Watchlist() {
  const { token } = useAuth();
  const [repos, setRepos] = useState<WatchedRepo[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (token) {
      getWatchlist(token)
        .then(setRepos)
        .catch(console.error)
        .finally(() => setLoading(false));
    }
  }, [token]);

  if (loading) return <div className="text-sm text-[#8b949e]">Loading watchlist...</div>;
  if (!repos || repos.length === 0) return <div className="text-sm text-[#8b949e]">You are not watching any repositories.</div>;

  return (
    
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        {repos.map((repo) => (
          <div key={repo.id} className="bg-[#161b22] border border-[#30363d] rounded-lg p-4 hover:border-[#58a6ff] transition-colors">
            <div className="flex justify-between items-center mb-2 gap-4">
              <div className="min-w-0 flex-1">
                <Link href={`https://github.com/${repo.repo_owner}/${repo.repo_name}`} target="_blank" className="font-semibold text-[#58a6ff] hover:underline flex items-center gap-1 truncate">
                  <span className="truncate">{repo.repo_owner}/{repo.repo_name}</span>
                  <ExternalLink className="w-3 h-3 shrink-0" />
                </Link>
                <div className="text-xs text-[#8b949e] mt-1 flex items-center gap-2">
                    <span className="flex items-center gap-1"><Bell className="w-3 h-3"/> Last checked: {repo.last_checked_at ? <TimeAgo date={repo.last_checked_at} />: 'Never'}</span>
                </div>
              </div>
              <div className="shrink-0">
                <WatchButton owner={repo.repo_owner} name={repo.repo_name} />
              </div>
            </div>
          </div>
        ))}
      </div>
  );
}