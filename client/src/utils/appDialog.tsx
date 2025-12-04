import { createContext, useContext, useState } from "react";


type DialogContextType = {
  showDialog: (status: number, message: string, fn?: () => void) => void;
};

const DialogContext = createContext<DialogContextType>({
    showDialog: () => {}
});

export function DialogProvider({ children }: { children: React.ReactNode }) {
  const [dialog, setDialog] = useState({
    open: false,
    status: 0,
    message: "",
  });
   const [callback, setCallback] = useState<(() => void) | null>(null);
    function showDialog(status: number, message: string, fn?: () => void) {
        setCallback(() => fn || null);
        setDialog({
            open: true,
            status,
            message,
        })
    }

    function handleClose() {
      setDialog(prev => ({ ...prev, open: false }));
      if (callback) callback();
    }
    return (
        <DialogContext.Provider value={{ showDialog: showDialog }}>
            {children}
            <ApiDialog open={dialog.open} onClose={handleClose} status={dialog.status} message={dialog.message} />
        </DialogContext.Provider>
    )
}


export default function ApiDialog({ open, onClose, status, message }: { open: boolean, onClose: () => void, status: number, message: string }) {
  if (!open) return null;

  return (
    <div className="fixed inset-0 flex items-center justify-center bg-black/50 z-50">
      <div className="bg-white dark:bg-gray-800 rounded-lg shadow-lg p-6 w-[90%] max-w-sm">

        <h2 className="text-xl font-semibold mb-2 text-gray-900 dark:text-white">
          Server Response
        </h2>

        <p className="text-gray-700 dark:text-gray-300">
          <strong>Status:</strong> {status}
        </p>
        <p className="text-gray-700 dark:text-gray-300">
          <strong>Message:</strong> {message}
        </p>

        <button
          onClick={onClose}
          className="mt-4 px-4 py-2 bg-indigo-600 hover:bg-indigo-700 text-white rounded-md"
        >
          Close
        </button>
      </div>
    </div>
  );
}

export function useDialog() {
    return useContext(DialogContext);
}
