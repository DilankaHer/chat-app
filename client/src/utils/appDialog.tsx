import { faCheckCircle, faExclamationCircle } from '@fortawesome/free-solid-svg-icons';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { createContext, useContext, useState } from 'react';

type DialogContextType = {
  showDialog: (dialogType: string, message: string, isError: boolean, fn?: () => void) => void;
};

const DialogContext = createContext<DialogContextType>({
  showDialog: () => {},
});

export function DialogProvider({ children }: { children: React.ReactNode }) {
  const [dialog, setDialog] = useState({
    open: false,
    dialogType: 'dialog',
    message: '',
    isError: true,
  });
  const [callback, setCallback] = useState<(() => void) | null>(null);
  function showDialog(dialogType: string, message: string, isError: boolean, fn?: () => void) {
    setCallback(() => fn || null);
    setDialog({
      open: true,
      dialogType,
      message,
      isError,
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
          isError={dialog.isError}
        />
      )}
      {dialog.dialogType === 'toast' && (
        <ApiToast
          open={dialog.open}
          onClose={handleClose}
          message={dialog.message}
          isError={dialog.isError}
        />
      )}
    </DialogContext.Provider>
  );
}

function ApiDialog({
  open,
  onClose,
  message,
  isError,
}: {
  open: boolean;
  onClose: () => void;
  message: string;
  isError?: boolean;
}) {
  if (!open) return null;

  return (
    <div className="fixed inset-0 flex items-center justify-center z-50 pointer-events-none">
      <div className="bg-white dark:bg-gray-800 rounded-lg shadow-lg p-6 w-[90%] max-w-sm">
        <div className= {`${isError ? 'text-red-500' : 'text-green-500'} mb-3 text-center text-2xl`}>
          <FontAwesomeIcon
            icon={isError ? faExclamationCircle : faCheckCircle}
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
  isError,
}: {
  open: boolean;
  onClose: () => void;
  message: string;
  isError: boolean;
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
        id="toast"
        className={`${isError ? 'bg-red-800' : 'bg-green-800'} flex items-center max-w-sm p-4 text-body rounded-base shadow-xs border border-default animate-fade-in`}
        role="alert"
      >
        <div className="inline-flex items-center justify-center shrink-0 w-7 h-7 text-fg-error text-lg bg-error-soft rounded">
          <FontAwesomeIcon
            icon={isError ? faExclamationCircle : faCheckCircle}
            fill="currentColor"
            size="xl"
          />
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
