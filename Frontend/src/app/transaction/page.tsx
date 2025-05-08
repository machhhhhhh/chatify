"use client";

import React, { ChangeEvent, FormEvent, useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import Image from "next/image";
import { GenerateClientJWT } from "@/libs/jwt";
import UserBadge from "@/components/UserBadge";

interface ITransactionFile {
  transaction_file_id: number;
  transaction_file_name: string;
  transaction_file_path: string;
  transaction_file_type: string;
}

interface ITransaction {
  transaction_id: number;
  transaction_number: string;
  transaction_description: string;
  creator: {
    account_first_name: string;
    account_last_name: string;
  };
  transaction_file: ITransactionFile[];
  created_at: string;
  transaction_reference: ITransaction[];
}

interface ITransactionSuccess {
  data: ITransaction[];
}

interface IRequestError {
  message: string;
}

export default function TransactionPage() {
  const router = useRouter();

  const [transactions, setTransactions] = useState<ITransaction[]>([]);
  const [error, setError] = useState<string>("");
  const [errorCreate, setErrorCreate] = useState<string>("");
  const [replyToId, setReplyToId] = useState<number | null>(null);
  const [replyText, setReplyText] = useState("");
  const [replyFiles, setReplyFiles] = useState<File[]>([]);
  const [sending, setSending] = useState(false);
  const [sendError, setSendError] = useState("");

  const [showCreate, setShowCreate] = useState(false);
  const [newDesc, setNewDesc] = useState("");
  const [newFiles, setNewFiles] = useState<File[]>([]);
  const [creating, setCreating] = useState(false);

  const rawName =
    typeof window !== "undefined"
      ? localStorage.getItem("display_name") || ""
      : "";
  const displayName = rawName || "User";

  const logout = () => {
    localStorage.removeItem("access_token");
    localStorage.removeItem("account_id");
    localStorage.removeItem("account_number");
    localStorage.removeItem("display_name");
    router.replace("/");
  };

  async function handleReplySubmit(e: FormEvent) {
    e.preventDefault();
    if (!replyToId) return;
    setSending(true);
    setSendError("");
    try {
      const jwt = await GenerateClientJWT({
        transaction_id: replyToId,
        transaction_description: replyText,
      });
      const fd = new FormData();
      fd.append("payload", jwt);
      replyFiles.forEach((f) => fd.append("files", f));

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
      // clear form
      setReplyText("");
      setReplyFiles([]);
      setReplyToId(null);
      // refresh thread
      loadTransactions();
    } catch (e) {
      setSendError((e as Error).message || "Network error");
    } finally {
      setSending(false);
    }
  }
  const loadTransactions = async () => {
    setError("");
    try {
      const payload = await GenerateClientJWT({
        is_show_comment: false,
      });
      const res = await fetch(
        `${process.env.NEXT_PUBLIC_API_URL}/manage-transaction/get-list-transaction`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${localStorage.getItem("access_token")}`,
          },
          body: JSON.stringify({ payload }),
        }
      );
      if (!res.ok) {
        const err: IRequestError = await res.json();
        setError(err.message || "Failed to load transactions");
        return;
      }
      const data: ITransactionSuccess = await res.json();
      setTransactions(data.data);
    } catch {
      setError("Network error, try again.");
    }
  };

  useEffect(() => {
    loadTransactions();
  }, []);

  const openCreate = () => setShowCreate(true);
  const closeCreate = () => {
    setShowCreate(false);
    setNewDesc("");
    setNewFiles([]);
  };

  const onFileChange = (e: ChangeEvent<HTMLInputElement>) => {
    const files = e.currentTarget.files;
    if (!files) return;
    setNewFiles((prev) => [...prev, ...Array.from(files)]);
  };

  const removeFile = (index: number) =>
    setNewFiles((files) => files.filter((_, i) => i !== index));

  const handleCreate = async (e: FormEvent) => {
    e.preventDefault();
    setCreating(true);
    setErrorCreate("");

    try {
      const payload = await GenerateClientJWT({
        transaction_description: newDesc,
        transaction_id: 0,
      });

      const formData = new FormData();
      formData.append("payload", payload);
      newFiles.forEach((f) => formData.append("files", f));

      const res = await fetch(
        `${process.env.NEXT_PUBLIC_API_URL}/manage-transaction/create-transaction`,
        {
          method: "POST",
          headers: {
            Authorization: `Bearer ${localStorage.getItem("access_token")}`,
          },
          body: formData,
        }
      );

      if (!res.ok) {
        const err: IRequestError = await res.json();
        setErrorCreate(err.message || "Create failed");
      } else {
        closeCreate();
        loadTransactions();
      }
    } catch {
      setErrorCreate("Network error, try again.");
    } finally {
      setCreating(false);
    }
  };

  const onReplyFileChange = (e: ChangeEvent<HTMLInputElement>) => {
    const fl = e.currentTarget.files;
    if (!fl) return;
    setReplyFiles((prev) => [...prev, ...Array.from(fl)]);
  };
  const removeReplyFile = (idx: number) =>
    setReplyFiles((f) => f.filter((_, i) => i !== idx));

  function Thread({ node, depth = 1 }: { node: ITransaction; depth?: number }) {
    const indent = depth * 16;
    return (
      <div
        className="bg-white bg-opacity-90 p-3 rounded-lg border border-slate-200 mt-3"
        style={{ marginLeft: indent }}
      >
        <div className="flex items-center space-x-2 mb-2">
          <div className="h-6 w-6 rounded-full bg-blue-400 flex items-center justify-center text-white text-sm font-bold">
            {node.creator.account_first_name.charAt(0).toUpperCase()}
          </div>
          <div>
            <p className="text-slate-900 text-sm font-medium">
              {node.creator.account_first_name} {node.creator.account_last_name}
            </p>
            <p className="text-xs text-slate-500">
              {new Date(node.created_at).toLocaleString()}
            </p>
          </div>
        </div>
        <p className="text-slate-700 text-sm mb-2">
          {node.transaction_description}
        </p>
        {node.transaction_file.length > 0 && (
          <div className="flex gap-2 flex-wrap">
            {node.transaction_file.map((f) => {
              const url =
                process.env.NEXT_PUBLIC_API_URL!.replace(/\/$/, "") +
                f.transaction_file_path;
              if (f.transaction_file_type.startsWith("image/")) {
                return (
                  <div
                    key={f.transaction_file_id}
                    className="h-16 w-16 relative overflow-hidden rounded-lg"
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
                  className="px-2 py-1 bg-slate-100 rounded text-xs"
                >
                  {f.transaction_file_name}
                </a>
              );
            })}
          </div>
        )}

        {/* Reply button */}
        <div className="mt-2 flex justify-end">
          <button
            onClick={() => setReplyToId(node.transaction_id)}
            className="text-sm text-sky-600 hover:underline"
          >
            Reply
          </button>
        </div>

        {/* Inline reply form */}
        {replyToId === node.transaction_id && (
          <form
            onSubmit={handleReplySubmit}
            className="mt-4 space-y-2 bg-slate-50 p-4 rounded-lg"
          >
            <textarea
              value={replyText}
              onChange={(e) => setReplyText(e.target.value)}
              rows={2}
              placeholder="Write a reply…"
              className="w-full border border-slate-300 rounded-lg p-2"
              required
            />
            <input type="file" multiple onChange={onReplyFileChange} />
            {replyFiles.length > 0 && (
              <ul className="space-y-1">
                {replyFiles.map((f, i) => (
                  <li key={i} className="flex justify-between text-sm">
                    <span className="truncate">{f.name}</span>
                    <button
                      type="button"
                      onClick={() => removeReplyFile(i)}
                      className="text-red-600"
                    >
                      ×
                    </button>
                  </li>
                ))}
              </ul>
            )}
            {sendError && <p className="text-red-600 text-sm">{sendError}</p>}
            <button
              type="submit"
              disabled={sending}
              className="px-4 py-2 bg-sky-500 hover:bg-sky-600 text-white rounded-lg"
            >
              {sending ? "Posting…" : "Submit Reply"}
            </button>
          </form>
        )}

        {node.transaction_reference?.map((r) => (
          <Thread key={r.transaction_id} node={r} depth={depth + 1} />
        ))}
      </div>
    );
  }

  return (
    <div
      className="min-h-screen flex flex-col
                    bg-gradient-to-tr from-sky-200 via-sky-100 to-sky-50
                    text-slate-900"
    >
      {/* Top bar */}
      <header className="flex items-center justify-between px-6 py-4 bg-white bg-opacity-80 border-b border-slate-200">
        <div className="flex items-center space-x-4">
          <h1 className="text-2xl font-bold">Transactions</h1>
          <button
            onClick={openCreate}
            className="px-4 py-2 bg-green-500 hover:bg-green-600 text-white rounded-md transition"
          >
            New Transaction
          </button>
        </div>
        <div className="flex items-center space-x-4">
          <UserBadge displayName={displayName} />
          <button
            onClick={logout}
            className="px-4 py-2 bg-red-500 hover:bg-red-600 text-white rounded-md transition"
          >
            Logout
          </button>
        </div>
      </header>

      {/* Error banner */}
      {error && (
        <div className="p-4 text-red-700 bg-red-100 m-6 rounded">{error}</div>
      )}

      {/* Transaction list */}
      <main className="flex-1 overflow-auto p-6">
        <ul className="space-y-4">
          {transactions.map((tx) => (
            <li
              key={tx.transaction_id}
              className="flex flex-col md:flex-row bg-white bg-opacity-80 p-4 rounded-lg border border-slate-200"
            >
              <div className="flex-1">
                <div className="flex items-center space-x-3">
                  <div className="h-10 w-10 rounded-full bg-blue-400 flex items-center justify-center text-white font-bold">
                    {tx.creator.account_first_name[0].toUpperCase()}
                  </div>
                  <div>
                    <p className="text-sm text-sky-600 font-semibold">
                      #{tx.transaction_number}
                    </p>
                    <p className="font-medium text-slate-900">
                      {tx.creator.account_first_name}{" "}
                      {tx.creator.account_last_name}
                    </p>
                    <p className="text-sm text-slate-500">
                      {new Date(tx.created_at).toLocaleString()}
                    </p>
                  </div>
                </div>

                <p className="mt-3 text-slate-700 whitespace-pre-wrap">
                  {tx.transaction_description}
                </p>

                {tx.transaction_file.length > 0 && (
                  <div className="mt-4 flex flex-wrap gap-4">
                    {tx.transaction_file.map((file) => {
                      const url =
                        process.env.NEXT_PUBLIC_API_URL!.replace(/\/$/, "") +
                        file.transaction_file_path;
                      const isImage =
                        file.transaction_file_type.startsWith("image/");

                      return isImage ? (
                        <div
                          key={file.transaction_file_id}
                          className="relative h-32 w-32 rounded-lg overflow-hidden shadow-md border border-slate-300"
                        >
                          <Image
                            src={url}
                            alt={file.transaction_file_name}
                            fill
                            style={{ objectFit: "cover" }}
                            unoptimized
                          />
                        </div>
                      ) : (
                        <a
                          key={file.transaction_file_id}
                          href={url}
                          target="_blank"
                          rel="noopener noreferrer"
                          className="flex items-center justify-center h-32 w-32 rounded-lg border border-slate-300 bg-slate-100 p-2 text-center text-sm hover:bg-slate-200 transition break-words"
                        >
                          {file.transaction_file_name}
                        </a>
                      );
                    })}
                  </div>
                )}

                {/* ← Inline comments */}
                {tx.transaction_reference.length > 0 && (
                  <div className="mt-4">
                    <h4 className="font-semibold text-slate-900 mb-2">
                      Comments
                    </h4>
                    {tx.transaction_reference.map((reply) => (
                      <Thread
                        key={reply.transaction_id}
                        node={reply}
                        depth={1}
                      />
                    ))}
                  </div>
                )}
              </div>

              <button
                onClick={() => router.push(`/transaction/${tx.transaction_id}`)}
                className="mt-4 md:mt-0 md:ml-6 px-4 py-2 bg-sky-400 hover:bg-sky-500 text-white rounded transition self-start"
              >
                Comment
              </button>
            </li>
          ))}
        </ul>
      </main>

      {/* Create Transaction Modal */}
      {showCreate && (
        <div className="fixed inset-0 flex items-center justify-center bg-black/50 p-4">
          <form
            onSubmit={handleCreate}
            className="bg-white rounded-xl shadow-xl w-[70vw] h-[70vh] overflow-auto p-6 space-y-4"
          >
            <h2 className="text-xl font-bold text-slate-900">
              Create Transaction
            </h2>
            {errorCreate && (
              <div className="p-4 text-red-700 bg-red-100 m-6 rounded">
                {errorCreate}
              </div>
            )}
            <label className="block">
              <span className="text-slate-700">Description</span>
              <textarea
                value={newDesc}
                onChange={(e) => setNewDesc(e.target.value)}
                rows={4}
                required
                className="mt-1 w-full px-3 py-2 border border-slate-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-sky-400 resize-vertical"
              />
            </label>
            <label className="block">
              <span className="text-slate-700">Files</span>
              <input
                type="file"
                multiple
                onChange={onFileChange}
                className="mt-1 w-full"
              />
            </label>
            {newFiles.length > 0 && (
              <ul className="space-y-2">
                {newFiles.map((file, idx) => (
                  <li
                    key={idx}
                    className="flex items-center justify-between bg-slate-100 p-2 rounded"
                  >
                    <span className="truncate">{file.name}</span>
                    <button
                      type="button"
                      onClick={() => removeFile(idx)}
                      className="text-red-600 hover:text-red-800"
                    >
                      Remove
                    </button>
                  </li>
                ))}
              </ul>
            )}
            <div className="flex justify-end space-x-2">
              <button
                type="button"
                onClick={closeCreate}
                className="px-4 py-2 bg-slate-300 hover:bg-slate-400 text-slate-800 rounded transition"
                disabled={creating}
              >
                Cancel
              </button>
              <button
                type="submit"
                disabled={creating}
                className="px-4 py-2 bg-green-500 hover:bg-green-600 text-white rounded transition disabled:opacity-50"
              >
                {creating ? "Creating…" : "Create"}
              </button>
            </div>
          </form>
        </div>
      )}
    </div>
  );
}
