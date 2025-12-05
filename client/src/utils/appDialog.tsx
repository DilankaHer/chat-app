import { faExclamationCircle } from '@fortawesome/free-solid-svg-icons';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { createContext, useContext, useState } from 'react';

type DialogContextType = {
  showDialog: (dialogType: string, message: string, fn?: () => void) => void;
};

const DialogContext = createContext<DialogContextType>({
  showDialog: () => {},
});

export function DialogProvider({ children }: { children: React.ReactNode }) {
  const [dialog, setDialog] = useState({
    open: false,
    dialogType: 'dialog',
    message: '',
  });
  const [callback, setCallback] = useState<(() => void) | null>(null);
  function showDialog(dialogType: string, message: string, fn?: () => void) {
    setCallback(() => fn || null);
    setDialog({
      open: true,
      dialogType,
      message,
    });
  }

  function handleClose() {
    setDialog((prev) => ({ ...prev, open: false }));

    if (callback) callback();
  }
  return (
    <DialogContext.Provider value={{ showDialog: showDialog }}>
      {children}
      {dialog.dialogType === 'dialog' && (
        <ApiDialog
          open={dialog.open}
          onClose={handleClose}
          message={dialog.message}
        />
      )}
      {dialog.dialogType === 'toast' && (
        <ApiToast
          open={dialog.open}
          onClose={handleClose}
          message={dialog.message}
        />
      )}
    </DialogContext.Provider>
  );
}

function ApiDialog({
  open,
  onClose,
  message,
}: {
  open: boolean;
  onClose: () => void;
  message: string;
}) {
  if (!open) return null;

  return (
    <div className="fixed inset-0 flex items-center justify-center z-50 pointer-events-none">
      <div className="bg-white dark:bg-gray-800 rounded-lg shadow-lg p-6 w-[90%] max-w-sm">
        <div className="text-red-500 mb-3 text-center text-2xl">
          <FontAwesomeIcon
            icon={faExclamationCircle}
            fill="currentColor"
            size="xl"
          />
        </div>

        <p className="text-gray-700 dark:text-gray-300 text-center font-bold text-xl">
          {message}
        </p>

        <div className="flex justify-center mt-4">
          <button
            onClick={onClose}
            className="mt-4 px-4 py-2 bg-indigo-600 hover:bg-indigo-700 text-white rounded-md text-center"
          >
            Close
          </button>
        </div>
      </div>
    </div>
  );
}

function ApiToast({
  open,
  onClose,
  message,
}: {
  open: boolean;
  onClose: () => void;
  message: string;
}) {
  if (!open) return null;

  return (
    <div
      className="fixed inset-0 flex items-start justify-center z-50 pointer-events-none pt-10"
      {...(() => {
        setTimeout(onClose, 3000);
        return {};
      })()}
    >
      <div
        id="toast-error"
        className="flex items-center max-w-sm p-4 text-body bg-red-800 rounded-base shadow-xs border border-default animate-fade-in"
        role="alert"
      >
        <div className="inline-flex items-center justify-center shrink-0 w-7 h-7 text-fg-error text-lg bg-error-soft rounded">
          <svg
            className="w-5 h-5"
            aria-hidden="true"
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
          >
            <path
              stroke="currentColor"
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth="2"
              d="M12 9v4m0 4h.01M5.07 19h13.86c1.54 0 2.5-1.67 1.66-3L13.66 4c-.77-1.33-2.55-1.33-3.32 0L3.41 16c-.84 1.33.12 3 1.66 3z"
            />
          </svg>
          <span className="sr-only">Error icon</span>
        </div>

        <div className="ml-3 text-lg font-normal">{message}</div>

        {/* <button
          type="button"
          onClick={onClose}
          className="ms-auto flex items-center justify-center text-body hover:text-heading bg-transparent box-border border border-transparent hover:bg-neutral-secondary-medium focus:ring-4 focus:ring-neutral-tertiary font-medium leading-5 rounded text-sm h-8 w-8 focus:outline-none"
          aria-label="Close"
        >
          <svg
            className="w-5 h-5"
            aria-hidden="true"
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
          >
            <path
              stroke="currentColor"
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth="2"
              d="M6 18 18 6M18 18 6 6"
            />
          </svg>
        </button> */}
      </div>
    </div>
  );
}

export function useDialog() {
  return useContext(DialogContext);
}
