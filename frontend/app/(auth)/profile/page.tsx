import GitHubProfile from "./components/githubprofile";
import TechPreferences from "./components/techpreferences";
import DomainPreferences from "./components/domainpreferences";
import DangerZone from "./components/dangerzone";

export default function ProfilePage() {
  return (
    <div className="space-y-10 max-w-3xl">
      {/* Header */}
      <div>
        <h1 className="text-2xl md:text-3xl font-semibold">
          Profile & preferences
        </h1>
        <p className="mt-2 text-[#8b949e]">
          Manage your GitHub connection and personalize your
          open-source recommendations.
        </p>
      </div>

      <GitHubProfile />
      <TechPreferences />
      <DomainPreferences />
      <DangerZone />
    </div>
  );
}
