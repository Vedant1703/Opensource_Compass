"use client";
import React, { useState } from "react";
import { Sidebar, SidebarBody, SidebarLink } from "./ui/sidebar";
import {
  IconArrowLeft,
  IconBrandTabler,
  IconBug,
  IconSearch,
  IconSettings,
  IconUserBolt,
  IconStars,
  IconBell,
} from "@tabler/icons-react";
import { motion } from "motion/react";
import { useAuth } from "@/contexts/auth-context";
import { useNotifications } from "@/contexts/notification-context";

export default function SidebarComponent() {
  const { logout, username, avatar } = useAuth();
  const { unreadCount } = useNotifications();
  
  const handleLogout = () => {
    logout();
  };
  
  const links = [
    {
      label: "Discover",
      href: "/discover",
      icon: (
        <IconSearch className="h-5 w-5 shrink-0 text-[#8b949e]" />
      ),
    },
    {
      label: "Recommendations",
      href: "/recommendations",
      icon: (
        <IconStars className="h-5 w-5 shrink-0 text-[#8b949e]" />
      ),
    },
    {
      label: "Issues",
      href: "/issues",
      icon: (
        <IconBug className="h-5 w-5 shrink-0 text-[#8b949e]" />
      ),
    },
    {
      label: "Notifications",
      href: "/notifications",
      icon: (
        <div className="relative">
          <IconBell className="h-5 w-5 shrink-0 text-[#8b949e]" />
          {unreadCount > 0 && (
            <span className="absolute -top-1 -right-1 h-4 w-4 rounded-full bg-red-500 text-[10px] font-semibold text-white flex items-center justify-center">
              {unreadCount > 9 ? '9+' : unreadCount}
            </span>
          )}
        </div>
      ),
    },
    {
      label: "Profile",
      href: "/profile",
      icon: (
        <IconUserBolt className="h-5 w-5 shrink-0 text-[#8b949e]" />
      ),
    },
    {
      label: "Settings",
      href: "/settings",
      icon: (
        <IconSettings className="h-5 w-5 shrink-0 text-[#8b949e]" />
      ),
    },
  ];
  const [open, setOpen] = useState(false);
  return (
    <Sidebar open={open} setOpen={setOpen}>
      <SidebarBody className="justify-between gap-10">
        <div className="flex flex-1 flex-col overflow-x-hidden overflow-y-auto">
          {open ? <Logo /> : <LogoIcon />}
          <div className="mt-8 flex flex-col gap-2">
            {links.map((link, idx) => (
              <SidebarLink key={idx} link={link} />
            ))}
            <button 
              onClick={handleLogout}
              className="flex items-center gap-2 px-2 py-2 text-sm text-[#8b949e] hover:text-[#c9d1d9] hover:bg-[#161b22] rounded-md group w-full text-left"
            >
              <IconArrowLeft className="h-5 w-5 shrink-0 text-[#8b949e]" />
              {open && <span>Logout</span>}
            </button>
          </div>
        </div>
        <div>
          <SidebarLink
            link={{
              label: username || "User",
              href: "/profile",
              icon: avatar ? (
                <img
                  src={avatar}
                  className="h-7 w-7 shrink-0 rounded-full"
                  width={50}
                  height={50}
                  alt="Avatar"
                />
              ) : (
                <div className="h-7 w-7 shrink-0 rounded-full bg-[#30363d] flex items-center justify-center text-xs">
                  {username?.[0]?.toUpperCase() || "U"}
                </div>
              ),
            }}
          />
        </div>
      </SidebarBody>
    </Sidebar>
  );
}
export const Logo = () => {
  return (
    <a
      href="/"
      className="relative z-20 flex items-center space-x-2 py-1 text-sm font-normal text-[#c9d1d9]"
    >
      <div className="h-5 w-6 shrink-0 rounded-tl-lg rounded-tr-sm rounded-br-lg rounded-bl-sm bg-[#c9d1d9]" />
      <motion.span
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        className="font-medium whitespace-pre text-[#c9d1d9]"
      >
        OS Compass
      </motion.span>
    </a>
  );
};
export const LogoIcon = () => {
  return (
    <a
      href="/"
      className="relative z-20 flex items-center space-x-2 py-1 text-sm font-normal text-[#c9d1d9]"
    >
      <div className="h-5 w-6 shrink-0 rounded-tl-lg rounded-tr-sm rounded-br-lg rounded-bl-sm bg-[#c9d1d9]" />
    </a>
  );
};


