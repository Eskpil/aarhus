export default async function fetcher<JSON = any>(
  input: RequestInfo,
  _init?: RequestInit
): Promise<JSON> {
  const res = await fetch(input, { credentials: "include" });
  return res.json();
}
