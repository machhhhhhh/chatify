"use client";

import { ReactNode, useEffect, useState } from "react";
import { useRouter } from "next/navigation";

export default function TransactionLayout({
  children,
}: {
  children: ReactNode;
}) {
  const router = useRouter();
  const [authorized, setAuthorized] = useState(false);

  useEffect(() => {
    const account_token = localStorage.getItem("access_token");
    const account_id = localStorage.getItem("account_id");
    if (!account_token || !account_id) {
      router.replace("/");
    } else {
      setAuthorized(true);
    }
  }, [router]);

  if (!authorized) {
    return <p className="p-8 text-center">Redirecting to login…</p>;
  }

  return <>{children}</>;
}
