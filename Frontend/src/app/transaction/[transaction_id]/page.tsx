"use client";

import React, { useEffect, useState, FormEvent, ChangeEvent } from "react";
import { useRouter, useParams } from "next/navigation";
import Image from "next/image";
import { GenerateClientJWT } from "@/libs/jwt";

interface IFile {
  transaction_file_id: number;
  transaction_file_name: string;
  transaction_file_path: string;
  transaction_file_type: string;
}

interface ITransaction {
  transaction_id: number;
  transaction_number: string;
  transaction_description: string;
  creator: { account_first_name: string; account_last_name: string };
  transaction_file: IFile[];
  transaction_reference?: ITransaction[];
  created_at: string;
}

interface IResponse<T> {
  data: T;
  message?: string;
}

export default function TransactionDetailPage() {
  const router = useRouter();
  const { transaction_id } = useParams<{ transaction_id: string }>();
  const id = parseInt(transaction_id || "", 10);

  const [transaction, setTransaction] = useState<ITransaction | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  const [reply, setReply] = useState("");
  const [files, setFiles] = useState<File[]>([]);
  const [sending, setSending] = useState(false);
  const [sendError, setSendError] = useState("");

  async function fetchDetail() {
    setLoading(true);
    setError("");
    if (isNaN(id)) {
      setError("Invalid transaction ID");
      setLoading(false);
      return;
    }
    try {
      const payload = await GenerateClientJWT({ transaction_id: id });
      const res = await fetch(
        `${process.env.NEXT_PUBLIC_API_URL}/manage-transaction/get-information-transaction`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${localStorage.getItem("access_token")}`,
          },
          body: JSON.stringify({ payload }),
        }
      );
      if (!res.ok) throw new Error("Fetch failed");
      const json: IResponse<ITransaction> = await res.json();
      setTransaction(json.data);
    } catch (e) {
      setError((e as Error).message || "Network error");
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    fetchDetail();
  }, [transaction_id]);

  const onFileChange = (e: ChangeEvent<HTMLInputElement>) => {
    const fl = e.currentTarget.files;
    if (!fl) return;
    setFiles((f) => [...f, ...Array.from(fl)]);
  };
  const removeFile = (idx: number) =>
    setFiles((f) => f.filter((_, i) => i !== idx));

  async function handleReply(e: FormEvent) {
    e.preventDefault();
    setSending(true);
    setSendError("");
    try {
      const jwt = await GenerateClientJWT({
        transaction_id: id,
        transaction_description: reply,
      });
      const fd = new FormData();
      fd.append("payload", jwt);
      files.forEach((f) => fd.append("files", f));

      const res = await fetch(
        `${process.env.NEXT_PUBLIC_API_URL}/manage-transaction/create-transaction`,
        {
          method: "POST",
          headers: {
            Authorization: `Bearer ${localStorage.getItem("access_token")}`,
          },
          body: fd,
        }
      );
      if (!res.ok) throw new Error("Send failed");
      setReply("");
      setFiles([]);
      fetchDetail();
    } catch (e) {
      setSendError((e as Error).message || "Network error");
    } finally {
      setSending(false);
    }
  }

  function Thread({ node, depth = 0 }: { node: ITransaction; depth?: number }) {
    const indent = depth * 24;
    return (
      <div
        className="bg-white bg-opacity-80 p-4 rounded-lg border border-slate-200 mt-4"
        style={{ marginLeft: indent }}
      >
        <div className="flex items-center justify-between">
          <div className="flex items-center space-x-3">
            <div className="h-8 w-8 rounded-full bg-blue-400 flex items-center justify-center text-white font-bold">
              {node.creator.account_first_name.charAt(0).toUpperCase()}
            </div>
            <div>
              <p className="text-slate-900 font-medium">
                {node.creator.account_first_name}{" "}
                {node.creator.account_last_name}
              </p>
              <p className="text-sm text-sky-600 font-semibold">
                #{node.transaction_number}
              </p>
              <p className="text-sm text-slate-500">
                {new Date(node.created_at).toLocaleString()}
              </p>
            </div>
          </div>
        </div>

        <p className="mt-2 text-slate-700 whitespace-pre-wrap">
          {node.transaction_description}
        </p>

        {node.transaction_file?.length > 0 && (
          <div className="mt-2 flex flex-wrap gap-3">
            {node.transaction_file.map((f) => {
              const url =
                process.env.NEXT_PUBLIC_API_URL!.replace(/\/$/, "") +
                f.transaction_file_path;
              if (f.transaction_file_type.startsWith("image/")) {
                return (
                  <div
                    key={f.transaction_file_id}
                    className="h-24 w-24 relative overflow-hidden rounded-lg border"
                  >
                    <Image
                      src={url}
                      alt={f.transaction_file_name}
                      fill
                      style={{ objectFit: "cover" }}
                      unoptimized
                    />
                  </div>
                );
              }
              return (
                <a
                  key={f.transaction_file_id}
                  href={url}
                  target="_blank"
                  rel="noreferrer"
                  className="px-2 py-1 bg-slate-100 rounded hover:bg-slate-200 text-sm"
                >
                  {f.transaction_file_name}
                </a>
              );
            })}
          </div>
        )}

        {node.transaction_reference?.map((r) => (
          <Thread key={r.transaction_id} node={r} depth={depth + 1} />
        ))}
      </div>
    );
  }

  if (loading) return <p className="p-6">Loading…</p>;
  if (error) return <p className="p-6 text-red-600">{error}</p>;
  if (!transaction) return <p className="p-6">Transaction not found</p>;
  console.log(transaction);
  return (
    <div className="min-h-screen flex flex-col bg-gradient-to-tr from-sky-200 via-sky-100 to-sky-50 text-slate-900">
      {/* ← Back link */}
      <div className="p-6">
        <button
          onClick={() => router.back()}
          className="text-sky-600 hover:underline"
        >
          ← Back to list
        </button>
      </div>

      {/* Mother card */}
      <div className="max-w-4xl mx-auto border-4 border-slate-900 rounded-xl overflow-hidden">
        {/* Green header */}
        <div className="bg-green-600 text-white text-center py-2 font-semibold">
          Transaction #{transaction.transaction_number}
        </div>

        {/* Body */}
        <div className="bg-white p-6 space-y-6">
          {/* Avatar + name + timestamp */}
          <div className="flex items-center space-x-4">
            <div className="h-12 w-12 rounded-full bg-blue-400 flex items-center justify-center text-white text-xl font-bold">
              {transaction.creator.account_first_name.charAt(0).toUpperCase()}
            </div>
            <div>
              <p className="text-lg font-medium text-slate-900">
                {transaction.creator.account_first_name}{" "}
                {transaction.creator.account_last_name}
              </p>
              <p className="text-sm text-slate-500">
                {new Date(transaction.created_at).toLocaleString()}
              </p>
            </div>
          </div>

          {/* Main attachments grid */}
          {transaction.transaction_file.length > 0 && (
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              {transaction.transaction_file.map((f) => {
                const url =
                  process.env.NEXT_PUBLIC_API_URL!.replace(/\/$/, "") +
                  f.transaction_file_path;
                if (f.transaction_file_type.startsWith("image/")) {
                  return (
                    <div
                      key={f.transaction_file_id}
                      className="w-full h-64 relative overflow-hidden rounded-lg border"
                    >
                      <Image
                        src={url}
                        alt={f.transaction_file_name}
                        fill
                        style={{ objectFit: "cover" }}
                        unoptimized
                      />
                    </div>
                  );
                }
                return (
                  <a
                    key={f.transaction_file_id}
                    href={url}
                    target="_blank"
                    rel="noreferrer"
                    className="block px-4 py-2 bg-slate-100 rounded hover:bg-slate-200 text-center"
                  >
                    {f.transaction_file_name}
                  </a>
                );
              })}
            </div>
          )}

          {/* Description */}
          <p className="text-slate-700 whitespace-pre-wrap">
            {transaction.transaction_description}
          </p>

          {/* Replies */}
          {transaction.transaction_reference &&
            transaction.transaction_reference.length > 0 && (
              <>
                <h3 className="font-semibold text-slate-900 mb-2">Replies</h3>
                {transaction.transaction_reference.map((r) => (
                  <Thread key={r.transaction_id} node={r} />
                ))}
              </>
            )}

          {/* Reply form */}
          <form onSubmit={handleReply} className="space-y-4">
            <textarea
              value={reply}
              onChange={(e) => setReply(e.target.value)}
              rows={3}
              placeholder="Write a reply…"
              className="w-full border border-slate-300 rounded-lg p-3 focus:ring-2 focus:ring-sky-400"
              required
            />
            <input
              type="file"
              multiple
              onChange={onFileChange}
              className="block"
            />
            {files.length > 0 && (
              <ul className="space-y-1">
                {files.map((f, i) => (
                  <li key={i} className="flex justify-between text-sm">
                    <span className="truncate">{f.name}</span>
                    <button
                      type="button"
                      onClick={() => removeFile(i)}
                      className="text-red-600"
                    >
                      ×
                    </button>
                  </li>
                ))}
              </ul>
            )}
            {sendError && <p className="text-red-600">{sendError}</p>}
            <button
              type="submit"
              disabled={sending}
              className="px-6 py-2 bg-sky-500 hover:bg-sky-600 text-white rounded-lg transition disabled:opacity-50"
            >
              {sending ? "Posting…" : "Comment"}
            </button>
          </form>
        </div>
      </div>
    </div>
  );
}
