"use client";

import React, { useState, FormEvent, JSX } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";

interface RegisterRequest {
  account_identify_number: string;
  account_username: string;
  account_password: string;
  account_email: string;
  account_phone_number: string;
  account_first_name: string;
  account_last_name: string;
}
interface ApiResponse {
  message: string;
}

export default function RegisterPage(): JSX.Element {
  const router = useRouter();
  const [form, setForm] = useState<RegisterRequest>({
    account_identify_number: "",
    account_username: "",
    account_password: "",
    account_email: "",
    account_phone_number: "",
    account_first_name: "",
    account_last_name: "",
  });
  const [error, setError] = useState<string>("");
  const [info, setInfo] = useState<string>("");
  const [loading, setLoading] = useState<boolean>(false);

  const update = (k: keyof RegisterRequest, v: string) =>
    setForm((f) => ({ ...f, [k]: v }));

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setError("");
    setInfo("");
    setLoading(true);

    try {
      const res = await fetch(
        `${process.env.NEXT_PUBLIC_API_URL}/manage-account/create-account`,
        {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify(form),
        }
      );
      const data: ApiResponse = await res.json();
      if (!res.ok) {
        setError(data.message || "Registration failed");
        return;
      }
      setInfo(data.message || "Account created!");
      setTimeout(() => {
        router.push("/");
        router.refresh();
      }, 1500);
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
        className="w-full max-w-lg bg-white bg-opacity-80 border border-slate-200 rounded-2xl shadow-xl p-8 space-y-5"
      >
        <h2 className="text-3xl font-bold text-slate-900 text-center">
          Register
        </h2>

        {error && (
          <p className="text-red-700 bg-red-100 p-2 rounded">{error}</p>
        )}
        {info && (
          <p className="text-green-700 bg-green-100 p-2 rounded">{info}</p>
        )}

        <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
          {/** ID Number **/}
          <label className="block">
            <span className="text-slate-700">ID Number</span>
            <input
              type="text"
              required
              value={form.account_identify_number}
              onChange={(e) =>
                update("account_identify_number", e.target.value)
              }
              className="mt-1 w-full px-4 py-2 rounded-lg border border-slate-300 bg-white focus:outline-none focus:ring-2 focus:ring-sky-400 text-slate-900"
            />
          </label>
          {/** Username **/}
          <label className="block">
            <span className="text-slate-700">Username</span>
            <input
              type="text"
              required
              value={form.account_username}
              onChange={(e) => update("account_username", e.target.value)}
              className="mt-1 w-full px-4 py-2 rounded-lg border border-slate-300 bg-white focus:outline-none focus:ring-2 focus:ring-sky-400 text-slate-900"
            />
          </label>
          {/** Password **/}
          <label className="block">
            <span className="text-slate-700">Password</span>
            <input
              type="password"
              required
              value={form.account_password}
              onChange={(e) => update("account_password", e.target.value)}
              className="mt-1 w-full px-4 py-2 rounded-lg border border-slate-300 bg-white focus:outline-none focus:ring-2 focus:ring-sky-400 text-slate-900"
            />
          </label>
          {/** Email **/}
          <label className="block">
            <span className="text-slate-700">Email</span>
            <input
              type="email"
              required
              value={form.account_email}
              onChange={(e) => update("account_email", e.target.value)}
              className="mt-1 w-full px-4 py-2 rounded-lg border border-slate-300 bg-white focus:outline-none focus:ring-2 focus:ring-sky-400 text-slate-900"
            />
          </label>
          {/** Phone **/}
          <label className="block">
            <span className="text-slate-700">Phone</span>
            <input
              type="tel"
              required
              value={form.account_phone_number}
              onChange={(e) => update("account_phone_number", e.target.value)}
              className="mt-1 w-full px-4 py-2 rounded-lg border border-slate-300 bg-white focus:outline-none focus:ring-2 focus:ring-sky-400 text-slate-900"
            />
          </label>
          {/** First Name **/}
          <label className="block">
            <span className="text-slate-700">First Name</span>
            <input
              type="text"
              required
              value={form.account_first_name}
              onChange={(e) => update("account_first_name", e.target.value)}
              className="mt-1 w-full px-4 py-2 rounded-lg border border-slate-300 bg-white focus:outline-none focus:ring-2 focus:ring-sky-400 text-slate-900"
            />
          </label>
          {/** Last Name **/}
          <label className="block">
            <span className="text-slate-700">Last Name</span>
            <input
              type="text"
              required
              value={form.account_last_name}
              onChange={(e) => update("account_last_name", e.target.value)}
              className="mt-1 w-full px-4 py-2 rounded-lg border border-slate-300 bg-white focus:outline-none focus:ring-2 focus:ring-sky-400 text-slate-900"
            />
          </label>
        </div>

        <button
          type="submit"
          disabled={loading}
          className="w-full py-3 rounded-lg bg-sky-500 text-white font-semibold text-lg hover:bg-sky-600 transition disabled:opacity-50 disabled:cursor-not-allowed"
        >
          {loading ? "Creating…" : "Create Account"}
        </button>

        <p className="text-center text-sm text-slate-700">
          Already have an account?{" "}
          <Link href="/">
            <span className="font-medium text-sky-600 hover:underline">
              Back to Login
            </span>
          </Link>
        </p>
      </form>
    </div>
  );
}
