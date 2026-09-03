import type { RequestHandler } from '@sveltejs/kit';

const BACKEND_URL = 'http://34.134.26.219:8000';

export const fallback: RequestHandler = async ({ params, request, url }) => {
  const path = params.path ?? '';
  const targetUrl = `${BACKEND_URL}/api/${path}${url.search}`;

  const headers = new Headers(request.headers);
  headers.delete('host');

  try {
    const response = await fetch(targetUrl, {
      method: request.method,
      headers,
      body: ['GET', 'HEAD'].includes(request.method) ? undefined : await request.blob(),
      duplex: 'half'
    } as RequestInit);

    return response;
  } catch (err) {
    return new Response(JSON.stringify({ error: 'Backend serverga ulanib bo‘lmadi' }), {
      status: 502,
      headers: { 'Content-Type': 'application/json' }
    });
  }
};