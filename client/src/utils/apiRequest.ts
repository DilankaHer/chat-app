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
  async function apiRequest<T>(req: ApiRequest): Promise<T> {
    try {
      const response = await fetch(req.url, {
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
          showDialog(req.dialogType, data.error, req.fn);
        }
        throw new Error(data.error);
      }
      return data.data as T;
    } catch (error) {
      showDialog('toast', 'Something went wrong', req.fn);
      throw error;
    }
  }
  return { apiRequest };
}
