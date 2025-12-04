import { useDialog } from "./appDialog";

type dataType = {
    data: any;
}

export interface ApiRequest {
    url: string;
    method: string;
    body?: any;
}

export function useApi() {
    const { showDialog } = useDialog();
    async function apiRequest(req: ApiRequest, fn?: () => void): Promise<dataType> {
    const response = await fetch(req.url, {
        method: req.method,
        headers: {
            "Content-Type": "application/json",
        },
        credentials: "include",
        body: JSON.stringify(req.body),
    });
    const data = await response.json().then((data) => data?.data).catch((error) => error);
    if (response.status !== 200) {
        showDialog(response.status, data, fn);
        throw new Error(data);
    }
    return { data };
}
    return { apiRequest };
}

