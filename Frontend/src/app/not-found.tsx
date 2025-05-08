"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";

export default function NotFound() {
  const router = useRouter();
  useEffect(() => {
    router.replace("/");
  }, [router]);

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-tr from-sky-200 via-sky-100 to-sky-50 p-6">
      <p className="text-slate-900 text-lg">
        Oops—page not found. Redirecting to login…
      </p>
    </div>
  );
}
