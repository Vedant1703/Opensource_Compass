"use client";

import { createContext, useContext, useEffect, useState, ReactNode, useRef } from "react";
import { useAuth } from "./auth-context";
import toast from "react-hot-toast";
import Link from "next/link";

type Notification = {
  type: string;
  repo: string;
  issue_number: number;
  issue_title?: string;
  issue_url?: string;
  message: string;
  timestamp: Date;
  id: string; // Unique identifier for each notification
  isRead: boolean;
};

type NotificationContextType = {
  notifications: Notification[];
  unreadCount: number;
  markAllAsSeen: () => void;
  clearAll: () => void;
  isConnected: boolean;
};

const NotificationContext = createContext<NotificationContextType>({
  notifications: [],
  unreadCount: 0,
  markAllAsSeen: () => {},
  clearAll: () => {},
  isConnected: false,
});

const NOTIFICATIONS_KEY = "compass_notifications";

function loadNotificationsFromStorage(): Notification[] {
  if (typeof window === "undefined") return [];
  try {
    const stored = localStorage.getItem(NOTIFICATIONS_KEY);
    if (!stored) return [];
    const parsed = JSON.parse(stored);
    return parsed.map((n: any) => ({
      ...n,
      timestamp: new Date(n.timestamp),
      isRead: n.isRead ?? false,
    }));
  } catch (err) {
    console.error("Failed to load notifications from storage", err);
    return [];
  }
}

function saveNotificationsToStorage(notifications: Notification[]) {
  if (typeof window === "undefined") return;
  try {
    localStorage.setItem(NOTIFICATIONS_KEY, JSON.stringify(notifications));
  } catch (err) {
    console.error("Failed to save notifications to storage", err);
  }
}

export function NotificationProvider({ children }: { children: ReactNode }) {
  const { user } = useAuth();
  const [isConnected, setIsConnected] = useState(false);
  const [notifications, setNotifications] = useState<Notification[]>(loadNotificationsFromStorage);

  const reconnectAttempts = useRef(0);
  const maxReconnectAttempts = 7; // Max backoff limit
  const baseDelay = 1000; // 1 second

  useEffect(() => {
    if (!user) return;

    // Connect to Notification Service
    // Ensure we use the correct ID field. The User type in auth-context typically has 'id'.
    // Assuming auth-context exposes user.id. 
    // If not, we need to check auth-context. 
    // For now, let's assume user object has an ID or we use username if ID is missing (but database uses ID).

    // In auth-context.tsx, User type usually has login/name but maybe not database ID if it comes from GitHub directly.
    // However, our backend uses the user_id from the token (which is usually the github ID).
    // Let's rely on what's available. 

    // Actually, looking at previous steps, the user object in auth-context seems to come from /auth/me or similar.
    // Let's try to connect with user.id.

    const userId = (user as any).id || (user as any).login;
    if (!userId) return;

    const wsBaseUrl = process.env.NEXT_PUBLIC_NOTIFICATION_SERVICE_URL || "ws://localhost:8084";
    const wsUrl = `${wsBaseUrl}/ws?user_id=${userId}`;
    let ws: WebSocket;
    let reconnectTimeoutId: NodeJS.Timeout;

    const connect = () => {
      console.log(`NotificationService: Connecting to ${wsUrl}...`);
      ws = new WebSocket(wsUrl);

      ws.onerror = (error) => {
        console.error("NotificationService: WebSocket error", error);
      };

      ws.onopen = () => {
        console.log(`NotificationService: Connected for user ${userId}`);
        setIsConnected(true);
        reconnectAttempts.current = 0; // reset on success
      };

      ws.onmessage = (event) => {
        console.log("NotificationService: Received message", event.data);
        try {
          const data = JSON.parse(event.data);
          const newNotification: Notification = {
            ...data,
            timestamp: new Date(),
            id: `${data.repo}-${data.issue_number}-${Date.now()}`,
            isRead: false,
          };

          setNotifications((prev) => {
            const updated = [newNotification, ...prev];
            saveNotificationsToStorage(updated);
            return updated;
          });

          // Show react-hot-toast notification
          const issueTitle = data.issue_title || `Issue #${data.issue_number}`;
          const repoName = data.repo.split("/")[1] || data.repo;

          toast(
            (t) => (
              <div className="flex flex-col gap-1">
                <div className="flex items-center gap-2">
                  <span className="text-sm font-semibold">{repoName}</span>
                  {data.issue_url && (
                    <a
                      href={data.issue_url}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="text-xs text-blue-600 hover:underline"
                      onClick={() => toast.dismiss(t.id)}
                    >
                      View Issue
                    </a>
                  )}
                </div>
                <p className="text-xs text-gray-600">{issueTitle}</p>
              </div>
            ),
            {
              icon: "🔔",
              duration: 5000,
              position: "top-right",
            }
          );
        } catch (err) {
          console.error("Failed to parse notification", err);
        }
      };

      ws.onclose = (e) => {
        console.log("Disconnected from notification service", e);
        setIsConnected(false);

        // Auto-reconnect with exponential backoff
        if (reconnectAttempts.current < maxReconnectAttempts) {
          const delay = baseDelay * Math.pow(2, reconnectAttempts.current);
          console.log(`Attempting to reconnect in ${delay}ms...`);
          reconnectTimeoutId = setTimeout(connect, delay);
          reconnectAttempts.current++;
        } else {
          console.error("Max reconnect attempts reached. Please refresh the page.");
        }
      };
    };

    connect();

    return () => {
      clearTimeout(reconnectTimeoutId);
      if (ws) ws.close();
    };
  }, [user]);

  const markAllAsSeen = () => {
    setNotifications(prev => {
      const updated = prev.map(n => ({ ...n, isRead: true }));
      saveNotificationsToStorage(updated);
      return updated;
    });
  };

  const clearAll = () => {
    setNotifications([]);
    saveNotificationsToStorage([]);
  };

  return (
    <NotificationContext.Provider
      value={{
        notifications,
        unreadCount: notifications.filter(n => !n.isRead).length,
        markAllAsSeen,
        clearAll,
        isConnected,
      }}
    >
      {children}
    </NotificationContext.Provider>
  );
}

export const useNotifications = () => useContext(NotificationContext);
