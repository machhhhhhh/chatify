const encoder = new TextEncoder();

function base64UrlEncode(bytes: Uint8Array): string {
  let str = "";
  for (let i = 0; i < bytes.byteLength; i++) {
    str += String.fromCharCode(bytes[i]);
  }
  return btoa(str).replace(/\+/g, "-").replace(/\//g, "_").replace(/=+$/, "");
}

async function importHmacKey(secret: string): Promise<CryptoKey> {
  return crypto.subtle.importKey(
    "raw",
    encoder.encode(secret),
    { name: "HMAC", hash: "SHA-256" },
    false,
    ["sign"]
  );
}

export async function GenerateClientJWT(
  payload: Record<string, unknown>,
  expiresInSec = 60 * 60 * 1 // 1 hour
): Promise<string> {
  if (typeof window === "undefined") {
    throw new Error("generateClientJWT() only in browser");
  }
  const secret = localStorage.getItem("access_token");
  if (!secret) {
    throw new Error("No access_token in localStorage to use as secret");
  }

  const header = { alg: "HS256", typ: "JWT" };
  const iat = Math.floor(Date.now() / 1000);
  const exp = iat + expiresInSec;
  const fullPayload = { ...payload, iat, exp };

  const encodedHeader = base64UrlEncode(encoder.encode(JSON.stringify(header)));
  const encodedPayload = base64UrlEncode(
    encoder.encode(JSON.stringify(fullPayload))
  );
  const data = `${encodedHeader}.${encodedPayload}`;

  const key = await importHmacKey(secret);
  const sigBuffer = await crypto.subtle.sign("HMAC", key, encoder.encode(data));
  const signature = base64UrlEncode(new Uint8Array(sigBuffer));

  return `${data}.${signature}`;
}
