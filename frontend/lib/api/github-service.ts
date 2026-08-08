// GitHub Service API Client

const GITHUB_SERVICE_URL = process.env.NEXT_PUBLIC_GITHUB_SERVICE_URL || 'http://localhost:8081';

export interface Repository {
    owner: string;
    name: string;
    full_name: string;
    description?: string;
    html_url: string;
    stargazers_count?: number;
    language?: string;
    updated_at?: string;
    topics?: string[];
}

export async function searchRepositories(
    languages: string[],
    frameworks: string[] = [],
    domains: string[] = [],
    page: number = 1,
    limit: number = 9
): Promise<Repository[]> {
    try {
        const params = new URLSearchParams();

        if (languages.length > 0) {
            params.append('languages', languages.map(l => l.toLowerCase()).join(','));
        }

        if (frameworks.length > 0) {
            params.append('frameworks', frameworks.join(','));
        }

        if (domains.length > 0) {
            params.append('domain', domains.map(d => d.toLowerCase()).join(','));
        }
        
        params.append('page', page.toString());
        params.append('limit', limit.toString());

        // Use full backend URL for production
        const url = `${GITHUB_SERVICE_URL}/repos/search?${params.toString()}`;

        console.log("Calling GitHub service:", url);

        const response = await fetch(url);

        if (!response.ok) {
            throw new Error(`Failed to fetch repositories: ${response.statusText}`);
        }

        const data = await response.json();
        console.log("GitHub service response:", data);
        return data;
    } catch (error) {
        console.error('Error fetching repositories:', error);
        throw error;
    }
}

export async function fetchRepository(owner: string, repo: string): Promise<Repository> {
    try {
        const url = `${GITHUB_SERVICE_URL}/repos/${owner}/${repo}`;

        const response = await fetch(url);

        if (!response.ok) {
            throw new Error(`Failed to fetch repository: ${response.statusText}`);
        }

        const data = await response.json();
        return data;
    } catch (error) {
        console.error('Error fetching repository:', error);
        throw error;
    }
}

export async function searchRepositoriesByName(query: string, page: number = 1, limit: number = 9): Promise<Repository[]> {
    try {
        const token = localStorage.getItem('auth_token');
        const headers: HeadersInit = {
            'Content-Type': 'application/json',
        };

        if (token) {
            headers['Authorization'] = `Bearer ${token}`;
        }

        const CORE_SERVICE_URL = process.env.NEXT_PUBLIC_CORE_SERVICE_URL;
        const url = `${CORE_SERVICE_URL}/repos/search?q=${encodeURIComponent(query)}&page=${page}&limit=${limit}`;
        console.log("Calling Core service:", url);

        const response = await fetch(url, {
            headers: headers
        });

        if (!response.ok) {
            throw new Error(`Failed to search repositories: ${response.statusText}`);
        }

        const data = await response.json();
        return data;
    } catch (error) {
        console.error('Error searching repositories:', error);
        throw error;
    }
}
