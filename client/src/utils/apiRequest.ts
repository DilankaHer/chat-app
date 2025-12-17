import { useDialog } from './appDialog';

export interface ApiRequest {
  url: string;
  method: string;
  body?: any;
  fn?: () => void;
  dialogType?: 'dialog' | 'toast';
}

export interface ApiResponse {
  data: any;
  message: string;
  error: string;
}

export function useApi() {
  const { showDialog } = useDialog();
  const baseURL = import.meta.env.VITE_BASE_URL;
  async function apiRequest<T>(req: ApiRequest, ifNull?: T): Promise<T> {
    let isUnexpectedError = true;
    try {
      const response = await fetch(baseURL + req.url, {
        method: req.method,
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include',
        body: JSON.stringify(req.body),
      });
      const data: ApiResponse = await response.json();
      if (response.status !== 200) {
        if (data.error !== 'missing auth token' && req.dialogType) {
          showDialog(req.dialogType, data.error, true, req.fn);
        }
        isUnexpectedError = false;
        throw new Error(data.error);
      }
      if (req.dialogType) {
        showDialog(req.dialogType, data.message, false);
      }
      return data.data == null && ifNull !== undefined ? ifNull as T: data.data as T;
    } catch (error) {
      if (isUnexpectedError) {
        showDialog('toast', 'Something went wrong', true, req.fn);
      }
      throw error;
    }
  }
  return { apiRequest };
}
