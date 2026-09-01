import { API_BASE_URL } from "./client";
import type { ImportResponse } from "$lib/types";

// xlsx faylni import qilish (multipart/form-data, faqat admin).
export async function importXlsx(file: File): Promise<ImportResponse> {
    const token = localStorage.getItem("token");
    const formData = new FormData();
    formData.append("file", file);

    const res = await fetch(`${API_BASE_URL}/api/import/xlsx`, {
        method: "POST",
        headers: token ? { Authorization: `Bearer ${token}` } : {},
        body: formData,
    });

    if (!res.ok) {
        let message = "import_failed";
        try {
            const data = await res.json();
            message = data?.error ?? message;
        } catch {
            // ignore parse errors
        }
        throw new Error(message);
    }

    return (await res.json()) as ImportResponse;
}
