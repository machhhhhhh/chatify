"use client";

import React, { useState, FormEvent, JSX, useEffect } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";

interface LoginSuccess {
  data: {
    access_token: string;
    account_id: number;
    account_number: string;
    display_name: string;
  };
}
interface LoginError {
  message: string;
}
interface LoginRequest {
  account_username: string;
  account_password: string;
}

export default function LoginPage(): JSX.Element {
  const router = useRouter();
  const [username, setUsername] = useState<string>("");
  const [password, setPassword] = useState<string>("");
  const [error, setError] = useState<string>("");
  const [loading, setLoading] = useState<boolean>(false);

  useEffect(() => {
    const access_token = localStorage.getItem("access_token");
    const account_id = localStorage.getItem("account_id");
    if (access_token && account_id) router.replace("/transaction");
  }, [router]);

  const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setError("");
    setLoading(true);

    try {
      const res = await fetch(
        `${process.env.NEXT_PUBLIC_API_URL}/authentication/system-login`,
        {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({
            account_username: username,
            account_password: password,
          } as LoginRequest),
        }
      );
      if (!res.ok) {
        const err: LoginError = await res.json();
        setError(err.message || "Login failed");
        return;
      }
      const data: LoginSuccess = await res.json();
      localStorage.setItem("access_token", String(data.data.access_token));
      localStorage.setItem("account_id", String(data.data.account_id));
      localStorage.setItem("account_number", String(data.data.account_number));
      localStorage.setItem("display_name", String(data.data.display_name));
      router.push("/transaction");
      router.refresh();
    } catch {
      setError("Network error, try again.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-tr from-sky-200 via-sky-100 to-sky-50 p-6">
      <form
        onSubmit={handleSubmit}
        className="w-full max-w-md bg-white bg-opacity-80 border border-slate-200 rounded-2xl shadow-xl p-8 space-y-6"
      >
        <h2 className="text-3xl font-bold text-slate-900 text-center">Login</h2>

        {error && (
          <p className="text-sm text-red-700 bg-red-100 p-2 rounded">{error}</p>
        )}

        <div className="space-y-4">
          <label className="block">
            <span className="text-slate-700">Username</span>
            <input
              type="text"
              className="mt-1 w-full px-4 py-2 rounded-lg border border-slate-300 bg-white placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-sky-400 text-slate-900"
              placeholder="your_username"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              required
            />
          </label>
          <label className="block">
            <span className="text-slate-700">Password</span>
            <input
              type="password"
              className="mt-1 w-full px-4 py-2 rounded-lg border border-slate-300 bg-white placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-sky-400 text-slate-900"
              placeholder="••••••••"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
            />
          </label>
        </div>

        <button
          type="submit"
          disabled={loading}
          className="w-full py-3 rounded-lg bg-sky-500 text-white font-semibold text-lg hover:bg-sky-600 transition disabled:opacity-50 disabled:cursor-not-allowed"
        >
          {loading ? "Signing In…" : "Sign In"}
        </button>

        <p className="text-center text-sm text-slate-700">
          Don’t have an account?{" "}
          <Link href="/register">
            <span className="font-medium text-sky-600 hover:underline">
              Register
            </span>
          </Link>
        </p>
      </form>
    </div>
  );
}
