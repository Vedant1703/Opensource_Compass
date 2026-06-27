// Core Service API Client

// Use environment variable for backend URL
const CORE_SERVICE_URL = process.env.NEXT_PUBLIC_CORE_SERVICE_URL;

export interface Issue {
    title: string;
    body: string;
    labels: string[];
    issue_number?: number;
    issue_url?: string;
    ai?: {
        summary?: string;
        difficulty?: string;
        estimated_time?: string;
    };
}

export interface Recommendation {
    repo_id: string;
    score: number;
    level: string;
    reasons: string[];
    suggested_repos: string[];
    issues: Issue[];
}

export async function getRecommendations(): Promise<Recommendation> {
    try {
        const token = typeof window !== 'undefined' ? localStorage.getItem('auth_token') : null;
        const headers: HeadersInit = {
            'Content-Type': 'application/json',
        };
        if (token) {
            headers['Authorization'] = `Bearer ${token}`;
        }

        const response = await fetch(`${CORE_SERVICE_URL}/recommendations`, {
            headers,
        });

        if (!response.ok) {
            if (response.status === 401) {
                throw new Error('Unauthorized - please login');
            }
            throw new Error(`Failed to fetch recommendations: ${response.statusText}`);
        }

        const data = await response.json();
        return data;
    } catch (error) {
        console.error('Error fetching recommendations:', error);
        throw error;
    }
}

export function getScoreColor(score: number): string {
    if (score >= 90) return '#3fb950'; // Excellent - green
    if (score >= 70) return '#58a6ff'; // Good - blue
    if (score >= 50) return '#d29922'; // Fair - yellow
    return '#f85149'; // Poor - red
}

export function getLevelBadgeColor(level: string): { bg: string; border: string; text: string } {
    const lowerLevel = level.toLowerCase();

    if (lowerLevel.includes('beginner')) {
        return {
            bg: '#3fb950/10',
            border: '#3fb950/20',
            text: '#3fb950'
        };
    }

    if (lowerLevel.includes('intermediate')) {
        return {
            bg: '#58a6ff/10',
            border: '#58a6ff/20',
            text: '#58a6ff'
        };
    }

    // Advanced
    return {
        bg: '#a371f7/10',
        border: '#a371f7/20',
        text: '#a371f7'
    };
}

export interface SetupGuide {
    prerequisites: string[];
    installation_steps: {
        description: string;
        command: string;
    }[];
    environment_variables: {
        name: string;
        description: string;
        required: boolean;
    }[];
    common_errors: {
        error: string;
        fix: string;
    }[];
    contribution_tips: string[];
}

export async function getSetupGuide(owner: string, repo: string, userId: string): Promise<SetupGuide> {
    const token = typeof window !== 'undefined' ? localStorage.getItem('auth_token') : null;
    const headers: HeadersInit = {
        'Content-Type': 'application/json',
    };
    if (token) {
        headers['Authorization'] = `Bearer ${token}`;
    }

    const encodedUserId = encodeURIComponent(userId);
    const response = await fetch(`${CORE_SERVICE_URL}/repo/setup-guide?repo=${owner}/${repo}&user_id=${encodedUserId}`, {
        method: 'GET',
        headers,
    });

    if (!response.ok) {
        throw new Error(`Failed to generate setup guide: ${response.statusText}`);
    }

    const data = await response.json();
    
    // Parse the markdown JSON string
    let jsonString = data.guide;

    // Remove markdown code block markers
    if (jsonString.startsWith('```json')) {
        jsonString = jsonString.replace(/^```json\s*/, '').replace(/\s*```$/, '');
    } else if (jsonString.startsWith('```')) {
         jsonString = jsonString.replace(/^```\s*/, '').replace(/\s*```$/, '');
    }

    try {
        return JSON.parse(jsonString);
    } catch (e) {
        console.error("Failed to parse guide JSON", e);
        throw new Error("Invalid response format from AI service");
    }
}