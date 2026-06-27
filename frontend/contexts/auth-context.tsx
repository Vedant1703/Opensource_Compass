"use client";

import { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { useRouter } from 'next/navigation';

interface AuthContextType {
  isAuthenticated: boolean;
  isLoading: boolean;
  userId: string | null;
  username: string | null;
  avatar: string | null;
  email: string | null;
  token: string | null;
  user: { username: string; avatar: string | null; email: string | null; id?: string } | null;
  checkAuth: () => Promise<void>;
  logout: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType | null>(null);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [username, setUsername] = useState<string | null>(null);
  const [userId, setUserId] = useState<string | null>(null);
  const [avatar, setAvatar] = useState<string | null>(null);
  const [email, setEmail] = useState<string | null>(null);
  const [token, setToken] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const router = useRouter();

  const checkAuth = async () => {
    try {
      setIsLoading(true);
      
      // Check for token in URL (from OAuth redirect)
      if (typeof window !== 'undefined') {
        const params = new URLSearchParams(window.location.search);
        const urlToken = params.get('token');
        if (urlToken) {
           localStorage.setItem('auth_token', urlToken);
           setToken(urlToken);
           // Clear token from URL
           window.history.replaceState({}, document.title, window.location.pathname);
        }
      }

      const storedToken = typeof window !== 'undefined' ? localStorage.getItem('auth_token') : null;
      if (storedToken) setToken(storedToken);

      const headers: HeadersInit = {};
      if (storedToken) {
          headers['Authorization'] = `Bearer ${storedToken}`;
      }

      const authUrl = `${process.env.NEXT_PUBLIC_AUTH_SERVICE_URL}/auth/me`;
      const res = await fetch(authUrl, { headers });
      if (res.ok) {
        const data = await res.json();
        setUserId(data.id);
        setUsername(data.username);
        setAvatar(data.avatar);
        setEmail(data.email);
        
        // Ensure token is consistent
        if (data.token) {
            setToken(data.token);
            localStorage.setItem('auth_token', data.token);
        }
        
        setIsAuthenticated(true);
      } else {
        // If 401, clear storage
        localStorage.removeItem('auth_token');
        setUserId(null);
        setUsername(null);
        setAvatar(null);
        setEmail(null);
        setUserId(null);
        setToken(null);
        setIsAuthenticated(false);
      }
    } catch (error) {
      console.error("Auth check failed:", error);
      setIsAuthenticated(false);
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    checkAuth();
  }, []);

  const logout = async () => {
    try {
      const logoutUrl = `${process.env.NEXT_PUBLIC_AUTH_SERVICE_URL}/auth/logout`;
      await fetch(logoutUrl, { method: 'POST' });
    } catch (error) {
      console.error("Logout failed:", error);
    } finally {
      // Always clear local state, even if the API call fails
      localStorage.removeItem('auth_token');
      setUserId(null);
      setUsername(null);
      setAvatar(null);
      setEmail(null);
      setToken(null);
      setIsAuthenticated(false);
      router.push('/');
    }
  };

  return (
    <AuthContext.Provider
      value={{
        isAuthenticated,
        isLoading,
        userId,
        username,
        avatar,
        email,
        token,
        user: username ? { username, avatar, email, id: userId || username } : null,
        checkAuth,
        logout,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within AuthProvider');
  }
  return context;
};
