"use client";

import React from "react";

export interface UserBadgeProps {
  displayName: string;
}

export default function UserBadge({ displayName }: UserBadgeProps) {
  return (
    <div className="flex items-center space-x-3">
      {/* Avatar */}
      <div
        className="w-8 h-8 rounded-full bg-sky-500 text-white
                   flex items-center justify-center font-semibold"
      >
        {displayName.charAt(0).toUpperCase() || "?"}
      </div>

      {/* Name */}
      <span className="font-medium text-slate-900">{displayName}</span>
    </div>
  );
}
