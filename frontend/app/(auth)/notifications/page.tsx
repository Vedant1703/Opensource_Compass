"use client";

import { useNotifications } from "@/contexts/notification-context";
import { Bell, ExternalLink, Trash2 } from "lucide-react";
import Link from "next/link";
import PageWrapper from "@/components/ui/page-wrapper";
import { TimeAgo } from "@/components/ui/time-ago";

export default function NotificationsPage() {
  const { notifications, unreadCount, markAllAsSeen, clearAll, isConnected } = useNotifications();


  return (
    <PageWrapper className="min-h-screen text-[#c9d1d9] p-6">
      <div className="max-w-4xl mx-auto">
        {/* Header */}
        <div className="flex items-center justify-between mb-8">
          <div className="flex items-center gap-3">
            <Bell className="w-8 h-8 text-blue-500" />
            <div>
              <h1 className="text-3xl font-bold">Notifications</h1>
              <div className="flex items-center gap-3 mt-1">
                <p className="text-sm text-gray-400">
                  {unreadCount > 0 ? `${unreadCount} notification${unreadCount > 1 ? 's' : ''}` : 'No new notifications'}
                </p>
                <span className={`flex items-center gap-1.5 text-xs ${isConnected ? 'text-green-400' : 'text-red-400'}`}>
                  <span className={`w-2 h-2 rounded-full ${isConnected ? 'bg-green-400' : 'bg-red-400'} animate-pulse`}></span>
                  {isConnected ? 'Connected' : 'Disconnected'}
                </span>
              </div>
            </div>
          </div>
          
          {notifications.length > 0 && (
            <button
              onClick={clearAll}
              className="flex items-center gap-2 px-4 py-2 bg-red-500/10 hover:bg-red-500/20 text-red-400 rounded-lg transition-colors"
            >
              <Trash2 className="w-4 h-4" />
              Clear All
            </button>
          )}
        </div>

        {/* Notifications List */}
        {notifications.length === 0 ? (
          <div className="text-center py-16">
            <Bell className="w-16 h-16 text-gray-600 mx-auto mb-4" />
            <p className="text-lg text-gray-400">No notifications yet</p>
            <p className="text-sm text-gray-500 mt-2">
              You'll see notifications here when new issues are created in your watched repositories
            </p>
          </div>
        ) : (
          <div className="space-y-3">
            {notifications.map((notification) => (
              <div
                key={notification.id}
                className={`bg-[#161b22]/60 backdrop-blur-md border border-[#30363d] rounded-lg p-4 hover:border-blue-500/50 transition-colors ${!notification.isRead ? 'border-l-2 border-l-[#2f81f7]' : ''}`}
              >
                <div className="flex items-start justify-between gap-4">
                  <div className="flex-1">
                    <div className="flex items-center gap-2 mb-2">
                      <Bell className="w-4 h-4 text-blue-500 flex-shrink-0" />
                      {!notification.isRead && <span className="inline-block w-1.5 h-1.5 rounded-full bg-[#2f81f7] mr-1" />}
                      <span className="text-sm font-semibold text-blue-400">
                        {notification.repo}
                      </span>
                      <TimeAgo 
                        date={notification.timestamp} 
                        className="text-xs text-gray-500" 
                      />
                    </div>
                    
                    <div className="mb-2">
                      <p className="text-base font-medium text-gray-200">
                        {notification.issue_title || `Issue #${notification.issue_number}`}
                      </p>
                      <p className="text-sm text-gray-400 mt-1">
                        {notification.message}
                      </p>
                    </div>

                    {notification.issue_url && (
                      <a
                        href={notification.issue_url}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="inline-flex items-center gap-1.5 text-sm text-blue-400 hover:text-blue-300 transition-colors"
                      >
                        View Issue #{notification.issue_number}
                        <ExternalLink className="w-3.5 h-3.5" />
                      </a>
                    )}
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </PageWrapper>
  );
}
